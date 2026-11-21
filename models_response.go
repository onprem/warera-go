package warera

import "encoding/json"

// Response models are maintained separately from the fixture recorder.
// Fixtures provide examples for checking decoding; they do not determine these Go types.

// ItemPrices represents market reference prices indexed by item code.
type ItemPrices map[string]float64

// Article is a news article published in the WarEra game.
type Article struct {
	ID          string       `json:"_id,omitempty"`
	Author      string       `json:"author,omitempty"`
	Category    string       `json:"category,omitempty"`
	Content     string       `json:"content,omitempty"`
	CreatedAt   string       `json:"createdAt,omitempty"`
	IsDeleted   bool         `json:"isDeleted,omitempty"`
	IsPublic    bool         `json:"isPublic,omitempty"`
	IsPublished bool         `json:"isPublished,omitempty"`
	Language    string       `json:"language,omitempty"`
	PublishedAt string       `json:"publishedAt,omitempty"`
	Slug        string       `json:"slug,omitempty"`
	Stats       ArticleStats `json:"stats,omitempty"`
	Title       string       `json:"title,omitempty"`
	UpdatedAt   string       `json:"updatedAt,omitempty"`
}

// ArticleStats is a nested component of Article.
type ArticleStats struct {
	Bookmarks int     `json:"bookmarks,omitempty"`
	Comments  int     `json:"comments,omitempty"`
	Dislikes  int     `json:"dislikes,omitempty"`
	GemTips   float64 `json:"gemTips,omitempty"`
	Likes     int     `json:"likes,omitempty"`
	Score     float64 `json:"score,omitempty"`
	Subs      int     `json:"subs,omitempty"`
	Tips      float64 `json:"tips,omitempty"`
	Views     int     `json:"views,omitempty"`
}

// Battle is a military conflict between attacker and defender sides.
type Battle struct {
	ID              string               `json:"_id,omitempty"`
	Attacker        BattleSideSummary    `json:"attacker,omitempty"`
	BadgesProcessed bool                 `json:"badgesProcessed,omitempty"`
	CreatedAt       string               `json:"createdAt,omitempty"`
	CurrentRound    *BattleRound         `json:"currentRound,omitempty"`
	Defender        BattleSideSummary    `json:"defender,omitempty"`
	EndedAt         string               `json:"endedAt,omitempty"`
	IsActive        bool                 `json:"isActive,omitempty"`
	Rounds          []string             `json:"rounds,omitempty"`
	RoundsHistory   []BattleRoundHistory `json:"roundsHistory,omitempty"`
	RoundsToWin     int                  `json:"roundsToWin,omitempty"`
	Stats           BattleStats          `json:"stats,omitempty"`
	Type            string               `json:"type,omitempty"`
	UpdatedAt       string               `json:"updatedAt,omitempty"`
	WarID           string               `json:"war,omitempty"`
	WonBy           string               `json:"wonBy,omitempty"`
}

// BattleRoundHistory summarizes the outcome of a completed battle round.
type BattleRoundHistory struct {
	AttackerDamages float64 `json:"attackerDamages,omitempty"`
	AttackerPoints  float64 `json:"attackerPoints,omitempty"`
	DefenderDamages float64 `json:"defenderDamages,omitempty"`
	DefenderPoints  float64 `json:"defenderPoints,omitempty"`
	WonBy           string  `json:"wonBy,omitempty"`
}

// BattleOrder is an order assigned to a country or military unit on a battle side.
type BattleOrder struct {
	ID          string     `json:"_id,omitempty"`
	Battle      string     `json:"battle,omitempty"`
	CountryID   string     `json:"country,omitempty"`
	CreatedAt   string     `json:"createdAt,omitempty"`
	IsActive    bool       `json:"isActive,omitempty"`
	MUID        string     `json:"mu,omitempty"`
	Priority    string     `json:"priority,omitempty"`
	Rank        int        `json:"rank,omitempty"`
	Side        BattleSide `json:"side,omitempty"`
	SideCountry string     `json:"sideCountry,omitempty"`
	Text        string     `json:"text,omitempty"`
	UpdatedAt   string     `json:"updatedAt,omitempty"`
	UserID      string     `json:"user,omitempty"`
}

// BattleRound is a nested component of Battle.
type BattleRound struct {
	ID        string               `json:"_id,omitempty"`
	Attacker  RoundSideCombatStats `json:"attacker,omitempty"`
	Battle    string               `json:"battle,omitempty"`
	CreatedAt string               `json:"createdAt,omitempty"`
	Defender  RoundSideCombatStats `json:"defender,omitempty"`
	IsActive  bool                 `json:"isActive,omitempty"`
	Live      RoundLiveState       `json:"live,omitempty"`
	Number    int                  `json:"number,omitempty"`
	UpdatedAt string               `json:"updatedAt,omitempty"`
}

// BattleSideSummary is a nested component of Battle.
type BattleSideSummary struct {
	BountyEffectiveAt string   `json:"bountyEffectiveAt,omitempty"`
	CountryID         string   `json:"country,omitempty"`
	CountryOrders     []string `json:"countryOrders,omitempty"`
	Damages           float64  `json:"damages,omitempty"`
	HitCount          int      `json:"hitCount,omitempty"`
	MoneyPer1kDamages float64  `json:"moneyPer1kDamages,omitempty"`
	MoneyPool         float64  `json:"moneyPool,omitempty"`
	MuOrders          []string `json:"muOrders,omitempty"`
	Region            string   `json:"region,omitempty"`
	WonRoundsCount    int      `json:"wonRoundsCount,omitempty"`
}

// BattleStats is a nested component of Battle.
type BattleStats struct {
	HitCount int `json:"hitCount,omitempty"`
}

// Company is a business that operates in a region and employs workers.
type Company struct {
	ID                  string                     `json:"_id,omitempty"`
	ActiveUpgradeLevels CompanyActiveUpgradeLevels `json:"activeUpgradeLevels,omitempty"`
	ConcreteInvested    float64                    `json:"concreteInvested,omitempty"`
	CreatedAt           string                     `json:"createdAt,omitempty"`
	EstimatedValue      float64                    `json:"estimatedValue,omitempty"`
	FreeMoves           int                        `json:"freeMoves,omitempty"`
	IsFull              bool                       `json:"isFull,omitempty"`
	ItemCode            string                     `json:"itemCode,omitempty"`
	Name                string                     `json:"name,omitempty"`
	Production          float64                    `json:"production,omitempty"`
	Region              string                     `json:"region,omitempty"`
	UpdatedAt           string                     `json:"updatedAt,omitempty"`
	UserID              string                     `json:"user,omitempty"`
}

// CompanyActiveUpgradeLevels is a nested component of Company.
type CompanyActiveUpgradeLevels struct {
	AutomatedEngine int `json:"automatedEngine,omitempty"`
	BreakRoom       int `json:"breakRoom,omitempty"`
	Storage         int `json:"storage,omitempty"`
}

// Country is a game nation with political, economic, and diplomatic state.
type Country struct {
	ID                 string                         `json:"_id,omitempty"`
	AllianceID         string                         `json:"allianceId,omitempty"`
	Allies             []string                       `json:"allies,omitempty"`
	AverageDevelopment float64                        `json:"averageDevelopment,omitempty"`
	BordersOpenUntil   string                         `json:"bordersOpenUntil,omitempty"`
	Code               string                         `json:"code,omitempty"`
	CoreDevelopment    float64                        `json:"coreDevelopment,omitempty"`
	CreatedAt          string                         `json:"createdAt,omitempty"`
	CurrentBattleOrder string                         `json:"currentBattleOrder,omitempty"`
	CurrentDevelopment float64                        `json:"currentDevelopment,omitempty"`
	CurrentPopulation  int                            `json:"currentPopulation,omitempty"`
	Dates              CountryDates                   `json:"dates,omitempty"`
	DefensivePacts     []string                       `json:"defensivePacts,omitempty"`
	DiscordURL         string                         `json:"discordUrl,omitempty"`
	Enemy              string                         `json:"enemy,omitempty"`
	GovType            string                         `json:"govType,omitempty"`
	Income             *IncomeBreakdown               `json:"income,omitempty"`
	MapAccent          string                         `json:"mapAccent,omitempty"`
	Money              float64                        `json:"money,omitempty"`
	Name               string                         `json:"name,omitempty"`
	NonAggressionUntil map[string]string              `json:"nonAggressionUntil,omitempty"`
	Orgs               []string                       `json:"orgs,omitempty"`
	PinnedArticle      string                         `json:"pinnedArticle,omitempty"`
	Rankings           map[string]RankingItem         `json:"rankings,omitempty"`
	RulingParty        string                         `json:"rulingParty,omitempty"`
	Scheme             string                         `json:"scheme,omitempty"`
	StrategicResources *CountryStrategicResourcesInfo `json:"strategicResources,omitempty"`
	Taxes              CountryTaxes                   `json:"taxes,omitempty"`
	Unrest             CountryUnrest                  `json:"unrest,omitempty"`
	UpdatedAt          string                         `json:"updatedAt,omitempty"`
	WarsWith           []string                       `json:"warsWith,omitempty"`
}

// CountryDates is a nested component of Country.
type CountryDates struct {
	GovTypeChangedAt   string `json:"govTypeChangedAt,omitempty"`
	LastAllianceLeftAt string `json:"lastAllianceLeftAt,omitempty"`
}

// CountryStrategicResourcesInfo is a nested component of Country.
type CountryStrategicResourcesInfo struct {
	Resources CountryStrategicResourcesInfoResources `json:"resources,omitempty"`
}

// CountryStrategicResourcesInfoResources is a nested component of CountryStrategicResourcesInfo.
type CountryStrategicResourcesInfoResources struct {
	Coal       []string `json:"coal,omitempty"`
	Diamonds   []string `json:"diamonds,omitempty"`
	Gold       []string `json:"gold,omitempty"`
	Lithium    []string `json:"lithium,omitempty"`
	RareEarths []string `json:"rareEarths,omitempty"`
	Uranium    []string `json:"uranium,omitempty"`
}

// CountryTaxes is a nested component of Country.
type CountryTaxes struct {
	Income  float64 `json:"income,omitempty"`
	Manager float64 `json:"manager,omitempty"`
	Market  float64 `json:"market,omitempty"`
}

// CountryUnrest is a nested component of Country.
type CountryUnrest struct {
	Bar                   float64 `json:"bar,omitempty"`
	BarMax                float64 `json:"barMax,omitempty"`
	LastContributionAt    string  `json:"lastContributionAt,omitempty"`
	LastRevolutionEndedAt string  `json:"lastRevolutionEndedAt,omitempty"`
}

