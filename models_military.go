package warera

// Response models for military unit members, leagues, sponsorships, state army
// spending, tournaments, contributions, and ranking history. WarEra's published
// OpenAPI document does not describe these procedures; the models follow
// recorded responses.

// MUMember is a member of a military unit and their damage and help in it.
type MUMember struct {
	ID                  string  `json:"_id"`
	CreatedAt           Time    `json:"createdAt"`
	MUID                string  `json:"mu"`
	MonthlyDamagesCount float64 `json:"monthlyDamagesCount"`
	MonthlyHelpCount    int     `json:"monthlyHelpCount"`
	TotalDamagesCount   float64 `json:"totalDamagesCount"`
	TotalHelpCount      int     `json:"totalHelpCount"`
	UpdatedAt           Time    `json:"updatedAt"`
	UserID              string  `json:"user"`
	WeeklyDamagesCount  float64 `json:"weeklyDamagesCount"`
	WeeklyHelpCount     int     `json:"weeklyHelpCount"`
}

// MULeague is a military unit's group in a week of the MU league.
type MULeague struct {
	// YearWeek identifies the league week, such as "2026-47".
	YearWeek   string `json:"yearWeek"`
	WeekEndsAt Time   `json:"weekEndsAt"`
	// League and Me are nil when the unit is not in a league this week.
	League *MULeagueGroup    `json:"league"`
	Me     *MULeagueStanding `json:"me"`
	Rows   []MULeagueRow     `json:"rows"`
	// Shares is how the unit's reward would be split among the users who
	// fought for it this week.
	Shares []MULeagueShare `json:"shares"`
	// PlacePots is the reward for each place in the group; places that earn
	// nothing are nil.
	PlacePots []*LeagueReward `json:"placePots"`
}

// MULeagueGroup is a group of military units competing in one division.
type MULeagueGroup struct {
	ID           string  `json:"_id"`
	DemoteCount  int     `json:"demoteCount"`
	GroupIndex   int     `json:"groupIndex"`
	GroupsCount  int     `json:"groupsCount"`
	MembersCount int     `json:"membersCount"`
	PromoteCount int     `json:"promoteCount"`
	StrengthAvg  float64 `json:"strengthAvg"`
}

// MULeagueStanding is a military unit's place in its league group.
type MULeagueStanding struct {
	Division int `json:"division"`
	// Move is "up", "stay", or "down": where the unit would move next week.
	Move       string  `json:"move"`
	Multiplier float64 `json:"multiplier"`
	Points     float64 `json:"points"`
	Rank       int     `json:"rank"`
}

// MULeagueRow is one military unit in a league group's standings.
type MULeagueRow struct {
	Division   int     `json:"division"`
	IsEligible bool    `json:"isEligible"`
	Move       string  `json:"move"`
	MUID       string  `json:"mu"`
	Points     float64 `json:"points"`
	// Pot is the reward for the unit's current place, or nil when it earns
	// nothing.
	Pot   *LeagueReward `json:"pot"`
	Rank  int           `json:"rank"`
	Seats int           `json:"seats"`
}

// MULeagueShare is one user's share of a military unit's league reward.
type MULeagueShare struct {
	// IsMember reports whether the user is still in the unit.
	IsMember bool                `json:"isMember"`
	Points   float64             `json:"points"`
	Reward   MULeagueShareReward `json:"reward"`
	// Share is the user's fraction of the unit's points, from 0 to 1.
	Share  float64 `json:"share"`
	UserID string  `json:"user"`
}

// MULeagueShareReward is the money and items a user would receive.
type MULeagueShareReward struct {
	// Items holds item counts keyed by item code.
	Items map[string]int `json:"items"`
	Money float64        `json:"money"`
}

// MULeagueDivisionCounts is the number of military units in each MU league
// division in a week.
type MULeagueDivisionCounts struct {
	// MUs holds the unit count of each division, indexed by division number
	// from 0, the lowest.
	MUs      []int  `json:"mus"`
	YearWeek string `json:"yearWeek"`
}

