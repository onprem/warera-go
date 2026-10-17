# Contributing to warera-go

Thank you for your interest in contributing to `warera-go`! This guide provides a quick overview to help you get started with development.

## Prerequisites

- **Go 1.23+**: The library utilizes Go 1.23+ range iterators (`iter.Seq2`).
- **WarEra API Key (Optional)**: Needed only if you wish to run manual live tests against `https://api2.warera.io/trpc`. Unit tests use in-memory mock transports and require no credentials or network access.

---

## Project Structure

```
warera-go/
├── internal/trpc/    # Low-level tRPC HTTP batching link, retry, and rate limiter
├── spec/             # Upstream WarEra OpenAPI 3.0 spec and response definitions
├── models.go         # Strongly typed request/response structs
├── iterator.go       # Generic cursor pagination iterator
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
# Run unit and integration tests
go test -v ./...

# Verify example builds
go build ./examples/...
```

### 2. Adding a New Endpoint or Service

1. **Check the Specification:** Inspect `spec/openapi.json` or `spec/Responses.d.ts` for parameter names and types.
2. **Define Models:** Add request and response structs in `models.go` with appropriate `json:"..."` tags.
3. **Implement Service Method:** Add the method in the respective service file (or create a new one if it's a new domain).
   ```go
   func (s *CountryService) GetCustom(ctx context.Context, id string) (*CustomResult, error) {
       var res CustomResult
       err := s.client.trpc.Call(ctx, "country.getCustom", customInput{ID: id}, &res)
       return &res, err
   }
   ```
4. **Write Unit Tests:** Add test cases in `client_test.go` or `internal/trpc/client_test.go` using the in-memory `roundTripFunc` pattern (no real network ports needed).

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
