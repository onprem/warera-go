package warera

import "encoding/json"

// Response models are maintained separately from the fixture recorder.
// Fixtures provide examples for checking decoding; they do not determine these Go types.

// ItemPrices represents market reference prices indexed by item code.
type ItemPrices map[string]float64

// Article is a news article published in the WarEra game.
type Article struct {
	ID          string       `json:"_id"`
	AuthorID    string       `json:"author"`
	Category    string       `json:"category"`
	Content     string       `json:"content"`
	CreatedAt   Time         `json:"createdAt"`
	IsDeleted   bool         `json:"isDeleted"`
	IsPublic    bool         `json:"isPublic"`
	IsPublished bool         `json:"isPublished"`
	Language    string       `json:"language"`
	PublishedAt Time         `json:"publishedAt"`
	Slug        string       `json:"slug"`
	Stats       ArticleStats `json:"stats"`
	Title       string       `json:"title"`
	UpdatedAt   Time         `json:"updatedAt"`
}

// ArticleStats is a nested component of Article.
type ArticleStats struct {
	Bookmarks int     `json:"bookmarks"`
	Comments  int     `json:"comments"`
	Dislikes  int     `json:"dislikes"`
	GemTips   float64 `json:"gemTips"`
	Likes     int     `json:"likes"`
	Score     float64 `json:"score"`
	Subs      int     `json:"subs"`
	Tips      float64 `json:"tips"`
	Views     int     `json:"views"`
}

// Battle is a military conflict between attacker and defender sides.
type Battle struct {
	ID              string               `json:"_id"`
	Attacker        BattleSideSummary    `json:"attacker"`
	BadgesProcessed bool                 `json:"badgesProcessed"`
	CreatedAt       Time                 `json:"createdAt"`
	CurrentRound    *BattleRound         `json:"currentRound"`
	Defender        BattleSideSummary    `json:"defender"`
	EndedAt         Time                 `json:"endedAt"`
	IsActive        bool                 `json:"isActive"`
	RoundIDs        []string             `json:"rounds"`
	RoundsHistory   []BattleRoundHistory `json:"roundsHistory"`
	RoundsToWin     int                  `json:"roundsToWin"`
	Stats           BattleStats          `json:"stats"`
	Type            string               `json:"type"`
	UpdatedAt       Time                 `json:"updatedAt"`
	WarID           string               `json:"war"`
	WonBy           string               `json:"wonBy"`
}

// BattleRoundHistory summarizes the outcome of a completed battle round.
type BattleRoundHistory struct {
	AttackerDamages float64 `json:"attackerDamages"`
	AttackerPoints  float64 `json:"attackerPoints"`
	DefenderDamages float64 `json:"defenderDamages"`
	DefenderPoints  float64 `json:"defenderPoints"`
	WonBy           string  `json:"wonBy"`
}

// BattleOrder is an order assigned to a country or military unit on a battle side.
type BattleOrder struct {
	ID          string     `json:"_id"`
	BattleID    string     `json:"battle"`
	CountryID   string     `json:"country"`
	CreatedAt   Time       `json:"createdAt"`
	IsActive    bool       `json:"isActive"`
	MUID        string     `json:"mu"`
	Priority    string     `json:"priority"`
	Rank        int        `json:"rank"`
	Side        BattleSide `json:"side"`
	SideCountry string     `json:"sideCountry"`
	Text        string     `json:"text"`
	UpdatedAt   Time       `json:"updatedAt"`
	UserID      string     `json:"user"`
}

// BattleRound is a nested component of Battle.
type BattleRound struct {
	ID        string               `json:"_id"`
	Attacker  RoundSideCombatStats `json:"attacker"`
	BattleID  string               `json:"battle"`
	CreatedAt Time                 `json:"createdAt"`
	Defender  RoundSideCombatStats `json:"defender"`
	IsActive  bool                 `json:"isActive"`
	Live      RoundLiveState       `json:"live"`
	Number    int                  `json:"number"`
	UpdatedAt Time                 `json:"updatedAt"`
}

// BattleSideSummary is a nested component of Battle.
type BattleSideSummary struct {
	BountyEffectiveAt Time     `json:"bountyEffectiveAt"`
	CountryID         string   `json:"country"`
	CountryOrderIDs   []string `json:"countryOrders"`
	Damages           float64  `json:"damages"`
	HitCount          int      `json:"hitCount"`
	MoneyPer1kDamages float64  `json:"moneyPer1kDamages"`
	MoneyPool         float64  `json:"moneyPool"`
	MUOrderIDs        []string `json:"muOrders"`
	RegionID          string   `json:"region"`
	WonRoundsCount    int      `json:"wonRoundsCount"`
}

// BattleStats is a nested component of Battle.
type BattleStats struct {
	HitCount int `json:"hitCount"`
}

// Company is a business that operates in a region and employs workers.
type Company struct {
	ID                  string                     `json:"_id"`
	ActiveUpgradeLevels CompanyActiveUpgradeLevels `json:"activeUpgradeLevels"`
	ConcreteInvested    float64                    `json:"concreteInvested"`
	CreatedAt           Time                       `json:"createdAt"`
	EstimatedValue      float64                    `json:"estimatedValue"`
	FreeMoves           int                        `json:"freeMoves"`
	IsFull              bool                       `json:"isFull"`
	ItemCode            string                     `json:"itemCode"`
	Name                string                     `json:"name"`
	Production          float64                    `json:"production"`
	RegionID            string                     `json:"region"`
	UpdatedAt           Time                       `json:"updatedAt"`
	UserID              string                     `json:"user"`
}

// CompanyActiveUpgradeLevels is a nested component of Company.
type CompanyActiveUpgradeLevels struct {
	AutomatedEngine int `json:"automatedEngine"`
	BreakRoom       int `json:"breakRoom"`
	Storage         int `json:"storage"`
}

// Country is a game nation with political, economic, and diplomatic state.
type Country struct {
	ID                   string                         `json:"_id"`
	AllianceID           string                         `json:"allianceId"`
	AllyIDs              []string                       `json:"allies"`
	AverageDevelopment   float64                        `json:"averageDevelopment"`
	BordersOpenUntil     Time                           `json:"bordersOpenUntil"`
	Code                 string                         `json:"code"`
	CoreDevelopment      float64                        `json:"coreDevelopment"`
	CreatedAt            Time                           `json:"createdAt"`
	CurrentBattleOrderID string                         `json:"currentBattleOrder"`
	CurrentDevelopment   float64                        `json:"currentDevelopment"`
	CurrentPopulation    int                            `json:"currentPopulation"`
	Dates                CountryDates                   `json:"dates"`
	DefensivePactIDs     []string                       `json:"defensivePacts"`
	DiscordURL           string                         `json:"discordUrl"`
	EnemyID              string                         `json:"enemy"`
	GovType              string                         `json:"govType"`
	Income               *IncomeBreakdown               `json:"income"`
	MapAccent            string                         `json:"mapAccent"`
	Money                float64                        `json:"money"`
	Name                 string                         `json:"name"`
	NonAggressionUntil   map[string]Time                `json:"nonAggressionUntil"`
	OrgIDs               []string                       `json:"orgs"`
	PinnedArticleID      string                         `json:"pinnedArticle"`
	Rankings             map[string]RankingItem         `json:"rankings"`
	RulingPartyID        string                         `json:"rulingParty"`
	Scheme               string                         `json:"scheme"`
	StrategicResources   *CountryStrategicResourcesInfo `json:"strategicResources"`
	Taxes                CountryTaxes                   `json:"taxes"`
	Unrest               CountryUnrest                  `json:"unrest"`
	UpdatedAt            Time                           `json:"updatedAt"`
	WarsWithIDs          []string                       `json:"warsWith"`
}

// CountryDates is a nested component of Country.
type CountryDates struct {
	GovTypeChangedAt   Time `json:"govTypeChangedAt"`
	LastAllianceLeftAt Time `json:"lastAllianceLeftAt"`
}

// CountryStrategicResourcesInfo is a nested component of Country.
type CountryStrategicResourcesInfo struct {
	Resources CountryStrategicResourcesInfoResources `json:"resources"`
}

// CountryStrategicResourcesInfoResources is a nested component of CountryStrategicResourcesInfo.
type CountryStrategicResourcesInfoResources struct {
	CoalRegionIDs       []string `json:"coal"`
	DiamondsRegionIDs   []string `json:"diamonds"`
	GoldRegionIDs       []string `json:"gold"`
	LithiumRegionIDs    []string `json:"lithium"`
	RareEarthsRegionIDs []string `json:"rareEarths"`
	UraniumRegionIDs    []string `json:"uranium"`
}

// CountryTaxes is a nested component of Country.
type CountryTaxes struct {
	Income  float64 `json:"income"`
	Manager float64 `json:"manager"`
	Market  float64 `json:"market"`
}

// CountryUnrest is a nested component of Country.
type CountryUnrest struct {
	Bar                   float64 `json:"bar"`
	BarMax                float64 `json:"barMax"`
	LastContributionAt    Time    `json:"lastContributionAt"`
	LastRevolutionEndedAt Time    `json:"lastRevolutionEndedAt"`
}

