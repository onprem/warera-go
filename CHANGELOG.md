# Changelog

All notable changes to this project are documented here. While the project is in `v0.x`, the public API may change between releases.

## [Unreleased]

## [0.1.0] - 2026-10-19

### Added

- Add an idiomatic Go SDK for WarEra's tRPC API, with services covering supported country, region, government, battle, military, economy, inventory, work, community, ranking, upgrade, and game-configuration procedures.
- Add cursor-based pagination, bulk lookup helpers, automatic request batching, rate limiting, retries, and structured API errors.
- Add fixture-backed response coverage and runnable, interactive workflows for traders, soldiers, government members, commanders, treasurers, and researchers.
- Generate request and filter models from the checked-in OpenAPI specification.

### Fixed

- Handle API responses whose JSON shapes vary, including nested battle history and fractional numeric skill values.

