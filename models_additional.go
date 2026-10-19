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
	CreatedAt              string               `json:"createdAt,omitempty"`
	Hits                   int                  `json:"hits,omitempty"`
	PoolLoot               []BattleLootPoolItem `json:"poolLoot,omitempty"`
	TotalDamage            int64                `json:"totalDmg,omitempty"`
	TotalMoneyFromBounty   float64              `json:"totalMoneyFromBounty,omitempty"`
	TotalMoneyFromContract float64              `json:"totalMoneyFromContract,omitempty"`
	UpdatedAt              string               `json:"updatedAt,omitempty"`
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
	Item  RoundWeapon `json:"item"`
	Pool  string      `json:"pool"`
	Rank  int         `json:"rank"`
	Round string      `json:"round,omitempty"`
}

// BattleRanking is the paginated response for a battle.
type BattleRanking struct {
	Items      []BattleRankingEntry `json:"items,omitempty"`
	NextCursor string               `json:"nextCursor,omitempty"`
	ItemCount  int                  `json:"itemCount,omitempty"`
}

// BattleRankingEntry is one battle ranking result.
type BattleRankingEntry struct {
	ID        string       `json:"_id,omitempty"`
	Ranking   string       `json:"ranking,omitempty"`
	UserID    string       `json:"user,omitempty"`
	CountryID string       `json:"country,omitempty"`
	MUID      string       `json:"mu,omitempty"`
	Value     float64      `json:"value,omitempty"`
	Rank      int          `json:"rank,omitempty"`
	Badge     string       `json:"badge,omitempty"`
	Version   int          `json:"__v,omitempty"`
	CreatedAt string       `json:"createdAt,omitempty"`
	UpdatedAt string       `json:"updatedAt,omitempty"`
	LootItem  *RoundWeapon `json:"lootItem,omitempty"`
}

// ItemOffer is a market offer for an item.
type ItemOffer struct {
	ID        string      `json:"_id,omitempty"`
	CountryID string      `json:"country,omitempty"`
	CreatedAt string      `json:"createdAt,omitempty"`
	Item      RoundWeapon `json:"item,omitempty"`
	ItemCode  string      `json:"itemCode,omitempty"`
	Price     float64     `json:"price,omitempty"`
	UpdatedAt string      `json:"updatedAt,omitempty"`
	UserID    string      `json:"user,omitempty"`
}

// MercenaryContractAuction is a mercenary contract auction listing.
type MercenaryContractAuction struct {
	ID                string                        `json:"_id,omitempty"`
	BattleID          string                        `json:"battle,omitempty"`
	Bids              []MercenaryContractAuctionBid `json:"bids,omitempty"`
	Budget            float64                       `json:"budget,omitempty"`
	CountryID         string                        `json:"country,omitempty"`
	CreatedAt         string                        `json:"createdAt,omitempty"`
	CreatedBy         string                        `json:"createdBy,omitempty"`
	CurrentPayout     float64                       `json:"currentPayout,omitempty"`
	CurrentPerK       float64                       `json:"currentPerK,omitempty"`
	CurrentWinner     *string                       `json:"currentWinner,omitempty"`
	CurrentWinnerUser *string                       `json:"currentWinnerUser,omitempty"`
	Duration          int                           `json:"duration,omitempty"`
	ExpiresAt         string                        `json:"expiresAt,omitempty"`
	ForCountry        string                        `json:"forCountry,omitempty"`
	ForCountrySide    string                        `json:"forCountrySide,omitempty"`
	InitialPerK       float64                       `json:"initialPerK,omitempty"`
	MinimumDamage     int64                         `json:"minimumDamage,omitempty"`
	ProfessionalsOnly bool                          `json:"professionalsOnly,omitempty"`
	RoundID           string                        `json:"round,omitempty"`
	RoundNumber       int                           `json:"roundNumber,omitempty"`
	Status            string                        `json:"status,omitempty"`
	UpdatedAt         string                        `json:"updatedAt,omitempty"`
}

// MercenaryContractAuctionBid is a bid in a mercenary contract auction.
type MercenaryContractAuctionBid struct {
	BidAt  string  `json:"bidAt,omitempty"`
	MUID   string  `json:"mu,omitempty"`
	Payout float64 `json:"payout,omitempty"`
	PerK   float64 `json:"perK,omitempty"`
	UserID string  `json:"user,omitempty"`
}

