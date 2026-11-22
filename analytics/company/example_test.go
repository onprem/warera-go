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

// Decide whether hiring a manager pays: the output their bonus adds,
// against their wage at the best offer's rate.
func ExampleExtraProduction() {
	cfg := must(waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig"))                             // client.GameConfig.Get(ctx)
	prices := must(waretest.Fixture[warera.ItemPrices]("itemTrading.getPrices"))                              // client.Trading.GetPrices(ctx)
	roster := must(waretest.Fixture[warera.CompanyRoster]("employee.getCompanyRoster"))                       // client.Employees.GetCompanyRoster(ctx, companyID)
	offers := must(waretest.Fixture[warera.CursorPage[warera.ManagerOffer]]("managerOffer.getManyPaginated")) // client.Managers.ListOffers(ctx, filter)

	item := roster.Hourly.ItemCode
	perPoint := prices[item] / cfg.Items[item].ProductionPoints // fish has no inputs
	rate := offers.Items[0].RateManagement
	for _, level := range []int{2, 6, 10} {
		gain := company.ExtraProduction(&roster, company.ManagerBonusPercent(level)) * perPoint
		wage := company.ManagerWage(rate, level)
		fmt.Printf("level %2d: +%.3f output, -%.3f wage, net %+.3f per hour\n", level, gain, wage, gain-wage)
	}
	// Output:
	// level  2: +0.083 output, -0.049 wage, net +0.034 per hour
	// level  6: +0.194 output, -0.114 wage, net +0.079 per hour
	// level 10: +0.305 output, -0.180 wage, net +0.125 per hour
}
