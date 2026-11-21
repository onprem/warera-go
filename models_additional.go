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
	ID                     string               `json:"_id,omitempty"`
	BattleID               string               `json:"battle,omitempty"`
	Case1Count             int                  `json:"case1Count,omitempty"`
	Case2Count             int                  `json:"case2Count,omitempty"`
	CreatedAt              Time                 `json:"createdAt,omitempty"`
	Hits                   int                  `json:"hits,omitempty"`
	PoolLoot               []BattleLootPoolItem `json:"poolLoot,omitempty"`
	TotalDamage            float64              `json:"totalDmg,omitempty"`
	TotalMoneyFromBounty   float64              `json:"totalMoneyFromBounty,omitempty"`
	TotalMoneyFromContract float64              `json:"totalMoneyFromContract,omitempty"`
	UpdatedAt              Time                 `json:"updatedAt,omitempty"`
	UserID                 string               `json:"user,omitempty"`
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
	Item  Item   `json:"item"`
	Pool  string `json:"pool"`
	Rank  int    `json:"rank"`
	Round string `json:"round,omitempty"`
}

// BattleRanking is the paginated response for a battle.
type BattleRanking struct {
	Items      []BattleRankingEntry `json:"items,omitempty"`
	NextCursor string               `json:"nextCursor,omitempty"`
	ItemCount  int                  `json:"itemCount,omitempty"`
}

// BattleRankingEntry is one battle ranking result.
type BattleRankingEntry struct {
	ID        string  `json:"_id,omitempty"`
	Ranking   string  `json:"ranking,omitempty"`
	UserID    string  `json:"user,omitempty"`
	CountryID string  `json:"country,omitempty"`
	MUID      string  `json:"mu,omitempty"`
	Value     float64 `json:"value,omitempty"`
	Rank      int     `json:"rank,omitempty"`
	Badge     string  `json:"badge,omitempty"`
	Version   int     `json:"__v,omitempty"`
	CreatedAt Time    `json:"createdAt,omitempty"`
	UpdatedAt Time    `json:"updatedAt,omitempty"`
	LootItem  *Item   `json:"lootItem,omitempty"`
}

// MercenaryContractAuction is a mercenary contract auction listing.
type MercenaryContractAuction struct {
	ID                string                        `json:"_id,omitempty"`
	BattleID          string                        `json:"battle,omitempty"`
	Bids              []MercenaryContractAuctionBid `json:"bids,omitempty"`
	Budget            float64                       `json:"budget,omitempty"`
	CountryID         string                        `json:"country,omitempty"`
	CreatedAt         Time                          `json:"createdAt,omitempty"`
	CreatedBy         string                        `json:"createdBy,omitempty"`
	CurrentPayout     float64                       `json:"currentPayout,omitempty"`
	CurrentPerK       float64                       `json:"currentPerK,omitempty"`
	CurrentWinner     *string                       `json:"currentWinner,omitempty"`
	CurrentWinnerUser *string                       `json:"currentWinnerUser,omitempty"`
	Duration          float64                       `json:"duration,omitempty"`
	ExpiresAt         Time                          `json:"expiresAt,omitempty"`
	ForCountry        string                        `json:"forCountry,omitempty"`
	ForCountrySide    string                        `json:"forCountrySide,omitempty"`
	InitialPerK       float64                       `json:"initialPerK,omitempty"`
	MinimumDamage     float64                       `json:"minimumDamage,omitempty"`
	ProfessionalsOnly bool                          `json:"professionalsOnly,omitempty"`
	RoundID           string                        `json:"round,omitempty"`
	RoundNumber       int                           `json:"roundNumber,omitempty"`
	Status            MercenaryAuctionStatus        `json:"status,omitempty"`
	UpdatedAt         Time                          `json:"updatedAt,omitempty"`
}

// MercenaryContractAuctionBid is a bid in a mercenary contract auction.
type MercenaryContractAuctionBid struct {
	BidAt  Time    `json:"bidAt,omitempty"`
	MUID   string  `json:"mu,omitempty"`
	Payout float64 `json:"payout,omitempty"`
	PerK   float64 `json:"perK,omitempty"`
	UserID string  `json:"user,omitempty"`
}

// TransactionListItem is one transaction history entry.
type TransactionListItem struct {
	ID              string          `json:"_id,omitempty"`
	BuyerID         string          `json:"buyerId,omitempty"`
	CreatedAt       Time            `json:"createdAt,omitempty"`
	Item            *Item           `json:"item,omitempty"`
	ItemCode        string          `json:"itemCode,omitempty"`
	Money           float64         `json:"money,omitempty"`
	OfferCreatedAt  Time            `json:"offerCreatedAt,omitempty"`
	Quantity        int             `json:"quantity,omitempty"`
	SellerCountryID string          `json:"sellerCountryId,omitempty"`
	SellerID        string          `json:"sellerId,omitempty"`
	SellerMUID      string          `json:"sellerMuId,omitempty"`
	TransactionType TransactionType `json:"transactionType,omitempty"`
	UpdatedAt       Time            `json:"updatedAt,omitempty"`
}

// Item is a weapon or equipment item, as attached to a transaction or awarded
// as battle loot.
type Item struct {
	ID                string     `json:"_id,omitempty"`
	Code              string     `json:"code,omitempty"`
	LastAcquisitionAt Time       `json:"lastAcquisitionAt,omitempty"`
	MaxState          float64    `json:"maxState,omitempty"`
	Quantity          int        `json:"quantity,omitempty"`
	Skills            ItemSkills `json:"skills,omitempty"`
	State             float64    `json:"state,omitempty"`
	Type              string     `json:"type,omitempty"`
}

