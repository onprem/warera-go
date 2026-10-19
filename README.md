# warera-go

[![Go Reference](https://pkg.go.dev/badge/github.com/onprem/warera-go.svg)](https://pkg.go.dev/github.com/onprem/warera-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

An idiomatic, high-performance Go SDK for the [WarEra.io](https://warera.io) online multiplayer strategy game API.

## Features

- **Transparent tRPC HTTP Batching:** Multiple procedure calls from concurrent goroutines are automatically coalesced into a single HTTP request within a configurable time window (default `25ms`), dramatically saving network bandwidth and adhering to API limits.
- **Zero External Dependencies:** Built entirely with the Go standard library (`net/http`, `sync`, `encoding/json`).
- **Go 1.23+ Iterators:** Native range-over-func support (`iter.Seq2[T, error]`) for effortless cursor-based pagination.
- **Smart Rate Limiting:** Built-in token-bucket rate limiter (200 req/min with API key, 100 req/min unauthenticated).
- **Transient Error Retries:** Automatic retries with exponential backoff and jitter for network drops and 429/5xx status codes.
- **URL Length Optimization:** Seamlessly switches large batch GET requests into POST to avoid URL length cutoffs.
- **Typed Response Models:** Go data structures for Countries, Companies, Battles, Users, Articles, Market Trading, and Game Clock, with nullable scalar pointers and opaque JSON fields where the API shape varies or is undocumented.

---

## Installation

```bash
go get github.com/onprem/warera-go
```

*Requires Go 1.23 or newer.*

---

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/onprem/warera-go"
)

func main() {
	// Initialize client with an optional API key.
	client := warera.NewClient(
		warera.WithAPIKey("your-api-key"),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Fetch official game rollover and energy regen times
	dates, err := client.GameConfig.GetDates(ctx)
	if err != nil {
		log.Fatalf("failed to fetch game dates: %v", err)
	}
	fmt.Printf("Next Day: %s, Next Energy Regen: %s\n", dates.NextDayAt, dates.NextRegenAt)

	// 2. Fetch all countries
	countries, err := client.Countries.GetAll(ctx)
	if err != nil {
		log.Fatalf("failed to fetch countries: %v", err)
	}
	fmt.Printf("Loaded %d countries\n", len(countries))
}
```

---

## Transparent Request Batching

WarEra uses **tRPC** and enforces strict request rate limits. Instead of making 10 separate HTTP round-trips for 10 queries, `warera-go` automatically coalesces concurrent requests across goroutines into a single HTTP batch request:

```go
client := warera.NewClient(
	warera.WithBatchWindow(25 * time.Millisecond), // buffer calls for 25ms
)
defer client.Close()

var wg sync.WaitGroup
var country *warera.Country
var dates *warera.GameDates
var prices warera.ItemPrices

wg.Add(3)
go func() { defer wg.Done(); country, _ = client.Countries.GetByID(ctx, "6813b6d446e731854c7ac7a6") }()
go func() { defer wg.Done(); dates, _ = client.GameConfig.GetDates(ctx) }()
go func() { defer wg.Done(); prices, _ = client.Trading.GetPrices(ctx) }()

wg.Wait()
// All 3 queries above were dispatched in 1 single HTTP request!
```

You can also run an explicit batch without waiting for the buffer:

```go
var outCountry warera.Country
var outCompany warera.Company

err := client.Batch(ctx, []warera.BatchCall{
	{Method: "country.getCountryById", Input: map[string]string{"countryId": "c1"}, Output: &outCountry},
	{Method: "company.getById", Input: map[string]string{"companyId": "comp1"}, Output: &outCompany},
})
```

---

## Auto-Pagination with Go 1.23 Iterators

For paginated resources (Articles, Companies, Battles, Users), iterate through every item across pages without writing manual cursor boilerplate:

```go
// Automatically fetches page by page as you loop
for article, err := range client.Articles.All(ctx, warera.ArticleFilter{Type: "last"}) {
	if err != nil {
		log.Fatalf("pagination error: %v", err)
	}
	fmt.Printf("[%s] %s (by %s)\n", article.PublishedAt, article.Title, article.Author)
}
```

---

## Configuration Options

Configure client behavior using functional options:

```go
client := warera.NewClient(
	warera.WithAPIKey("your-api-key"),             // WarEra API Key
	warera.WithBaseURL("https://api2.warera.io/trpc"), // Custom endpoint / gateway proxy
	warera.WithBatchWindow(25 * time.Millisecond), // 0 to disable batch buffering
	warera.WithMaxBatchSize(50),                   // Max calls per batch (default 50)
	warera.WithRateLimit(200),                     // Requests per minute (default 200 with key, 100 without)
	warera.WithMaxRetries(3),                      // Retries on transient errors (default 3)
	warera.WithHTTPClient(&http.Client{}),         // Custom *http.Client
)
```

---

## Available Services

| Service | Methods | Description |
|---|---|---|
| `client.Countries` | `GetAll`, `GetByID` | Country stats, taxes, unrest, allies, wars |
| `client.Companies` | `GetByID`, `GetCompanies`, `All` | Company details; listing methods return company IDs |
| `client.Battles` | `GetByID`, `GetBattles`, `All`, `GetLiveBattleData` | War campaigns, active rounds, live combat |
| `client.Articles` | `GetByID`, `GetLiteByID`, `GetArticlesPaginated`, `All` | In-game newspaper articles and comments |
| `client.Users` | `GetByID`, `GetLite`, `GetUsersByCountry`, `AllByCountry` | Player profiles, level, wealth, damage |
| `client.Regions` | `GetByID`, `GetAll` | Global map regions, resources, development |
| `client.Government` | `GetByCountryID` | Presidents, ministers, congress members |
| `client.GameConfig` | `GetDates`, `GetGameConfig` | Game clock, rollover times, configs |
| `client.Trading` | `GetPrices`, `GetTopOrders` | Item exchange prices and orderbook |
| `client.Search` | `SearchAnything` | Universal search across players, companies, countries |
| `client.Events` | `GetEventsPaginated`, `All` | Paginated world events |
| `client.BattleOrders` | `GetByBattle` | Battle side orders |
| `client.BattleRankings` | `GetRanking` | Battle and war rankings |
| `client.LootSummaries` | `GetByBattleAndUser` | Player battle loot summaries |
| `client.MilitaryUnits` | `GetByID`, `GetManyPaginated`, `All` | Military unit profiles and listings |
| `client.Rounds` | `GetByID`, `GetLastHits` | Battle round details and hit histories |
| `client.Inventory` | `GetCurrentEquipment` | User equipment |
| `client.Rankings` | `GetRanking` | Global game rankings |
| `client.Upgrades` | `GetByTypeAndEntity` | Region, company, and unit upgrades |
| `client.ItemOffers` | `GetByID` | Marketplace item offers |
| `client.MercenaryAuctions` | `GetPaginatedAuctions`, `All` | Mercenary contract auctions |
| `client.Transactions` | `GetPaginatedTransactions`, `All` | Transaction history |
| `client.WorkOffers` | `GetByID`, `GetByCompanyID`, `GetWorkOffersPaginated`, `All` | Company work offers |
| `client.Workers` | `GetTotalWorkersCount`, `GetWorkers` | Worker counts and company groupings |
| `client.RawCall` | `RawCall(ctx, method, input, output)` | Call any unmapped or custom tRPC procedure |

## Model generation and fixtures

Run `make update-specs` to refresh and pretty-print WarEra's upstream OpenAPI document. Then run `make generate` to regenerate request and filter types from `spec/openapi.json`. Response models in `models_response.go` and `models_additional.go` are maintained as Go types in this repository, informed by observed API behavior and recorded fixtures. We do not generate response types from fixtures. See [`spec/SOURCES.md`](spec/SOURCES.md) for source attribution and model policy.

Common development commands are available through `make help`, including `make test`, `make build`, `make fmt`, `make check-generated`, and the fixture recording targets. Live fixture recording uses `WARERA_API_KEY` when provided; endpoint-specific IDs can be supplied to the recorder with `-seeds`.

---

## License

MIT License. See [LICENSE](LICENSE) for details.
