package trpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
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
)

// VoidInputProcedures contains endpoints where the server rejects `{}` and requires void/no input.
var VoidInputProcedures = map[string]bool{
	"gameStat.getWorldDevelopment": true,
}

// ClientConfig holds configuration for the low-level tRPC client.
type ClientConfig struct {
	BaseURL             string
	APIKey              string
	HTTPClient          *http.Client
	RateLimitPerMinute  int           // e.g. 100 or 200
	BatchWindow         time.Duration // 0 disables transparent batching
	MaxBatchSize        int           // max items per batch, default 50
	MaxRetries          int
	MaxURLLength        int
	VoidInputProcedures map[string]bool
}

// Client is a low-level tRPC client that handles query encoding, batching, rate limiting, and retries.
type Client struct {
	baseURL             string
	apiKey              string
	httpClient          *http.Client
	limiter             *RateLimiter
	batcher             *Batcher
	maxRetries          int
	maxURLLength        int
	voidInputProcedures map[string]bool
	mu                  sync.RWMutex
}

// BatchCall describes a single query within an explicit batch.
type BatchCall struct {
	Method string
	Input  any
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
	if rateLimit <= 0 {
		if cfg.APIKey != "" {
			rateLimit = 200
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

	voidProcs := make(map[string]bool)
	for k, v := range VoidInputProcedures {
		voidProcs[k] = v
	}
	for k, v := range cfg.VoidInputProcedures {
		voidProcs[k] = v
	}

	c := &Client{
		baseURL:             baseURL,
		apiKey:              cfg.APIKey,
		httpClient:          httpClient,
		limiter:             NewRateLimiter(rateLimit),
		maxRetries:          maxRetries,
		maxURLLength:        maxURLLength,
		voidInputProcedures: voidProcs,
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

// isVoidInput reports whether the given method requires no input parameters.
func (c *Client) isVoidInput(method string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.voidInputProcedures[method]
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
	if !c.isVoidInput(method) {
		if input == nil {
			inputParam = "{}"
		} else {
			data, err := json.Marshal(input)
			if err != nil {
				return fmt.Errorf("failed to marshal input: %w", err)
			}
			inputParam = string(data)
		}
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

	body, err := c.doWithRetry(req)
	if err != nil {
		return err
	}

	var env ResponseEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("failed to decode response envelope: %w (body: %s)", err, string(body))
	}

	if env.Error != nil {
		return env.Error
	}

	if output != nil && len(env.Result.Data) > 0 {
		if err := json.Unmarshal(env.Result.Data, output); err != nil {
			return fmt.Errorf("failed to decode response data: %w", err)
		}
	}

	return nil
}

// CallBatch executes a slice of calls as a single tRPC batch request.
func (c *Client) CallBatch(ctx context.Context, calls []BatchCall) error {
	if len(calls) == 0 {
		return nil
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
		if c.isVoidInput(call.Method) {
			continue
		}
		if call.Input == nil {
			inputsMap[fmt.Sprintf("%d", i)] = struct{}{}
		} else {
			inputsMap[fmt.Sprintf("%d", i)] = call.Input
		}
	}

	methodsJoined := strings.Join(methods, ",")
	inputJSON, err := json.Marshal(inputsMap)
	if err != nil {
		return fmt.Errorf("failed to marshal batch inputs: %w", err)
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
		return fmt.Errorf("failed to create batch request: %w", err)
	}

	c.setHeaders(req)

	body, err := c.doWithRetry(req)
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
		return fmt.Errorf("failed to decode batch response envelopes: %w (body: %s)", err, string(body))
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
		if env.Error != nil {
			calls[i].Err = env.Error
			if firstErr == nil {
				firstErr = env.Error
			}
			continue
		}

		if calls[i].Output != nil && len(env.Result.Data) > 0 {
			if err := json.Unmarshal(env.Result.Data, calls[i].Output); err != nil {
				calls[i].Err = err
				if firstErr == nil {
					firstErr = err
				}
				continue
			}
		}
	}

	return firstErr
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}
}

func (c *Client) doWithRetry(req *http.Request) ([]byte, error) {
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
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			delay := c.getBackoffDelay(attempt)
			select {
			case <-req.Context().Done():
				return nil, req.Context().Err()
			case <-time.After(delay):
			}
		}

		if c.limiter != nil {
			if err := c.limiter.Wait(req.Context()); err != nil {
				return nil, err
			}
		}

		var reqToSend *http.Request
		if len(bodyBytes) > 0 {
			reqToSend = req.Clone(req.Context())
			reqToSend.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		} else {
			reqToSend = req.Clone(req.Context())
		}

		resp, err := c.httpClient.Do(reqToSend)
		if err != nil {
			lastErr = err
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
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
			continue
		}

		// For non-retryable errors (e.g. 400 Bad Request or 401 Unauthorized), check for tRPC error payload
		var env ResponseEnvelope
		if err := json.Unmarshal(respBody, &env); err == nil && env.Error != nil {
			return nil, env.Error
		}

		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			StatusText: resp.Status,
			Body:       respBody,
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", c.maxRetries+1, lastErr)
}

func isRetryableStatusCode(code int) bool {
	switch code {
	case http.StatusRequestTimeout, // 408
		http.StatusConflict,            // 409
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
