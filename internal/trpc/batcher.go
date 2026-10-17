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
	client   *Client
	window   time.Duration
	maxSize  int
	inCh     chan *batchEntry
	closeCh  chan struct{}
	closeOnce sync.Once
}

func newBatcher(client *Client, window time.Duration, maxSize int) *Batcher {
	b := &Batcher{
		client:  client,
		window:  window,
		maxSize: maxSize,
		inCh:    make(chan *batchEntry, 256),
		closeCh: make(chan struct{}),
	}
	go b.loop()
	return b
}

func (b *Batcher) close() {
	b.closeOnce.Do(func() {
		close(b.closeCh)
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

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.closeCh:
		// Fallback to direct call if batcher has closed
		return b.client.CallDirect(ctx, method, input, output)
	case b.inCh <- entry:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-entry.doneCh:
		return err
	}
}

func (b *Batcher) loop() {
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
		go b.dispatchBatch(toFlush)
	}

	for {
		select {
		case <-b.closeCh:
			flush()
			return

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

	// Use background context for the batch call so one caller's cancel doesn't abort the entire batch
	_ = b.client.CallBatch(context.Background(), calls)

	for i, entry := range activeEntries {
		entry.doneCh <- calls[i].Err
	}
}
