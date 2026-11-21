package trpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DefaultBaseURL      = "https://api2.warera.io/trpc"
	DefaultMaxBatchSize = 50
	DefaultBatchWindow  = 25 * time.Millisecond
	DefaultMaxRetries   = 3
	DefaultInitialDelay = 250 * time.Millisecond
	DefaultMaxDelay     = 5 * time.Second
	DefaultMaxURLLength = 2048
	DefaultMaxRetryWait = time.Minute
)

// ClientConfig holds configuration for the low-level tRPC client.
type ClientConfig struct {
	BaseURL            string
	APIKey             string
	HTTPClient         *http.Client
	RateLimitPerMinute int           // negative auto-detects, zero disables, positive sets a limit
	BatchWindow        time.Duration // 0 disables transparent batching
	MaxBatchSize       int           // max items per batch, default 50
	MaxRetries         int
	MaxURLLength       int
	// UserAgent is sent with every request.
	UserAgent string
	// MaxRetryWait caps how long one retry waits for the server's rate limit
	// to reset; 0 selects DefaultMaxRetryWait.
	MaxRetryWait time.Duration
	Hooks        Hooks
}

// Hooks observe the client's HTTP traffic. Hooks run synchronously on the
// requesting goroutine and must be safe for concurrent use.
type Hooks struct {
	// OnRequest is called after every HTTP attempt, including retries.
	OnRequest func(RequestInfo)
	// OnWait is called before the client sleeps, either to back off before a
	// retry or because the server's rate limit is exhausted.
	OnWait func(WaitInfo)
}

// RequestInfo describes one HTTP attempt.
type RequestInfo struct {
	// Procedures lists the procedures in the request; more than one for a
	// batch.
	Procedures []string
	// Attempt is 0 for the first try and counts retries after it.
	Attempt int
	// StatusCode is the HTTP status, or 0 when no response arrived.
	StatusCode int
	Duration   time.Duration
	// Err is the transport error when no response arrived.
	Err error
	// RateLimitRemaining is the server's remaining request quota, or -1 when
	// the response did not report it.
	RateLimitRemaining int
}

// WaitInfo describes a pause before a request.
type WaitInfo struct {
	// Reason is "retry" for a backoff before retrying, "server-rate-limit" when
	// the server reported its quota exhausted, or "rate-limit" when the
	// client's own limiter delayed the request.
	Reason     string
	Procedures []string
	Duration   time.Duration
}

// Client is a low-level tRPC client that handles query encoding, batching, rate limiting, and retries.
type Client struct {
	baseURL      string
	apiKey       string
	httpClient   *http.Client
	limiter      *RateLimiter
	batcher      *Batcher
	maxRetries   int
	maxBatchSize int
	maxURLLength int
	userAgent    string
	maxRetryWait time.Duration
	hooks        Hooks

	pauseMu    sync.Mutex
	pauseUntil time.Time // set when the server reports its quota exhausted
}

// BatchCall describes a single query within an explicit batch.
type BatchCall struct {
	Method string
	Input  any // nil sends no input
	Output any
	Err    error
}

// NewClient creates a new tRPC client.
func NewClient(cfg ClientConfig) *Client {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	rateLimit := cfg.RateLimitPerMinute
	if rateLimit < 0 {
		if cfg.APIKey != "" {
			rateLimit = 500
		} else {
			rateLimit = 100
		}
	}

	maxBatchSize := cfg.MaxBatchSize
	if maxBatchSize <= 0 || maxBatchSize > DefaultMaxBatchSize {
		maxBatchSize = DefaultMaxBatchSize
	}

	maxRetries := cfg.MaxRetries
	if maxRetries < 0 {
		maxRetries = DefaultMaxRetries
	}

	maxURLLength := cfg.MaxURLLength
	if maxURLLength <= 0 {
		maxURLLength = DefaultMaxURLLength
	}

	c := &Client{
		baseURL:      baseURL,
		apiKey:       cfg.APIKey,
		httpClient:   httpClient,
		limiter:      NewRateLimiter(rateLimit),
		maxRetries:   maxRetries,
		maxBatchSize: maxBatchSize,
		maxURLLength: maxURLLength,
		userAgent:    cfg.UserAgent,
		maxRetryWait: cfg.MaxRetryWait,
		hooks:        cfg.Hooks,
	}
	if c.maxRetryWait <= 0 {
		c.maxRetryWait = DefaultMaxRetryWait
	}

	if cfg.BatchWindow > 0 {
		c.batcher = newBatcher(c, cfg.BatchWindow, maxBatchSize)
	}

	return c
}

