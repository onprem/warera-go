# Changelog

All notable changes to this project are documented here. While the project is in `v0.x`, the public API may change between releases.

## [Unreleased]

### Changed

- The rate limiter now enforces a sliding 60-second window. Previously a new client could send up to twice its per-minute limit in its first minute.

### Fixed

- A batched call whose context ended no longer has its output written after it returns.
- Caller deadlines now bound batched requests: an in-flight batch is canceled once every caller in it has gone away.

## [0.1.0] - 2026-10-19

### Added

- Add an idiomatic Go SDK for WarEra's tRPC API, with services covering supported country, region, government, battle, military, economy, inventory, work, community, ranking, upgrade, and game-configuration procedures.
- Add cursor-based pagination, bulk lookup helpers, automatic request batching, rate limiting, retries, and structured API errors.
- Add fixture-backed response coverage and runnable, interactive workflows for traders, soldiers, government members, commanders, treasurers, and researchers.
- Generate request and filter models from the checked-in OpenAPI specification.

### Fixed

- Handle API responses whose JSON shapes vary, including nested battle history and fractional numeric skill values.

