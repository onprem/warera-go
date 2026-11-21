package warera

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// These response models cover procedures whose contracts are not described by
// the upstream OpenAPI document. They are maintained here from observed API
// responses, endpoint behavior, and recorded fixtures.

// BattleLootSummary is a player's loot summary for a battle.
type BattleLootSummary struct {
	ID                     string               `json:"_id"`
	BattleID               string               `json:"battle"`
	Case1Count             int                  `json:"case1Count"`
	Case2Count             int                  `json:"case2Count"`
	CreatedAt              Time                 `json:"createdAt"`
	Hits                   int                  `json:"hits"`
	PoolLoot               []BattleLootPoolItem `json:"poolLoot"`
	TotalDamage            float64              `json:"totalDmg"`
	TotalMoneyFromBounty   float64              `json:"totalMoneyFromBounty"`
	TotalMoneyFromContract float64              `json:"totalMoneyFromContract"`
	UpdatedAt              Time                 `json:"updatedAt"`
	UserID                 string               `json:"user"`
}

// RegionActiveBattleValue represents the two shapes returned for a region's
// active battle: a battle ID or an expanded battle object.
type RegionActiveBattleValue struct {
	BattleID string
	Details  *RegionActiveBattle
}

// UnmarshalJSON accepts an active battle as an ID, an expanded object, or null.
func (b *RegionActiveBattleValue) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	*b = RegionActiveBattleValue{}
	if len(data) == 0 {
		return fmt.Errorf("decode active battle: empty JSON value")
	}
	switch data[0] {
	case '"':
		return json.Unmarshal(data, &b.BattleID)
	case '{':
		var battle RegionActiveBattle
		if err := json.Unmarshal(data, &battle); err != nil {
			return fmt.Errorf("decode expanded active battle: %w", err)
		}
		b.Details = &battle
		return nil
	case 'n':
		if bytes.Equal(data, []byte("null")) {
			*b = RegionActiveBattleValue{}
			return nil
		}
	}
	return fmt.Errorf("decode active battle: expected ID, object, or null")
}

// MarshalJSON writes the expanded battle when available, otherwise its ID or null.
func (b RegionActiveBattleValue) MarshalJSON() ([]byte, error) {
	if b.Details != nil {
		return json.Marshal(b.Details)
	}
	if b.BattleID != "" {
		return json.Marshal(b.BattleID)
	}
	return []byte("null"), nil
}

// BattleLootPoolItem is one item in a battle loot pool.
type BattleLootPoolItem struct {
	Item    Item   `json:"item"`
	Pool    string `json:"pool"`
	Rank    int    `json:"rank"`
	RoundID string `json:"round"`
}

// BattleRanking is the paginated response for a battle.
type BattleRanking struct {
	CursorPage[BattleRankingEntry]
	ItemCount int `json:"itemCount"`
}

// BattleRankingEntry is one battle ranking result.
type BattleRankingEntry struct {
	ID        string  `json:"_id"`
	Ranking   string  `json:"ranking"`
	UserID    string  `json:"user"`
	CountryID string  `json:"country"`
	MUID      string  `json:"mu"`
	Value     float64 `json:"value"`
	Rank      int     `json:"rank"`
	Badge     string  `json:"badge"`
	Version   int     `json:"__v"`
	CreatedAt Time    `json:"createdAt"`
	UpdatedAt Time    `json:"updatedAt"`
	LootItem  *Item   `json:"lootItem"`
}

// MercenaryContractAuction is a mercenary contract auction listing.
type MercenaryContractAuction struct {
	ID                  string                        `json:"_id"`
	BattleID            string                        `json:"battle"`
	Bids                []MercenaryContractAuctionBid `json:"bids"`
	Budget              float64                       `json:"budget"`
	CountryID           string                        `json:"country"`
	CreatedAt           Time                          `json:"createdAt"`
	CreatedByID         string                        `json:"createdBy"`
	CurrentPayout       float64                       `json:"currentPayout"`
	CurrentPerK         float64                       `json:"currentPerK"`
	CurrentWinnerID     *string                       `json:"currentWinner"`
	CurrentWinnerUserID *string                       `json:"currentWinnerUser"`
	Duration            float64                       `json:"duration"`
	ExpiresAt           Time                          `json:"expiresAt"`
	ForCountryID        string                        `json:"forCountry"`
	ForCountrySide      string                        `json:"forCountrySide"`
	InitialPerK         float64                       `json:"initialPerK"`
	MinimumDamage       float64                       `json:"minimumDamage"`
	ProfessionalsOnly   bool                          `json:"professionalsOnly"`
	RoundID             string                        `json:"round"`
	RoundNumber         int                           `json:"roundNumber"`
	Status              MercenaryAuctionStatus        `json:"status"`
	UpdatedAt           Time                          `json:"updatedAt"`
}

// MercenaryContractAuctionBid is a bid in a mercenary contract auction.
type MercenaryContractAuctionBid struct {
	BidAt  Time    `json:"bidAt"`
	MUID   string  `json:"mu"`
	Payout float64 `json:"payout"`
	PerK   float64 `json:"perK"`
	UserID string  `json:"user"`
}

