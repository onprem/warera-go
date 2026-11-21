package combat

import (
	"errors"
	"fmt"
	"math"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
	"github.com/onprem/warera-go/analytics/character"
)

// Assumptions are the parts of the fight model the game does not publish.
// [DefaultAssumptions] returns the values the Warera Simulator uses.
type Assumptions struct {
	// MissDamage is the share of attack a missed swing still deals.
	MissDamage float64
	// ItemUses is how many swings a weapon, or how many undodged swings an
	// armor piece, lasts before it breaks: one durability point per use,
	// from a MaxState of 100.
	ItemUses float64
}

// DefaultAssumptions returns a miss dealing half of attack and items lasting
// 100 uses.
func DefaultAssumptions() Assumptions {
	return Assumptions{MissDamage: 0.5, ItemUses: 100}
}

// Swing is the expected outcome of one attack.
type Swing struct {
	// Damage is the expected damage, averaged over misses, hits, and
	// critical hits.
	Damage float64
	// HitHealth is the health one swing costs when not dodged.
	HitHealth float64
	// Health is the expected health one swing costs, counting dodges.
	Health float64
}

// ExpectedSwing returns the expected outcome of one swing with the given
// stats. A swing hits with Precision percent chance and otherwise deals
// MissDamage of attack. A hit is critical with CriticalChance percent chance
// and then deals CriticalDamages percent more. Each swing costs the
// game config's HealthCost reduced by ArmorReduction, rounded down to 0.1,
// and costs nothing when dodged, with DodgeChance.
//
// Model: the API does not publish damage per swing. HealthCost is game
// config; the rest follows the Warera Simulator.
func ExpectedSwing(cfg *warera.GameConfig, s character.Stats, a Assumptions) Swing {
	hit := s.Precision / 100
	crit := s.CriticalChance / 100
	damage := s.Attack * ((1-hit)*a.MissDamage + hit*(1-crit) + hit*crit*(1+s.CriticalDamages/100))
	// The small epsilon keeps 10·(1-0.68) at 3.2 despite float rounding.
	hitHealth := math.Floor(cfg.Battle.HealthCost*(1-s.ArmorReduction)*10+1e-9) / 10
	return Swing{Damage: damage, HitHealth: hitHealth, Health: hitHealth * (1 - s.DodgeChance)}
}

// Fight is the expected result of fighting until health and food run out.
type Fight struct {
	Hours float64
	// Swings is the number of attacks, each using one ammo when ammo is
	// loaded.
	Swings float64
	// Hits is the swings that did not miss, which can drop cases.
	Hits float64
	// Undodged is the swings that cost health and armor durability.
	Undodged float64
	Damage   float64
	// Health is the health the fight can spend: the bar, food, and what
	// both regenerate during Hours.
	Health float64
	Food   float64
}

// ErrNoHealthCost is returned when a swing costs no health, so a fight
// never ends.
var ErrNoHealthCost = errors.New("combat: swings cost no health")

// Window returns the expected fight of a player who starts with full health
// and hunger, eats one food per hunger point, and fights for the given
// hours. Zero hours is a burst: the health bar plus a full stomach. Health
// and hunger each regenerate their maximum divided by the game config's
// RegenDividedBy per hour, and regenerated hunger counts in whole foods.
// Food is an item code, such as "cookedFish", that heals its
// "healthRegenPercent" of maximum health; empty means no food.
//
// Model: the regeneration and healing rates are game config; spending all
// of it in one window follows the Warera Simulator.
func Window(cfg *warera.GameConfig, s character.Stats, food string, hours float64, a Assumptions) (Fight, error) {
	var perFood float64
	if food != "" {
		item, ok := cfg.Items[food]
		if !ok {
			return Fight{}, fmt.Errorf("%w %q", analytics.ErrUnknownItem, food)
		}
		perFood = item.FlatStats["healthRegenPercent"] / 100 * s.Health
	}
	swing := ExpectedSwing(cfg, s, a)
	if swing.Health <= 0 {
		return Fight{}, ErrNoHealthCost
	}
	regen := hours / cfg.User.RegenDividedBy
	f := Fight{Hours: hours, Health: s.Health * (1 + regen)}
	if perFood > 0 {
		f.Food = s.Hunger + math.Floor(s.Hunger*regen)
		f.Health += f.Food * perFood
	}
	f.Swings = math.Ceil(f.Health / swing.Health)
	f.Hits = f.Swings * s.Precision / 100
	f.Undodged = f.Swings * (1 - s.DodgeChance)
	f.Damage = f.Swings * swing.Damage
	return f, nil
}

