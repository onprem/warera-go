package warera

import (
	"context"
	"iter"
)

// The procedures in this file are served by the WarEra API but are not listed
// in its published OpenAPI document, so they may change without notice. They
// require an API key.

type itemCodeInput struct {
	ItemCode string `json:"itemCode"`
}

type getTopOrdersPerItemCodeInput struct {
	ItemCodes []string `json:"itemCodes"`
	Limit     int      `json:"limit,omitempty"`
}

// GetItemTrading returns an item's daily market trading history: volume,
// turnover, and average price for each of the most recent trading days.
func (s *TradingService) GetItemTrading(ctx context.Context, itemCode string) (*ItemTrading, error) {
	var trading ItemTrading
	if err := s.client.trpc.Call(ctx, "itemTrading.getItemTrading", itemCodeInput{ItemCode: itemCode}, &trading); err != nil {
		return nil, err
	}
	return &trading, nil
}

// GetTopOrdersPerItem returns the best buy and sell orders for each of the
// given item codes in one request, keyed by item code. Limit caps the orders
// per side and item (the API allows 1 to 100 and defaults to 1). An empty
// input returns an empty map without making a request.
func (s *TradingService) GetTopOrdersPerItem(ctx context.Context, itemCodes []string, limit int) (map[string]TopOrders, error) {
	if len(itemCodes) == 0 {
		return map[string]TopOrders{}, nil
	}
	var orders map[string]TopOrders
	if err := s.client.trpc.Call(ctx, "tradingOrder.getTopOrdersPerItemCode", getTopOrdersPerItemCodeInput{ItemCodes: itemCodes, Limit: limit}, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// GetPublicOrdersByOwner returns the open market orders of a user, country,
// military unit, or party, with their totals. An empty filter returns no
// orders.
func (s *TradingService) GetPublicOrdersByOwner(ctx context.Context, filter OwnerOrdersFilter) (*OwnerOrders, error) {
	var orders OwnerOrders
	if err := s.client.trpc.Call(ctx, "tradingOrder.getPublicOrdersByOwner", filter, &orders); err != nil {
		return nil, err
	}
	return &orders, nil
}

// GetMaxBountyRate returns the highest bounty rate a mercenary contract may
// offer. It follows the PP money value (see
// [MarketStatsService.GetPPMoneyValue]) times the game config's
// mercenaryContract.maxBountyPpMultiplier.
func (s *GameConfigService) GetMaxBountyRate(ctx context.Context) (float64, error) {
	var rate float64
	if err := s.client.trpc.Call(ctx, "gameConfig.getMaxBountyRate", struct{}{}, &rate); err != nil {
		return 0, err
	}
	return rate, nil
}

// MarketStatsService handles game-wide economy statistics.
type MarketStatsService struct {
	client *Client
}

// GetPPMoneyValue returns the game-wide value of a production point in money.
func (s *MarketStatsService) GetPPMoneyValue(ctx context.Context) (float64, error) {
	var value float64
	if err := s.client.trpc.Call(ctx, "gameStat.getPpMoneyValue", nil, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// GetPPMoneyValueByItem returns, for each producible item, the money a
// production point earns when spent producing it at current market prices.
func (s *MarketStatsService) GetPPMoneyValueByItem(ctx context.Context) (PPValueByItem, error) {
	var values PPValueByItem
	if err := s.client.trpc.Call(ctx, "gameStat.getPpMoneyValueByItem", nil, &values); err != nil {
		return nil, err
	}
	return values, nil
}

// GetWorldDevelopment returns the total development of all regions.
func (s *MarketStatsService) GetWorldDevelopment(ctx context.Context) (float64, error) {
	var development float64
	if err := s.client.trpc.Call(ctx, "gameStat.getWorldDevelopment", nil, &development); err != nil {
		return 0, err
	}
	return development, nil
}

// GetEquipmentAveragePrice returns the average market sale price of an
// equipment item code, such as "knife" or "helmet3". It returns 0 for codes
// that are not tradable equipment.
func (s *MarketStatsService) GetEquipmentAveragePrice(ctx context.Context, itemCode string) (float64, error) {
	var price float64
	if err := s.client.trpc.Call(ctx, "gameStat.getEquipmentAvgByCode", itemCodeInput{ItemCode: itemCode}, &price); err != nil {
		return 0, err
	}
	return price, nil
}

// DonationService handles money donations to countries, military units, and
// parties.
type DonationService struct {
	client *Client
}

// List returns one page of donors to the recipient in the filter, largest
// total first. An empty filter returns no donors.
func (s *DonationService) List(ctx context.Context, filter DonationFilter) (*CursorPage[Donation], error) {
	var page CursorPage[Donation]
	if err := s.client.trpc.Call(ctx, "donation.getManyPaginated", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over every donor to the recipient in the filter.
func (s *DonationService) All(ctx context.Context, filter DonationFilter) iter.Seq2[Donation, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[Donation], error) {
		f := filter
		f.Cursor = cursor
		return s.List(c, f)
	})
}

// GetTotal returns the total donated to the recipient in the filter and the
// number of donors.
func (s *DonationService) GetTotal(ctx context.Context, filter DonationTotalFilter) (*DonationTotal, error) {
	var total DonationTotal
	if err := s.client.trpc.Call(ctx, "donation.getTotalDonations", filter, &total); err != nil {
		return nil, err
	}
	return &total, nil
}
