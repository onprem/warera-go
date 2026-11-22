package optimize_test

import (
	"fmt"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/character"
	"github.com/onprem/warera-go/analytics/combat"
	"github.com/onprem/warera-go/analytics/optimize"
	"github.com/onprem/warera-go/waretest"
)

// Plan the next player levels to raise 8-hour damage with cooked fish.
func ExampleRoadmap() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig"))                   // client.GameConfig.Get(ctx)
	user := must(waretest.Fixture[warera.User]("user.getUserById"))                                 // client.Users.GetByID(ctx, id)
	equipment := must(waretest.Fixture[warera.CurrentEquipment]("inventory.fetchCurrentEquipment")) // client.Inventory.GetCurrentEquipment(ctx, id)

	loadout := character.LoadoutOf(user.Skills, &equipment)
	metric := optimize.Damage(cfg, "cookedFish", 8, combat.DefaultAssumptions())
	steps, _, _ := optimize.Roadmap(cfg, loadout, user.Leveling.Level, user.Leveling.Level+4, metric)
	for _, s := range steps {
		fmt.Printf("level %d:", s.PlayerLevel)
		for _, c := range s.Raised {
			fmt.Printf(" %s %d", c.Skill, c.Level)
		}
		fmt.Printf(" -> %.0fk damage, %v points spare\n", s.Score/1000, s.SparePoints)
	}
	// Output:
	// level 37: hunger 5 -> 620k damage, 0 points spare
	// level 38: -> 620k damage, 4 points spare
	// level 39: health 5 -> 662k damage, 3 points spare
	// level 40: hunger 6 -> 715k damage, 1 points spare
}

// Find the gear swap that adds the most 8-hour damage per gold, and the
// cheapest one to sell down.
func ExampleGearChanges() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig"))                   // client.GameConfig.Get(ctx)
	user := must(waretest.Fixture[warera.User]("user.getUserById"))                                 // client.Users.GetByID(ctx, id)
	equipment := must(waretest.Fixture[warera.CurrentEquipment]("inventory.fetchCurrentEquipment")) // client.Inventory.GetCurrentEquipment(ctx, id)
	// client.MarketStats.GetEquipmentAveragePrice(ctx, code) for each code.
	prices := map[string]float64{
		"sniper": 56.69, "tank": 157.51, "jet": 515.74,
		"helmet4": 62.03, "helmet5": 132.35, "helmet6": 406.77,
		"chest4": 56.52, "chest5": 138.96, "chest6": 383.33,
		"gloves4": 51.28, "gloves5": 177.56, "gloves6": 461.08,
	}

	loadout := character.LoadoutOf(user.Skills, &equipment)
	metric := optimize.Damage(cfg, "cookedFish", 8, combat.DefaultAssumptions())
	changes := must(optimize.GearChanges(cfg, loadout, &equipment, prices, metric))
	for _, c := range changes {
		fmt.Printf("%-7s %-7s -> %-7s %+6.0fk damage for %+7.2f gold\n", c.Slot, c.From, c.To, c.Gain/1000, c.Cost)
	}
	// Output:
	// weapon  tank    -> jet       +232k damage for +358.23 gold
	// chest   chest5  -> chest6     +81k damage for +244.37 gold
	// gloves  gloves5 -> gloves6    +90k damage for +283.52 gold
	// helmet  helmet5 -> helmet6    +42k damage for +274.42 gold
	// helmet  helmet5 -> helmet4    -25k damage for  -70.32 gold
	// gloves  gloves5 -> gloves4    -49k damage for -126.28 gold
	// chest   chest5  -> chest4     -78k damage for  -82.44 gold
	// weapon  tank    -> sniper    -105k damage for -100.82 gold
}
