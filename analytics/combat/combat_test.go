package combat_test

import (
	"errors"
	"math"
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
	"github.com/onprem/warera-go/analytics/character"
	"github.com/onprem/warera-go/analytics/combat"
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

func near(a, b float64) bool { return math.Abs(a-b) < 1e-6 }

// recordedStats are the fixture player's stats: attack 347.82 under the
// pill's after-effect, precision 100, critical chance 74, critical damage
// 317, armor 68%, dodge 56%, health 150, hunger 8, loot chance 13.
func recordedStats(t *testing.T, cfg *warera.GameConfig) character.Stats {
	user := fixture[warera.User](t, "user.getUserById")
	eq := fixture[warera.CurrentEquipment](t, "inventory.fetchCurrentEquipment")
	s, err := character.StatsOf(cfg, character.LoadoutOf(user.Skills, &eq))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestExpectedSwing(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	a := combat.DefaultAssumptions()

	// Every swing hits: 26% normal, 74% critical at 4.17×.
	s := combat.ExpectedSwing(cfg, recordedStats(t, cfg), a)
	if want := 347.82 * (0.26 + 0.74*4.17); !near(s.Damage, want) {
		t.Errorf("damage = %v, want %v", s.Damage, want)
	}
	// 10 health × (1 - 0.68) = 3.2, paid on the 44% of swings not dodged.
	if !near(s.HitHealth, 3.2) || !near(s.Health, 3.2*0.44) {
		t.Errorf("health = %v per hit, %v expected; want 3.2 and 1.408", s.HitHealth, s.Health)
	}

	// 40% misses at half damage, 30% hits, 30% critical hits at 2×.
	s = combat.ExpectedSwing(cfg, character.Stats{Attack: 100, Precision: 60, CriticalChance: 50, CriticalDamages: 100}, a)
	if !near(s.Damage, 110) || s.HitHealth != 10 {
		t.Errorf("damage %v, hit health %v; want 110 and 10", s.Damage, s.HitHealth)
	}
}

func TestWindowCostAndIncome(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	stats := recordedStats(t, cfg)
	a := combat.DefaultAssumptions()
	swing := combat.ExpectedSwing(cfg, stats, a)

	// Cooked fish heals 20% of 150. Burst: 150 + 8 foods × 30 = 390 health,
	// 390 / 1.408 → 277 swings.
	burst, err := combat.Window(cfg, stats, "cookedFish", 0, a)
	if err != nil {
		t.Fatal(err)
	}
	if burst.Health != 390 || burst.Food != 8 || burst.Swings != 277 || !near(burst.Damage, 277*swing.Damage) {
		t.Errorf("burst = %+v", burst)
	}
	// 8 hours regenerate 80% of each bar: 150 + 120 health and 8 + 6 whole
	// foods, 270 + 14 × 30 = 690 health, 690 / 1.408 → 491 swings.
	day, err := combat.Window(cfg, stats, "cookedFish", 8, a)
	if err != nil {
		t.Fatal(err)
	}
	if day.Health != 690 || day.Food != 14 || day.Swings != 491 {
		t.Errorf("8h = %+v", day)
	}
	if _, err := combat.Window(cfg, stats, "caviar", 0, a); !errors.Is(err, analytics.ErrUnknownItem) {
		t.Errorf("unknown food error = %v", err)
	}
	if _, err := combat.Window(cfg, character.Stats{Attack: 1, ArmorReduction: 1}, "", 0, a); !errors.Is(err, combat.ErrNoHealthCost) {
		t.Errorf("free swings error = %v", err)
	}

	// 277 swings wear 2.77% of a 100 weapon, 277 × 0.44 undodged swings
	// 1.2188% of 200 armor; 277 ammo at 0.2, 8 food at 7, and a 35 pill.
	cost := combat.CostOf(burst, combat.Prices{Weapon: 100, Armor: 200, Ammo: 0.2, Food: 7, Boosts: 35}, a)
	if !near(cost.Wear, 277+243.76) || !near(cost.Ammo, 55.4) || cost.Food != 56 || !near(cost.Total, 520.76+55.4+56+35) {
		t.Errorf("cost = %+v", cost)
	}
	// 277 hits at 13% drop 36.01 case1 and 0.3601 case2.
	in := combat.IncomeOf(burst, stats.LootChance, combat.Rewards{Case1: 4, Case2: 28, BountyPer1k: 1})
	if !near(in.Case1, 36.01) || !near(in.Case2, 0.3601) || !near(in.Cases, 36.01*4+0.3601*28) || !near(in.Bounty, burst.Damage/1000) {
		t.Errorf("income = %+v", in)
	}
	if got := burst.PerThousand(cost.Total); !near(got, cost.Total/burst.Damage*1000) {
		t.Errorf("PerThousand = %v", got)
	}
}

func TestBattlePercent(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	// Country order 10, patriotic 20, headquarters level 2 10, bunker
	// level 1 5, not linked to the capital -25.
	got, err := combat.BattlePercent(cfg, combat.Bonuses{
		CountryOrder: true, Patriotic: true, HeadquartersLevel: 2, BunkerLevel: 1, NotLinkedToCapital: true,
	})
	if err != nil || got != 20 {
		t.Errorf("BattlePercent = %v, %v; want 20", got, err)
	}
	// -50 - 25 - 25 - 20 stops at -100.
	got, _ = combat.BattlePercent(cfg, combat.Bonuses{OccupyingYourRegions: true, LostAttackingRegion: true, NotLinkedToCapital: true, Extra: -20})
	if got != -100 {
		t.Errorf("BattlePercent floor = %v, want -100", got)
	}
	if _, err := combat.BattlePercent(cfg, combat.Bonuses{HeadquartersLevel: 9}); !errors.Is(err, analytics.ErrUnknownUpgrade) {
		t.Errorf("unknown level error = %v", err)
	}
}

func TestCraftCost(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	// Legendary: 32 steel at 2 and 486 scraps at 0.25.
	got, err := combat.CraftCost(cfg, "legendary", warera.ItemPrices{"steel": 2, "scraps": 0.25})
	if err != nil || got != 185.5 {
		t.Errorf("CraftCost = %v, %v; want 185.5", got, err)
	}
	if _, err := combat.CraftCost(cfg, "shiny", nil); err == nil {
		t.Error("unknown rarity: want an error")
	}
}
