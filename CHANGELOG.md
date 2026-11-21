# Changelog

All notable changes to this project are documented here. While the project is in `v0.x`, the public API may change between releases.

## [Unreleased]

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

### Changed

- **Breaking:** Response fields for amounts, money, damage, points, skill and stat values, costs, resources, durations, and percentages are now `float64`, so fractional API values no longer fail the whole call. Counts, levels, ranks, round numbers, quantities, and skill prestige levels remain `int`.
- **Breaking:** Rename `TransactionItem` and `TransactionItemSkills` to `Item` and `ItemSkills`. Battle loot (`BattleLootPoolItem.Item`, `BattleRankingItem.LootItem`) now uses `Item`, which carries the item type and equipment skills that `RoundWeapon` dropped.
- **Breaking:** `GameConfig.Items`, `GameConfig.Badge`, and `GameConfig.Skills` are now maps keyed by item code, badge code, and skill name, with one `GameConfigItem`, `GameConfigBadge`, and `GameConfigSkill` type each, so new game content no longer needs a model change. The per-skill `UserSkills*` types are replaced by one `UserSkill` type.
- **Breaking:** `Users.List` and `Users.All` return `UserListItem` values (ID and join date), which is all `user.getUsersByCountry` returns; they were typed as full `User` profiles with every other field empty.
- **Breaking:** `Ranking` drops `ID` and `IsGlobal`, adds `Count` and `UpdatedAt`, and lists `RankingEntry` values, which identify users, countries, military units, or alliances.
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

