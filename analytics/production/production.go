package production

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
)

// Cost returns the production points needed to make one unit of an
// item from raw resources: the item's own ProductionPoints plus, for each
// input, the quantity needed times that input's production cost. Steel, which
// takes 10 production points and 10 iron at 1 point each, costs 20.
func Cost(cfg *warera.GameConfig, itemCode string) (float64, error) {
	return cost(cfg, itemCode, map[string]bool{})
}

func cost(cfg *warera.GameConfig, itemCode string, visiting map[string]bool) (float64, error) {
	item, ok := cfg.Items[itemCode]
	if !ok {
		return 0, fmt.Errorf("%w %q", analytics.ErrUnknownItem, itemCode)
	}
	if visiting[itemCode] {
		return 0, fmt.Errorf("production: item %q is its own input", itemCode)
	}
	visiting[itemCode] = true
	defer delete(visiting, itemCode)
	total := item.ProductionPoints
	for input, quantity := range item.ProductionNeeds {
		inputCost, err := cost(cfg, input, visiting)
		if err != nil {
			return 0, err
		}
		total += quantity * inputCost
	}
	return total, nil
}

// RawMaterials returns the raw resources, items with no inputs of their own,
// that making quantity units of an item consumes, keyed by item code. A raw
// resource itself returns only itself.
func RawMaterials(cfg *warera.GameConfig, itemCode string, quantity float64) (map[string]float64, error) {
	materials := map[string]float64{}
	if err := addRawMaterials(cfg, itemCode, quantity, materials, map[string]bool{}); err != nil {
		return nil, err
	}
	return materials, nil
}

func addRawMaterials(cfg *warera.GameConfig, itemCode string, quantity float64, materials map[string]float64, visiting map[string]bool) error {
	item, ok := cfg.Items[itemCode]
	if !ok {
		return fmt.Errorf("%w %q", analytics.ErrUnknownItem, itemCode)
	}
	if len(item.ProductionNeeds) == 0 {
		materials[itemCode] += quantity
		return nil
	}
	if visiting[itemCode] {
		return fmt.Errorf("production: item %q is its own input", itemCode)
	}
	visiting[itemCode] = true
	defer delete(visiting, itemCode)
	for input, needed := range item.ProductionNeeds {
		if err := addRawMaterials(cfg, input, quantity*needed, materials, visiting); err != nil {
			return err
		}
	}
	return nil
}

// Value is what producing one item earns at given market prices.
type Value struct {
	ItemCode string
	// Price is the market price of one unit.
	Price float64
	// ProductionPoints is the points one unit takes in its own production
	// step, and TotalProductionPoints the points including its inputs (see
	// [Cost]).
	ProductionPoints      float64
	TotalProductionPoints float64
	// InputCost is the market price of the inputs one unit consumes.
	InputCost float64
	// MoneyPerPoint is what a production point earns when the inputs are
	// bought at market prices: (Price - InputCost) / ProductionPoints. Compare
	// it across items to choose what a company should make, and with wages
	// to see whether production pays.
	MoneyPerPoint float64
	// MoneyPerPointFromRaw is what a production point earns when every step
	// is produced in-house: Price / TotalProductionPoints. The server's
	// MarketStats.GetPPMoneyValueByItem reports the same figure.
	MoneyPerPointFromRaw float64
}

// Values returns the [Value] of every item that has
// production points and a price, along with prices for all its inputs,
// sorted by MoneyPerPoint, highest first.
func Values(cfg *warera.GameConfig, prices warera.ItemPrices) []Value {
	var values []Value
	for code, item := range cfg.Items {
		if v, ok := value(cfg, prices, code, item); ok {
			values = append(values, v)
		}
	}
	slices.SortFunc(values, func(a, b Value) int {
		return cmp.Or(cmp.Compare(b.MoneyPerPoint, a.MoneyPerPoint), cmp.Compare(a.ItemCode, b.ItemCode))
	})
	return values
}

func value(cfg *warera.GameConfig, prices warera.ItemPrices, code string, item warera.GameConfigItem) (Value, bool) {
	price, ok := prices[code]
	if !ok || item.ProductionPoints <= 0 {
		return Value{}, false
	}
	total, err := Cost(cfg, code)
	if err != nil {
		return Value{}, false
	}
	v := Value{ItemCode: code, Price: price, ProductionPoints: item.ProductionPoints, TotalProductionPoints: total}
	for input, quantity := range item.ProductionNeeds {
		inputPrice, ok := prices[input]
		if !ok {
			return Value{}, false
		}
		v.InputCost += quantity * inputPrice
	}
	v.MoneyPerPoint = (price - v.InputCost) / item.ProductionPoints
	v.MoneyPerPointFromRaw = price / total
	return v, true
}
