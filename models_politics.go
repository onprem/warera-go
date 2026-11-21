package warera

// Response models for alliances, wars, parties, elections, and country
// politics. WarEra's published OpenAPI document does not describe these
// procedures; the models follow recorded responses.

// Alliance is a group of countries that pool their strategic resources.
type Alliance struct {
	ID                 string  `json:"_id"`
	AvatarURL          string  `json:"avatarUrl"`
	AverageDevelopment float64 `json:"averageDevelopment"`
	CoreDevelopment    float64 `json:"coreDevelopment"`
	CreatedAt          Time    `json:"createdAt"`
	CurrentDevelopment float64 `json:"currentDevelopment"`
	// LeaderID is the user ID of the alliance's leader.
	LeaderID        string                 `json:"leader"`
	MapAccent       string                 `json:"mapAccent"`
	MemberCountries []AllianceMember       `json:"memberCountries"`
	Name            string                 `json:"name"`
	Rankings        map[string]RankingItem `json:"rankings"`
	Scheme          string                 `json:"scheme"`
	UpdatedAt       Time                   `json:"updatedAt"`
}

// AllianceMember is a country in an alliance.
type AllianceMember struct {
	AverageDevelopment float64 `json:"averageDevelopment"`
	CoreDevelopment    float64 `json:"coreDevelopment"`
	CountryID          string  `json:"country"`
	Suspended          bool    `json:"suspended"`
}

// StrategicResourceBonuses lists the strategic resource regions held by a
// country or its alliance and the bonuses they grant.
type StrategicResourceBonuses struct {
	// AllianceID is empty for a country outside any alliance.
	AllianceID string   `json:"allianceId"`
	CountryIDs []string `json:"countryIds"`
	// RegionIDs lists the regions producing each resource, keyed by resource
	// name such as "gold" or "rareEarths".
	RegionIDs map[string][]string          `json:"regionIds"`
	Bonuses   StrategicResourceBonusTotals `json:"bonuses"`
}

// StrategicResourceBonusTotals is the total bonus from strategic resources.
type StrategicResourceBonusTotals struct {
	DevelopmentPercent float64 `json:"developmentPercent"`
	ProductionPercent  float64 `json:"productionPercent"`
	// Skills holds skill bonuses keyed by skill name, such as "health".
	Skills map[string]float64 `json:"skills"`
}

// War is a war between two countries.
type War struct {
	ID       string  `json:"_id"`
	Attacker WarSide `json:"attacker"`
	// BattleIDs has been empty in every recorded response.
	BattleIDs []string `json:"battles"`
	CreatedAt Time     `json:"createdAt"`
	Defender  WarSide  `json:"defender"`
	IsActive  bool     `json:"isActive"`
	// PriorityCountryID is the attacker or defender country holding priority
	// until PriorityEndAt.
	PriorityCountryID string `json:"priority"`
	PriorityEndAt     Time   `json:"priorityEndAt"`
	UpdatedAt         Time   `json:"updatedAt"`
}

// WarSide is one side's results in a war.
type WarSide struct {
	CountryID       string  `json:"country"`
	Damages         float64 `json:"damages"`
	WonBattlesCount int     `json:"wonBattlesCount"`
	WonRoundsCount  int     `json:"wonRoundsCount"`
}

// Party is a political party. Member, council, and congress list entries are
// user IDs.
type Party struct {
	ID        string `json:"_id"`
	AvatarURL string `json:"avatarUrl"`
	// CongressListIDs is the party's ordered list of congress candidates.
	CongressListIDs         []string `json:"congressList"`
	CongressListUpdatedAt   Time     `json:"congressListUpdatedAt"`
	CongressListUpdatedByID string   `json:"congressListUpdatedBy"`
	CouncilMemberIDs        []string `json:"councilMembers"`
	CountryID               string   `json:"country"`
	CreatedAt               Time     `json:"createdAt"`
	Description             string   `json:"description"`
	FavoriteGovType         string   `json:"favoriteGovType"`
	LeaderID                string   `json:"leader"`
	MemberIDs               []string `json:"members"`
	Name                    string   `json:"name"`
	// PrimaryWinnerID is the user who won the party's latest primary
	// election.
	PrimaryWinnerID string `json:"primaryWinner"`
	RegionID        string `json:"region"`
	TreasurerID     string `json:"treasurer"`
	UpdatedAt       Time   `json:"updatedAt"`
}

