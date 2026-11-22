# warera-go

[![Go Reference](https://pkg.go.dev/badge/github.com/onprem/warera-go.svg)](https://pkg.go.dev/github.com/onprem/warera-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

An idiomatic Go client for the [WarEra.io](https://warera.io) online multiplayer strategy game API.

The SDK provides typed service clients and response models, cursor-based pagination, and a shared HTTP transport with automatic batching, rate limiting, and retries. It uses only the Go standard library and requires Go 1.23 or newer.

## Install

```bash
go get github.com/onprem/warera-go
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/onprem/warera-go"
)

func main() {
	client := warera.NewClient()
	defer client.Close()

	countries, err := client.Countries.GetAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d countries\n", len(countries))
}
```

## Coverage

The client exposes one service per area of the game:

| Area | Services |
| --- | --- |
| Countries and politics | `Countries` (including diplomacy, unrest, and strategic resources), `Government`, `Alliances`, `Wars`, `Parties`, `Elections`, `Regions` |
| Economy and markets | `Companies`, `Employees`, `Managers`, `Trading` (prices, order books, trading history, public orders), `MarketStats`, `Transactions`, `Donations` |
| Military | `Battles`, `Rounds`, `BattleOrders`, `BattleRankings`, `LootSummaries`, `MilitaryUnits` (including members, leagues, and sponsorships), `MercenaryAuctions`, `Tournaments`, `StateArmy`, `Contributions` |
| Players and community | `Users`, `Search`, `Articles`, `Polls`, `Events`, `Rankings` (including history), `Leagues` |
| Game data | `GameConfig`, `Upgrades`, `Inventory` |

Most procedures come from WarEra's [published OpenAPI document](https://api2.warera.io/openapi.json). Employees, managers, politics, markets beyond prices and top orders, the map, the moderation log, and the newer military procedures are served by the API without being listed there; their documentation is marked `Unstable:` because they may change without notice. Only the employee procedures and `Me` need an API key. `Me` reads the key owner's own employees and their events; the SDK does not wrap procedures that change the account or that API keys are not allowed to call.

## Documentation and examples

See the [Go Reference](https://pkg.go.dev/github.com/onprem/warera-go) for the package overview, client options, service methods, filters, and response models. The [`examples`](examples) directory contains runnable programs:

- `go run ./examples/trader` prompts for a player handle and shows companies with their NPC employees and managers, plus each product's price history, order book, and production-point value.
- `go run ./examples/soldier` prompts for a player, then selects an active or recent battle involving their country and shows profile, equipment, and battle data.
- `go run ./examples/government` prompts for a country and collects a briefing: government, alliance and diplomacy, wars, unrest, congress and parties, citizens, and battles.
- `go run ./examples/commander` prompts for a military unit and shows each member's damage, the unit's league standing and sponsorship, and its orders in an active battle.
- `go run ./examples/treasurer` prompts for a country and shows its treasury, taxes, strategic-resource bonuses, state army spending, donations, and one transaction page.
- `go run ./examples/researcher` prompts for a country and follows its event feed and pinned-article reference.
- `go run ./examples/planner` ranks what a production point earns on each item, quotes input costs from the order book, and values your own workforce per hour (with an API key).
- `go run ./examples/advisor <username>` prints a war report for a player: their fight stats against the best war build for their level, expected and simulated damage, what an 8-hour fight costs and earns, and the gear swaps and skill points that gain the most.
- `go run ./examples/collector` is a long-running collector that appends prices, new trades, and new events to JSON Lines files and resumes where it stopped.
- `go run ./examples/basic` and `go run ./examples/batching` show basic client and batching usage.

## Building tools

- **Errors:** `IsNotFound`, `IsRateLimited`, and `IsTemporary` tell a missing entity from a failure worth retrying.
- **Rate limits:** the client keeps to the API's limits, waits when the server says its quota is used up, and reports requests and waits through `WithHooks`.
- **Polling and resuming:** `Until` stops a newest-first listing at the last record you stored; `Pages` exposes cursors so a collector can resume after a restart.
- **Names:** `GetMany` methods resolve many IDs at once, and `Directory` caches countries, regions, game config, and user profiles.
- **Your own account:** `Me` reads the API key owner's employees, workforce totals, and employee events; create one client per player's key.
- **Economy and combat math:** the `analytics` packages compute production costs, money per production point, employee output and wages, manager pay, company profit, skill build costs, fight stats, expected damage and fight cost, battle bonuses, loot pools and round timing, skill plans and gear swaps, equipment roll quality, upgrade costs, and order-book fills from API records.
- **Testing:** `waretest.NewServer(t)` is a fake API that answers every wrapped procedure with a recorded response; override procedures to simulate errors.

Set `WARERA_API_KEY` in the environment; the persona examples use procedures that require it. They are read-only. [`CONTRIBUTING.md`](CONTRIBUTING.md) describes the development workflow.

## License

MIT. See [LICENSE](LICENSE).
