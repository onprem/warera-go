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

See the [Go Reference](https://pkg.go.dev/github.com/onprem/warera-go) for the package overview, client options, service methods, filters, and response models. The [`examples`](examples) directory contains runnable programs, and [`CONTRIBUTING.md`](CONTRIBUTING.md) describes the development workflow.

## License

MIT. See [LICENSE](LICENSE).
