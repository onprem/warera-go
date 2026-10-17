package warera

import (
	"context"

	"github.com/onprem/warera-go/internal/trpc"
)

// Client is the main client for interacting with the WarEra API.
type Client struct {
	trpc *trpc.Client

	// Services
	Countries  *CountryService
	Companies  *CompanyService
	Battles    *BattleService
	Articles   *ArticleService
	Users      *UserService
	Regions    *RegionService
	Government *GovernmentService
	GameConfig *GameConfigService
	Trading    *TradingService
	Search     *SearchService
}

// NewClient creates a new WarEra API client with the given options.
func NewClient(opts ...Option) *Client {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	trpcClient := trpc.NewClient(trpc.ClientConfig{
		BaseURL:            o.baseURL,
		APIKey:             o.apiKey,
		HTTPClient:         o.httpClient,
		RateLimitPerMinute: o.rateLimitPerMinute,
		BatchWindow:        o.batchWindow,
		MaxBatchSize:       o.maxBatchSize,
		MaxRetries:         o.maxRetries,
	})

	c := &Client{
		trpc: trpcClient,
	}

	c.Countries = &CountryService{client: c}
	c.Companies = &CompanyService{client: c}
	c.Battles = &BattleService{client: c}
	c.Articles = &ArticleService{client: c}
	c.Users = &UserService{client: c}
	c.Regions = &RegionService{client: c}
	c.Government = &GovernmentService{client: c}
	c.GameConfig = &GameConfigService{client: c}
	c.Trading = &TradingService{client: c}
	c.Search = &SearchService{client: c}

	return c
}

// Close gracefully releases any background workers (e.g. batch buffers).
func (c *Client) Close() {
	if c.trpc != nil {
		c.trpc.Close()
	}
}

// RawCall allows invoking an arbitrary or unmapped tRPC procedure directly.
func (c *Client) RawCall(ctx context.Context, method string, input any, output any) error {
	return c.trpc.Call(ctx, method, input, output)
}

// BatchCall represents a single procedure call within an explicit Batch execution.
type BatchCall = trpc.BatchCall

// Batch executes multiple distinct procedure calls in a single batched HTTP request.
func (c *Client) Batch(ctx context.Context, calls []BatchCall) error {
	return c.trpc.CallBatch(ctx, calls)
}
