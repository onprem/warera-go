package battle_test

import (
	"encoding/json"
	"math"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/battle"
	"github.com/onprem/warera-go/waretest"
)

func fixture[T any](t *testing.T, procedure string) T {
	t.Helper()
	v, err := waretest.Fixture[T](procedure)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// testdata/battle.json holds the fixture battle's top 45 damage dealers,
// counting both sides, with the item each got from the battle pool, and its
// two ended rounds, recorded from the live API.
type recorded struct {
	Rounds  []warera.Round `json:"rounds"`
	Ranking []struct {
		Rank           int     `json:"rank"`
		Damage         float64 `json:"damage"`
		BattlePoolItem string  `json:"battlePoolItem"`
	} `json:"ranking"`
}

func load(t *testing.T) recorded {
	data, err := os.ReadFile("testdata/battle.json")
	if err != nil {
		t.Fatal(err)
	}
	var r recorded
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatal(err)
	}
	return r
}

func TestPoolsMatchARecordedBattle(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	b := fixture[warera.Battle](t, "battle.getById")
	r := load(t)

	// 9,967,809 + 4,586,082 damage: 29 items of 500k.
	size := battle.PoolSize(cfg, battle.PoolBattle, b.Defender.Damages+b.Attacker.Damages)
	if size != 29 {
		t.Fatalf("battle pool = %d, want 29", size)
	}
	rarityRank := map[string]int{"common": 0, "uncommon": 1, "rare": 2, "epic": 3, "legendary": 4, "mythic": 5}
	prev := math.MaxInt
	for _, e := range r.Ranking {
		got := e.BattlePoolItem != ""
		if got != (e.Rank <= size) {
			t.Errorf("rank %d got a battle pool item: %v; pool size %d", e.Rank, got, size)
		}
		if got {
			rr := rarityRank[cfg.Items[e.BattlePoolItem].Rarity]
			if rr > prev {
				t.Errorf("rank %d got a better item (%s) than the rank above", e.Rank, e.BattlePoolItem)
			}
			prev = rr
		}
	}
	for _, rd := range r.Rounds {
		total := rd.Attacker.Damages + rd.Defender.Damages
		if want := int(total / 200000); battle.PoolSize(cfg, battle.PoolRound, total) != want {
			t.Errorf("round pool for %v damage, want %d", total, want)
		}
	}
	if battle.PoolSize(cfg, battle.PoolRound, 199999) != 0 {
		t.Error("a round below 200k damage has no pool")
	}
}

func TestPoolRarity(t *testing.T) {
	odds := battle.DefaultPoolOdds()
	near := func(a, b, tol float64) bool { return math.Abs(a-b) < tol }
	if got := battle.PoolRarity(1, 1, odds); !near(got[0].Chance, 0.005, 1e-12) || !near(got[4].Chance, 0.5, 1e-12) {
		t.Errorf("one-item pool = %v, want the odds", got)
	}
	// Best of two draws: mythic 1 - 0.995², at least legendary 1 - 0.97².
	best := battle.PoolRarity(2, 1, odds)
	if !near(best[0].Chance, 1-0.995*0.995, 1e-12) || !near(best[1].Chance, (1-0.97*0.97)-(1-0.995*0.995), 1e-12) {
		t.Errorf("rank 1 of 2 = %v", best)
	}
	// Worst of two: mythic 0.005², at least legendary 0.03².
	worst := battle.PoolRarity(2, 2, odds)
	if !near(worst[0].Chance, 0.005*0.005, 1e-12) || !near(worst[1].Chance, 0.03*0.03-0.005*0.005, 1e-12) {
		t.Errorf("rank 2 of 2 = %v", worst)
	}
	if battle.PoolRarity(29, 30, odds) != nil || battle.PoolRarity(29, 0, odds) != nil {
		t.Error("ranks outside the pool get nothing")
	}

	// Rank 13 of 29 against a simulation of sorting 29 draws.
	const trials = 200000
	rng := rand.New(rand.NewPCG(1, 2))
	counts := make([]float64, len(odds))
	draws := make([]int, 29)
	for range trials {
		for i := range draws {
			u, d := rng.Float64(), 0
			for acc := odds[0].Chance; u >= acc && d < len(odds)-1; acc += odds[d].Chance {
				d++
			}
			draws[i] = d
		}
		// The 13th best is the 13th smallest index.
		n := [5]int{}
		for _, d := range draws {
			n[d]++
		}
		seen := 0
		for d := range n {
			if seen += n[d]; seen >= 13 {
				counts[d]++
				break
			}
		}
	}
	var sum float64
	for i, c := range battle.PoolRarity(29, 13, odds) {
		sum += c.Chance
		if !near(c.Chance, counts[i]/trials, 0.004) {
			t.Errorf("rank 13 of 29, %s: %v, simulated %v", c.Rarity, c.Chance, counts[i]/trials)
		}
	}
	if !near(sum, 1, 1e-9) {
		t.Errorf("chances sum to %v", sum)
	}
	values := map[string]float64{"mythic": 1000, "legendary": 100, "epic": 10, "rare": 1}
	if got := battle.ExpectedPoolValue(1, 1, odds, values); !near(got, 5+2.5+1.5+0.32, 1e-9) {
		t.Errorf("ExpectedPoolValue = %v, want 9.32", got)
	}
}

func TestRoundTicks(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	rounds := append(load(t).Rounds, fixture[warera.Round](t, "round.getById"))
	for _, rd := range rounds {
		total := int(rd.Attacker.Points + rd.Defender.Points)
		if got := battle.RoundTicks(cfg, total); got != rd.Live.TicksCount {
			t.Errorf("%d points: %d ticks, API %d", total, got, rd.Live.TicksCount)
		}
		if got := battle.TickPoints(cfg, total); float64(got) != rd.Live.ActualTickPoints {
			t.Errorf("%d points: %d per tick, API %v", total, got, rd.Live.ActualTickPoints)
		}
	}
	// From 0-0, winning every tick: 100 ticks at 1, 50 at 2, 34 at 3.
	if got := battle.TicksToWin(cfg, 0, 0); got != 184 {
		t.Errorf("TicksToWin(0, 0) = %d, want 184", got)
	}
	// At 290-100 the total is 390: 4 per tick, 3 ticks to 302.
	if got := battle.TicksToWin(cfg, 290, 100); got != 3 {
		t.Errorf("TicksToWin(290, 100) = %d, want 3", got)
	}
}
