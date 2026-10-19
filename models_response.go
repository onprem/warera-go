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
	Comments int `json:"comments,omitempty"`
	Dislikes int `json:"dislikes,omitempty"`
	GemTips  int `json:"gemTips,omitempty"`
	Likes    int `json:"likes,omitempty"`
	Score    int `json:"score,omitempty"`
	Subs     int `json:"subs,omitempty"`
	Tips     int `json:"tips,omitempty"`
	Views    int `json:"views,omitempty"`
}

// Battle is a military conflict between attacker and defender sides.
type Battle struct {
	ID            string               `json:"_id,omitempty"`
	Attacker      BattleSideSummary    `json:"attacker,omitempty"`
	CreatedAt     string               `json:"createdAt,omitempty"`
	CurrentRound  *BattleRound         `json:"currentRound,omitempty"`
	Defender      BattleSideSummary    `json:"defender,omitempty"`
	IsActive      bool                 `json:"isActive,omitempty"`
	Rounds        []string             `json:"rounds,omitempty"`
	RoundsHistory []BattleRoundHistory `json:"roundsHistory,omitempty"`
	RoundsToWin   int                  `json:"roundsToWin,omitempty"`
	Stats         BattleStats          `json:"stats,omitempty"`
	Type          string               `json:"type,omitempty"`
	UpdatedAt     string               `json:"updatedAt,omitempty"`
	WarID         string               `json:"war,omitempty"`
}

// BattleRoundHistory summarizes the outcome of a completed battle round.
type BattleRoundHistory struct {
	AttackerDamages int    `json:"attackerDamages,omitempty"`
	AttackerPoints  int    `json:"attackerPoints,omitempty"`
	DefenderDamages int    `json:"defenderDamages,omitempty"`
	DefenderPoints  int    `json:"defenderPoints,omitempty"`
	WonBy           string `json:"wonBy,omitempty"`
}

// BattleOrder is an order assigned to a country or military unit on a battle side.
type BattleOrder struct {
	ID          string `json:"_id,omitempty"`
	Battle      string `json:"battle,omitempty"`
	CountryID   string `json:"country,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
	IsActive    bool   `json:"isActive,omitempty"`
	MUID        string `json:"mu,omitempty"`
	Priority    string `json:"priority,omitempty"`
	Rank        int    `json:"rank,omitempty"`
	Side        string `json:"side,omitempty"`
	SideCountry string `json:"sideCountry,omitempty"`
	Text        string `json:"text,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
	UserID      string `json:"user,omitempty"`
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
	Damages           int      `json:"damages,omitempty"`
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
	ConcreteInvested    int                        `json:"concreteInvested,omitempty"`
	CreatedAt           string                     `json:"createdAt,omitempty"`
	EstimatedValue      float64                    `json:"estimatedValue,omitempty"`
	IsFull              bool                       `json:"isFull,omitempty"`
	ItemCode            string                     `json:"itemCode,omitempty"`
	Name                string                     `json:"name,omitempty"`
	Production          float64                    `json:"production,omitempty"`
	Region              string                     `json:"region,omitempty"`
	UpdatedAt           string                     `json:"updatedAt,omitempty"`
	UserID              string                     `json:"user,omitempty"`
	WorkerCount         int                        `json:"workerCount,omitempty"`
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
	MapAccent          string                         `json:"mapAccent,omitempty"`
	Money              float64                        `json:"money,omitempty"`
	Name               string                         `json:"name,omitempty"`
	NonAggressionUntil map[string]string              `json:"nonAggressionUntil,omitempty"`
	Orgs               []string                       `json:"orgs,omitempty"`
	PinnedArticle      string                         `json:"pinnedArticle,omitempty"`
	Rankings           map[string]RankingItem         `json:"rankings,omitempty"`
	RulingParty        string                         `json:"rulingParty,omitempty"`
	Scheme             string                         `json:"scheme,omitempty"`
	SpecializedItem    string                         `json:"specializedItem,omitempty"`
	StrategicResources *CountryStrategicResourcesInfo `json:"strategicResources,omitempty"`
	Taxes              CountryTaxes                   `json:"taxes,omitempty"`
	Unrest             CountryUnrest                  `json:"unrest,omitempty"`
	UpdatedAt          string                         `json:"updatedAt,omitempty"`
	WarsWith           []string                       `json:"warsWith,omitempty"`
}

// CountryDates is a nested component of Country.
type CountryDates struct {
	LastAllianceLeftAt string `json:"lastAllianceLeftAt,omitempty"`
}

// CountryStrategicResourcesInfo is a nested component of Country.
type CountryStrategicResourcesInfo struct {
	Bonuses   CountryStrategicResourcesInfoBonuses   `json:"bonuses,omitempty"`
	Resources CountryStrategicResourcesInfoResources `json:"resources,omitempty"`
}

// CountryStrategicResourcesInfoBonuses is a nested component of CountryStrategicResourcesInfo.
type CountryStrategicResourcesInfoBonuses struct {
	DevelopmentPercent float64 `json:"developmentPercent,omitempty"`
	ProductionPercent  float64 `json:"productionPercent,omitempty"`
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
	Income   float64 `json:"income,omitempty"`
	Market   float64 `json:"market,omitempty"`
	SelfWork float64 `json:"selfWork,omitempty"`
}

// CountryUnrest is a nested component of Country.
type CountryUnrest struct {
	Bar                   int    `json:"bar,omitempty"`
	BarMax                int    `json:"barMax,omitempty"`
	LastContributionAt    string `json:"lastContributionAt,omitempty"`
	LastRevolutionEndedAt string `json:"lastRevolutionEndedAt,omitempty"`
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
	MaxState          int             `json:"maxState,omitempty"`
	Quantity          int             `json:"quantity,omitempty"`
	Skills            EquipmentSkills `json:"skills,omitempty"`
	State             int             `json:"state,omitempty"`
	Type              string          `json:"type,omitempty"`
}

// EquipmentSkills is a nested component of RoundEquipment.
type EquipmentSkills struct {
	Armor           int `json:"armor,omitempty"`
	CriticalDamages int `json:"criticalDamages,omitempty"`
	Dodge           int `json:"dodge,omitempty"`
	Precision       int `json:"precision,omitempty"`
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
	AllianceID      string   `json:"allianceId,omitempty"`
	AllianceName    string   `json:"allianceName,omitempty"`
	AttackerCountry string   `json:"attackerCountry,omitempty"`
	AttackerRegion  string   `json:"attackerRegion,omitempty"`
	Battle          string   `json:"battle,omitempty"`
	BonusPercent    int      `json:"bonusPercent,omitempty"`
	Countries       []string `json:"countries,omitempty"`
	CountryID       string   `json:"country,omitempty"`
	DefenderCountry string   `json:"defenderCountry,omitempty"`
	DefenderRegion  string   `json:"defenderRegion,omitempty"`
	DurationDays    int      `json:"durationDays,omitempty"`
	ItemCode        string   `json:"itemCode,omitempty"`
	Region          string   `json:"region,omitempty"`
	Type            string   `json:"type,omitempty"`
	Wars            []string `json:"wars,omitempty"`
	WonBy           string   `json:"wonBy,omitempty"`
}

// GameConfig contains server settings for game systems, skills, items, and rewards.
type GameConfig struct {
	Alliance               GameConfigAlliance               `json:"alliance,omitempty"`
	Announcement           GameConfigAnnouncement           `json:"announcement,omitempty"`
	Badge                  GameConfigBadge                  `json:"badge,omitempty"`
	Battle                 GameConfigBattle                 `json:"battle,omitempty"`
	CitizenshipApplication GameConfigCitizenshipApplication `json:"citizenshipApplication,omitempty"`
	Company                GameConfigCompany                `json:"company,omitempty"`
	Country                GameConfigCountry                `json:"country,omitempty"`
	Election               GameConfigElection               `json:"election,omitempty"`
	Government             GameConfigGovernment             `json:"government,omitempty"`
	Items                  GameConfigItems                  `json:"items,omitempty"`
	Law                    GameConfigLaw                    `json:"law,omitempty"`
	Loot                   GameConfigLoot                   `json:"loot,omitempty"`
	MercenaryContract      GameConfigMercenaryContract      `json:"mercenaryContract,omitempty"`
	MergingCost            GameConfigMergingCost            `json:"mergingCost,omitempty"`
	Mission                GameConfigMission                `json:"mission,omitempty"`
	MU                     GameConfigMu                     `json:"mu,omitempty"`
	Newspaper              GameConfigNewspaper              `json:"newspaper,omitempty"`
	Org                    GameConfigOrg                    `json:"org,omitempty"`
	Party                  GameConfigParty                  `json:"party,omitempty"`
	Referral               GameConfigReferral               `json:"referral,omitempty"`
	Region                 GameConfigRegion                 `json:"region,omitempty"`
	Skills                 GameConfigSkills                 `json:"skills,omitempty"`
	SubSkinReward          GameConfigSubSkinReward          `json:"subSkinReward,omitempty"`
	Unrest                 GameConfigUnrest                 `json:"unrest,omitempty"`
	Upgrade                GameConfigUpgrade                `json:"upgrade,omitempty"`
	UpgradesConfig         GameConfigUpgradesConfig         `json:"upgradesConfig,omitempty"`
	User                   GameConfigUser                   `json:"user,omitempty"`
	Worker                 GameConfigWorker                 `json:"worker,omitempty"`
}

// GameConfigAlliance is a nested component of GameConfig.
type GameConfigAlliance struct {
	LeaveCooldownDays          int    `json:"leaveCooldownDays,omitempty"`
	LeaveCooldownDisabledUntil string `json:"leaveCooldownDisabledUntil,omitempty"`
}

// GameConfigAnnouncement is a nested component of GameConfig.
type GameConfigAnnouncement struct {
	DurationHours int `json:"durationHours,omitempty"`
	MaxLength     int `json:"maxLength,omitempty"`
	MaxPerHour    int `json:"maxPerHour,omitempty"`
}

// GameConfigBadge is a nested component of GameConfig.
type GameConfigBadge struct {
	AlphaTester             GameConfigBadgeAlphaTester             `json:"alphaTester,omitempty"`
	BabyBoomer              GameConfigBadgeBabyBoomer              `json:"babyBoomer,omitempty"`
	BattleTerrainTop1       GameConfigBadgeBattleTerrainTop1       `json:"battleTerrainTop1,omitempty"`
	BattleTop1              GameConfigBadgeBattleTop1              `json:"battleTop1,omitempty"`
	BugFinder               GameConfigBadgeBugFinder               `json:"bugFinder,omitempty"`
	Coffee                  GameConfigBadgeCoffee                  `json:"coffee,omitempty"`
	CongressMember          GameConfigBadgeCongressMember          `json:"congressMember,omitempty"`
	CountryPresident        GameConfigBadgeCountryPresident        `json:"countryPresident,omitempty"`
	CountryTournamentWinner GameConfigBadgeCountryTournamentWinner `json:"countryTournamentWinner,omitempty"`
	ExploitFinder           GameConfigBadgeExploitFinder           `json:"exploitFinder,omitempty"`
	FoundingFather          GameConfigBadgeFoundingFather          `json:"foundingFather,omitempty"`
	GiftPremium             GameConfigBadgeGiftPremium             `json:"giftPremium,omitempty"`
	GiveawayWinner          GameConfigBadgeGiveawayWinner          `json:"giveawayWinner,omitempty"`
	GovMember               GameConfigBadgeGovMember               `json:"govMember,omitempty"`
	HardWorker              GameConfigBadgeHardWorker              `json:"hardWorker,omitempty"`
	MuTournamentWinner      GameConfigBadgeMuTournamentWinner      `json:"muTournamentWinner,omitempty"`
	PopularArticle          GameConfigBadgePopularArticle          `json:"popularArticle,omitempty"`
	Premium                 GameConfigBadgePremium                 `json:"premium,omitempty"`
	Prestige                GameConfigBadgePrestige                `json:"prestige,omitempty"`
	Referral                GameConfigBadgeReferral                `json:"referral,omitempty"`
	ResetSurvivor           GameConfigBadgeResetSurvivor           `json:"resetSurvivor,omitempty"`
	RoundTerrainTop1        GameConfigBadgeRoundTerrainTop1        `json:"roundTerrainTop1,omitempty"`
	RoundTop1               GameConfigBadgeRoundTop1               `json:"roundTop1,omitempty"`
	Staff                   GameConfigBadgeStaff                   `json:"staff,omitempty"`
	SugarDaddy              GameConfigBadgeSugarDaddy              `json:"sugarDaddy,omitempty"`
	Translator              GameConfigBadgeTranslator              `json:"translator,omitempty"`
	VicePresident           GameConfigBadgeVicePresident           `json:"vicePresident,omitempty"`
	Voted                   GameConfigBadgeVoted                   `json:"voted,omitempty"`
}

