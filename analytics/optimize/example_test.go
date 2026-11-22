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
