package warera

// Response models for employees and managers. WarEra's published OpenAPI
// document does not describe these procedures; the models follow recorded
// responses.

// CompanyRoster is a company's employees and its hourly output.
type CompanyRoster struct {
	Employees []Employee          `json:"employees,omitempty"`
	Hourly    CompanyHourlyOutput `json:"hourly,omitempty"`
}

// CompanyHourlyOutput is a company's estimated output and costs per hour.
type CompanyHourlyOutput struct {
	// Cap is the number of employees the company's break room allows.
	Cap              int     `json:"cap,omitempty"`
	EngineMultiplier float64 `json:"engineMultiplier,omitempty"`
	ItemCode         string  `json:"itemCode,omitempty"`
	ManagerCost      float64 `json:"managerCost,omitempty"`
	Production       float64 `json:"production,omitempty"`
	Toolboxes        float64 `json:"toolboxes,omitempty"`
	WageCost         float64 `json:"wageCost,omitempty"`
	// IfEnabled is the output a disabled company would have once enabled.
	IfEnabled *CompanyHourlyEstimate `json:"ifEnabled,omitempty"`
}

// CompanyHourlyEstimate is an estimated hourly output and cost.
type CompanyHourlyEstimate struct {
	ManagerCost float64 `json:"managerCost,omitempty"`
	Production  float64 `json:"production,omitempty"`
	Toolboxes   float64 `json:"toolboxes,omitempty"`
	WageCost    float64 `json:"wageCost,omitempty"`
}

// Employee is an NPC employee working in a player's company.
type Employee struct {
	ID              string `json:"_id,omitempty"`
	Age             int    `json:"age,omitempty"`
	CompanyDisabled bool   `json:"companyDisabled,omitempty"`
	CompanyID       string `json:"company,omitempty"`
	CompanyMovedAt  Time   `json:"companyMovedAt,omitempty"`
	CreatedAt       Time   `json:"createdAt,omitempty"`
	Gender          string `json:"gender,omitempty"`
	Name            string `json:"name,omitempty"`
	NationalityID   string `json:"nationality,omitempty"`
	OwnerID         string `json:"owner,omitempty"`
	UpdatedAt       Time   `json:"updatedAt,omitempty"`

	// Base holds the stats the employee was hired with, and Trained the
	// points added by training.
	Base    EmployeeStats   `json:"base,omitempty"`
	Trained EmployeeStats   `json:"trained,omitempty"`
	Traits  []EmployeeTrait `json:"traits,omitempty"`
	// Breakdown explains each stat's total by source, keyed by stat name.
	Breakdown       map[string]EmployeeStatBreakdown `json:"breakdown,omitempty"`
	Potential       int                              `json:"potential,omitempty"`
	PotentialMax    int                              `json:"potentialMax,omitempty"`
	Specializations map[string]float64               `json:"specializations,omitempty"`
	RetireInDays    float64                          `json:"retireInDays,omitempty"`

	Wage            float64 `json:"wage,omitempty"`
	WageExpectation float64 `json:"wageExpectation,omitempty"`
	// WageNegotiated is the wage reduction won through training, if any.
	WageNegotiated  float64                 `json:"wageNegotiated,omitempty"`
	Estimate        EmployeeEstimate        `json:"estimate,omitempty"`
	ProductionBonus EmployeeProductionBonus `json:"productionBonus,omitempty"`

	HappinessBreakdown EmployeeHappiness       `json:"happinessBreakdown,omitempty"`
	EventModifiers     []EmployeeEventModifier `json:"eventModifiers,omitempty"`
	CharmLuck          *EmployeeCharmLuck      `json:"charmLuck,omitempty"`
	// CareBoosts lists when a manager's daily care boosted this employee.
	CareBoosts       []Time `json:"careBoosts,omitempty"`
	CoworkerFiredAts []Time `json:"coworkerFiredAts,omitempty"`

	ToolID     string              `json:"tool,omitempty"`
	ToolItem   *EmployeeTool       `json:"toolItem,omitempty"`
	ToolStatus *EmployeeToolStatus `json:"toolStatus,omitempty"`

	TrainingHistory       []EmployeeTraining `json:"trainingHistory,omitempty"`
	TrainingsPerDay       int                `json:"trainingsPerDay,omitempty"`
	TrainingsToday        int                `json:"trainingsToday,omitempty"`
	TrainerSessionsUsed   int                `json:"trainerSessionsUsed,omitempty"`
	TrainerSessionsUsedAt Time               `json:"trainerSessionsUsedAt,omitempty"`

	LastAgedDay        Date `json:"lastAgedDay,omitempty"`
	LastEventAt        Time `json:"lastEventAt,omitempty"`
	LastMovedAt        Time `json:"lastMovedAt,omitempty"`
	LastSpecializedDay Date `json:"lastSpecializedDay,omitempty"`
	LastTrainedAt      Time `json:"lastTrainedAt,omitempty"`
}

