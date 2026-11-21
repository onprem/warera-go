package upgrades_test

import (
	"fmt"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/upgrades"
	"github.com/onprem/warera-go/waretest"
)

// must stands in for error handling, keeping the examples short.
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func ExampleCost() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig")) // client.GameConfig.Get(ctx)

	cost, _ := upgrades.Cost(cfg, warera.UpgradeTypeBunker, 2, 4)
	fmt.Printf("bunker 2 -> 4: %v steel\n", cost.Steel)
	// Output: bunker 2 -> 4: 800 steel
}
