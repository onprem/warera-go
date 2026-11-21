package production_test

import (
	"errors"
	"math"
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
	"github.com/onprem/warera-go/analytics/production"
	"github.com/onprem/warera-go/waretest"
)

func fixture[T any](t *testing.T, procedure string) T {
	t.Helper()
	v, err := waretest.Fixture[T](procedure)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

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

func TestProductionCostAndRawMaterials(t *testing.T) {
	cfg := chainConfig()
	// plate = 3 + 5*2 = 13; gear = 1 + 2*13 + 1*2 = 29.
	for code, want := range map[string]float64{"ore": 2, "plate": 13, "gear": 29} {
		if got, err := production.Cost(cfg, code); err != nil || got != want {
			t.Errorf("ProductionCost(%s) = %v, %v; want %v", code, got, err, want)
		}
	}
	// 3 gears need 3*(2*5 + 1) = 33 ore.
	if got, err := production.RawMaterials(cfg, "gear", 3); err != nil || len(got) != 1 || got["ore"] != 33 {
		t.Errorf("RawMaterials(gear, 3) = %v, %v; want map[ore:33]", got, err)
	}
	if _, err := production.Cost(cfg, "nope"); !errors.Is(err, analytics.ErrUnknownItem) {
		t.Errorf("unknown item error = %v", err)
	}
	if _, err := production.Cost(cfg, "loop"); err == nil {
		t.Error("an item that is its own input must fail instead of recursing forever")
	}
}

func TestProductionValues(t *testing.T) {
	cfg := chainConfig()
	prices := warera.ItemPrices{"ore": 0.5, "plate": 9, "gear": 40}
	values := production.Values(cfg, prices)
	// ore:   (0.5 - 0) / 2 = 0.25 per point; from raw 0.5/2 = 0.25.
	// plate: (9 - 5*0.5) / 3 = 2.1667; from raw 9/13 = 0.6923.
	// gear:  (40 - 2*9 - 0.5) / 1 = 21.5; from raw 40/29 = 1.3793.
	// Ordering by MoneyPerPoint differs from ordering by MoneyPerPointFromRaw,
	// and loop (a cycle) and cases (no points) are left out.
	want := []production.Value{
		{ItemCode: "gear", Price: 40, ProductionPoints: 1, TotalProductionPoints: 29, InputCost: 18.5, MoneyPerPoint: 21.5, MoneyPerPointFromRaw: 40.0 / 29},
		{ItemCode: "plate", Price: 9, ProductionPoints: 3, TotalProductionPoints: 13, InputCost: 2.5, MoneyPerPoint: 6.5 / 3, MoneyPerPointFromRaw: 9.0 / 13},
		{ItemCode: "ore", Price: 0.5, ProductionPoints: 2, TotalProductionPoints: 2, MoneyPerPoint: 0.25, MoneyPerPointFromRaw: 0.25},
	}
	if len(values) != len(want) {
		t.Fatalf("got %d values, want %d: %+v", len(values), len(want), values)
	}
	for i, w := range want {
		g := values[i]
		if g.ItemCode != w.ItemCode || !near(g.InputCost, w.InputCost) || !near(g.MoneyPerPoint, w.MoneyPerPoint) ||
			!near(g.MoneyPerPointFromRaw, w.MoneyPerPointFromRaw) || g.TotalProductionPoints != w.TotalProductionPoints {
			t.Errorf("values[%d] = %+v, want %+v", i, g, w)
		}
	}
	// Without an input's price an item cannot be valued: with plate unpriced,
	// gear is left out too.
	if got := production.Values(cfg, warera.ItemPrices{"gear": 40, "ore": 0.5}); len(got) != 1 || got[0].ItemCode != "ore" {
		t.Errorf("with plate unpriced, got %+v; want only ore", got)
	}
}

// The server's money-per-point figures use the same formula as
// MoneyPerPointFromRaw. The recorded prices and server values were fetched
// minutes apart, so allow for prices moving in between (they differ by
// under 2%); the other plausible formulas miss by 14% or more.
func TestMoneyPerPointFromRawMatchesServer(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	prices := fixture[warera.ItemPrices](t, "itemTrading.getPrices")
	server := fixture[warera.PPValueByItem](t, "gameStat.getPpMoneyValueByItem")
	checked := 0
	for _, v := range production.Values(cfg, prices) {
		want, ok := server[v.ItemCode]
		if !ok {
			continue
		}
		checked++
		if math.Abs(v.MoneyPerPointFromRaw/want-1) > 0.03 {
			t.Errorf("%s: MoneyPerPointFromRaw = %.4f, server says %.4f", v.ItemCode, v.MoneyPerPointFromRaw, want)
		}
	}
	if checked < 15 {
		t.Errorf("compared only %d items with the server", checked)
	}
}
