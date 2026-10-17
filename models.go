package warera

import "time"

// CursorPage represents a generic cursor-paginated page of items.
type CursorPage[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"nextCursor,omitempty"`
}

// HasNext reports whether there is a subsequent page.
func (p *CursorPage[T]) HasNext() bool {
	return p != nil && p.NextCursor != "" && p.NextCursor != "undefined"
}

// --- Country Models ---

type Country struct {
	ID                   string                         `json:"_id"`
	Name                 string                         `json:"name"`
	Code                 string                         `json:"code"`
	Money                float64                        `json:"money"`
	Scheme               string                         `json:"scheme,omitempty"`
	MapAccent            string                         `json:"mapAccent,omitempty"`
	Taxes                CountryTaxes                   `json:"taxes"`
	Unrest               CountryUnrest                  `json:"unrest"`
	Rankings             map[string]RankingItem         `json:"rankings,omitempty"`
	Allies               []string                       `json:"allies,omitempty"`
	WarsWith             []string                       `json:"warsWith,omitempty"`
	DefensivePacts       []string                       `json:"defensivePacts,omitempty"`
	NonAggressionUntil   map[string]string              `json:"nonAggressionUntil,omitempty"`
	CoreDevelopment      float64                        `json:"coreDevelopment,omitempty"`
	AverageDevelopment   float64                        `json:"averageDevelopment,omitempty"`
	CurrentDevelopment   float64                        `json:"currentDevelopment,omitempty"`
	CurrentPopulation    int                            `json:"currentPopulation,omitempty"`
	SpecializedItem      string                         `json:"specializedItem,omitempty"`
	Enemy                string                         `json:"enemy,omitempty"`
	AllianceID           string                         `json:"allianceId,omitempty"`
	RulingParty          string                         `json:"rulingParty,omitempty"`
	DiscordURL           string                         `json:"discordUrl,omitempty"`
	PinnedArticle        string                         `json:"pinnedArticle,omitempty"`
	StrategicResources   *CountryStrategicResourcesInfo `json:"strategicResources,omitempty"`
	UpdatedAt            string                         `json:"updatedAt,omitempty"`
}

type CountryTaxes struct {
	Income   float64 `json:"income"`
	Market   float64 `json:"market"`
	SelfWork float64 `json:"selfWork"`
}

type CountryUnrest struct {
	Bar                int    `json:"bar"`
	BarMax             int    `json:"barMax"`
	LastContributionAt string `json:"lastContributionAt,omitempty"`
}

type RankingItem struct {
	Value any    `json:"value"`
	Rank  int    `json:"rank"`
	Tier  string `json:"tier"`
}

type CountryStrategicResourcesInfo struct {
	Resources map[string][]string `json:"resources,omitempty"`
	Bonuses   map[string]float64  `json:"bonuses,omitempty"`
}

// --- Government Models ---

type Government struct {
	CountryID string         `json:"countryId,omitempty"`
	President *GovernmentRole `json:"president,omitempty"`
	Ministers []GovernmentRole `json:"ministers,omitempty"`
	Congress  []GovernmentRole `json:"congress,omitempty"`
}

type GovernmentRole struct {
	UserID    string `json:"userId"`
	Username  string `json:"username,omitempty"`
	PartyID   string `json:"partyId,omitempty"`
	RoleTitle string `json:"roleTitle,omitempty"`
}

// --- Company Models ---

type Company struct {
	ID           string         `json:"_id"`
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	Code         string         `json:"code,omitempty"`
	CountryID    string         `json:"country,omitempty"`
	RegionID     string         `json:"region,omitempty"`
	UserID       string         `json:"user,omitempty"`
	Level        int            `json:"level,omitempty"`
	Quality      int            `json:"quality,omitempty"`
	Money        float64        `json:"money,omitempty"`
	WorkersCount int            `json:"workersCount,omitempty"`
	CreatedAt    string         `json:"createdAt,omitempty"`
	UpdatedAt    string         `json:"updatedAt,omitempty"`
}

type CompanyFilter struct {
	UserID  string `json:"userId,omitempty"`
	PerPage int    `json:"perPage,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
}

