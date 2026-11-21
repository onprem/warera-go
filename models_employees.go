package warera

// Response models for employees and managers. WarEra's published OpenAPI
// document does not describe these procedures; the models follow recorded
// responses.

// CompanyRoster is a company's employees and its hourly output.
type CompanyRoster struct {
	Employees []Employee          `json:"employees"`
	Hourly    CompanyHourlyOutput `json:"hourly"`
}

// CompanyHourlyOutput is a company's estimated output and costs per hour.
type CompanyHourlyOutput struct {
	// Cap is the number of employees the company's break room allows.
	Cap              int     `json:"cap"`
	EngineMultiplier float64 `json:"engineMultiplier"`
	ItemCode         string  `json:"itemCode"`
	ManagerCost      float64 `json:"managerCost"`
	Production       float64 `json:"production"`
	Toolboxes        float64 `json:"toolboxes"`
	WageCost         float64 `json:"wageCost"`
	// IfEnabled is the output a disabled company would have once enabled.
	IfEnabled *CompanyHourlyEstimate `json:"ifEnabled"`
}

// CompanyHourlyEstimate is an estimated hourly output and cost.
type CompanyHourlyEstimate struct {
	ManagerCost float64 `json:"managerCost"`
	Production  float64 `json:"production"`
	Toolboxes   float64 `json:"toolboxes"`
	WageCost    float64 `json:"wageCost"`
}

// Employee is an NPC employee working in a player's company.
type Employee struct {
	ID              string `json:"_id"`
	Age             int    `json:"age"`
	CompanyDisabled bool   `json:"companyDisabled"`
	CompanyID       string `json:"company"`
	CompanyMovedAt  Time   `json:"companyMovedAt"`
	CreatedAt       Time   `json:"createdAt"`
	Gender          string `json:"gender"`
	Name            string `json:"name"`
	NationalityID   string `json:"nationality"`
	OwnerID         string `json:"owner"`
	UpdatedAt       Time   `json:"updatedAt"`

	// Base holds the stats the employee was hired with, and Trained the
	// points added by training.
	Base    EmployeeStats   `json:"base"`
	Trained EmployeeStats   `json:"trained"`
	Traits  []EmployeeTrait `json:"traits"`
	// Breakdown explains each stat's total by source, keyed by stat name.
	Breakdown       map[string]EmployeeStatBreakdown `json:"breakdown"`
	Potential       int                              `json:"potential"`
	PotentialMax    int                              `json:"potentialMax"`
	Specializations map[string]float64               `json:"specializations"`
	RetireInDays    float64                          `json:"retireInDays"`

	Wage            float64 `json:"wage"`
	WageExpectation float64 `json:"wageExpectation"`
	// WageNegotiated is the wage reduction won through training, if any.
	WageNegotiated  float64                 `json:"wageNegotiated"`
	Estimate        EmployeeEstimate        `json:"estimate"`
	ProductionBonus EmployeeProductionBonus `json:"productionBonus"`

	HappinessBreakdown EmployeeHappiness       `json:"happinessBreakdown"`
	EventModifiers     []EmployeeEventModifier `json:"eventModifiers"`
	CharmLuck          *EmployeeCharmLuck      `json:"charmLuck"`
	// CareBoosts lists when a manager's daily care boosted this employee.
	CareBoosts       []Time `json:"careBoosts"`
	CoworkerFiredAts []Time `json:"coworkerFiredAts"`

	ToolID     string              `json:"tool"`
	ToolItem   *EmployeeTool       `json:"toolItem"`
	ToolStatus *EmployeeToolStatus `json:"toolStatus"`

	TrainingHistory       []EmployeeTraining `json:"trainingHistory"`
	TrainingsPerDay       int                `json:"trainingsPerDay"`
	TrainingsToday        int                `json:"trainingsToday"`
	TrainerSessionsUsed   int                `json:"trainerSessionsUsed"`
	TrainerSessionsUsedAt Time               `json:"trainerSessionsUsedAt"`

	LastAgedDay        Date `json:"lastAgedDay"`
	LastEventAt        Time `json:"lastEventAt"`
	LastMovedAt        Time `json:"lastMovedAt"`
	LastSpecializedDay Date `json:"lastSpecializedDay"`
	LastTrainedAt      Time `json:"lastTrainedAt"`
}

