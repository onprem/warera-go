# Changelog

All notable changes to this project are documented here. While the project is in `v0.x`, the public API may change between releases.

## [Unreleased]

### Added

- `analytics` packages: economy and character math on API records, with no requests of their own, each function documented as verified against the API, read from the game config, or a community model. Formulas and verification approach credit the community Warera Simulator. `production.Cost` and `RawMaterials` follow recipes down to raw resources; `production.Values` ranks items by what a production point earns, buying inputs or making them in-house (the latter matches the server's `GetPPMoneyValueByItem`); `company.ProfitPerHour` values a workforce such as `Me.GetWorkforceSummary` at market prices; `character.Build`, `SkillCost`, and `SkillPointsSpent` price skill builds (matching a player's spent points); `character.GradeEquipment` and `RollQuality` place rolled stats within their ranges; `upgrades.Cost` sums upgrade costs (matching invested steel); `market.BuyCost`, `SellProceeds`, `Spread`, and `SummarizeTrading` work on order books and price history.
- `character.StatsOf` and `LoadoutOf` compute a player's fight stats from skill and prestige levels, strategic bonuses, equipment rolls, and attack modifiers, reproducing the totals the API reports (checked on 31 recorded players); change the loadout to compare builds and kits.
- `analytics/combat`: expected damage and health per swing, burst and timed fight windows, wear and supply cost, case and bounty income, `BattlePercent` from the game config's battle bonuses and building levels, and `CraftCost` to value gear. The damage model follows the Warera Simulator, and its unpublished parts are `Assumptions` with the simulator's defaults.
- `analytics/company`: `EmployeeProduction`, `EmployeeWage`, `ManagerWage`, `ManagerBonusPercent`, and the happiness effects of income tax, responsibility, and happiness itself reproduce the API's numbers for recorded rosters, so tools can ask what-if questions such as what a manager would add (`ExtraProduction`). `EngineProduction` reads the automated engine from the game config, and `HourlyOf` gives a company's hourly revenue, costs, profit, break-even price, and hours until storage is full.
- `combat.PillDay` gives the expected fighting of a day on a pill: boosted from full bars, then the after-effect (`PillAfterEffectPercent`) and the rest of the day on what regenerates.
- `combat.Simulate` fights swing by swing with seeded random rolls (dodge, hit, critical hit, damage spread, eating at low health), whose average matches `Window`.
- `analytics/optimize`: `SkillGains` ranks raising each fight skill one level by gain per skill point, and `Roadmap` plans skill purchases level by level, both against a `Metric` such as `Damage`, the expected damage of a fight.
- `optimize.GearChanges` weighs swapping each equipped item one rarity up or down: the metric gained or lost against the price difference, with the swapped loadout for further costing, from prices such as `MarketStats.GetEquipmentAveragePrice`. `market.RollPrice` prices a specific roll from recent equipment sales in `Transactions.List`.
- `analytics/battle`: `PoolSize` gives a battle's or round's loot pool from the damage dealt, which goes to that many top damage dealers counting both sides; `PoolRarity` and `ExpectedPoolValue` give the rarity a rank can expect, as the rank-th best of independent draws at `DefaultPoolOdds`; `TickPoints`, `RoundTicks`, and `TicksToWin` give round timing. Pool sizes, the ranks that get items, and tick counts are checked against recorded battles.
- `market.Basis` values trade at reference prices (`Reference`) or at the order books' best bid and ask (`BidAsk`).
- `examples/planner`: a production planner built on the `analytics` packages.
- `examples/advisor`: a war report for a player by username, comparing their build with the best war build for their level, with damage over a day on a pill, simulations, fight economics, and gear swaps judged by what their extra damage costs. Flags set the pill, battle bonuses, bounty rate, and goal (damage or cost).
- `waretest.Fixture` decodes a recorded API response, for tests and examples of code that computes from records.

### Fixed

- `ManagerWageRange` documentation: its values are the hourly wage at Management value 11, not a rate per Management point.

## [0.3.0] - 2026-11-21

