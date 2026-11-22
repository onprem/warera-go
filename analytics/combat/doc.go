// Package combat estimates what fighting deals, costs, and earns: expected
// damage and health per swing, how far health and food last, wear and
// supplies, case drops and bounties, and the battle bonuses that change
// attack.
//
// Start from a player's stats,
// [github.com/onprem/warera-go/analytics/character.StatsOf] of their
// LoadoutOf, which reproduces the totals the API reports. Set the
// loadout's BattlePercent to [BattlePercent] for a battle, then
// pass the stats to [Window] for a fight and to [CostOf] and [IncomeOf] for
// its gold. [Simulate] fights swing by swing with seeded random rolls, for
// the spread of outcomes around the expected fight.
//
// The game does not publish its damage formula, so everything after the
// stats is a Model in the sense of the analytics package: it follows the
// community Warera Simulator (https://war-era.vercel.app/), uses game
// config values where the config has them, and takes the rest as
// [Assumptions].
package combat
