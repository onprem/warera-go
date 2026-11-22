package character_test

import (
	"encoding/json"
	"math"
	"os"
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/character"
)

// testdata/players.json holds the skills and equipment of 31 players recorded
// from the live API, including the fixture player.
func TestStatsMatchTheAPI(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	data, err := os.ReadFile("testdata/players.json")
	if err != nil {
		t.Fatal(err)
	}
	var players []struct {
		Skills    warera.UserSkills       `json:"skills"`
		Equipment warera.CurrentEquipment `json:"equipment"`
	}
	if err := json.Unmarshal(data, &players); err != nil {
		t.Fatal(err)
	}
	if len(players) < 30 {
		t.Fatalf("decoded %d players", len(players))
	}
	for i, p := range players {
		s, err := character.StatsOf(cfg, character.LoadoutOf(p.Skills, &p.Equipment))
		if err != nil {
			t.Fatal(err)
		}
		k := p.Skills
		// The API rounds attack, and seemingly the military rank
		// percentage before applying it.
		if math.Abs(s.Attack-k.Attack.Total) >= 1 {
			t.Errorf("player %d: attack %v, API %v", i, s.Attack, k.Attack.Total)
		}
		for name, c := range map[string][2]float64{
			"precision":       {s.Precision, k.Precision.Total},
			"criticalChance":  {s.CriticalChance, k.CriticalChance.Total},
			"criticalDamages": {s.CriticalDamages, k.CriticalDamages.Total},
			"armor":           {s.Armor, k.Armor.Total},
			"armorAfterCap":   {math.Round(s.ArmorReduction * 100), deref(k.Armor.TotalAfterSoftCap)},
			"dodge":           {s.Dodge, k.Dodge.Total},
			"dodgeAfterCap":   {math.Round(s.DodgeChance * 100), deref(k.Dodge.TotalAfterSoftCap)},
			"health":          {s.Health, k.Health.Total},
			"hunger":          {s.Hunger, k.Hunger.Total},
			"lootChance":      {s.LootChance, k.LootChance.Total},
		} {
			if c[0] != c[1] {
				t.Errorf("player %d: %s %v, API %v", i, name, c[0], c[1])
			}
		}
	}
}

func deref(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

func TestStatsWhatIf(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	// Precision 50+10·5 = 100 plus 30 from gloves overflows 30 points into
	// 30·4 = 120 attack. Attack 100+10·25 = 350 + 50 weapon + 120 = 520,
	// then ×1.4 heavy ammo, ×1.6 pill, ×1.5 rank = 1747.2.
	l := character.Loadout{
		Levels:              map[string]int{"attack": 10, "precision": 10},
		Equipment:           map[string]float64{"attack": 50, "precision": 30},
		AmmoPercent:         40,
		BuffPercent:         60,
		MilitaryRankPercent: 50,
	}
	s, err := character.StatsOf(cfg, l)
	if err != nil {
		t.Fatal(err)
	}
	if want := 520 * 1.4 * 1.6 * 1.5; math.Abs(s.Attack-want) > 1e-9 || s.Precision != 100 {
		t.Errorf("attack %v precision %v, want %v and 100", s.Attack, s.Precision, want)
	}
	// Two prestige levels above the top level keep the step: attack 10+2
	// is 350 + 2·25.
	l = character.Loadout{Levels: map[string]int{"attack": 10}, Prestige: map[string]int{"attack": 2}}
	if s, _ := character.StatsOf(cfg, l); s.Attack != 400 {
		t.Errorf("attack at level 12 = %v, want 400", s.Attack)
	}
}

func TestMilitaryRankPercent(t *testing.T) {
	// Recorded from live players.
	for rank, want := range map[int]float64{3: 1, 4: 1.5, 9: 5, 13: 7.5, 25: 15, 40: 24, 89: 55, 92: 56.5, 93: 57.5, 101: 62.5, 108: 66.5} {
		if got := character.MilitaryRankPercent(rank); got != want {
			t.Errorf("MilitaryRankPercent(%d) = %v, want %v", rank, got, want)
		}
	}
	user := fixture[warera.User](t, "user.getUserById")
	if got := character.MilitaryRankPercent(user.MilitaryRank); got != user.Skills.Attack.MilitaryRankPercent {
		t.Errorf("rank %d: %v, API %v", user.MilitaryRank, got, user.Skills.Attack.MilitaryRankPercent)
	}
}
