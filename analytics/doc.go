// Package analytics is the root of a family of packages that compute economy,
// character, and battle numbers from records the WarEra API returns. Each
// subpackage covers one area:
//
//   - [github.com/onprem/warera-go/analytics/production]: what an item costs
//     to produce and what a production point earns on it.
//   - [github.com/onprem/warera-go/analytics/market]: order-book fills,
//     spreads, price bases, and price history.
//   - [github.com/onprem/warera-go/analytics/company]: employee output and
//     wages, manager pay, and what companies and a workforce make per hour.
//   - [github.com/onprem/warera-go/analytics/character]: skill builds and their
//     cost, fight stats, and equipment roll quality.
//   - [github.com/onprem/warera-go/analytics/combat]: expected damage, how far
//     health and food last, what a fight costs and earns, and battle bonuses.
//   - [github.com/onprem/warera-go/analytics/battle]: loot pool sizes, the
//     ranks that get items and the rarity to expect, and round ticks.
//   - [github.com/onprem/warera-go/analytics/optimize]: which skill levels and
//     gear swaps gain the most per point or gold, and a skill plan for the
//     levels ahead.
//   - [github.com/onprem/warera-go/analytics/upgrades]: what region and
//     company upgrades cost.
//
// Every function is pure: it takes records fetched with the warera package,
// such as [warera.GameConfig] from GameConfig.Get and [warera.ItemPrices]
// from Trading.GetPrices, and makes no requests. Fetch the game
// configuration once and refresh it rarely; [warera.Directory] caches it.
//
// # How far to trust a number
//
// The game publishes some rules in its configuration and keeps others on its
// server. The documentation of each function says which kind it computes:
//
//   - Verified: recomputed from API records and checked against numbers the
//     API itself reports, with a test against recorded responses.
//   - Game config: read from a rule [warera.GameConfig] publishes.
//   - Model: a community assumption the API cannot confirm, such as how much
//     damage a missed hit deals. Models take their assumptions as parameters
//     with documented defaults, so tools can tune them.
//
// Where a community model and the game configuration disagree, the packages
// follow the configuration and offer the community value as an option.
//
// # Credits
//
// Many formulas in these packages, and the idea of checking them against the
// API, follow the community Warera Simulator (https://war-era.vercel.app/),
// which models combat, economy, and battles in far more depth. Its rules for
// skill totals, employee output and wages, manager pay, and battle loot were
// recomputed here against live API data before being adopted.
//
// Errors for codes and levels the configuration does not define are shared
// across the packages: [ErrUnknownItem], [ErrUnknownSkill], and
// [ErrUnknownUpgrade].
package analytics
