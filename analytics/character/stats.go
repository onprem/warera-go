package character

import (
	"fmt"
	"math"
	"strconv"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
)

// percentCap is the value above which a percentage skill, such as precision
// or critical chance, overflows into the skill its SkillOverflow names.
const percentCap = 100

// Loadout is everything that sets a player's fight stats: skill levels,
// strategic bonuses, equipment rolls, and attack modifiers. [LoadoutOf] reads
// a player's current loadout from API records; change its fields to try a
// different build or kit.
type Loadout struct {
	// Levels are regular skill levels keyed by skill name, and Prestige the
	// levels bought with prestige points on top. Stats use their sum; a sum
	// above the highest configured level continues the configured step per
	// level.
	Levels   map[string]int
	Prestige map[string]int
	// Strategic is the flat bonus per skill from strategic resources, as in
	// [warera.UserSkill.Strategic].
	Strategic map[string]float64
	// Equipment is the rolled stats of the weapon and armor, summed per stat
	// name, such as "attack", "criticalChance", "armor", or "precision".
	Equipment map[string]float64
	// AmmoPercent is the attack bonus of the loaded ammo, the item's
	// "percentAttack" flat stat: 10 for lightAmmo, 20 for ammo, 40 for
	// heavyAmmo.
	AmmoPercent float64
	// BuffPercent is the attack bonus of an active boost: a pill (cocain)
	// gives its "percentAttack", 60, for its "buffDurationHours".
	BuffPercent float64
	// DebuffPercent is the attack lost to a boost's after-effect: 60 for
	// the pill's "debuffDurationHours" after it ends.
	DebuffPercent float64
	// MilitaryRankPercent is the attack bonus of the player's military rank,
	// as in [warera.UserSkill.MilitaryRankPercent].
	MilitaryRankPercent float64
	// BattlePercent is the sum of battle bonuses and maluses, such as
	// orders and buildings. The API reports stats without it; the combat
	// package computes it for a battle.
	BattlePercent float64
}

// LoadoutOf returns a player's current loadout from their skills, as in
// [warera.User.Skills], and equipment, as returned by
// Inventory.GetCurrentEquipment. Equipment may be nil for no equipment.
func LoadoutOf(skills warera.UserSkills, equipment *warera.CurrentEquipment) Loadout {
	l := Loadout{
		Levels:              map[string]int{},
		Prestige:            map[string]int{},
		Strategic:           map[string]float64{},
		Equipment:           map[string]float64{},
		AmmoPercent:         skills.Attack.AmmoPercent,
		BuffPercent:         skills.Attack.BuffsPercent,
		DebuffPercent:       skills.Attack.DebuffsPercent,
		MilitaryRankPercent: skills.Attack.MilitaryRankPercent,
	}
	for name, s := range userSkills(skills) {
		l.Levels[name] = s.Level
		if s.Prestige != 0 {
			l.Prestige[name] = s.Prestige
		}
		if s.Strategic != 0 {
			l.Strategic[name] = s.Strategic
		}
	}
	if equipment == nil {
		return l
	}
	if w := equipment.Weapon; w != nil {
		l.Equipment["attack"] += w.Skills.Attack
		l.Equipment["criticalChance"] += w.Skills.CriticalChance
	}
	for _, item := range []warera.EquipmentItem{equipment.Helmet, equipment.Chest, equipment.Pants, equipment.Boots, equipment.Gloves} {
		l.Equipment["armor"] += item.Skills.Armor
		l.Equipment["criticalDamages"] += item.Skills.CriticalDamages
		l.Equipment["dodge"] += item.Skills.Dodge
		l.Equipment["precision"] += item.Skills.Precision
	}
	return l
}

func userSkills(skills warera.UserSkills) map[string]warera.UserSkill {
	return map[string]warera.UserSkill{
		"armor":           skills.Armor,
		"attack":          skills.Attack,
		"companies":       skills.Companies,
		"criticalChance":  skills.CriticalChance,
		"criticalDamages": skills.CriticalDamages,
		"dodge":           skills.Dodge,
		"energy":          skills.Energy,
		"health":          skills.Health,
		"hunger":          skills.Hunger,
		"lootChance":      skills.LootChance,
		"management":      skills.Management,
		"precision":       skills.Precision,
		"production":      skills.Production,
		"stamina":         skills.Stamina,
	}
}