// CurrentEquipment contains the equipment currently worn by a user.
type CurrentEquipment struct {
	Ammo   string        `json:"ammo,omitempty"`
	Boots  EquipmentItem `json:"boots,omitempty"`
	Chest  EquipmentItem `json:"chest,omitempty"`
	Gloves EquipmentItem `json:"gloves,omitempty"`
	Helmet EquipmentItem `json:"helmet,omitempty"`
	Pants  EquipmentItem `json:"pants,omitempty"`
	Weapon *WeaponItem   `json:"weapon,omitempty"`
}

// EquipmentItem is a nested component of CurrentEquipment.
type EquipmentItem struct {
	ID                string          `json:"_id,omitempty"`
	Code              string          `json:"code,omitempty"`
	LastAcquisitionAt string          `json:"lastAcquisitionAt,omitempty"`
	MaxState          float64         `json:"maxState,omitempty"`
	Quantity          int             `json:"quantity,omitempty"`
	Skills            EquipmentSkills `json:"skills,omitempty"`
	State             float64         `json:"state,omitempty"`
	Type              string          `json:"type,omitempty"`
}

// EquipmentSkills is a nested component of RoundEquipment.
type EquipmentSkills struct {
	Armor           float64 `json:"armor,omitempty"`
	CriticalDamages float64 `json:"criticalDamages,omitempty"`
	Dodge           float64 `json:"dodge,omitempty"`
	Precision       float64 `json:"precision,omitempty"`
}

// Event is a world or country event with event-specific data.
type Event struct {
	ID        string    `json:"_id,omitempty"`
	Countries []string  `json:"countries,omitempty"`
	CreatedAt string    `json:"createdAt,omitempty"`
	Data      EventData `json:"data,omitempty"`
	Priority  string    `json:"priority,omitempty"`
	UpdatedAt string    `json:"updatedAt,omitempty"`
}

// EventData is a nested component of Event.
type EventData struct {
	AllianceID      string    `json:"allianceId,omitempty"`
	AllianceName    string    `json:"allianceName,omitempty"`
	AttackerCountry string    `json:"attackerCountry,omitempty"`
	AttackerRegion  string    `json:"attackerRegion,omitempty"`
	Battle          string    `json:"battle,omitempty"`
	BonusPercent    float64   `json:"bonusPercent,omitempty"`
	Countries       []string  `json:"countries,omitempty"`
	CountryID       string    `json:"country,omitempty"`
	DefenderCountry string    `json:"defenderCountry,omitempty"`
	DefenderRegion  string    `json:"defenderRegion,omitempty"`
	DurationDays    float64   `json:"durationDays,omitempty"`
	ItemCode        string    `json:"itemCode,omitempty"`
	Money           float64   `json:"money,omitempty"`
	Region          string    `json:"region,omitempty"`
	Type            EventType `json:"type,omitempty"`
	Wars            []string  `json:"wars,omitempty"`
	WonBy           string    `json:"wonBy,omitempty"`
}

// GameConfig contains server settings for game systems, skills, items, and rewards.
type GameConfig struct {
	Alliance               GameConfigAlliance               `json:"alliance,omitempty"`
	Announcement           GameConfigAnnouncement           `json:"announcement,omitempty"`
	Badge                  map[string]GameConfigBadge       `json:"badge,omitempty"`
	Battle                 GameConfigBattle                 `json:"battle,omitempty"`
	CitizenshipApplication GameConfigCitizenshipApplication `json:"citizenshipApplication,omitempty"`
	Company                GameConfigCompany                `json:"company,omitempty"`
	Country                GameConfigCountry                `json:"country,omitempty"`
	Election               GameConfigElection               `json:"election,omitempty"`
	Government             GameConfigGovernment             `json:"government,omitempty"`
	Items                  map[string]GameConfigItem        `json:"items,omitempty"`
	Law                    GameConfigLaw                    `json:"law,omitempty"`
	Loot                   GameConfigLoot                   `json:"loot,omitempty"`
	MapTravel              GameConfigMapTravel              `json:"mapTravel,omitempty"`
	MercenaryContract      GameConfigMercenaryContract      `json:"mercenaryContract,omitempty"`
	MergingCost            GameConfigMergingCost            `json:"mergingCost,omitempty"`
	Mission                GameConfigMission                `json:"mission,omitempty"`
	MU                     GameConfigMu                     `json:"mu,omitempty"`
	Newspaper              GameConfigNewspaper              `json:"newspaper,omitempty"`
	Org                    GameConfigOrg                    `json:"org,omitempty"`
	Party                  GameConfigParty                  `json:"party,omitempty"`
	Referral               GameConfigReferral               `json:"referral,omitempty"`
	Region                 GameConfigRegion                 `json:"region,omitempty"`
	Skills                 map[string]GameConfigSkill       `json:"skills,omitempty"`
	SubSkinReward          GameConfigSubSkinReward          `json:"subSkinReward,omitempty"`
	Unrest                 GameConfigUnrest                 `json:"unrest,omitempty"`
	Upgrade                GameConfigUpgrade                `json:"upgrade,omitempty"`
	UpgradesConfig         GameConfigUpgradesConfig         `json:"upgradesConfig,omitempty"`
	User                   GameConfigUser                   `json:"user,omitempty"`
}

// GameConfigAlliance is a nested component of GameConfig.
type GameConfigAlliance struct {
	LeaveCooldownDays          float64 `json:"leaveCooldownDays,omitempty"`
	LeaveCooldownDisabledUntil string  `json:"leaveCooldownDisabledUntil,omitempty"`
}

// GameConfigAnnouncement is a nested component of GameConfig.
type GameConfigAnnouncement struct {
	DurationHours float64 `json:"durationHours,omitempty"`
	MaxLength     int     `json:"maxLength,omitempty"`
	MaxPerHour    float64 `json:"maxPerHour,omitempty"`
}

// GameConfigBattle is a nested component of GameConfig.
type GameConfigBattle struct {
	AllianceDamagesBonusPercent          float64        `json:"allianceDamagesBonusPercent,omitempty"`
	CasesPer1kDamagesInPool              float64        `json:"casesPer1kDamagesInPool,omitempty"`
	CountryOrderBonusPercent             float64        `json:"countryOrderBonusPercent,omitempty"`
	EnemyDamagesBonusPercent             float64        `json:"enemyDamagesBonusPercent,omitempty"`
	GovMemberBountyMinPpPer1k            float64        `json:"govMemberBountyMinPpPer1k,omitempty"`
	GovMemberBountyRewardPercent         float64        `json:"govMemberBountyRewardPercent,omitempty"`
	HealthCost                           float64        `json:"healthCost,omitempty"`
	HitFor1CaseInPool                    float64        `json:"hitFor1CaseInPool,omitempty"`
	LocationGateEnabled                  bool           `json:"locationGateEnabled,omitempty"`
	LostAttackingRegionMalusPercent      float64        `json:"lostAttackingRegionMalusPercent,omitempty"`
	MaxRounds                            int            `json:"maxRounds,omitempty"`
	MuOrderBonusPercent                  float64        `json:"muOrderBonusPercent,omitempty"`
	OccupyingYourRegionsMalusPercent     float64        `json:"occupyingYourRegionsMalusPercent,omitempty"`
	PatrioticBonusPercent                float64        `json:"patrioticBonusPercent,omitempty"`
	PointsToWinRound                     float64        `json:"pointsToWinRound,omitempty"`
	RankingsLootPercentPer1kDmg          float64        `json:"rankingsLootPercentPer1kDmg,omitempty"`
	RegionNotLinkedToCapitalMalusPercent float64        `json:"regionNotLinkedToCapitalMalusPercent,omitempty"`
	RoundsToWin                          int            `json:"roundsToWin,omitempty"`
	SetCountryOrderMoneyCost             float64        `json:"setCountryOrderMoneyCost,omitempty"`
	SetMuOrderMoneyCost                  float64        `json:"setMuOrderMoneyCost,omitempty"`
	SetOrderMoneyCost                    float64        `json:"setOrderMoneyCost,omitempty"`
	TickPoints                           map[string]int `json:"tickPoints,omitempty"`
}

// GameConfigCitizenshipApplication is a nested component of GameConfig.
type GameConfigCitizenshipApplication struct {
	AutoApprovalEnabled       bool    `json:"autoApprovalEnabled,omitempty"`
	AutoApprovalMaxPopulation float64 `json:"autoApprovalMaxPopulation,omitempty"`
	NewCountryBordersOpenDays float64 `json:"newCountryBordersOpenDays,omitempty"`
}

// GameConfigCompany is a nested component of GameConfig.
type GameConfigCompany struct {
	ChangeItemCost                     float64 `json:"changeItemCost,omitempty"`
	ConstructionCostIncreasePerCompany float64 `json:"constructionCostIncreasePerCompany,omitempty"`
	DepositResourceBonus               float64 `json:"depositResourceBonus,omitempty"`
	DestructionValuePercent            float64 `json:"destructionValuePercent,omitempty"`
	MoveCost                           float64 `json:"moveCost,omitempty"`
}

// GameConfigCountry is a nested component of GameConfig.
type GameConfigCountry struct {
	ConqueredRegionCoreTaxPercent float64 `json:"conqueredRegionCoreTaxPercent,omitempty"`
	MaxTaxAmount                  float64 `json:"maxTaxAmount,omitempty"`
}

// GameConfigElection is a nested component of GameConfig.
type GameConfigElection struct {
	CandidateDurationHours                  float64 `json:"candidateDurationHours,omitempty"`
	CandidateMinLevel                       int     `json:"candidateMinLevel,omitempty"`
	CongressMaxSeats                        int     `json:"congressMaxSeats,omitempty"`
	ElectionVoteDurationHours               float64 `json:"electionVoteDurationHours,omitempty"`
	VoteCitizenshipLockNewCountryMultiplier float64 `json:"voteCitizenshipLockNewCountryMultiplier,omitempty"`
	VoteMinDaysSinceCitizenshipChange       float64 `json:"voteMinDaysSinceCitizenshipChange,omitempty"`
	VoteMinLevel                            int     `json:"voteMinLevel,omitempty"`
}

// GameConfigGovernment is a nested component of GameConfig.
type GameConfigGovernment struct {
	NominationCooldownDays   float64 `json:"nominationCooldownDays,omitempty"`
	PresidentWagePercent     float64 `json:"presidentWagePercent,omitempty"`
	VicePresidentWagePercent float64 `json:"vicePresidentWagePercent,omitempty"`
}

