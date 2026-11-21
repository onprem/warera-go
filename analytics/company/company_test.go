package company_test

import (
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/company"
)

// A chain with quantities other than 1 and an input used at two depths:
// ore is raw (2 points), plate takes 3 points and 5 ore, and gear takes 1
// point, 2 plates, and 1 ore.
func chainConfig() *warera.GameConfig {
	return &warera.GameConfig{Items: map[string]warera.GameConfigItem{
		"ore":   {ProductionPoints: 2},
		"plate": {ProductionPoints: 3, ProductionNeeds: map[string]float64{"ore": 5}},
		"gear":  {ProductionPoints: 1, ProductionNeeds: map[string]float64{"plate": 2, "ore": 1}},
		"loop":  {ProductionPoints: 1, ProductionNeeds: map[string]float64{"loop": 1}},
		"cases": {},
	}}
}

func TestProfitPerHour(t *testing.T) {
	cfg := chainConfig()
	prices := warera.ItemPrices{"ore": 0.5, "plate": 9}
	summary := &warera.WorkforceSummary{
		ProductionByItem: map[string]float64{"plate": 30, "ore": 10},
		WageCost:         4, ManagerCost: 1,
	}
	got, err := company.ProfitPerHour(cfg, prices, summary)
	if err != nil {
		t.Fatal(err)
	}
	// plate: 30 points / 3 = 10 units, revenue 90, inputs 10*5*0.5 = 25.
	// ore:   10 points / 2 = 5 units, revenue 2.5.
	// profit = 92.5 - 25 - 4 - 1 = 62.5.
	plate := got.ByItem["plate"]
	if plate.Units != 10 || plate.Revenue != 90 || plate.InputCost != 25 || got.Revenue != 92.5 || got.Profit != 62.5 {
		t.Errorf("ProfitPerHour = %+v (plate %+v); want revenue 92.5, profit 62.5", got, plate)
	}
	if _, err := company.ProfitPerHour(cfg, warera.ItemPrices{"plate": 9}, summary); err == nil {
		t.Error("an unpriced input must be an error, not a silent zero cost")
	}
}