// EmployeeStats are an employee's production, luck, and responsibility.
type EmployeeStats struct {
	Luck           float64 `json:"luck,omitempty"`
	Production     float64 `json:"production,omitempty"`
	Responsibility float64 `json:"responsibility,omitempty"`
}

// EmployeeTrait is a personality trait. Revealed is reported when training
// uncovers a hidden trait.
type EmployeeTrait struct {
	// Arg parameterizes some traits, such as the item code a trait refers to.
	Arg      string `json:"arg,omitempty"`
	Code     string `json:"code,omitempty"`
	Revealed bool   `json:"revealed,omitempty"`
}

// EmployeeStatBreakdown explains one stat's total.
type EmployeeStatBreakdown struct {
	Base          float64                `json:"base,omitempty"`
	Contributions []EmployeeContribution `json:"contributions,omitempty"`
	Total         float64                `json:"total,omitempty"`
	Trained       float64                `json:"trained,omitempty"`
}

// EmployeeContribution is one source of a stat change. Percent reports
// whether Delta is a percentage rather than flat points.
type EmployeeContribution struct {
	Code    string  `json:"code,omitempty"`
	Delta   float64 `json:"delta,omitempty"`
	Percent bool    `json:"percent,omitempty"`
}

// EmployeeEstimate is an employee's estimated hourly production and wage.
type EmployeeEstimate struct {
	ItemCode   string  `json:"itemCode,omitempty"`
	Production float64 `json:"production,omitempty"`
	Wage       float64 `json:"wage,omitempty"`
}

// EmployeeProductionBonus is an employee's total production bonus percent by
// source.
type EmployeeProductionBonus struct {
	Contributions []EmployeeBonusSource `json:"contributions,omitempty"`
	Total         float64               `json:"total,omitempty"`
}

// EmployeeBonusSource is one source of a production bonus.
type EmployeeBonusSource struct {
	// EventCode and Until identify a temporary bonus from an employee event.
	EventCode string  `json:"eventCode,omitempty"`
	Percent   float64 `json:"percent,omitempty"`
	Source    string  `json:"source,omitempty"`
	Until     Time    `json:"until,omitempty"`
}

// EmployeeHappiness is an employee's happiness by source, and the chance they
// quit.
type EmployeeHappiness struct {
	Base          float64                   `json:"base,omitempty"`
	Contributions []EmployeeHappinessReason `json:"contributions,omitempty"`
	Mitigated     float64                   `json:"mitigated,omitempty"`
	QuitChance    float64                   `json:"quitChance,omitempty"`
	Value         float64                   `json:"value,omitempty"`
}