// GameConfigLaw is a nested component of GameConfig.
type GameConfigLaw struct {
	AbusiveLawPossibleVotersNeeded float64                           `json:"abusiveLawPossibleVotersNeeded,omitempty"`
	AbusiveLawsCooldownInDays      float64                           `json:"abusiveLawsCooldownInDays,omitempty"`
	Accept_join_alliance           GameConfigLawAccept_join_alliance `json:"accept_join_alliance,omitempty"`
	Create_alliance                GameConfigLawCreate_alliance      `json:"create_alliance,omitempty"`
	Define_enemy_country           GameConfigLawDefine_enemy_country `json:"define_enemy_country,omitempty"`
	LawVotesDurationHours          float64                           `json:"lawVotesDurationHours,omitempty"`
	Leave_alliance                 GameConfigLawLeave_alliance       `json:"leave_alliance,omitempty"`
	SendMoneyToCountry             GameConfigLawSendMoneyToCountry   `json:"sendMoneyToCountry,omitempty"`
	Set_color_scheme               GameConfigLawSet_color_scheme     `json:"set_color_scheme,omitempty"`
	VotersRatioNeeded              float64                           `json:"votersRatioNeeded,omitempty"`
}

// GameConfigLawAccept_join_alliance is a nested component of GameConfigLaw.
type GameConfigLawAccept_join_alliance struct {
	Cost float64 `json:"cost,omitempty"`
}

// GameConfigLawCreate_alliance is a nested component of GameConfigLaw.
type GameConfigLawCreate_alliance struct {
	Cost            float64 `json:"cost,omitempty"`
	MaintenanceCost float64 `json:"maintenanceCost,omitempty"`
}

// GameConfigLawDefine_enemy_country is a nested component of GameConfigLaw.
type GameConfigLawDefine_enemy_country struct {
	Cost            float64 `json:"cost,omitempty"`
	MaintenanceCost float64 `json:"maintenanceCost,omitempty"`
}

// GameConfigLawLeave_alliance is a nested component of GameConfigLaw.
type GameConfigLawLeave_alliance struct {
	Cost float64 `json:"cost,omitempty"`
}

// GameConfigLawSendMoneyToCountry is a nested component of GameConfigLaw.
type GameConfigLawSendMoneyToCountry struct {
	AllianceTaxRate float64 `json:"allianceTaxRate,omitempty"`
	ExternalTaxRate float64 `json:"externalTaxRate,omitempty"`
}

// GameConfigLawSet_color_scheme is a nested component of GameConfigLaw.
type GameConfigLawSet_color_scheme struct {
	Cost float64 `json:"cost,omitempty"`
}

// GameConfigLoot is a nested component of GameConfig.
type GameConfigLoot struct {
	BattleLootDamagePerLootItem float64 `json:"battleLootDamagePerLootItem,omitempty"`
	DamagePerLootItem           float64 `json:"damagePerLootItem,omitempty"`
	WeaponChancePercent         float64 `json:"weaponChancePercent,omitempty"`
}

// GameConfigMercenaryContract is a nested component of GameConfig.
type GameConfigMercenaryContract struct {
	AcceptanceFeePercent            float64                               `json:"acceptanceFeePercent,omitempty"`
	Auction                         GameConfigMercenaryContractAuction    `json:"auction,omitempty"`
	BountyCooldownMinutes           float64                               `json:"bountyCooldownMinutes,omitempty"`
	BountyMinActiveCitizens         float64                               `json:"bountyMinActiveCitizens,omitempty"`
	CancellationPenaltyPercent      float64                               `json:"cancellationPenaltyPercent,omitempty"`
	Enabled                         bool                                  `json:"enabled,omitempty"`
	GovMemberBidMaxPpPer1k          float64                               `json:"govMemberBidMaxPpPer1k,omitempty"`
	MaxBountyPpMultiplier           float64                               `json:"maxBountyPpMultiplier,omitempty"`
	NationalBountyEnabled           bool                                  `json:"nationalBountyEnabled,omitempty"`
	NationalBountyMinActiveCitizens float64                               `json:"nationalBountyMinActiveCitizens,omitempty"`
	Reputation                      GameConfigMercenaryContractReputation `json:"reputation,omitempty"`
}

// GameConfigMercenaryContractAuction is a nested component of GameConfigMercenaryContract.
type GameConfigMercenaryContractAuction struct {
	BidStep                 float64 `json:"bidStep,omitempty"`
	MaxDamagePerContract    float64 `json:"maxDamagePerContract,omitempty"`
	MaxDuration             float64 `json:"maxDuration,omitempty"`
	MaxPerKMultiplier       float64 `json:"maxPerKMultiplier,omitempty"`
	MinDuration             float64 `json:"minDuration,omitempty"`
	MinPerK                 float64 `json:"minPerK,omitempty"`
	TimerExtensionAmount    float64 `json:"timerExtensionAmount,omitempty"`
	TimerExtensionThreshold float64 `json:"timerExtensionThreshold,omitempty"`
}

// GameConfigMercenaryContractReputation is a nested component of GameConfigMercenaryContract.
type GameConfigMercenaryContractReputation struct {
	FailurePerDamage           float64 `json:"failurePerDamage,omitempty"`
	FailurePerDollar           float64 `json:"failurePerDollar,omitempty"`
	NegativeBuyAmount          float64 `json:"negativeBuyAmount,omitempty"`
	NegativeBuyCooldownHours   float64 `json:"negativeBuyCooldownHours,omitempty"`
	NegativeBuyCost            float64 `json:"negativeBuyCost,omitempty"`
	ProfessionalsOnlyThreshold float64 `json:"professionalsOnlyThreshold,omitempty"`
	SuccessPerDollar           float64 `json:"successPerDollar,omitempty"`
	WeeklyDecayPercent         float64 `json:"weeklyDecayPercent,omitempty"`
}

// GameConfigMergingCost is a nested component of GameConfig.
type GameConfigMergingCost struct {
	Common    float64 `json:"common,omitempty"`
	Epic      float64 `json:"epic,omitempty"`
	Legendary float64 `json:"legendary,omitempty"`
	Mythic    float64 `json:"mythic,omitempty"`
	Rare      float64 `json:"rare,omitempty"`
	Uncommon  float64 `json:"uncommon,omitempty"`
}

// GameConfigMission is a nested component of GameConfig.
type GameConfigMission struct {
	RerollMissionCost map[string]int          `json:"rerollMissionCost,omitempty"`
	Reward            GameConfigMissionReward `json:"reward,omitempty"`
}

// GameConfigMissionReward is a nested component of GameConfigMission.
type GameConfigMissionReward struct {
	Daily    GameConfigMissionRewardDaily    `json:"daily,omitempty"`
	Monthly  GameConfigMissionRewardMonthly  `json:"monthly,omitempty"`
	Starting GameConfigMissionRewardStarting `json:"starting,omitempty"`
	Weekly   GameConfigMissionRewardWeekly   `json:"weekly,omitempty"`
}

// GameConfigMissionRewardDaily is a nested component of GameConfigMissionReward.
type GameConfigMissionRewardDaily struct {
	Cases          float64 `json:"cases,omitempty"`
	Money          float64 `json:"money,omitempty"`
	XP             float64 `json:"xp,omitempty"`
	XpWhenFinished float64 `json:"xpWhenFinished,omitempty"`
}

// GameConfigMissionRewardMonthly is a nested component of GameConfigMissionReward.
type GameConfigMissionRewardMonthly struct {
	Cases          float64 `json:"cases,omitempty"`
	Money          float64 `json:"money,omitempty"`
	XP             float64 `json:"xp,omitempty"`
	XpWhenFinished float64 `json:"xpWhenFinished,omitempty"`
}

// GameConfigMissionRewardStarting is a nested component of GameConfigMissionReward.
type GameConfigMissionRewardStarting struct {
	Cases          float64 `json:"cases,omitempty"`
	Money          float64 `json:"money,omitempty"`
	XP             float64 `json:"xp,omitempty"`
	XpWhenFinished float64 `json:"xpWhenFinished,omitempty"`
}

// GameConfigMissionRewardWeekly is a nested component of GameConfigMissionReward.
type GameConfigMissionRewardWeekly struct {
	Cases          float64 `json:"cases,omitempty"`
	Money          float64 `json:"money,omitempty"`
	XP             float64 `json:"xp,omitempty"`
	XpWhenFinished float64 `json:"xpWhenFinished,omitempty"`
}

// GameConfigMu is a nested component of GameConfig.
type GameConfigMu struct {
	ConstructionCost              float64 `json:"constructionCost,omitempty"`
	DestructionValuePercent       float64 `json:"destructionValuePercent,omitempty"`
	HealthPerHelp                 float64 `json:"healthPerHelp,omitempty"`
	HelpCooldownHours             float64 `json:"helpCooldownHours,omitempty"`
	HelperHealthPerHelp           float64 `json:"helperHealthPerHelp,omitempty"`
	HelpValue                     float64 `json:"helpValue,omitempty"`
	MaxOwnedMus                   int     `json:"maxOwnedMus,omitempty"`
	MoveCost                      float64 `json:"moveCost,omitempty"`
	NationalityChangeConcreteCost float64 `json:"nationalityChangeConcreteCost,omitempty"`
	NationalityChangeCooldownDays float64 `json:"nationalityChangeCooldownDays,omitempty"`
}

// GameConfigNewspaper is a nested component of GameConfig.
type GameConfigNewspaper struct {
	CreateArticleMinLevel int     `json:"createArticleMinLevel,omitempty"`
	GemTipValue           float64 `json:"gemTipValue,omitempty"`
	PublishCost           float64 `json:"publishCost,omitempty"`
	TipMinLevel           int     `json:"tipMinLevel,omitempty"`
	TipValue              float64 `json:"tipValue,omitempty"`
}

// GameConfigOrg is a nested component of GameConfig.
type GameConfigOrg struct {
	ConstructionCost float64 `json:"constructionCost,omitempty"`
	MoveCost         float64 `json:"moveCost,omitempty"`
}

// GameConfigParty is a nested component of GameConfig.
type GameConfigParty struct {
	CreateCost float64 `json:"createCost,omitempty"`
}

// GameConfigReferral is a nested component of GameConfig.
type GameConfigReferral struct {
	CanSetReferrerBeforeOrAtLevel  int     `json:"canSetReferrerBeforeOrAtLevel,omitempty"`
	LevelNeededForBadge            float64 `json:"levelNeededForBadge,omitempty"`
	LifeTimeBadgeMoneySharePercent float64 `json:"lifeTimeBadgeMoneySharePercent,omitempty"`
	MoneyForBeingReferred          float64 `json:"moneyForBeingReferred,omitempty"`
}

