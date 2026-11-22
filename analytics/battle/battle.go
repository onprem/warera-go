package battle

import (
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/onprem/warera-go"
)

// Pool is one of a battle's loot pools: the battle pool, shared out when the
// battle ends, or a round pool, shared out when a round ends.
type Pool int

const (
	PoolBattle Pool = iota
	PoolRound
)

// PoolSize returns how many items a loot pool holds for the damage both
// sides dealt in the battle or the round: one item per
// BattleLootDamagePerLootItem (battle) or DamagePerLootItem (round) of the
// game config's Loot, rounded down. The top PoolSize damage dealers of the
// battle or round, counting both sides, each get one item.
//
// Verified: a recorded battle's pool of 29 items went to exactly ranks 1
// to 29.
func PoolSize(cfg *warera.GameConfig, pool Pool, damage float64) int {
	per := cfg.Loot.BattleLootDamagePerLootItem
	if pool == PoolRound {
		per = cfg.Loot.DamagePerLootItem
	}
	if per <= 0 || damage <= 0 {
		return 0
	}
	return int(damage / per)
}

// RarityChance is the chance of one rarity.
type RarityChance struct {
	Rarity string
	Chance float64
}

// DefaultPoolOdds returns the rarity odds of a pool item, best rarity first:
// 0.5% mythic, 2.5% legendary, 15% epic, 32% rare, and 50% uncommon.
//
// Model: the Warera Simulator's odds. 2,604 recorded pool items came out
// 0.31%, 2.61%, 15.55%, 31.84%, and 49.69%.
func DefaultPoolOdds() []RarityChance {
	return []RarityChance{
		{"mythic", 0.005}, {"legendary", 0.025}, {"epic", 0.15}, {"rare", 0.32}, {"uncommon", 0.5},
	}
}

// PoolRarity returns the chance of each rarity for the item that the
// rank-th damage dealer gets from a pool of size items, given each item's
// odds, best rarity first. Items are handed out best first, so rank r gets
// the r-th best of size independent draws. It returns nil when rank is not
// between 1 and size.
//
// Verified: items go to ranks best first. That each item is an
// independent draw follows the Warera Simulator.
func PoolRarity(size, rank int, odds []RarityChance) []RarityChance {
	if rank < 1 || rank > size {
		return nil
	}
	out := make([]RarityChance, len(odds))
	var atLeast, below float64
	for i, o := range odds {
		atLeast += o.Chance
		// The r-th best is at least this good when at least r of the size
		// draws are.
		p := binomialTail(size, rank, atLeast)
		out[i] = RarityChance{Rarity: o.Rarity, Chance: max(0, p-below)}
		below = p
	}
	return out
}

// ExpectedPoolValue returns the expected value of the item the rank-th
// damage dealer gets from a pool of size items, with each rarity's value
// in values. It is zero when rank does not get an item.
func ExpectedPoolValue(size, rank int, odds []RarityChance, values map[string]float64) float64 {
	var total float64
	for _, c := range PoolRarity(size, rank, odds) {
		total += c.Chance * values[c.Rarity]
	}
	return total
}

// binomialTail returns P(X ≥ k) for X ~ Binomial(n, p).
func binomialTail(n, k int, p float64) float64 {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return 1
	}
	lnN, _ := math.Lgamma(float64(n + 1))
	var sum float64
	for i := k; i <= n; i++ {
		lnI, _ := math.Lgamma(float64(i + 1))
		lnRest, _ := math.Lgamma(float64(n - i + 1))
		sum += math.Exp(lnN - lnI - lnRest + float64(i)*math.Log(p) + float64(n-i)*math.Log1p(-p))
	}
	return min(1, sum)
}

// TickDuration is the time between round ticks, as observed on live rounds;
// the game config does not publish it.
const TickDuration = 2 * time.Minute

// TickPoints returns the points the winner of a round tick gets when both
// sides together have totalPoints: the game config's TickPoints value at
// the highest threshold not above totalPoints.
//
// Verified: ended rounds' tick counts match [RoundTicks].
func TickPoints(cfg *warera.GameConfig, totalPoints int) int {
	thresholds := make([]int, 0, len(cfg.Battle.TickPoints))
	for key := range cfg.Battle.TickPoints {
		if t, err := strconv.Atoi(key); err == nil {
			thresholds = append(thresholds, t)
		}
	}
	if len(thresholds) == 0 {
		return 1
	}
	slices.Sort(thresholds)
	points := cfg.Battle.TickPoints[strconv.Itoa(thresholds[0])]
	for _, t := range thresholds {
		if t <= totalPoints {
			points = cfg.Battle.TickPoints[strconv.Itoa(t)]
		}
	}
	return points
}

// RoundTicks returns how many ticks a round takes for both sides together
// to score totalPoints, whoever wins each tick.
//
// Verified: equals [warera.RoundLiveState] TicksCount of ended rounds.
func RoundTicks(cfg *warera.GameConfig, totalPoints int) int {
	ticks := 0
	for total := 0; total < totalPoints; ticks++ {
		total += TickPoints(cfg, total)
	}
	return ticks
}

// TicksToWin returns how many ticks a side with points needs to reach the
// game config's PointsToWinRound if it wins every tick from now on, against
// a side with otherPoints. Multiply by [TickDuration] for the time.
func TicksToWin(cfg *warera.GameConfig, points, otherPoints int) int {
	target := int(cfg.Battle.PointsToWinRound)
	ticks := 0
	for ; points < target; ticks++ {
		points += TickPoints(cfg, points+otherPoints)
	}
	return ticks
}
