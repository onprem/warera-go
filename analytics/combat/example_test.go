package combat_test

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/character"
	"github.com/onprem/warera-go/analytics/combat"
	"github.com/onprem/warera-go/waretest"
)

// must stands in for error handling, keeping the examples short.
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Estimate a player's 8-hour fight under a country order, with cooked fish,
// and what it nets at market prices.
func ExampleWindow() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig"))                   // client.GameConfig.Get(ctx)
	user := must(waretest.Fixture[warera.User]("user.getUserById"))                                 // client.Users.GetByID(ctx, id)
	equipment := must(waretest.Fixture[warera.CurrentEquipment]("inventory.fetchCurrentEquipment")) // client.Inventory.GetCurrentEquipment(ctx, id)
	prices := must(waretest.Fixture[warera.ItemPrices]("itemTrading.getPrices"))                    // client.Trading.GetPrices(ctx)

	loadout := character.LoadoutOf(user.Skills, &equipment)
	loadout.BattlePercent = must(combat.BattlePercent(cfg, combat.Bonuses{CountryOrder: true}))
	stats := must(character.StatsOf(cfg, loadout))

	// Value the gear at what crafting it takes.
	gearValue := func(code string) float64 { return must(combat.CraftCost(cfg, cfg.Items[code].Rarity, prices)) }
	var armor float64
	for _, piece := range []warera.EquipmentItem{equipment.Helmet, equipment.Chest, equipment.Pants, equipment.Boots, equipment.Gloves} {
		armor += gearValue(piece.Code)
	}

	a := combat.DefaultAssumptions()
	fight := must(combat.Window(cfg, stats, "cookedFish", 8, a))
	cost := combat.CostOf(fight, combat.Prices{
		Weapon: gearValue(equipment.Weapon.Code),
		Armor:  armor,
		Ammo:   prices[equipment.Ammo],
		Food:   prices["cookedFish"],
	}, a)
	income := combat.IncomeOf(fight, stats.LootChance, combat.Rewards{Case1: prices["case1"], Case2: prices["case2"]})

	fmt.Printf("attack %.0f, %.0f swings, %.0fk damage\n", stats.Attack, fight.Swings, fight.Damage/1000)
	fmt.Printf("cost %.0f, cases %.0f, net %.0f gold (%.2f per 1k)\n",
		cost.Total, income.Total, income.Total-cost.Total, fight.PerThousand(cost.Total-income.Total))
	// Output:
	// attack 383, 491 swings, 629k damage
	// cost 2743, cases 266, net -2477 gold (3.94 per 1k)
}

// See how much an 8-hour fight varies, over 1,000 seeded fights.
func ExampleSimulate() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig"))                   // client.GameConfig.Get(ctx)
	user := must(waretest.Fixture[warera.User]("user.getUserById"))                                 // client.Users.GetByID(ctx, id)
	equipment := must(waretest.Fixture[warera.CurrentEquipment]("inventory.fetchCurrentEquipment")) // client.Inventory.GetCurrentEquipment(ctx, id)
	stats := must(character.StatsOf(cfg, character.LoadoutOf(user.Skills, &equipment)))

	var damages []float64
	for seed := range uint64(1000) {
		f := must(combat.Simulate(cfg, stats, "cookedFish", 8, combat.DefaultAssumptions(), rand.New(rand.NewPCG(seed, 0))))
		damages = append(damages, f.Damage)
	}
	slices.Sort(damages)
	fmt.Printf("10%%: %.0fk, median: %.0fk, 90%%: %.0fk\n", damages[100]/1000, damages[500]/1000, damages[900]/1000)
	// Output: 10%: 534k, median: 572k, 90%: 610k
}