// ItemSkills contains an item's skill modifiers.
type ItemSkills struct {
	Armor           float64 `json:"armor,omitempty"`
	Attack          float64 `json:"attack,omitempty"`
	CriticalChance  float64 `json:"criticalChance,omitempty"`
	CriticalDamages float64 `json:"criticalDamages,omitempty"`
	Dodge           float64 `json:"dodge,omitempty"`
	Precision       float64 `json:"precision,omitempty"`
	// Production is the employee production bonus percent of a tool.
	Production float64 `json:"production,omitempty"`
}

// RankingPage is one page of ranking entries from ranking.getRankingPage.
type RankingPage struct {
	CursorPage[RankingEntry]
	// Ranking holds the ranking's tier cutoffs and entry count. Only the first
	// page includes it, and its Items are empty.
	Ranking *Ranking `json:"ranking,omitempty"`
}

// Poll is a poll attached to an article.
type Poll struct {
	ID               string  `json:"_id,omitempty"`
	AllowMultiple    bool    `json:"allowMultiple,omitempty"`
	ArticleID        string  `json:"article,omitempty"`
	AuthorID         string  `json:"author,omitempty"`
	CanVote          bool    `json:"canVote,omitempty"`
	CannotVoteReason string  `json:"cannotVoteReason,omitempty"`
	CreatedAt        Time    `json:"createdAt,omitempty"`
	DurationHours    float64 `json:"durationHours,omitempty"`
	EndsAt           Time    `json:"endsAt,omitempty"`
	IsEnded          bool    `json:"isEnded,omitempty"`
	// MyVote lists the option IDs the caller picked. It is empty when the
	// caller has not voted or the request is not authenticated.
	MyVote      []int        `json:"myVote,omitempty"`
	Options     []PollOption `json:"options,omitempty"`
	Question    string       `json:"question,omitempty"`
	StartedAt   Time         `json:"startedAt,omitempty"`
	TotalVoters int          `json:"totalVoters,omitempty"`
	UpdatedAt   Time         `json:"updatedAt,omitempty"`
}

// PollOption is one answer in a poll with its vote count.
type PollOption struct {
	ID    int    `json:"id"`
	Text  string `json:"text,omitempty"`
	Votes int    `json:"votes,omitempty"`
}

// PollVote is one user's vote in a poll.
type PollVote struct {
	ID        string `json:"_id,omitempty"`
	CreatedAt Time   `json:"createdAt,omitempty"`
	Options   []int  `json:"options,omitempty"`
	PollID    string `json:"poll,omitempty"`
	UpdatedAt Time   `json:"updatedAt,omitempty"`
	UserID    string `json:"user,omitempty"`
	Version   int    `json:"version,omitempty"`
}

// GameConfigItem describes one item in the game configuration, keyed by item
// code in [GameConfig.Items]. Which fields are set depends on the item type:
// equipment and tools have DynamicStats, consumables have FlatStats, and
// produced goods have ProductionPoints and ProductionNeeds.
type GameConfigItem struct {
	Code   string `json:"code,omitempty"`
	Type   string `json:"type,omitempty"`
	Rarity string `json:"rarity,omitempty"`
	Usage  string `json:"usage,omitempty"`
	// DynamicStats maps a stat to the [min, max] range an item can roll.
	DynamicStats map[string][]float64 `json:"dynamicStats,omitempty"`
	// FlatStats maps a stat to the fixed value the item gives.
	FlatStats        map[string]float64 `json:"flatStats,omitempty"`
	ProductionPoints float64            `json:"productionPoints,omitempty"`
	// ProductionNeeds maps an input item code to the quantity one unit needs.
	ProductionNeeds map[string]float64 `json:"productionNeeds,omitempty"`
	Climates        []string           `json:"climates,omitempty"`
	IsConsumable    bool               `json:"isConsumable,omitempty"`
	IsDeposit       bool               `json:"isDeposit,omitempty"`
	IsTradable      bool               `json:"isTradable,omitempty"`
	SkinSlot        string             `json:"skinSlot,omitempty"`
	IconImg         string             `json:"iconImg,omitempty"`
	// IconComponent is web-client rendering metadata.
	IconComponent json.RawMessage `json:"IconComponent,omitempty"`
}

// GameConfigBadge describes one badge in the game configuration, keyed by
// badge code in [GameConfig.Badge].
type GameConfigBadge struct {
	Availabilities       []string `json:"availabilities,omitempty"`
	CooldownDays         float64  `json:"cooldownDays,omitempty"`
	PreserveBetweenReset bool     `json:"preserveBetweenReset,omitempty"`
	Reward               float64  `json:"reward,omitempty"`
	UniqueMetadataKey    string   `json:"uniqueMetadataKey,omitempty"`
}

// UserSkillProfile is a player's skills, as returned by user.getUserSkills.
type UserSkillProfile struct {
	ID       string       `json:"_id,omitempty"`
	Dates    UserDates    `json:"dates,omitempty"`
	Leveling UserLeveling `json:"leveling,omitempty"`
	Skills   UserSkills   `json:"skills,omitempty"`
}

// UserListItem is a citizen in a country's user listing.
type UserListItem struct {
	ID        string `json:"_id,omitempty"`
	CreatedAt Time   `json:"createdAt,omitempty"`
}