// GameConfigBadgeAlphaTester is a nested component of GameConfigBadge.
type GameConfigBadgeAlphaTester struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeBabyBoomer is a nested component of GameConfigBadge.
type GameConfigBadgeBabyBoomer struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeBattleTerrainTop1 is a nested component of GameConfigBadge.
type GameConfigBadgeBattleTerrainTop1 struct {
	Availabilities []string `json:"availabilities,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBadgeBattleTop1 is a nested component of GameConfigBadge.
type GameConfigBadgeBattleTop1 struct {
	Availabilities []string `json:"availabilities,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBadgeBugFinder is a nested component of GameConfigBadge.
type GameConfigBadgeBugFinder struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeCoffee is a nested component of GameConfigBadge.
type GameConfigBadgeCoffee struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeCongressMember is a nested component of GameConfigBadge.
type GameConfigBadgeCongressMember struct {
	Availabilities []string `json:"availabilities,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBadgeCountryPresident is a nested component of GameConfigBadge.
type GameConfigBadgeCountryPresident struct {
	Availabilities []string `json:"availabilities,omitempty"`
	CooldownDays   int      `json:"cooldownDays,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBadgeCountryTournamentWinner is a nested component of GameConfigBadge.
type GameConfigBadgeCountryTournamentWinner struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeExploitFinder is a nested component of GameConfigBadge.
type GameConfigBadgeExploitFinder struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeFoundingFather is a nested component of GameConfigBadge.
type GameConfigBadgeFoundingFather struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeGiftPremium is a nested component of GameConfigBadge.
type GameConfigBadgeGiftPremium struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeGiveawayWinner is a nested component of GameConfigBadge.
type GameConfigBadgeGiveawayWinner struct {
	Availabilities []string `json:"availabilities,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBadgeGovMember is a nested component of GameConfigBadge.
type GameConfigBadgeGovMember struct {
	Availabilities []string `json:"availabilities,omitempty"`
	CooldownDays   int      `json:"cooldownDays,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBadgeHardWorker is a nested component of GameConfigBadge.
type GameConfigBadgeHardWorker struct {
	Availabilities []string `json:"availabilities,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBadgeMuTournamentWinner is a nested component of GameConfigBadge.
type GameConfigBadgeMuTournamentWinner struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgePopularArticle is a nested component of GameConfigBadge.
type GameConfigBadgePopularArticle struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
	UniqueMetadataKey    string   `json:"uniqueMetadataKey,omitempty"`
}

// GameConfigBadgePremium is a nested component of GameConfigBadge.
type GameConfigBadgePremium struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgePrestige is a nested component of GameConfigBadge.
type GameConfigBadgePrestige struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeReferral is a nested component of GameConfigBadge.
type GameConfigBadgeReferral struct {
	Availabilities []string `json:"availabilities,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBadgeResetSurvivor is a nested component of GameConfigBadge.
type GameConfigBadgeResetSurvivor struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeRoundTerrainTop1 is a nested component of GameConfigBadge.
type GameConfigBadgeRoundTerrainTop1 struct {
	Availabilities []string `json:"availabilities,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBadgeRoundTop1 is a nested component of GameConfigBadge.
type GameConfigBadgeRoundTop1 struct {
	Availabilities []string `json:"availabilities,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBadgeStaff is a nested component of GameConfigBadge.
type GameConfigBadgeStaff struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeSugarDaddy is a nested component of GameConfigBadge.
type GameConfigBadgeSugarDaddy struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeTranslator is a nested component of GameConfigBadge.
type GameConfigBadgeTranslator struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               int      `json:"reward,omitempty"`
}

// GameConfigBadgeVicePresident is a nested component of GameConfigBadge.
type GameConfigBadgeVicePresident struct {
	Availabilities []string `json:"availabilities,omitempty"`
	CooldownDays   int      `json:"cooldownDays,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBadgeVoted is a nested component of GameConfigBadge.
type GameConfigBadgeVoted struct {
	Availabilities []string `json:"availabilities,omitempty"`
	Reward         int      `json:"reward,omitempty"`
}

// GameConfigBattle is a nested component of GameConfig.
type GameConfigBattle struct {
	AllianceDamagesBonusPercent          int            `json:"allianceDamagesBonusPercent,omitempty"`
	CasesPer1kDamagesInPool              int            `json:"casesPer1kDamagesInPool,omitempty"`
	CountryOrderBonusPercent             int            `json:"countryOrderBonusPercent,omitempty"`
	EnemyDamagesBonusPercent             int            `json:"enemyDamagesBonusPercent,omitempty"`
	GovMemberBountyRewardPercent         int            `json:"govMemberBountyRewardPercent,omitempty"`
	HealthCost                           int            `json:"healthCost,omitempty"`
	HitFor1CaseInPool                    int            `json:"hitFor1CaseInPool,omitempty"`
	LostAttackingRegionMalusPercent      int            `json:"lostAttackingRegionMalusPercent,omitempty"`
	MaxRounds                            int            `json:"maxRounds,omitempty"`
	MuOrderBonusPercent                  int            `json:"muOrderBonusPercent,omitempty"`
	OccupyingYourRegionsMalusPercent     int            `json:"occupyingYourRegionsMalusPercent,omitempty"`
	PatrioticBonusPercent                int            `json:"patrioticBonusPercent,omitempty"`
	PointsToWinRound                     int            `json:"pointsToWinRound,omitempty"`
	RankingsLootPercentPer1kDmg          float64        `json:"rankingsLootPercentPer1kDmg,omitempty"`
	RegionNotLinkedToCapitalMalusPercent int            `json:"regionNotLinkedToCapitalMalusPercent,omitempty"`
	RoundsToWin                          int            `json:"roundsToWin,omitempty"`
	SetCountryOrderMoneyCost             int            `json:"setCountryOrderMoneyCost,omitempty"`
	SetMuOrderMoneyCost                  int            `json:"setMuOrderMoneyCost,omitempty"`
	SetOrderMoneyCost                    int            `json:"setOrderMoneyCost,omitempty"`
	TickPoints                           map[string]int `json:"tickPoints,omitempty"`
}

// GameConfigCitizenshipApplication is a nested component of GameConfig.
type GameConfigCitizenshipApplication struct {
	AutoApprovalEnabled       bool `json:"autoApprovalEnabled,omitempty"`
	AutoApprovalMaxPopulation int  `json:"autoApprovalMaxPopulation,omitempty"`
}

// GameConfigCompany is a nested component of GameConfig.
type GameConfigCompany struct {
	ChangeItemCost                     int `json:"changeItemCost,omitempty"`
	ConstructionCostIncreasePerCompany int `json:"constructionCostIncreasePerCompany,omitempty"`
	DepositResourceBonus               int `json:"depositResourceBonus,omitempty"`
	DestructionValuePercent            int `json:"destructionValuePercent,omitempty"`
	MoveCost                           int `json:"moveCost,omitempty"`
}

// GameConfigCountry is a nested component of GameConfig.
type GameConfigCountry struct {
	MaxTaxAmount int `json:"maxTaxAmount,omitempty"`
}

// GameConfigElection is a nested component of GameConfig.
type GameConfigElection struct {
	CandidateDurationHours            int `json:"candidateDurationHours,omitempty"`
	CandidateMinLevel                 int `json:"candidateMinLevel,omitempty"`
	ElectionVoteDurationHours         int `json:"electionVoteDurationHours,omitempty"`
	VoteMinDaysSinceCitizenshipChange int `json:"voteMinDaysSinceCitizenshipChange,omitempty"`
	VoteMinLevel                      int `json:"voteMinLevel,omitempty"`
}

// GameConfigGovernment is a nested component of GameConfig.
type GameConfigGovernment struct {
	MemberWagePercent                  int `json:"memberWagePercent,omitempty"`
	NominationCooldownDays             int `json:"nominationCooldownDays,omitempty"`
	PresidentWagePercent               int `json:"presidentWagePercent,omitempty"`
	ResistanceDecreasedCooldownInHours int `json:"resistanceDecreasedCooldownInHours,omitempty"`
	ResistanceIncreasedCooldownInHours int `json:"resistanceIncreasedCooldownInHours,omitempty"`
}

// GameConfigItems is a nested component of GameConfig.
type GameConfigItems struct {
	Ammo       GameConfigItemsAmmo       `json:"ammo,omitempty"`
	Boots1     GameConfigItemsBoots1     `json:"boots1,omitempty"`
	Boots2     GameConfigItemsBoots2     `json:"boots2,omitempty"`
	Boots3     GameConfigItemsBoots3     `json:"boots3,omitempty"`
	Boots4     GameConfigItemsBoots4     `json:"boots4,omitempty"`
	Boots5     GameConfigItemsBoots5     `json:"boots5,omitempty"`
	Boots6     GameConfigItemsBoots6     `json:"boots6,omitempty"`
	Bread      GameConfigItemsBread      `json:"bread,omitempty"`
	Case1      GameConfigItemsCase1      `json:"case1,omitempty"`
	Case2      GameConfigItemsCase2      `json:"case2,omitempty"`
	Chest1     GameConfigItemsChest1     `json:"chest1,omitempty"`
	Chest2     GameConfigItemsChest2     `json:"chest2,omitempty"`
	Chest3     GameConfigItemsChest3     `json:"chest3,omitempty"`
	Chest4     GameConfigItemsChest4     `json:"chest4,omitempty"`
	Chest5     GameConfigItemsChest5     `json:"chest5,omitempty"`
	Chest6     GameConfigItemsChest6     `json:"chest6,omitempty"`
	Coca       GameConfigItemsCoca       `json:"coca,omitempty"`
	Cocain     GameConfigItemsCocain     `json:"cocain,omitempty"`
	Concrete   GameConfigItemsConcrete   `json:"concrete,omitempty"`
	CookedFish GameConfigItemsCookedFish `json:"cookedFish,omitempty"`
	Fish       GameConfigItemsFish       `json:"fish,omitempty"`
	Gloves1    GameConfigItemsGloves1    `json:"gloves1,omitempty"`
	Gloves2    GameConfigItemsGloves2    `json:"gloves2,omitempty"`
	Gloves3    GameConfigItemsGloves3    `json:"gloves3,omitempty"`
	Gloves4    GameConfigItemsGloves4    `json:"gloves4,omitempty"`
	Gloves5    GameConfigItemsGloves5    `json:"gloves5,omitempty"`
	Gloves6    GameConfigItemsGloves6    `json:"gloves6,omitempty"`
	Grain      GameConfigItemsGrain      `json:"grain,omitempty"`
	Gun        GameConfigItemsGun        `json:"gun,omitempty"`
	HeavyAmmo  GameConfigItemsHeavyAmmo  `json:"heavyAmmo,omitempty"`
	Helmet1    GameConfigItemsHelmet1    `json:"helmet1,omitempty"`
	Helmet2    GameConfigItemsHelmet2    `json:"helmet2,omitempty"`
	Helmet3    GameConfigItemsHelmet3    `json:"helmet3,omitempty"`
	Helmet4    GameConfigItemsHelmet4    `json:"helmet4,omitempty"`
	Helmet5    GameConfigItemsHelmet5    `json:"helmet5,omitempty"`
	Helmet6    GameConfigItemsHelmet6    `json:"helmet6,omitempty"`
	Iron       GameConfigItemsIron       `json:"iron,omitempty"`
	Jet        GameConfigItemsJet        `json:"jet,omitempty"`
	Knife      GameConfigItemsKnife      `json:"knife,omitempty"`
	Lead       GameConfigItemsLead       `json:"lead,omitempty"`
	LightAmmo  GameConfigItemsLightAmmo  `json:"lightAmmo,omitempty"`
	Limestone  GameConfigItemsLimestone  `json:"limestone,omitempty"`
	Livestock  GameConfigItemsLivestock  `json:"livestock,omitempty"`
	Oil        GameConfigItemsOil        `json:"oil,omitempty"`
	Pants1     GameConfigItemsPants1     `json:"pants1,omitempty"`
	Pants2     GameConfigItemsPants2     `json:"pants2,omitempty"`
	Pants3     GameConfigItemsPants3     `json:"pants3,omitempty"`
	Pants4     GameConfigItemsPants4     `json:"pants4,omitempty"`
	Pants5     GameConfigItemsPants5     `json:"pants5,omitempty"`
	Pants6     GameConfigItemsPants6     `json:"pants6,omitempty"`
	Paper      GameConfigItemsPaper      `json:"paper,omitempty"`
	Petroleum  GameConfigItemsPetroleum  `json:"petroleum,omitempty"`
	Rifle      GameConfigItemsRifle      `json:"rifle,omitempty"`
	Scraps     GameConfigItemsScraps     `json:"scraps,omitempty"`
	Sniper     GameConfigItemsSniper     `json:"sniper,omitempty"`
	Steak      GameConfigItemsSteak      `json:"steak,omitempty"`
	Steel      GameConfigItemsSteel      `json:"steel,omitempty"`
	Tank       GameConfigItemsTank       `json:"tank,omitempty"`
	Wood       GameConfigItemsWood       `json:"wood,omitempty"`
	WoodenCase GameConfigItemsWoodenCase `json:"woodenCase,omitempty"`
}

// GameConfigItemsAmmo is a nested component of GameConfigItems.
type GameConfigItemsAmmo struct {
	Code             string                             `json:"code,omitempty"`
	FlatStats        GameConfigItemsAmmoFlatStats       `json:"flatStats,omitempty"`
	IsTradable       bool                               `json:"isTradable,omitempty"`
	ProductionNeeds  GameConfigItemsAmmoProductionNeeds `json:"productionNeeds,omitempty"`
	ProductionPoints int                                `json:"productionPoints,omitempty"`
	Rarity           string                             `json:"rarity,omitempty"`
	SkinSlot         string                             `json:"skinSlot,omitempty"`
	Type             string                             `json:"type,omitempty"`
	Usage            string                             `json:"usage,omitempty"`
}

// GameConfigItemsAmmoFlatStats is a nested component of GameConfigItemsAmmo.
type GameConfigItemsAmmoFlatStats struct {
	PercentAttack int `json:"percentAttack,omitempty"`
}

// GameConfigItemsAmmoProductionNeeds is a nested component of GameConfigItemsAmmo.
type GameConfigItemsAmmoProductionNeeds struct {
	Lead int `json:"lead,omitempty"`
}

// GameConfigItemsBoots1 is a nested component of GameConfigItems.
type GameConfigItemsBoots1 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsBoots1DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsBoots1DynamicStats is a nested component of GameConfigItemsBoots1.
type GameConfigItemsBoots1DynamicStats struct {
	Dodge []int `json:"dodge,omitempty"`
}

// GameConfigItemsBoots2 is a nested component of GameConfigItems.
type GameConfigItemsBoots2 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsBoots2DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsBoots2DynamicStats is a nested component of GameConfigItemsBoots2.
type GameConfigItemsBoots2DynamicStats struct {
	Dodge []int `json:"dodge,omitempty"`
}

// GameConfigItemsBoots3 is a nested component of GameConfigItems.
type GameConfigItemsBoots3 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsBoots3DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsBoots3DynamicStats is a nested component of GameConfigItemsBoots3.
type GameConfigItemsBoots3DynamicStats struct {
	Dodge []int `json:"dodge,omitempty"`
}

// GameConfigItemsBoots4 is a nested component of GameConfigItems.
type GameConfigItemsBoots4 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsBoots4DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsBoots4DynamicStats is a nested component of GameConfigItemsBoots4.
type GameConfigItemsBoots4DynamicStats struct {
	Dodge []int `json:"dodge,omitempty"`
}

// GameConfigItemsBoots5 is a nested component of GameConfigItems.
type GameConfigItemsBoots5 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsBoots5DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsBoots5DynamicStats is a nested component of GameConfigItemsBoots5.
type GameConfigItemsBoots5DynamicStats struct {
	Dodge []int `json:"dodge,omitempty"`
}

// GameConfigItemsBoots6 is a nested component of GameConfigItems.
type GameConfigItemsBoots6 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsBoots6DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsBoots6DynamicStats is a nested component of GameConfigItemsBoots6.
type GameConfigItemsBoots6DynamicStats struct {
	Dodge []int `json:"dodge,omitempty"`
}

// GameConfigItemsBread is a nested component of GameConfigItems.
type GameConfigItemsBread struct {
	IconComponent    GameConfigItemsBreadIconComponent   `json:"IconComponent,omitempty"`
	Code             string                              `json:"code,omitempty"`
	FlatStats        GameConfigItemsBreadFlatStats       `json:"flatStats,omitempty"`
	IconImg          string                              `json:"iconImg,omitempty"`
	IsConsumable     bool                                `json:"isConsumable,omitempty"`
	IsTradable       bool                                `json:"isTradable,omitempty"`
	ProductionNeeds  GameConfigItemsBreadProductionNeeds `json:"productionNeeds,omitempty"`
	ProductionPoints int                                 `json:"productionPoints,omitempty"`
	Rarity           string                              `json:"rarity,omitempty"`
	Type             string                              `json:"type,omitempty"`
}

// GameConfigItemsBreadFlatStats is a nested component of GameConfigItemsBread.
type GameConfigItemsBreadFlatStats struct {
	HealthRegenPercent int `json:"healthRegenPercent,omitempty"`
}

// GameConfigItemsBreadIconComponent is a nested component of GameConfigItemsBread.
type GameConfigItemsBreadIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsBreadProductionNeeds is a nested component of GameConfigItemsBread.
type GameConfigItemsBreadProductionNeeds struct {
	Grain int `json:"grain,omitempty"`
}

// GameConfigItemsCase1 is a nested component of GameConfigItems.
type GameConfigItemsCase1 struct {
	Code       string `json:"code,omitempty"`
	IsTradable bool   `json:"isTradable,omitempty"`
	Rarity     string `json:"rarity,omitempty"`
	Type       string `json:"type,omitempty"`
	Usage      string `json:"usage,omitempty"`
}

// GameConfigItemsCase2 is a nested component of GameConfigItems.
type GameConfigItemsCase2 struct {
	Code       string `json:"code,omitempty"`
	IsTradable bool   `json:"isTradable,omitempty"`
	Rarity     string `json:"rarity,omitempty"`
	Type       string `json:"type,omitempty"`
	Usage      string `json:"usage,omitempty"`
}

// GameConfigItemsChest1 is a nested component of GameConfigItems.
type GameConfigItemsChest1 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsChest1DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsChest1DynamicStats is a nested component of GameConfigItemsChest1.
type GameConfigItemsChest1DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsChest2 is a nested component of GameConfigItems.
type GameConfigItemsChest2 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsChest2DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsChest2DynamicStats is a nested component of GameConfigItemsChest2.
type GameConfigItemsChest2DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsChest3 is a nested component of GameConfigItems.
type GameConfigItemsChest3 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsChest3DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsChest3DynamicStats is a nested component of GameConfigItemsChest3.
type GameConfigItemsChest3DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsChest4 is a nested component of GameConfigItems.
type GameConfigItemsChest4 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsChest4DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsChest4DynamicStats is a nested component of GameConfigItemsChest4.
type GameConfigItemsChest4DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsChest5 is a nested component of GameConfigItems.
type GameConfigItemsChest5 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsChest5DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsChest5DynamicStats is a nested component of GameConfigItemsChest5.
type GameConfigItemsChest5DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsChest6 is a nested component of GameConfigItems.
type GameConfigItemsChest6 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsChest6DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsChest6DynamicStats is a nested component of GameConfigItemsChest6.
type GameConfigItemsChest6DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsCoca is a nested component of GameConfigItems.
type GameConfigItemsCoca struct {
	IconComponent    GameConfigItemsCocaIconComponent `json:"IconComponent,omitempty"`
	Climates         []string                         `json:"climates,omitempty"`
	Code             string                           `json:"code,omitempty"`
	IsDeposit        bool                             `json:"isDeposit,omitempty"`
	IsTradable       bool                             `json:"isTradable,omitempty"`
	ProductionPoints int                              `json:"productionPoints,omitempty"`
	Rarity           string                           `json:"rarity,omitempty"`
	Type             string                           `json:"type,omitempty"`
}

// GameConfigItemsCocaIconComponent is a nested component of GameConfigItemsCoca.
type GameConfigItemsCocaIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsCocain is a nested component of GameConfigItems.
type GameConfigItemsCocain struct {
	IconComponent    GameConfigItemsCocainIconComponent   `json:"IconComponent,omitempty"`
	Code             string                               `json:"code,omitempty"`
	FlatStats        GameConfigItemsCocainFlatStats       `json:"flatStats,omitempty"`
	IsTradable       bool                                 `json:"isTradable,omitempty"`
	ProductionNeeds  GameConfigItemsCocainProductionNeeds `json:"productionNeeds,omitempty"`
	ProductionPoints int                                  `json:"productionPoints,omitempty"`
	Rarity           string                               `json:"rarity,omitempty"`
	Type             string                               `json:"type,omitempty"`
}

// GameConfigItemsCocainFlatStats is a nested component of GameConfigItemsCocain.
type GameConfigItemsCocainFlatStats struct {
	BuffDurationHours   int     `json:"buffDurationHours,omitempty"`
	DebuffDurationHours float64 `json:"debuffDurationHours,omitempty"`
	PercentAttack       int     `json:"percentAttack,omitempty"`
}

// GameConfigItemsCocainIconComponent is a nested component of GameConfigItemsCocain.
type GameConfigItemsCocainIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsCocainProductionNeeds is a nested component of GameConfigItemsCocain.
type GameConfigItemsCocainProductionNeeds struct {
	Coca int `json:"coca,omitempty"`
}

// GameConfigItemsConcrete is a nested component of GameConfigItems.
type GameConfigItemsConcrete struct {
	Code             string                                 `json:"code,omitempty"`
	IsTradable       bool                                   `json:"isTradable,omitempty"`
	ProductionNeeds  GameConfigItemsConcreteProductionNeeds `json:"productionNeeds,omitempty"`
	ProductionPoints int                                    `json:"productionPoints,omitempty"`
	Rarity           string                                 `json:"rarity,omitempty"`
	Type             string                                 `json:"type,omitempty"`
}

// GameConfigItemsConcreteProductionNeeds is a nested component of GameConfigItemsConcrete.
type GameConfigItemsConcreteProductionNeeds struct {
	Limestone int `json:"limestone,omitempty"`
}

// GameConfigItemsCookedFish is a nested component of GameConfigItems.
type GameConfigItemsCookedFish struct {
	Code             string                                   `json:"code,omitempty"`
	FlatStats        GameConfigItemsCookedFishFlatStats       `json:"flatStats,omitempty"`
	IsConsumable     bool                                     `json:"isConsumable,omitempty"`
	IsTradable       bool                                     `json:"isTradable,omitempty"`
	ProductionNeeds  GameConfigItemsCookedFishProductionNeeds `json:"productionNeeds,omitempty"`
	ProductionPoints int                                      `json:"productionPoints,omitempty"`
	Rarity           string                                   `json:"rarity,omitempty"`
	Type             string                                   `json:"type,omitempty"`
}

// GameConfigItemsCookedFishFlatStats is a nested component of GameConfigItemsCookedFish.
type GameConfigItemsCookedFishFlatStats struct {
	HealthRegenPercent int `json:"healthRegenPercent,omitempty"`
}

// GameConfigItemsCookedFishProductionNeeds is a nested component of GameConfigItemsCookedFish.
type GameConfigItemsCookedFishProductionNeeds struct {
	Fish int `json:"fish,omitempty"`
}

// GameConfigItemsFish is a nested component of GameConfigItems.
type GameConfigItemsFish struct {
	IconComponent    GameConfigItemsFishIconComponent `json:"IconComponent,omitempty"`
	Climates         []string                         `json:"climates,omitempty"`
	Code             string                           `json:"code,omitempty"`
	IsDeposit        bool                             `json:"isDeposit,omitempty"`
	IsTradable       bool                             `json:"isTradable,omitempty"`
	ProductionPoints int                              `json:"productionPoints,omitempty"`
	Rarity           string                           `json:"rarity,omitempty"`
	Type             string                           `json:"type,omitempty"`
}

// GameConfigItemsFishIconComponent is a nested component of GameConfigItemsFish.
type GameConfigItemsFishIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsGloves1 is a nested component of GameConfigItems.
type GameConfigItemsGloves1 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsGloves1DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsGloves1DynamicStats is a nested component of GameConfigItemsGloves1.
type GameConfigItemsGloves1DynamicStats struct {
	Precision []int `json:"precision,omitempty"`
}

// GameConfigItemsGloves2 is a nested component of GameConfigItems.
type GameConfigItemsGloves2 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsGloves2DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsGloves2DynamicStats is a nested component of GameConfigItemsGloves2.
type GameConfigItemsGloves2DynamicStats struct {
	Precision []int `json:"precision,omitempty"`
}

// GameConfigItemsGloves3 is a nested component of GameConfigItems.
type GameConfigItemsGloves3 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsGloves3DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsGloves3DynamicStats is a nested component of GameConfigItemsGloves3.
type GameConfigItemsGloves3DynamicStats struct {
	Precision []int `json:"precision,omitempty"`
}

// GameConfigItemsGloves4 is a nested component of GameConfigItems.
type GameConfigItemsGloves4 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsGloves4DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsGloves4DynamicStats is a nested component of GameConfigItemsGloves4.
type GameConfigItemsGloves4DynamicStats struct {
	Precision []int `json:"precision,omitempty"`
}

// GameConfigItemsGloves5 is a nested component of GameConfigItems.
type GameConfigItemsGloves5 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsGloves5DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsGloves5DynamicStats is a nested component of GameConfigItemsGloves5.
type GameConfigItemsGloves5DynamicStats struct {
	Precision []int `json:"precision,omitempty"`
}

// GameConfigItemsGloves6 is a nested component of GameConfigItems.
type GameConfigItemsGloves6 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsGloves6DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsGloves6DynamicStats is a nested component of GameConfigItemsGloves6.
type GameConfigItemsGloves6DynamicStats struct {
	Precision []int `json:"precision,omitempty"`
}

// GameConfigItemsGrain is a nested component of GameConfigItems.
type GameConfigItemsGrain struct {
	IconComponent    GameConfigItemsGrainIconComponent `json:"IconComponent,omitempty"`
	Climates         []string                          `json:"climates,omitempty"`
	Code             string                            `json:"code,omitempty"`
	IsDeposit        bool                              `json:"isDeposit,omitempty"`
	IsTradable       bool                              `json:"isTradable,omitempty"`
	ProductionPoints int                               `json:"productionPoints,omitempty"`
	Rarity           string                            `json:"rarity,omitempty"`
	Type             string                            `json:"type,omitempty"`
}

// GameConfigItemsGrainIconComponent is a nested component of GameConfigItemsGrain.
type GameConfigItemsGrainIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsGun is a nested component of GameConfigItems.
type GameConfigItemsGun struct {
	IconComponent GameConfigItemsGunIconComponent `json:"IconComponent,omitempty"`
	Code          string                          `json:"code,omitempty"`
	DynamicStats  GameConfigItemsGunDynamicStats  `json:"dynamicStats,omitempty"`
	Rarity        string                          `json:"rarity,omitempty"`
	SkinSlot      string                          `json:"skinSlot,omitempty"`
	Type          string                          `json:"type,omitempty"`
	Usage         string                          `json:"usage,omitempty"`
}

// GameConfigItemsGunDynamicStats is a nested component of GameConfigItemsGun.
type GameConfigItemsGunDynamicStats struct {
	Attack         []int `json:"attack,omitempty"`
	CriticalChance []int `json:"criticalChance,omitempty"`
}

// GameConfigItemsGunIconComponent is a nested component of GameConfigItemsGun.
type GameConfigItemsGunIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsHeavyAmmo is a nested component of GameConfigItems.
type GameConfigItemsHeavyAmmo struct {
	IconComponent    GameConfigItemsHeavyAmmoIconComponent   `json:"IconComponent,omitempty"`
	Code             string                                  `json:"code,omitempty"`
	FlatStats        GameConfigItemsHeavyAmmoFlatStats       `json:"flatStats,omitempty"`
	IsTradable       bool                                    `json:"isTradable,omitempty"`
	ProductionNeeds  GameConfigItemsHeavyAmmoProductionNeeds `json:"productionNeeds,omitempty"`
	ProductionPoints int                                     `json:"productionPoints,omitempty"`
	Rarity           string                                  `json:"rarity,omitempty"`
	SkinSlot         string                                  `json:"skinSlot,omitempty"`
	Type             string                                  `json:"type,omitempty"`
	Usage            string                                  `json:"usage,omitempty"`
}

// GameConfigItemsHeavyAmmoFlatStats is a nested component of GameConfigItemsHeavyAmmo.
type GameConfigItemsHeavyAmmoFlatStats struct {
	PercentAttack int `json:"percentAttack,omitempty"`
}

// GameConfigItemsHeavyAmmoIconComponent is a nested component of GameConfigItemsHeavyAmmo.
type GameConfigItemsHeavyAmmoIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsHeavyAmmoProductionNeeds is a nested component of GameConfigItemsHeavyAmmo.
type GameConfigItemsHeavyAmmoProductionNeeds struct {
	Lead int `json:"lead,omitempty"`
}

// GameConfigItemsHelmet1 is a nested component of GameConfigItems.
type GameConfigItemsHelmet1 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsHelmet1DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsHelmet1DynamicStats is a nested component of GameConfigItemsHelmet1.
type GameConfigItemsHelmet1DynamicStats struct {
	CriticalDamages []int `json:"criticalDamages,omitempty"`
}

// GameConfigItemsHelmet2 is a nested component of GameConfigItems.
type GameConfigItemsHelmet2 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsHelmet2DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsHelmet2DynamicStats is a nested component of GameConfigItemsHelmet2.
type GameConfigItemsHelmet2DynamicStats struct {
	CriticalDamages []int `json:"criticalDamages,omitempty"`
}

// GameConfigItemsHelmet3 is a nested component of GameConfigItems.
type GameConfigItemsHelmet3 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsHelmet3DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsHelmet3DynamicStats is a nested component of GameConfigItemsHelmet3.
type GameConfigItemsHelmet3DynamicStats struct {
	CriticalDamages []int `json:"criticalDamages,omitempty"`
}

// GameConfigItemsHelmet4 is a nested component of GameConfigItems.
type GameConfigItemsHelmet4 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsHelmet4DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsHelmet4DynamicStats is a nested component of GameConfigItemsHelmet4.
type GameConfigItemsHelmet4DynamicStats struct {
	CriticalDamages []int `json:"criticalDamages,omitempty"`
}

// GameConfigItemsHelmet5 is a nested component of GameConfigItems.
type GameConfigItemsHelmet5 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsHelmet5DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsHelmet5DynamicStats is a nested component of GameConfigItemsHelmet5.
type GameConfigItemsHelmet5DynamicStats struct {
	CriticalDamages []int `json:"criticalDamages,omitempty"`
}

// GameConfigItemsHelmet6 is a nested component of GameConfigItems.
type GameConfigItemsHelmet6 struct {
	Code         string                             `json:"code,omitempty"`
	DynamicStats GameConfigItemsHelmet6DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                             `json:"iconImg,omitempty"`
	Rarity       string                             `json:"rarity,omitempty"`
	SkinSlot     string                             `json:"skinSlot,omitempty"`
	Type         string                             `json:"type,omitempty"`
	Usage        string                             `json:"usage,omitempty"`
}

// GameConfigItemsHelmet6DynamicStats is a nested component of GameConfigItemsHelmet6.
type GameConfigItemsHelmet6DynamicStats struct {
	CriticalDamages []int `json:"criticalDamages,omitempty"`
}

// GameConfigItemsIron is a nested component of GameConfigItems.
type GameConfigItemsIron struct {
	Climates         []string `json:"climates,omitempty"`
	Code             string   `json:"code,omitempty"`
	IsDeposit        bool     `json:"isDeposit,omitempty"`
	IsTradable       bool     `json:"isTradable,omitempty"`
	ProductionPoints int      `json:"productionPoints,omitempty"`
	Rarity           string   `json:"rarity,omitempty"`
	Type             string   `json:"type,omitempty"`
}

// GameConfigItemsJet is a nested component of GameConfigItems.
type GameConfigItemsJet struct {
	Code         string                         `json:"code,omitempty"`
	DynamicStats GameConfigItemsJetDynamicStats `json:"dynamicStats,omitempty"`
	Rarity       string                         `json:"rarity,omitempty"`
	SkinSlot     string                         `json:"skinSlot,omitempty"`
	Type         string                         `json:"type,omitempty"`
	Usage        string                         `json:"usage,omitempty"`
}

// GameConfigItemsJetDynamicStats is a nested component of GameConfigItemsJet.
type GameConfigItemsJetDynamicStats struct {
	Attack         []int `json:"attack,omitempty"`
	CriticalChance []int `json:"criticalChance,omitempty"`
}

// GameConfigItemsKnife is a nested component of GameConfigItems.
type GameConfigItemsKnife struct {
	IconComponent GameConfigItemsKnifeIconComponent `json:"IconComponent,omitempty"`
	Code          string                            `json:"code,omitempty"`
	DynamicStats  GameConfigItemsKnifeDynamicStats  `json:"dynamicStats,omitempty"`
	Rarity        string                            `json:"rarity,omitempty"`
	SkinSlot      string                            `json:"skinSlot,omitempty"`
	Type          string                            `json:"type,omitempty"`
	Usage         string                            `json:"usage,omitempty"`
}

// GameConfigItemsKnifeDynamicStats is a nested component of GameConfigItemsKnife.
type GameConfigItemsKnifeDynamicStats struct {
	Attack         []int `json:"attack,omitempty"`
	CriticalChance []int `json:"criticalChance,omitempty"`
}

// GameConfigItemsKnifeIconComponent is a nested component of GameConfigItemsKnife.
type GameConfigItemsKnifeIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsLead is a nested component of GameConfigItems.
type GameConfigItemsLead struct {
	Climates         []string `json:"climates,omitempty"`
	Code             string   `json:"code,omitempty"`
	IsDeposit        bool     `json:"isDeposit,omitempty"`
	IsTradable       bool     `json:"isTradable,omitempty"`
	ProductionPoints int      `json:"productionPoints,omitempty"`
	Rarity           string   `json:"rarity,omitempty"`
	Type             string   `json:"type,omitempty"`
}

// GameConfigItemsLightAmmo is a nested component of GameConfigItems.
type GameConfigItemsLightAmmo struct {
	Code             string                                  `json:"code,omitempty"`
	FlatStats        GameConfigItemsLightAmmoFlatStats       `json:"flatStats,omitempty"`
	IsTradable       bool                                    `json:"isTradable,omitempty"`
	ProductionNeeds  GameConfigItemsLightAmmoProductionNeeds `json:"productionNeeds,omitempty"`
	ProductionPoints int                                     `json:"productionPoints,omitempty"`
	Rarity           string                                  `json:"rarity,omitempty"`
	SkinSlot         string                                  `json:"skinSlot,omitempty"`
	Type             string                                  `json:"type,omitempty"`
	Usage            string                                  `json:"usage,omitempty"`
}

// GameConfigItemsLightAmmoFlatStats is a nested component of GameConfigItemsLightAmmo.
type GameConfigItemsLightAmmoFlatStats struct {
	PercentAttack int `json:"percentAttack,omitempty"`
}

// GameConfigItemsLightAmmoProductionNeeds is a nested component of GameConfigItemsLightAmmo.
type GameConfigItemsLightAmmoProductionNeeds struct {
	Lead int `json:"lead,omitempty"`
}

// GameConfigItemsLimestone is a nested component of GameConfigItems.
type GameConfigItemsLimestone struct {
	Climates         []string `json:"climates,omitempty"`
	Code             string   `json:"code,omitempty"`
	IsDeposit        bool     `json:"isDeposit,omitempty"`
	IsTradable       bool     `json:"isTradable,omitempty"`
	ProductionPoints int      `json:"productionPoints,omitempty"`
	Rarity           string   `json:"rarity,omitempty"`
	Type             string   `json:"type,omitempty"`
}

// GameConfigItemsLivestock is a nested component of GameConfigItems.
type GameConfigItemsLivestock struct {
	IconComponent    GameConfigItemsLivestockIconComponent `json:"IconComponent,omitempty"`
	Climates         []string                              `json:"climates,omitempty"`
	Code             string                                `json:"code,omitempty"`
	IsDeposit        bool                                  `json:"isDeposit,omitempty"`
	IsTradable       bool                                  `json:"isTradable,omitempty"`
	ProductionPoints int                                   `json:"productionPoints,omitempty"`
	Rarity           string                                `json:"rarity,omitempty"`
	Type             string                                `json:"type,omitempty"`
}

// GameConfigItemsLivestockIconComponent is a nested component of GameConfigItemsLivestock.
type GameConfigItemsLivestockIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsOil is a nested component of GameConfigItems.
type GameConfigItemsOil struct {
	IconComponent    GameConfigItemsOilIconComponent   `json:"IconComponent,omitempty"`
	Code             string                            `json:"code,omitempty"`
	IsTradable       bool                              `json:"isTradable,omitempty"`
	ProductionNeeds  GameConfigItemsOilProductionNeeds `json:"productionNeeds,omitempty"`
	ProductionPoints int                               `json:"productionPoints,omitempty"`
	Rarity           string                            `json:"rarity,omitempty"`
	Type             string                            `json:"type,omitempty"`
}

// GameConfigItemsOilIconComponent is a nested component of GameConfigItemsOil.
type GameConfigItemsOilIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsOilProductionNeeds is a nested component of GameConfigItemsOil.
type GameConfigItemsOilProductionNeeds struct {
	Petroleum int `json:"petroleum,omitempty"`
}

// GameConfigItemsPants1 is a nested component of GameConfigItems.
type GameConfigItemsPants1 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsPants1DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsPants1DynamicStats is a nested component of GameConfigItemsPants1.
type GameConfigItemsPants1DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsPants2 is a nested component of GameConfigItems.
type GameConfigItemsPants2 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsPants2DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsPants2DynamicStats is a nested component of GameConfigItemsPants2.
type GameConfigItemsPants2DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsPants3 is a nested component of GameConfigItems.
type GameConfigItemsPants3 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsPants3DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsPants3DynamicStats is a nested component of GameConfigItemsPants3.
type GameConfigItemsPants3DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsPants4 is a nested component of GameConfigItems.
type GameConfigItemsPants4 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsPants4DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsPants4DynamicStats is a nested component of GameConfigItemsPants4.
type GameConfigItemsPants4DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsPants5 is a nested component of GameConfigItems.
type GameConfigItemsPants5 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsPants5DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsPants5DynamicStats is a nested component of GameConfigItemsPants5.
type GameConfigItemsPants5DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsPants6 is a nested component of GameConfigItems.
type GameConfigItemsPants6 struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsPants6DynamicStats `json:"dynamicStats,omitempty"`
	IconImg      string                            `json:"iconImg,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsPants6DynamicStats is a nested component of GameConfigItemsPants6.
type GameConfigItemsPants6DynamicStats struct {
	Armor []int `json:"armor,omitempty"`
}

// GameConfigItemsPaper is a nested component of GameConfigItems.
type GameConfigItemsPaper struct {
	IconComponent    GameConfigItemsPaperIconComponent   `json:"IconComponent,omitempty"`
	Code             string                              `json:"code,omitempty"`
	IsTradable       bool                                `json:"isTradable,omitempty"`
	ProductionNeeds  GameConfigItemsPaperProductionNeeds `json:"productionNeeds,omitempty"`
	ProductionPoints int                                 `json:"productionPoints,omitempty"`
	Rarity           string                              `json:"rarity,omitempty"`
	Type             string                              `json:"type,omitempty"`
}

// GameConfigItemsPaperIconComponent is a nested component of GameConfigItemsPaper.
type GameConfigItemsPaperIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsPaperProductionNeeds is a nested component of GameConfigItemsPaper.
type GameConfigItemsPaperProductionNeeds struct {
	Wood int `json:"wood,omitempty"`
}

// GameConfigItemsPetroleum is a nested component of GameConfigItems.
type GameConfigItemsPetroleum struct {
	IconComponent    GameConfigItemsPetroleumIconComponent `json:"IconComponent,omitempty"`
	Climates         []string                              `json:"climates,omitempty"`
	Code             string                                `json:"code,omitempty"`
	IsDeposit        bool                                  `json:"isDeposit,omitempty"`
	IsTradable       bool                                  `json:"isTradable,omitempty"`
	ProductionPoints int                                   `json:"productionPoints,omitempty"`
	Rarity           string                                `json:"rarity,omitempty"`
	Type             string                                `json:"type,omitempty"`
}

// GameConfigItemsPetroleumIconComponent is a nested component of GameConfigItemsPetroleum.
type GameConfigItemsPetroleumIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsRifle is a nested component of GameConfigItems.
type GameConfigItemsRifle struct {
	Code         string                           `json:"code,omitempty"`
	DynamicStats GameConfigItemsRifleDynamicStats `json:"dynamicStats,omitempty"`
	Rarity       string                           `json:"rarity,omitempty"`
	SkinSlot     string                           `json:"skinSlot,omitempty"`
	Type         string                           `json:"type,omitempty"`
	Usage        string                           `json:"usage,omitempty"`
}

// GameConfigItemsRifleDynamicStats is a nested component of GameConfigItemsRifle.
type GameConfigItemsRifleDynamicStats struct {
	Attack         []int `json:"attack,omitempty"`
	CriticalChance []int `json:"criticalChance,omitempty"`
}

// GameConfigItemsScraps is a nested component of GameConfigItems.
type GameConfigItemsScraps struct {
	Code       string `json:"code,omitempty"`
	IsTradable bool   `json:"isTradable,omitempty"`
	Rarity     string `json:"rarity,omitempty"`
	Type       string `json:"type,omitempty"`
}

// GameConfigItemsSniper is a nested component of GameConfigItems.
type GameConfigItemsSniper struct {
	Code         string                            `json:"code,omitempty"`
	DynamicStats GameConfigItemsSniperDynamicStats `json:"dynamicStats,omitempty"`
	Rarity       string                            `json:"rarity,omitempty"`
	SkinSlot     string                            `json:"skinSlot,omitempty"`
	Type         string                            `json:"type,omitempty"`
	Usage        string                            `json:"usage,omitempty"`
}

// GameConfigItemsSniperDynamicStats is a nested component of GameConfigItemsSniper.
type GameConfigItemsSniperDynamicStats struct {
	Attack         []int `json:"attack,omitempty"`
	CriticalChance []int `json:"criticalChance,omitempty"`
}

// GameConfigItemsSteak is a nested component of GameConfigItems.
type GameConfigItemsSteak struct {
	IconComponent    GameConfigItemsSteakIconComponent   `json:"IconComponent,omitempty"`
	Code             string                              `json:"code,omitempty"`
	FlatStats        GameConfigItemsSteakFlatStats       `json:"flatStats,omitempty"`
	IsConsumable     bool                                `json:"isConsumable,omitempty"`
	IsTradable       bool                                `json:"isTradable,omitempty"`
	ProductionNeeds  GameConfigItemsSteakProductionNeeds `json:"productionNeeds,omitempty"`
	ProductionPoints int                                 `json:"productionPoints,omitempty"`
	Rarity           string                              `json:"rarity,omitempty"`
	Type             string                              `json:"type,omitempty"`
}

// GameConfigItemsSteakFlatStats is a nested component of GameConfigItemsSteak.
type GameConfigItemsSteakFlatStats struct {
	HealthRegenPercent int `json:"healthRegenPercent,omitempty"`
}

// GameConfigItemsSteakIconComponent is a nested component of GameConfigItemsSteak.
type GameConfigItemsSteakIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsSteakProductionNeeds is a nested component of GameConfigItemsSteak.
type GameConfigItemsSteakProductionNeeds struct {
	Livestock int `json:"livestock,omitempty"`
}

// GameConfigItemsSteel is a nested component of GameConfigItems.
type GameConfigItemsSteel struct {
	Code             string                              `json:"code,omitempty"`
	IsTradable       bool                                `json:"isTradable,omitempty"`
	ProductionNeeds  GameConfigItemsSteelProductionNeeds `json:"productionNeeds,omitempty"`
	ProductionPoints int                                 `json:"productionPoints,omitempty"`
	Rarity           string                              `json:"rarity,omitempty"`
	Type             string                              `json:"type,omitempty"`
}

// GameConfigItemsSteelProductionNeeds is a nested component of GameConfigItemsSteel.
type GameConfigItemsSteelProductionNeeds struct {
	Iron int `json:"iron,omitempty"`
}

// GameConfigItemsTank is a nested component of GameConfigItems.
type GameConfigItemsTank struct {
	IconComponent GameConfigItemsTankIconComponent `json:"IconComponent,omitempty"`
	Code          string                           `json:"code,omitempty"`
	DynamicStats  GameConfigItemsTankDynamicStats  `json:"dynamicStats,omitempty"`
	Rarity        string                           `json:"rarity,omitempty"`
	SkinSlot      string                           `json:"skinSlot,omitempty"`
	Type          string                           `json:"type,omitempty"`
	Usage         string                           `json:"usage,omitempty"`
}

// GameConfigItemsTankDynamicStats is a nested component of GameConfigItemsTank.
type GameConfigItemsTankDynamicStats struct {
	Attack         []int `json:"attack,omitempty"`
	CriticalChance []int `json:"criticalChance,omitempty"`
}

// GameConfigItemsTankIconComponent is a nested component of GameConfigItemsTank.
type GameConfigItemsTankIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsWood is a nested component of GameConfigItems.
type GameConfigItemsWood struct {
	IconComponent    GameConfigItemsWoodIconComponent `json:"IconComponent,omitempty"`
	Climates         []string                         `json:"climates,omitempty"`
	Code             string                           `json:"code,omitempty"`
	IsDeposit        bool                             `json:"isDeposit,omitempty"`
	IsTradable       bool                             `json:"isTradable,omitempty"`
	ProductionPoints int                              `json:"productionPoints,omitempty"`
	Rarity           string                           `json:"rarity,omitempty"`
	Type             string                           `json:"type,omitempty"`
}

// GameConfigItemsWoodIconComponent is a nested component of GameConfigItemsWood.
type GameConfigItemsWoodIconComponent struct {
	Compare json.RawMessage `json:"compare,omitempty"`
}

// GameConfigItemsWoodenCase is a nested component of GameConfigItems.
type GameConfigItemsWoodenCase struct {
	Code       string `json:"code,omitempty"`
	IsTradable bool   `json:"isTradable,omitempty"`
	Rarity     string `json:"rarity,omitempty"`
	Type       string `json:"type,omitempty"`
	Usage      string `json:"usage,omitempty"`
}

// GameConfigLaw is a nested component of GameConfig.
type GameConfigLaw struct {
	AbusiveLawPossibleVotersNeeded int                               `json:"abusiveLawPossibleVotersNeeded,omitempty"`
	AbusiveLawsCooldownInDays      int                               `json:"abusiveLawsCooldownInDays,omitempty"`
	Accept_join_alliance           GameConfigLawAccept_join_alliance `json:"accept_join_alliance,omitempty"`
	Create_alliance                GameConfigLawCreate_alliance      `json:"create_alliance,omitempty"`
	Define_enemy_country           GameConfigLawDefine_enemy_country `json:"define_enemy_country,omitempty"`
	LawVotesDurationHours          int                               `json:"lawVotesDurationHours,omitempty"`
	Leave_alliance                 GameConfigLawLeave_alliance       `json:"leave_alliance,omitempty"`
	SendMoneyToCountry             GameConfigLawSendMoneyToCountry   `json:"sendMoneyToCountry,omitempty"`
	Set_color_scheme               GameConfigLawSet_color_scheme     `json:"set_color_scheme,omitempty"`
	VotersRatioNeeded              float64                           `json:"votersRatioNeeded,omitempty"`
}

// GameConfigLawAccept_join_alliance is a nested component of GameConfigLaw.
type GameConfigLawAccept_join_alliance struct {
	Cost int `json:"cost,omitempty"`
}

// GameConfigLawCreate_alliance is a nested component of GameConfigLaw.
type GameConfigLawCreate_alliance struct {
	Cost            int `json:"cost,omitempty"`
	MaintenanceCost int `json:"maintenanceCost,omitempty"`
}

// GameConfigLawDefine_enemy_country is a nested component of GameConfigLaw.
type GameConfigLawDefine_enemy_country struct {
	Cost            int `json:"cost,omitempty"`
	MaintenanceCost int `json:"maintenanceCost,omitempty"`
}

// GameConfigLawLeave_alliance is a nested component of GameConfigLaw.
type GameConfigLawLeave_alliance struct {
	Cost int `json:"cost,omitempty"`
}

// GameConfigLawSendMoneyToCountry is a nested component of GameConfigLaw.
type GameConfigLawSendMoneyToCountry struct {
	AllianceTaxRate float64 `json:"allianceTaxRate,omitempty"`
	ExternalTaxRate int     `json:"externalTaxRate,omitempty"`
}

// GameConfigLawSet_color_scheme is a nested component of GameConfigLaw.
type GameConfigLawSet_color_scheme struct {
	Cost int `json:"cost,omitempty"`
}

// GameConfigLoot is a nested component of GameConfig.
type GameConfigLoot struct {
	BattleLootDamagePerLootItem int `json:"battleLootDamagePerLootItem,omitempty"`
	DamagePerLootItem           int `json:"damagePerLootItem,omitempty"`
	WeaponChancePercent         int `json:"weaponChancePercent,omitempty"`
}

// GameConfigMercenaryContract is a nested component of GameConfig.
type GameConfigMercenaryContract struct {
	AcceptanceFeePercent            int                                   `json:"acceptanceFeePercent,omitempty"`
	Auction                         GameConfigMercenaryContractAuction    `json:"auction,omitempty"`
	BountyCooldownMinutes           int                                   `json:"bountyCooldownMinutes,omitempty"`
	BountyMinActiveCitizens         int                                   `json:"bountyMinActiveCitizens,omitempty"`
	CancellationPenaltyPercent      int                                   `json:"cancellationPenaltyPercent,omitempty"`
	Enabled                         bool                                  `json:"enabled,omitempty"`
	MaxBountyWageMultiplier         int                                   `json:"maxBountyWageMultiplier,omitempty"`
	NationalBountyEnabled           bool                                  `json:"nationalBountyEnabled,omitempty"`
	NationalBountyMinActiveCitizens int                                   `json:"nationalBountyMinActiveCitizens,omitempty"`
	Reputation                      GameConfigMercenaryContractReputation `json:"reputation,omitempty"`
}

// GameConfigMercenaryContractAuction is a nested component of GameConfigMercenaryContract.
type GameConfigMercenaryContractAuction struct {
	BidStep                 float64 `json:"bidStep,omitempty"`
	MaxDamagePerContract    int     `json:"maxDamagePerContract,omitempty"`
	MaxDuration             int     `json:"maxDuration,omitempty"`
	MaxPerKMultiplier       int     `json:"maxPerKMultiplier,omitempty"`
	MinDuration             int     `json:"minDuration,omitempty"`
	MinPerK                 float64 `json:"minPerK,omitempty"`
	TimerExtensionAmount    int     `json:"timerExtensionAmount,omitempty"`
	TimerExtensionThreshold int     `json:"timerExtensionThreshold,omitempty"`
}

// GameConfigMercenaryContractReputation is a nested component of GameConfigMercenaryContract.
type GameConfigMercenaryContractReputation struct {
	FailurePerDamage           int `json:"failurePerDamage,omitempty"`
	FailurePerDollar           int `json:"failurePerDollar,omitempty"`
	NegativeBuyAmount          int `json:"negativeBuyAmount,omitempty"`
	NegativeBuyCooldownHours   int `json:"negativeBuyCooldownHours,omitempty"`
	NegativeBuyCost            int `json:"negativeBuyCost,omitempty"`
	ProfessionalsOnlyThreshold int `json:"professionalsOnlyThreshold,omitempty"`
	SuccessPerDollar           int `json:"successPerDollar,omitempty"`
	WeeklyDecayPercent         int `json:"weeklyDecayPercent,omitempty"`
}

// GameConfigMergingCost is a nested component of GameConfig.
type GameConfigMergingCost struct {
	Common    int `json:"common,omitempty"`
	Epic      int `json:"epic,omitempty"`
	Legendary int `json:"legendary,omitempty"`
	Mythic    int `json:"mythic,omitempty"`
	Rare      int `json:"rare,omitempty"`
	Uncommon  int `json:"uncommon,omitempty"`
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
	Cases          int `json:"cases,omitempty"`
	Money          int `json:"money,omitempty"`
	XP             int `json:"xp,omitempty"`
	XpWhenFinished int `json:"xpWhenFinished,omitempty"`
}

// GameConfigMissionRewardMonthly is a nested component of GameConfigMissionReward.
type GameConfigMissionRewardMonthly struct {
	Cases          int `json:"cases,omitempty"`
	Money          int `json:"money,omitempty"`
	XP             int `json:"xp,omitempty"`
	XpWhenFinished int `json:"xpWhenFinished,omitempty"`
}

// GameConfigMissionRewardStarting is a nested component of GameConfigMissionReward.
type GameConfigMissionRewardStarting struct {
	Cases          int `json:"cases,omitempty"`
	Money          int `json:"money,omitempty"`
	XP             int `json:"xp,omitempty"`
	XpWhenFinished int `json:"xpWhenFinished,omitempty"`
}

// GameConfigMissionRewardWeekly is a nested component of GameConfigMissionReward.
type GameConfigMissionRewardWeekly struct {
	Cases          int `json:"cases,omitempty"`
	Money          int `json:"money,omitempty"`
	XP             int `json:"xp,omitempty"`
	XpWhenFinished int `json:"xpWhenFinished,omitempty"`
}

// GameConfigMu is a nested component of GameConfig.
type GameConfigMu struct {
	ConstructionCost              int `json:"constructionCost,omitempty"`
	DestructionValuePercent       int `json:"destructionValuePercent,omitempty"`
	HealthPerHelp                 int `json:"healthPerHelp,omitempty"`
	HelpCooldownHours             int `json:"helpCooldownHours,omitempty"`
	HelpValue                     int `json:"helpValue,omitempty"`
	MaxOwnedMus                   int `json:"maxOwnedMus,omitempty"`
	MoveCost                      int `json:"moveCost,omitempty"`
	NationalityChangeCooldownDays int `json:"nationalityChangeCooldownDays,omitempty"`
}

// GameConfigNewspaper is a nested component of GameConfig.
type GameConfigNewspaper struct {
	CommentMinLevel       int `json:"commentMinLevel,omitempty"`
	CreateArticleMinLevel int `json:"createArticleMinLevel,omitempty"`
	GemTipValue           int `json:"gemTipValue,omitempty"`
	PublishCost           int `json:"publishCost,omitempty"`
	TipMinLevel           int `json:"tipMinLevel,omitempty"`
	TipValue              int `json:"tipValue,omitempty"`
}

// GameConfigOrg is a nested component of GameConfig.
type GameConfigOrg struct {
	ConstructionCost int `json:"constructionCost,omitempty"`
	MoveCost         int `json:"moveCost,omitempty"`
}

// GameConfigParty is a nested component of GameConfig.
type GameConfigParty struct {
	CreateCost int `json:"createCost,omitempty"`
}

// GameConfigReferral is a nested component of GameConfig.
type GameConfigReferral struct {
	CanSetReferrerBeforeOrAtLevel  int `json:"canSetReferrerBeforeOrAtLevel,omitempty"`
	LevelNeededForBadge            int `json:"levelNeededForBadge,omitempty"`
	LifeTimeBadgeMoneySharePercent int `json:"lifeTimeBadgeMoneySharePercent,omitempty"`
	MoneyForBeingReferred          int `json:"moneyForBeingReferred,omitempty"`
}

// GameConfigRegion is a nested component of GameConfig.
type GameConfigRegion struct {
	BattleCooldownHours                             int                `json:"battleCooldownHours,omitempty"`
	CitizenIncomeTaxPercent                         int                `json:"citizenIncomeTaxPercent,omitempty"`
	DecreaseBy                                      int                `json:"decreaseBy,omitempty"`
	DecreaseResistanceCost                          int                `json:"decreaseResistanceCost,omitempty"`
	DepleteHourlyPercent                            float64            `json:"depleteHourlyPercent,omitempty"`
	FreeTravelHomeCooldownHours                     int                `json:"freeTravelHomeCooldownHours,omitempty"`
	HomeMoveConcreteCost                            int                `json:"homeMoveConcreteCost,omitempty"`
	IncreaseBy                                      int                `json:"increaseBy,omitempty"`
	IncreaseResistanceCost                          int                `json:"increaseResistanceCost,omitempty"`
	LiberationDaysCooldown                          int                `json:"liberationDaysCooldown,omitempty"`
	MaxDailyResistance                              int                `json:"maxDailyResistance,omitempty"`
	MaxResistance                                   int                `json:"maxResistance,omitempty"`
	MinDailyResistance                              int                `json:"minDailyResistance,omitempty"`
	NonAggressionHoursAfterLiberation               int                `json:"nonAggressionHoursAfterLiberation,omitempty"`
	NonAggressionHoursAfterPeace                    int                `json:"nonAggressionHoursAfterPeace,omitempty"`
	RegionNotLinkedToCapitalMalusDevelopmentPercent int                `json:"regionNotLinkedToCapitalMalusDevelopmentPercent,omitempty"`
	ResistanceAllyBonusPercent                      int                `json:"resistanceAllyBonusPercent,omitempty"`
	ResistanceBarMultiplier                         int                `json:"resistanceBarMultiplier,omitempty"`
	ResistanceBattleCooldownHours                   int                `json:"resistanceBattleCooldownHours,omitempty"`
	ResistanceBattleStartCostMultiplier             int                `json:"resistanceBattleStartCostMultiplier,omitempty"`
	ResistanceCitizenBonus                          int                `json:"resistanceCitizenBonus,omitempty"`
	ResistanceCitizenBonusPercent                   int                `json:"resistanceCitizenBonusPercent,omitempty"`
	ResistanceContributionCooldownAfterRevoltHours  int                `json:"resistanceContributionCooldownAfterRevoltHours,omitempty"`
	ResistanceContributionCost                      int                `json:"resistanceContributionCost,omitempty"`
	ResistanceContributionMinLevel                  int                `json:"resistanceContributionMinLevel,omitempty"`
	ResistanceContributionValue                     int                `json:"resistanceContributionValue,omitempty"`
	ResistanceDecayPercent                          float64            `json:"resistanceDecayPercent,omitempty"`
	ResistanceForeignGovCostMultiplier              int                `json:"resistanceForeignGovCostMultiplier,omitempty"`
	ResistanceGrowthPercentMax                      float64            `json:"resistanceGrowthPercentMax,omitempty"`
	ResistanceGrowthPercentMin                      float64            `json:"resistanceGrowthPercentMin,omitempty"`
	ResistancePassiveGrowthPercent                  float64            `json:"resistancePassiveGrowthPercent,omitempty"`
	ResourcesBonus                                  map[string]float64 `json:"resourcesBonus,omitempty"`
	TransferDaysCooldown                            int                `json:"transferDaysCooldown,omitempty"`
}

// GameConfigSkills is a nested component of GameConfig.
type GameConfigSkills struct {
	Armor            GameConfigSkillsArmor            `json:"armor,omitempty"`
	Attack           GameConfigSkillsAttack           `json:"attack,omitempty"`
	Companies        GameConfigSkillsCompanies        `json:"companies,omitempty"`
	CriticalChance   GameConfigSkillsCriticalChance   `json:"criticalChance,omitempty"`
	CriticalDamages  GameConfigSkillsCriticalDamages  `json:"criticalDamages,omitempty"`
	Dodge            GameConfigSkillsDodge            `json:"dodge,omitempty"`
	Energy           GameConfigSkillsEnergy           `json:"energy,omitempty"`
	Entrepreneurship GameConfigSkillsEntrepreneurship `json:"entrepreneurship,omitempty"`
	Health           GameConfigSkillsHealth           `json:"health,omitempty"`
	Hunger           GameConfigSkillsHunger           `json:"hunger,omitempty"`
	LootChance       GameConfigSkillsLootChance       `json:"lootChance,omitempty"`
	Management       GameConfigSkillsManagement       `json:"management,omitempty"`
	Precision        GameConfigSkillsPrecision        `json:"precision,omitempty"`
	Production       GameConfigSkillsProduction       `json:"production,omitempty"`
	Stamina          GameConfigSkillsStamina          `json:"stamina,omitempty"`
}

// GameConfigSkillLevel describes a skill's cost and value at a given level.
type GameConfigSkillLevel struct {
	Cost          float64 `json:"cost,omitempty"`
	IsABar        bool    `json:"isABar,omitempty"`
	TotalCost     float64 `json:"totalCost,omitempty"`
	UnlockAtLevel int     `json:"unlockAtLevel,omitempty"`
	Value         float64 `json:"value,omitempty"`
}

// GameConfigSkillsArmor is a nested component of GameConfigSkills.
type GameConfigSkillsArmor struct {
	Levels  map[string]GameConfigSkillLevel `json:"levels,omitempty"`
	SoftCap int                             `json:"softCap,omitempty"`
}

// GameConfigSkillsAttack is a nested component of GameConfigSkills.
type GameConfigSkillsAttack struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
}

// GameConfigSkillsCompanies is a nested component of GameConfigSkills.
type GameConfigSkillsCompanies struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
}

// GameConfigSkillsCriticalChance is a nested component of GameConfigSkills.
type GameConfigSkillsCriticalChance struct {
	Levels             map[string]GameConfigSkillLevel `json:"levels,omitempty"`
	SkillOverflow      string                          `json:"skillOverflow,omitempty"`
	SkillOverflowValue int                             `json:"skillOverflowValue,omitempty"`
}

// GameConfigSkillsCriticalDamages is a nested component of GameConfigSkills.
type GameConfigSkillsCriticalDamages struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
}

// GameConfigSkillsDodge is a nested component of GameConfigSkills.
type GameConfigSkillsDodge struct {
	Levels  map[string]GameConfigSkillLevel `json:"levels,omitempty"`
	SoftCap int                             `json:"softCap,omitempty"`
}

// GameConfigSkillsEnergy is a nested component of GameConfigSkills.
type GameConfigSkillsEnergy struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
}

// GameConfigSkillsEntrepreneurship is a nested component of GameConfigSkills.
type GameConfigSkillsEntrepreneurship struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
}

