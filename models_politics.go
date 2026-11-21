package warera

// Response models for alliances, wars, parties, elections, and country
// politics. WarEra's published OpenAPI document does not describe these
// procedures; the models follow recorded responses.

// Alliance is a group of countries that pool their strategic resources.
type Alliance struct {
	ID                 string  `json:"_id,omitempty"`
	AvatarURL          string  `json:"avatarUrl,omitempty"`
	AverageDevelopment float64 `json:"averageDevelopment,omitempty"`
	CoreDevelopment    float64 `json:"coreDevelopment,omitempty"`
	CreatedAt          Time    `json:"createdAt,omitempty"`
	CurrentDevelopment float64 `json:"currentDevelopment,omitempty"`
	// LeaderID is the user ID of the alliance's leader.
	LeaderID        string                 `json:"leader,omitempty"`
	MapAccent       string                 `json:"mapAccent,omitempty"`
	MemberCountries []AllianceMember       `json:"memberCountries,omitempty"`
	Name            string                 `json:"name,omitempty"`
	Rankings        map[string]RankingItem `json:"rankings,omitempty"`
	Scheme          string                 `json:"scheme,omitempty"`
	UpdatedAt       Time                   `json:"updatedAt,omitempty"`
}

// AllianceMember is a country in an alliance.
type AllianceMember struct {
	AverageDevelopment float64 `json:"averageDevelopment,omitempty"`
	CoreDevelopment    float64 `json:"coreDevelopment,omitempty"`
	CountryID          string  `json:"country,omitempty"`
	Suspended          bool    `json:"suspended,omitempty"`
}

// StrategicResourceBonuses lists the strategic resource regions held by a
// country or its alliance and the bonuses they grant.
type StrategicResourceBonuses struct {
	// AllianceID is empty for a country outside any alliance.
	AllianceID string   `json:"allianceId,omitempty"`
	CountryIDs []string `json:"countryIds,omitempty"`
	// RegionIDs lists the regions producing each resource, keyed by resource
	// name such as "gold" or "rareEarths".
	RegionIDs map[string][]string          `json:"regionIds,omitempty"`
	Bonuses   StrategicResourceBonusTotals `json:"bonuses,omitempty"`
}

// StrategicResourceBonusTotals is the total bonus from strategic resources.
type StrategicResourceBonusTotals struct {
	DevelopmentPercent float64 `json:"developmentPercent,omitempty"`
	ProductionPercent  float64 `json:"productionPercent,omitempty"`
	// Skills holds skill bonuses keyed by skill name, such as "health".
	Skills map[string]float64 `json:"skills,omitempty"`
}

// War is a war between two countries.
type War struct {
	ID       string  `json:"_id,omitempty"`
	Attacker WarSide `json:"attacker,omitempty"`
	// BattleIDs has been empty in every recorded response.
	BattleIDs []string `json:"battles,omitempty"`
	CreatedAt Time     `json:"createdAt,omitempty"`
	Defender  WarSide  `json:"defender,omitempty"`
	IsActive  bool     `json:"isActive,omitempty"`
	// PriorityCountryID is the attacker or defender country holding priority
	// until PriorityEndAt.
	PriorityCountryID string `json:"priority,omitempty"`
	PriorityEndAt     Time   `json:"priorityEndAt,omitempty"`
	UpdatedAt         Time   `json:"updatedAt,omitempty"`
}

// WarSide is one side's results in a war.
type WarSide struct {
	CountryID       string  `json:"country,omitempty"`
	Damages         float64 `json:"damages,omitempty"`
	WonBattlesCount int     `json:"wonBattlesCount,omitempty"`
	WonRoundsCount  int     `json:"wonRoundsCount,omitempty"`
}

