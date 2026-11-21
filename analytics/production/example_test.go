package production_test

import (
	"fmt"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/production"
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

	// Steel takes 10 production points plus 10 iron at 1 point each.
	steel, _ := production.Cost(cfg, "steel")
	fmt.Println("steel:", steel)
	// Output: steel: 20
}

func ExampleRawMaterials() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig")) // client.GameConfig.Get(ctx)

	materials, _ := production.RawMaterials(cfg, "cocain", 5)
	fmt.Println(materials)
	// Output: map[coca:1000]
}

// Rank what a company could make by what a production point earns.
func ExampleValues() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig")) // client.GameConfig.Get(ctx)
	prices := must(waretest.Fixture[warera.ItemPrices]("itemTrading.getPrices"))  // client.Trading.GetPrices(ctx)

	for _, v := range production.Values(cfg, prices)[:3] {
		fmt.Printf("%-9s %.4f per point (%.4f making the inputs too)\n", v.ItemCode, v.MoneyPerPoint, v.MoneyPerPointFromRaw)
	}
	// Output:
	// wood      0.1033 per point (0.1033 making the inputs too)
	// ammo      0.1015 per point (0.0996 making the inputs too)
	// limestone 0.0993 per point (0.0993 making the inputs too)
}