// GameConfigSkillsHealth is a nested component of GameConfigSkills.
type GameConfigSkillsHealth struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
}

// GameConfigSkillsHunger is a nested component of GameConfigSkills.
type GameConfigSkillsHunger struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
}

// GameConfigSkillsLootChance is a nested component of GameConfigSkills.
type GameConfigSkillsLootChance struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
}

// GameConfigSkillsManagement is a nested component of GameConfigSkills.
type GameConfigSkillsManagement struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
}

// GameConfigSkillsPrecision is a nested component of GameConfigSkills.
type GameConfigSkillsPrecision struct {
	Levels             map[string]GameConfigSkillLevel `json:"levels,omitempty"`
	SkillOverflow      string                          `json:"skillOverflow,omitempty"`
	SkillOverflowValue int                             `json:"skillOverflowValue,omitempty"`
}

// GameConfigSkillsProduction is a nested component of GameConfigSkills.
type GameConfigSkillsProduction struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
}

// GameConfigSkillsStamina is a nested component of GameConfigSkills.
type GameConfigSkillsStamina struct {
	Levels map[string]GameConfigSkillLevel `json:"levels,omitempty"`
}

// GameConfigSubSkinReward is a nested component of GameConfig.
type GameConfigSubSkinReward struct {
	Deadline string `json:"deadline,omitempty"`
	SkinKey  string `json:"skinKey,omitempty"`
}

