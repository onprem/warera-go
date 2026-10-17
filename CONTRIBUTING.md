# Contributing to warera-go

Thank you for your interest in contributing to `warera-go`! This guide provides a quick overview to help you get started with development.

## Prerequisites

- **Go 1.23+**: The library utilizes Go 1.23+ range iterators (`iter.Seq2`).
- **WarEra API Key (Optional)**: Needed only if you wish to run manual live tests against `https://api2.warera.io/trpc` or record new fixtures via `cmd/record`. Unit tests use recorded fixtures and in-memory mock transports with zero credentials or network access.

---

## Project Structure

```
warera-go/
├── cmd/
│   ├── generate/     # Model generator from OpenAPI spec and production fixtures
│   └── record/       # Utility to record real live WarEra tRPC responses
├── internal/trpc/    # Low-level tRPC HTTP batching link, retry, and rate limiter
├── spec/
│   ├── fixtures/     # Production JSON response fixtures for tests and codegen
│   ├── Responses.d.ts # WarEra frontend TypeScript interface reference
│   └── openapi.json  # Upstream WarEra OpenAPI 3.0 specification
├── models_gen.go     # Generated domain models (CursorPage, Country, User, Battle, etc.)
├── models.go         # Custom domain helper methods & //go:generate directive
├── iterator.go       # Generic cursor pagination iterator (iter.Seq2)
├── client.go         # Top-level Client and service wiring
├── client_options.go # Functional options (WithAPIKey, WithBatchWindow, etc.)
├── world.go          # Countries, Regions, and Government services
├── military.go       # Battles and Combat services
├── economy.go        # Companies and Market/Trading services
├── community.go      # Articles, Users, and Search services
├── game_config.go    # Game clock & server parameters
└── examples/         # Runnable sample programs
```

---

## Development Workflow

### 1. Run Tests and Builds

Make sure all tests pass and examples build:

```bash
# Run unit and fixture integration tests
go test -v ./...

# Verify example builds
go build ./examples/...
```

### 2. Recording Fixtures & Generating Models

When adding or updating endpoints:
1. **Record a Live Response Fixture:**
   ```bash
   go run ./cmd/record -method <method.name> [-input '<json>']
   ```
   Or record all standard public endpoints:
   ```bash
   go run ./cmd/record -all
   ```
2. **Regenerate Models:**
   ```bash
   go generate ./...
   ```
3. **Implement Service Method:** Add the method in the respective domain service file (e.g. `world.go`, `community.go`, etc.).
4. **Write Unit Tests:** Add test cases in `client_test.go` loading the recorded fixture with `loadFixture(t, "<method>.json")`.

---

## Guidelines

- **Zero External Dependencies:** We strive to keep `warera-go` lightweight by relying solely on the Go standard library. Please avoid adding third-party dependencies unless strictly necessary.
- **Code Formatting:** Format code using standard Go tooling:
  ```bash
  gofmt -s -w .
  ```
- **Commit Messages:** Follow [Conventional Commits](https://www.conventionalcommits.org/) (e.g. `feat:`, `fix:`, `docs:`, `chore:`).

---

## Questions or Issues?

Feel free to open an issue or pull request on GitHub:
👉 https://github.com/onprem/warera-go
