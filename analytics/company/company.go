package company

import (
	"fmt"
	"math"
	"slices"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics"
	"github.com/onprem/warera-go/analytics/market"
	"github.com/onprem/warera-go/analytics/upgrades"
)

// productionStatPerPoint is the production stat that makes one production
// point per hour, before bonuses.
const productionStatPerPoint = 4

// modestWage is the share of the usual wage a modest employee asks.
const modestWage = 0.85

// EmployeeProduction returns the production points an employee makes per
// hour: their production stat total divided by 4, times their production
// bonus. extraBonusPercent adds percentage points to the bonus, for
// example to try a manager with a different bonus; pass 0 for the
// employee as they are.
//
// Verified: equals [warera.EmployeeEstimate.Production] for every employee
// of two recorded rosters.
func EmployeeProduction(e warera.Employee, extraBonusPercent float64) float64 {
	return e.Breakdown["production"].Total / productionStatPerPoint * (1 + (e.ProductionBonus.Total+extraBonusPercent)/100)
}

// EmployeeWage returns an employee's hourly wage: ppMoneyValue, from
// MarketStats.GetPPMoneyValue, times their wage expectation, times the
// production points their trained production stat makes, before bonuses.
// A modest employee asks 85% of that. Trait bonuses to production are not
// paid for.
//
// Verified: equals [warera.Employee.Wage] for every employee of two
// recorded rosters, including a modest one.
func EmployeeWage(e warera.Employee, ppMoneyValue float64) float64 {
	wage := ppMoneyValue * e.WageExpectation * e.Base.Production / productionStatPerPoint
	if slices.ContainsFunc(e.Traits, func(t warera.EmployeeTrait) bool { return t.Code == "modest" }) {
		wage *= modestWage
	}
	return wage
}

// ExtraProduction returns the production points per hour a roster's
// employees would gain from bonusPercent more production bonus, such as
// [ManagerBonusPercent] for hiring a manager.
//
// Verified: production bonuses add up, as in [EmployeeProduction].
func ExtraProduction(roster *warera.CompanyRoster, bonusPercent float64) float64 {
	var total float64
	for _, e := range roster.Employees {
		total += EmployeeProduction(e, bonusPercent) - EmployeeProduction(e, 0)
	}
	return total
}

// ManagerWage returns the hourly wage of a manager with a Management level
// at a pay rate per Management point: rate × (level + 1).
//
// Verified: equals [warera.Manager.HourlyWage] of recorded managers.
func ManagerWage(rateManagement float64, managementLevel int) float64 {
	return rateManagement * float64(managementLevel+1)
}

// ManagerBonusPercent returns the production bonus a manager with a
// Management level gives each employee: 2 × (level + 1) percent.
//
// Verified: equals [warera.Manager.ProductionBonusPercent].
func ManagerBonusPercent(managementLevel int) float64 {
	return 2 * float64(managementLevel+1)
}

// IncomeTaxHappiness returns the happiness an income tax in percent costs
// each employee: tax² / 60. A 20% tax costs 6.67.
//
// Verified against recorded happiness breakdowns.
func IncomeTaxHappiness(taxPercent float64) float64 {
	return -taxPercent * taxPercent / 60
}

// ResponsibilityHappiness returns the happiness an employee's
// responsibility stat gives: 3 per point.
//
// Verified against recorded happiness breakdowns.
func ResponsibilityHappiness(responsibility float64) float64 {
	return 3 * responsibility
}

// HappinessBonusPercent returns the production bonus an employee's
// happiness gives: happiness - 50, negative below 50.
//
// Verified against recorded production bonuses.
func HappinessBonusPercent(happiness float64) float64 {
	return happiness - 50
}

// EngineProduction returns the production points per hour of a company's
// automated engine at a level: the level's "dailyProd" stat divided by 24,
// times the roster's [warera.CompanyHourlyOutput.EngineMultiplier] and the
// company's production bonus percent, from Companies.GetProductionBonus.
//
// Game config: the daily production is config; applying the multiplier and
// the bonus follows the Warera Simulator.
func EngineProduction(cfg *warera.GameConfig, level int, engineMultiplier, bonusPercent float64) (float64, error) {
	daily, err := upgradeStat(cfg, warera.UpgradeTypeAutomatedEngine, level, "dailyProd")
	return daily / 24 * engineMultiplier * (1 + bonusPercent/100), err
}