// --- Region Models ---

type Region struct {
	ID          string   `json:"_id"`
	Name        string   `json:"name"`
	CountryID   string   `json:"country"`
	CoreCountry string   `json:"coreCountry,omitempty"`
	Development float64  `json:"development,omitempty"`
	Resource    string   `json:"resource,omitempty"`
	Neighbors   []string `json:"neighbors,omitempty"`
}

// --- Battle Models ---

type Battle struct {
	ID               string              `json:"_id"`
	AttackerCountry  string              `json:"attackerCountry"`
	DefenderCountry  string              `json:"defenderCountry"`
	DefenderRegion   string              `json:"defenderRegion"`
	WarID            string              `json:"war,omitempty"`
	CurrentRound     int                 `json:"currentRound,omitempty"`
	IsActive         bool                `json:"isActive"`
	CreatedAt        string              `json:"createdAt,omitempty"`
	AttackerSummary  *BattleSideSummary  `json:"attackerSummary,omitempty"`
	DefenderSummary  *BattleSideSummary  `json:"defenderSummary,omitempty"`
}

type BattleSideSummary struct {
	Country        string   `json:"country"`
	Damages        int64    `json:"damages"`
	HitCount       int      `json:"hitCount"`
	WonRoundsCount int      `json:"wonRoundsCount"`
	CountryOrders  []string `json:"countryOrders,omitempty"`
	MUOrders       []string `json:"muOrders,omitempty"`
	Region         string   `json:"region,omitempty"`
}

type BattleFilter struct {
	IsActive         *bool  `json:"isActive,omitempty"`
	Limit            int    `json:"limit,omitempty"`
	Cursor           string `json:"cursor,omitempty"`
	Direction        string `json:"direction,omitempty"`
	Filter           string `json:"filter,omitempty"`
	DefenderRegionID string `json:"defenderRegionId,omitempty"`
	WarID            string `json:"warId,omitempty"`
	CountryID        string `json:"countryId,omitempty"`
}

type LiveBattleData struct {
	BattleID string         `json:"battleId"`
	Round    int            `json:"round"`
	Attacker LiveSideCombat `json:"attacker"`
	Defender LiveSideCombat `json:"defender"`
}

type LiveSideCombat struct {
	Country  string `json:"country"`
	Damages  int64  `json:"damages"`
	HitCount int    `json:"hitCount"`
}

// --- Article Models ---

type Article struct {
	ID          string       `json:"_id"`
	Title       string       `json:"title"`
	Content     string       `json:"content,omitempty"`
	Author      string       `json:"author"`
	Category    string       `json:"category"`
	Language    string       `json:"language"`
	Stats       ArticleStats `json:"stats"`
	IsPublished bool         `json:"isPublished"`
	IsDeleted   bool         `json:"isDeleted"`
	CreatedAt   string       `json:"createdAt"`
	PublishedAt string       `json:"publishedAt,omitempty"`
	UpdatedAt   string       `json:"updatedAt,omitempty"`
}

type ArticleStats struct {
	Comments int `json:"comments"`
	Likes    int `json:"likes"`
	Dislikes int `json:"dislikes"`
	Views    int `json:"views"`
	Subs     int `json:"subs"`
	Tips     int `json:"tips"`
	GemTips  int `json:"gemTips"`
	Score    int `json:"score"`
}

type ArticleFilter struct {
	Type              string   `json:"type,omitempty"` // "last", "top", "trending"
	Limit             int      `json:"limit,omitempty"`
	Cursor            string   `json:"cursor,omitempty"`
	UserID            string   `json:"userId,omitempty"`
	Categories        []string `json:"categories,omitempty"`
	Languages         []string `json:"languages,omitempty"`
	PositiveScoreOnly *bool    `json:"positiveScoreOnly,omitempty"`
}

// --- User Models ---