// GameConfigUnrest is a nested component of GameConfig.
type GameConfigUnrest struct {
	BarMultiplier                            int `json:"barMultiplier,omitempty"`
	BattleCooldownHours                      int `json:"battleCooldownHours,omitempty"`
	BattleStartCost                          int `json:"battleStartCost,omitempty"`
	BordersOpenDays                          int `json:"bordersOpenDays,omitempty"`
	ContributionCooldownAfterRevolutionHours int `json:"contributionCooldownAfterRevolutionHours,omitempty"`
	ContributionCost                         int `json:"contributionCost,omitempty"`
	ContributionMinLevel                     int `json:"contributionMinLevel,omitempty"`
	ContributionValue                        int `json:"contributionValue,omitempty"`
	NominationPeriodHours                    int `json:"nominationPeriodHours,omitempty"`
}

// GameConfigUpgrade is a nested component of GameConfig.
type GameConfigUpgrade struct {
	ConstructionContributionMinLevel  int `json:"constructionContributionMinLevel,omitempty"`
	ConstructionPointsPerContribution int `json:"constructionPointsPerContribution,omitempty"`
	RefundPercent                     int `json:"refundPercent,omitempty"`
	RegionDowngradeCooldownHours      int `json:"regionDowngradeCooldownHours,omitempty"`
	RegionUpgradeCooldownHours        int `json:"regionUpgradeCooldownHours,omitempty"`
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
	PendingDurationHours int                               `json:"pendingDurationHours,omitempty"`
}