// CurrentEquipment contains the equipment currently worn by a user.
type CurrentEquipment struct {
	Ammo   string        `json:"ammo"`
	Boots  EquipmentItem `json:"boots"`
	Chest  EquipmentItem `json:"chest"`
	Gloves EquipmentItem `json:"gloves"`
	Helmet EquipmentItem `json:"helmet"`
	Pants  EquipmentItem `json:"pants"`
	Weapon *WeaponItem   `json:"weapon"`
}

// EquipmentItem is a nested component of CurrentEquipment.
type EquipmentItem struct {
	ID                string          `json:"_id"`
	Code              string          `json:"code"`
	LastAcquisitionAt Time            `json:"lastAcquisitionAt"`
	MaxState          float64         `json:"maxState"`
	Quantity          int             `json:"quantity"`
	Skills            EquipmentSkills `json:"skills"`
	State             float64         `json:"state"`
	Type              string          `json:"type"`
}

// EquipmentSkills is a nested component of RoundEquipment.
type EquipmentSkills struct {
	Armor           float64 `json:"armor"`
	CriticalDamages float64 `json:"criticalDamages"`
	Dodge           float64 `json:"dodge"`
	Precision       float64 `json:"precision"`
}

// Event is a world or country event with event-specific data.
type Event struct {
	ID         string    `json:"_id"`
	CountryIDs []string  `json:"countries"`
	CreatedAt  Time      `json:"createdAt"`
	Data       EventData `json:"data"`
	Priority   string    `json:"priority"`
	UpdatedAt  Time      `json:"updatedAt"`
}

// EventData is a nested component of Event.
type EventData struct {
	AllianceID        string    `json:"allianceId"`
	AllianceName      string    `json:"allianceName"`
	AttackerCountryID string    `json:"attackerCountry"`
	AttackerRegionID  string    `json:"attackerRegion"`
	BattleID          string    `json:"battle"`
	BonusPercent      float64   `json:"bonusPercent"`
	CountryIDs        []string  `json:"countries"`
	CountryID         string    `json:"country"`
	DefenderCountryID string    `json:"defenderCountry"`
	DefenderRegionID  string    `json:"defenderRegion"`
	DurationDays      float64   `json:"durationDays"`
	ItemCode          string    `json:"itemCode"`
	Money             float64   `json:"money"`
	RegionID          string    `json:"region"`
	Type              EventType `json:"type"`
	WarIDs            []string  `json:"wars"`
	WonBy             string    `json:"wonBy"`
}

// GameConfig contains server settings for game systems, skills, items, and rewards.
type GameConfig struct {
	Alliance               GameConfigAlliance               `json:"alliance"`
	Announcement           GameConfigAnnouncement           `json:"announcement"`
	Badge                  map[string]GameConfigBadge       `json:"badge"`
	Battle                 GameConfigBattle                 `json:"battle"`
	CitizenshipApplication GameConfigCitizenshipApplication `json:"citizenshipApplication"`
	Company                GameConfigCompany                `json:"company"`
	Country                GameConfigCountry                `json:"country"`
	Election               GameConfigElection               `json:"election"`
	Government             GameConfigGovernment             `json:"government"`
	Items                  map[string]GameConfigItem        `json:"items"`
	Law                    GameConfigLaw                    `json:"law"`
	Loot                   GameConfigLoot                   `json:"loot"`
	MapTravel              GameConfigMapTravel              `json:"mapTravel"`
	MercenaryContract      GameConfigMercenaryContract      `json:"mercenaryContract"`
	MergingCost            GameConfigMergingCost            `json:"mergingCost"`
	Mission                GameConfigMission                `json:"mission"`
	MU                     GameConfigMu                     `json:"mu"`
	Newspaper              GameConfigNewspaper              `json:"newspaper"`
	Org                    GameConfigOrg                    `json:"org"`
	Party                  GameConfigParty                  `json:"party"`
	Referral               GameConfigReferral               `json:"referral"`
	Region                 GameConfigRegion                 `json:"region"`
	Skills                 map[string]GameConfigSkill       `json:"skills"`
	SubSkinReward          GameConfigSubSkinReward          `json:"subSkinReward"`
	Unrest                 GameConfigUnrest                 `json:"unrest"`
	Upgrade                GameConfigUpgrade                `json:"upgrade"`
	UpgradesConfig         GameConfigUpgradesConfig         `json:"upgradesConfig"`
	User                   GameConfigUser                   `json:"user"`
}

// GameConfigAlliance is a nested component of GameConfig.
type GameConfigAlliance struct {
	LeaveCooldownDays          float64 `json:"leaveCooldownDays"`
	LeaveCooldownDisabledUntil Time    `json:"leaveCooldownDisabledUntil"`
}

// GameConfigAnnouncement is a nested component of GameConfig.
type GameConfigAnnouncement struct {
	DurationHours float64 `json:"durationHours"`
	MaxLength     int     `json:"maxLength"`
	MaxPerHour    float64 `json:"maxPerHour"`
}

// GameConfigBattle is a nested component of GameConfig.
type GameConfigBattle struct {
	AllianceDamagesBonusPercent          float64        `json:"allianceDamagesBonusPercent"`
	CasesPer1kDamagesInPool              float64        `json:"casesPer1kDamagesInPool"`
	CountryOrderBonusPercent             float64        `json:"countryOrderBonusPercent"`
	EnemyDamagesBonusPercent             float64        `json:"enemyDamagesBonusPercent"`
	GovMemberBountyMinPpPer1k            float64        `json:"govMemberBountyMinPpPer1k"`
	GovMemberBountyRewardPercent         float64        `json:"govMemberBountyRewardPercent"`
	HealthCost                           float64        `json:"healthCost"`
	HitFor1CaseInPool                    float64        `json:"hitFor1CaseInPool"`
	LocationGateEnabled                  bool           `json:"locationGateEnabled"`
	LostAttackingRegionMalusPercent      float64        `json:"lostAttackingRegionMalusPercent"`
	MaxRounds                            int            `json:"maxRounds"`
	MuOrderBonusPercent                  float64        `json:"muOrderBonusPercent"`
	OccupyingYourRegionsMalusPercent     float64        `json:"occupyingYourRegionsMalusPercent"`
	PatrioticBonusPercent                float64        `json:"patrioticBonusPercent"`
	PointsToWinRound                     float64        `json:"pointsToWinRound"`
	RankingsLootPercentPer1kDmg          float64        `json:"rankingsLootPercentPer1kDmg"`
	RegionNotLinkedToCapitalMalusPercent float64        `json:"regionNotLinkedToCapitalMalusPercent"`
	RoundsToWin                          int            `json:"roundsToWin"`
	SetCountryOrderMoneyCost             float64        `json:"setCountryOrderMoneyCost"`
	SetMuOrderMoneyCost                  float64        `json:"setMuOrderMoneyCost"`
	SetOrderMoneyCost                    float64        `json:"setOrderMoneyCost"`
	TickPoints                           map[string]int `json:"tickPoints"`
}

// GameConfigCitizenshipApplication is a nested component of GameConfig.
type GameConfigCitizenshipApplication struct {
	AutoApprovalEnabled       bool    `json:"autoApprovalEnabled"`
	AutoApprovalMaxPopulation float64 `json:"autoApprovalMaxPopulation"`
	NewCountryBordersOpenDays float64 `json:"newCountryBordersOpenDays"`
}

// GameConfigCompany is a nested component of GameConfig.
type GameConfigCompany struct {
	ChangeItemCost                     float64 `json:"changeItemCost"`
	ConstructionCostIncreasePerCompany float64 `json:"constructionCostIncreasePerCompany"`
	DepositResourceBonus               float64 `json:"depositResourceBonus"`
	DestructionValuePercent            float64 `json:"destructionValuePercent"`
	MoveCost                           float64 `json:"moveCost"`
}

// GameConfigCountry is a nested component of GameConfig.
type GameConfigCountry struct {
	ConqueredRegionCoreTaxPercent float64 `json:"conqueredRegionCoreTaxPercent"`
	MaxTaxAmount                  float64 `json:"maxTaxAmount"`
}

// GameConfigElection is a nested component of GameConfig.
type GameConfigElection struct {
	CandidateDurationHours                  float64 `json:"candidateDurationHours"`
	CandidateMinLevel                       int     `json:"candidateMinLevel"`
	CongressMaxSeats                        int     `json:"congressMaxSeats"`
	ElectionVoteDurationHours               float64 `json:"electionVoteDurationHours"`
	VoteCitizenshipLockNewCountryMultiplier float64 `json:"voteCitizenshipLockNewCountryMultiplier"`
	VoteMinDaysSinceCitizenshipChange       float64 `json:"voteMinDaysSinceCitizenshipChange"`
	VoteMinLevel                            int     `json:"voteMinLevel"`
}

// GameConfigGovernment is a nested component of GameConfig.
type GameConfigGovernment struct {
	NominationCooldownDays   float64 `json:"nominationCooldownDays"`
	PresidentWagePercent     float64 `json:"presidentWagePercent"`
	VicePresidentWagePercent float64 `json:"vicePresidentWagePercent"`
}

