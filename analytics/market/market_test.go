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

func TestBidAsk(t *testing.T) {
	ref := warera.ItemPrices{"fish": 1, "steel": 2, "iron": 0.2}
	b := market.BidAsk(ref, map[string]warera.TopOrders{
		"fish":  {BuyOrders: []warera.TopOrder{{Price: 0.8}, {Price: 0.9}}, SellOrders: []warera.TopOrder{{Price: 1.2}, {Price: 1.1}}},
		"steel": {SellOrders: []warera.TopOrder{{Price: 2.5}}}, // nobody bids: sell at reference
	})
	if b.Sell["fish"] != 0.9 || b.Buy["fish"] != 1.1 || b.Sell["steel"] != 2 || b.Buy["steel"] != 2.5 || b.Sell["iron"] != 0.2 || b.Buy["iron"] != 0.2 {
		t.Errorf("BidAsk = %+v", b)
	}
	if ref["fish"] != 1 {
		t.Error("BidAsk changed the reference prices")
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

func TestRollPrice(t *testing.T) {
	cfg := &warera.GameConfig{Items: map[string]warera.GameConfigItem{
		"helmet4": {DynamicStats: map[string][]float64{"criticalDamages": {71, 90}}},
		"sniper":  {DynamicStats: map[string][]float64{"attack": {101, 130}, "criticalChance": {16, 20}}},
	}}
	sale := func(code string, money float64, skills warera.ItemSkills, state float64) warera.TransactionListItem {
		return warera.TransactionListItem{
			TransactionType: warera.TransactionTypeItemMarket, Money: money, Quantity: 1,
			Item: &warera.Item{Code: code, Skills: skills, State: state, MaxState: 100},
		}
	}
	// Recorded helmet4 sales, plus ones that must be ignored.
	sales := []warera.TransactionListItem{
		sale("helmet4", 47.99, warera.ItemSkills{CriticalDamages: 74}, 100),
		sale("helmet4", 47.98, warera.ItemSkills{CriticalDamages: 74}, 100),
		sale("helmet4", 48.10, warera.ItemSkills{CriticalDamages: 74}, 100),
		sale("helmet4", 57.30, warera.ItemSkills{CriticalDamages: 79}, 100),
		sale("helmet4", 64.70, warera.ItemSkills{CriticalDamages: 84}, 100),
		sale("helmet4", 20.00, warera.ItemSkills{CriticalDamages: 80}, 40),  // worn
		sale("helmet5", 99.00, warera.ItemSkills{CriticalDamages: 80}, 100), // other code
		{TransactionType: warera.TransactionTypeDismantleItem, Item: &warera.Item{Code: "helmet4", Skills: warera.ItemSkills{CriticalDamages: 80}, State: 100, MaxState: 100}},
	}
	if p, ok := market.RollPrice(cfg, sales, "helmet4", map[string]float64{"criticalDamages": 80}); !ok || p != 57.30 {
		t.Errorf("RollPrice(80) = %v, %v; want the nearest roll's 57.30", p, ok)
	}
	if p, _ := market.RollPrice(cfg, sales, "helmet4", map[string]float64{"criticalDamages": 74}); p != 47.99 {
		t.Errorf("RollPrice(74) = %v; want the median 47.99", p)
	}
	if _, ok := market.RollPrice(cfg, sales, "helmet6", map[string]float64{"criticalDamages": 130}); ok {
		t.Error("no sales must report !ok")
	}
	// 8 attack of a 29 range is nearer than 2 critical chance of a 4 range.
	weapons := []warera.TransactionListItem{
		sale("sniper", 70, warera.ItemSkills{Attack: 120, CriticalChance: 16}, 100),
		sale("sniper", 60, warera.ItemSkills{Attack: 112, CriticalChance: 18}, 100),
	}
	if p, _ := market.RollPrice(cfg, weapons, "sniper", map[string]float64{"attack": 120, "criticalChance": 18}); p != 60 {
		t.Errorf("RollPrice(sniper) = %v, want 60", p)
	}
}
