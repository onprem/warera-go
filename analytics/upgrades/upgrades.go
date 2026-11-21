package upgrades

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
)

// Levels returns the levels of an upgrade type, lowest first, with
// each level's cost, maintenance, and stats.
func Levels(cfg *warera.GameConfig, upgradeType warera.UpgradeType) ([]warera.GameConfigUpgradeLevel, error) {
	u := cfg.UpgradesConfig
	byType := map[warera.UpgradeType]map[string]warera.GameConfigUpgradeLevel{
		warera.UpgradeTypeAutomatedEngine:    u.AutomatedEngine.Levels,
		warera.UpgradeTypeBase:               u.Base.Levels,
		warera.UpgradeTypeBreakRoom:          u.BreakRoom.Levels,
		warera.UpgradeTypeBunker:             u.Bunker.Levels,
		warera.UpgradeTypeDormitories:        u.Dormitories.Levels,
		warera.UpgradeTypeHeadquarters:       u.Headquarters.Levels,
		warera.UpgradeTypePacificationCenter: u.PacificationCenter.Levels,
		warera.UpgradeTypeStorage:            u.Storage.Levels,
	}
	levels, ok := byType[upgradeType]
	if !ok || len(levels) == 0 {
		return nil, fmt.Errorf("%w %q", analytics.ErrUnknownUpgrade, upgradeType)
	}
	sorted := make([]warera.GameConfigUpgradeLevel, 0, len(levels))
	for _, l := range levels {
		sorted = append(sorted, l)
	}
	slices.SortFunc(sorted, func(a, b warera.GameConfigUpgradeLevel) int { return cmp.Compare(a.Level, b.Level) })
	return sorted, nil
}

// Costs is what raising an upgrade costs.
type Costs struct {
	Steel              float64
	ConstructionPoints float64
}

// Cost returns what raising an upgrade from one level to a higher one
// costs: the sum of the SteelCost and ConstructionPointsCost of each level
// above from, up to and including to. Region upgrades on the live API
// report an InvestedSteel equal to Cost from level 1, so the first
// level appears to be built with construction points rather than steel.
func Cost(cfg *warera.GameConfig, upgradeType warera.UpgradeType, from, to int) (Costs, error) {
	levels, err := Levels(cfg, upgradeType)
	if err != nil {
		return Costs{}, err
	}
	if from < 0 || to < from || to > levels[len(levels)-1].Level {
		return Costs{}, fmt.Errorf("%w: %q cannot go from level %d to %d", analytics.ErrUnknownUpgrade, upgradeType, from, to)
	}
	var c Costs
	for _, l := range levels {
		if l.Level > from && l.Level <= to {
			c.Steel += l.SteelCost
			c.ConstructionPoints += l.ConstructionPointsCost
		}
	}
	return c, nil
}