// CongressSeat is a seat in a country's congress.
type CongressSeat struct {
	// HasLeftParty reports whether the member left the party they were
	// elected for.
	HasLeftParty bool `json:"hasLeftParty"`
	// PartyID is empty for members elected without a party.
	PartyID string `json:"party"`
	UserID  string `json:"user"`
}

// Diplomacy is a country's sworn enemy and defensive pacts.
type Diplomacy struct {
	DefensivePacts []DiplomacyPact `json:"defensivePacts"`
	// SwornEnemy is nil when the country has none.
	SwornEnemy *DiplomacySwornEnemy `json:"swornEnemy"`
}

// DiplomacyPact is a defensive pact with a partner country and the damage
// bonus it gives when fighting for that partner.
type DiplomacyPact struct {
	BonusPercent float64 `json:"bonusPercent"`
	DamagesDealt float64 `json:"damagesDealt"`
	PartnerID    string  `json:"partner"`
}

// DiplomacySwornEnemy is a country's sworn enemy and the damage bonus against
// it.
type DiplomacySwornEnemy struct {
	BonusPercent float64 `json:"bonusPercent"`
	DamagesDealt float64 `json:"damagesDealt"`
	EnemyID      string  `json:"enemy"`
}

// UnrestStatus is a country's unrest bar. A full bar starts a revolution.
type UnrestStatus struct {
	Bar                      float64 `json:"bar"`
	BarMax                   float64 `json:"barMax"`
	IsFull                   bool    `json:"isFull"`
	IsOnRevolutionCooldown   bool    `json:"isOnRevolutionCooldown"`
	LastContributionAt       Time    `json:"lastContributionAt"`
	RevolutionCooldownEndsAt Time    `json:"revolutionCooldownEndsAt"`
}

// Election is a country election (Type "president", "congress", or
// "congressList") or a party election (such as "partyLeader",
// "partyCouncil", or "partyPrimary").
type Election struct {
	ID         string              `json:"_id"`
	Candidates []ElectionCandidate `json:"candidates"`
	// CountryID is set for country elections and PartyID for party
	// elections.
	CountryID           string         `json:"country"`
	CreatedAt           Time           `json:"createdAt"`
	ElectedCandidateIDs []string       `json:"electedCandidates"`
	ElectedCount        int            `json:"electedCount"`
	IsActive            bool           `json:"isActive"`
	Lists               []ElectionList `json:"lists"`
	PartyID             string         `json:"party"`
	// Status is the election phase, such as "candidate" while candidates
	// register, or "finished".
	Status string `json:"status"`
	Type   string `json:"type"`
	// Votes holds vote counts keyed by candidate user ID.
	Votes        map[string]int `json:"votes"`
	VotesCount   int            `json:"votesCount"`
	VotesEndAt   Time           `json:"votesEndAt"`
	VotesStartAt Time           `json:"votesStartAt"`
}

// ElectionCandidate is a candidate in an election.
type ElectionCandidate struct {
	// ArticleID is the candidate's campaign article, if any.
	ArticleID string `json:"article"`
	IsElected bool   `json:"isElected"`
	PartyID   string `json:"party"`
	UserID    string `json:"user"`
	// VoteCount has been 0 in every recorded response; Election.Votes holds
	// the tally.
	VoteCount int `json:"voteCount"`
}

// ElectionList is a party's list of candidates in a congress list election.
type ElectionList struct {
	ArticleID      string   `json:"article"`
	MemberIDs      []string `json:"members"`
	PartyID        string   `json:"party"`
	RegisteredAt   Time     `json:"registeredAt"`
	RegisteredByID string   `json:"registeredBy"`
	Seats          int      `json:"seats"`
}
