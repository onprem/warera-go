package character_test

import (
	"fmt"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/character"
	"github.com/onprem/warera-go/waretest"
)

// must stands in for error handling, keeping the examples short.
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Check what a player's build cost and what raising a skill would cost.
func ExampleBuild() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig")) // client.GameConfig.Get(ctx)
	user := must(waretest.Fixture[warera.User]("user.getUserById"))               // client.Users.GetByID(ctx, id)

	build, _ := character.Build(cfg, character.SkillLevels(user.Skills))
	fmt.Printf("spent %v points; needs player level %d\n", build.SkillPoints, build.PlayerLevel)

	next, _ := character.SkillCost(cfg, "attack", user.Skills.Attack.Level, user.Skills.Attack.Level+1)
	fmt.Printf("attack %d -> %d costs %v points\n", user.Skills.Attack.Level, user.Skills.Attack.Level+1, next)
	// Output:
	// spent 143 points; needs player level 10
	// attack 7 -> 8 costs 8 points
}

func ExampleGradeEquipment() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig"))                   // client.GameConfig.Get(ctx)
	equipment := must(waretest.Fixture[warera.CurrentEquipment]("inventory.fetchCurrentEquipment")) // client.Inventory.GetCurrentEquipment(ctx, userID)

	grades, _ := character.GradeEquipment(cfg, &equipment)
	for _, g := range grades {
		for _, r := range g.Rolls {
			fmt.Printf("%-6s %-8s %-15s %3.0f in [%3.0f, %3.0f]: %3.0f%%\n", g.Slot, g.ItemCode, r.Stat, r.Value, r.Min, r.Max, 100*r.Quality)
		}
	}
	// Output:
	// weapon tank     attack          163 in [141, 170]:  76%
	// weapon tank     criticalChance   29 in [ 26,  35]:  33%
	// helmet helmet5  criticalDamages  97 in [ 91, 110]:  32%
	// chest  chest5   armor            42 in [ 36,  50]:  43%
	// pants  pants4   armor            24 in [ 21,  30]:  33%
	// boots  boots5   dodge            34 in [ 31,  40]:  33%
	// gloves gloves5  precision        33 in [ 31,  40]:  22%
}