// Close shuts down background batching workers if any.
func (c *Client) Close() {
	if c.batcher != nil {
		c.batcher.close()
	}
}

// Call performs a tRPC procedure call. If transparent batching is enabled,
// it joins the current batch window. Otherwise, it executes directly.
func (c *Client) Call(ctx context.Context, method string, input any, output any) error {
	if c.batcher != nil {
		return c.batcher.submit(ctx, method, input, output)
	}
	return c.CallDirect(ctx, method, input, output)
}

// CallDirect performs an immediate single unbatched call to the tRPC server.
func (c *Client) CallDirect(ctx context.Context, method string, input any, output any) error {
	var inputParam string
	if input != nil {
		data, err := json.Marshal(input)
		if err != nil {
			return fmt.Errorf("failed to marshal input: %w", err)
		}
		inputParam = string(data)
	}

	endpointURL := fmt.Sprintf("%s/%s", c.baseURL, method)
	if inputParam != "" {
		endpointURL = fmt.Sprintf("%s?input=%s", endpointURL, url.QueryEscape(inputParam))
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpointURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(req)

	body, err := c.doWithRetry(req, []string{method})
	if err != nil {
		return err
	}

	var env ResponseEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("failed to decode response envelope: %w (body: %s)", err, string(body))
	}

	return decodeResult(env, output)
}

