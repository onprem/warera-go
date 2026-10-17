package warera

// These response models cover procedures whose contracts are not described by
// the upstream OpenAPI document. Their field shapes follow spec/Responses.d.ts.

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

// BattleLootPoolItem is one item in a battle loot pool.
type BattleLootPoolItem struct {
	Item  RoundWeapon `json:"item"`
	Pool  string      `json:"pool"`
	Rank  int         `json:"rank"`
	Round string      `json:"round,omitempty"`
}

// BattleRanking contains ranking entries for a battle or war.
type BattleRanking struct {
	Rankings []BattleRankingEntry `json:"rankings,omitempty"`
}

// BattleRankingEntry is one battle ranking result.
type BattleRankingEntry struct {
	CountryID string  `json:"country,omitempty"`
	Value     float64 `json:"value,omitempty"`
	Rank      int     `json:"rank,omitempty"`
	ID        string  `json:"_id,omitempty"`
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

// WorkOffer is an offer to work for a company.
type WorkOffer struct {
	ID              string  `json:"_id,omitempty"`
	Citizenship     string  `json:"citizenship,omitempty"`
	CompanyID       string  `json:"company,omitempty"`
	CreatedAt       string  `json:"createdAt,omitempty"`
	InitialQuantity int     `json:"initialQuantity,omitempty"`
	MinEnergy       int     `json:"minEnergy,omitempty"`
	MinProduction   int     `json:"minProduction,omitempty"`
	Quantity        int     `json:"quantity,omitempty"`
	RegionID        string  `json:"region,omitempty"`
	Text            string  `json:"text,omitempty"`
	UpdatedAt       string  `json:"updatedAt,omitempty"`
	UserID          string  `json:"user,omitempty"`
	Wage            float64 `json:"wage,omitempty"`
}

// WorkOfferListItem is a work offer returned in a paginated list.
type WorkOfferListItem = WorkOffer

// Workers groups a user's workers by company.
type Workers struct {
	Type              string                  `json:"type,omitempty"`
	WorkersPerCompany []WorkersPerCompanyItem `json:"workersPerCompany,omitempty"`
}

// WorkersPerCompanyItem contains a company and its workers.
type WorkersPerCompanyItem struct {
	Company WorkerCompany `json:"company,omitempty"`
	Workers []any         `json:"workers,omitempty"`
}

// WorkerCompany is identifying information for a company in a worker listing.
type WorkerCompany struct {
	ID       string `json:"_id,omitempty"`
	ItemCode string `json:"itemCode,omitempty"`
	Name     string `json:"name,omitempty"`
}