// TransactionListItem is one transaction history entry.
type TransactionListItem struct {
	ID              string          `json:"_id"`
	BuyerID         string          `json:"buyerId"`
	CreatedAt       Time            `json:"createdAt"`
	Item            *Item           `json:"item"`
	ItemCode        string          `json:"itemCode"`
	Money           float64         `json:"money"`
	OfferCreatedAt  Time            `json:"offerCreatedAt"`
	Quantity        int             `json:"quantity"`
	SellerCountryID string          `json:"sellerCountryId"`
	SellerID        string          `json:"sellerId"`
	SellerMUID      string          `json:"sellerMuId"`
	TransactionType TransactionType `json:"transactionType"`
	UpdatedAt       Time            `json:"updatedAt"`
}

// Item is a weapon or equipment item, as attached to a transaction or awarded
// as battle loot.
type Item struct {
	ID                string     `json:"_id"`
	Code              string     `json:"code"`
	LastAcquisitionAt Time       `json:"lastAcquisitionAt"`
	MaxState          float64    `json:"maxState"`
	Quantity          int        `json:"quantity"`
	Skills            ItemSkills `json:"skills"`
	State             float64    `json:"state"`
	Type              string     `json:"type"`
}

// ItemSkills contains an item's skill modifiers.
type ItemSkills struct {
	Armor           float64 `json:"armor"`
	Attack          float64 `json:"attack"`
	CriticalChance  float64 `json:"criticalChance"`
	CriticalDamages float64 `json:"criticalDamages"`
	Dodge           float64 `json:"dodge"`
	Precision       float64 `json:"precision"`
	// Production is the employee production bonus percent of a tool.
	Production float64 `json:"production"`
}

// RankingPage is one page of ranking entries from ranking.getRankingPage.
type RankingPage struct {
	CursorPage[RankingEntry]
	// Ranking holds the ranking's tier cutoffs and entry count. Only the first
	// page includes it, and its Items are empty.
	Ranking *Ranking `json:"ranking"`
}

// Poll is a poll attached to an article.
type Poll struct {
	ID               string  `json:"_id"`
	AllowMultiple    bool    `json:"allowMultiple"`
	ArticleID        string  `json:"article"`
	AuthorID         string  `json:"author"`
	CanVote          bool    `json:"canVote"`
	CannotVoteReason string  `json:"cannotVoteReason"`
	CreatedAt        Time    `json:"createdAt"`
	DurationHours    float64 `json:"durationHours"`
	EndsAt           Time    `json:"endsAt"`
	IsEnded          bool    `json:"isEnded"`
	// MyVote lists the option IDs the caller picked. It is empty when the
	// caller has not voted or the request is not authenticated.
	MyVote      []int        `json:"myVote"`
	Options     []PollOption `json:"options"`
	Question    string       `json:"question"`
	StartedAt   Time         `json:"startedAt"`
	TotalVoters int          `json:"totalVoters"`
	UpdatedAt   Time         `json:"updatedAt"`
}

// PollOption is one answer in a poll with its vote count.
type PollOption struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

// PollVote is one user's vote in a poll.
type PollVote struct {
	ID        string `json:"_id"`
	CreatedAt Time   `json:"createdAt"`
	Options   []int  `json:"options"`
	PollID    string `json:"poll"`
	UpdatedAt Time   `json:"updatedAt"`
	UserID    string `json:"user"`
	Version   int    `json:"version"`
}

// GameConfigItem describes one item in the game configuration, keyed by item
// code in [GameConfig.Items]. Which fields are set depends on the item type:
// equipment and tools have DynamicStats, consumables have FlatStats, and
// produced goods have ProductionPoints and ProductionNeeds.
type GameConfigItem struct {
	Code   string `json:"code"`
	Type   string `json:"type"`
	Rarity string `json:"rarity"`
	Usage  string `json:"usage"`
	// DynamicStats maps a stat to the [min, max] range an item can roll.
	DynamicStats map[string][]float64 `json:"dynamicStats"`
	// FlatStats maps a stat to the fixed value the item gives.
	FlatStats        map[string]float64 `json:"flatStats"`
	ProductionPoints float64            `json:"productionPoints"`
	// ProductionNeeds maps an input item code to the quantity one unit needs.
	ProductionNeeds map[string]float64 `json:"productionNeeds"`
	Climates        []string           `json:"climates"`
	IsConsumable    bool               `json:"isConsumable"`
	IsDeposit       bool               `json:"isDeposit"`
	IsTradable      bool               `json:"isTradable"`
	SkinSlot        string             `json:"skinSlot"`
	IconImg         string             `json:"iconImg"`
	// IconComponent is web-client rendering metadata.
	IconComponent json.RawMessage `json:"IconComponent"`
}

// GameConfigBadge describes one badge in the game configuration, keyed by
// badge code in [GameConfig.Badge].
type GameConfigBadge struct {
	Availabilities       []string `json:"availabilities"`
	CooldownDays         float64  `json:"cooldownDays"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset"`
	Reward               float64  `json:"reward"`
	UniqueMetadataKey    string   `json:"uniqueMetadataKey"`
}

// UserSkillProfile is a player's skills, as returned by user.getUserSkills.
type UserSkillProfile struct {
	ID       string       `json:"_id"`
	Dates    UserDates    `json:"dates"`
	Leveling UserLeveling `json:"leveling"`
	Skills   UserSkills   `json:"skills"`
}

// UserListItem is a citizen in a country's user listing.
type UserListItem struct {
	ID        string `json:"_id"`
	CreatedAt Time   `json:"createdAt"`
}
