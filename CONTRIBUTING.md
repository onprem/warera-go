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
│   ├── generate/     # Request model generator from the OpenAPI spec
│   ├── record/       # Utility to record real live WarEra tRPC responses
│   └── update-specs/ # Refreshes upstream contract snapshots
├── internal/trpc/    # Low-level tRPC HTTP batching link, retry, and rate limiter
├── spec/
│   ├── fixtures/     # Recorded JSON responses used by tests and recorder seeds
│   ├── Responses.d.ts # Generated response declarations from WarEraProjects/TRPC
│   ├── openapi.json  # WarEra's published OpenAPI 3 specification
│   └── SOURCES.md    # Source attribution and refresh instructions
├── models_gen.go     # Generated request/filter types and CursorPage
├── models.go         # Custom domain helper methods & //go:generate directive
├── models_response.go # Maintained response models
├── models_additional.go # Response models based on upstream Responses.d.ts
├── iterator.go       # Generic cursor pagination iterator (iter.Seq2)
├── client.go         # Top-level Client and service wiring
├── client_options.go # Functional options (WithAPIKey, WithBatchWindow, etc.)
├── world.go          # Countries, Regions, and Government services
├── military.go       # Battles and Combat services
├── economy.go        # Companies and Market/Trading services
├── inventory.go      # User equipment service
├── ranking.go        # Global ranking service
├── upgrades.go       # Upgrade service
├── work.go           # Work offers and workers services
├── community.go      # Articles, Users, and Search services
├── game_config.go    # Game clock & server parameters
└── examples/         # Runnable sample programs
```

---

## Development Workflow

### Refreshing upstream specifications

Run this from the repository root when you want to update the checked-in contract snapshots:

```bash
make update-specs
```

This requires Go 1.23+ and network access. It refreshes and formats `spec/openapi.json` from WarEra and refreshes `spec/Responses.d.ts` from WarEraProjects/TRPC with a pinned source commit and license link in its header. See [`spec/SOURCES.md`](spec/SOURCES.md) for attribution. Review these files and their diff before regenerating models.

### 1. Run Tests and Builds

Make sure all tests pass and examples build:

```bash
make test
make build
```

### 2. Recording Fixtures & Generating Models

When adding or updating endpoints:
1. **Record Live Response Fixtures:**
   - **Discover & record missing fixtures** (skips existing files on disk):
     ```bash
     make record
     ```
   - **Re-record all fixtures cleanly** (overwriting disk):
     ```bash
     make record-force
     ```
   - **Target a specific endpoint** (replaces local fixture by default):
     ```bash
     go run ./cmd/record -method <method.name> [-input '<json>']
     ```
   - **Supply IDs for bulk endpoints without a local fixture seed**:
     ```bash
     go run ./cmd/record -seeds '{"itemOfferId":"...","workOfferId":"..."}'
     ```
     The recorder derives `companyId` from `company.getById.json` and `workOfferId` from a non-empty work-offer list fixture when available.
2. **Regenerate Request Models:**
   ```bash
   make generate
   ```
   Run `make check-generated` to confirm checked-in generated code matches the OpenAPI spec.
3. **Implement Service Method:** Add the method in the respective domain service file (e.g. `world.go`, `community.go`, etc.).
4. **Write Unit Tests:** Add test cases in `client_test.go` loading the recorded fixture with `loadFixture(t, "<method>.json")`.

---

## Guidelines

- **Zero External Dependencies:** We strive to keep `warera-go` lightweight by relying solely on the Go standard library. Please avoid adding third-party dependencies unless strictly necessary.
- **Code Formatting:** Use the Make target for consistent formatting:
  ```bash
  make fmt
  ```
- **Commit Messages:** Follow [Conventional Commits](https://www.conventionalcommits.org/) (e.g. `feat:`, `fix:`, `docs:`, `chore:`).

---

## Questions or Issues?

Feel free to open an issue or pull request on GitHub:
👉 https://github.com/onprem/warera-go