// GameConfigRegion is a nested component of GameConfig.
type GameConfigRegion struct {
	BattleCooldownHours                            float64 `json:"battleCooldownHours,omitempty"`
	DepleteHourlyPercent                           float64 `json:"depleteHourlyPercent,omitempty"`
	DepositsEnabled                                bool    `json:"depositsEnabled,omitempty"`
	DevelopmentIncomeMultiplier                    float64 `json:"developmentIncomeMultiplier,omitempty"`
	DevelopmentMultiplier                          float64 `json:"developmentMultiplier,omitempty"`
	HomeMoveConcreteCost                           float64 `json:"homeMoveConcreteCost,omitempty"`
	LiberationDaysCooldown                         float64 `json:"liberationDaysCooldown,omitempty"`
	LostResistanceOtherRegionsBarMultiplier        float64 `json:"lostResistanceOtherRegionsBarMultiplier,omitempty"`
	MaxDailyResistance                             float64 `json:"maxDailyResistance,omitempty"`
	MaxResistance                                  float64 `json:"maxResistance,omitempty"`
	MinDailyResistance                             float64 `json:"minDailyResistance,omitempty"`
	NonAggressionHoursAfterLiberation              float64 `json:"nonAggressionHoursAfterLiberation,omitempty"`
	NonAggressionHoursAfterPeace                   float64 `json:"nonAggressionHoursAfterPeace,omitempty"`
	ResistanceAllyBonusPercent                     float64 `json:"resistanceAllyBonusPercent,omitempty"`
	ResistanceBarMultiplier                        float64 `json:"resistanceBarMultiplier,omitempty"`
	ResistanceBattleCooldownHours                  float64 `json:"resistanceBattleCooldownHours,omitempty"`
	ResistanceBattleStartCostMultiplier            float64 `json:"resistanceBattleStartCostMultiplier,omitempty"`
	ResistanceCitizenBonus                         float64 `json:"resistanceCitizenBonus,omitempty"`
	ResistanceCitizenBonusPercent                  float64 `json:"resistanceCitizenBonusPercent,omitempty"`
	ResistanceContributionCooldownAfterRevoltHours float64 `json:"resistanceContributionCooldownAfterRevoltHours,omitempty"`
	ResistanceContributionCost                     float64 `json:"resistanceContributionCost,omitempty"`
	ResistanceContributionMinLevel                 int     `json:"resistanceContributionMinLevel,omitempty"`
	ResistanceContributionValue                    float64 `json:"resistanceContributionValue,omitempty"`
	ResistanceDecayPercent                         float64 `json:"resistanceDecayPercent,omitempty"`
	ResistanceForeignGovCostMultiplier             float64 `json:"resistanceForeignGovCostMultiplier,omitempty"`
	ResistanceGrowthPercentMax                     float64 `json:"resistanceGrowthPercentMax,omitempty"`
	ResistanceGrowthPercentMin                     float64 `json:"resistanceGrowthPercentMin,omitempty"`
	ResistancePassiveGrowthPercent                 float64 `json:"resistancePassiveGrowthPercent,omitempty"`
	StrategicResourceIncomeMultiplier              float64 `json:"strategicResourceIncomeMultiplier,omitempty"`
	TransferDaysCooldown                           float64 `json:"transferDaysCooldown,omitempty"`
	WonResistanceBarMultiplier                     float64 `json:"wonResistanceBarMultiplier,omitempty"`
}

// GameConfigSkillLevel describes a skill's cost and value at a given level.
type GameConfigSkillLevel struct {
	Cost          float64 `json:"cost,omitempty"`
	IsABar        bool    `json:"isABar,omitempty"`
	TotalCost     float64 `json:"totalCost,omitempty"`
	UnlockAtLevel int     `json:"unlockAtLevel,omitempty"`
	Value         float64 `json:"value,omitempty"`
}

// GameConfigSubSkinReward is a nested component of GameConfig.
type GameConfigSubSkinReward struct {
	Deadline string `json:"deadline,omitempty"`
	SkinKey  string `json:"skinKey,omitempty"`
}

// GameConfigUnrest is a nested component of GameConfig.
type GameConfigUnrest struct {
	BarMultiplier                            float64 `json:"barMultiplier,omitempty"`
	BattleCooldownHours                      float64 `json:"battleCooldownHours,omitempty"`
	BattleStartCost                          float64 `json:"battleStartCost,omitempty"`
	BordersOpenDays                          float64 `json:"bordersOpenDays,omitempty"`
	ContributionCooldownAfterRevolutionHours float64 `json:"contributionCooldownAfterRevolutionHours,omitempty"`
	ContributionCost                         float64 `json:"contributionCost,omitempty"`
	ContributionMinLevel                     int     `json:"contributionMinLevel,omitempty"`
	ContributionValue                        float64 `json:"contributionValue,omitempty"`
	NominationPeriodHours                    float64 `json:"nominationPeriodHours,omitempty"`
}

// GameConfigUpgrade is a nested component of GameConfig.
type GameConfigUpgrade struct {
	ConstructionContributionMinLevel  int     `json:"constructionContributionMinLevel,omitempty"`
	ConstructionPointsPerContribution float64 `json:"constructionPointsPerContribution,omitempty"`
	RefundPercent                     float64 `json:"refundPercent,omitempty"`
	RegionDowngradeCooldownHours      float64 `json:"regionDowngradeCooldownHours,omitempty"`
	RegionUpgradeCooldownHours        float64 `json:"regionUpgradeCooldownHours,omitempty"`
}

// GameConfigUpgradesConfig is a nested component of GameConfig.
type GameConfigUpgradesConfig struct {
	AutomatedEngine    GameConfigUpgradesConfigAutomatedEngine    `json:"automatedEngine,omitempty"`
	Base               GameConfigUpgradesConfigBase               `json:"base,omitempty"`
	BreakRoom          GameConfigUpgradesConfigBreakRoom          `json:"breakRoom,omitempty"`
	Bunker             GameConfigUpgradesConfigBunker             `json:"bunker,omitempty"`
	Dormitories        GameConfigUpgradesConfigDormitories        `json:"dormitories,omitempty"`
	Headquarters       GameConfigUpgradesConfigHeadquarters       `json:"headquarters,omitempty"`
	PacificationCenter GameConfigUpgradesConfigPacificationCenter `json:"pacificationCenter,omitempty"`
	Storage            GameConfigUpgradesConfigStorage            `json:"storage,omitempty"`
}

// GameConfigUpgradeLevel describes the costs and effects of an upgrade level.
type GameConfigUpgradeLevel struct {
	ConstructionPointsCost         float64            `json:"constructionPointsCost,omitempty"`
	Level                          int                `json:"level,omitempty"`
	MaintenanceCost                float64            `json:"maintenanceCost,omitempty"`
	MaintenanceCostCountryDevScale float64            `json:"maintenanceCostCountryDevScale,omitempty"`
	MaintenanceCostRegionDevScale  float64            `json:"maintenanceCostRegionDevScale,omitempty"`
	MinimumMaintenanceCost         float64            `json:"minimumMaintenanceCost,omitempty"`
	Stats                          map[string]float64 `json:"stats,omitempty"`
	SteelCost                      float64            `json:"steelCost,omitempty"`
}

// GameConfigUpgradesConfigAutomatedEngine is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigAutomatedEngine struct {
	CanDowngrade         bool                              `json:"canDowngrade,omitempty"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels,omitempty"`
	PendingDurationHours float64                           `json:"pendingDurationHours,omitempty"`
}

// GameConfigUpgradesConfigBase is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigBase struct {
	CanBeDestroyed       bool                              `json:"canBeDestroyed,omitempty"`
	CanBeDisabled        bool                              `json:"canBeDisabled,omitempty"`
	CanDowngrade         bool                              `json:"canDowngrade,omitempty"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels,omitempty"`
	PendingDurationHours float64                           `json:"pendingDurationHours,omitempty"`
}

// GameConfigUpgradesConfigBreakRoom is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigBreakRoom struct {
	CanDowngrade bool                              `json:"canDowngrade,omitempty"`
	Levels       map[string]GameConfigUpgradeLevel `json:"levels,omitempty"`
}

// GameConfigUpgradesConfigBunker is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigBunker struct {
	CanBeDestroyed       bool                              `json:"canBeDestroyed,omitempty"`
	CanBeDisabled        bool                              `json:"canBeDisabled,omitempty"`
	CanDowngrade         bool                              `json:"canDowngrade,omitempty"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels,omitempty"`
	PendingDurationHours float64                           `json:"pendingDurationHours,omitempty"`
}

// GameConfigUpgradesConfigDormitories is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigDormitories struct {
	CanDowngrade bool                              `json:"canDowngrade,omitempty"`
	Levels       map[string]GameConfigUpgradeLevel `json:"levels,omitempty"`
}

// GameConfigUpgradesConfigHeadquarters is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigHeadquarters struct {
	CanBeDisabled        bool                              `json:"canBeDisabled,omitempty"`
	CanDowngrade         bool                              `json:"canDowngrade,omitempty"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels,omitempty"`
	PendingDurationHours float64                           `json:"pendingDurationHours,omitempty"`
}

// GameConfigUpgradesConfigPacificationCenter is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigPacificationCenter struct {
	CanBeDestroyed       bool                              `json:"canBeDestroyed,omitempty"`
	CanBeDisabled        bool                              `json:"canBeDisabled,omitempty"`
	CanDowngrade         bool                              `json:"canDowngrade,omitempty"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels,omitempty"`
	PendingDurationHours float64                           `json:"pendingDurationHours,omitempty"`
}

// GameConfigUpgradesConfigStorage is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigStorage struct {
	CanDowngrade         bool                              `json:"canDowngrade,omitempty"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels,omitempty"`
	PendingDurationHours float64                           `json:"pendingDurationHours,omitempty"`
}

// GameConfigUser is a nested component of GameConfig.
type GameConfigUser struct {
	ActiveCitizenMinLevel     int                         `json:"activeCitizenMinLevel,omitempty"`
	CanTakeControlAtLevel     int                         `json:"canTakeControlAtLevel,omitempty"`
	ChatSlowmode              GameConfigChatSlowmode      `json:"chatSlowmode,omitempty"`
	CitizenshipDaysCooldown   float64                     `json:"citizenshipDaysCooldown,omitempty"`
	ConstructionEnergyCost    float64                     `json:"constructionEnergyCost,omitempty"`
	DailyReward               GameConfigUserDailyReward   `json:"dailyReward,omitempty"`
	DailyXp                   float64                     `json:"dailyXp,omitempty"`
	DonationMinLevel          int                         `json:"donationMinLevel,omitempty"`
	EnergyCostPerAction       float64                     `json:"energyCostPerAction,omitempty"`
	EquipmentSets             GameConfigUserEquipmentSets `json:"equipmentSets,omitempty"`
	FeatureUnlockLevels       map[string]int              `json:"featureUnlockLevels,omitempty"`
	FieldsToPopulate          string                      `json:"fieldsToPopulate,omitempty"`
	IsInactiveAfterDays       float64                     `json:"isInactiveAfterDays,omitempty"`
	Kits                      GameConfigUserKits          `json:"kits,omitempty"`
	MarketMinLevel            int                         `json:"marketMinLevel,omitempty"`
	MaxConstructionPoints     float64                     `json:"maxConstructionPoints,omitempty"`
	MaxEnergy                 float64                     `json:"maxEnergy,omitempty"`
	MaxHunger                 float64                     `json:"maxHunger,omitempty"`
	Prestige                  GameConfigUserPrestige      `json:"prestige,omitempty"`
	RegenDividedBy            float64                     `json:"regenDividedBy,omitempty"`
	ResetSkillDaysCooldown    float64                     `json:"resetSkillDaysCooldown,omitempty"`
	ResetSkillsCostPerPoint   float64                     `json:"resetSkillsCostPerPoint,omitempty"`
	SkillPresets              GameConfigSkillPresets      `json:"skillPresets,omitempty"`
	TakeControlCooldownInDays float64                     `json:"takeControlCooldownInDays,omitempty"`
	XpPerAction               float64                     `json:"xpPerAction,omitempty"`
}