// Stats are a player's fight stats, the totals the API reports in each
// [warera.UserSkill.Total] and TotalAfterSoftCap.
type Stats struct {
	// Attack is the damage of a normal hit: base attack, weapon, and
	// precision overflow, times every attack modifier. The API reports it
	// rounded.
	Attack float64
	// Precision and CriticalChance are percentages capped at 100; the excess
	// adds to attack and critical damage.
	Precision      float64
	CriticalChance float64
	// CriticalDamages is the extra damage of a critical hit, in percent.
	CriticalDamages float64
	// Armor and Dodge are the raw totals. ArmorReduction and DodgeChance are
	// their effect after the soft cap, as fractions: x / (x + softCap),
	// rounded to 0.01. The API reports them as TotalAfterSoftCap percent.
	Armor          float64
	ArmorReduction float64
	Dodge          float64
	DodgeChance    float64
	Health         float64
	Hunger         float64
	LootChance     float64
}

// Fight stat names, the keys of [warera.GameConfig.Skills].
var fightSkills = []string{"attack", "precision", "criticalChance", "criticalDamages", "armor", "dodge", "health", "hunger", "lootChance"}

// StatsOf returns the fight stats of a loadout. Each skill's total is its
// configured value at its level, plus strategic and equipment bonuses.
// Precision and critical chance above 100 overflow into attack and critical
// damage at their SkillOverflowValue per point. Attack is then multiplied
// by each modifier in turn: battle bonuses, ammo, boost, after-effect, and
// military rank.
//
// Verified: for 31 players, [LoadoutOf] their records gives every total the
// API reports, and attack within one point; the API appears to round
// military rank percentages before reporting them.
func StatsOf(cfg *warera.GameConfig, l Loadout) (Stats, error) {
	total := map[string]float64{}
	for _, name := range fightSkills {
		v, err := skillValue(cfg, name, l.Levels[name]+l.Prestige[name])
		if err != nil {
			return Stats{}, err
		}
		total[name] = v + l.Strategic[name] + l.Equipment[name]
	}
	for _, name := range fightSkills {
		s := cfg.Skills[name]
		if s.SkillOverflow == "" || total[name] <= percentCap {
			continue
		}
		total[s.SkillOverflow] += (total[name] - percentCap) * s.SkillOverflowValue
		total[name] = percentCap
	}
	attack := total["attack"]
	for _, percent := range []float64{l.BattlePercent, l.AmmoPercent, l.BuffPercent, -l.DebuffPercent, l.MilitaryRankPercent} {
		attack *= 1 + percent/100
	}
	return Stats{
		Attack:          attack,
		Precision:       total["precision"],
		CriticalChance:  total["criticalChance"],
		CriticalDamages: total["criticalDamages"],
		Armor:           total["armor"],
		ArmorReduction:  afterSoftCap(total["armor"], cfg.Skills["armor"].SoftCap),
		Dodge:           total["dodge"],
		DodgeChance:     afterSoftCap(total["dodge"], cfg.Skills["dodge"].SoftCap),
		Health:          total["health"],
		Hunger:          total["hunger"],
		LootChance:      total["lootChance"],
	}, nil
}

// skillValue returns a skill's value at a level, continuing the configured
// step per level above the highest configured level.
func skillValue(cfg *warera.GameConfig, skill string, level int) (float64, error) {
	s, ok := cfg.Skills[skill]
	if !ok {
		return 0, fmt.Errorf("%w %q", analytics.ErrUnknownSkill, skill)
	}
	if l, ok := s.Levels[strconv.Itoa(level)]; ok {
		return l.Value, nil
	}
	highest, err := MaxSkillLevel(cfg, skill)
	if err != nil {
		return 0, err
	}
	top, below := s.Levels[strconv.Itoa(highest)], s.Levels[strconv.Itoa(highest-1)]
	if level < 0 || highest == 0 {
		return 0, fmt.Errorf("%w: %q has no level %d", analytics.ErrUnknownSkill, skill, level)
	}
	return top.Value + float64(level-highest)*(top.Value-below.Value), nil
}

func afterSoftCap(x, softCap float64) float64 {
	if x <= 0 {
		return 0
	}
	return math.Round(x/(x+softCap)*100) / 100
}
