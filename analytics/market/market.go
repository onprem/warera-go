package market

import (
	"cmp"
	"maps"
	"math"
	"slices"

	"github.com/onprem/warera-go"
)

// Fill is the result of walking an order book for a quantity.
type Fill struct {
	// Quantity is how much the listed orders cover, at most the quantity
	// asked for, and Complete reports whether they cover all of it.
	Quantity int
	Complete bool
	// Total is the money paid or received for Quantity, and AveragePrice
	// Total divided by Quantity.
	Total        float64
	AveragePrice float64
}

// BuyCost returns what buying quantity units costs right now by taking the
// cheapest sell orders first. Order books from Trading.GetTopOrders list only
// the best orders, so a large quantity may come back incomplete.
func BuyCost(book *warera.TopOrders, quantity int) Fill {
	return walk(book.SellOrders, quantity, func(a, b warera.TopOrder) int { return cmp.Compare(a.Price, b.Price) })
}

// SellProceeds returns what selling quantity units earns right now by filling
// the highest buy orders first.
func SellProceeds(book *warera.TopOrders, quantity int) Fill {
	return walk(book.BuyOrders, quantity, func(a, b warera.TopOrder) int { return cmp.Compare(b.Price, a.Price) })
}

func walk(orders []warera.TopOrder, quantity int, better func(a, b warera.TopOrder) int) Fill {
	sorted := slices.Clone(orders)
	slices.SortStableFunc(sorted, better)
	var f Fill
	for _, o := range sorted {
		if f.Quantity >= quantity {
			break
		}
		take := min(o.Quantity, quantity-f.Quantity)
		f.Quantity += take
		f.Total += float64(take) * o.Price
	}
	f.Complete = f.Quantity >= quantity
	if f.Quantity > 0 {
		f.AveragePrice = f.Total / float64(f.Quantity)
	}
	return f
}

// Spread returns the highest buy price and the lowest sell price in an order
// book. ok is false when either side is empty.
func Spread(book *warera.TopOrders) (bid, ask float64, ok bool) {
	if len(book.BuyOrders) == 0 || len(book.SellOrders) == 0 {
		return 0, 0, false
	}
	bid = slices.MaxFunc(book.BuyOrders, byPrice).Price
	ask = slices.MinFunc(book.SellOrders, byPrice).Price
	return bid, ask, true
}

// PriceSummary summarizes a run of trading days.
type PriceSummary struct {
	Days int
	// AveragePrice is the volume-weighted average over all days:
	// Turnover / Volume.
	AveragePrice float64
	// Low and High are the lowest and highest daily average prices, and
	// First and Last the daily averages of the first and last day.
	Low, High   float64
	First, Last float64
	// Change is Last relative to First, as a fraction: 0.05 is a 5% rise.
	Change   float64
	Volume   int
	Turnover float64
}

// SummarizeTrading summarizes trading days, such as
// [warera.ItemTrading.Values] from Trading.GetItemTrading. The last day of
// that history is usually today and still incomplete; drop it to compare
// whole days. Days without trades are skipped.
func SummarizeTrading(days []warera.ItemTradingDay) PriceSummary {
	var s PriceSummary
	for _, d := range days {
		if d.TotalQuantity == 0 {
			continue
		}
		if s.Days == 0 {
			s.First, s.Low, s.High = d.AvgValue, d.AvgValue, d.AvgValue
		}
		s.Days++
		s.Last = d.AvgValue
		s.Low = min(s.Low, d.AvgValue)
		s.High = max(s.High, d.AvgValue)
		s.Volume += d.TotalQuantity
		s.Turnover += d.TotalValue
	}
	if s.Volume > 0 {
		s.AveragePrice = s.Turnover / float64(s.Volume)
	}
	if s.First > 0 {
		s.Change = s.Last/s.First - 1
	}
	return s
}

// Basis is the prices a tool values trade at: Sell for what it sells, such
// as a company's output, and Buy for what it buys, such as inputs.
type Basis struct {
	Sell, Buy warera.ItemPrices
}

// Reference values both selling and buying at reference prices, such as
// those from Trading.GetPrices.
func Reference(prices warera.ItemPrices) Basis {
	return Basis{Sell: prices, Buy: prices}
}

// BidAsk values selling at each order book's best bid and buying at its best
// ask, what trading right now gets, with books keyed by item code as
// returned by Trading.GetTopOrdersPerItem. An item without a book, or a
// side without orders, falls back to its reference price.
func BidAsk(prices warera.ItemPrices, books map[string]warera.TopOrders) Basis {
	b := Basis{Sell: warera.ItemPrices{}, Buy: warera.ItemPrices{}}
	maps.Copy(b.Sell, prices)
	maps.Copy(b.Buy, prices)
	for code, book := range books {
		if len(book.BuyOrders) > 0 {
			b.Sell[code] = slices.MaxFunc(book.BuyOrders, byPrice).Price
		}
		if len(book.SellOrders) > 0 {
			b.Buy[code] = slices.MinFunc(book.SellOrders, byPrice).Price
		}
	}
	return b
}

func byPrice(a, b warera.TopOrder) int { return cmp.Compare(a.Price, b.Price) }

// RollPrice estimates what an equipment item with given rolled stats sells
// for, from recent sales such as Transactions.List with TransactionTypes
// itemMarket and the item code. It uses the sales of that code at full
// durability whose rolls are nearest to stats, each stat's difference
// measured as a share of its range in the game config, and returns their
// median price. ok is false when no such sale is listed.
func RollPrice(cfg *warera.GameConfig, sales []warera.TransactionListItem, itemCode string, stats map[string]float64) (price float64, ok bool) {
	ranges := cfg.Items[itemCode].DynamicStats
	best := math.Inf(1)
	var prices []float64
	for _, s := range sales {
		it := s.Item
		if s.TransactionType != warera.TransactionTypeItemMarket || it == nil || it.Code != itemCode || it.State < it.MaxState || s.Quantity < 1 {
			continue
		}
		rolled := itemStats(it.Skills)
		var d float64
		for stat, r := range ranges {
			if len(r) == 2 && r[1] > r[0] {
				d += math.Abs(rolled[stat]-stats[stat]) / (r[1] - r[0])
			}
		}
		switch {
		case d < best-1e-9:
			best, prices = d, []float64{s.Money / float64(s.Quantity)}
		case d < best+1e-9:
			prices = append(prices, s.Money/float64(s.Quantity))
		}
	}
	if len(prices) == 0 {
		return 0, false
	}
	slices.Sort(prices)
	if n := len(prices); n%2 == 0 {
		return (prices[n/2-1] + prices[n/2]) / 2, true
	}
	return prices[len(prices)/2], true
}

// itemStats returns an item's rolled stats keyed by stat name.
func itemStats(s warera.ItemSkills) map[string]float64 {
	stats := map[string]float64{}
	for name, v := range map[string]float64{
		"armor": s.Armor, "attack": s.Attack, "criticalChance": s.CriticalChance,
		"criticalDamages": s.CriticalDamages, "dodge": s.Dodge, "precision": s.Precision,
	} {
		if v != 0 {
			stats[name] = v
		}
	}
	return stats
}
