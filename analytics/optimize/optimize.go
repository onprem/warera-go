package optimize

import (
	"cmp"
	"maps"
	"slices"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/character"
	"github.com/onprem/warera-go/analytics/combat"
)

// PointsPerLevel is the skill points a player gets per player level, as
// [warera.UserLeveling.TotalSkillPoints] reports: 4 × level.
const PointsPerLevel = 4

// Metric scores fight stats; higher is better.
type Metric func(character.Stats) (float64, error)

// Damage returns a Metric of the expected damage of a fight, a
// [combat.Window] with the given food and hours.
func Damage(cfg *warera.GameConfig, food string, hours float64, a combat.Assumptions) Metric {
	return func(s character.Stats) (float64, error) {
		f, err := combat.Window(cfg, s, food, hours, a)
		return f.Damage, err
	}
}

// fightSkills are the skills that change fight stats.
var fightSkills = []string{"attack", "precision", "criticalChance", "criticalDamages", "armor", "dodge", "health", "hunger", "lootChance"}

// SkillChange is raising one skill by one regular level.
type SkillChange struct {
	Skill string
	// Level is the regular level after the change.
	Level int
	// Points is the skill points the level costs, and Gain how much it
	// raises the metric.
	Points       float64
	Gain         float64
	GainPerPoint float64
}

// SkillGains returns, for each fight skill, what raising it one regular
// level would gain, best GainPerPoint first. Skills at their highest
// level, or whose next level playerLevel does not unlock, are left out.
func SkillGains(cfg *warera.GameConfig, l character.Loadout, playerLevel int, m Metric) ([]SkillChange, error) {
	base, err := score(cfg, l, m)
	if err != nil {
		return nil, err
	}
	var changes []SkillChange
	for _, skill := range fightSkills {
		next := l.Levels[skill] + 1
		top, err := character.MaxSkillLevel(cfg, skill)
		if err != nil {
			return nil, err
		}
		if next > top {
			continue
		}
		level, err := character.SkillLevel(cfg, skill, next)
		if err != nil {
			return nil, err
		}
		if level.UnlockAtLevel > playerLevel {
			continue
		}
		raised := l
		raised.Levels = cloneLevels(l.Levels)
		raised.Levels[skill] = next
		s, err := score(cfg, raised, m)
		if err != nil {
			return nil, err
		}
		c := SkillChange{Skill: skill, Level: next, Points: level.Cost, Gain: s - base}
		if c.Points > 0 {
			c.GainPerPoint = c.Gain / c.Points
		}
		changes = append(changes, c)
	}
	slices.SortStableFunc(changes, func(a, b SkillChange) int { return cmp.Compare(b.GainPerPoint, a.GainPerPoint) })
	return changes, nil
}

func cloneLevels(levels map[string]int) map[string]int {
	c := map[string]int{}
	maps.Copy(c, levels)
	return c
}

func score(cfg *warera.GameConfig, l character.Loadout, m Metric) (float64, error) {
	s, err := character.StatsOf(cfg, l)
	if err != nil {
		return 0, err
	}
	return m(s)
}

// Step is what [Roadmap] does at one player level.
type Step struct {
	PlayerLevel int
	// Raised is the skill levels bought at this player level, in order.
	Raised []SkillChange
	// Score is the metric after them, and SparePoints the skill points
	// saved for later.
	Score       float64
	SparePoints float64
}

// Roadmap plans how to spend skill points from player level from to level
// to, starting from a loadout whose regular levels the player has at level
// from. At each level it gains [PointsPerLevel] points and buys the
// [SkillGains] change with the best gain per point while it can afford it;
// when it cannot, it saves the points for that change rather than buy a
// worse one. Changes that gain nothing are never bought. It returns one
// step per level after from, and the loadout at level to.
//
// Points the player had unspent at level from count:
// PointsPerLevel × from minus what the loadout's levels cost.
func Roadmap(cfg *warera.GameConfig, l character.Loadout, from, to int, m Metric) ([]Step, character.Loadout, error) {
	build, err := character.Build(cfg, l.Levels)
	if err != nil {
		return nil, l, err
	}
	l.Levels = cloneLevels(l.Levels)
	spare := float64(PointsPerLevel*from) - build.SkillPoints
	var steps []Step
	for level := from + 1; level <= to; level++ {
		spare += PointsPerLevel
		step := Step{PlayerLevel: level}
		for {
			changes, err := SkillGains(cfg, l, level, m)
			if err != nil {
				return nil, l, err
			}
			if len(changes) == 0 || changes[0].Gain <= 0 || changes[0].Points > spare {
				break
			}
			best := changes[0]
			l.Levels[best.Skill] = best.Level
			spare -= best.Points
			step.Raised = append(step.Raised, best)
		}
		if step.Score, err = score(cfg, l, m); err != nil {
			return nil, l, err
		}
		step.SparePoints = spare
		steps = append(steps, step)
	}
	return steps, l, nil
}