// League is a user's group in a week of the player league.
type League struct {
	// YearWeek identifies the league week, such as "2026-47".
	YearWeek   string `json:"yearWeek"`
	WeekEndsAt Time   `json:"weekEndsAt"`
	// League and Me are nil when the user is not in a league this week.
	League *LeagueGroup    `json:"league"`
	Me     *LeagueStanding `json:"me"`
	Rows   []LeagueRow     `json:"rows"`
	// PlaceRewards is the reward for each place in the group; places that earn
	// nothing are nil.
	PlaceRewards []*LeagueReward `json:"placeRewards"`
}

// LeagueGroup is a group of players competing in a league week.
type LeagueGroup struct {
	ID         string `json:"_id"`
	GroupIndex int    `json:"groupIndex"`
	// Kind is "regular", or "late" for groups of players who joined after
	// the week started.
	Kind         string  `json:"kind"`
	LevelAvg     float64 `json:"levelAvg"`
	MembersCount int     `json:"membersCount"`
	Multiplier   float64 `json:"multiplier"`
	RatingAvg    float64 `json:"ratingAvg"`
}

// LeagueStanding is a user's place in their league group.
type LeagueStanding struct {
	Points float64 `json:"points"`
	Rank   int     `json:"rank"`
}

// LeagueRow is one user in a league group's standings.
type LeagueRow struct {
	IsEligible bool    `json:"isEligible"`
	Points     float64 `json:"points"`
	Rank       int     `json:"rank"`
	UserID     string  `json:"user"`
}

// LeagueReward is the reward for a place in a player or MU league group.
type LeagueReward struct {
	Case1    int     `json:"case1"`
	Case2    int     `json:"case2"`
	HasBadge bool    `json:"hasBadge"`
	Money    float64 `json:"money"`
}

// MUSponsorship is a country's sponsorship of a military unit, offered
// through a law.
type MUSponsorship struct {
	ID         string `json:"_id"`
	AcceptedAt Time   `json:"acceptedAt"`
	// AcceptedBy is the user who accepted the offer for the unit.
	AcceptedBy string `json:"acceptedBy"`
	CountryID  string `json:"country"`
	CreatedAt  Time   `json:"createdAt"`
	// EndReason explains why an ended sponsorship stopped, such as
	// "citizenship" or "abandoned".
	EndReason      string   `json:"endReason"`
	EndedAt        Time     `json:"endedAt"`
	KickedUserIDs  []string `json:"kickedUsers"`
	LawID          string   `json:"law"`
	Live           bool     `json:"live"`
	MUID           string   `json:"mu"`
	OfferExpiresAt Time     `json:"offerExpiresAt"`
	OfferedAt      Time     `json:"offeredAt"`
	// Status is "offered", "active", or "ended".
	Status    string `json:"status"`
	UpdatedAt Time   `json:"updatedAt"`
}

// MUSponsorshipsByCountry is a country's military unit sponsorships.
type MUSponsorshipsByCountry struct {
	Active  []MUSponsorship `json:"active"`
	Offered []MUSponsorship `json:"offered"`
	Past    []MUSponsorship `json:"past"`
	// MUs describes the units in the sponsorships, keyed by unit ID.
	MUs map[string]MUSponsorshipUnit `json:"mus"`
}

// MUSponsorshipPage is one page of a country's active sponsorships.
type MUSponsorshipPage struct {
	CursorPage[MUSponsorship]
	// MUs describes the units on this page, keyed by unit ID.
	MUs map[string]MUSponsorshipUnit `json:"mus"`
}

// MUSponsorshipUnit is a summary of a sponsored military unit.
type MUSponsorshipUnit struct {
	HQLevel      int    `json:"hqLevel"`
	MembersCount int    `json:"membersCount"`
	Name         string `json:"name"`
}