type User struct {
	ID        string       `json:"_id"`
	Username  string       `json:"username"`
	CountryID string       `json:"country,omitempty"`
	AvatarURL string       `json:"avatarUrl,omitempty"`
	Leveling  UserLeveling `json:"leveling,omitempty"`
	IsActive  bool         `json:"isActive,omitempty"`
	CreatedAt string       `json:"createdAt,omitempty"`
}

type UserLeveling struct {
	Level               int   `json:"level"`
	XP                  int64 `json:"xp"`
	TotalXP             int64 `json:"totalXp"`
	DailyXPLeft         int64 `json:"dailyXpLeft"`
	AvailableSkillPoints int  `json:"availableSkillPoints"`
	SpentSkillPoints    int   `json:"spentSkillPoints"`
	TotalSkillPoints    int   `json:"totalSkillPoints"`
}

type UserLite struct {
	ID        string       `json:"_id"`
	Username  string       `json:"username"`
	CountryID string       `json:"country,omitempty"`
	AvatarURL string       `json:"avatarUrl,omitempty"`
	Leveling  UserLeveling `json:"leveling,omitempty"`
}

// Level returns the player's level from Leveling.
func (u *User) Level() int {
	return u.Leveling.Level
}

// Level returns the player's level from Leveling.
func (u *UserLite) Level() int {
	return u.Leveling.Level
}

type UserFilter struct {
	CountryID string `json:"countryId"`
	Limit     int    `json:"limit,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
}

// --- Game Config & Dates Models ---

type GameDates struct {
	NextDayAt                 string `json:"nextDayAt,omitempty"`
	NextRegenAt               string `json:"nextRegenAt,omitempty"`
	PreviousDayAt             string `json:"previousDayAt,omitempty"`
	NextCongressElectionsAt   string `json:"nextCongressElectionsAt,omitempty"`
	NextPresidentialElectionsAt string `json:"nextPresidentialElectionsAt,omitempty"`
	NextMonthAt               string `json:"nextMonthAt,omitempty"`
	DailyMissionRegenAt       string `json:"dailyMissionRegenAt,omitempty"`
	WeeklyMissionRegenAt      string `json:"weeklyMissionRegenAt,omitempty"`
	MonthlyMissionRegenAt     string `json:"monthlyMissionRegenAt,omitempty"`
}

// NextDayTime parses NextDayAt as time.Time.
func (d *GameDates) NextDayTime() (time.Time, error) {
	return time.Parse(time.RFC3339, d.NextDayAt)
}

// NextRegenTime parses NextRegenAt as time.Time.
func (d *GameDates) NextRegenTime() (time.Time, error) {
	return time.Parse(time.RFC3339, d.NextRegenAt)
}

// --- Search Models ---

type SearchResult struct {
	UserIDs     []string `json:"userIds,omitempty"`
	MUIDs       []string `json:"muIds,omitempty"`
	CountryIDs  []string `json:"countryIds,omitempty"`
	RegionIDs   []string `json:"regionIds,omitempty"`
	PartyIDs    []string `json:"partyIds,omitempty"`
	AllianceIDs []string `json:"allianceIds,omitempty"`
	HasData     bool     `json:"hasData"`
}

// --- Market & Trading Models ---

type ItemPrice struct {
	ItemCode string  `json:"itemCode"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency,omitempty"`
}

type TopOrder struct {
	ID       string  `json:"_id"`
	ItemCode string  `json:"itemCode"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	UserID   string  `json:"user"`
	Type     string  `json:"type"` // "buy" or "sell"
}

// --- Military Unit (MU) Models ---

type MilitaryUnit struct {
	ID        string `json:"_id"`
	Name      string `json:"name"`
	CountryID string `json:"country"`
	LeaderID  string `json:"leader,omitempty"`
	Level     int    `json:"level,omitempty"`
	Members   int    `json:"membersCount,omitempty"`
}

type MUFilter struct {
	Limit    int    `json:"limit,omitempty"`
	Cursor   string `json:"cursor,omitempty"`
	MemberID string `json:"memberId,omitempty"`
	UserID   string `json:"userId,omitempty"`
	Search   string `json:"search,omitempty"`
}
