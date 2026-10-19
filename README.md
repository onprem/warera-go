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

## Documentation and examples

See the [Go Reference](https://pkg.go.dev/github.com/onprem/warera-go) for the package overview, client options, service methods, filters, and response models. The [`examples`](examples) directory contains runnable programs:

- `go run ./examples/trader` prompts for a player handle and shows companies, workers, wages, and current market context.
- `go run ./examples/soldier` prompts for a player, then selects an active or recent battle involving their country and shows profile, equipment, and battle data.
- `go run ./examples/government` prompts for a country and collects a country briefing.
- `go run ./examples/basic` and `go run ./examples/batching` show basic client and batching usage.

Set `WARERA_API_KEY` in the environment when the API requires authentication. The persona examples are read-only. [`CONTRIBUTING.md`](CONTRIBUTING.md) describes the development workflow.

## License

MIT. See [LICENSE](LICENSE).
