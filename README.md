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

Most procedures come from WarEra's [published OpenAPI document](https://api2.warera.io/openapi.json). Employees, managers, politics, markets beyond prices and top orders, and the newer military procedures are served by the API without being listed there; they need an API key and may change without notice. The SDK does not wrap procedures that act on the caller's own account or that API keys are not allowed to call.

## Documentation and examples

See the [Go Reference](https://pkg.go.dev/github.com/onprem/warera-go) for the package overview, client options, service methods, filters, and response models. The [`examples`](examples) directory contains runnable programs:

- `go run ./examples/trader` prompts for a player handle and shows companies with their NPC employees and managers, plus each product's price history, order book, and production-point value.
- `go run ./examples/soldier` prompts for a player, then selects an active or recent battle involving their country and shows profile, equipment, and battle data.
- `go run ./examples/government` prompts for a country and collects a briefing: government, alliance and diplomacy, wars, unrest, congress and parties, citizens, and battles.
- `go run ./examples/commander` prompts for a military unit and shows each member's damage, the unit's league standing and sponsorship, and its orders in an active battle.
- `go run ./examples/treasurer` prompts for a country and shows its treasury, taxes, strategic-resource bonuses, state army spending, donations, and one transaction page.
- `go run ./examples/researcher` prompts for a country and follows its event feed and pinned-article reference.
- `go run ./examples/basic` and `go run ./examples/batching` show basic client and batching usage.

Set `WARERA_API_KEY` in the environment; the persona examples use procedures that require it. They are read-only. [`CONTRIBUTING.md`](CONTRIBUTING.md) describes the development workflow.

## License

MIT. See [LICENSE](LICENSE).
