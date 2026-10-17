package warera

import (
	"net/http"
	"time"
)

// Option configures a WarEra Client.
type Option func(*options)

type options struct {
	baseURL            string
	apiKey             string
	httpClient         *http.Client
	rateLimitPerMinute int
	batchWindow        time.Duration
	maxBatchSize       int
	maxRetries         int
}

func defaultOptions() *options {
	return &options{
		baseURL:            "https://api2.warera.io/trpc",
		rateLimitPerMinute: 0, // auto-detect: 200 with apiKey, 100 without
		batchWindow:        25 * time.Millisecond,
		maxBatchSize:       50,
		maxRetries:         3,
	}
}

// WithAPIKey sets the WarEra API key (sent as x-api-key header).
func WithAPIKey(apiKey string) Option {
	return func(o *options) {
		o.apiKey = apiKey
	}
}

// WithBaseURL overrides the default tRPC base URL (e.g. for proxy, gateway, or tests).
func WithBaseURL(url string) Option {
	return func(o *options) {
		o.baseURL = url
	}
}

// WithHTTPClient provides a custom *http.Client for requests.
func WithHTTPClient(client *http.Client) Option {
	return func(o *options) {
		o.httpClient = client
	}
}

// WithRateLimit sets an explicit request rate limit (requests per minute).
// If set to 0, rate limiting is disabled.
func WithRateLimit(requestsPerMinute int) Option {
	return func(o *options) {
		o.rateLimitPerMinute = requestsPerMinute
	}
}

// WithBatchWindow sets the duration to buffer requests before flushing a batch.
// Set to 0 to disable automatic batching and send queries immediately.
func WithBatchWindow(window time.Duration) Option {
	return func(o *options) {
		o.batchWindow = window
	}
}

// WithMaxBatchSize sets the maximum number of queries in a single tRPC batch.
func WithMaxBatchSize(size int) Option {
	return func(o *options) {
		o.maxBatchSize = size
	}
}

// WithMaxRetries sets the maximum number of retry attempts for transient errors.
func WithMaxRetries(retries int) Option {
	return func(o *options) {
		o.maxRetries = retries
	}
}
