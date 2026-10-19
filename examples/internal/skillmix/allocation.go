// Package skillmix summarizes a player's regular skill-point allocation for
// the runnable examples. It is an example helper, not a game rule in the SDK.
package skillmix

import "github.com/onprem/warera-go"

// Mode is a heuristic label based on which group has more allocated points.
type Mode string

const (
	ModeWar      Mode = "war"
	ModeEco      Mode = "eco"
	ModeBalanced Mode = "balanced"
	ModeUnknown  Mode = "unknown"
)

// Allocation contains the regular skill points grouped into war and economy
// skills, their shares of classified points, and separate prestige levels.
type Allocation struct {
	WarPoints           int
	EcoPoints           int
	ClassifiedPoints    int
	WarPercent          float64
	EcoPercent          float64
	WarPrestigeLevels   int
	EcoPrestigeLevels   int
	PrestigeLevelTotals int
}

// Analyze converts skill levels into the regular skill points spent in each
// group. The cost of level N is N points, so raising a skill to level N costs
// N*(N+1)/2 points total. Percentages use only classified regular points;
// prestige levels are kept separate because they use a separate point pool.
//
// War skills here are attack, precision, critical chance and damage, armor,
// dodge, health, hunger, loot chance, and stamina (used for travel to a battle).
// Economy skills are energy, entrepreneurship, production, companies, and
// management. This grouping is an example convention for exploration.
func Analyze(skills warera.UserSkills) Allocation {
	warPoints := sumPoints(
		skills.Attack.Level,
		skills.Precision.Level,
		skills.CriticalChance.Level,
		skills.CriticalDamages.Level,
		skills.Armor.Level,
		skills.Dodge.Level,
		skills.Health.Level,
		skills.Hunger.Level,
		skills.LootChance.Level,
		skills.Stamina.Level,
	)
	ecoPoints := sumPoints(
		skills.Energy.Level,
		skills.Entrepreneurship.Level,
		skills.Production.Level,
		skills.Companies.Level,
		skills.Management.Level,
	)
	warPrestige := sumInts(
		skills.Attack.Prestige,
		skills.Precision.Prestige,
		skills.CriticalChance.Prestige,
		skills.CriticalDamages.Prestige,
		skills.Armor.Prestige,
		skills.Dodge.Prestige,
		skills.Health.Prestige,
		skills.Hunger.Prestige,
		skills.LootChance.Prestige,
		skills.Stamina.Prestige,
	)
	ecoPrestige := sumInts(
		skills.Energy.Prestige,
		skills.Entrepreneurship.Prestige,
		skills.Production.Prestige,
		skills.Companies.Prestige,
		skills.Management.Prestige,
	)
	classified := warPoints + ecoPoints

	allocation := Allocation{
		WarPoints:           warPoints,
		EcoPoints:           ecoPoints,
		ClassifiedPoints:    classified,
		WarPrestigeLevels:   warPrestige,
		EcoPrestigeLevels:   ecoPrestige,
		PrestigeLevelTotals: warPrestige + ecoPrestige,
	}
	if classified > 0 {
		allocation.WarPercent = float64(warPoints) / float64(classified) * 100
		allocation.EcoPercent = float64(ecoPoints) / float64(classified) * 100
	}
	return allocation
}

// Mode labels the larger regular skill-point group; ties are balanced.
func (a Allocation) Mode() Mode {
	switch {
	case a.ClassifiedPoints == 0:
		return ModeUnknown
	case a.WarPoints > a.EcoPoints:
		return ModeWar
	case a.EcoPoints > a.WarPoints:
		return ModeEco
	default:
		return ModeBalanced
	}
}

func sumPoints(levels ...int) int {
	total := 0
	for _, level := range levels {
		if level > 0 {
			total += level * (level + 1) / 2
		}
	}
	return total
}

func sumInts(values ...int) int {
	total := 0
	for _, value := range values {
		if value > 0 {
			total += value
		}
	}
	return total
}