// GameConfigUpgradesConfigBase is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigBase struct {
	CanBeDestroyed       bool                              `json:"canBeDestroyed,omitempty"`
	CanBeDisabled        bool                              `json:"canBeDisabled,omitempty"`
	CanDowngrade         bool                              `json:"canDowngrade,omitempty"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels,omitempty"`
	PendingDurationHours int                               `json:"pendingDurationHours,omitempty"`
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
	PendingDurationHours int                               `json:"pendingDurationHours,omitempty"`
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
	PendingDurationHours int                               `json:"pendingDurationHours,omitempty"`
}

// GameConfigUpgradesConfigPacificationCenter is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigPacificationCenter struct {
	CanBeDestroyed       bool                              `json:"canBeDestroyed,omitempty"`
	CanBeDisabled        bool                              `json:"canBeDisabled,omitempty"`
	CanDowngrade         bool                              `json:"canDowngrade,omitempty"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels,omitempty"`
	PendingDurationHours int                               `json:"pendingDurationHours,omitempty"`
}

// GameConfigUpgradesConfigStorage is a nested component of GameConfigUpgradesConfig.
type GameConfigUpgradesConfigStorage struct {
	CanDowngrade         bool                              `json:"canDowngrade,omitempty"`
	Levels               map[string]GameConfigUpgradeLevel `json:"levels,omitempty"`
	PendingDurationHours int                               `json:"pendingDurationHours,omitempty"`
}

