package market_test

import (
	"fmt"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/market"
	"github.com/onprem/warera-go/waretest"
)

// must stands in for error handling, keeping the examples short.
func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func ExampleBuyCost() {
	book := must(waretest.Fixture[*warera.TopOrders]("tradingOrder.getTopOrders")) // client.Trading.GetTopOrders(ctx, "ammo", 10)

	fill := market.BuyCost(book, 1000)
	fmt.Printf("1000 ammo: %d available for %.2f, %.4f each\n", fill.Quantity, fill.Total, fill.AveragePrice)
	bid, ask, _ := market.Spread(book)
	fmt.Printf("spread %.3f - %.3f\n", bid, ask)
	// Output:
	// 1000 ammo: 1000 available for 799.22, 0.7992 each
	// spread 0.796 - 0.799
}
