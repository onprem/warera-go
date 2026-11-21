package company_test

import (
	"fmt"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/company"
	"github.com/onprem/warera-go/waretest"
)

// must stands in for error handling, keeping the examples short.
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Value a workforce's hourly output, such as the API key owner's from
// client.Me.GetWorkforceSummary(ctx).
func ExampleProfitPerHour() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig")) // client.GameConfig.Get(ctx)
	prices := must(waretest.Fixture[warera.ItemPrices]("itemTrading.getPrices"))  // client.Trading.GetPrices(ctx)
	summary := &warera.WorkforceSummary{
		ProductionByItem: map[string]float64{"iron": 30, "steel": 50},
		WageCost:         1.5,
	}

	p, _ := company.ProfitPerHour(cfg, prices, summary)
	fmt.Printf("revenue %.2f, inputs %.2f, wages %.2f: profit %.2f per hour\n", p.Revenue, p.InputCost, p.Wages, p.Profit)
	// Output: revenue 11.08, inputs 4.30, wages 1.50: profit 5.28 per hour
}