// GameConfigUser is a nested component of GameConfig.
type GameConfigUser struct {
	ActiveCitizenMinLevel     int                         `json:"activeCitizenMinLevel,omitempty"`
	CanTakeControlAtLevel     int                         `json:"canTakeControlAtLevel,omitempty"`
	ChatMinLevel              int                         `json:"chatMinLevel,omitempty"`
	CitizenshipDaysCooldown   int                         `json:"citizenshipDaysCooldown,omitempty"`
	ConstructionEnergyCost    int                         `json:"constructionEnergyCost,omitempty"`
	DailyReward               GameConfigUserDailyReward   `json:"dailyReward,omitempty"`
	DailyXp                   int                         `json:"dailyXp,omitempty"`
	DonationMinLevel          int                         `json:"donationMinLevel,omitempty"`
	EnergyCostPerAction       int                         `json:"energyCostPerAction,omitempty"`
	EquipmentSets             GameConfigUserEquipmentSets `json:"equipmentSets,omitempty"`
	FieldsToPopulate          string                      `json:"fieldsToPopulate,omitempty"`
	IsInactiveAfterDays       int                         `json:"isInactiveAfterDays,omitempty"`
	Kits                      GameConfigUserKits          `json:"kits,omitempty"`
	MarketMinLevel            int                         `json:"marketMinLevel,omitempty"`
	MaxConstructionPoints     int                         `json:"maxConstructionPoints,omitempty"`
	MaxEnergy                 int                         `json:"maxEnergy,omitempty"`
	MaxHunger                 int                         `json:"maxHunger,omitempty"`
	Prestige                  GameConfigUserPrestige      `json:"prestige,omitempty"`
	RegenDividedBy            int                         `json:"regenDividedBy,omitempty"`
	ResetSkillDaysCooldown    int                         `json:"resetSkillDaysCooldown,omitempty"`
	ResetSkillsCostPerPoint   float64                     `json:"resetSkillsCostPerPoint,omitempty"`
	TakeControlCooldownInDays int                         `json:"takeControlCooldownInDays,omitempty"`
	XpPerAction               int                         `json:"xpPerAction,omitempty"`
}