// GameConfigLaw is a nested component of GameConfig.
type GameConfigLaw struct {
	AbusiveLawPossibleVotersNeeded float64                           `json:"abusiveLawPossibleVotersNeeded"`
	AbusiveLawsCooldownInDays      float64                           `json:"abusiveLawsCooldownInDays"`
	Accept_join_alliance           GameConfigLawAccept_join_alliance `json:"accept_join_alliance"`
	Create_alliance                GameConfigLawCreate_alliance      `json:"create_alliance"`
	Define_enemy_country           GameConfigLawDefine_enemy_country `json:"define_enemy_country"`
	LawVotesDurationHours          float64                           `json:"lawVotesDurationHours"`
	Leave_alliance                 GameConfigLawLeave_alliance       `json:"leave_alliance"`
	SendMoneyToCountry             GameConfigLawSendMoneyToCountry   `json:"sendMoneyToCountry"`
	Set_color_scheme               GameConfigLawSet_color_scheme     `json:"set_color_scheme"`
	VotersRatioNeeded              float64                           `json:"votersRatioNeeded"`
}

// GameConfigLawAccept_join_alliance is a nested component of GameConfigLaw.
type GameConfigLawAccept_join_alliance struct {
	Cost float64 `json:"cost"`
}

// GameConfigLawCreate_alliance is a nested component of GameConfigLaw.
type GameConfigLawCreate_alliance struct {
	Cost            float64 `json:"cost"`
	MaintenanceCost float64 `json:"maintenanceCost"`
}

// GameConfigLawDefine_enemy_country is a nested component of GameConfigLaw.
type GameConfigLawDefine_enemy_country struct {
	Cost            float64 `json:"cost"`
	MaintenanceCost float64 `json:"maintenanceCost"`
}

// GameConfigLawLeave_alliance is a nested component of GameConfigLaw.
type GameConfigLawLeave_alliance struct {
	Cost float64 `json:"cost"`
}

// GameConfigLawSendMoneyToCountry is a nested component of GameConfigLaw.
type GameConfigLawSendMoneyToCountry struct {
	AllianceTaxRate float64 `json:"allianceTaxRate"`
	ExternalTaxRate float64 `json:"externalTaxRate"`
}

// GameConfigLawSet_color_scheme is a nested component of GameConfigLaw.
type GameConfigLawSet_color_scheme struct {
	Cost float64 `json:"cost"`
}

// GameConfigLoot is a nested component of GameConfig.
type GameConfigLoot struct {
	BattleLootDamagePerLootItem float64 `json:"battleLootDamagePerLootItem"`
	DamagePerLootItem           float64 `json:"damagePerLootItem"`
	WeaponChancePercent         float64 `json:"weaponChancePercent"`
}

// GameConfigMercenaryContract is a nested component of GameConfig.
type GameConfigMercenaryContract struct {
	AcceptanceFeePercent            float64                               `json:"acceptanceFeePercent"`
	Auction                         GameConfigMercenaryContractAuction    `json:"auction"`
	BountyCooldownMinutes           float64                               `json:"bountyCooldownMinutes"`
	BountyMinActiveCitizens         float64                               `json:"bountyMinActiveCitizens"`
	CancellationPenaltyPercent      float64                               `json:"cancellationPenaltyPercent"`
	Enabled                         bool                                  `json:"enabled"`
	GovMemberBidMaxPpPer1k          float64                               `json:"govMemberBidMaxPpPer1k"`
	MaxBountyPpMultiplier           float64                               `json:"maxBountyPpMultiplier"`
	NationalBountyEnabled           bool                                  `json:"nationalBountyEnabled"`
	NationalBountyMinActiveCitizens float64                               `json:"nationalBountyMinActiveCitizens"`
	Reputation                      GameConfigMercenaryContractReputation `json:"reputation"`
}

// GameConfigMercenaryContractAuction is a nested component of GameConfigMercenaryContract.
type GameConfigMercenaryContractAuction struct {
	BidStep                 float64 `json:"bidStep"`
	MaxDamagePerContract    float64 `json:"maxDamagePerContract"`
	MaxDuration             float64 `json:"maxDuration"`
	MaxPerKMultiplier       float64 `json:"maxPerKMultiplier"`
	MinDuration             float64 `json:"minDuration"`
	MinPerK                 float64 `json:"minPerK"`
	TimerExtensionAmount    float64 `json:"timerExtensionAmount"`
	TimerExtensionThreshold float64 `json:"timerExtensionThreshold"`
}

// GameConfigMercenaryContractReputation is a nested component of GameConfigMercenaryContract.
type GameConfigMercenaryContractReputation struct {
	FailurePerDamage           float64 `json:"failurePerDamage"`
	FailurePerDollar           float64 `json:"failurePerDollar"`
	NegativeBuyAmount          float64 `json:"negativeBuyAmount"`
	NegativeBuyCooldownHours   float64 `json:"negativeBuyCooldownHours"`
	NegativeBuyCost            float64 `json:"negativeBuyCost"`
	ProfessionalsOnlyThreshold float64 `json:"professionalsOnlyThreshold"`
	SuccessPerDollar           float64 `json:"successPerDollar"`
	WeeklyDecayPercent         float64 `json:"weeklyDecayPercent"`
}

// GameConfigMergingCost is a nested component of GameConfig.
type GameConfigMergingCost struct {
	Common    float64 `json:"common"`
	Epic      float64 `json:"epic"`
	Legendary float64 `json:"legendary"`
	Mythic    float64 `json:"mythic"`
	Rare      float64 `json:"rare"`
	Uncommon  float64 `json:"uncommon"`
}

// GameConfigMission is a nested component of GameConfig.
type GameConfigMission struct {
	RerollMissionCost map[string]int          `json:"rerollMissionCost"`
	Reward            GameConfigMissionReward `json:"reward"`
}

// GameConfigMissionReward is a nested component of GameConfigMission.
type GameConfigMissionReward struct {
	Daily    GameConfigMissionRewardDaily    `json:"daily"`
	Monthly  GameConfigMissionRewardMonthly  `json:"monthly"`
	Starting GameConfigMissionRewardStarting `json:"starting"`
	Weekly   GameConfigMissionRewardWeekly   `json:"weekly"`
}

// GameConfigMissionRewardDaily is a nested component of GameConfigMissionReward.
type GameConfigMissionRewardDaily struct {
	Cases          float64 `json:"cases"`
	Money          float64 `json:"money"`
	XP             float64 `json:"xp"`
	XpWhenFinished float64 `json:"xpWhenFinished"`
}

// GameConfigMissionRewardMonthly is a nested component of GameConfigMissionReward.
type GameConfigMissionRewardMonthly struct {
	Cases          float64 `json:"cases"`
	Money          float64 `json:"money"`
	XP             float64 `json:"xp"`
	XpWhenFinished float64 `json:"xpWhenFinished"`
}

// GameConfigMissionRewardStarting is a nested component of GameConfigMissionReward.
type GameConfigMissionRewardStarting struct {
	Cases          float64 `json:"cases"`
	Money          float64 `json:"money"`
	XP             float64 `json:"xp"`
	XpWhenFinished float64 `json:"xpWhenFinished"`
}

// GameConfigMissionRewardWeekly is a nested component of GameConfigMissionReward.
type GameConfigMissionRewardWeekly struct {
	Cases          float64 `json:"cases"`
	Money          float64 `json:"money"`
	XP             float64 `json:"xp"`
	XpWhenFinished float64 `json:"xpWhenFinished"`
}

// GameConfigMu is a nested component of GameConfig.
type GameConfigMu struct {
	ConstructionCost              float64 `json:"constructionCost"`
	DestructionValuePercent       float64 `json:"destructionValuePercent"`
	HealthPerHelp                 float64 `json:"healthPerHelp"`
	HelpCooldownHours             float64 `json:"helpCooldownHours"`
	HelperHealthPerHelp           float64 `json:"helperHealthPerHelp"`
	HelpValue                     float64 `json:"helpValue"`
	MaxOwnedMus                   int     `json:"maxOwnedMus"`
	MoveCost                      float64 `json:"moveCost"`
	NationalityChangeConcreteCost float64 `json:"nationalityChangeConcreteCost"`
	NationalityChangeCooldownDays float64 `json:"nationalityChangeCooldownDays"`
}

// GameConfigNewspaper is a nested component of GameConfig.
type GameConfigNewspaper struct {
	CreateArticleMinLevel int     `json:"createArticleMinLevel"`
	GemTipValue           float64 `json:"gemTipValue"`
	PublishCost           float64 `json:"publishCost"`
	TipMinLevel           int     `json:"tipMinLevel"`
	TipValue              float64 `json:"tipValue"`
}

// GameConfigOrg is a nested component of GameConfig.
type GameConfigOrg struct {
	ConstructionCost float64 `json:"constructionCost"`
	MoveCost         float64 `json:"moveCost"`
}

// GameConfigParty is a nested component of GameConfig.
type GameConfigParty struct {
	CreateCost float64 `json:"createCost"`
}

// GameConfigReferral is a nested component of GameConfig.
type GameConfigReferral struct {
	CanSetReferrerBeforeOrAtLevel  int     `json:"canSetReferrerBeforeOrAtLevel"`
	LevelNeededForBadge            float64 `json:"levelNeededForBadge"`
	LifeTimeBadgeMoneySharePercent float64 `json:"lifeTimeBadgeMoneySharePercent"`
	MoneyForBeingReferred          float64 `json:"moneyForBeingReferred"`
}

