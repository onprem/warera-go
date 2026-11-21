package character_test

import (
	"errors"
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
	"github.com/onprem/warera-go/analytics/character"
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

func TestSkillsAgainstRecordedPlayer(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	user := fixture[warera.User](t, "user.getUserById")
	spent, err := character.SkillPointsSpent(cfg, user.Skills)
	if err != nil {
		t.Fatal(err)
	}
	if want := float64(user.Leveling.SpentSkillPoints); spent != want {
		t.Errorf("SkillPointsSpent = %v, the API says the player spent %v", spent, want)
	}
	req, err := character.Build(cfg, character.SkillLevels(user.Skills))
	if err != nil {
		t.Fatal(err)
	}
	if req.PlayerLevel > user.Leveling.Level {
		t.Errorf("Build needs player level %d, but the player is level %d", req.PlayerLevel, user.Leveling.Level)
	}

	// Attack's levels cost 1, 2, 3, ... points, so 3 -> 5 costs 4 + 5.
	if got, err := character.SkillCost(cfg, "attack", 3, 5); err != nil || got != 9 {
		t.Errorf("SkillCost(attack, 3, 5) = %v, %v; want 9", got, err)
	}
	if got, _ := character.SkillCost(cfg, "attack", 5, 3); got != -9 {
		t.Errorf("SkillCost(attack, 5, 3) = %v; want -9", got)
	}
	if _, err := character.SkillCost(cfg, "attack", 0, 99); !errors.Is(err, analytics.ErrUnknownSkill) {
		t.Errorf("missing level error = %v", err)
	}
	if got, err := character.MaxSkillLevel(cfg, "attack"); err != nil || got != 10 {
		t.Errorf("MaxSkillLevel(attack) = %v, %v; want 10", got, err)
	}
}

func TestGradeEquipment(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	eq := fixture[warera.CurrentEquipment](t, "inventory.fetchCurrentEquipment")
	grades, err := character.GradeEquipment(cfg, &eq)
	if err != nil {
		t.Fatal(err)
	}
	if len(grades) != 6 || grades[0].Slot != "weapon" || grades[0].ItemCode != "tank" {
		t.Fatalf("grades = %+v; want 6 slots starting with the tank", grades)
	}
	// The tank rolled attack 163 in [141, 170] and critical chance 29 in [26, 35].
	want := []character.StatRoll{
		{Stat: "attack", Value: 163, Min: 141, Max: 170, Quality: 22.0 / 29},
		{Stat: "criticalChance", Value: 29, Min: 26, Max: 35, Quality: 3.0 / 9},
	}
	if got := grades[0].Rolls; len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("tank rolls = %+v, want %+v", got, want)
	}
	// Equipment rolls one stat; the other three zero-valued stats are ignored.
	for _, g := range grades[1:] {
		if len(g.Rolls) != 1 || g.Rolls[0].Quality < 0 || g.Rolls[0].Quality > 1 {
			t.Errorf("%s rolls = %+v; want one stat within its range", g.Slot, g.Rolls)
		}
	}
}
