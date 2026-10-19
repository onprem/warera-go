package warera

import (
	"net/http"
	"time"
)

// Option configures a WarEra Client.
type Option func(*options)

const autoRateLimit = -1

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
		rateLimitPerMinute: autoRateLimit,
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

// WithBaseURL overrides the default tRPC endpoint, https://api2.warera.io/trpc.
// This is useful for a proxy, gateway, or test server.
func WithBaseURL(url string) Option {
	return func(o *options) {
		o.baseURL = url
	}
}

// WithHTTPClient provides the HTTP client used for requests. When unset, the
// SDK uses an http.Client with a 30-second timeout.
func WithHTTPClient(client *http.Client) Option {
	return func(o *options) {
		o.httpClient = client
	}
}

// WithRateLimit sets the request rate limit in requests per minute. A positive
// value sets an explicit limit, 0 disables rate limiting, and a negative value
// enables automatic selection (200 with an API key, 100 without one).
func WithRateLimit(requestsPerMinute int) Option {
	return func(o *options) {
		o.rateLimitPerMinute = requestsPerMinute
	}
}

// WithBatchWindow sets how long concurrent calls wait to be grouped into a
// batch. The default is 25 milliseconds. Set it to 0 to disable automatic
// batching and send calls immediately.
func WithBatchWindow(window time.Duration) Option {
	return func(o *options) {
		o.batchWindow = window
	}
}

// WithMaxBatchSize sets the maximum number of queries in one tRPC batch. The
// default and server-supported maximum is 50.
func WithMaxBatchSize(size int) Option {
	return func(o *options) {
		o.maxBatchSize = size
	}
}

// WithMaxRetries sets the maximum number of retries for transient errors. The
// default is 3; set it to 0 to disable retries.
func WithMaxRetries(retries int) Option {
	return func(o *options) {
		o.maxRetries = retries
	}
}