This release is for tools that run for a long time: errors and rate limits behave the same whether or not calls were batched, inputs and timestamps are typed, and new helpers resume listings, cache names, and fake the API in tests. It also wraps the world map, the moderation log, and the API key owner's employees. Many changes are breaking; they are marked below.

### Added

- Typed enums with constants for every value the API's inputs accept, generated from `spec/openapi.json`: `RankingType`, `TransactionType`, `EventType`, `BattleRankingDataType`, `BattleRankingType`, `BattleRankingSide`, `BattleSide`, `UpgradeType`, `MercenaryAuctionStatus`, `ArticleListType`, `BattleDirection`, and `BattleListFilter`; plus hand-written `RankingHistoryRange` and `StateArmyExpenseKind`. Filters, method parameters, and the matching response fields use them, so a typo fails to compile instead of returning a 400.
- `Time` and `Date` types for API timestamps and `YYYY-MM-DD` days. They embed `time.Time`, decode empty strings and null as the zero value, print as RFC 3339, and encode back to the API's formats.
- `Pages` iterates whole pages from any `List` method, so long-running collectors can save each page's `NextCursor` and resume after a restart. `Collect` gathers an iterator's values (optionally up to a limit), and `Until` stops an iterator at a value already seen, so polling the newest-first event and transaction listings fetches only what is new.
- `Directory` caches countries, regions, the game configuration, and compact user profiles with a TTL, so tools can turn IDs into names without refetching. Concurrent callers share one fetch, and user lookups fetch only IDs not already cached.
- `waretest` package: a fake WarEra API server for testing tools without network access. It answers every wrapped procedure with a response recorded from the live API, speaks tRPC's batch protocol, and lets tests override procedures with their own data, null results, or API errors (`waretest.NotFound`, `RateLimited`, …). `Server.Calls` and `Server.Inputs` report what a tool sent. The recorded fixtures moved from `spec/fixtures` to `waretest/fixtures`.
- `examples/collector`: a long-running, read-only collector that stores prices, new trades, and new events as JSON Lines and resumes after a restart.
- `Battles.GetMany`, `Articles.GetMany`, and `Tournaments.GetMany`.
- `Contributions.AllCountryUnrest` and `Contributions.AllRegion` iterate the page-numbered contribution listings.
- `WithHooks` observes every HTTP attempt (procedures, status, duration, remaining server quota) and every pause for a retry or rate limit, for logging and metrics.
- `WithUserAgent` identifies your tool to the API. Requests now send a `warera-go/<version>` User-Agent instead of Go's default.
- `WithMaxRetryWait` caps how long one retry waits for the server.
- `IsNotFound`, `IsRateLimited`, and `IsTemporary` classify errors, so tools can tell a missing entity from a failure worth retrying. `ErrNotFound` is returned when a lookup's result is null.
- `Regions.GetMapData` returns the world map the game draws: region and country shapes as TopoJSON, sea links, and the arrows of ongoing attacks.
- `Upgrades.GetMapUpgrades` returns every region's bunker, base, and pacification center with its level and status in one call.
- `Sanctions` service (`List`, `All`) reads the public moderation log: bans, mutes, warnings, fines, and username resets, filterable by player and `SanctionType`.
- `Managers.GetPointValue` returns how much a company produces per manager point.
- `Me` service reads the account that owns the API key: `GetEmployees` (all its employees, with `Employee.FoundToolboxes`), `GetWorkforceSummary` (hourly production and costs), `GetPendingEmployeeEvents`, `GetPendingEmployeeEventCount`, and `GetEmployeeEventHistory` (an employee's story events with their choices, success chances, and outcomes). They need an API key; the API forbids keys from the owner's other data, such as the profile and inventory.
- `MercenaryAuctions.GetBidEligibility` reports whether a military unit may bid on each auction, and `MilitaryUnits.GetReputationBuyStatus` what buying mercenary reputation costs.
- `UserBuffs.BuffCodes` and `BuffEndAt` report a player's active boosts, such as pills; only the after-effects were modelled before. `User` and `UserLite` also gain `AnimatedAvatarURL`, `Infos.IsBanned`, `Infos.Description`, `Infos.BannerID`, `Mute`, and newer `Dates` fields, and `UserStats` gains `EstimatedWealth`, `EstimatedInventoryValue`, `EstimatedCompanyValues`, `DamagesRanking`, `CasesOpenedCount`, and `GemsPurchased`. Checked by strictly decoding 66 live profiles, including banned and muted players.

### Changed

- The client follows the server's rate-limit headers: when the server reports its quota exhausted, every request waits until the window resets, so several processes sharing a key slow down instead of failing. Rate-limited and unavailable responses are retried after the server's `Retry-After` or reset time (capped at one minute) instead of a fixed backoff of about 2 seconds.
- The automatic rate limit with an API key is now 500 requests a minute, the API's own limit (was 200).
- 409 Conflict responses are no longer retried.
- **Breaking:** Every timestamp field in the response models is now a `Time` (and every day field, such as `ItemTradingDay.ValueAt`, a `Date`) instead of a string, so tools can sort, compare, and bucket times without parsing. `GameDates.NextDayTime` and `NextRegenTime` are removed; use the `NextDayAt` and `NextRegenAt` fields.
- **Breaking:** Response models no longer use `omitempty`, so re-encoding a record (for example to store it) keeps `false`, `0`, and empty values instead of dropping them. Request filters keep it.
- **Breaking:** Fields holding entity IDs now end in `ID` or `IDs`, like the rest of the models. For example `Country.Enemy` → `EnemyID`, `Country.Allies` → `AllyIDs`, `Government.President` → `PresidentID`, `MilitaryUnit.Members` → `MemberIDs`, `EventData.AttackerCountry` → `AttackerCountryID`, `User.Party` → `PartyID`, `Company.Region` → `RegionID`, and `UserEquipment.Helmet` → `HelmetID`; 62 fields in all.
- **Breaking:** Every `GetMany` method (and `Users.GetLiteMany` and `GetSkillsMany`) returns a map keyed by ID and ignores duplicate and empty IDs. When some lookups fail it returns the records it found together with a `BatchLookupError` naming only the failed IDs, so one deleted account no longer discards hundreds of resolved names. `Alliances.GetMany` follows the same contract, splits inputs over 100 IDs, and reports IDs the API omits as not found.
- **Breaking:** `BattleRanking` embeds `CursorPage[BattleRankingEntry]` like the other page types; its `Items` and `NextCursor` fields are unchanged for readers.
- **Breaking:** `GameConfig.GetGameConfig` is now `GameConfig.Get`, and `CompanyFilter.PerPage` is `Limit` like every other filter's page size.
- **Breaking:** `BattleLiveBattleDataFilter`, `BattleLootSummaryByBattleAndUserFilter`, `BattleOrderByBattleFilter`, and `TradingOrderTopOrdersFilter` are no longer exported; no public method accepted them.
- **Breaking:** `TransactionFilter.TransactionType *StringOrStrings` is now `TransactionTypes []TransactionType`; `StringOrStrings`, `SingleString`, and `StringList` are removed.
- **Breaking:** `BattleOrders.GetByBattle` takes a `BattleSide` and `Rankings.GetRanking` a `RankingType`; `RankingHistoryFilter`, `StateArmyLedgerFilter`, and the generated filters use the new enum types.
- **Breaking:** `Countries.GetStrategicResourceBonuses` returns nil and no error for a country with no strategic resources; the API returns null for these.
- **Breaking:** `UserStats.Case1`, `Case2`, and `Toolbox` are `CaseStats`, with counts keyed by rarity, and `WoodenCase` is `WoodenCaseStats`, with counts keyed by item code, so new rarities and items decode without a release. The per-case `UserStats*` types are removed.
- **Breaking:** `Client.RawCall` and `Client.Batch` send no input when the input is nil, which procedures that take no input require; pass `struct{}{}` to send `{}`. The client no longer keeps a list of such procedures.
- Methods that wrap procedures missing from the published OpenAPI document say `Unstable:` in their documentation, and only those that need an API key say so. Only the employee procedures do; the docs used to say every unlisted procedure needed one.

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