// Prices are the gold values a fight uses up.
type Prices struct {
	// Weapon is the value of the weapon, and Armor of all armor pieces
	// together. [CraftCost] estimates them from material prices.
	Weapon, Armor float64
	// Ammo and Food are the price of one unit, zero when none is used.
	Ammo, Food float64
	// Boosts is the value of boosts taken for this fight, such as a pill.
	Boosts float64
}

// Cost is the gold a fight uses up.
type Cost struct {
	// Wear is the share of the weapon and armor value the fight wears out.
	Wear  float64
	Ammo  float64
	Food  float64
	Total float64
}

// CostOf returns what a fight uses up: one ammo per swing, its food, its
// boosts, and wear of 1/ItemUses of the weapon per swing and of the armor
// per undodged swing.
//
// Model: wear per swing follows the Warera Simulator.
func CostOf(f Fight, p Prices, a Assumptions) Cost {
	c := Cost{
		Wear: f.Swings/a.ItemUses*p.Weapon + f.Undodged/a.ItemUses*p.Armor,
		Ammo: f.Swings * p.Ammo,
		Food: f.Food * p.Food,
	}
	c.Total = c.Wear + c.Ammo + c.Food + p.Boosts
	return c
}

// Rewards are the gold values of what a fight earns.
type Rewards struct {
	// Case1 and Case2 are the prices of the two case types.
	Case1, Case2 float64
	// BountyPer1k is the bounty paid per 1,000 damage.
	BountyPer1k float64
}

// Income is what a fight earns.
type Income struct {
	// Case1 and Case2 are the expected number of cases dropped.
	Case1, Case2 float64
	Cases        float64
	Bounty       float64
	Total        float64
}

// IncomeOf returns what a fight earns with the given loot chance percent,
// [character.Stats.LootChance]: each hit drops a case1 with LootChance
// percent chance and a case2 with LootChance/100 percent chance, and the
// bounty pays per 1,000 damage. Battle loot pools are not included.
//
// Model: case drop rates follow the Warera Simulator.
func IncomeOf(f Fight, lootChance float64, r Rewards) Income {
	in := Income{Case1: f.Hits * lootChance / 100, Case2: f.Hits * lootChance / 10000}
	in.Cases = in.Case1*r.Case1 + in.Case2*r.Case2
	in.Bounty = f.Damage / 1000 * r.BountyPer1k
	in.Total = in.Cases + in.Bounty
	return in
}

// PerThousand returns gold per 1,000 damage of a fight, such as the cost
// of 1k damage from a [Cost] total.
func (f Fight) PerThousand(gold float64) float64 {
	if f.Damage == 0 {
		return 0
	}
	return gold / f.Damage * 1000
}

// craftScraps is the scraps crafting an item of each rarity takes, as seen
// in the game; the game config does not publish it.
var craftScraps = map[string]float64{
	"common": 6, "uncommon": 18, "rare": 54, "epic": 162, "legendary": 486, "mythic": 1460,
}

// CraftCost estimates the gold value of a weapon or armor piece of a rarity
// from what crafting it takes: the steel of the game config's MergingCost
// and the scraps seen in the game, at prices["steel"] and prices["scraps"].
// It values gear that has no market price.
//
// Model: the scraps counts are observed, and reading MergingCost as the
// steel crafting takes follows the Warera Simulator.
func CraftCost(cfg *warera.GameConfig, rarity string, prices warera.ItemPrices) (float64, error) {
	m := cfg.MergingCost
	steel, ok := map[string]float64{
		"common": m.Common, "uncommon": m.Uncommon, "rare": m.Rare,
		"epic": m.Epic, "legendary": m.Legendary, "mythic": m.Mythic,
	}[rarity]
	if !ok {
		return 0, fmt.Errorf("combat: unknown rarity %q", rarity)
	}
	return steel*prices["steel"] + craftScraps[rarity]*prices["scraps"], nil
}
