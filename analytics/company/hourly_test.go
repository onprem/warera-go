package company_test

import (
	"encoding/json"
	"math"
	"os"
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/company"
	"github.com/onprem/warera-go/analytics/market"
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

// rybka is a public company recorded from the live API: cooked fish,
// engine level 7, storage level 2, a level 8 manager, a 20% income tax,
// and a modest employee.
type rybka struct {
	Company         warera.Company                `json:"company"`
	Roster          warera.CompanyRoster          `json:"roster"`
	Manager         warera.Manager                `json:"manager"`
	ProductionBonus warera.CompanyProductionBonus `json:"productionBonus"`
}

func loadRybka(t *testing.T) rybka {
	data, err := os.ReadFile("testdata/rybka.json")
	if err != nil {
		t.Fatal(err)
	}
	var r rybka
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatal(err)
	}
	return r
}

func TestEmployeesMatchTheAPI(t *testing.T) {
	pp := fixture[float64](t, "gameStat.getPpMoneyValue")
	r := loadRybka(t)
	rosters := []warera.CompanyRoster{fixture[warera.CompanyRoster](t, "employee.getCompanyRoster"), r.Roster}
	for _, roster := range rosters {
		var production float64
		for _, e := range roster.Employees {
			p := company.EmployeeProduction(e, 0)
			production += p
			if math.Round(p*100)/100 != e.Estimate.Production {
				t.Errorf("%s: production %v, API %v", e.Name, p, e.Estimate.Production)
			}
			if w := company.EmployeeWage(e, pp); math.Abs(w-e.Wage) > 1e-5 {
				t.Errorf("%s: wage %v, API %v", e.Name, w, e.Wage)
			}
			h := e.HappinessBreakdown
			for _, c := range h.Contributions {
				if c.Reason == "responsibility" && c.Delta != company.ResponsibilityHappiness(e.Breakdown["responsibility"].Total) {
					t.Errorf("%s: responsibility happiness %v", e.Name, c.Delta)
				}
				if c.Reason == "incomeTax" && math.Round(company.IncomeTaxHappiness(20)*100)/100 != c.Delta {
					t.Errorf("%s: income tax happiness %v at 20%%", e.Name, c.Delta)
				}
			}
			for _, c := range e.ProductionBonus.Contributions {
				if c.Source == "happiness" && math.Abs(c.Percent-company.HappinessBonusPercent(h.Value)) > 1e-9 {
					t.Errorf("%s: happiness bonus %v at happiness %v", e.Name, c.Percent, h.Value)
				}
			}
		}
		if math.Abs(production-roster.Hourly.Production) > 0.02 {
			t.Errorf("roster production %v, API %v", production, roster.Hourly.Production)
		}
	}
	for _, m := range []warera.Manager{fixture[warera.Manager](t, "manager.getByCompanyId"), r.Manager} {
		if w := company.ManagerWage(m.RateManagement, m.ManagementLevel); math.Abs(w-m.HourlyWage) > 1e-9 {
			t.Errorf("manager wage %v, API %v", w, m.HourlyWage)
		}
		if b := company.ManagerBonusPercent(m.ManagementLevel); b != m.ProductionBonusPercent {
			t.Errorf("manager bonus %v, API %v", b, m.ProductionBonusPercent)
		}
	}
	// Without its manager's 18%, Rybka's stats of 12, 15, and 15 lose
	// 42 / 4 × 0.18 = 1.89 points an hour.
	if got := company.ExtraProduction(&r.Roster, -18); math.Abs(got+1.89) > 1e-9 {
		t.Errorf("ExtraProduction(-18) = %v, want -1.89", got)
	}
}

func TestHourlyOf(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	r := loadRybka(t)
	basis := market.Basis{Sell: warera.ItemPrices{"cookedFish": 8}, Buy: warera.ItemPrices{"fish": 0.5}}
	h, err := company.HourlyOf(cfg, &r.Company, &r.Roster, &r.ProductionBonus, basis)
	if err != nil {
		t.Fatal(err)
	}
	// Engine level 7 makes 168 a day, 7 an hour; employees 16.72. 23.72
	// points make 0.593 cooked fish (40 points and 1 fish each).
	// Revenue 4.744, fish 0.2965, wages 0.505053, manager 0.1467.
	costs := 0.2965 + 0.505053 + 0.1467
	near := func(a, b float64) bool { return math.Abs(a-b) < 1e-9 }
	if !near(h.EngineProduction, 7) || !near(h.Production, 23.72) || !near(h.Units, 0.593) ||
		!near(h.Revenue, 4.744) || !near(h.InputCost, 0.2965) || !near(h.Profit, 4.744-costs) || !near(h.BreakEvenPrice, costs/0.593) {
		t.Errorf("HourlyOf = %+v", h)
	}
	// Storage level 2 holds 400 points; 114.61 are stored.
	if h.StorageCap != 400 || !near(h.HoursToFull, (400-114.61)/23.72) {
		t.Errorf("storage %v, full in %v hours", h.StorageCap, h.HoursToFull)
	}
	if _, err := company.HourlyOf(cfg, &r.Company, &r.Roster, &r.ProductionBonus, market.Basis{Sell: basis.Sell}); err == nil {
		t.Error("an unpriced input must be an error")
	}
}
func TestHourlyOfDisabledCompany(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	// testdata/liscie.json is a public company recorded while disabled: its
	// roster reports no production or wages, but it has an engine at level 7.
	data, err := os.ReadFile("testdata/liscie.json")
	if err != nil {
		t.Fatal(err)
	}
	var r rybka
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatal(err)
	}
	if r.Company.DisabledAt.IsZero() {
		t.Fatal("the recorded company should be disabled")
	}
	h, err := company.HourlyOf(cfg, &r.Company, &r.Roster, &r.ProductionBonus, market.Reference(warera.ItemPrices{"coca": 0.1}))
	if err != nil {
		t.Fatal(err)
	}
	if !h.Disabled || h.Production != 0 || h.EngineProduction != 0 || h.Profit != 0 || !math.IsInf(h.HoursToFull, 1) {
		t.Errorf("HourlyOf a disabled company = %+v; want no production or profit", h)
	}
}

func TestWorkPerDay(t *testing.T) {
	cfg := fixture[*warera.GameConfig](t, "gameConfig.getGameConfig")
	// 100 energy regenerates 240 a day: 24 actions of 10, each making 33.
	if got := company.WorkPerDay(cfg, 100, 33); got != 792 {
		t.Errorf("WorkPerDay = %v, want 792", got)
	}
}
