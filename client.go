package warera

import (
	"context"

	"github.com/onprem/warera-go/internal/trpc"
)

// Client is the main client for interacting with the WarEra API.
type Client struct {
	trpc *trpc.Client

	// Services
	Countries         *CountryService
	Companies         *CompanyService
	Battles           *BattleService
	Articles          *ArticleService
	Users             *UserService
	Regions           *RegionService
	Government        *GovernmentService
	GameConfig        *GameConfigService
	Trading           *TradingService
	Search            *SearchService
	Events            *EventService
	BattleOrders      *BattleOrderService
	MilitaryUnits     *MilitaryUnitService
	Inventory         *InventoryService
	Rankings          *RankingService
	Rounds            *RoundService
	Upgrades          *UpgradeService
	LootSummaries     *BattleLootSummaryService
	BattleRankings    *BattleRankingService
	ItemOffers        *ItemOfferService
	MercenaryAuctions *MercenaryAuctionService
	Transactions      *TransactionService
	WorkOffers        *WorkOfferService
	Workers           *WorkerService
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
	c.Events = &EventService{client: c}
	c.BattleOrders = &BattleOrderService{client: c}
	c.MilitaryUnits = &MilitaryUnitService{client: c}
	c.Inventory = &InventoryService{client: c}
	c.Rankings = &RankingService{client: c}
	c.Rounds = &RoundService{client: c}
	c.Upgrades = &UpgradeService{client: c}
	c.LootSummaries = &BattleLootSummaryService{client: c}
	c.BattleRankings = &BattleRankingService{client: c}
	c.ItemOffers = &ItemOfferService{client: c}
	c.MercenaryAuctions = &MercenaryAuctionService{client: c}
	c.Transactions = &TransactionService{client: c}
	c.WorkOffers = &WorkOfferService{client: c}
	c.Workers = &WorkerService{client: c}

	return c
}

// Close flushes queued calls, cancels any batches still in flight, and waits for batch workers to stop.
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
