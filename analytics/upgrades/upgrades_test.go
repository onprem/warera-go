package upgrades_test

import (
	"errors"
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
	"github.com/onprem/warera-go/analytics/upgrades"
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

func TestUpgradeCost(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	// Live level-5 bunkers report InvestedSteel 2000: levels 2-5 cost
	// 200 + 300 + 500 + 1000.
	got, err := upgrades.Cost(cfg, warera.UpgradeTypeBunker, 1, 5)
	if err != nil || got.Steel != 2000 {
		t.Errorf("UpgradeCost(bunker, 1, 5) = %+v, %v; want 2000 steel", got, err)
	}
	if got, _ := upgrades.Cost(cfg, warera.UpgradeTypeBunker, 0, 1); got.Steel != 100 || got.ConstructionPoints != 100 {
		t.Errorf("UpgradeCost(bunker, 0, 1) = %+v; want 100 steel and 100 construction points", got)
	}
	if _, err := upgrades.Cost(cfg, warera.UpgradeTypeBunker, 3, 2); !errors.Is(err, analytics.ErrUnknownUpgrade) {
		t.Errorf("downward range error = %v", err)
	}
	if _, err := upgrades.Levels(cfg, "moat"); !errors.Is(err, analytics.ErrUnknownUpgrade) {
		t.Errorf("unknown type error = %v", err)
	}
}