// GameConfigRegion is a nested component of GameConfig.
type GameConfigRegion struct {
	BattleCooldownHours                            float64 `json:"battleCooldownHours"`
	DepleteHourlyPercent                           float64 `json:"depleteHourlyPercent"`
	DepositsEnabled                                bool    `json:"depositsEnabled"`
	DevelopmentIncomeMultiplier                    float64 `json:"developmentIncomeMultiplier"`
	DevelopmentMultiplier                          float64 `json:"developmentMultiplier"`
	HomeMoveConcreteCost                           float64 `json:"homeMoveConcreteCost"`
	LiberationDaysCooldown                         float64 `json:"liberationDaysCooldown"`
	LostResistanceOtherRegionsBarMultiplier        float64 `json:"lostResistanceOtherRegionsBarMultiplier"`
	MaxDailyResistance                             float64 `json:"maxDailyResistance"`
	MaxResistance                                  float64 `json:"maxResistance"`
	MinDailyResistance                             float64 `json:"minDailyResistance"`
	NonAggressionHoursAfterLiberation              float64 `json:"nonAggressionHoursAfterLiberation"`
	NonAggressionHoursAfterPeace                   float64 `json:"nonAggressionHoursAfterPeace"`
	ResistanceAllyBonusPercent                     float64 `json:"resistanceAllyBonusPercent"`
	ResistanceBarMultiplier                        float64 `json:"resistanceBarMultiplier"`
	ResistanceBattleCooldownHours                  float64 `json:"resistanceBattleCooldownHours"`
	ResistanceBattleStartCostMultiplier            float64 `json:"resistanceBattleStartCostMultiplier"`
	ResistanceCitizenBonus                         float64 `json:"resistanceCitizenBonus"`
	ResistanceCitizenBonusPercent                  float64 `json:"resistanceCitizenBonusPercent"`
	ResistanceContributionCooldownAfterRevoltHours float64 `json:"resistanceContributionCooldownAfterRevoltHours"`
	ResistanceContributionCost                     float64 `json:"resistanceContributionCost"`
	ResistanceContributionMinLevel                 int     `json:"resistanceContributionMinLevel"`
	ResistanceContributionValue                    float64 `json:"resistanceContributionValue"`
	ResistanceDecayPercent                         float64 `json:"resistanceDecayPercent"`
	ResistanceForeignGovCostMultiplier             float64 `json:"resistanceForeignGovCostMultiplier"`
	ResistanceGrowthPercentMax                     float64 `json:"resistanceGrowthPercentMax"`
	ResistanceGrowthPercentMin                     float64 `json:"resistanceGrowthPercentMin"`
	ResistancePassiveGrowthPercent                 float64 `json:"resistancePassiveGrowthPercent"`
	StrategicResourceIncomeMultiplier              float64 `json:"strategicResourceIncomeMultiplier"`
	TransferDaysCooldown                           float64 `json:"transferDaysCooldown"`
	WonResistanceBarMultiplier                     float64 `json:"wonResistanceBarMultiplier"`
}

// GameConfigSkillLevel describes a skill's cost and value at a given level.
type GameConfigSkillLevel struct {
	Cost          float64 `json:"cost"`
	IsABar        bool    `json:"isABar"`
	TotalCost     float64 `json:"totalCost"`
	UnlockAtLevel int     `json:"unlockAtLevel"`
	Value         float64 `json:"value"`
}

// GameConfigSubSkinReward is a nested component of GameConfig.
type GameConfigSubSkinReward struct {
	Deadline Time   `json:"deadline"`
	SkinKey  string `json:"skinKey"`
}

// GameConfigUnrest is a nested component of GameConfig.
type GameConfigUnrest struct {
	BarMultiplier                            float64 `json:"barMultiplier"`
	BattleCooldownHours                      float64 `json:"battleCooldownHours"`
	BattleStartCost                          float64 `json:"battleStartCost"`
	BordersOpenDays                          float64 `json:"bordersOpenDays"`
	ContributionCooldownAfterRevolutionHours float64 `json:"contributionCooldownAfterRevolutionHours"`
	ContributionCost                         float64 `json:"contributionCost"`
	ContributionMinLevel                     int     `json:"contributionMinLevel"`
	ContributionValue                        float64 `json:"contributionValue"`
	NominationPeriodHours                    float64 `json:"nominationPeriodHours"`
}

// GameConfigUpgrade is a nested component of GameConfig.
type GameConfigUpgrade struct {
	ConstructionContributionMinLevel  int     `json:"constructionContributionMinLevel"`
	ConstructionPointsPerContribution float64 `json:"constructionPointsPerContribution"`
	RefundPercent                     float64 `json:"refundPercent"`
	RegionDowngradeCooldownHours      float64 `json:"regionDowngradeCooldownHours"`
	RegionUpgradeCooldownHours        float64 `json:"regionUpgradeCooldownHours"`
}

// GameConfigUpgradesConfig is a nested component of GameConfig.
type GameConfigUpgradesConfig struct {
	AutomatedEngine    GameConfigUpgradesConfigAutomatedEngine    `json:"automatedEngine"`
	Base               GameConfigUpgradesConfigBase               `json:"base"`
	BreakRoom          GameConfigUpgradesConfigBreakRoom          `json:"breakRoom"`
	Bunker             GameConfigUpgradesConfigBunker             `json:"bunker"`
	Dormitories        GameConfigUpgradesConfigDormitories        `json:"dormitories"`
	Headquarters       GameConfigUpgradesConfigHeadquarters       `json:"headquarters"`
	PacificationCenter GameConfigUpgradesConfigPacificationCenter `json:"pacificationCenter"`
	Storage            GameConfigUpgradesConfigStorage            `json:"storage"`
}

// GameConfigUpgradeLevel describes the costs and effects of an upgrade level.
type GameConfigUpgradeLevel struct {
	ConstructionPointsCost         float64            `json:"constructionPointsCost"`
	Level                          int                `json:"level"`
	MaintenanceCost                float64            `json:"maintenanceCost"`
	MaintenanceCostCountryDevScale float64            `json:"maintenanceCostCountryDevScale"`
	MaintenanceCostRegionDevScale  float64            `json:"maintenanceCostRegionDevScale"`
	MinimumMaintenanceCost         float64            `json:"minimumMaintenanceCost"`
	Stats                          map[string]float64 `json:"stats"`
	SteelCost                      float64            `json:"steelCost"`
}

// GameConfigUpgradesConfigAutomatedEngine is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigAutomatedEngine struct {
	CanDowngrade         bool                              `json:"canDowngrade"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels"`
	PendingDurationHours float64                           `json:"pendingDurationHours"`
}

// GameConfigUpgradesConfigBase is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigBase struct {
	CanBeDestroyed       bool                              `json:"canBeDestroyed"`
	CanBeDisabled        bool                              `json:"canBeDisabled"`
	CanDowngrade         bool                              `json:"canDowngrade"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels"`
	PendingDurationHours float64                           `json:"pendingDurationHours"`
}

// GameConfigUpgradesConfigBreakRoom is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigBreakRoom struct {
	CanDowngrade bool                              `json:"canDowngrade"`
	Levels       map[string]GameConfigUpgradeLevel `json:"levels"`
}

// GameConfigUpgradesConfigBunker is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigBunker struct {
	CanBeDestroyed       bool                              `json:"canBeDestroyed"`
	CanBeDisabled        bool                              `json:"canBeDisabled"`
	CanDowngrade         bool                              `json:"canDowngrade"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels"`
	PendingDurationHours float64                           `json:"pendingDurationHours"`
}

// GameConfigUpgradesConfigDormitories is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigDormitories struct {
	CanDowngrade bool                              `json:"canDowngrade"`
	Levels       map[string]GameConfigUpgradeLevel `json:"levels"`
}

// GameConfigUpgradesConfigHeadquarters is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigHeadquarters struct {
	CanBeDisabled        bool                              `json:"canBeDisabled"`
	CanDowngrade         bool                              `json:"canDowngrade"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels"`
	PendingDurationHours float64                           `json:"pendingDurationHours"`
}

// GameConfigUpgradesConfigPacificationCenter is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigPacificationCenter struct {
	CanBeDestroyed       bool                              `json:"canBeDestroyed"`
	CanBeDisabled        bool                              `json:"canBeDisabled"`
	CanDowngrade         bool                              `json:"canDowngrade"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels"`
	PendingDurationHours float64                           `json:"pendingDurationHours"`
}

// GameConfigUpgradesConfigStorage is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigStorage struct {
	CanDowngrade         bool                              `json:"canDowngrade"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels"`
	PendingDurationHours float64                           `json:"pendingDurationHours"`
}