// GameConfigUserDailyReward is a nested component of GameConfigUser.
type GameConfigUserDailyReward struct {
	Case1 float64 `json:"case1,omitempty"`
	Money float64 `json:"money,omitempty"`
	XP    float64 `json:"xp,omitempty"`
}

// GameConfigUserEquipmentSets is a nested component of GameConfigUser.
type GameConfigUserEquipmentSets struct {
	NonPremiumMax float64 `json:"nonPremiumMax,omitempty"`
	PremiumMax    float64 `json:"premiumMax,omitempty"`
}

// GameConfigUserKits is a nested component of GameConfigUser.
type GameConfigUserKits struct {
	NonPremiumMax float64 `json:"nonPremiumMax,omitempty"`
	PremiumMax    float64 `json:"premiumMax,omitempty"`
}

// GameConfigUserPrestige is a nested component of GameConfigUser.
type GameConfigUserPrestige struct {
	Enabled                     bool    `json:"enabled,omitempty"`
	EnabledLevel                int     `json:"enabledLevel,omitempty"`
	SkillBonusLevels            int     `json:"skillBonusLevels,omitempty"`
	XpMultiplierPerLevelPercent float64 `json:"xpMultiplierPerLevelPercent,omitempty"`
	XpMultiplierPercent         float64 `json:"xpMultiplierPercent,omitempty"`
	XpPenalty                   float64 `json:"xpPenalty,omitempty"`
}

// GameDates contains server timestamps for daily rollover and energy regeneration.
type GameDates struct {
	DailyMissionRegenAt         string `json:"dailyMissionRegenAt,omitempty"`
	MonthlyMissionRegenAt       string `json:"monthlyMissionRegenAt,omitempty"`
	NextCongressElectionsAt     string `json:"nextCongressElectionsAt,omitempty"`
	NextDayAt                   string `json:"nextDayAt,omitempty"`
	NextMonthAt                 string `json:"nextMonthAt,omitempty"`
	NextPresidentialElectionsAt string `json:"nextPresidentialElectionsAt,omitempty"`
	NextRegenAt                 string `json:"nextRegenAt,omitempty"`
	PreviousDayAt               string `json:"previousDayAt,omitempty"`
	WeeklyMissionRegenAt        string `json:"weeklyMissionRegenAt,omitempty"`
}

// Government contains the leadership and legislative positions for a country.
type Government struct {
	ID              string               `json:"_id,omitempty"`
	CongressMembers []string             `json:"congressMembers,omitempty"`
	CountryID       string               `json:"country,omitempty"`
	Dates           GovernmentDates      `json:"dates,omitempty"`
	Ministers       []GovernmentMinister `json:"ministers,omitempty"`
	President       string               `json:"president,omitempty"`
	VicePresident   string               `json:"vicePresident,omitempty"`
}

// GovernmentDates is a nested component of Government.
type GovernmentDates struct {
	AnnouncementCreatedAts []string `json:"announcementCreatedAts,omitempty"`
	LastAnnouncementAt     string   `json:"lastAnnouncementAt,omitempty"`
}

// LiveBattleData represents data returned by the WarEra API.
type LiveBattleData struct {
	Battle LiveBattleSummary `json:"battle,omitempty"`
	Round  LiveRoundSummary  `json:"round,omitempty"`
}

// LiveBattleSummary is a nested component of LiveBattleData.
type LiveBattleSummary struct {
	AttackerCountryOrders     []string             `json:"attackerCountryOrders,omitempty"`
	AttackerMoneyPer1kDamages float64              `json:"attackerMoneyPer1kDamages,omitempty"`
	DefenderCountryOrders     []string             `json:"defenderCountryOrders,omitempty"`
	DefenderMoneyPer1kDamages float64              `json:"defenderMoneyPer1kDamages,omitempty"`
	IsActive                  bool                 `json:"isActive,omitempty"`
	RoundHistory              []BattleRoundHistory `json:"roundHistory,omitempty"`
	RoundIds                  []string             `json:"roundIds,omitempty"`
}

// LiveRoundSummary is a nested component of LiveBattleData.
type LiveRoundSummary struct {
	ActualTickPoints float64 `json:"actualTickPoints,omitempty"`
	AttackerDamages  float64 `json:"attackerDamages,omitempty"`
	AttackerPoints   float64 `json:"attackerPoints,omitempty"`
	DefenderDamages  float64 `json:"defenderDamages,omitempty"`
	DefenderPoints   float64 `json:"defenderPoints,omitempty"`
	IsActive         bool    `json:"isActive,omitempty"`
	NextTickAt       string  `json:"nextTickAt,omitempty"`
	RoundID          string  `json:"roundId,omitempty"`
}

// MULeveling is a nested component of MilitaryUnit.
type MULeveling struct {
	Level          int     `json:"level,omitempty"`
	MonthlyDamages float64 `json:"monthlyDamages,omitempty"`
}

// MURoles is a nested component of MilitaryUnit.
type MURoles struct {
	Commanders []string `json:"commanders,omitempty"`
	Managers   []string `json:"managers,omitempty"`
}

// MilitaryUnit is a player-managed military organization.
type MilitaryUnit struct {
	ID                   string                          `json:"_id,omitempty"`
	ActiveUpgradeLevels  MilitaryUnitActiveUpgradeLevels `json:"activeUpgradeLevels,omitempty"`
	AvatarURL            string                          `json:"avatarUrl,omitempty"`
	CountryID            string                          `json:"country,omitempty"`
	CreatedAt            string                          `json:"createdAt,omitempty"`
	InvestedMoneyByUsers map[string]float64              `json:"investedMoneyByUsers,omitempty"`
	LastAnnouncementAt   string                          `json:"lastAnnouncementAt,omitempty"`
	Leveling             MULeveling                      `json:"leveling,omitempty"`
	Members              []string                        `json:"members,omitempty"`
	MercenaryReputation  float64                         `json:"mercenaryReputation,omitempty"`
	Name                 string                          `json:"name,omitempty"`
	Rankings             map[string]RankingItem          `json:"rankings,omitempty"`
	Region               string                          `json:"region,omitempty"`
	Roles                MURoles                         `json:"roles,omitempty"`
	SponsorCountryID     string                          `json:"sponsorCountry,omitempty"`
	Type                 string                          `json:"type,omitempty"`
	UpdatedAt            string                          `json:"updatedAt,omitempty"`
	UserID               string                          `json:"user,omitempty"`
}

// MilitaryUnitActiveUpgradeLevels is a nested component of MilitaryUnit.
type MilitaryUnitActiveUpgradeLevels struct {
	Dormitories  int `json:"dormitories,omitempty"`
	Headquarters int `json:"headquarters,omitempty"`
}

// Ranking contains a global ranking and its entries.
type Ranking struct {
	Count      int               `json:"count,omitempty"`
	Items      []RankingEntry    `json:"items,omitempty"`
	TierValues RankingTierValues `json:"tierValues,omitempty"`
	Type       string            `json:"type,omitempty"`
	UpdatedAt  string            `json:"updatedAt,omitempty"`
}

// RankingItem represents a tiered leaderboard rank entry.
type RankingItem struct {
	Rank  int     `json:"rank,omitempty"`
	Tier  string  `json:"tier,omitempty"`
	Value float64 `json:"value,omitempty"`
}

// RankingEntry is one user, country, military unit, or alliance in a ranking.
// Which entity IDs are set depends on the ranking type; user rankings also
// carry the user's country and military unit.
type RankingEntry struct {
	ID         string  `json:"_id,omitempty"`
	AllianceID string  `json:"alliance,omitempty"`
	CountryID  string  `json:"country,omitempty"`
	MUID       string  `json:"mu,omitempty"`
	Rank       int     `json:"rank,omitempty"`
	Tier       string  `json:"tier,omitempty"`
	Type       string  `json:"type,omitempty"`
	UserID     string  `json:"user,omitempty"`
	Value      float64 `json:"value,omitempty"`
}

// RankingTierValues is a nested component of Ranking.
type RankingTierValues struct {
	Bronze   float64 `json:"bronze,omitempty"`
	Diamond  float64 `json:"diamond,omitempty"`
	Gold     float64 `json:"gold,omitempty"`
	Master   float64 `json:"master,omitempty"`
	Platinum float64 `json:"platinum,omitempty"`
	Silver   float64 `json:"silver,omitempty"`
}