// EmployeeHappinessReason is one source of a happiness change: a general
// Reason such as income tax, or the Code of a trait.
type EmployeeHappinessReason struct {
	Code  string  `json:"code,omitempty"`
	Delta float64 `json:"delta,omitempty"`
	// EventCode and Until identify a temporary change from an employee event.
	EventCode string `json:"eventCode,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Until     Time   `json:"until,omitempty"`
}

// EmployeeEventModifier is a temporary effect from an employee event.
type EmployeeEventModifier struct {
	Code  string  `json:"code,omitempty"`
	Delta float64 `json:"delta,omitempty"`
	Kind  string  `json:"kind,omitempty"`
	Until Time    `json:"until,omitempty"`
}

// EmployeeCharmLuck is a temporary luck bonus.
type EmployeeCharmLuck struct {
	Amount float64 `json:"amount,omitempty"`
	Until  Time    `json:"until,omitempty"`
}

// EmployeeTool is the tool an employee holds.
type EmployeeTool struct {
	Item
	CreatedAt  Time   `json:"createdAt,omitempty"`
	EmployeeID string `json:"employee,omitempty"`
	OwnerID    string `json:"owner,omitempty"`
	UpdatedAt  Time   `json:"updatedAt,omitempty"`
}

// EmployeeToolStatus reports whether an employee's tool works and the bonus
// it gives.
type EmployeeToolStatus struct {
	BonusPercent  float64 `json:"bonusPercent,omitempty"`
	IsBroken      bool    `json:"isBroken,omitempty"`
	IsUsable      bool    `json:"isUsable,omitempty"`
	RolledPercent float64 `json:"rolledPercent,omitempty"`
}

// EmployeeTraining is one training card an employee received.
type EmployeeTraining struct {
	Card     string                  `json:"card,omitempty"`
	Delta    float64                 `json:"delta,omitempty"`
	Outcome  EmployeeTrainingOutcome `json:"outcome,omitempty"`
	PickedAt Time                    `json:"pickedAt,omitempty"`
	Rarity   string                  `json:"rarity,omitempty"`
	Stat     string                  `json:"stat,omitempty"`
	// Trait is the trait a card removes or changes; TraitTo and TraitTo2 are
	// the traits it offers instead.
	Trait    string `json:"trait,omitempty"`
	TraitTo  string `json:"traitTo,omitempty"`
	TraitTo2 string `json:"traitTo2,omitempty"`
}

// EmployeeTrainingOutcome is the result of a training card.
type EmployeeTrainingOutcome struct {
	Gained          []EmployeeTrait `json:"gained,omitempty"`
	Lost            []string        `json:"lost,omitempty"`
	PotentialMax    int             `json:"potentialMax,omitempty"`
	Stat            string          `json:"stat,omitempty"`
	StatGain        float64         `json:"statGain,omitempty"`
	WageExpectation float64         `json:"wageExpectation,omitempty"`
}

// Manager is a player hired to run another player's company.
type Manager struct {
	ID                      string  `json:"_id,omitempty"`
	Accrued                 float64 `json:"accrued,omitempty"`
	CanEndAt                Time    `json:"canEndAt,omitempty"`
	CompanyID               string  `json:"company,omitempty"`
	CreatedAt               Time    `json:"createdAt,omitempty"`
	EmployerID              string  `json:"employer,omitempty"`
	HasPendingRateReduction bool    `json:"hasPendingRateReduction,omitempty"`
	HiredAt                 Time    `json:"hiredAt,omitempty"`
	HourlyWage              float64 `json:"hourlyWage,omitempty"`
	IsInactive              bool    `json:"isInactive,omitempty"`
	LastAccruedAt           Time    `json:"lastAccruedAt,omitempty"`
	LockedUntil             Time    `json:"lockedUntil,omitempty"`
	ManagementLevel         int     `json:"managementLevel,omitempty"`
	MaxHourlyWage           float64 `json:"maxHourlyWage,omitempty"`
	ProductionBonusPercent  float64 `json:"productionBonusPercent,omitempty"`
	// RateManagement is the pay per Management point per hour.
	RateManagement float64 `json:"rateManagement,omitempty"`
	TotalEarned    float64 `json:"totalEarned,omitempty"`
	UnpaidHours    float64 `json:"unpaidHours,omitempty"`
	UpdatedAt      Time    `json:"updatedAt,omitempty"`
	UserID         string  `json:"user,omitempty"`
}

// ManagerOffer is an open manager job at a company.
type ManagerOffer struct {
	ID            string  `json:"_id,omitempty"`
	CompanyID     string  `json:"company,omitempty"`
	CreatedAt     Time    `json:"createdAt,omitempty"`
	MaxHourlyWage float64 `json:"maxHourlyWage,omitempty"`
	MinManagement float64 `json:"minManagement,omitempty"`
	// RateManagement is the pay per Management point per hour.
	RateManagement float64 `json:"rateManagement,omitempty"`
	RegionID       string  `json:"region,omitempty"`
	UpdatedAt      Time    `json:"updatedAt,omitempty"`
	UserID         string  `json:"user,omitempty"`
}

// ManagerOfferFilter specifies options for listing manager offers.
type ManagerOfferFilter struct {
	CountryID string `json:"countryId,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	RegionID  string `json:"regionId,omitempty"`
}

// ManagerWageStats is the allowed range of manager pay rates.
type ManagerWageStats struct {
	AllowedRange ManagerWageRange `json:"allowedRange,omitempty"`
	// TopEligibleOffer is the best offer the caller qualifies for, when known.
	TopEligibleOffer *float64 `json:"topEligibleOffer,omitempty"`
	TopOffer         float64  `json:"topOffer,omitempty"`
}

// ManagerWageRange bounds manager pay per Management point per hour.
type ManagerWageRange struct {
	Average float64 `json:"average,omitempty"`
	Max     float64 `json:"max,omitempty"`
	Min     float64 `json:"min,omitempty"`
	Source  string  `json:"source,omitempty"`
}

// CompanyProductionBonus is a company's production bonus percent by source.
type CompanyProductionBonus struct {
	DepositBonus            float64 `json:"depositBonus,omitempty"`
	StrategicResourcesBonus float64 `json:"strategicResourcesBonus,omitempty"`
	Total                   float64 `json:"total,omitempty"`
}
