package warera

// Response models for military unit members, leagues, sponsorships, state army
// spending, tournaments, contributions, and ranking history. WarEra's published
// OpenAPI document does not describe these procedures; the models follow
// recorded responses.

// MUMember is a member of a military unit and their damage and help in it.
type MUMember struct {
	ID                  string  `json:"_id,omitempty"`
	CreatedAt           Time    `json:"createdAt,omitempty"`
	MUID                string  `json:"mu,omitempty"`
	MonthlyDamagesCount float64 `json:"monthlyDamagesCount,omitempty"`
	MonthlyHelpCount    int     `json:"monthlyHelpCount,omitempty"`
	TotalDamagesCount   float64 `json:"totalDamagesCount,omitempty"`
	TotalHelpCount      int     `json:"totalHelpCount,omitempty"`
	UpdatedAt           Time    `json:"updatedAt,omitempty"`
	UserID              string  `json:"user,omitempty"`
	WeeklyDamagesCount  float64 `json:"weeklyDamagesCount,omitempty"`
	WeeklyHelpCount     int     `json:"weeklyHelpCount,omitempty"`
}

// MULeague is a military unit's group in a week of the MU league.
type MULeague struct {
	// YearWeek identifies the league week, such as "2026-47".
	YearWeek   string `json:"yearWeek,omitempty"`
	WeekEndsAt Time   `json:"weekEndsAt,omitempty"`
	// League and Me are nil when the unit is not in a league this week.
	League *MULeagueGroup    `json:"league,omitempty"`
	Me     *MULeagueStanding `json:"me,omitempty"`
	Rows   []MULeagueRow     `json:"rows,omitempty"`
	// Shares is how the unit's reward would be split among the users who
	// fought for it this week.
	Shares []MULeagueShare `json:"shares,omitempty"`
	// PlacePots is the reward for each place in the group; places that earn
	// nothing are nil.
	PlacePots []*LeagueReward `json:"placePots,omitempty"`
}

// MULeagueGroup is a group of military units competing in one division.
type MULeagueGroup struct {
	ID           string  `json:"_id,omitempty"`
	DemoteCount  int     `json:"demoteCount,omitempty"`
	GroupIndex   int     `json:"groupIndex,omitempty"`
	GroupsCount  int     `json:"groupsCount,omitempty"`
	MembersCount int     `json:"membersCount,omitempty"`
	PromoteCount int     `json:"promoteCount,omitempty"`
	StrengthAvg  float64 `json:"strengthAvg,omitempty"`
}

// MULeagueStanding is a military unit's place in its league group.
type MULeagueStanding struct {
	Division int `json:"division,omitempty"`
	// Move is "up", "stay", or "down": where the unit would move next week.
	Move       string  `json:"move,omitempty"`
	Multiplier float64 `json:"multiplier,omitempty"`
	Points     float64 `json:"points,omitempty"`
	Rank       int     `json:"rank,omitempty"`
}

// MULeagueRow is one military unit in a league group's standings.
type MULeagueRow struct {
	Division   int     `json:"division,omitempty"`
	IsEligible bool    `json:"isEligible,omitempty"`
	Move       string  `json:"move,omitempty"`
	MUID       string  `json:"mu,omitempty"`
	Points     float64 `json:"points,omitempty"`
	// Pot is the reward for the unit's current place, or nil when it earns
	// nothing.
	Pot   *LeagueReward `json:"pot,omitempty"`
	Rank  int           `json:"rank,omitempty"`
	Seats int           `json:"seats,omitempty"`
}

// MULeagueShare is one user's share of a military unit's league reward.
type MULeagueShare struct {
	// IsMember reports whether the user is still in the unit.
	IsMember bool                `json:"isMember,omitempty"`
	Points   float64             `json:"points,omitempty"`
	Reward   MULeagueShareReward `json:"reward,omitempty"`
	// Share is the user's fraction of the unit's points, from 0 to 1.
	Share  float64 `json:"share,omitempty"`
	UserID string  `json:"user,omitempty"`
}

// MULeagueShareReward is the money and items a user would receive.
type MULeagueShareReward struct {
	// Items holds item counts keyed by item code.
	Items map[string]int `json:"items,omitempty"`
	Money float64        `json:"money,omitempty"`
}

// MULeagueDivisionCounts is the number of military units in each MU league
// division in a week.
type MULeagueDivisionCounts struct {
	// MUs holds the unit count of each division, indexed by division number
	// from 0, the lowest.
	MUs      []int  `json:"mus,omitempty"`
	YearWeek string `json:"yearWeek,omitempty"`
}

