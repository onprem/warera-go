package trpc

import (
	"context"
	"sync"
	"time"
)

type batchEntry struct {
	ctx    context.Context
	method string
	input  any
	output any
	doneCh chan error
}

// Batcher coalesces concurrent procedure calls into single HTTP tRPC batches.
type Batcher struct {
	client     *Client
	ctx        context.Context
	cancel     context.CancelFunc
	window     time.Duration
	maxSize    int
	inCh       chan *batchEntry
	closeCh    chan struct{}
	closeOnce  sync.Once
	stateMu    sync.Mutex
	closed     bool
	loopDone   chan struct{}
	dispatchWG sync.WaitGroup
}

func newBatcher(client *Client, window time.Duration, maxSize int) *Batcher {
	ctx, cancel := context.WithCancel(context.Background())
	b := &Batcher{
		client:   client,
		ctx:      ctx,
		cancel:   cancel,
		window:   window,
		maxSize:  maxSize,
		inCh:     make(chan *batchEntry, 256),
		closeCh:  make(chan struct{}),
		loopDone: make(chan struct{}),
	}
	go b.loop()
	return b
}

func (b *Batcher) close() {
	b.closeOnce.Do(func() {
		b.stateMu.Lock()
		b.closed = true
		close(b.closeCh)
		b.stateMu.Unlock()
		<-b.loopDone
		b.cancel()
		b.dispatchWG.Wait()
	})
}

func (b *Batcher) submit(ctx context.Context, method string, input any, output any) error {
	entry := &batchEntry{
		ctx:    ctx,
		method: method,
		input:  input,
		output: output,
		doneCh: make(chan error, 1),
	}

	b.stateMu.Lock()
	if b.closed {
		b.stateMu.Unlock()
		return b.client.CallDirect(ctx, method, input, output)
	}
	select {
	case <-ctx.Done():
		b.stateMu.Unlock()
		return ctx.Err()
	case b.inCh <- entry:
		b.stateMu.Unlock()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-entry.doneCh:
		return err
	}
}

func (b *Batcher) loop() {
	defer close(b.loopDone)
	var (
		pending []*batchEntry
		timer   *time.Timer
		timerCh <-chan time.Time
	)

	flush := func() {
		if len(pending) == 0 {
			return
		}
		toFlush := pending
		pending = nil
		if timer != nil {
			timer.Stop()
			timerCh = nil
		}
		// Run batch execution in a separate goroutine so queue collection isn't blocked by network I/O
		b.dispatchWG.Add(1)
		go func() {
			defer b.dispatchWG.Done()
			b.dispatchBatch(toFlush)
		}()
	}

	for {
		select {
		case <-b.closeCh:
			// close marks the batcher closed while holding stateMu, so every
			// submission accepted before shutdown is already queued here.
			for {
				select {
				case entry := <-b.inCh:
					pending = append(pending, entry)
					if len(pending) >= b.maxSize {
						flush()
					}
				default:
					flush()
					return
				}
			}

		case <-timerCh:
			flush()

		case entry := <-b.inCh:
			pending = append(pending, entry)
			if len(pending) == 1 {
				timer = time.NewTimer(b.window)
				timerCh = timer.C
			} else if len(pending) >= b.maxSize {
				flush()
			}
		}
	}
}

func (b *Batcher) dispatchBatch(entries []*batchEntry) {
	// Filter out already-canceled entries
	activeEntries := make([]*batchEntry, 0, len(entries))
	for _, entry := range entries {
		if err := entry.ctx.Err(); err != nil {
			entry.doneCh <- err
		} else {
			activeEntries = append(activeEntries, entry)
		}
	}

	if len(activeEntries) == 0 {
		return
	}

	calls := make([]BatchCall, len(activeEntries))
	for i, entry := range activeEntries {
		calls[i] = BatchCall{
			Method: entry.method,
			Input:  entry.input,
			Output: entry.output,
		}
	}

	// Use the batcher's shared context so one caller's cancellation does not
	// abort other calls, while Client.Close can still cancel outstanding I/O.
	_ = b.client.CallBatch(b.ctx, calls)

	for i, entry := range activeEntries {
		entry.doneCh <- calls[i].Err
	}
}