// Region is a map area with economic, geographic, and military state.
type Region struct {
	ID                           string                    `json:"_id,omitempty"`
	ActiveBattle                 *RegionActiveBattleValue  `json:"activeBattle,omitempty"`
	ActiveUpgradeLevels          RegionActiveUpgradeLevels `json:"activeUpgradeLevels,omitempty"`
	BaseDevelopment              float64                   `json:"baseDevelopment,omitempty"`
	Biome                        string                    `json:"biome,omitempty"`
	Climate                      string                    `json:"climate,omitempty"`
	Code                         string                    `json:"code,omitempty"`
	CompanyCount                 int                       `json:"companyCount,omitempty"`
	CountryCode                  string                    `json:"countryCode,omitempty"`
	CountryID                    string                    `json:"country,omitempty"`
	CurrentPopulation            int                       `json:"currentPopulation,omitempty"`
	Dates                        RegionDates               `json:"dates,omitempty"`
	Development                  float64                   `json:"development,omitempty"`
	EmployeeCount                int                       `json:"employeeCount,omitempty"`
	HasCoast                     bool                      `json:"hasCoast,omitempty"`
	Income                       *IncomeBreakdown          `json:"income,omitempty"`
	InitialCountry               string                    `json:"initialCountry,omitempty"`
	IsCapital                    bool                      `json:"isCapital,omitempty"`
	IsLinkedToCapital            bool                      `json:"isLinkedToCapital,omitempty"`
	LastBattleEndedAt            string                    `json:"lastBattleEndedAt,omitempty"`
	LastResistanceContributionAt string                    `json:"lastResistanceContributionAt,omitempty"`
	LastRevoltEndedAt            string                    `json:"lastRevoltEndedAt,omitempty"`
	MainCity                     string                    `json:"mainCity,omitempty"`
	Name                         string                    `json:"name,omitempty"`
	Neighbors                    []string                  `json:"neighbors,omitempty"`
	Population                   int                       `json:"population,omitempty"`
	Position                     []float64                 `json:"position,omitempty"`
	Residents                    int                       `json:"residents,omitempty"`
	Resistance                   float64                   `json:"resistance,omitempty"`
	ResistanceMax                float64                   `json:"resistanceMax,omitempty"`
	Stats                        RegionStats               `json:"stats,omitempty"`
	StrategicResource            string                    `json:"strategicResource,omitempty"`
	UpgradesV2                   RegionUpgradesV2          `json:"upgradesV2,omitempty"`
}

// RegionActiveBattle is a nested component of Region.
type RegionActiveBattle struct {
	ID            string                  `json:"_id,omitempty"`
	Attacker      BattleSideSummary       `json:"attacker,omitempty"`
	CreatedAt     string                  `json:"createdAt,omitempty"`
	CurrentRound  *BattleRound            `json:"currentRound,omitempty"`
	Defender      BattleSideSummary       `json:"defender,omitempty"`
	IsActive      bool                    `json:"isActive,omitempty"`
	IsBigBattle   bool                    `json:"isBigBattle,omitempty"`
	Rounds        []string                `json:"rounds,omitempty"`
	RoundsHistory []BattleRoundHistory    `json:"roundsHistory,omitempty"`
	RoundsToWin   int                     `json:"roundsToWin,omitempty"`
	Stats         RegionActiveBattleStats `json:"stats,omitempty"`
	Type          string                  `json:"type,omitempty"`
	UpdatedAt     string                  `json:"updatedAt,omitempty"`
	WarID         string                  `json:"war,omitempty"`
}

// RegionActiveBattleRoundsHistory is an alias for [BattleRoundHistory].
type RegionActiveBattleRoundsHistory = BattleRoundHistory

// RegionActiveBattleStats is a nested component of RegionActiveBattle.
type RegionActiveBattleStats struct {
	HitCount int `json:"hitCount,omitempty"`
}

// RegionActiveUpgradeLevels is a nested component of Region.
type RegionActiveUpgradeLevels struct {
	Base               int `json:"base,omitempty"`
	Bunker             int `json:"bunker,omitempty"`
	PacificationCenter int `json:"pacificationCenter,omitempty"`
}

// RegionDates is a nested component of Region.
type RegionDates struct {
	LastOwnershipChangeAt string `json:"lastOwnershipChangeAt,omitempty"`
}

// RegionDeposit is a nested component of Region.
type RegionDeposit struct {
	BonusPercent float64 `json:"bonusPercent,omitempty"`
	EndsAt       string  `json:"endsAt,omitempty"`
	StartsAt     string  `json:"startsAt,omitempty"`
	Type         string  `json:"type,omitempty"`
}

// RegionStats is a nested component of Region.
type RegionStats struct {
	InvestedMoney float64 `json:"investedMoney,omitempty"`
}

// RegionUpgradesV2 is a nested component of Region.
type RegionUpgradesV2 struct {
	ActiveConstructionCount int                      `json:"activeConstructionCount,omitempty"`
	Upgrades                RegionUpgradesV2Upgrades `json:"upgrades,omitempty"`
}

// RegionUpgradesV2Upgrades is a nested component of RegionUpgradesV2.
type RegionUpgradesV2Upgrades struct {
	Base   RegionUpgradesV2UpgradesBase   `json:"base,omitempty"`
	Bunker RegionUpgradesV2UpgradesBunker `json:"bunker,omitempty"`
}

// RegionUpgradesV2UpgradesBase is a nested component of RegionUpgradesV2Upgrades.
type RegionUpgradesV2UpgradesBase struct {
	ConstructionEndedAt   string                          `json:"constructionEndedAt,omitempty"`
	ConstructionPoints    float64                         `json:"constructionPoints,omitempty"`
	ConstructionStartedAt string                          `json:"constructionStartedAt,omitempty"`
	InvestedMoney         float64                         `json:"investedMoney,omitempty"`
	IsUnderConstruction   *bool                           `json:"isUnderConstruction,omitempty"`
	LastConstructions     []RegionUpgradeLastConstruction `json:"lastConstructions,omitempty"`
	Level                 int                             `json:"level,omitempty"`
	Status                string                          `json:"status,omitempty"`
	StatusChangedAt       string                          `json:"statusChangedAt,omitempty"`
}

// RegionUpgradesV2UpgradesBunker is a nested component of RegionUpgradesV2Upgrades.
type RegionUpgradesV2UpgradesBunker struct {
	ConstructionEndedAt   string                          `json:"constructionEndedAt,omitempty"`
	ConstructionPoints    float64                         `json:"constructionPoints,omitempty"`
	ConstructionStartedAt string                          `json:"constructionStartedAt,omitempty"`
	InvestedMoney         float64                         `json:"investedMoney,omitempty"`
	IsUnderConstruction   *bool                           `json:"isUnderConstruction,omitempty"`
	LastConstructions     []RegionUpgradeLastConstruction `json:"lastConstructions,omitempty"`
	Level                 int                             `json:"level,omitempty"`
	Status                string                          `json:"status,omitempty"`
	StatusChangedAt       string                          `json:"statusChangedAt,omitempty"`
}

// RegionUpgradeLastConstruction records a completed upgrade construction.
type RegionUpgradeLastConstruction struct {
	Construction   float64 `json:"construction,omitempty"`
	ConstructionAt string  `json:"constructionAt,omitempty"`
	UserID         string  `json:"user,omitempty"`
}

// RegionUpgradesV2UpgradesBunkerLastConstructions is an alias for [RegionUpgradeLastConstruction].
type RegionUpgradesV2UpgradesBunkerLastConstructions = RegionUpgradeLastConstruction

// Round is one stage of a battle and includes combat state and results.
type Round struct {
	ID        string               `json:"_id,omitempty"`
	Attacker  RoundSideCombatStats `json:"attacker,omitempty"`
	Battle    string               `json:"battle,omitempty"`
	CreatedAt string               `json:"createdAt,omitempty"`
	Defender  RoundSideCombatStats `json:"defender,omitempty"`
	IsActive  bool                 `json:"isActive,omitempty"`
	Live      RoundLiveState       `json:"live,omitempty"`
	Number    int                  `json:"number,omitempty"`
	UpdatedAt string               `json:"updatedAt,omitempty"`
}

// RoundEquipment is a nested component of RoundLastHit.
type RoundEquipment struct {
	ID                string          `json:"_id,omitempty"`
	Code              string          `json:"code,omitempty"`
	LastAcquisitionAt string          `json:"lastAcquisitionAt,omitempty"`
	MaxState          float64         `json:"maxState,omitempty"`
	Quantity          int             `json:"quantity,omitempty"`
	Skills            EquipmentSkills `json:"skills,omitempty"`
	State             float64         `json:"state,omitempty"`
	Type              string          `json:"type,omitempty"`
}

// RoundLastHit is a nested component of RoundSideCombatStats.
type RoundLastHit struct {
	ID            string           `json:"_id,omitempty"`
	Ammo          string           `json:"ammo,omitempty"`
	Damages       float64          `json:"damages,omitempty"`
	Equipments    []RoundEquipment `json:"equipments,omitempty"`
	HitAt         string           `json:"hitAt,omitempty"`
	IsCriticalHit bool             `json:"isCriticalHit,omitempty"`
	IsMissed      bool             `json:"isMissed,omitempty"`
	MUID          string           `json:"mu,omitempty"`
	UserID        string           `json:"user,omitempty"`
	Weapon        *RoundWeapon     `json:"weapon,omitempty"`
}

// RoundLastHits represents data returned by the WarEra API.
type RoundLastHits struct {
	Attacker []RoundLastHit `json:"attacker,omitempty"`
	Defender []RoundLastHit `json:"defender,omitempty"`
}

// RoundLiveState is a nested component of BattleRound.
type RoundLiveState struct {
	ActualTickPoints float64 `json:"actualTickPoints,omitempty"`
	NextTickAt       string  `json:"nextTickAt,omitempty"`
	TicksCount       int     `json:"ticksCount,omitempty"`
}

// RoundSideCombatStats is a nested component of BattleRound.
type RoundSideCombatStats struct {
	CountryID string         `json:"country,omitempty"`
	Damages   float64        `json:"damages,omitempty"`
	HitCount  int            `json:"hitCount,omitempty"`
	LastHits  []RoundLastHit `json:"lastHits,omitempty"`
	Points    float64        `json:"points,omitempty"`
}

// RoundWeapon is a nested component of RoundLastHit.
type RoundWeapon struct {
	ID                string       `json:"_id,omitempty"`
	Code              string       `json:"code,omitempty"`
	LastAcquisitionAt string       `json:"lastAcquisitionAt,omitempty"`
	MaxState          float64      `json:"maxState,omitempty"`
	Quantity          int          `json:"quantity,omitempty"`
	Skills            WeaponSkills `json:"skills,omitempty"`
	State             float64      `json:"state,omitempty"`
}

// SearchResult contains matching entity IDs grouped by entity type.
type SearchResult struct {
	AllianceIDs []string `json:"allianceIds,omitempty"`
	CountryIDs  []string `json:"countryIds,omitempty"`
	HasData     bool     `json:"hasData,omitempty"`
	MUIDs       []string `json:"muIds,omitempty"`
	PartyIDs    []string `json:"partyIds,omitempty"`
	RegionIDs   []string `json:"regionIds,omitempty"`
	UserIDs     []string `json:"userIds,omitempty"`
}

// TopOrder is a nested component of TopOrders.
type TopOrder struct {
	ID        string  `json:"_id,omitempty"`
	CountryID string  `json:"country,omitempty"`
	ItemCode  string  `json:"itemCode,omitempty"`
	MUID      string  `json:"mu,omitempty"`
	OfferAt   string  `json:"offerAt,omitempty"`
	PartyID   string  `json:"party,omitempty"`
	Price     float64 `json:"price,omitempty"`
	Quantity  int     `json:"quantity,omitempty"`
	Type      string  `json:"type,omitempty"`
	UserID    string  `json:"user,omitempty"`
}

