package battle_test

import (
	"fmt"
	"time"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/battle"
	"github.com/onprem/warera-go/waretest"
)

// must stands in for error handling, keeping the examples short.
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// See what the 13th damage dealer of a battle can expect from its pool.
func ExamplePoolRarity() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig")) // client.GameConfig.Get(ctx)
	b := must(waretest.Fixture[warera.Battle]("battle.getById"))                  // client.Battles.GetByID(ctx, id)

	size := battle.PoolSize(cfg, battle.PoolBattle, b.Attacker.Damages+b.Defender.Damages)
	fmt.Printf("%d items for the top %d\n", size, size)
	for _, c := range battle.PoolRarity(size, 13, battle.DefaultPoolOdds()) {
		fmt.Printf("%-9s %5.1f%%\n", c.Rarity, c.Chance*100)
	}
	// Output:
	// 29 items for the top 29
	// mythic      0.0%
	// legendary   0.0%
	// epic        0.1%
	// rare       77.0%
	// uncommon   22.9%
}

// Estimate how soon a side can win the current round.
func ExampleTicksToWin() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig")) // client.GameConfig.Get(ctx)
	round := must(waretest.Fixture[warera.Round]("round.getById"))                // client.Rounds.GetByID(ctx, id)

	ticks := battle.TicksToWin(cfg, int(round.Attacker.Points), int(round.Defender.Points))
	fmt.Printf("attacker at %v-%v wins in %d ticks (%v) at the earliest\n",
		round.Attacker.Points, round.Defender.Points, ticks, battle.TickDuration*time.Duration(ticks))
	// Output: attacker at 23-0 wins in 161 ticks (5h22m0s) at the earliest
}