// MUSponsorshipCandidate is a military unit a country could sponsor.
type MUSponsorshipCandidate struct {
	ID           string `json:"_id"`
	MembersCount int    `json:"membersCount"`
	Name         string `json:"name"`
}

// StateArmyOverview is a country's state army spending over a period.
type StateArmyOverview struct {
	// ByDay totals spending per day and kind.
	ByDay []StateArmyExpenseTotal `json:"byDay"`
	// ByKind totals spending per kind.
	ByKind []StateArmyExpenseTotal `json:"byKind"`
	// ByUnit totals spending per military unit and kind.
	ByUnit []StateArmyExpenseTotal `json:"byUnit"`
}

// StateArmyExpenseTotal is a total of state army spending. Day and MUID are
// set only in per-day and per-unit totals.
type StateArmyExpenseTotal struct {
	Day  Date                 `json:"day"`
	Kind StateArmyExpenseKind `json:"kind"`
	MUID string               `json:"mu"`
	// Money is the money spent, and Quantity the items given.
	Money    float64 `json:"money"`
	Quantity int     `json:"quantity"`
}

// StateArmyExpense is one state army expense: oil sent to a military unit's
// headquarters, a mercenary contract payment, or items given under a
// sponsorship law.
type StateArmyExpense struct {
	ID        string               `json:"_id"`
	At        Time                 `json:"at"`
	CountryID string               `json:"country"`
	CreatedAt Time                 `json:"createdAt"`
	Day       Date                 `json:"day"`
	ItemCode  string               `json:"itemCode"`
	Kind      StateArmyExpenseKind `json:"kind"`
	MUID      string               `json:"mu"`
	Money     float64              `json:"money"`
	Quantity  int                  `json:"quantity"`
	// RefID is the contract or law the expense came from.
	RefID     string `json:"refId"`
	UpdatedAt Time   `json:"updatedAt"`
}

// Tournament is an MU or country tournament and its bracket.
type Tournament struct {
	ID          string `json:"_id"`
	ActiveRound int    `json:"activeRound"`
	// AutoQualifyFirstRound lists the teams that skip the first round.
	AutoQualifyFirstRound []string `json:"autoQualify1stRound"`
	CreatedAt             Time     `json:"createdAt"`
	Description           string   `json:"description"`
	IsActive              bool     `json:"isActive"`
	// MaxRarity is the rarest equipment allowed.
	MaxRarity  string                  `json:"maxRarity"`
	Name       string                  `json:"name"`
	Registered TournamentRegistrations `json:"registered"`
	// Rounds holds the bracket, keyed by round number.
	Rounds      map[string]TournamentRound `json:"rounds"`
	RoundsCount int                        `json:"roundsCount"`
	// SkillKey is the skill the tournament's rounds limit.
	SkillKey  string `json:"skillKey"`
	StartAt   Time   `json:"startAt"`
	Status    string `json:"status"`
	TeamCount int    `json:"teamCount"`
	// TeamSize is the number of players per team in MU tournaments; country
	// tournaments report a fractional average.
	TeamSize float64 `json:"teamSize"`
	// Type is "mu" or "country".
	Type         string `json:"type"`
	UpdatedAt    Time   `json:"updatedAt"`
	WinnerTeamID string `json:"winnerTournamentTeam"`
}

// TournamentRegistrations lists who registered for a tournament.
type TournamentRegistrations struct {
	CountryIDs []string `json:"countries"`
	MUIDs      []string `json:"mus"`
	UserIDs    []string `json:"users"`
}

// TournamentRound is one round of a tournament bracket.
type TournamentRound struct {
	// Cases is the number of cases awarded in the round.
	Cases                int               `json:"cases"`
	IsQualificationRound bool              `json:"isQualificationRound"`
	Matches              []TournamentMatch `json:"matches"`
	RoundNumber          int               `json:"roundNumber"`
	// SkillValue is the round's limit on the tournament's skill, or nil when
	// the round has none.
	SkillValue *float64 `json:"skillValue"`
}