// GameConfigUserDailyReward is a nested component of GameConfigUser.
type GameConfigUserDailyReward struct {
	Case1 int `json:"case1,omitempty"`
	Money int `json:"money,omitempty"`
	XP    int `json:"xp,omitempty"`
}

// GameConfigUserEquipmentSets is a nested component of GameConfigUser.
type GameConfigUserEquipmentSets struct {
	NonPremiumMax int `json:"nonPremiumMax,omitempty"`
	PremiumMax    int `json:"premiumMax,omitempty"`
}

// GameConfigUserKits is a nested component of GameConfigUser.
type GameConfigUserKits struct {
	NonPremiumMax int `json:"nonPremiumMax,omitempty"`
	PremiumMax    int `json:"premiumMax,omitempty"`
}

// GameConfigUserPrestige is a nested component of GameConfigUser.
type GameConfigUserPrestige struct {
	Enabled                     bool `json:"enabled,omitempty"`
	EnabledLevel                int  `json:"enabledLevel,omitempty"`
	SkillBonusLevels            int  `json:"skillBonusLevels,omitempty"`
	XpMultiplierPerLevelPercent int  `json:"xpMultiplierPerLevelPercent,omitempty"`
	XpMultiplierPercent         int  `json:"xpMultiplierPercent,omitempty"`
	XpPenalty                   int  `json:"xpPenalty,omitempty"`
}

// GameConfigWorker is a nested component of GameConfig.
type GameConfigWorker struct {
	FidelityProductionBonusPercent int `json:"fidelityProductionBonusPercent,omitempty"`
	MaxFidelity                    int `json:"maxFidelity,omitempty"`
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
	ID                  string          `json:"_id,omitempty"`
	CongressMembers     []string        `json:"congressMembers,omitempty"`
	CountryID           string          `json:"country,omitempty"`
	Dates               GovernmentDates `json:"dates,omitempty"`
	MinOfDefense        string          `json:"minOfDefense,omitempty"`
	MinOfEconomy        string          `json:"minOfEconomy,omitempty"`
	MinOfForeignAffairs string          `json:"minOfForeignAffairs,omitempty"`
	President           string          `json:"president,omitempty"`
	VicePresident       string          `json:"vicePresident,omitempty"`
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
	IsActive                  bool                 `json:"isActive,omitempty"`
	RoundHistory              []BattleRoundHistory `json:"roundHistory,omitempty"`
	RoundIds                  []string             `json:"roundIds,omitempty"`
}

// LiveRoundSummary is a nested component of LiveBattleData.
type LiveRoundSummary struct {
	ActualTickPoints int    `json:"actualTickPoints,omitempty"`
	AttackerDamages  int    `json:"attackerDamages,omitempty"`
	AttackerPoints   int    `json:"attackerPoints,omitempty"`
	DefenderDamages  int    `json:"defenderDamages,omitempty"`
	DefenderPoints   int    `json:"defenderPoints,omitempty"`
	IsActive         bool   `json:"isActive,omitempty"`
	NextTickAt       string `json:"nextTickAt,omitempty"`
	RoundID          string `json:"roundId,omitempty"`
}

// MULeveling is a nested component of MilitaryUnit.
type MULeveling struct {
	Level          int `json:"level,omitempty"`
	MonthlyDamages int `json:"monthlyDamages,omitempty"`
}

// MURoles is a nested component of MilitaryUnit.
type MURoles struct {
	Commanders []string `json:"commanders,omitempty"`
	Managers   []string `json:"managers,omitempty"`
}

// MilitaryUnit is a player-managed military organization.
type MilitaryUnit struct {
	ID                  string                          `json:"_id,omitempty"`
	ActiveUpgradeLevels MilitaryUnitActiveUpgradeLevels `json:"activeUpgradeLevels,omitempty"`
	AvatarURL           string                          `json:"avatarUrl,omitempty"`
	CountryID           string                          `json:"country,omitempty"`
	CreatedAt           string                          `json:"createdAt,omitempty"`
	Leveling            MULeveling                      `json:"leveling,omitempty"`
	Members             []string                        `json:"members,omitempty"`
	MercenaryReputation float64                         `json:"mercenaryReputation,omitempty"`
	Name                string                          `json:"name,omitempty"`
	Rankings            map[string]RankingItem          `json:"rankings,omitempty"`
	Region              string                          `json:"region,omitempty"`
	Roles               MURoles                         `json:"roles,omitempty"`
	UpdatedAt           string                          `json:"updatedAt,omitempty"`
	UserID              string                          `json:"user,omitempty"`
}

// MilitaryUnitActiveUpgradeLevels is a nested component of MilitaryUnit.
type MilitaryUnitActiveUpgradeLevels struct {
	Dormitories int `json:"dormitories,omitempty"`
}

// Ranking contains a global ranking and its entries.
type Ranking struct {
	ID         string            `json:"_id,omitempty"`
	IsGlobal   bool              `json:"isGlobal,omitempty"`
	Items      []RankingItems    `json:"items,omitempty"`
	TierValues RankingTierValues `json:"tierValues,omitempty"`
	Type       string            `json:"type,omitempty"`
}

// RankingItem represents a tiered leaderboard rank entry.
type RankingItem struct {
	Rank  int     `json:"rank,omitempty"`
	Tier  string  `json:"tier,omitempty"`
	Value float64 `json:"value,omitempty"`
}

// RankingItems is a nested component of Ranking.
type RankingItems struct {
	ID        string `json:"_id,omitempty"`
	CountryID string `json:"country,omitempty"`
	Rank      int    `json:"rank,omitempty"`
	Tier      string `json:"tier,omitempty"`
	Value     int64  `json:"value,omitempty"`
}

// RankingTierValues is a nested component of Ranking.
type RankingTierValues struct {
	Bronze   int   `json:"bronze,omitempty"`
	Diamond  int64 `json:"diamond,omitempty"`
	Gold     int   `json:"gold,omitempty"`
	Master   int64 `json:"master,omitempty"`
	Platinum int64 `json:"platinum,omitempty"`
	Silver   int   `json:"silver,omitempty"`
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
	CountryID                    string                    `json:"country,omitempty"`
	CountryCode                  string                    `json:"countryCode,omitempty"`
	CurrentPopulation            int                       `json:"currentPopulation,omitempty"`
	Dates                        RegionDates               `json:"dates,omitempty"`
	Deposit                      *RegionDeposit            `json:"deposit,omitempty"`
	Development                  float64                   `json:"development,omitempty"`
	HasCoast                     bool                      `json:"hasCoast,omitempty"`
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
	ResistanceMax                int                       `json:"resistanceMax,omitempty"`
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
	BonusPercent int    `json:"bonusPercent,omitempty"`
	EndsAt       string `json:"endsAt,omitempty"`
	StartsAt     string `json:"startsAt,omitempty"`
	Type         string `json:"type,omitempty"`
}

// RegionStats is a nested component of Region.
type RegionStats struct {
	InvestedMoney int `json:"investedMoney,omitempty"`
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
	ConstructionPoints    int                             `json:"constructionPoints,omitempty"`
	ConstructionStartedAt string                          `json:"constructionStartedAt,omitempty"`
	InvestedMoney         int                             `json:"investedMoney,omitempty"`
	IsUnderConstruction   *bool                           `json:"isUnderConstruction,omitempty"`
	LastConstructions     []RegionUpgradeLastConstruction `json:"lastConstructions,omitempty"`
	Level                 int                             `json:"level,omitempty"`
	Status                string                          `json:"status,omitempty"`
	StatusChangedAt       string                          `json:"statusChangedAt,omitempty"`
}

// RegionUpgradesV2UpgradesBunker is a nested component of RegionUpgradesV2Upgrades.
type RegionUpgradesV2UpgradesBunker struct {
	ConstructionEndedAt   string                          `json:"constructionEndedAt,omitempty"`
	ConstructionPoints    int                             `json:"constructionPoints,omitempty"`
	ConstructionStartedAt string                          `json:"constructionStartedAt,omitempty"`
	InvestedMoney         int                             `json:"investedMoney,omitempty"`
	IsUnderConstruction   *bool                           `json:"isUnderConstruction,omitempty"`
	LastConstructions     []RegionUpgradeLastConstruction `json:"lastConstructions,omitempty"`
	Level                 int                             `json:"level,omitempty"`
	Status                string                          `json:"status,omitempty"`
	StatusChangedAt       string                          `json:"statusChangedAt,omitempty"`
}

// RegionUpgradeLastConstruction records a completed upgrade construction.
type RegionUpgradeLastConstruction struct {
	Construction   int    `json:"construction,omitempty"`
	ConstructionAt string `json:"constructionAt,omitempty"`
	UserID         string `json:"user,omitempty"`
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
	MaxState          int             `json:"maxState,omitempty"`
	Quantity          int             `json:"quantity,omitempty"`
	Skills            EquipmentSkills `json:"skills,omitempty"`
	State             int             `json:"state,omitempty"`
	Type              string          `json:"type,omitempty"`
}

// RoundLastHit is a nested component of RoundSideCombatStats.
type RoundLastHit struct {
	ID            string           `json:"_id,omitempty"`
	Ammo          string           `json:"ammo,omitempty"`
	Damages       int              `json:"damages,omitempty"`
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
	ActualTickPoints int    `json:"actualTickPoints,omitempty"`
	NextTickAt       string `json:"nextTickAt,omitempty"`
	TicksCount       int    `json:"ticksCount,omitempty"`
}

// RoundSideCombatStats is a nested component of BattleRound.
type RoundSideCombatStats struct {
	CountryID string         `json:"country,omitempty"`
	Damages   int            `json:"damages,omitempty"`
	HitCount  int            `json:"hitCount,omitempty"`
	LastHits  []RoundLastHit `json:"lastHits,omitempty"`
	Points    int            `json:"points,omitempty"`
}

// RoundWeapon is a nested component of RoundLastHit.
type RoundWeapon struct {
	ID                string       `json:"_id,omitempty"`
	Code              string       `json:"code,omitempty"`
	LastAcquisitionAt string       `json:"lastAcquisitionAt,omitempty"`
	MaxState          int          `json:"maxState,omitempty"`
	Quantity          int          `json:"quantity,omitempty"`
	Skills            WeaponSkills `json:"skills,omitempty"`
	State             int          `json:"state,omitempty"`
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
	ID                  string `json:"_id,omitempty"`
	CreatedAt           string `json:"createdAt,omitempty"`
	DependantUsersCount int    `json:"dependantUsersCount,omitempty"`
	InvestedConcrete    int    `json:"investedConcrete,omitempty"`
	InvestedMoney       int    `json:"investedMoney,omitempty"`
	InvestedSteel       int    `json:"investedSteel,omitempty"`
	LastDowngradeAt     string `json:"lastDowngradeAt,omitempty"`
	LastUpgradeAt       string `json:"lastUpgradeAt,omitempty"`
	Level               int    `json:"level,omitempty"`
	Region              string `json:"region,omitempty"`
	Status              string `json:"status,omitempty"`
	StatusChangedAt     string `json:"statusChangedAt,omitempty"`
	UpdatedAt           string `json:"updatedAt,omitempty"`
	UpgradeType         string `json:"upgradeType,omitempty"`
	WillBeActiveAt      string `json:"willBeActiveAt,omitempty"`
}

// User is a full player profile, including skills, statistics, and equipment.
type User struct {
	ID                    string                 `json:"_id,omitempty"`
	AvailableColorSchemes []string               `json:"availableColorSchemes,omitempty"`
	AvatarURL             string                 `json:"avatarUrl,omitempty"`
	CanOnboard            bool                   `json:"canOnboard,omitempty"`
	Company               string                 `json:"company,omitempty"`
	CountryID             string                 `json:"country,omitempty"`
	CreatedAt             string                 `json:"createdAt,omitempty"`
	Dates                 UserDates              `json:"dates,omitempty"`
	Discord               UserDiscord            `json:"discord,omitempty"`
	EmailVerified         bool                   `json:"emailVerified,omitempty"`
	Equipment             UserEquipment          `json:"equipment,omitempty"`
	FinishedTours         UserFinishedTours      `json:"finishedTours,omitempty"`
	IsActive              bool                   `json:"isActive,omitempty"`
	Leveling              UserLeveling           `json:"leveling,omitempty"`
	Location              string                 `json:"location,omitempty"`
	MilitaryRank          int                    `json:"militaryRank,omitempty"`
	Missions              UserMissions           `json:"missions,omitempty"`
	MUID                  string                 `json:"mu,omitempty"`
	MuMaxLevelRewarded    int                    `json:"muMaxLevelRewarded,omitempty"`
	Party                 string                 `json:"party,omitempty"`
	Preferences           UserPreferences        `json:"preferences,omitempty"`
	Rankings              map[string]RankingItem `json:"rankings,omitempty"`
	Region                string                 `json:"region,omitempty"`
	ShouldUpdateProfile   bool                   `json:"shouldUpdateProfile,omitempty"`
	Skills                UserSkills             `json:"skills,omitempty"`
	Stats                 UserStats              `json:"stats,omitempty"`
	UpdatedAt             string                 `json:"updatedAt,omitempty"`
	Username              string                 `json:"username,omitempty"`
	UsernameLower         string                 `json:"usernameLower,omitempty"`
}

