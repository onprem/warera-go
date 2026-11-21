package warera

import "encoding/json"

// Response models for market and economy statistics. WarEra's published
// OpenAPI document does not describe these procedures; the models follow
// recorded responses.

// ItemTrading is an item's daily market trading history.
type ItemTrading struct {
	ID       string `json:"_id"`
	ItemCode string `json:"itemCode"`
	// CurrentValue is the average price of the most recent trading day, the
	// same value [TradingService.GetPrices] reports for the item.
	CurrentValue float64 `json:"currentValue"`
	// MinPrice and MaxPrice are price bounds stored for some items (zero when
	// absent). Market prices are not held within them.
	MinPrice float64 `json:"minPrice"`
	MaxPrice float64 `json:"maxPrice"`
	// Values holds up to 30 of the most recent trading days, oldest first.
	// The last day is usually the current, still incomplete, day.
	Values []ItemTradingDay `json:"values"`
	// BuyOrders and SellOrders have been empty in every observed response, so
	// their element shape is unknown. Use [TradingService.GetTopOrders] for
	// the order book.
	BuyOrders  []json.RawMessage `json:"buyOrders"`
	SellOrders []json.RawMessage `json:"sellOrders"`
}

// ItemTradingDay is one day of market trades for an item.
type ItemTradingDay struct {
	// ValueAt is the trading day as YYYY-MM-DD.
	ValueAt Date `json:"valueAt"`
	// AvgValue is the volume-weighted average price: TotalValue divided by
	// TotalQuantity.
	AvgValue float64 `json:"avgValue"`
	// TotalValue is the money that changed hands.
	TotalValue        float64 `json:"totalValue"`
	TotalQuantity     int     `json:"totalQuantity"`
	TransactionsCount int     `json:"transactionsCount"`
}

// OwnerOrdersFilter selects whose open market orders
// [TradingService.GetPublicOrdersByOwner] returns. A user's orders include
// those they placed on behalf of a country, military unit, or party.
type OwnerOrdersFilter struct {
	UserID    string `json:"userId,omitempty"`
	CountryID string `json:"countryId,omitempty"`
	MUID      string `json:"muId,omitempty"`
	PartyID   string `json:"partyId,omitempty"`
}

// OwnerOrders is an owner's open market orders and their totals.
type OwnerOrders struct {
	BuyOrders  []TopOrder `json:"buyOrders"`
	SellOrders []TopOrder `json:"sellOrders"`
	// AllOrders holds the buy orders followed by the sell orders.
	AllOrders []TopOrder `json:"allOrders"`
	// TotalBuyMoneyInvested is the money locked in the buy orders.
	TotalBuyMoneyInvested float64 `json:"totalBuyMoneyInvested"`
	// TotalSellMoneyExpected is the money the sell orders would earn if filled.
	TotalSellMoneyExpected float64 `json:"totalSellMoneyExpected"`
	// TotalSellQuantities is the quantity offered for sale, keyed by item code.
	TotalSellQuantities map[string]int `json:"totalSellQuantities"`
}

// PPValueByItem is the money one production point earns when spent
// producing each item at current market prices, keyed by item code.
type PPValueByItem map[string]float64

// DonationFilter selects the recipient whose donors
// [DonationService.List] returns. Set one of CountryID, MUID, or PartyID.
type DonationFilter struct {
	CountryID string `json:"countryId,omitempty"`
	MUID      string `json:"muId,omitempty"`
	PartyID   string `json:"partyId,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
	// Limit is the page size, at most 100.
	Limit int `json:"limit,omitempty"`
}

// DonationTotalFilter selects the recipient [DonationService.GetTotal]
// totals. Set one of CountryID, MUID, or PartyID.
type DonationTotalFilter struct {
	CountryID string `json:"countryId,omitempty"`
	MUID      string `json:"muId,omitempty"`
	PartyID   string `json:"partyId,omitempty"`
}

// Donation is one user's cumulative money donations to one country, military
// unit, or party. Only the recipient's ID field is set.
type Donation struct {
	ID        string  `json:"_id"`
	UserID    string  `json:"userId"`
	CountryID string  `json:"countryId"`
	MUID      string  `json:"muId"`
	PartyID   string  `json:"partyId"`
	Amount    float64 `json:"amount"`
	// CreatedAt is the time of the first donation and UpdatedAt of the latest.
	CreatedAt Time `json:"createdAt"`
	UpdatedAt Time `json:"updatedAt"`
}

// DonationTotal is the total donated to a recipient.
type DonationTotal struct {
	TotalAmount float64 `json:"totalAmount"`
	DonorCount  int     `json:"donorCount"`
}