// League is a user's group in a week of the player league.
type League struct {
	// YearWeek identifies the league week, such as "2026-47".
	YearWeek   string `json:"yearWeek,omitempty"`
	WeekEndsAt Time   `json:"weekEndsAt,omitempty"`
	// League and Me are nil when the user is not in a league this week.
	League *LeagueGroup    `json:"league,omitempty"`
	Me     *LeagueStanding `json:"me,omitempty"`
	Rows   []LeagueRow     `json:"rows,omitempty"`
	// PlaceRewards is the reward for each place in the group; places that earn
	// nothing are nil.
	PlaceRewards []*LeagueReward `json:"placeRewards,omitempty"`
}

// LeagueGroup is a group of players competing in a league week.
type LeagueGroup struct {
	ID         string `json:"_id,omitempty"`
	GroupIndex int    `json:"groupIndex,omitempty"`
	// Kind is "regular", or "late" for groups of players who joined after
	// the week started.
	Kind         string  `json:"kind,omitempty"`
	LevelAvg     float64 `json:"levelAvg,omitempty"`
	MembersCount int     `json:"membersCount,omitempty"`
	Multiplier   float64 `json:"multiplier,omitempty"`
	RatingAvg    float64 `json:"ratingAvg,omitempty"`
}

// LeagueStanding is a user's place in their league group.
type LeagueStanding struct {
	Points float64 `json:"points,omitempty"`
	Rank   int     `json:"rank,omitempty"`
}

// LeagueRow is one user in a league group's standings.
type LeagueRow struct {
	IsEligible bool    `json:"isEligible,omitempty"`
	Points     float64 `json:"points,omitempty"`
	Rank       int     `json:"rank,omitempty"`
	UserID     string  `json:"user,omitempty"`
}

// LeagueReward is the reward for a place in a player or MU league group.
type LeagueReward struct {
	Case1    int     `json:"case1,omitempty"`
	Case2    int     `json:"case2,omitempty"`
	HasBadge bool    `json:"hasBadge,omitempty"`
	Money    float64 `json:"money,omitempty"`
}

// MUSponsorship is a country's sponsorship of a military unit, offered
// through a law.
type MUSponsorship struct {
	ID         string `json:"_id,omitempty"`
	AcceptedAt Time   `json:"acceptedAt,omitempty"`
	// AcceptedBy is the user who accepted the offer for the unit.
	AcceptedBy string `json:"acceptedBy,omitempty"`
	CountryID  string `json:"country,omitempty"`
	CreatedAt  Time   `json:"createdAt,omitempty"`
	// EndReason explains why an ended sponsorship stopped, such as
	// "citizenship" or "abandoned".
	EndReason      string   `json:"endReason,omitempty"`
	EndedAt        Time     `json:"endedAt,omitempty"`
	KickedUserIDs  []string `json:"kickedUsers,omitempty"`
	LawID          string   `json:"law,omitempty"`
	Live           bool     `json:"live,omitempty"`
	MUID           string   `json:"mu,omitempty"`
	OfferExpiresAt Time     `json:"offerExpiresAt,omitempty"`
	OfferedAt      Time     `json:"offeredAt,omitempty"`
	// Status is "offered", "active", or "ended".
	Status    string `json:"status,omitempty"`
	UpdatedAt Time   `json:"updatedAt,omitempty"`
}

// MUSponsorshipsByCountry is a country's military unit sponsorships.
type MUSponsorshipsByCountry struct {
	Active  []MUSponsorship `json:"active,omitempty"`
	Offered []MUSponsorship `json:"offered,omitempty"`
	Past    []MUSponsorship `json:"past,omitempty"`
	// MUs describes the units in the sponsorships, keyed by unit ID.
	MUs map[string]MUSponsorshipUnit `json:"mus,omitempty"`
}

// MUSponsorshipPage is one page of a country's active sponsorships.
type MUSponsorshipPage struct {
	CursorPage[MUSponsorship]
	// MUs describes the units on this page, keyed by unit ID.
	MUs map[string]MUSponsorshipUnit `json:"mus,omitempty"`
}

// MUSponsorshipUnit is a summary of a sponsored military unit.
type MUSponsorshipUnit struct {
	HQLevel      int    `json:"hqLevel,omitempty"`
	MembersCount int    `json:"membersCount,omitempty"`
	Name         string `json:"name,omitempty"`
}

// MUSponsorshipCandidate is a military unit a country could sponsor.
type MUSponsorshipCandidate struct {
	ID           string `json:"_id,omitempty"`
	MembersCount int    `json:"membersCount,omitempty"`
	Name         string `json:"name,omitempty"`
}