// UserDates is a nested component of User.
type UserDates struct {
	LastCompanyJoinedAt        string   `json:"lastCompanyJoinedAt,omitempty"`
	LastConnectionAt           string   `json:"lastConnectionAt,omitempty"`
	LastCountryMessageCheckAt  string   `json:"lastCountryMessageCheckAt,omitempty"`
	LastDailyCalendarClaimedAt string   `json:"lastDailyCalendarClaimedAt,omitempty"`
	LastEventsCheckAt          string   `json:"lastEventsCheckAt,omitempty"`
	LastFreeTravelHomeAt       string   `json:"lastFreeTravelHomeAt,omitempty"`
	LastGlobalMessageCheckAt   string   `json:"lastGlobalMessageCheckAt,omitempty"`
	LastHelpAskedAt            string   `json:"lastHelpAskedAt,omitempty"`
	LastHiresAt                []string `json:"lastHiresAt,omitempty"`
	LastNotificationsCheckAt   string   `json:"lastNotificationsCheckAt,omitempty"`
	LastSeenAnnouncementAt     string   `json:"lastSeenAnnouncementAt,omitempty"`
	LastSkillsResetAt          string   `json:"lastSkillsResetAt,omitempty"`
	LastWorkAt                 string   `json:"lastWorkAt,omitempty"`
	LastWorkOfferApplications  []string `json:"lastWorkOfferApplications,omitempty"`
}

// UserDiscord is a nested component of User.
type UserDiscord struct {
	IsVerified bool `json:"isVerified,omitempty"`
}

// UserEquipment is a nested component of User.
type UserEquipment struct {
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
	AvailableSkillPoints int `json:"availableSkillPoints,omitempty"`
	DailyXPLeft          int `json:"dailyXpLeft,omitempty"`
	FreeReset            int `json:"freeReset,omitempty"`
	Level                int `json:"level,omitempty"`
	Prestige             int `json:"prestige,omitempty"`
	PrestigeLevel        int `json:"prestigeLevel,omitempty"`
	SpentSkillPoints     int `json:"spentSkillPoints,omitempty"`
	TotalSkillPoints     int `json:"totalSkillPoints,omitempty"`
	TotalXP              int `json:"totalXp,omitempty"`
	XP                   int `json:"xp,omitempty"`
}

// UserLite is a compact player profile for listings and search results.
type UserLite struct {
	ID           string                 `json:"_id,omitempty"`
	AvatarURL    string                 `json:"avatarUrl,omitempty"`
	CountryID    string                 `json:"country,omitempty"`
	CreatedAt    string                 `json:"createdAt,omitempty"`
	Dates        UserDates              `json:"dates,omitempty"`
	IsActive     bool                   `json:"isActive,omitempty"`
	Leveling     UserLeveling           `json:"leveling,omitempty"`
	MilitaryRank int                    `json:"militaryRank,omitempty"`
	MUID         string                 `json:"mu,omitempty"`
	Rankings     map[string]RankingItem `json:"rankings,omitempty"`
	Skills       UserSkills             `json:"skills,omitempty"`
	Stats        UserStats              `json:"stats,omitempty"`
	Username     string                 `json:"username,omitempty"`
}

// UserMissions is a nested component of User.
type UserMissions struct {
	ClaimedAt               UserMissionsClaimedAt `json:"claimedAt,omitempty"`
	ClaimedCount            int                   `json:"claimedCount,omitempty"`
	RerolledDailyMissions   int                   `json:"rerolledDailyMissions,omitempty"`
	RerolledMonthlyMissions int                   `json:"rerolledMonthlyMissions,omitempty"`
	RerolledWeeklyMissions  int                   `json:"rerolledWeeklyMissions,omitempty"`
	TotalXpClaimed          int                   `json:"totalXpClaimed,omitempty"`
}

// UserMissionsClaimedAt is a nested component of UserMissions.
type UserMissionsClaimedAt struct {
	Daily    string `json:"daily,omitempty"`
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

// UserSkills is a nested component of User.
type UserSkills struct {
	Armor            UserSkillsArmor            `json:"armor,omitempty"`
	Attack           UserSkillsAttack           `json:"attack,omitempty"`
	Companies        UserSkillsCompanies        `json:"companies,omitempty"`
	CriticalChance   UserSkillsCriticalChance   `json:"criticalChance,omitempty"`
	CriticalDamages  UserSkillsCriticalDamages  `json:"criticalDamages,omitempty"`
	Dodge            UserSkillsDodge            `json:"dodge,omitempty"`
	Energy           UserSkillsEnergy           `json:"energy,omitempty"`
	Entrepreneurship UserSkillsEntrepreneurship `json:"entrepreneurship,omitempty"`
	Health           UserSkillsHealth           `json:"health,omitempty"`
	Hunger           UserSkillsHunger           `json:"hunger,omitempty"`
	LootChance       UserSkillsLootChance       `json:"lootChance,omitempty"`
	Management       UserSkillsManagement       `json:"management,omitempty"`
	Precision        UserSkillsPrecision        `json:"precision,omitempty"`
	Production       UserSkillsProduction       `json:"production,omitempty"`
	Stamina          UserSkillsStamina          `json:"stamina,omitempty"`
}

// UserSkillsArmor is a nested component of UserSkills.
type UserSkillsArmor struct {
	Equipment         int             `json:"equipment,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap int             `json:"totalAfterSoftCap,omitempty"`
	Value             *float64        `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsAttack is a nested component of UserSkills.
type UserSkillsAttack struct {
	AmmoPercent         int             `json:"ammoPercent,omitempty"`
	BuffsPercent        int             `json:"buffsPercent,omitempty"`
	DebuffsPercent      int             `json:"debuffsPercent,omitempty"`
	Equipment           *float64        `json:"equipment,omitempty"`
	Level               int             `json:"level,omitempty"`
	Limited             json.RawMessage `json:"limited,omitempty"`
	MilitaryRankPercent float64         `json:"militaryRankPercent,omitempty"`
	Overflow            json.RawMessage `json:"overflow,omitempty"`
	Prestige            int             `json:"prestige,omitempty"`
	Total               int             `json:"total,omitempty"`
	TotalAfterSoftCap   *float64        `json:"totalAfterSoftCap,omitempty"`
	Value               int             `json:"value,omitempty"`
	Weapon              int             `json:"weapon,omitempty"`
}

// UserSkillsCompanies is a nested component of UserSkills.
type UserSkillsCompanies struct {
	Equipment         *float64        `json:"equipment,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsCriticalChance is a nested component of UserSkills.
type UserSkillsCriticalChance struct {
	Equipment         *float64        `json:"equipment,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            int             `json:"weapon,omitempty"`
}

// UserSkillsCriticalDamages is a nested component of UserSkills.
type UserSkillsCriticalDamages struct {
	Equipment         int             `json:"equipment,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsDodge is a nested component of UserSkills.
type UserSkillsDodge struct {
	Equipment         int             `json:"equipment,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap int             `json:"totalAfterSoftCap,omitempty"`
	Value             *float64        `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsEnergy is a nested component of UserSkills.
type UserSkillsEnergy struct {
	CurrentBarValue   float64         `json:"currentBarValue,omitempty"`
	Equipment         *float64        `json:"equipment,omitempty"`
	HourlyBarRegen    int             `json:"hourlyBarRegen,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsEntrepreneurship is a nested component of UserSkills.
type UserSkillsEntrepreneurship struct {
	CurrentBarValue   float64         `json:"currentBarValue,omitempty"`
	Equipment         *float64        `json:"equipment,omitempty"`
	HourlyBarRegen    int             `json:"hourlyBarRegen,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsHealth is a nested component of UserSkills.
type UserSkillsHealth struct {
	CurrentBarValue   float64         `json:"currentBarValue,omitempty"`
	Equipment         *float64        `json:"equipment,omitempty"`
	HourlyBarRegen    int             `json:"hourlyBarRegen,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsHunger is a nested component of UserSkills.
type UserSkillsHunger struct {
	CurrentBarValue   float64         `json:"currentBarValue,omitempty"`
	Equipment         *float64        `json:"equipment,omitempty"`
	HourlyBarRegen    float64         `json:"hourlyBarRegen,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsLootChance is a nested component of UserSkills.
type UserSkillsLootChance struct {
	Equipment         *float64        `json:"equipment,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsManagement is a nested component of UserSkills.
type UserSkillsManagement struct {
	Equipment         *float64        `json:"equipment,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsPrecision is a nested component of UserSkills.
type UserSkillsPrecision struct {
	Equipment         int             `json:"equipment,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsProduction is a nested component of UserSkills.
type UserSkillsProduction struct {
	CurrentBarValue   float64         `json:"currentBarValue,omitempty"`
	Equipment         *float64        `json:"equipment,omitempty"`
	HourlyBarRegen    float64         `json:"hourlyBarRegen,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserSkillsStamina is a nested component of UserSkills.
type UserSkillsStamina struct {
	CurrentBarValue   float64         `json:"currentBarValue,omitempty"`
	Equipment         *float64        `json:"equipment,omitempty"`
	HourlyBarRegen    int             `json:"hourlyBarRegen,omitempty"`
	Level             int             `json:"level,omitempty"`
	Limited           json.RawMessage `json:"limited,omitempty"`
	Overflow          json.RawMessage `json:"overflow,omitempty"`
	Prestige          int             `json:"prestige,omitempty"`
	Total             int             `json:"total,omitempty"`
	TotalAfterSoftCap *float64        `json:"totalAfterSoftCap,omitempty"`
	Value             int             `json:"value,omitempty"`
	Weapon            *float64        `json:"weapon,omitempty"`
}

// UserStats is a nested component of User.
type UserStats struct {
	Case1                 UserStatsCase1      `json:"case1,omitempty"`
	Case2                 UserStatsCase2      `json:"case2,omitempty"`
	DailyRewardBestStreak int                 `json:"dailyRewardBestStreak,omitempty"`
	DailyRewardStreak     int                 `json:"dailyRewardStreak,omitempty"`
	DamagesCount          int                 `json:"damagesCount,omitempty"`
	MilestoneTiersClaimed int                 `json:"milestoneTiersClaimed,omitempty"`
	MilestoneXpClaimed    int                 `json:"milestoneXpClaimed,omitempty"`
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
	Common   int `json:"common,omitempty"`
	Rare     int `json:"rare,omitempty"`
	Uncommon int `json:"uncommon,omitempty"`
}

// UserStatsCase2 is a nested component of UserStats.
type UserStatsCase2 struct {
	ByRarity    UserStatsCase2ByRarity `json:"byRarity,omitempty"`
	OpenedCount int                    `json:"openedCount,omitempty"`
}

// UserStatsCase2ByRarity is a nested component of UserStatsCase2.
type UserStatsCase2ByRarity struct {
	Rare     int `json:"rare,omitempty"`
	Uncommon int `json:"uncommon,omitempty"`
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
	Ammo       int `json:"ammo,omitempty"`
	Bread      int `json:"bread,omitempty"`
	Coca       int `json:"coca,omitempty"`
	Cocain     int `json:"cocain,omitempty"`
	Concrete   int `json:"concrete,omitempty"`
	CookedFish int `json:"cookedFish,omitempty"`
	Fish       int `json:"fish,omitempty"`
	Grain      int `json:"grain,omitempty"`
	Iron       int `json:"iron,omitempty"`
	Lead       int `json:"lead,omitempty"`
	LightAmmo  int `json:"lightAmmo,omitempty"`
	Limestone  int `json:"limestone,omitempty"`
	Livestock  int `json:"livestock,omitempty"`
	Oil        int `json:"oil,omitempty"`
	Paper      int `json:"paper,omitempty"`
	Petroleum  int `json:"petroleum,omitempty"`
	Steak      int `json:"steak,omitempty"`
	Steel      int `json:"steel,omitempty"`
	Wood       int `json:"wood,omitempty"`
}

// WeaponItem is a nested component of CurrentEquipment.
type WeaponItem struct {
	ID                string       `json:"_id,omitempty"`
	Code              string       `json:"code,omitempty"`
	LastAcquisitionAt string       `json:"lastAcquisitionAt,omitempty"`
	MaxState          int          `json:"maxState,omitempty"`
	Quantity          int          `json:"quantity,omitempty"`
	Skills            WeaponSkills `json:"skills,omitempty"`
	State             int          `json:"state,omitempty"`
}

// WeaponSkills is a nested component of RoundWeapon.
type WeaponSkills struct {
	Attack         int `json:"attack,omitempty"`
	CriticalChance int `json:"criticalChance,omitempty"`
}