// TournamentMatch is a battle between two tournament teams.
type TournamentMatch struct {
	AttackerTeamID       string `json:"attacker"`
	BattleID             string `json:"battle"`
	DefenderTeamID       string `json:"defender"`
	IsQualificationRound bool   `json:"isQualificationRound"`
	MatchIndex           int    `json:"matchIndex"`
	// PossibleAttackerTeamIDs and PossibleDefenderTeamIDs list the teams that
	// could still fill each side before earlier matches are decided.
	PossibleAttackerTeamIDs []string `json:"possibleAttackerTeamIds"`
	PossibleDefenderTeamIDs []string `json:"possibleDefenderTeamIds"`
	// Predecessors holds the previous round's match indexes whose winners
	// fill each side.
	Predecessors *TournamentMatchPredecessors `json:"predecessors"`
	// WonBy is "attacker" or "defender" once the match is decided.
	WonBy string `json:"wonBy"`
}

// TournamentMatchPredecessors are the match indexes feeding a match.
type TournamentMatchPredecessors struct {
	Attacker int `json:"attacker"`
	Defender int `json:"defender"`
}

// TournamentTeam is a team of users representing military units or
// countries in a tournament.
type TournamentTeam struct {
	ID          string   `json:"_id"`
	ColorScheme string   `json:"colorScheme"`
	CountryIDs  []string `json:"countries"`
	CreatedAt   Time     `json:"createdAt"`
	// DamageByEntity holds the damage dealt for the team, keyed by user,
	// military unit, or country ID.
	DamageByEntity map[string]float64 `json:"damageByEntity"`
	EstimatedUsers int                `json:"estimatedUsers"`
	MUIDs          []string           `json:"mus"`
	// Number is the team's number within its tournament.
	Number         int      `json:"number"`
	ParticipantIDs []string `json:"participants"`
	// Status is "active", "eliminated", or "winner".
	Status       string   `json:"status"`
	TotalDamage  float64  `json:"totalDamage"`
	TournamentID string   `json:"tournament"`
	UpdatedAt    Time     `json:"updatedAt"`
	UserIDs      []string `json:"users"`
}

// BattleMercenaryDamage is the damage that mercenaries under private
// contracts dealt for each side of a battle round.
type BattleMercenaryDamage struct {
	Attacker float64 `json:"attacker"`
	Defender float64 `json:"defender"`
}

// ContributionPage is one numbered page of contributions.
type ContributionPage struct {
	Count int `json:"count"`
	Limit int `json:"limit"`
	// NextPage is the next page number, or 0 on the last page.
	NextPage int            `json:"nextPage"`
	Page     int            `json:"page"`
	Pages    int            `json:"pages"`
	Results  []Contribution `json:"results"`
}

// Contribution is health or energy a user spent on a country's unrest or a
// region's resistance.
type Contribution struct {
	ID string `json:"_id"`
	// Amount is the points contributed; unrest contributions can be negative.
	Amount    float64 `json:"amount"`
	CountryID string  `json:"country"`
	CreatedAt Time    `json:"createdAt"`
	RegionID  string  `json:"region"`
	// Resource is "health" or "energy".
	Resource string `json:"resource"`
	// Type is "unrest" or "resistance".
	Type      string `json:"type"`
	UpdatedAt Time   `json:"updatedAt"`
	UserID    string `json:"user"`
}

// RankingSnapshot is an entity's place in a ranking on one day.
type RankingSnapshot struct {
	// Count is the number of entries in the ranking that day.
	Count int     `json:"count"`
	Day   Date    `json:"day"`
	Rank  int     `json:"rank"`
	Value float64 `json:"value"`
}

// RankingSnapshotRecentRanks is an entity's daily ranks over recent days.
type RankingSnapshotRecentRanks struct {
	// Days lists the dates covered, oldest first.
	Days []Date `json:"days"`
	// Ranks holds, per ranking type, the entity's rank on each day in Days;
	// an entry is nil when there is no rank for that day.
	Ranks map[string][]*int `json:"ranks"`
}