// TopOrders contains the best buy and sell orders for a tradeable item.
type TopOrders struct {
	BuyOrders  []TopOrder `json:"buyOrders,omitempty"`
	SellOrders []TopOrder `json:"sellOrders,omitempty"`
}

// Upgrade represents data returned by the WarEra API.
type Upgrade struct {
	ID                  string      `json:"_id,omitempty"`
	CreatedAt           string      `json:"createdAt,omitempty"`
	DependantUsersCount int         `json:"dependantUsersCount,omitempty"`
	InvestedConcrete    float64     `json:"investedConcrete,omitempty"`
	InvestedMoney       float64     `json:"investedMoney,omitempty"`
	InvestedSteel       float64     `json:"investedSteel,omitempty"`
	LastDowngradeAt     string      `json:"lastDowngradeAt,omitempty"`
	LastUpgradeAt       string      `json:"lastUpgradeAt,omitempty"`
	Level               int         `json:"level,omitempty"`
	Region              string      `json:"region,omitempty"`
	Status              string      `json:"status,omitempty"`
	StatusChangedAt     string      `json:"statusChangedAt,omitempty"`
	UpdatedAt           string      `json:"updatedAt,omitempty"`
	UpgradeType         UpgradeType `json:"upgradeType,omitempty"`
	WillBeActiveAt      string      `json:"willBeActiveAt,omitempty"`
}

// User is a full player profile, including skills, statistics, and equipment.
type User struct {
	ID                    string                 `json:"_id,omitempty"`
	AvailableColorSchemes []string               `json:"availableColorSchemes,omitempty"`
	AvatarURL             string                 `json:"avatarUrl,omitempty"`
	Buffs                 *UserBuffs             `json:"buffs,omitempty"`
	CanOnboard            bool                   `json:"canOnboard,omitempty"`
	CountryID             string                 `json:"country,omitempty"`
	CreatedAt             string                 `json:"createdAt,omitempty"`
	Dates                 UserDates              `json:"dates,omitempty"`
	Discord               UserDiscord            `json:"discord,omitempty"`
	EmailVerified         bool                   `json:"emailVerified,omitempty"`
	EnergyGift            *UserEnergyGift        `json:"energyGift,omitempty"`
	Equipment             UserEquipment          `json:"equipment,omitempty"`
	EquippedSkinKeys      json.RawMessage        `json:"equippedSkinKeys,omitempty"`
	FinishedTours         UserFinishedTours      `json:"finishedTours,omitempty"`
	Infos                 UserInfos              `json:"infos,omitempty"`
	IsActive              bool                   `json:"isActive,omitempty"`
	Leveling              UserLeveling           `json:"leveling,omitempty"`
	Location              string                 `json:"location,omitempty"`
	MilitaryRank          int                    `json:"militaryRank,omitempty"`
	Missions              UserMissions           `json:"missions,omitempty"`
	MUID                  string                 `json:"mu,omitempty"`
	MuMaxLevelRewarded    float64                `json:"muMaxLevelRewarded,omitempty"`
	Party                 string                 `json:"party,omitempty"`
	Preferences           UserPreferences        `json:"preferences,omitempty"`
	Rankings              map[string]RankingItem `json:"rankings,omitempty"`
	Region                string                 `json:"region,omitempty"`
	ShouldUpdateProfile   bool                   `json:"shouldUpdateProfile,omitempty"`
	Skills                UserSkills             `json:"skills,omitempty"`
	SkillVersion          int                    `json:"skillVersion,omitempty"`
	Stats                 UserStats              `json:"stats,omitempty"`
	UpdatedAt             string                 `json:"updatedAt,omitempty"`
	Username              string                 `json:"username,omitempty"`
	UsernameLower         string                 `json:"usernameLower,omitempty"`
	Warning               *UserWarning           `json:"warning,omitempty"`
}

// UserDates is a nested component of User.
type UserDates struct {
	CitizenSince               string `json:"citizenSince,omitempty"`
	LastCitizenshipChangeAt    string `json:"lastCitizenshipChangeAt,omitempty"`
	LastConnectionAt           string `json:"lastConnectionAt,omitempty"`
	LastCountryMessageCheckAt  string `json:"lastCountryMessageCheckAt,omitempty"`
	LastDailyCalendarClaimedAt string `json:"lastDailyCalendarClaimedAt,omitempty"`
	LastEventsCheckAt          string `json:"lastEventsCheckAt,omitempty"`
	LastFiredFromGovAt         string `json:"lastFiredFromGovAt,omitempty"`
	LastFreeTravelHomeAt       string `json:"lastFreeTravelHomeAt,omitempty"`
	LastGlobalMessageCheckAt   string `json:"lastGlobalMessageCheckAt,omitempty"`
	LastHelpAskedAt            string `json:"lastHelpAskedAt,omitempty"`
	LastManagerCareAt          string `json:"lastManagerCareAt,omitempty"`
	LastNotificationsCheckAt   string `json:"lastNotificationsCheckAt,omitempty"`
	LastSeenAnnouncementAt     string `json:"lastSeenAnnouncementAt,omitempty"`
	LastSkillsResetAt          string `json:"lastSkillsResetAt,omitempty"`
	LastWorkAt                 string `json:"lastWorkAt,omitempty"`
}

// UserDiscord is a nested component of User.
type UserDiscord struct {
	IsVerified bool `json:"isVerified,omitempty"`
}

// UserEquipment is a nested component of User.
type UserEquipment struct {
	Ammo   string `json:"ammo,omitempty"`
	Boots  string `json:"boots,omitempty"`
	Chest  string `json:"chest,omitempty"`
	Gloves string `json:"gloves,omitempty"`
	Helmet string `json:"helmet,omitempty"`
	Pants  string `json:"pants,omitempty"`
	Weapon string `json:"weapon,omitempty"`
}

// UserFinishedTours is a nested component of User.
type UserFinishedTours struct {
	Onboarding bool `json:"onboarding,omitempty"`
}

// UserLeveling is a nested component of User.
type UserLeveling struct {
	AvailableSkillPoints int     `json:"availableSkillPoints,omitempty"`
	DailyXPLeft          float64 `json:"dailyXpLeft,omitempty"`
	FreeReset            float64 `json:"freeReset,omitempty"`
	Level                int     `json:"level,omitempty"`
	Prestige             int     `json:"prestige,omitempty"`
	PrestigeLevel        int     `json:"prestigeLevel,omitempty"`
	SpentSkillPoints     int     `json:"spentSkillPoints,omitempty"`
	TotalSkillPoints     int     `json:"totalSkillPoints,omitempty"`
	TotalXP              float64 `json:"totalXp,omitempty"`
	XP                   float64 `json:"xp,omitempty"`
}

// UserLite is a compact player profile for listings and search results.
type UserLite struct {
	ID               string                 `json:"_id,omitempty"`
	AvatarURL        string                 `json:"avatarUrl,omitempty"`
	Buffs            *UserBuffs             `json:"buffs,omitempty"`
	CountryID        string                 `json:"country,omitempty"`
	CreatedAt        string                 `json:"createdAt,omitempty"`
	Dates            UserDates              `json:"dates,omitempty"`
	EquippedSkinKeys json.RawMessage        `json:"equippedSkinKeys,omitempty"`
	Infos            UserInfos              `json:"infos,omitempty"`
	IsActive         bool                   `json:"isActive,omitempty"`
	Leveling         UserLeveling           `json:"leveling,omitempty"`
	MilitaryRank     int                    `json:"militaryRank,omitempty"`
	MUID             string                 `json:"mu,omitempty"`
	Rankings         map[string]RankingItem `json:"rankings,omitempty"`
	Stats            UserStats              `json:"stats,omitempty"`
	Username         string                 `json:"username,omitempty"`
}

// UserMissions is a nested component of User.
type UserMissions struct {
	ClaimedAt               UserMissionsClaimedAt   `json:"claimedAt,omitempty"`
	ClaimedCount            int                     `json:"claimedCount,omitempty"`
	GeneratedAt             UserMissionsGeneratedAt `json:"generatedAt,omitempty"`
	RerolledDailyMissions   int                     `json:"rerolledDailyMissions,omitempty"`
	RerolledMonthlyMissions int                     `json:"rerolledMonthlyMissions,omitempty"`
	RerolledWeeklyMissions  int                     `json:"rerolledWeeklyMissions,omitempty"`
	TotalXpClaimed          float64                 `json:"totalXpClaimed,omitempty"`
}

// UserMissionsClaimedAt is a nested component of UserMissions.
type UserMissionsClaimedAt struct {
	Daily    string `json:"daily,omitempty"`
	Monthly  string `json:"monthly,omitempty"`
	Starting string `json:"starting,omitempty"`
	Weekly   string `json:"weekly,omitempty"`
}

// UserPreferences is a nested component of User.
type UserPreferences struct {
	AppFont            string `json:"appFont,omitempty"`
	AppPattern         string `json:"appPattern,omitempty"`
	AutoEquipMode      string `json:"autoEquipMode,omitempty"`
	AutoReplaceOnBreak bool   `json:"autoReplaceOnBreak,omitempty"`
	Colorblind         bool   `json:"colorblind,omitempty"`
	Locale             string `json:"locale,omitempty"`
	Sfx                bool   `json:"sfx,omitempty"`
	Vibrations         bool   `json:"vibrations,omitempty"`
}

// UserSkills contains a player's skills.
type UserSkills struct {
	Armor           UserSkill `json:"armor,omitempty"`
	Attack          UserSkill `json:"attack,omitempty"`
	Companies       UserSkill `json:"companies,omitempty"`
	CriticalChance  UserSkill `json:"criticalChance,omitempty"`
	CriticalDamages UserSkill `json:"criticalDamages,omitempty"`
	Dodge           UserSkill `json:"dodge,omitempty"`
	Energy          UserSkill `json:"energy,omitempty"`
	Health          UserSkill `json:"health,omitempty"`
	Hunger          UserSkill `json:"hunger,omitempty"`
	LootChance      UserSkill `json:"lootChance,omitempty"`
	Management      UserSkill `json:"management,omitempty"`
	Precision       UserSkill `json:"precision,omitempty"`
	Production      UserSkill `json:"production,omitempty"`
	Stamina         UserSkill `json:"stamina,omitempty"`
}