// StateArmyOverview is a country's state army spending over a period.
type StateArmyOverview struct {
	// ByDay totals spending per day and kind.
	ByDay []StateArmyExpenseTotal `json:"byDay,omitempty"`
	// ByKind totals spending per kind.
	ByKind []StateArmyExpenseTotal `json:"byKind,omitempty"`
	// ByUnit totals spending per military unit and kind.
	ByUnit []StateArmyExpenseTotal `json:"byUnit,omitempty"`
}

// StateArmyExpenseTotal is a total of state army spending. Day and MUID are
// set only in per-day and per-unit totals.
type StateArmyExpenseTotal struct {
	Day  Date                 `json:"day,omitempty"`
	Kind StateArmyExpenseKind `json:"kind,omitempty"`
	MUID string               `json:"mu,omitempty"`
	// Money is the money spent, and Quantity the items given.
	Money    float64 `json:"money,omitempty"`
	Quantity int     `json:"quantity,omitempty"`
}

// StateArmyExpense is one state army expense: oil sent to a military unit's
// headquarters, a mercenary contract payment, or items given under a
// sponsorship law.
type StateArmyExpense struct {
	ID        string               `json:"_id,omitempty"`
	At        Time                 `json:"at,omitempty"`
	CountryID string               `json:"country,omitempty"`
	CreatedAt Time                 `json:"createdAt,omitempty"`
	Day       Date                 `json:"day,omitempty"`
	ItemCode  string               `json:"itemCode,omitempty"`
	Kind      StateArmyExpenseKind `json:"kind,omitempty"`
	MUID      string               `json:"mu,omitempty"`
	Money     float64              `json:"money,omitempty"`
	Quantity  int                  `json:"quantity,omitempty"`
	// RefID is the contract or law the expense came from.
	RefID     string `json:"refId,omitempty"`
	UpdatedAt Time   `json:"updatedAt,omitempty"`
}

// Tournament is an MU or country tournament and its bracket.
type Tournament struct {
	ID          string `json:"_id,omitempty"`
	ActiveRound int    `json:"activeRound,omitempty"`
	// AutoQualifyFirstRound lists the teams that skip the first round.
	AutoQualifyFirstRound []string `json:"autoQualify1stRound,omitempty"`
	CreatedAt             Time     `json:"createdAt,omitempty"`
	Description           string   `json:"description,omitempty"`
	IsActive              bool     `json:"isActive,omitempty"`
	// MaxRarity is the rarest equipment allowed.
	MaxRarity  string                  `json:"maxRarity,omitempty"`
	Name       string                  `json:"name,omitempty"`
	Registered TournamentRegistrations `json:"registered,omitempty"`
	// Rounds holds the bracket, keyed by round number.
	Rounds      map[string]TournamentRound `json:"rounds,omitempty"`
	RoundsCount int                        `json:"roundsCount,omitempty"`
	// SkillKey is the skill the tournament's rounds limit.
	SkillKey  string `json:"skillKey,omitempty"`
	StartAt   Time   `json:"startAt,omitempty"`
	Status    string `json:"status,omitempty"`
	TeamCount int    `json:"teamCount,omitempty"`
	// TeamSize is the number of players per team in MU tournaments; country
	// tournaments report a fractional average.
	TeamSize float64 `json:"teamSize,omitempty"`
	// Type is "mu" or "country".
	Type         string `json:"type,omitempty"`
	UpdatedAt    Time   `json:"updatedAt,omitempty"`
	WinnerTeamID string `json:"winnerTournamentTeam,omitempty"`
}

// TournamentRegistrations lists who registered for a tournament.
type TournamentRegistrations struct {
	CountryIDs []string `json:"countries,omitempty"`
	MUIDs      []string `json:"mus,omitempty"`
	UserIDs    []string `json:"users,omitempty"`
}

// TournamentRound is one round of a tournament bracket.
type TournamentRound struct {
	// Cases is the number of cases awarded in the round.
	Cases                int               `json:"cases,omitempty"`
	IsQualificationRound bool              `json:"isQualificationRound,omitempty"`
	Matches              []TournamentMatch `json:"matches,omitempty"`
	RoundNumber          int               `json:"roundNumber,omitempty"`
	// SkillValue is the round's limit on the tournament's skill, or nil when
	// the round has none.
	SkillValue *float64 `json:"skillValue,omitempty"`
}

