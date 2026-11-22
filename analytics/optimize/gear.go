package optimize

import (
	"cmp"
	"maps"
	"math"
	"slices"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/character"
)

// rarities are the item rarities, worst first.
var rarities = []string{"common", "uncommon", "rare", "epic", "legendary", "mythic"}

// GearChange is swapping one equipped item for the same kind of item one
// rarity better or worse.
type GearChange struct {
	Slot     string
	From, To string
	// Stats is the roll assumed for To: each stat at the same place in its
	// range as the current item's roll, rounded.
	Stats map[string]float64
	// Gain is how much the swap changes the metric, and Cost how much
	// more To costs than From; both are negative for a downgrade.
	Gain float64
	Cost float64
	// Loadout is the loadout after the swap, to weigh it further, such as
	// what a fight with it costs.
	Loadout character.Loadout
}

// GearChanges returns, for each equipped item, what swapping it for the
// same kind of item one rarity better and one rarity worse would gain and
// cost. prices maps item codes to their market price, such as
// MarketStats.GetEquipmentAveragePrice for each code; swaps whose items
// have no price are left out. The loadout must be [character.LoadoutOf]
// the player with this equipment.
//
// Upgrades come first, most gain per gold first; downgrades follow, least
// loss per gold saved first. Cost compares purchase prices only; better
// gear also wears out more value per fight, which
// [github.com/onprem/warera-go/analytics/combat.CostOf] prices.
func GearChanges(cfg *warera.GameConfig, l character.Loadout, equipment *warera.CurrentEquipment, prices map[string]float64, m Metric) ([]GearChange, error) {
	base, err := score(cfg, l, m)
	if err != nil {
		return nil, err
	}
	grades, err := character.GradeEquipment(cfg, equipment)
	if err != nil {
		return nil, err
	}
	var changes []GearChange
	for _, g := range grades {
		from, ok := prices[g.ItemCode]
		if !ok {
			continue
		}
		r := slices.Index(rarities, g.Rarity)
		for _, step := range []int{1, -1} {
			if r+step < 0 || r+step >= len(rarities) {
				continue
			}
			code, item, ok := sameKind(cfg, cfg.Items[g.ItemCode].Usage, rarities[r+step])
			to, priced := prices[code]
			if !ok || !priced {
				continue
			}
			c := GearChange{Slot: g.Slot, From: g.ItemCode, To: code, Stats: map[string]float64{}, Cost: to - from}
			swapped := l
			swapped.Equipment = maps.Clone(l.Equipment)
			for _, roll := range g.Rolls {
				swapped.Equipment[roll.Stat] -= roll.Value
				if bounds := item.DynamicStats[roll.Stat]; len(bounds) == 2 {
					c.Stats[roll.Stat] = math.Round(bounds[0] + roll.Quality*(bounds[1]-bounds[0]))
					swapped.Equipment[roll.Stat] += c.Stats[roll.Stat]
				}
			}
			s, err := score(cfg, swapped, m)
			if err != nil {
				return nil, err
			}
			c.Gain = s - base
			c.Loadout = swapped
			changes = append(changes, c)
		}
	}
	slices.SortStableFunc(changes, func(a, b GearChange) int {
		if up, other := a.Cost > 0, b.Cost > 0; up != other {
			if up {
				return -1
			}
			return 1
		}
		// For downgrades both are negative, so the ratio is loss per gold
		// saved and the smaller comes first.
		if a.Cost > 0 {
			return cmp.Compare(b.Gain/b.Cost, a.Gain/a.Cost)
		}
		return cmp.Compare(a.Gain/a.Cost, b.Gain/b.Cost)
	})
	return changes, nil
}

// sameKind returns the item with a usage, such as "weapon" or "helmet", and
// a rarity.
func sameKind(cfg *warera.GameConfig, usage, rarity string) (string, warera.GameConfigItem, bool) {
	for code, item := range cfg.Items {
		if item.Usage == usage && item.Rarity == rarity && len(item.DynamicStats) > 0 {
			return code, item, true
		}
	}
	return "", warera.GameConfigItem{}, false
}
