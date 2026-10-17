package trpc

import (
	"context"
	"sync"
	"time"
)

// RateLimiter is a thread-safe token bucket rate limiter with zero external dependencies.
type RateLimiter struct {
	mu       sync.Mutex
	rate     float64 // tokens added per second
	burst    float64 // max burst capacity
	tokens   float64 // current available tokens
	lastTime time.Time
}

// NewRateLimiter creates a new rate limiter allowing requestsPerMinute requests.
// If requestsPerMinute <= 0, rate limiting is disabled (returns nil).
func NewRateLimiter(requestsPerMinute int) *RateLimiter {
	if requestsPerMinute <= 0 {
		return nil
	}
	r := float64(requestsPerMinute) / 60.0
	burst := float64(requestsPerMinute)
	return &RateLimiter{
		rate:     r,
		burst:    burst,
		tokens:   burst,
		lastTime: time.Now(),
	}
}

// Wait blocks until a token is available or the context is cancelled.
func (l *RateLimiter) Wait(ctx context.Context) error {
	if l == nil {
		return nil
	}

	for {
		l.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(l.lastTime).Seconds()
		l.lastTime = now

		l.tokens += elapsed * l.rate
		if l.tokens > l.burst {
			l.tokens = l.burst
		}

		if l.tokens >= 1.0 {
			l.tokens -= 1.0
			l.mu.Unlock()
			return nil
		}

		needed := 1.0 - l.tokens
		waitDuration := time.Duration(needed / l.rate * float64(time.Second))
		l.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitDuration):
		}
	}
}