// GameConfigUser is a nested component of GameConfig.
type GameConfigUser struct {
	ActiveCitizenMinLevel     int                         `json:"activeCitizenMinLevel"`
	CanTakeControlAtLevel     int                         `json:"canTakeControlAtLevel"`
	ChatSlowmode              GameConfigChatSlowmode      `json:"chatSlowmode"`
	CitizenshipDaysCooldown   float64                     `json:"citizenshipDaysCooldown"`
	ConstructionEnergyCost    float64                     `json:"constructionEnergyCost"`
	DailyReward               GameConfigUserDailyReward   `json:"dailyReward"`
	DailyXp                   float64                     `json:"dailyXp"`
	DonationMinLevel          int                         `json:"donationMinLevel"`
	EnergyCostPerAction       float64                     `json:"energyCostPerAction"`
	EquipmentSets             GameConfigUserEquipmentSets `json:"equipmentSets"`
	FeatureUnlockLevels       map[string]int              `json:"featureUnlockLevels"`
	FieldsToPopulate          string                      `json:"fieldsToPopulate"`
	IsInactiveAfterDays       float64                     `json:"isInactiveAfterDays"`
	Kits                      GameConfigUserKits          `json:"kits"`
	MarketMinLevel            int                         `json:"marketMinLevel"`
	MaxConstructionPoints     float64                     `json:"maxConstructionPoints"`
	MaxEnergy                 float64                     `json:"maxEnergy"`
	MaxHunger                 float64                     `json:"maxHunger"`
	Prestige                  GameConfigUserPrestige      `json:"prestige"`
	RegenDividedBy            float64                     `json:"regenDividedBy"`
	ResetSkillDaysCooldown    float64                     `json:"resetSkillDaysCooldown"`
	ResetSkillsCostPerPoint   float64                     `json:"resetSkillsCostPerPoint"`
	SkillPresets              GameConfigSkillPresets      `json:"skillPresets"`
	TakeControlCooldownInDays float64                     `json:"takeControlCooldownInDays"`
	XpPerAction               float64                     `json:"xpPerAction"`
}

// GameConfigUserDailyReward is a nested component of GameConfigUser.
type GameConfigUserDailyReward struct {
	Case1 float64 `json:"case1"`
	Money float64 `json:"money"`
	XP    float64 `json:"xp"`
}

// GameConfigUserEquipmentSets is a nested component of GameConfigUser.
type GameConfigUserEquipmentSets struct {
	NonPremiumMax float64 `json:"nonPremiumMax"`
	PremiumMax    float64 `json:"premiumMax"`
}

// GameConfigUserKits is a nested component of GameConfigUser.
type GameConfigUserKits struct {
	NonPremiumMax float64 `json:"nonPremiumMax"`
	PremiumMax    float64 `json:"premiumMax"`
}

// GameConfigUserPrestige is a nested component of GameConfigUser.
type GameConfigUserPrestige struct {
	Enabled                     bool    `json:"enabled"`
	EnabledLevel                int     `json:"enabledLevel"`
	SkillBonusLevels            int     `json:"skillBonusLevels"`
	XpMultiplierPerLevelPercent float64 `json:"xpMultiplierPerLevelPercent"`
	XpMultiplierPercent         float64 `json:"xpMultiplierPercent"`
	XpPenalty                   float64 `json:"xpPenalty"`
}

// GameDates contains server timestamps for daily rollover and energy regeneration.
type GameDates struct {
	DailyMissionRegenAt         Time `json:"dailyMissionRegenAt"`
	MonthlyMissionRegenAt       Time `json:"monthlyMissionRegenAt"`
	NextCongressElectionsAt     Time `json:"nextCongressElectionsAt"`
	NextDayAt                   Time `json:"nextDayAt"`
	NextMonthAt                 Time `json:"nextMonthAt"`
	NextPresidentialElectionsAt Time `json:"nextPresidentialElectionsAt"`
	NextRegenAt                 Time `json:"nextRegenAt"`
	PreviousDayAt               Time `json:"previousDayAt"`
	WeeklyMissionRegenAt        Time `json:"weeklyMissionRegenAt"`
}

// Government contains the leadership and legislative positions for a country.
type Government struct {
	ID                string               `json:"_id"`
	CongressMemberIDs []string             `json:"congressMembers"`
	CountryID         string               `json:"country"`
	Dates             GovernmentDates      `json:"dates"`
	Ministers         []GovernmentMinister `json:"ministers"`
	PresidentID       string               `json:"president"`
	VicePresidentID   string               `json:"vicePresident"`
}

// GovernmentDates is a nested component of Government.
type GovernmentDates struct {
	AnnouncementCreatedAts []Time `json:"announcementCreatedAts"`
	LastAnnouncementAt     Time   `json:"lastAnnouncementAt"`
}

// LiveBattleData represents data returned by the WarEra API.
type LiveBattleData struct {
	Battle LiveBattleSummary `json:"battle"`
	Round  LiveRoundSummary  `json:"round"`
}

// LiveBattleSummary is a nested component of LiveBattleData.
type LiveBattleSummary struct {
	AttackerCountryOrders     []string             `json:"attackerCountryOrders"`
	AttackerMoneyPer1kDamages float64              `json:"attackerMoneyPer1kDamages"`
	DefenderCountryOrders     []string             `json:"defenderCountryOrders"`
	DefenderMoneyPer1kDamages float64              `json:"defenderMoneyPer1kDamages"`
	IsActive                  bool                 `json:"isActive"`
	RoundHistory              []BattleRoundHistory `json:"roundHistory"`
	RoundIDs                  []string             `json:"roundIds"`
}

// LiveRoundSummary is a nested component of LiveBattleData.
type LiveRoundSummary struct {
	ActualTickPoints float64 `json:"actualTickPoints"`
	AttackerDamages  float64 `json:"attackerDamages"`
	AttackerPoints   float64 `json:"attackerPoints"`
	DefenderDamages  float64 `json:"defenderDamages"`
	DefenderPoints   float64 `json:"defenderPoints"`
	IsActive         bool    `json:"isActive"`
	NextTickAt       Time    `json:"nextTickAt"`
	RoundID          string  `json:"roundId"`
}

// MULeveling is a nested component of MilitaryUnit.
type MULeveling struct {
	Level          int     `json:"level"`
	MonthlyDamages float64 `json:"monthlyDamages"`
}

// MURoles is a nested component of MilitaryUnit.
type MURoles struct {
	CommanderIDs []string `json:"commanders"`
	ManagerIDs   []string `json:"managers"`
}

// MilitaryUnit is a player-managed military organization.
type MilitaryUnit struct {
	ID                   string                          `json:"_id"`
	ActiveUpgradeLevels  MilitaryUnitActiveUpgradeLevels `json:"activeUpgradeLevels"`
	AvatarURL            string                          `json:"avatarUrl"`
	CountryID            string                          `json:"country"`
	CreatedAt            Time                            `json:"createdAt"`
	InvestedMoneyByUsers map[string]float64              `json:"investedMoneyByUsers"`
	LastAnnouncementAt   Time                            `json:"lastAnnouncementAt"`
	Leveling             MULeveling                      `json:"leveling"`
	MemberIDs            []string                        `json:"members"`
	MercenaryReputation  float64                         `json:"mercenaryReputation"`
	Name                 string                          `json:"name"`
	Rankings             map[string]RankingItem          `json:"rankings"`
	RegionID             string                          `json:"region"`
	Roles                MURoles                         `json:"roles"`
	SponsorCountryID     string                          `json:"sponsorCountry"`
	Type                 string                          `json:"type"`
	UpdatedAt            Time                            `json:"updatedAt"`
	UserID               string                          `json:"user"`
}

// MilitaryUnitActiveUpgradeLevels is a nested component of MilitaryUnit.
type MilitaryUnitActiveUpgradeLevels struct {
	Dormitories  int `json:"dormitories"`
	Headquarters int `json:"headquarters"`
}

// Ranking contains a global ranking and its entries.
type Ranking struct {
	Count      int               `json:"count"`
	Items      []RankingEntry    `json:"items"`
	TierValues RankingTierValues `json:"tierValues"`
	Type       string            `json:"type"`
	UpdatedAt  Time              `json:"updatedAt"`
}

// RankingItem represents a tiered leaderboard rank entry.
type RankingItem struct {
	Rank  int     `json:"rank"`
	Tier  string  `json:"tier"`
	Value float64 `json:"value"`
}

// RankingEntry is one user, country, military unit, or alliance in a ranking.
// Which entity IDs are set depends on the ranking type; user rankings also
// carry the user's country and military unit.
type RankingEntry struct {
	ID         string  `json:"_id"`
	AllianceID string  `json:"alliance"`
	CountryID  string  `json:"country"`
	MUID       string  `json:"mu"`
	Rank       int     `json:"rank"`
	Tier       string  `json:"tier"`
	Type       string  `json:"type"`
	UserID     string  `json:"user"`
	Value      float64 `json:"value"`
}

// RankingTierValues is a nested component of Ranking.
type RankingTierValues struct {
	Bronze   float64 `json:"bronze"`
	Diamond  float64 `json:"diamond"`
	Gold     float64 `json:"gold"`
	Master   float64 `json:"master"`
	Platinum float64 `json:"platinum"`
	Silver   float64 `json:"silver"`
}