// EmployeeStats are an employee's production, luck, and responsibility.
type EmployeeStats struct {
	Luck           float64 `json:"luck"`
	Production     float64 `json:"production"`
	Responsibility float64 `json:"responsibility"`
}

// EmployeeTrait is a personality trait. Revealed is reported when training
// uncovers a hidden trait.
type EmployeeTrait struct {
	// Arg parameterizes some traits, such as the item code a trait refers to.
	Arg      string `json:"arg"`
	Code     string `json:"code"`
	Revealed bool   `json:"revealed"`
}

// EmployeeStatBreakdown explains one stat's total.
type EmployeeStatBreakdown struct {
	Base          float64                `json:"base"`
	Contributions []EmployeeContribution `json:"contributions"`
	Total         float64                `json:"total"`
	Trained       float64                `json:"trained"`
}

// EmployeeContribution is one source of a stat change. Percent reports
// whether Delta is a percentage rather than flat points.
type EmployeeContribution struct {
	Code    string  `json:"code"`
	Delta   float64 `json:"delta"`
	Percent bool    `json:"percent"`
}

// EmployeeEstimate is an employee's estimated hourly production and wage.
type EmployeeEstimate struct {
	ItemCode   string  `json:"itemCode"`
	Production float64 `json:"production"`
	Wage       float64 `json:"wage"`
}

// EmployeeProductionBonus is an employee's total production bonus percent by
// source.
type EmployeeProductionBonus struct {
	Contributions []EmployeeBonusSource `json:"contributions"`
	Total         float64               `json:"total"`
}

// EmployeeBonusSource is one source of a production bonus.
type EmployeeBonusSource struct {
	// EventCode and Until identify a temporary bonus from an employee event.
	EventCode string  `json:"eventCode"`
	Percent   float64 `json:"percent"`
	Source    string  `json:"source"`
	Until     Time    `json:"until"`
}

// EmployeeHappiness is an employee's happiness by source, and the chance they
// quit.
type EmployeeHappiness struct {
	Base          float64                   `json:"base"`
	Contributions []EmployeeHappinessReason `json:"contributions"`
	Mitigated     float64                   `json:"mitigated"`
	QuitChance    float64                   `json:"quitChance"`
	Value         float64                   `json:"value"`
}

// EmployeeHappinessReason is one source of a happiness change: a general
// Reason such as income tax, or the Code of a trait.
type EmployeeHappinessReason struct {
	Code  string  `json:"code"`
	Delta float64 `json:"delta"`
	// EventCode and Until identify a temporary change from an employee event.
	EventCode string `json:"eventCode"`
	Reason    string `json:"reason"`
	Until     Time   `json:"until"`
}

// EmployeeEventModifier is a temporary effect from an employee event.
type EmployeeEventModifier struct {
	Code  string  `json:"code"`
	Delta float64 `json:"delta"`
	Kind  string  `json:"kind"`
	Until Time    `json:"until"`
}

// EmployeeCharmLuck is a temporary luck bonus.
type EmployeeCharmLuck struct {
	Amount float64 `json:"amount"`
	Until  Time    `json:"until"`
}

// EmployeeTool is the tool an employee holds.
type EmployeeTool struct {
	Item
	CreatedAt  Time   `json:"createdAt"`
	EmployeeID string `json:"employee"`
	OwnerID    string `json:"owner"`
	UpdatedAt  Time   `json:"updatedAt"`
}

