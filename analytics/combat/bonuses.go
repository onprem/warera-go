package combat

import (
	"fmt"
	"slices"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
	"github.com/onprem/warera-go/analytics/upgrades"
)

// Bonuses are the battle conditions that change a player's attack. Each
// true field adds the percentage the game config gives for it.
type Bonuses struct {
	// CountryOrder and MUOrder are the country's and the military unit's
	// orders to fight in this battle: CountryOrderBonusPercent and
	// MuOrderBonusPercent.
	CountryOrder bool
	MUOrder      bool
	// Patriotic is fighting for the player's own country:
	// PatrioticBonusPercent.
	Patriotic bool
	// Alliance is fighting for an ally: AllianceDamagesBonusPercent.
	Alliance bool
	// Enemy is fighting the country's sworn enemy: EnemyDamagesBonusPercent.
	Enemy bool
	// NotLinkedToCapital, LostAttackingRegion, and OccupyingYourRegions are
	// the maluses RegionNotLinkedToCapitalMalusPercent,
	// LostAttackingRegionMalusPercent, and OccupyingYourRegionsMalusPercent.
	NotLinkedToCapital   bool
	LostAttackingRegion  bool
	OccupyingYourRegions bool
	// HeadquartersLevel is the level of the military unit's headquarters,
	// which gives its "attackBonus" stat.
	HeadquartersLevel int
	// BaseLevel and BunkerLevel are the levels of the military base and the
	// bunker in the battle's region, which give their "attackBonus" and
	// "defenseBonus" stats.
	BaseLevel   int
	BunkerLevel int
	// Extra is any other bonus or malus in percent, such as resistance,
	// government type, or a defensive pact, which the game config does not
	// publish.
	Extra float64
}

// BattlePercent returns the sum of the battle bonuses and maluses in
// percent, at least -100, to use as [character.Loadout.BattlePercent].
//
// Game config: every percentage comes from the game config. That the
// bonuses add up rather than multiply is the community model.
func BattlePercent(cfg *warera.GameConfig, b Bonuses) (float64, error) {
	c := cfg.Battle
	var total float64
	for _, t := range []struct {
		on      bool
		percent float64
	}{
		{b.CountryOrder, c.CountryOrderBonusPercent},
		{b.MUOrder, c.MuOrderBonusPercent},
		{b.Patriotic, c.PatrioticBonusPercent},
		{b.Alliance, c.AllianceDamagesBonusPercent},
		{b.Enemy, c.EnemyDamagesBonusPercent},
		{b.NotLinkedToCapital, c.RegionNotLinkedToCapitalMalusPercent},
		{b.LostAttackingRegion, c.LostAttackingRegionMalusPercent},
		{b.OccupyingYourRegions, c.OccupyingYourRegionsMalusPercent},
	} {
		if t.on {
			total += t.percent
		}
	}
	for _, u := range []struct {
		upgrade warera.UpgradeType
		level   int
		stat    string
	}{
		{warera.UpgradeTypeHeadquarters, b.HeadquartersLevel, "attackBonus"},
		{warera.UpgradeTypeBase, b.BaseLevel, "attackBonus"},
		{warera.UpgradeTypeBunker, b.BunkerLevel, "defenseBonus"},
	} {
		if u.level <= 0 {
			continue
		}
		levels, err := upgrades.Levels(cfg, u.upgrade)
		if err != nil {
			return 0, err
		}
		i := slices.IndexFunc(levels, func(l warera.GameConfigUpgradeLevel) bool { return l.Level == u.level })
		if i < 0 {
			return 0, fmt.Errorf("%w: %q has no level %d", analytics.ErrUnknownUpgrade, u.upgrade, u.level)
		}
		total += levels[i].Stats[u.stat]
	}
	return max(-100, total+b.Extra), nil
}
