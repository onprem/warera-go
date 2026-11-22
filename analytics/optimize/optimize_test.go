package optimize_test

import (
	"math"
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/character"
	"github.com/onprem/warera-go/analytics/optimize"
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

func attack(s character.Stats) (float64, error) { return s.Attack, nil }

// The fixture player is level 36 with 1 point spare. Their attack modifiers
// multiply by 1.1 × 0.4 × 1.55 = 0.682.
func recorded(t *testing.T) (*warera.GameConfig, character.Loadout) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	user := fixture[warera.User](t, "user.getUserById")
	eq := fixture[warera.CurrentEquipment](t, "inventory.fetchCurrentEquipment")
	return cfg, character.LoadoutOf(user.Skills, &eq)
}

func TestSkillGains(t *testing.T) {
	cfg, l := recorded(t)
	gains, err := optimize.SkillGains(cfg, l, 36, attack)
	if err != nil {
		t.Fatal(err)
	}
	near := func(a, b float64) bool { return math.Abs(a-b) < 1e-9 }
	// Attack 7 → 8 adds 25 × 0.682 for 8 points; precision 6 → 7 adds 5
	// points of overflow, 20 × 0.682 attack, for 7 points.
	if len(gains) < 2 || gains[0].Skill != "attack" || gains[0].Level != 8 || gains[0].Points != 8 || !near(gains[0].Gain, 25*0.682) ||
		gains[1].Skill != "precision" || gains[1].Points != 7 || !near(gains[1].Gain, 20*0.682) {
		t.Fatalf("SkillGains = %+v", gains)
	}
	for _, g := range gains[2:] {
		if g.Gain != 0 {
			t.Errorf("%s changes attack by %v", g.Skill, g.Gain)
		}
	}
	// Armor's levels unlock at player level 5.
	for _, g := range must(optimize.SkillGains(cfg, character.Loadout{}, 4, attack)) {
		if g.Skill == "armor" {
			t.Error("armor level 1 is not unlocked at player level 4")
		}
	}
}

func TestRoadmap(t *testing.T) {
	cfg, l := recorded(t)
	steps, final, err := optimize.Roadmap(cfg, l, 36, 38, attack)
	if err != nil {
		t.Fatal(err)
	}
	// Level 37: 5 points, saved for attack 8 (8 points). Level 38: 9
	// points buy attack 8; the next best, precision 7, costs 7 > 1.
	if len(steps) != 2 || len(steps[0].Raised) != 0 || steps[0].SparePoints != 5 ||
		len(steps[1].Raised) != 1 || steps[1].Raised[0].Skill != "attack" || steps[1].SparePoints != 1 {
		t.Fatalf("Roadmap = %+v", steps)
	}
	if final.Levels["attack"] != 8 || l.Levels["attack"] != 7 {
		t.Errorf("final attack %d, original %d; want 8 and an untouched 7", final.Levels["attack"], l.Levels["attack"])
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