// TournamentMatch is a battle between two tournament teams.
type TournamentMatch struct {
	AttackerTeamID       string `json:"attacker,omitempty"`
	BattleID             string `json:"battle,omitempty"`
	DefenderTeamID       string `json:"defender,omitempty"`
	IsQualificationRound bool   `json:"isQualificationRound,omitempty"`
	MatchIndex           int    `json:"matchIndex,omitempty"`
	// PossibleAttackerTeamIDs and PossibleDefenderTeamIDs list the teams that
	// could still fill each side before earlier matches are decided.
	PossibleAttackerTeamIDs []string `json:"possibleAttackerTeamIds,omitempty"`
	PossibleDefenderTeamIDs []string `json:"possibleDefenderTeamIds,omitempty"`
	// Predecessors holds the previous round's match indexes whose winners
	// fill each side.
	Predecessors *TournamentMatchPredecessors `json:"predecessors,omitempty"`
	// WonBy is "attacker" or "defender" once the match is decided.
	WonBy string `json:"wonBy,omitempty"`
}

// TournamentMatchPredecessors are the match indexes feeding a match.
type TournamentMatchPredecessors struct {
	Attacker int `json:"attacker,omitempty"`
	Defender int `json:"defender,omitempty"`
}

// TournamentTeam is a team of users representing military units or
// countries in a tournament.
type TournamentTeam struct {
	ID          string   `json:"_id,omitempty"`
	ColorScheme string   `json:"colorScheme,omitempty"`
	CountryIDs  []string `json:"countries,omitempty"`
	CreatedAt   Time     `json:"createdAt,omitempty"`
	// DamageByEntity holds the damage dealt for the team, keyed by user,
	// military unit, or country ID.
	DamageByEntity map[string]float64 `json:"damageByEntity,omitempty"`
	EstimatedUsers int                `json:"estimatedUsers,omitempty"`
	MUIDs          []string           `json:"mus,omitempty"`
	// Number is the team's number within its tournament.
	Number         int      `json:"number,omitempty"`
	ParticipantIDs []string `json:"participants,omitempty"`
	// Status is "active", "eliminated", or "winner".
	Status       string   `json:"status,omitempty"`
	TotalDamage  float64  `json:"totalDamage,omitempty"`
	TournamentID string   `json:"tournament,omitempty"`
	UpdatedAt    Time     `json:"updatedAt,omitempty"`
	UserIDs      []string `json:"users,omitempty"`
}

// BattleMercenaryDamage is the damage that mercenaries under private
// contracts dealt for each side of a battle round.
type BattleMercenaryDamage struct {
	Attacker float64 `json:"attacker,omitempty"`
	Defender float64 `json:"defender,omitempty"`
}

// ContributionPage is one numbered page of contributions.
type ContributionPage struct {
	Count int `json:"count,omitempty"`
	Limit int `json:"limit,omitempty"`
	// NextPage is the next page number, or 0 on the last page.
	NextPage int            `json:"nextPage,omitempty"`
	Page     int            `json:"page,omitempty"`
	Pages    int            `json:"pages,omitempty"`
	Results  []Contribution `json:"results,omitempty"`
}

// Contribution is health or energy a user spent on a country's unrest or a
// region's resistance.
type Contribution struct {
	ID string `json:"_id,omitempty"`
	// Amount is the points contributed; unrest contributions can be negative.
	Amount    float64 `json:"amount,omitempty"`
	CountryID string  `json:"country,omitempty"`
	CreatedAt Time    `json:"createdAt,omitempty"`
	RegionID  string  `json:"region,omitempty"`
	// Resource is "health" or "energy".
	Resource string `json:"resource,omitempty"`
	// Type is "unrest" or "resistance".
	Type      string `json:"type,omitempty"`
	UpdatedAt Time   `json:"updatedAt,omitempty"`
	UserID    string `json:"user,omitempty"`
}

// RankingSnapshot is an entity's place in a ranking on one day.
type RankingSnapshot struct {
	// Count is the number of entries in the ranking that day.
	Count int     `json:"count,omitempty"`
	Day   Date    `json:"day,omitempty"`
	Rank  int     `json:"rank,omitempty"`
	Value float64 `json:"value,omitempty"`
}

// RankingSnapshotRecentRanks is an entity's daily ranks over recent days.
type RankingSnapshotRecentRanks struct {
	// Days lists the dates covered, oldest first.
	Days []Date `json:"days,omitempty"`
	// Ranks holds, per ranking type, the entity's rank on each day in Days;
	// an entry is nil when there is no rank for that day.
	Ranks map[string][]*int `json:"ranks,omitempty"`
}
