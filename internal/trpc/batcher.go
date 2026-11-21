package trpc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type batchEntry struct {
	ctx    context.Context
	method string
	input  any
	output any
	doneCh chan batchResult
}

// batchResult carries undecoded response data so that the waiting caller, not
// the batch worker, writes to the caller's output value.
type batchResult struct {
	data json.RawMessage
	err  error
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
		doneCh: make(chan batchResult, 1),
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
	case res := <-entry.doneCh:
		if res.err != nil {
			return res.err
		}
		if output != nil {
			if err := json.Unmarshal(res.data, output); err != nil {
				return fmt.Errorf("failed to decode response data: %w", err)
			}
		}
		return nil
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
			entry.doneCh <- batchResult{err: err}
		} else {
			activeEntries = append(activeEntries, entry)
		}
	}

	if len(activeEntries) == 0 {
		return
	}

	results := make([]json.RawMessage, len(activeEntries))
	calls := make([]BatchCall, len(activeEntries))
	for i, entry := range activeEntries {
		calls[i] = BatchCall{
			Method: entry.method,
			Input:  entry.input,
			Output: &results[i],
		}
	}

	// One caller's cancellation must not abort the others, so the request is
	// canceled only once every caller has gone away, or when Client.Close
	// cancels the batcher.
	ctx, cancel := context.WithCancel(b.ctx)
	defer cancel()
	var waiting atomic.Int32
	waiting.Store(int32(len(activeEntries)))
	for _, entry := range activeEntries {
		stop := context.AfterFunc(entry.ctx, func() {
			if waiting.Add(-1) == 0 {
				cancel()
			}
		})
		defer stop()
	}

	_ = b.client.CallBatch(ctx, calls)

	for i, entry := range activeEntries {
		entry.doneCh <- batchResult{data: results[i], err: calls[i].Err}
	}
}
