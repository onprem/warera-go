# Changelog

All notable changes to this project are documented here. While the project is in `v0.x`, the public API may change between releases.

## [Unreleased]

### Added

- Typed enums with constants for every value the API's inputs accept, generated from `spec/openapi.json`: `RankingType`, `TransactionType`, `EventType`, `BattleRankingDataType`, `BattleRankingType`, `BattleRankingSide`, `BattleSide`, `UpgradeType`, `MercenaryAuctionStatus`, `ArticleListType`, `BattleDirection`, and `BattleListFilter`; plus hand-written `RankingHistoryRange` and `StateArmyExpenseKind`. Filters, method parameters, and the matching response fields use them, so a typo fails to compile instead of returning a 400.
- `Time` and `Date` types for API timestamps and `YYYY-MM-DD` days. They embed `time.Time`, decode empty strings and null as the zero value, print as RFC 3339, and encode back to the API's formats.
- `WithHooks` observes every HTTP attempt (procedures, status, duration, remaining server quota) and every pause for a retry or rate limit, for logging and metrics.
- `WithUserAgent` identifies your tool to the API. Requests now send a `warera-go/<version>` User-Agent instead of Go's default.
- `WithMaxRetryWait` caps how long one retry waits for the server.
- `IsNotFound`, `IsRateLimited`, and `IsTemporary` classify errors, so tools can tell a missing entity from a failure worth retrying. `ErrNotFound` is returned when a lookup's result is null.

### Changed

- The client follows the server's rate-limit headers: when the server reports its quota exhausted, every request waits until the window resets, so several processes sharing a key slow down instead of failing. Rate-limited and unavailable responses are retried after the server's `Retry-After` or reset time (capped at one minute) instead of a fixed backoff of about 2 seconds.
- The automatic rate limit with an API key is now 500 requests a minute, the API's own limit (was 200).
- 409 Conflict responses are no longer retried.
- **Breaking:** Every timestamp field in the response models is now a `Time` (and every day field, such as `ItemTradingDay.ValueAt`, a `Date`) instead of a string, so tools can sort, compare, and bucket times without parsing. `GameDates.NextDayTime` and `NextRegenTime` are removed; use the `NextDayAt` and `NextRegenAt` fields.
- **Breaking:** `TransactionFilter.TransactionType *StringOrStrings` is now `TransactionTypes []TransactionType`; `StringOrStrings`, `SingleString`, and `StringList` are removed.
- **Breaking:** `BattleOrders.GetByBattle` takes a `BattleSide` and `Rankings.GetRanking` a `RankingType`; `RankingHistoryFilter`, `StateArmyLedgerFilter`, and the generated filters use the new enum types.
- **Breaking:** `Countries.GetStrategicResourceBonuses` returns nil and no error for a country with no strategic resources; the API returns null for these.

### Fixed

- When every call in a batch failed with the same status, such as `GetMany` with only unknown IDs or concurrent lookups that all missed, each call got a shared `*HTTPError` instead of its own `*TRPCError`. Error handling no longer depends on whether calls happened to be batched.
- Lookups whose result is null, such as `Tournaments.GetByID`, `Tournaments.GetTeam`, and `Employees.GetByID` with an unknown ID, returned an empty record and no error. They now return `ErrNotFound`, and bulk lookups report those IDs as failed.

## [0.2.0] - 2026-11-21

This release follows the game's v0.27–v0.28 updates, which replaced player workers with NPC employees and managers, removed Entrepreneurship, and introduced government types, polls, and paged rankings.

### Added

