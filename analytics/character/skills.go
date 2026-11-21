package character

import (
	"fmt"
	"strconv"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
)

// SkillLevels returns a player's skill levels keyed by skill name, the names
// [warera.GameConfig.Skills] uses, such as "attack" or "criticalChance".
func SkillLevels(skills warera.UserSkills) map[string]int {
	return map[string]int{
		"armor":           skills.Armor.Level,
		"attack":          skills.Attack.Level,
		"companies":       skills.Companies.Level,
		"criticalChance":  skills.CriticalChance.Level,
		"criticalDamages": skills.CriticalDamages.Level,
		"dodge":           skills.Dodge.Level,
		"energy":          skills.Energy.Level,
		"health":          skills.Health.Level,
		"hunger":          skills.Hunger.Level,
		"lootChance":      skills.LootChance.Level,
		"management":      skills.Management.Level,
		"precision":       skills.Precision.Level,
		"production":      skills.Production.Level,
		"stamina":         skills.Stamina.Level,
	}
}

// SkillLevel returns the configuration of one level of a skill: its value,
// its cost, the total cost from level 0, and the player level that unlocks it.
func SkillLevel(cfg *warera.GameConfig, skill string, level int) (warera.GameConfigSkillLevel, error) {
	s, ok := cfg.Skills[skill]
	if !ok {
		return warera.GameConfigSkillLevel{}, fmt.Errorf("%w %q", analytics.ErrUnknownSkill, skill)
	}
	l, ok := s.Levels[strconv.Itoa(level)]
	if !ok {
		return warera.GameConfigSkillLevel{}, fmt.Errorf("%w: %q has no level %d", analytics.ErrUnknownSkill, skill, level)
	}
	return l, nil
}

// MaxSkillLevel returns the highest level of a skill.
func MaxSkillLevel(cfg *warera.GameConfig, skill string) (int, error) {
	s, ok := cfg.Skills[skill]
	if !ok {
		return 0, fmt.Errorf("%w %q", analytics.ErrUnknownSkill, skill)
	}
	highest := 0
	for key := range s.Levels {
		if level, err := strconv.Atoi(key); err == nil && level > highest {
			highest = level
		}
	}
	return highest, nil
}

// SkillCost returns the skill points needed to raise a skill from one level
// to another. Lowering a skill returns a negative cost, the points freed.
func SkillCost(cfg *warera.GameConfig, skill string, from, to int) (float64, error) {
	start, err := SkillLevel(cfg, skill, from)
	if err != nil {
		return 0, err
	}
	end, err := SkillLevel(cfg, skill, to)
	if err != nil {
		return 0, err
	}
	return end.TotalCost - start.TotalCost, nil
}

// BuildRequirements is what a set of skill levels needs.
type BuildRequirements struct {
	// SkillPoints is the regular skill points the levels cost from level 0.
	SkillPoints float64
	// PlayerLevel is the lowest player level at which every chosen skill
	// level is unlocked.
	PlayerLevel int
}

// Build returns what a build, skill levels keyed by skill name, needs. Pass
// [SkillLevels] of a player to see what their current build cost; its
// SkillPoints equals [warera.UserLeveling.SpentSkillPoints].
func Build(cfg *warera.GameConfig, levels map[string]int) (BuildRequirements, error) {
	var req BuildRequirements
	for skill, level := range levels {
		l, err := SkillLevel(cfg, skill, level)
		if err != nil {
			return BuildRequirements{}, err
		}
		req.SkillPoints += l.TotalCost
		if level > 0 && l.UnlockAtLevel > req.PlayerLevel {
			req.PlayerLevel = l.UnlockAtLevel
		}
	}
	return req, nil
}

// SkillPointsSpent returns the regular skill points a player's skills cost,
// which equals [warera.UserLeveling.SpentSkillPoints].
func SkillPointsSpent(cfg *warera.GameConfig, skills warera.UserSkills) (float64, error) {
	req, err := Build(cfg, SkillLevels(skills))
	return req.SkillPoints, err
}
