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
	// Employees provides access to the NPC employees of companies.
	Employees *EmployeeService
	// Managers provides access to company managers and manager job offers.
	Managers *ManagerService
	// Battles provides access to battles and live battle data.
	Battles *BattleService
	// Articles provides access to news articles.
	Articles *ArticleService
	// Polls provides access to article polls and their voters.
	Polls *PollService
	// Users provides access to player profiles and country listings.
	Users *UserService
	// Regions provides access to regions and their details.
	Regions *RegionService
	// Government provides access to country government data.
	Government *GovernmentService
	// Alliances provides access to alliances and their strategic resource bonuses.
	Alliances *AllianceService
	// Wars provides access to wars between countries.
	Wars *WarService
	// Parties provides access to political parties.
	Parties *PartyService
	// Elections provides access to country and party elections.
	Elections *ElectionService
	// GameConfig provides access to game configuration and timing data.
	GameConfig *GameConfigService
	// Trading provides access to market prices, order books, and trading history.
	Trading *TradingService
	// MarketStats provides access to game-wide economy statistics such as production point values.
	MarketStats *MarketStatsService
	// Donations provides access to money donations to countries, military units, and parties.
	Donations *DonationService
	// Search searches across users and other game entities.
	Search *SearchService
	// Events provides access to world and country events.
	Events *EventService
	// BattleOrders provides access to battle side orders.
	BattleOrders *BattleOrderService
	// MilitaryUnits provides access to military unit profiles and listings.
	MilitaryUnits *MilitaryUnitService
	// Tournaments provides access to MU and country tournaments and their teams.
	Tournaments *TournamentService
	// StateArmy provides access to countries' state army spending.
	StateArmy *StateArmyService
	// Contributions provides access to contributions to country unrest and region resistance.
	Contributions *ContributionService
	// Leagues provides access to the weekly player leagues.
	Leagues *LeagueService
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
	// MercenaryAuctions provides access to mercenary contract auctions.
	MercenaryAuctions *MercenaryAuctionService
	// Transactions provides access to paginated transaction history.
	Transactions *TransactionService
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
		UserAgent:          userAgent(o.userAgent),
		MaxRetryWait:       o.maxRetryWait,
		Hooks:              o.hooks,
	})

	c := &Client{
		trpc: trpcClient,
	}

	c.Countries = &CountryService{client: c}
	c.Companies = &CompanyService{client: c}
	c.Employees = &EmployeeService{client: c}
	c.Managers = &ManagerService{client: c}
	c.Battles = &BattleService{client: c}
	c.Articles = &ArticleService{client: c}
	c.Polls = &PollService{client: c}
	c.Users = &UserService{client: c}
	c.Regions = &RegionService{client: c}
	c.Government = &GovernmentService{client: c}
	c.Alliances = &AllianceService{client: c}
	c.Wars = &WarService{client: c}
	c.Parties = &PartyService{client: c}
	c.Elections = &ElectionService{client: c}
	c.GameConfig = &GameConfigService{client: c}
	c.Trading = &TradingService{client: c}
	c.MarketStats = &MarketStatsService{client: c}
	c.Donations = &DonationService{client: c}
	c.Search = &SearchService{client: c}
	c.Events = &EventService{client: c}
	c.BattleOrders = &BattleOrderService{client: c}
	c.MilitaryUnits = &MilitaryUnitService{client: c}
	c.Tournaments = &TournamentService{client: c}
	c.StateArmy = &StateArmyService{client: c}
	c.Contributions = &ContributionService{client: c}
	c.Leagues = &LeagueService{client: c}
	c.Inventory = &InventoryService{client: c}
	c.Rankings = &RankingService{client: c}
	c.Rounds = &RoundService{client: c}
	c.Upgrades = &UpgradeService{client: c}
	c.LootSummaries = &BattleLootSummaryService{client: c}
	c.BattleRankings = &BattleRankingService{client: c}
	c.MercenaryAuctions = &MercenaryAuctionService{client: c}
	c.Transactions = &TransactionService{client: c}

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