// CallBatch executes a slice of calls as a single tRPC batch request.
func (c *Client) CallBatch(ctx context.Context, calls []BatchCall) error {
	if len(calls) == 0 {
		return nil
	}
	if len(calls) > c.maxBatchSize {
		var firstErr error
		for start := 0; start < len(calls); start += c.maxBatchSize {
			end := start + c.maxBatchSize
			if end > len(calls) {
				end = len(calls)
			}
			if err := c.CallBatch(ctx, calls[start:end]); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		return firstErr
	}

	if len(calls) == 1 {
		err := c.CallDirect(ctx, calls[0].Method, calls[0].Input, calls[0].Output)
		calls[0].Err = err
		return err
	}

	methods := make([]string, len(calls))
	inputsMap := make(map[string]any, len(calls))

	for i, call := range calls {
		methods[i] = call.Method
		if call.Input != nil {
			inputsMap[fmt.Sprintf("%d", i)] = call.Input
		}
	}

	methodsJoined := strings.Join(methods, ",")
	inputJSON, err := json.Marshal(inputsMap)
	if err != nil {
		return setBatchError(calls, fmt.Errorf("failed to marshal batch inputs: %w", err))
	}

	getURL := fmt.Sprintf("%s/%s?batch=1&input=%s", c.baseURL, methodsJoined, url.QueryEscape(string(inputJSON)))

	var req *http.Request
	if len(getURL) <= c.maxURLLength {
		req, err = http.NewRequestWithContext(ctx, "GET", getURL, nil)
	} else {
		postURL := fmt.Sprintf("%s/%s?batch=1", c.baseURL, methodsJoined)
		req, err = http.NewRequestWithContext(ctx, "POST", postURL, bytes.NewReader(inputJSON))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	if err != nil {
		return setBatchError(calls, fmt.Errorf("failed to create batch request: %w", err))
	}

	c.setHeaders(req)

	body, err := c.doWithRetry(req, methods)
	if err != nil {
		for i := range calls {
			calls[i].Err = err
		}
		return err
	}

	var envelopes []ResponseEnvelope
	if err := json.Unmarshal(body, &envelopes); err != nil {
		var single ResponseEnvelope
		if errSingle := json.Unmarshal(body, &single); errSingle == nil && single.Error != nil {
			for i := range calls {
				calls[i].Err = single.Error
			}
			return single.Error
		}
		return setBatchError(calls, fmt.Errorf("failed to decode batch response envelopes: %w (body: %s)", err, string(body)))
	}

	if len(envelopes) != len(calls) {
		errMismatch := fmt.Errorf("mismatched batch response count: expected %d, got %d", len(calls), len(envelopes))
		for i := range calls {
			calls[i].Err = errMismatch
		}
		return errMismatch
	}

	var firstErr error
	for i, env := range envelopes {
		if err := decodeResult(env, calls[i].Output); err != nil {
			calls[i].Err = err
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

func setBatchError(calls []BatchCall, err error) error {
	for i := range calls {
		calls[i].Err = err
	}
	return err
}

// ErrNotFound is returned when a procedure's result is null but the caller's
// output cannot represent null, such as a struct or a number. Outputs that can
// (pointers, slices, maps, and interfaces) decode null as their zero value.
var ErrNotFound = errors.New("warera: not found")

func decodeResult(env ResponseEnvelope, output any) error {
	if env.Error != nil {
		return env.Error
	}
	if env.Result == nil {
		return fmt.Errorf("invalid tRPC response: missing result and error")
	}
	if len(env.Result.Data) == 0 {
		return fmt.Errorf("invalid tRPC response: missing result data")
	}
	return decodeData(env.Result.Data, output)
}

// decodeData decodes a procedure's result into output, returning ErrNotFound
// for a null result that output cannot represent.
func decodeData(data json.RawMessage, output any) error {
	if output == nil {
		return nil
	}
	if string(bytes.TrimSpace(data)) == "null" && !acceptsNull(output) {
		return ErrNotFound
	}
	if err := json.Unmarshal(data, output); err != nil {
		return fmt.Errorf("failed to decode response data: %w", err)
	}
	return nil
}

// acceptsNull reports whether decoding JSON null into output yields a
// meaningful "none" value rather than silently leaving it unchanged.
func acceptsNull(output any) bool {
	v := reflect.ValueOf(output)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return true
	}
	switch v.Elem().Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Interface:
		return true
	}
	return false
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}
}

func (c *Client) doWithRetry(req *http.Request, procedures []string) ([]byte, error) {
	ctx := req.Context()
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		req.Body.Close()
	}

	var lastErr error
	var serverWait time.Duration // how long the last response asked us to wait
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			delay := max(c.getBackoffDelay(attempt), min(serverWait, c.maxRetryWait))
			if err := c.sleep(ctx, WaitInfo{Reason: "retry", Procedures: procedures, Duration: delay}); err != nil {
				return nil, err
			}
		}

		if err := c.waitForServerQuota(ctx, procedures); err != nil {
			return nil, err
		}
		if c.limiter != nil {
			started := time.Now()
			if err := c.limiter.Wait(ctx); err != nil {
				return nil, err
			}
			if waited := time.Since(started); waited >= time.Millisecond && c.hooks.OnWait != nil {
				c.hooks.OnWait(WaitInfo{Reason: "rate-limit", Procedures: procedures, Duration: waited})
			}
		}

		reqToSend := req.Clone(ctx)
		if len(bodyBytes) > 0 {
			reqToSend.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		started := time.Now()
		resp, err := c.httpClient.Do(reqToSend)
		if err != nil {
			c.reportRequest(procedures, attempt, 0, started, err, -1)
			lastErr = err
			serverWait = 0
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		remaining, reset := c.observeRateLimit(resp.Header)
		c.reportRequest(procedures, attempt, resp.StatusCode, started, err, remaining)
		if err != nil {
			lastErr = err
			serverWait = 0
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return respBody, nil
		}

		if isRetryableStatusCode(resp.StatusCode) {
			lastErr = &HTTPError{
				StatusCode: resp.StatusCode,
				StatusText: resp.Status,
				Body:       respBody,
			}
			serverWait = retryAfter(resp.Header)
			if resp.StatusCode == http.StatusTooManyRequests && reset > serverWait {
				serverWait = reset
			}
			continue
		}

		// For non-retryable errors (e.g. 400 Bad Request or 401 Unauthorized), check for tRPC error payload
		var env ResponseEnvelope
		if err := json.Unmarshal(respBody, &env); err == nil && env.Error != nil {
			return nil, env.Error
		}
		// When every call in a batch fails with the same status, tRPC responds
		// with that status and one error envelope per call. Return the body so
		// each caller receives its own tRPC error.
		var envs []ResponseEnvelope
		if err := json.Unmarshal(respBody, &envs); err == nil && len(envs) > 0 {
			return respBody, nil
		}

		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			StatusText: resp.Status,
			Body:       respBody,
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.maxRetries+1, lastErr)
}

func (c *Client) sleep(ctx context.Context, info WaitInfo) error {
	if c.hooks.OnWait != nil {
		c.hooks.OnWait(info)
	}
	timer := time.NewTimer(info.Duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (c *Client) reportRequest(procedures []string, attempt, status int, started time.Time, err error, remaining int) {
	if c.hooks.OnRequest == nil {
		return
	}
	c.hooks.OnRequest(RequestInfo{
		Procedures:         procedures,
		Attempt:            attempt,
		StatusCode:         status,
		Duration:           time.Since(started),
		Err:                err,
		RateLimitRemaining: remaining,
	})
}

// observeRateLimit reads the server's ratelimit-remaining and ratelimit-reset
// headers. When the quota is exhausted, it pauses all requests from this client
// until the server's window resets, which also keeps several processes sharing
// one API key from failing. It returns the remaining quota (-1 when absent)
// and the time until reset.
func (c *Client) observeRateLimit(h http.Header) (int, time.Duration) {
	remaining, err := strconv.Atoi(h.Get("Ratelimit-Remaining"))
	if err != nil {
		remaining = -1
	}
	var reset time.Duration
	if seconds, err := strconv.Atoi(h.Get("Ratelimit-Reset")); err == nil && seconds > 0 {
		reset = time.Duration(seconds) * time.Second
	}
	if remaining == 0 && reset > 0 {
		until := time.Now().Add(min(reset, c.maxRetryWait))
		c.pauseMu.Lock()
		if until.After(c.pauseUntil) {
			c.pauseUntil = until
		}
		c.pauseMu.Unlock()
	}
	return remaining, reset
}

func (c *Client) waitForServerQuota(ctx context.Context, procedures []string) error {
	c.pauseMu.Lock()
	wait := time.Until(c.pauseUntil)
	c.pauseMu.Unlock()
	if wait <= 0 {
		return nil
	}
	return c.sleep(ctx, WaitInfo{Reason: "server-rate-limit", Procedures: procedures, Duration: wait})
}

// retryAfter parses a Retry-After header given in seconds or as an HTTP date.
func retryAfter(h http.Header) time.Duration {
	value := h.Get("Retry-After")
	if value == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(value); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	if at, err := http.ParseTime(value); err == nil {
		return time.Until(at)
	}
	return 0
}

func isRetryableStatusCode(code int) bool {
	switch code {
	case http.StatusRequestTimeout, // 408
		http.StatusTooEarly,            // 425
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	default:
		return false
	}
}

func (c *Client) getBackoffDelay(attempt int) time.Duration {
	shift := 0
	if attempt > 1 {
		shift = attempt - 1
		if shift > 10 {
			shift = 10
		}
	}
	factor := 1
	for i := 0; i < shift; i++ {
		factor *= 2
	}
	base := float64(DefaultInitialDelay) * float64(factor)
	max := float64(DefaultMaxDelay)
	if base > max {
		base = max
	}
	// Add 50% - 150% jitter
	jitter := (0.5 + rand.Float64()) * base
	return time.Duration(jitter)
}