// TransactionListItem is one transaction history entry.
type TransactionListItem struct {
	ID              string           `json:"_id,omitempty"`
	BuyerID         string           `json:"buyerId,omitempty"`
	CreatedAt       string           `json:"createdAt,omitempty"`
	Item            *TransactionItem `json:"item,omitempty"`
	ItemCode        string           `json:"itemCode,omitempty"`
	Money           float64          `json:"money,omitempty"`
	OfferCreatedAt  string           `json:"offerCreatedAt,omitempty"`
	Quantity        int              `json:"quantity,omitempty"`
	SellerCountryID string           `json:"sellerCountryId,omitempty"`
	SellerID        string           `json:"sellerId,omitempty"`
	SellerMUID      string           `json:"sellerMuId,omitempty"`
	TransactionType string           `json:"transactionType,omitempty"`
	UpdatedAt       string           `json:"updatedAt,omitempty"`
}

// TransactionItem describes the item attached to a transaction.
type TransactionItem struct {
	ID                string                `json:"_id,omitempty"`
	Code              string                `json:"code,omitempty"`
	LastAcquisitionAt string                `json:"lastAcquisitionAt,omitempty"`
	MaxState          int                   `json:"maxState,omitempty"`
	Quantity          int                   `json:"quantity,omitempty"`
	Skills            TransactionItemSkills `json:"skills,omitempty"`
	State             int                   `json:"state,omitempty"`
	Type              string                `json:"type,omitempty"`
}

// TransactionItemSkills contains item skill modifiers in a transaction.
type TransactionItemSkills struct {
	Armor           int `json:"armor,omitempty"`
	Attack          int `json:"attack,omitempty"`
	CriticalChance  int `json:"criticalChance,omitempty"`
	CriticalDamages int `json:"criticalDamages,omitempty"`
	Dodge           int `json:"dodge,omitempty"`
	Precision       int `json:"precision,omitempty"`
}

// WorkOffer is an offer to work for a company. It includes wageAfterTax and
// minLevel fields present in current live responses but absent upstream.
type WorkOffer struct {
	ID              string  `json:"_id,omitempty"`
	Citizenship     string  `json:"citizenship,omitempty"`
	CompanyID       string  `json:"company,omitempty"`
	CreatedAt       string  `json:"createdAt,omitempty"`
	InitialQuantity int     `json:"initialQuantity,omitempty"`
	MinEnergy       float64 `json:"minEnergy,omitempty"`
	MinLevel        float64 `json:"minLevel,omitempty"`
	MinProduction   float64 `json:"minProduction,omitempty"`
	Quantity        int     `json:"quantity,omitempty"`
	RegionID        string  `json:"region,omitempty"`
	Text            string  `json:"text,omitempty"`
	UpdatedAt       string  `json:"updatedAt,omitempty"`
	UserID          string  `json:"user,omitempty"`
	Wage            float64 `json:"wage,omitempty"`
	WageAfterTax    float64 `json:"wageAfterTax,omitempty"`
}

// WorkOfferListItem is a work offer returned in a paginated list.
type WorkOfferListItem = WorkOffer

// Workers contains the response to worker.getWorkers. The API returns a flat
// Workers slice when filtering by company, or WorkersPerCompany when filtering
// by user.
type Workers struct {
	Type              string                  `json:"type,omitempty"`
	Workers           []Worker                `json:"workers,omitempty"`
	WorkersPerCompany []WorkersPerCompanyItem `json:"workersPerCompany,omitempty"`
}

// Worker is an employment record returned by worker.getWorkers.
type Worker struct {
	ID                     string  `json:"_id,omitempty"`
	UserID                 string  `json:"user,omitempty"`
	CompanyID              string  `json:"company,omitempty"`
	EmployerID             string  `json:"employer,omitempty"`
	Wage                   float64 `json:"wage,omitempty"`
	JoinedAt               string  `json:"joinedAt,omitempty"`
	Fidelity               int     `json:"fidelity,omitempty"`
	LastFidelityIncreaseAt string  `json:"lastFidelityIncreaseAt,omitempty"`
	CreatedAt              string  `json:"createdAt,omitempty"`
	UpdatedAt              string  `json:"updatedAt,omitempty"`
	Version                int     `json:"__v,omitempty"`
}

// WorkersPerCompanyItem contains a company and its workers.
type WorkersPerCompanyItem struct {
	Company WorkerCompany `json:"company,omitempty"`
	Workers []Worker      `json:"workers,omitempty"`
}

// WorkerCompany is identifying information for a company in a worker listing.
type WorkerCompany struct {
	ID       string `json:"_id,omitempty"`
	ItemCode string `json:"itemCode,omitempty"`
	Name     string `json:"name,omitempty"`
}
