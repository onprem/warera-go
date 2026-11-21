package company

import (
	"fmt"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
)

// WorkforceProfit is a workforce's hourly result at market prices.
type WorkforceProfit struct {
	// Revenue is the market value of the items produced per hour, and
	// InputCost the market value of the inputs they consume.
	Revenue   float64
	InputCost float64
	// Wages and ManagerCost are the hourly pay of employees and managers.
	Wages       float64
	ManagerCost float64
	// Profit is Revenue minus InputCost, Wages, and ManagerCost.
	Profit float64
	// ByItem is each produced item's share, keyed by item code.
	ByItem map[string]ItemOutput
}

// ItemOutput is what a workforce makes of one item per hour.
type ItemOutput struct {
	ProductionPoints float64
	Units            float64
	Revenue          float64
	InputCost        float64
}

// ProfitPerHour values a workforce's hourly production, such as the summary
// from Me.GetWorkforceSummary, at market prices. Inputs are valued at market
// prices even when the owner's own companies supply them, because selling
// them is the alternative. It returns [analytics.ErrUnknownItem] for an item the
// configuration does not define and an error for an item or input without a
// price.
func ProfitPerHour(cfg *warera.GameConfig, prices warera.ItemPrices, summary *warera.WorkforceSummary) (*WorkforceProfit, error) {
	p := &WorkforceProfit{Wages: summary.WageCost, ManagerCost: summary.ManagerCost, ByItem: map[string]ItemOutput{}}
	for code, points := range summary.ProductionByItem {
		item, ok := cfg.Items[code]
		if !ok {
			return nil, fmt.Errorf("%w %q", analytics.ErrUnknownItem, code)
		}
		price, ok := prices[code]
		if !ok || item.ProductionPoints <= 0 {
			return nil, fmt.Errorf("company: no price or production points for %q", code)
		}
		out := ItemOutput{ProductionPoints: points, Units: points / item.ProductionPoints}
		out.Revenue = out.Units * price
		for input, quantity := range item.ProductionNeeds {
			inputPrice, ok := prices[input]
			if !ok {
				return nil, fmt.Errorf("company: no price for %q, an input of %q", input, code)
			}
			out.InputCost += out.Units * quantity * inputPrice
		}
		p.ByItem[code] = out
		p.Revenue += out.Revenue
		p.InputCost += out.InputCost
	}
	p.Profit = p.Revenue - p.InputCost - p.Wages - p.ManagerCost
	return p, nil
}
