package warera

import (
	"context"

	"github.com/onprem/warera-go/internal/trpc"
)

// Client is the main client for the WarEra API. It exposes typed services and
// shares batching, rate limiting, and retry behavior across their requests.
type Client struct {
	trpc *trpc.Client

	// Countries provides access to country details and listings.
	Countries *CountryService
	// Companies provides access to companies and company listings.
	Companies *CompanyService
	// Battles provides access to battles and live battle data.
	Battles *BattleService
	// Articles provides access to news articles.
	Articles *ArticleService
	// Users provides access to player profiles and country listings.
	Users *UserService
	// Regions provides access to regions and their details.
	Regions *RegionService
	// Government provides access to country government data.
	Government *GovernmentService
	// GameConfig provides access to game configuration and timing data.
	GameConfig *GameConfigService
	// Trading provides access to market prices and top orders.
	Trading *TradingService
	// Search searches across users and other game entities.
	Search *SearchService
	// Events provides access to world and country events.
	Events *EventService
	// BattleOrders provides access to battle side orders.
	BattleOrders *BattleOrderService
	// MilitaryUnits provides access to military unit profiles and listings.
	MilitaryUnits *MilitaryUnitService
	// Inventory provides access to user equipment.
	Inventory *InventoryService
	// Rankings provides access to global game rankings.
	Rankings *RankingService
	// Rounds provides access to battle rounds and hit histories.
	Rounds *RoundService
	// Upgrades provides access to region, company, and military unit upgrades.
	Upgrades *UpgradeService
	// LootSummaries provides access to battle loot summaries.
	LootSummaries *BattleLootSummaryService
	// BattleRankings provides access to rankings within battles and wars.
	BattleRankings *BattleRankingService
	// ItemOffers provides access to marketplace item offers.
	ItemOffers *ItemOfferService
	// MercenaryAuctions provides access to mercenary contract auctions.
	MercenaryAuctions *MercenaryAuctionService
	// Transactions provides access to paginated transaction history.
	Transactions *TransactionService
	// WorkOffers provides access to work offers.
	WorkOffers *WorkOfferService
	// Workers provides access to worker counts and company groupings.
	Workers *WorkerService
}

// NewClient creates a client for the WarEra API. It uses the default API
// endpoint, batches concurrent calls, applies rate limiting, and retries
// transient errors. Use the options to customize these behaviors.
//
// Call Close when the client is no longer needed.
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

// RawCall invokes a tRPC procedure that does not have a typed service method.
// The input is encoded as JSON and output must be a pointer to the destination.
func (c *Client) RawCall(ctx context.Context, method string, input any, output any) error {
	return c.trpc.Call(ctx, method, input, output)
}

// BatchCall describes one procedure in an explicit Client.Batch call. Set
// Output to a pointer where the decoded response should be stored.
type BatchCall = trpc.BatchCall

// Batch executes the supplied procedure calls together. Calls are split into
// requests when the configured batch-size limit is exceeded.
func (c *Client) Batch(ctx context.Context, calls []BatchCall) error {
	return c.trpc.CallBatch(ctx, calls)
}
