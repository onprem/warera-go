package market_test

import (
	"math"
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/market"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestOrderBook(t *testing.T) {
	// Listed out of price order, as nothing guarantees otherwise.
	book := &warera.TopOrders{
		SellOrders: []warera.TopOrder{{Price: 1.2, Quantity: 10}, {Price: 1.0, Quantity: 5}, {Price: 1.1, Quantity: 3}},
		BuyOrders:  []warera.TopOrder{{Price: 0.8, Quantity: 4}, {Price: 0.9, Quantity: 2}},
	}
	// 5 at 1.0 + 3 at 1.1 + 2 at 1.2 = 10.7.
	if got := market.BuyCost(book, 10); got.Quantity != 10 || !got.Complete || !near(got.Total, 10.7) || !near(got.AveragePrice, 1.07) {
		t.Errorf("BuyCost(10) = %+v; want 10 units for 10.7", got)
	}
	// Only 6 units are bid for: 2 at 0.9 + 4 at 0.8 = 5.0.
	if got := market.SellProceeds(book, 10); got.Quantity != 6 || got.Complete || !near(got.Total, 5) {
		t.Errorf("SellProceeds(10) = %+v; want an incomplete fill of 6 for 5.0", got)
	}
	if bid, ask, ok := market.Spread(book); !ok || bid != 0.9 || ask != 1.0 {
		t.Errorf("Spread = %v, %v, %v; want 0.9, 1.0", bid, ask, ok)
	}
	if _, _, ok := market.Spread(&warera.TopOrders{}); ok {
		t.Error("Spread of an empty book must report !ok")
	}
}

func TestSummarizeTrading(t *testing.T) {
	days := []warera.ItemTradingDay{
		{AvgValue: 2, TotalValue: 20, TotalQuantity: 10},
		{AvgValue: 0, TotalValue: 0, TotalQuantity: 0}, // no trades: skipped
		{AvgValue: 1, TotalValue: 30, TotalQuantity: 30},
		{AvgValue: 2.5, TotalValue: 25, TotalQuantity: 10},
	}
	s := market.SummarizeTrading(days)
	// Volume-weighted: 75 / 50 = 1.5, unlike the plain mean of 1.83.
	if s.Days != 3 || s.Volume != 50 || s.AveragePrice != 1.5 || s.Low != 1 || s.High != 2.5 || s.First != 2 || s.Last != 2.5 || !near(s.Change, 0.25) {
		t.Errorf("SummarizeTrading = %+v", s)
	}
}
