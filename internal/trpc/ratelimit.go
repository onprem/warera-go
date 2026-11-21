package trpc

import (
	"context"
	"sync"
	"time"
)

// RateLimiter is a thread-safe sliding-window rate limiter with zero external
// dependencies. It never admits more than its limit within any one window.
type RateLimiter struct {
	mu     sync.Mutex
	window time.Duration
	// sent holds the times of the most recent requests as a ring buffer; next
	// indexes the oldest. Zero times stand for requests that never happened.
	sent []time.Time
	next int
}

// NewRateLimiter creates a new rate limiter allowing requestsPerMinute requests
// in any 60-second window. If requestsPerMinute <= 0, rate limiting is disabled
// (returns nil).
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	return newRateLimiter(requestsPerMinute, time.Minute)
}

func newRateLimiter(limit int, window time.Duration) *RateLimiter {
	if limit <= 0 {
		return nil
	}
	return &RateLimiter{window: window, sent: make([]time.Time, limit)}
}

// Wait blocks until a request is allowed or the context is cancelled.
func (l *RateLimiter) Wait(ctx context.Context) error {
	if l == nil {
		return nil
	}

	for {
		l.mu.Lock()
		now := time.Now()
		wait := l.sent[l.next].Add(l.window).Sub(now)
		if wait <= 0 {
			l.sent[l.next] = now
			l.next = (l.next + 1) % len(l.sent)
			l.mu.Unlock()
			return nil
		}
		l.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
}