// EmployeeToolStatus reports whether an employee's tool works and the bonus
// it gives.
type EmployeeToolStatus struct {
	BonusPercent  float64 `json:"bonusPercent"`
	IsBroken      bool    `json:"isBroken"`
	IsUsable      bool    `json:"isUsable"`
	RolledPercent float64 `json:"rolledPercent"`
}

// EmployeeTraining is one training card an employee received.
type EmployeeTraining struct {
	Card     string                  `json:"card"`
	Delta    float64                 `json:"delta"`
	Outcome  EmployeeTrainingOutcome `json:"outcome"`
	PickedAt Time                    `json:"pickedAt"`
	Rarity   string                  `json:"rarity"`
	Stat     string                  `json:"stat"`
	// Trait is the trait a card removes or changes; TraitTo and TraitTo2 are
	// the traits it offers instead.
	Trait    string `json:"trait"`
	TraitTo  string `json:"traitTo"`
	TraitTo2 string `json:"traitTo2"`
}

// EmployeeTrainingOutcome is the result of a training card.
type EmployeeTrainingOutcome struct {
	Gained          []EmployeeTrait `json:"gained"`
	Lost            []string        `json:"lost"`
	PotentialMax    int             `json:"potentialMax"`
	Stat            string          `json:"stat"`
	StatGain        float64         `json:"statGain"`
	WageExpectation float64         `json:"wageExpectation"`
}

// Manager is a player hired to run another player's company.
type Manager struct {
	ID                      string  `json:"_id"`
	Accrued                 float64 `json:"accrued"`
	CanEndAt                Time    `json:"canEndAt"`
	CompanyID               string  `json:"company"`
	CreatedAt               Time    `json:"createdAt"`
	EmployerID              string  `json:"employer"`
	HasPendingRateReduction bool    `json:"hasPendingRateReduction"`
	HiredAt                 Time    `json:"hiredAt"`
	HourlyWage              float64 `json:"hourlyWage"`
	IsInactive              bool    `json:"isInactive"`
	LastAccruedAt           Time    `json:"lastAccruedAt"`
	LockedUntil             Time    `json:"lockedUntil"`
	ManagementLevel         int     `json:"managementLevel"`
	MaxHourlyWage           float64 `json:"maxHourlyWage"`
	ProductionBonusPercent  float64 `json:"productionBonusPercent"`
	// RateManagement is the pay per Management point per hour.
	RateManagement float64 `json:"rateManagement"`
	TotalEarned    float64 `json:"totalEarned"`
	UnpaidHours    float64 `json:"unpaidHours"`
	UpdatedAt      Time    `json:"updatedAt"`
	UserID         string  `json:"user"`
}

// ManagerOffer is an open manager job at a company.
type ManagerOffer struct {
	ID            string  `json:"_id"`
	CompanyID     string  `json:"company"`
	CreatedAt     Time    `json:"createdAt"`
	MaxHourlyWage float64 `json:"maxHourlyWage"`
	MinManagement float64 `json:"minManagement"`
	// RateManagement is the pay per Management point per hour.
	RateManagement float64 `json:"rateManagement"`
	RegionID       string  `json:"region"`
	UpdatedAt      Time    `json:"updatedAt"`
	UserID         string  `json:"user"`
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
	AllowedRange ManagerWageRange `json:"allowedRange"`
	// TopEligibleOffer is the best offer the caller qualifies for, when known.
	TopEligibleOffer *float64 `json:"topEligibleOffer"`
	TopOffer         float64  `json:"topOffer"`
}

// ManagerWageRange bounds manager pay per Management point per hour.
type ManagerWageRange struct {
	Average float64 `json:"average"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	Source  string  `json:"source"`
}

// CompanyProductionBonus is a company's production bonus percent by source.
type CompanyProductionBonus struct {
	DepositBonus            float64 `json:"depositBonus"`
	StrategicResourcesBonus float64 `json:"strategicResourcesBonus"`
	Total                   float64 `json:"total"`
}
