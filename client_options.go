package warera

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/onprem/warera-go/internal/trpc"
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
	userAgent          string
	maxRetryWait       time.Duration
	hooks              Hooks
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

// WithRateLimit sets the request rate limit in requests per minute: the client
// never sends more than this many HTTP requests in any 60-second window. A
// positive value sets an explicit limit, 0 disables rate limiting, and a
// negative value enables automatic selection, matching the API's own limits
// (500 with an API key, 100 without one).
//
// Independently of this limit, the client follows the server's rate-limit
// headers: when the server reports its quota exhausted, every request waits
// until the server's window resets. This keeps several clients or processes
// sharing one key from failing.
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

// WithMaxRetryWait caps how long the client waits before one retry when the
// server asks it to slow down with Retry-After or its rate-limit headers. The
// default is one minute, the length of the API's rate-limit window.
func WithMaxRetryWait(d time.Duration) Option {
	return func(o *options) {
		o.maxRetryWait = d
	}
}

// WithUserAgent identifies your tool to the WarEra API, such as
// "market-bot/1.2 (+https://example.com)". The SDK appends its own product
// token, so requests are sent as "market-bot/1.2 (+https://example.com)
// warera-go/v0.3.0".
func WithUserAgent(userAgent string) Option {
	return func(o *options) {
		o.userAgent = userAgent
	}
}

// Hooks observe the client's HTTP traffic, for logging or metrics. Hooks run
// synchronously on the requesting goroutine and must be safe for concurrent
// use. Set them with [WithHooks].
type Hooks = trpc.Hooks

// RequestInfo describes one HTTP attempt reported to [Hooks.OnRequest].
type RequestInfo = trpc.RequestInfo

// WaitInfo describes a pause reported to [Hooks.OnWait].
type WaitInfo = trpc.WaitInfo

// WithHooks sets hooks that observe every HTTP attempt and every pause for a
// retry or a rate limit. For example, to log retries with log/slog:
//
//	warera.WithHooks(warera.Hooks{
//		OnWait: func(w warera.WaitInfo) {
//			slog.Info("warera wait", "reason", w.Reason, "duration", w.Duration)
//		},
//	})
func WithHooks(hooks Hooks) Option {
	return func(o *options) {
		o.hooks = hooks
	}
}

// userAgent returns the User-Agent header: the caller's product, if any,
// followed by this module's name and version.
func userAgent(product string) string {
	sdk := "warera-go"
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Path == "github.com/onprem/warera-go" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			sdk += "/" + info.Main.Version
		}
		for _, dep := range info.Deps {
			if dep.Path == "github.com/onprem/warera-go" {
				sdk += "/" + dep.Version
			}
		}
	}
	if product == "" {
		return sdk
	}
	return product + " " + sdk
}