// Region is a map area with economic, geographic, and military state.
type Region struct {
	ID                           string                    `json:"_id"`
	ActiveBattle                 *RegionActiveBattleValue  `json:"activeBattle"`
	ActiveUpgradeLevels          RegionActiveUpgradeLevels `json:"activeUpgradeLevels"`
	BaseDevelopment              float64                   `json:"baseDevelopment"`
	Biome                        string                    `json:"biome"`
	Climate                      string                    `json:"climate"`
	Code                         string                    `json:"code"`
	CompanyCount                 int                       `json:"companyCount"`
	CountryCode                  string                    `json:"countryCode"`
	CountryID                    string                    `json:"country"`
	CurrentPopulation            int                       `json:"currentPopulation"`
	Dates                        RegionDates               `json:"dates"`
	Development                  float64                   `json:"development"`
	EmployeeCount                int                       `json:"employeeCount"`
	HasCoast                     bool                      `json:"hasCoast"`
	Income                       *IncomeBreakdown          `json:"income"`
	InitialCountryID             string                    `json:"initialCountry"`
	IsCapital                    bool                      `json:"isCapital"`
	IsLinkedToCapital            bool                      `json:"isLinkedToCapital"`
	LastBattleEndedAt            Time                      `json:"lastBattleEndedAt"`
	LastResistanceContributionAt Time                      `json:"lastResistanceContributionAt"`
	LastRevoltEndedAt            Time                      `json:"lastRevoltEndedAt"`
	MainCity                     string                    `json:"mainCity"`
	Name                         string                    `json:"name"`
	NeighborIDs                  []string                  `json:"neighbors"`
	Population                   int                       `json:"population"`
	Position                     []float64                 `json:"position"`
	Residents                    int                       `json:"residents"`
	Resistance                   float64                   `json:"resistance"`
	ResistanceMax                float64                   `json:"resistanceMax"`
	Stats                        RegionStats               `json:"stats"`
	StrategicResource            string                    `json:"strategicResource"`
	UpgradesV2                   RegionUpgradesV2          `json:"upgradesV2"`
}

// RegionActiveBattle is a nested component of Region.
type RegionActiveBattle struct {
	ID            string                  `json:"_id"`
	Attacker      BattleSideSummary       `json:"attacker"`
	CreatedAt     Time                    `json:"createdAt"`
	CurrentRound  *BattleRound            `json:"currentRound"`
	Defender      BattleSideSummary       `json:"defender"`
	IsActive      bool                    `json:"isActive"`
	IsBigBattle   bool                    `json:"isBigBattle"`
	RoundIDs      []string                `json:"rounds"`
	RoundsHistory []BattleRoundHistory    `json:"roundsHistory"`
	RoundsToWin   int                     `json:"roundsToWin"`
	Stats         RegionActiveBattleStats `json:"stats"`
	Type          string                  `json:"type"`
	UpdatedAt     Time                    `json:"updatedAt"`
	WarID         string                  `json:"war"`
}

// RegionActiveBattleRoundsHistory is an alias for [BattleRoundHistory].
type RegionActiveBattleRoundsHistory = BattleRoundHistory

// RegionActiveBattleStats is a nested component of RegionActiveBattle.
type RegionActiveBattleStats struct {
	HitCount int `json:"hitCount"`
}

// RegionActiveUpgradeLevels is a nested component of Region.
type RegionActiveUpgradeLevels struct {
	Base               int `json:"base"`
	Bunker             int `json:"bunker"`
	PacificationCenter int `json:"pacificationCenter"`
}

// RegionDates is a nested component of Region.
type RegionDates struct {
	LastOwnershipChangeAt Time `json:"lastOwnershipChangeAt"`
}

// RegionDeposit is a nested component of Region.
type RegionDeposit struct {
	BonusPercent float64 `json:"bonusPercent"`
	EndsAt       Time    `json:"endsAt"`
	StartsAt     Time    `json:"startsAt"`
	Type         string  `json:"type"`
}

// RegionStats is a nested component of Region.
type RegionStats struct {
	InvestedMoney float64 `json:"investedMoney"`
}

// RegionUpgradesV2 is a nested component of Region.
type RegionUpgradesV2 struct {
	ActiveConstructionCount int                      `json:"activeConstructionCount"`
	Upgrades                RegionUpgradesV2Upgrades `json:"upgrades"`
}

// RegionUpgradesV2Upgrades is a nested component of RegionUpgradesV2.
type RegionUpgradesV2Upgrades struct {
	Base   RegionUpgradesV2UpgradesBase   `json:"base"`
	Bunker RegionUpgradesV2UpgradesBunker `json:"bunker"`
}

// RegionUpgradesV2UpgradesBase is a nested component of RegionUpgradesV2Upgrades.
type RegionUpgradesV2UpgradesBase struct {
	ConstructionEndedAt   Time                            `json:"constructionEndedAt"`
	ConstructionPoints    float64                         `json:"constructionPoints"`
	ConstructionStartedAt Time                            `json:"constructionStartedAt"`
	InvestedMoney         float64                         `json:"investedMoney"`
	IsUnderConstruction   *bool                           `json:"isUnderConstruction"`
	LastConstructions     []RegionUpgradeLastConstruction `json:"lastConstructions"`
	Level                 int                             `json:"level"`
	Status                string                          `json:"status"`
	StatusChangedAt       Time                            `json:"statusChangedAt"`
}

// RegionUpgradesV2UpgradesBunker is a nested component of RegionUpgradesV2Upgrades.
type RegionUpgradesV2UpgradesBunker struct {
	ConstructionEndedAt   Time                            `json:"constructionEndedAt"`
	ConstructionPoints    float64                         `json:"constructionPoints"`
	ConstructionStartedAt Time                            `json:"constructionStartedAt"`
	InvestedMoney         float64                         `json:"investedMoney"`
	IsUnderConstruction   *bool                           `json:"isUnderConstruction"`
	LastConstructions     []RegionUpgradeLastConstruction `json:"lastConstructions"`
	Level                 int                             `json:"level"`
	Status                string                          `json:"status"`
	StatusChangedAt       Time                            `json:"statusChangedAt"`
}

// RegionUpgradeLastConstruction records a completed upgrade construction.
type RegionUpgradeLastConstruction struct {
	Construction   float64 `json:"construction"`
	ConstructionAt Time    `json:"constructionAt"`
	UserID         string  `json:"user"`
}

// RegionUpgradesV2UpgradesBunkerLastConstructions is an alias for [RegionUpgradeLastConstruction].
type RegionUpgradesV2UpgradesBunkerLastConstructions = RegionUpgradeLastConstruction

// Round is one stage of a battle and includes combat state and results.
type Round struct {
	ID        string               `json:"_id"`
	Attacker  RoundSideCombatStats `json:"attacker"`
	BattleID  string               `json:"battle"`
	CreatedAt Time                 `json:"createdAt"`
	Defender  RoundSideCombatStats `json:"defender"`
	IsActive  bool                 `json:"isActive"`
	Live      RoundLiveState       `json:"live"`
	Number    int                  `json:"number"`
	UpdatedAt Time                 `json:"updatedAt"`
}

// RoundEquipment is a nested component of RoundLastHit.
type RoundEquipment struct {
	ID                string          `json:"_id"`
	Code              string          `json:"code"`
	LastAcquisitionAt Time            `json:"lastAcquisitionAt"`
	MaxState          float64         `json:"maxState"`
	Quantity          int             `json:"quantity"`
	Skills            EquipmentSkills `json:"skills"`
	State             float64         `json:"state"`
	Type              string          `json:"type"`
}

// RoundLastHit is a nested component of RoundSideCombatStats.
type RoundLastHit struct {
	ID            string           `json:"_id"`
	Ammo          string           `json:"ammo"`
	Damages       float64          `json:"damages"`
	Equipments    []RoundEquipment `json:"equipments"`
	HitAt         Time             `json:"hitAt"`
	IsCriticalHit bool             `json:"isCriticalHit"`
	IsMissed      bool             `json:"isMissed"`
	MUID          string           `json:"mu"`
	UserID        string           `json:"user"`
	Weapon        *RoundWeapon     `json:"weapon"`
}

// RoundLastHits represents data returned by the WarEra API.
type RoundLastHits struct {
	Attacker []RoundLastHit `json:"attacker"`
	Defender []RoundLastHit `json:"defender"`
}

// RoundLiveState is a nested component of BattleRound.
type RoundLiveState struct {
	ActualTickPoints float64 `json:"actualTickPoints"`
	NextTickAt       Time    `json:"nextTickAt"`
	TicksCount       int     `json:"ticksCount"`
}

// RoundSideCombatStats is a nested component of BattleRound.
type RoundSideCombatStats struct {
	CountryID string         `json:"country"`
	Damages   float64        `json:"damages"`
	HitCount  int            `json:"hitCount"`
	LastHits  []RoundLastHit `json:"lastHits"`
	Points    float64        `json:"points"`
}

// RoundWeapon is a nested component of RoundLastHit.
type RoundWeapon struct {
	ID                string       `json:"_id"`
	Code              string       `json:"code"`
	LastAcquisitionAt Time         `json:"lastAcquisitionAt"`
	MaxState          float64      `json:"maxState"`
	Quantity          int          `json:"quantity"`
	Skills            WeaponSkills `json:"skills"`
	State             float64      `json:"state"`
}

// SearchResult contains matching entity IDs grouped by entity type.
type SearchResult struct {
	AllianceIDs []string `json:"allianceIds"`
	CountryIDs  []string `json:"countryIds"`
	HasData     bool     `json:"hasData"`
	MUIDs       []string `json:"muIds"`
	PartyIDs    []string `json:"partyIds"`
	RegionIDs   []string `json:"regionIds"`
	UserIDs     []string `json:"userIds"`
}

