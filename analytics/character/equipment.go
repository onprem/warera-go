package character

import (
	"fmt"
	"slices"
	"strings"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
)

// StatRoll is how one rolled stat of an item compares with the range its item
// code can roll.
type StatRoll struct {
	Stat  string
	Value float64
	// Min and Max are the range from [warera.GameConfigItem.DynamicStats].
	Min, Max float64
	// Quality is where Value falls in the range, from 0 at Min to 1 at Max.
	Quality float64
}

// RollQuality compares an item's rolled stats, keyed by stat name, with the
// ranges its item code can roll. Stats the item code does not roll are
// ignored, and the result is sorted by stat name.
func RollQuality(cfg *warera.GameConfig, itemCode string, stats map[string]float64) ([]StatRoll, error) {
	item, ok := cfg.Items[itemCode]
	if !ok {
		return nil, fmt.Errorf("%w %q", analytics.ErrUnknownItem, itemCode)
	}
	var rolls []StatRoll
	for stat, bounds := range item.DynamicStats {
		value, ok := stats[stat]
		if !ok || len(bounds) != 2 {
			continue
		}
		r := StatRoll{Stat: stat, Value: value, Min: bounds[0], Max: bounds[1], Quality: 1}
		if r.Max > r.Min {
			r.Quality = (value - r.Min) / (r.Max - r.Min)
		}
		rolls = append(rolls, r)
	}
	slices.SortFunc(rolls, func(a, b StatRoll) int { return strings.Compare(a.Stat, b.Stat) })
	return rolls, nil
}

// EquipmentGrade grades one equipped item.
type EquipmentGrade struct {
	// Slot is "weapon", "helmet", "chest", "pants", "boots", or "gloves".
	Slot     string
	ItemCode string
	Rarity   string
	Rolls    []StatRoll
	// Condition is the item's remaining durability, State / MaxState.
	Condition float64
}

// GradeEquipment grades each item a player has equipped, as returned by
// Inventory.GetCurrentEquipment, in slot order. Empty slots are skipped.
func GradeEquipment(cfg *warera.GameConfig, equipment *warera.CurrentEquipment) ([]EquipmentGrade, error) {
	var grades []EquipmentGrade
	add := func(slot, code string, state, maxState float64, stats map[string]float64) error {
		if code == "" {
			return nil
		}
		rolls, err := RollQuality(cfg, code, stats)
		if err != nil {
			return err
		}
		g := EquipmentGrade{Slot: slot, ItemCode: code, Rarity: cfg.Items[code].Rarity, Rolls: rolls}
		if maxState > 0 {
			g.Condition = state / maxState
		}
		grades = append(grades, g)
		return nil
	}
	if w := equipment.Weapon; w != nil {
		if err := add("weapon", w.Code, w.State, w.MaxState, map[string]float64{
			"attack": w.Skills.Attack, "criticalChance": w.Skills.CriticalChance,
		}); err != nil {
			return nil, err
		}
	}
	for _, slot := range []struct {
		name string
		item warera.EquipmentItem
	}{
		{"helmet", equipment.Helmet}, {"chest", equipment.Chest}, {"pants", equipment.Pants},
		{"boots", equipment.Boots}, {"gloves", equipment.Gloves},
	} {
		s := slot.item.Skills
		if err := add(slot.name, slot.item.Code, slot.item.State, slot.item.MaxState, map[string]float64{
			"armor": s.Armor, "criticalDamages": s.CriticalDamages, "dodge": s.Dodge, "precision": s.Precision,
		}); err != nil {
			return nil, err
		}
	}
	return grades, nil
}