- `Polls` service: `GetByID`, `ListVoters`, and `AllVoters` for article polls.
- `Rankings.List`, `Rankings.All`, and `Rankings.GetAround` for paged rankings beyond the top 100 and for an entity's neighbours.
- `Users.GetSkills` and `Users.GetSkillsMany`, backed by `user.getUserSkills`, because compact user profiles no longer include skills.
- `Employees` service (`GetCompanyRoster`, `GetByID`), `Managers` service (`GetByCompanyID`, `GetOfferByCompanyID`, `ListOffers`, `AllOffers`, `GetWageStats`), and `Companies.GetProductionBonus` for the NPC employee economy. These procedures require an API key and are not in WarEra's published OpenAPI document, so they may change without notice.
- Politics services: `Alliances` (`GetByID`, `GetMany`, `GetStrategicResourceBonuses`), `Wars` (`GetByID`, `GetMany`), `Parties` (`GetByID`, `GetMany`, `List`, `All`, `GetInactiveMemberIDs`), and `Elections` (`GetByID`, `GetMany`, `List`, `All`), plus `Countries.GetDiplomacy`, `Countries.GetUnrest`, `Countries.GetStrategicResourceBonuses`, and `Government.GetCongressSeats`.
- Market and economy procedures: `Trading.GetItemTrading` (30-day price, volume, and trade-count history), `Trading.GetTopOrdersPerItem`, `Trading.GetPublicOrdersByOwner`, `GameConfig.GetMaxBountyRate`, a `MarketStats` service (`GetPPMoneyValue`, `GetPPMoneyValueByItem`, `GetWorldDevelopment`, `GetEquipmentAveragePrice`), and a `Donations` service (`List`, `All`, `GetTotal`). `TopOrder.PartyID` identifies orders placed for a party.
- Military procedures: `MilitaryUnits.GetMembers`, `GetLeague`, `GetLeagueDivisionCounts`, `GetSponsorship`, `GetSponsorshipsByCountry`, `ListActiveSponsorships`, `AllActiveSponsorships`, and `GetSponsorable`; `Tournaments` (`GetByID`, `List`, `All`, `GetTeam`, `GetTeams`), `StateArmy` (`GetOverview`, `ListLedger`, `AllLedger`), `Contributions` (`ListCountryUnrest`, `ListRegion`), and `Leagues` (`GetUserLeague`) services; `Battles.GetActiveIDsByCountry`, `BattleRankings.GetPrivateMercenaryDamage`, `Rankings.GetHistory`, and `Rankings.GetRecentRanks`.
- The politics, market, and military procedures above are not in WarEra's published OpenAPI document and may change without notice.
- `BattleFilter.AllianceID`, plus response fields for government types, income breakdowns, manager taxes, ministers, region company and employee counts, battle outcomes, military-unit types and sponsors, and other values the API now returns.

### Removed

- **Breaking:** `Workers` and `WorkOffers` services and their models and filters. The API no longer serves these procedures.
- **Breaking:** Response fields for removed game features: company worker counts, country specialization and self-work tax, region deposits, strategic-resource production bonuses, fixed minister slots on `Government` (use `Ministers`), worker-era user fields, the Entrepreneurship skill, and game-config settings the API no longer returns.
- **Breaking:** `UserLite.Skills`. Use `Users.GetSkills`.
- **Breaking:** `ItemOffers` service. API keys may no longer list item offers, so offer IDs cannot be discovered.

### Changed

- **Breaking:** Response fields for amounts, money, damage, points, skill and stat values, costs, resources, durations, and percentages are now `float64`, so fractional API values no longer fail the whole call. Counts, levels, ranks, round numbers, quantities, and skill prestige levels remain `int`.
- **Breaking:** Rename `TransactionItem` and `TransactionItemSkills` to `Item` and `ItemSkills`. Battle loot (`BattleLootPoolItem.Item`, `BattleRankingItem.LootItem`) now uses `Item`, which carries the item type and equipment skills that `RoundWeapon` dropped.
- **Breaking:** `GameConfig.Items`, `GameConfig.Badge`, and `GameConfig.Skills` are now maps keyed by item code, badge code, and skill name, with one `GameConfigItem`, `GameConfigBadge`, and `GameConfigSkill` type each, so new game content no longer needs a model change. The per-skill `UserSkills*` types are replaced by one `UserSkill` type.
- **Breaking:** `Users.List` and `Users.All` return `UserListItem` values (ID and join date), which is all `user.getUsersByCountry` returns; they were typed as full `User` profiles with every other field empty.
- **Breaking:** `Ranking` drops `ID` and `IsGlobal`, adds `Count` and `UpdatedAt`, and lists `RankingEntry` values, which identify users, countries, military units, or alliances.
- The trader, government, commander, and treasurer examples show the new employee, politics, market, and military data, and the README lists the services by game area. The trader example reads `WARERA_ITEM_CODES` (comma-separated) instead of `WARERA_ITEM_CODE`.
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