func upgradeStat(cfg *warera.GameConfig, upgrade warera.UpgradeType, level int, stat string) (float64, error) {
	levels, err := upgrades.Levels(cfg, upgrade)
	if err != nil {
		return 0, err
	}
	i := slices.IndexFunc(levels, func(l warera.GameConfigUpgradeLevel) bool { return l.Level == level })
	if i < 0 {
		return 0, fmt.Errorf("%w: %q has no level %d", analytics.ErrUnknownUpgrade, upgrade, level)
	}
	return levels[i].Stats[stat], nil
}

// Hourly is a company's hourly result.
type Hourly struct {
	ItemCode string
	// Disabled reports a disabled company, which produces and costs
	// nothing.
	Disabled bool
	// EmployeeProduction and EngineProduction are production points per
	// hour, adding up to Production. Units is the items they make.
	EmployeeProduction float64
	EngineProduction   float64
	Production         float64
	Units              float64
	// Revenue is Units at the basis's sell price, and InputCost their
	// inputs at its buy price.
	Revenue     float64
	InputCost   float64
	Wages       float64
	ManagerCost float64
	Profit      float64
	// BreakEvenPrice is the sell price at which Profit is zero.
	BreakEvenPrice float64
	// StorageCap is the production points the company stores, and
	// HoursToFull how long until production fills it; +Inf when nothing is
	// produced.
	StorageCap  float64
	HoursToFull float64
}

// HourlyOf returns a company's hourly result: employee production, wages,
// and manager cost as its roster reports them, Employees.GetCompanyRoster, plus
// [EngineProduction], with output and inputs valued at a [market.Basis].
// bonus is from Companies.GetProductionBonus. Company upkeep, tool wear,
// and market fees are not included.
//
// Inputs are valued at the buy price even when the owner makes them,
// because selling them is the alternative. A disabled company, one with a
// [warera.Company.DisabledAt], produces and costs nothing: its roster
// reports no production or wages, and its engine stops too.
func HourlyOf(cfg *warera.GameConfig, c *warera.Company, roster *warera.CompanyRoster, bonus *warera.CompanyProductionBonus, basis market.Basis) (*Hourly, error) {
	item, ok := cfg.Items[c.ItemCode]
	if !ok {
		return nil, fmt.Errorf("%w %q", analytics.ErrUnknownItem, c.ItemCode)
	}
	if item.ProductionPoints <= 0 {
		return nil, fmt.Errorf("company: %q has no production points", c.ItemCode)
	}
	price, ok := basis.Sell[c.ItemCode]
	if !ok {
		return nil, fmt.Errorf("company: no price for %q", c.ItemCode)
	}
	engine, err := EngineProduction(cfg, c.ActiveUpgradeLevels.AutomatedEngine, roster.Hourly.EngineMultiplier, bonus.Total)
	if err != nil {
		return nil, err
	}
	disabled := !c.DisabledAt.IsZero()
	if disabled {
		engine = 0
	}
	storage, err := upgradeStat(cfg, warera.UpgradeTypeStorage, c.ActiveUpgradeLevels.Storage, "maxProduction")
	if err != nil {
		return nil, err
	}
	h := &Hourly{
		ItemCode:           c.ItemCode,
		Disabled:           disabled,
		EmployeeProduction: roster.Hourly.Production,
		EngineProduction:   engine,
		Production:         roster.Hourly.Production + engine,
		Wages:              roster.Hourly.WageCost,
		ManagerCost:        roster.Hourly.ManagerCost,
		StorageCap:         storage,
		HoursToFull:        math.Inf(1),
	}
	h.Units = h.Production / item.ProductionPoints
	h.Revenue = h.Units * price
	for input, quantity := range item.ProductionNeeds {
		inputPrice, ok := basis.Buy[input]
		if !ok {
			return nil, fmt.Errorf("company: no price for %q, an input of %q", input, c.ItemCode)
		}
		h.InputCost += h.Units * quantity * inputPrice
	}
	costs := h.InputCost + h.Wages + h.ManagerCost
	h.Profit = h.Revenue - costs
	if h.Units > 0 {
		h.BreakEvenPrice = costs / h.Units
	}
	if h.Production > 0 {
		h.HoursToFull = max(0, storage-c.Production) / h.Production
	}
	return h, nil
}