// TopOrder is a nested component of TopOrders.
type TopOrder struct {
	ID        string  `json:"_id"`
	CountryID string  `json:"country"`
	ItemCode  string  `json:"itemCode"`
	MUID      string  `json:"mu"`
	OfferAt   Time    `json:"offerAt"`
	PartyID   string  `json:"party"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Type      string  `json:"type"`
	UserID    string  `json:"user"`
}

// TopOrders contains the best buy and sell orders for a tradeable item.
type TopOrders struct {
	BuyOrders  []TopOrder `json:"buyOrders"`
	SellOrders []TopOrder `json:"sellOrders"`
}

// Upgrade represents data returned by the WarEra API.
type Upgrade struct {
	ID                  string      `json:"_id"`
	CreatedAt           Time        `json:"createdAt"`
	DependantUsersCount int         `json:"dependantUsersCount"`
	InvestedConcrete    float64     `json:"investedConcrete"`
	InvestedMoney       float64     `json:"investedMoney"`
	InvestedSteel       float64     `json:"investedSteel"`
	LastDowngradeAt     Time        `json:"lastDowngradeAt"`
	LastUpgradeAt       Time        `json:"lastUpgradeAt"`
	Level               int         `json:"level"`
	RegionID            string      `json:"region"`
	Status              string      `json:"status"`
	StatusChangedAt     Time        `json:"statusChangedAt"`
	UpdatedAt           Time        `json:"updatedAt"`
	UpgradeType         UpgradeType `json:"upgradeType"`
	WillBeActiveAt      Time        `json:"willBeActiveAt"`
}

// User is a full player profile, including skills, statistics, and equipment.
type User struct {
	ID                    string                 `json:"_id"`
	AvailableColorSchemes []string               `json:"availableColorSchemes"`
	AvatarURL             string                 `json:"avatarUrl"`
	Buffs                 *UserBuffs             `json:"buffs"`
	CanOnboard            bool                   `json:"canOnboard"`
	CountryID             string                 `json:"country"`
	CreatedAt             Time                   `json:"createdAt"`
	Dates                 UserDates              `json:"dates"`
	Discord               UserDiscord            `json:"discord"`
	EmailVerified         bool                   `json:"emailVerified"`
	EnergyGift            *UserEnergyGift        `json:"energyGift"`
	Equipment             UserEquipment          `json:"equipment"`
	EquippedSkinKeys      json.RawMessage        `json:"equippedSkinKeys"`
	FinishedTours         UserFinishedTours      `json:"finishedTours"`
	Infos                 UserInfos              `json:"infos"`
	IsActive              bool                   `json:"isActive"`
	Leveling              UserLeveling           `json:"leveling"`
	LocationID            string                 `json:"location"`
	MilitaryRank          int                    `json:"militaryRank"`
	Missions              UserMissions           `json:"missions"`
	MUID                  string                 `json:"mu"`
	MuMaxLevelRewarded    float64                `json:"muMaxLevelRewarded"`
	PartyID               string                 `json:"party"`
	Preferences           UserPreferences        `json:"preferences"`
	Rankings              map[string]RankingItem `json:"rankings"`
	RegionID              string                 `json:"region"`
	ShouldUpdateProfile   bool                   `json:"shouldUpdateProfile"`
	Skills                UserSkills             `json:"skills"`
	SkillVersion          int                    `json:"skillVersion"`
	Stats                 UserStats              `json:"stats"`
	UpdatedAt             Time                   `json:"updatedAt"`
	Username              string                 `json:"username"`
	UsernameLower         string                 `json:"usernameLower"`
	Warning               *UserWarning           `json:"warning"`
}

// UserDates is a nested component of User.
type UserDates struct {
	CitizenSince               Time `json:"citizenSince"`
	LastCitizenshipChangeAt    Time `json:"lastCitizenshipChangeAt"`
	LastConnectionAt           Time `json:"lastConnectionAt"`
	LastCountryMessageCheckAt  Time `json:"lastCountryMessageCheckAt"`
	LastDailyCalendarClaimedAt Time `json:"lastDailyCalendarClaimedAt"`
	LastEventsCheckAt          Time `json:"lastEventsCheckAt"`
	LastFiredFromGovAt         Time `json:"lastFiredFromGovAt"`
	LastFreeTravelHomeAt       Time `json:"lastFreeTravelHomeAt"`
	LastGlobalMessageCheckAt   Time `json:"lastGlobalMessageCheckAt"`
	LastHelpAskedAt            Time `json:"lastHelpAskedAt"`
	LastManagerCareAt          Time `json:"lastManagerCareAt"`
	LastNotificationsCheckAt   Time `json:"lastNotificationsCheckAt"`
	LastSeenAnnouncementAt     Time `json:"lastSeenAnnouncementAt"`
	LastSkillsResetAt          Time `json:"lastSkillsResetAt"`
	LastWorkAt                 Time `json:"lastWorkAt"`
}

// UserDiscord is a nested component of User.
type UserDiscord struct {
	IsVerified bool `json:"isVerified"`
}

// UserEquipment is a nested component of User.
type UserEquipment struct {
	Ammo     string `json:"ammo"`
	BootsID  string `json:"boots"`
	ChestID  string `json:"chest"`
	GlovesID string `json:"gloves"`
	HelmetID string `json:"helmet"`
	PantsID  string `json:"pants"`
	WeaponID string `json:"weapon"`
}

// UserFinishedTours is a nested component of User.
type UserFinishedTours struct {
	Onboarding bool `json:"onboarding"`
}

// UserLeveling is a nested component of User.
type UserLeveling struct {
	AvailableSkillPoints int     `json:"availableSkillPoints"`
	DailyXPLeft          float64 `json:"dailyXpLeft"`
	FreeReset            float64 `json:"freeReset"`
	Level                int     `json:"level"`
	Prestige             int     `json:"prestige"`
	PrestigeLevel        int     `json:"prestigeLevel"`
	SpentSkillPoints     int     `json:"spentSkillPoints"`
	TotalSkillPoints     int     `json:"totalSkillPoints"`
	TotalXP              float64 `json:"totalXp"`
	XP                   float64 `json:"xp"`
}

// UserLite is a compact player profile for listings and search results.
type UserLite struct {
	ID               string                 `json:"_id"`
	AvatarURL        string                 `json:"avatarUrl"`
	Buffs            *UserBuffs             `json:"buffs"`
	CountryID        string                 `json:"country"`
	CreatedAt        Time                   `json:"createdAt"`
	Dates            UserDates              `json:"dates"`
	EquippedSkinKeys json.RawMessage        `json:"equippedSkinKeys"`
	Infos            UserInfos              `json:"infos"`
	IsActive         bool                   `json:"isActive"`
	Leveling         UserLeveling           `json:"leveling"`
	MilitaryRank     int                    `json:"militaryRank"`
	MUID             string                 `json:"mu"`
	Rankings         map[string]RankingItem `json:"rankings"`
	Stats            UserStats              `json:"stats"`
	Username         string                 `json:"username"`
}

// UserMissions is a nested component of User.
type UserMissions struct {
	ClaimedAt               UserMissionsClaimedAt   `json:"claimedAt"`
	ClaimedCount            int                     `json:"claimedCount"`
	GeneratedAt             UserMissionsGeneratedAt `json:"generatedAt"`
	RerolledDailyMissions   int                     `json:"rerolledDailyMissions"`
	RerolledMonthlyMissions int                     `json:"rerolledMonthlyMissions"`
	RerolledWeeklyMissions  int                     `json:"rerolledWeeklyMissions"`
	TotalXpClaimed          float64                 `json:"totalXpClaimed"`
}

// UserMissionsClaimedAt is a nested component of UserMissions.
type UserMissionsClaimedAt struct {
	Daily    Time `json:"daily"`
	Monthly  Time `json:"monthly"`
	Starting Time `json:"starting"`
	Weekly   Time `json:"weekly"`
}

// UserPreferences is a nested component of User.
type UserPreferences struct {
	AppFont            string `json:"appFont"`
	AppPattern         string `json:"appPattern"`
	AutoEquipMode      string `json:"autoEquipMode"`
	AutoReplaceOnBreak bool   `json:"autoReplaceOnBreak"`
	Colorblind         bool   `json:"colorblind"`
	Locale             string `json:"locale"`
	Sfx                bool   `json:"sfx"`
	Vibrations         bool   `json:"vibrations"`
}

// UserSkills contains a player's skills.
type UserSkills struct {
	Armor           UserSkill `json:"armor"`
	Attack          UserSkill `json:"attack"`
	Companies       UserSkill `json:"companies"`
	CriticalChance  UserSkill `json:"criticalChance"`
	CriticalDamages UserSkill `json:"criticalDamages"`
	Dodge           UserSkill `json:"dodge"`
	Energy          UserSkill `json:"energy"`
	Health          UserSkill `json:"health"`
	Hunger          UserSkill `json:"hunger"`
	LootChance      UserSkill `json:"lootChance"`
	Management      UserSkill `json:"management"`
	Precision       UserSkill `json:"precision"`
	Production      UserSkill `json:"production"`
	Stamina         UserSkill `json:"stamina"`
}

// UserStats is a nested component of User.
type UserStats struct {
	Case1                 UserStatsCase1      `json:"case1"`
	Case2                 UserStatsCase2      `json:"case2"`
	DailyRewardBestStreak int                 `json:"dailyRewardBestStreak"`
	DailyRewardStreak     int                 `json:"dailyRewardStreak"`
	DamagesCount          int                 `json:"damagesCount"`
	MilestoneTiersClaimed int                 `json:"milestoneTiersClaimed"`
	MilestoneXpClaimed    float64             `json:"milestoneXpClaimed"`
	Toolbox               UserStatsToolbox    `json:"toolbox"`
	Wealth                UserStatsWealth     `json:"wealth"`
	WoodenCase            UserStatsWoodenCase `json:"woodenCase"`
	WorksCount            int                 `json:"worksCount"`
}

// UserStatsCase1 is a nested component of UserStats.
type UserStatsCase1 struct {
	ByRarity    UserStatsCase1ByRarity `json:"byRarity"`
	OpenedCount int                    `json:"openedCount"`
}

// UserStatsCase1ByRarity is a nested component of UserStatsCase1.
type UserStatsCase1ByRarity struct {
	Common    float64 `json:"common"`
	Epic      float64 `json:"epic"`
	Legendary float64 `json:"legendary"`
	Mythic    float64 `json:"mythic"`
	Rare      float64 `json:"rare"`
	Uncommon  float64 `json:"uncommon"`
}

// UserStatsCase2 is a nested component of UserStats.
type UserStatsCase2 struct {
	ByRarity    UserStatsCase2ByRarity `json:"byRarity"`
	OpenedCount int                    `json:"openedCount"`
}

// UserStatsCase2ByRarity is a nested component of UserStatsCase2.
type UserStatsCase2ByRarity struct {
	Epic      float64 `json:"epic"`
	Legendary float64 `json:"legendary"`
	Rare      float64 `json:"rare"`
	Uncommon  float64 `json:"uncommon"`
}

// UserStatsWealth is a nested component of UserStats.
type UserStatsWealth struct {
	Companies  float64 `json:"companies"`
	Equipments float64 `json:"equipments"`
	Items      float64 `json:"items"`
	Money      float64 `json:"money"`
	Total      float64 `json:"total"`
	Weapons    float64 `json:"weapons"`
}

// UserStatsWoodenCase is a nested component of UserStats.
type UserStatsWoodenCase struct {
	ByItem      UserStatsWoodenCaseByItem `json:"byItem"`
	OpenedCount int                       `json:"openedCount"`
}

// UserStatsWoodenCaseByItem is a nested component of UserStatsWoodenCase.
type UserStatsWoodenCaseByItem struct {
	Ammo       float64 `json:"ammo"`
	Bread      float64 `json:"bread"`
	Coca       float64 `json:"coca"`
	Cocain     float64 `json:"cocain"`
	Concrete   float64 `json:"concrete"`
	CookedFish float64 `json:"cookedFish"`
	Fish       float64 `json:"fish"`
	Grain      float64 `json:"grain"`
	HeavyAmmo  float64 `json:"heavyAmmo"`
	Iron       float64 `json:"iron"`
	Lead       float64 `json:"lead"`
	LightAmmo  float64 `json:"lightAmmo"`
	Limestone  float64 `json:"limestone"`
	Livestock  float64 `json:"livestock"`
	Oil        float64 `json:"oil"`
	Paper      float64 `json:"paper"`
	Petroleum  float64 `json:"petroleum"`
	Steak      float64 `json:"steak"`
	Steel      float64 `json:"steel"`
	Wood       float64 `json:"wood"`
}

// WeaponItem is a nested component of CurrentEquipment.
type WeaponItem struct {
	ID                string       `json:"_id"`
	Code              string       `json:"code"`
	LastAcquisitionAt Time         `json:"lastAcquisitionAt"`
	MaxState          float64      `json:"maxState"`
	Quantity          int          `json:"quantity"`
	Skills            WeaponSkills `json:"skills"`
	State             float64      `json:"state"`
}

// WeaponSkills is a nested component of RoundWeapon.
type WeaponSkills struct {
	Attack         float64 `json:"attack"`
	CriticalChance float64 `json:"criticalChance"`
}

// UserSkill is one skill of a player. Fight skills report their weapon,
// equipment, and percentage modifiers; bar skills such as energy, health,
// hunger, and stamina report their current bar value and hourly regeneration.
type UserSkill struct {
	AmmoPercent         float64         `json:"ammoPercent"`
	BuffsPercent        float64         `json:"buffsPercent"`
	CurrentBarValue     float64         `json:"currentBarValue"`
	DebuffsPercent      float64         `json:"debuffsPercent"`
	Equipment           *float64        `json:"equipment"`
	HourlyBarRegen      float64         `json:"hourlyBarRegen"`
	Level               int             `json:"level"`
	Limited             json.RawMessage `json:"limited"`
	MilitaryRankPercent float64         `json:"militaryRankPercent"`
	Overflow            json.RawMessage `json:"overflow"`
	Prestige            int             `json:"prestige"`
	// Strategic is the bonus from strategic resources held by the player's country.
	Strategic         float64  `json:"strategic"`
	Total             float64  `json:"total"`
	TotalAfterSoftCap *float64 `json:"totalAfterSoftCap"`
	Value             float64  `json:"value"`
	Weapon            *float64 `json:"weapon"`
}

// IncomeBreakdown is a country's or region's daily income by source.
type IncomeBreakdown struct {
	Development             float64 `json:"development"`
	EmployeeIncomeTax       float64 `json:"employeeIncomeTax"`
	GovType                 float64 `json:"govType"`
	ManagerIncomeTax        float64 `json:"managerIncomeTax"`
	StrategicResourcesBonus float64 `json:"strategicResourcesBonus"`
	Total                   float64 `json:"total"`
}

// GovernmentMinister is a minister appointed to a country's government.
type GovernmentMinister struct {
	AppointedAt    Time `json:"appointedAt"`
	BadgeGrantedAt Time `json:"badgeGrantedAt"`
	// CustomTitle is a display title the president gave the minister.
	CustomTitle string   `json:"customTitle"`
	Rights      []string `json:"rights"`
	Title       string   `json:"title"`
	UserID      string   `json:"user"`
}

// GameConfigMapTravel reports which map travel features are enabled.
type GameConfigMapTravel struct {
	MapCasesEnabled       bool `json:"mapCasesEnabled"`
	MovementEnabled       bool `json:"movementEnabled"`
	PositionMarkerEnabled bool `json:"positionMarkerEnabled"`
	StaminaBarEnabled     bool `json:"staminaBarEnabled"`
}

// GameConfigChatSlowmode configures chat posting limits for low-level players.
type GameConfigChatSlowmode struct {
	DuplicateWindowSeconds float64                      `json:"duplicateWindowSeconds"`
	LiftedAtLevel          int                          `json:"liftedAtLevel"`
	Tiers                  []GameConfigChatSlowmodeTier `json:"tiers"`
}

// GameConfigChatSlowmodeTier is the chat limit that applies from MinLevel.
type GameConfigChatSlowmodeTier struct {
	CooldownSeconds float64 `json:"cooldownSeconds"`
	MaxPerHour      float64 `json:"maxPerHour"`
	MinLevel        int     `json:"minLevel"`
	NewDirectPerDay float64 `json:"newDirectPerDay"`
}

// GameConfigSkillPresets sets how many skill presets a player can save.
type GameConfigSkillPresets struct {
	NonPremiumMax int `json:"nonPremiumMax"`
	PremiumMax    int `json:"premiumMax"`
}

// UserBuffs lists a player's active debuffs.
type UserBuffs struct {
	DebuffCodes []string `json:"debuffCodes"`
	DebuffEndAt Time     `json:"debuffEndAt"`
}

// UserEnergyGift is an energy gift granted to a player.
type UserEnergyGift struct {
	Amount    float64 `json:"amount"`
	ClaimedAt Time    `json:"claimedAt"`
	GrantedAt Time    `json:"grantedAt"`
}

// UserInfos holds a player's profile cosmetics and government roles.
type UserInfos struct {
	ColorScheme               string `json:"colorScheme"`
	CongressMemberOfCountryID string `json:"congressMemberOf"`
	Font                      string `json:"font"`
	IsPremium                 bool   `json:"isPremium"`
	MinisterOfCountryID       string `json:"ministerOf"`
	MinisterTitle             string `json:"ministerTitle"`
	PremiumMonthsCount        int    `json:"premiumMonthsCount"`
	UsernameEffect            string `json:"usernameEffect"`
}

// UserWarning is a moderation warning shown on a player's profile.
type UserWarning struct {
	IsAccepted bool   `json:"isAccepted"`
	Message    string `json:"message"`
}

// UserMissionsGeneratedAt records when each mission set was generated.
type UserMissionsGeneratedAt struct {
	Daily  Time `json:"daily"`
	Weekly Time `json:"weekly"`
}

// UserStatsToolbox counts the toolboxes a player opened.
type UserStatsToolbox struct {
	ByRarity    map[string]int `json:"byRarity"`
	OpenedCount int            `json:"openedCount"`
}

// GameConfigSkill holds the per-level configuration of one skill, keyed by
// skill name in [GameConfig.Skills] and by level in Levels.
type GameConfigSkill struct {
	Levels map[string]GameConfigSkillLevel `json:"levels"`
	// SkillOverflow names the skill that receives this skill's value above its
	// cap, at SkillOverflowValue per point.
	SkillOverflow      string  `json:"skillOverflow"`
	SkillOverflowValue float64 `json:"skillOverflowValue"`
	SoftCap            float64 `json:"softCap"`
}