// UserStats is a nested component of User.
type UserStats struct {
	Case1                 UserStatsCase1      `json:"case1,omitempty"`
	Case2                 UserStatsCase2      `json:"case2,omitempty"`
	DailyRewardBestStreak int                 `json:"dailyRewardBestStreak,omitempty"`
	DailyRewardStreak     int                 `json:"dailyRewardStreak,omitempty"`
	DamagesCount          int                 `json:"damagesCount,omitempty"`
	MilestoneTiersClaimed int                 `json:"milestoneTiersClaimed,omitempty"`
	MilestoneXpClaimed    float64             `json:"milestoneXpClaimed,omitempty"`
	Toolbox               UserStatsToolbox    `json:"toolbox,omitempty"`
	Wealth                UserStatsWealth     `json:"wealth,omitempty"`
	WoodenCase            UserStatsWoodenCase `json:"woodenCase,omitempty"`
	WorksCount            int                 `json:"worksCount,omitempty"`
}

// UserStatsCase1 is a nested component of UserStats.
type UserStatsCase1 struct {
	ByRarity    UserStatsCase1ByRarity `json:"byRarity,omitempty"`
	OpenedCount int                    `json:"openedCount,omitempty"`
}

// UserStatsCase1ByRarity is a nested component of UserStatsCase1.
type UserStatsCase1ByRarity struct {
	Common    float64 `json:"common,omitempty"`
	Epic      float64 `json:"epic,omitempty"`
	Legendary float64 `json:"legendary,omitempty"`
	Mythic    float64 `json:"mythic,omitempty"`
	Rare      float64 `json:"rare,omitempty"`
	Uncommon  float64 `json:"uncommon,omitempty"`
}

// UserStatsCase2 is a nested component of UserStats.
type UserStatsCase2 struct {
	ByRarity    UserStatsCase2ByRarity `json:"byRarity,omitempty"`
	OpenedCount int                    `json:"openedCount,omitempty"`
}

// UserStatsCase2ByRarity is a nested component of UserStatsCase2.
type UserStatsCase2ByRarity struct {
	Epic      float64 `json:"epic,omitempty"`
	Legendary float64 `json:"legendary,omitempty"`
	Rare      float64 `json:"rare,omitempty"`
	Uncommon  float64 `json:"uncommon,omitempty"`
}

// UserStatsWealth is a nested component of UserStats.
type UserStatsWealth struct {
	Companies  float64 `json:"companies,omitempty"`
	Equipments float64 `json:"equipments,omitempty"`
	Items      float64 `json:"items,omitempty"`
	Money      float64 `json:"money,omitempty"`
	Total      float64 `json:"total,omitempty"`
	Weapons    float64 `json:"weapons,omitempty"`
}

// UserStatsWoodenCase is a nested component of UserStats.
type UserStatsWoodenCase struct {
	ByItem      UserStatsWoodenCaseByItem `json:"byItem,omitempty"`
	OpenedCount int                       `json:"openedCount,omitempty"`
}

// UserStatsWoodenCaseByItem is a nested component of UserStatsWoodenCase.
type UserStatsWoodenCaseByItem struct {
	Ammo       float64 `json:"ammo,omitempty"`
	Bread      float64 `json:"bread,omitempty"`
	Coca       float64 `json:"coca,omitempty"`
	Cocain     float64 `json:"cocain,omitempty"`
	Concrete   float64 `json:"concrete,omitempty"`
	CookedFish float64 `json:"cookedFish,omitempty"`
	Fish       float64 `json:"fish,omitempty"`
	Grain      float64 `json:"grain,omitempty"`
	HeavyAmmo  float64 `json:"heavyAmmo,omitempty"`
	Iron       float64 `json:"iron,omitempty"`
	Lead       float64 `json:"lead,omitempty"`
	LightAmmo  float64 `json:"lightAmmo,omitempty"`
	Limestone  float64 `json:"limestone,omitempty"`
	Livestock  float64 `json:"livestock,omitempty"`
	Oil        float64 `json:"oil,omitempty"`
	Paper      float64 `json:"paper,omitempty"`
	Petroleum  float64 `json:"petroleum,omitempty"`
	Steak      float64 `json:"steak,omitempty"`
	Steel      float64 `json:"steel,omitempty"`
	Wood       float64 `json:"wood,omitempty"`
}

// WeaponItem is a nested component of CurrentEquipment.
type WeaponItem struct {
	ID                string       `json:"_id,omitempty"`
	Code              string       `json:"code,omitempty"`
	LastAcquisitionAt string       `json:"lastAcquisitionAt,omitempty"`
	MaxState          float64      `json:"maxState,omitempty"`
	Quantity          int          `json:"quantity,omitempty"`
	Skills            WeaponSkills `json:"skills,omitempty"`
	State             float64      `json:"state,omitempty"`
}

// WeaponSkills is a nested component of RoundWeapon.
type WeaponSkills struct {
	Attack         float64 `json:"attack,omitempty"`
	CriticalChance float64 `json:"criticalChance,omitempty"`
}

// UserSkill is one skill of a player. Fight skills report their weapon,
// equipment, and percentage modifiers; bar skills such as energy, health,
// hunger, and stamina report their current bar value and hourly regeneration.
type UserSkill struct {
	AmmoPercent         float64         `json:"ammoPercent,omitempty"`
	BuffsPercent        float64         `json:"buffsPercent,omitempty"`
	CurrentBarValue     float64         `json:"currentBarValue,omitempty"`
	DebuffsPercent      float64         `json:"debuffsPercent,omitempty"`
	Equipment           *float64        `json:"equipment,omitempty"`
	HourlyBarRegen      float64         `json:"hourlyBarRegen,omitempty"`
	Level               int             `json:"level,omitempty"`
	Limited             json.RawMessage `json:"limited,omitempty"`
	MilitaryRankPercent float64         `json:"militaryRankPercent,omitempty"`
	Overflow            json.RawMessage `json:"overflow,omitempty"`
	Prestige            int             `json:"prestige,omitempty"`
	// Strategic is the bonus from strategic resources held by the player's country.
	Strategic         float64  `json:"strategic,omitempty"`
	Total             float64  `json:"total,omitempty"`
	TotalAfterSoftCap *float64 `json:"totalAfterSoftCap,omitempty"`
	Value             float64  `json:"value,omitempty"`
	Weapon            *float64 `json:"weapon,omitempty"`
}

// IncomeBreakdown is a country's or region's daily income by source.
type IncomeBreakdown struct {
	Development             float64 `json:"development,omitempty"`
	EmployeeIncomeTax       float64 `json:"employeeIncomeTax,omitempty"`
	GovType                 float64 `json:"govType,omitempty"`
	ManagerIncomeTax        float64 `json:"managerIncomeTax,omitempty"`
	StrategicResourcesBonus float64 `json:"strategicResourcesBonus,omitempty"`
	Total                   float64 `json:"total,omitempty"`
}

// GovernmentMinister is a minister appointed to a country's government.
type GovernmentMinister struct {
	AppointedAt    string `json:"appointedAt,omitempty"`
	BadgeGrantedAt string `json:"badgeGrantedAt,omitempty"`
	// CustomTitle is a display title the president gave the minister.
	CustomTitle string   `json:"customTitle,omitempty"`
	Rights      []string `json:"rights,omitempty"`
	Title       string   `json:"title,omitempty"`
	UserID      string   `json:"user,omitempty"`
}

// GameConfigMapTravel reports which map travel features are enabled.
type GameConfigMapTravel struct {
	MapCasesEnabled       bool `json:"mapCasesEnabled,omitempty"`
	MovementEnabled       bool `json:"movementEnabled,omitempty"`
	PositionMarkerEnabled bool `json:"positionMarkerEnabled,omitempty"`
	StaminaBarEnabled     bool `json:"staminaBarEnabled,omitempty"`
}

// GameConfigChatSlowmode configures chat posting limits for low-level players.
type GameConfigChatSlowmode struct {
	DuplicateWindowSeconds float64                      `json:"duplicateWindowSeconds,omitempty"`
	LiftedAtLevel          int                          `json:"liftedAtLevel,omitempty"`
	Tiers                  []GameConfigChatSlowmodeTier `json:"tiers,omitempty"`
}

// GameConfigChatSlowmodeTier is the chat limit that applies from MinLevel.
type GameConfigChatSlowmodeTier struct {
	CooldownSeconds float64 `json:"cooldownSeconds,omitempty"`
	MaxPerHour      float64 `json:"maxPerHour,omitempty"`
	MinLevel        int     `json:"minLevel,omitempty"`
	NewDirectPerDay float64 `json:"newDirectPerDay,omitempty"`
}

// GameConfigSkillPresets sets how many skill presets a player can save.
type GameConfigSkillPresets struct {
	NonPremiumMax int `json:"nonPremiumMax,omitempty"`
	PremiumMax    int `json:"premiumMax,omitempty"`
}

// UserBuffs lists a player's active debuffs.
type UserBuffs struct {
	DebuffCodes []string `json:"debuffCodes,omitempty"`
	DebuffEndAt string   `json:"debuffEndAt,omitempty"`
}

// UserEnergyGift is an energy gift granted to a player.
type UserEnergyGift struct {
	Amount    float64 `json:"amount,omitempty"`
	ClaimedAt string  `json:"claimedAt,omitempty"`
	GrantedAt string  `json:"grantedAt,omitempty"`
}

// UserInfos holds a player's profile cosmetics and government roles.
type UserInfos struct {
	ColorScheme        string `json:"colorScheme,omitempty"`
	CongressMemberOf   string `json:"congressMemberOf,omitempty"`
	Font               string `json:"font,omitempty"`
	IsPremium          bool   `json:"isPremium,omitempty"`
	MinisterOf         string `json:"ministerOf,omitempty"`
	MinisterTitle      string `json:"ministerTitle,omitempty"`
	PremiumMonthsCount int    `json:"premiumMonthsCount,omitempty"`
	UsernameEffect     string `json:"usernameEffect,omitempty"`
}

// UserWarning is a moderation warning shown on a player's profile.
type UserWarning struct {
	IsAccepted bool   `json:"isAccepted,omitempty"`
	Message    string `json:"message,omitempty"`
}

// UserMissionsGeneratedAt records when each mission set was generated.
type UserMissionsGeneratedAt struct {
	Daily  string `json:"daily,omitempty"`
	Weekly string `json:"weekly,omitempty"`
}

// UserStatsToolbox counts the toolboxes a player opened.
type UserStatsToolbox struct {
	ByRarity    map[string]int `json:"byRarity,omitempty"`
	OpenedCount int            `json:"openedCount,omitempty"`
}

// GameConfigSkill holds the per-level configuration of one skill, keyed by
// skill name in [GameConfig.Skills] and by level in Levels.
type GameConfigSkill struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
	// SkillOverflow names the skill that receives this skill's value above its
	// cap, at SkillOverflowValue per point.
	SkillOverflow      string  `json:"skillOverflow,omitempty"`
	SkillOverflowValue float64 `json:"skillOverflowValue,omitempty"`
	SoftCap            float64 `json:"softCap,omitempty"`
}