// Party is a political party. Member, council, and congress list entries are
// user IDs.
type Party struct {
	ID        string `json:"_id,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
	// CongressListIDs is the party's ordered list of congress candidates.
	CongressListIDs         []string `json:"congressList,omitempty"`
	CongressListUpdatedAt   Time     `json:"congressListUpdatedAt,omitempty"`
	CongressListUpdatedByID string   `json:"congressListUpdatedBy,omitempty"`
	CouncilMemberIDs        []string `json:"councilMembers,omitempty"`
	CountryID               string   `json:"country,omitempty"`
	CreatedAt               Time     `json:"createdAt,omitempty"`
	Description             string   `json:"description,omitempty"`
	FavoriteGovType         string   `json:"favoriteGovType,omitempty"`
	LeaderID                string   `json:"leader,omitempty"`
	MemberIDs               []string `json:"members,omitempty"`
	Name                    string   `json:"name,omitempty"`
	// PrimaryWinnerID is the user who won the party's latest primary
	// election.
	PrimaryWinnerID string `json:"primaryWinner,omitempty"`
	RegionID        string `json:"region,omitempty"`
	TreasurerID     string `json:"treasurer,omitempty"`
	UpdatedAt       Time   `json:"updatedAt,omitempty"`
}

// CongressSeat is a seat in a country's congress.
type CongressSeat struct {
	// HasLeftParty reports whether the member left the party they were
	// elected for.
	HasLeftParty bool `json:"hasLeftParty,omitempty"`
	// PartyID is empty for members elected without a party.
	PartyID string `json:"party,omitempty"`
	UserID  string `json:"user,omitempty"`
}

// Diplomacy is a country's sworn enemy and defensive pacts.
type Diplomacy struct {
	DefensivePacts []DiplomacyPact `json:"defensivePacts,omitempty"`
	// SwornEnemy is nil when the country has none.
	SwornEnemy *DiplomacySwornEnemy `json:"swornEnemy,omitempty"`
}

// DiplomacyPact is a defensive pact with a partner country and the damage
// bonus it gives when fighting for that partner.
type DiplomacyPact struct {
	BonusPercent float64 `json:"bonusPercent,omitempty"`
	DamagesDealt float64 `json:"damagesDealt,omitempty"`
	PartnerID    string  `json:"partner,omitempty"`
}

// DiplomacySwornEnemy is a country's sworn enemy and the damage bonus against
// it.
type DiplomacySwornEnemy struct {
	BonusPercent float64 `json:"bonusPercent,omitempty"`
	DamagesDealt float64 `json:"damagesDealt,omitempty"`
	EnemyID      string  `json:"enemy,omitempty"`
}

// UnrestStatus is a country's unrest bar. A full bar starts a revolution.
type UnrestStatus struct {
	Bar                      float64 `json:"bar,omitempty"`
	BarMax                   float64 `json:"barMax,omitempty"`
	IsFull                   bool    `json:"isFull,omitempty"`
	IsOnRevolutionCooldown   bool    `json:"isOnRevolutionCooldown,omitempty"`
	LastContributionAt       Time    `json:"lastContributionAt,omitempty"`
	RevolutionCooldownEndsAt Time    `json:"revolutionCooldownEndsAt,omitempty"`
}

// Election is a country election (Type "president", "congress", or
// "congressList") or a party election (such as "partyLeader",
// "partyCouncil", or "partyPrimary").
type Election struct {
	ID         string              `json:"_id,omitempty"`
	Candidates []ElectionCandidate `json:"candidates,omitempty"`
	// CountryID is set for country elections and PartyID for party
	// elections.
	CountryID           string         `json:"country,omitempty"`
	CreatedAt           Time           `json:"createdAt,omitempty"`
	ElectedCandidateIDs []string       `json:"electedCandidates,omitempty"`
	ElectedCount        int            `json:"electedCount,omitempty"`
	IsActive            bool           `json:"isActive,omitempty"`
	Lists               []ElectionList `json:"lists,omitempty"`
	PartyID             string         `json:"party,omitempty"`
	// Status is the election phase, such as "candidate" while candidates
	// register, or "finished".
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
	// Votes holds vote counts keyed by candidate user ID.
	Votes        map[string]int `json:"votes,omitempty"`
	VotesCount   int            `json:"votesCount,omitempty"`
	VotesEndAt   Time           `json:"votesEndAt,omitempty"`
	VotesStartAt Time           `json:"votesStartAt,omitempty"`
}

// ElectionCandidate is a candidate in an election.
type ElectionCandidate struct {
	// ArticleID is the candidate's campaign article, if any.
	ArticleID string `json:"article,omitempty"`
	IsElected bool   `json:"isElected,omitempty"`
	PartyID   string `json:"party,omitempty"`
	UserID    string `json:"user,omitempty"`
	// VoteCount has been 0 in every recorded response; Election.Votes holds
	// the tally.
	VoteCount int `json:"voteCount,omitempty"`
}

// ElectionList is a party's list of candidates in a congress list election.
type ElectionList struct {
	ArticleID      string   `json:"article,omitempty"`
	MemberIDs      []string `json:"members,omitempty"`
	PartyID        string   `json:"party,omitempty"`
	RegisteredAt   Time     `json:"registeredAt,omitempty"`
	RegisteredByID string   `json:"registeredBy,omitempty"`
	Seats          int      `json:"seats,omitempty"`
}
