package warera

import "context"

// TradingService handles market prices and exchange trading orders.
type TradingService struct {
	client *Client
}

type getTopOrdersInput struct {
	ItemCode string `json:"itemCode"`
	Limit    int    `json:"limit,omitempty"`
}

// GetPrices returns reference/current prices for tradeable items.
func (s *TradingService) GetPrices(ctx context.Context) ([]ItemPrice, error) {
	var prices []ItemPrice
	err := s.client.trpc.Call(ctx, "itemTrading.getPrices", struct{}{}, &prices)
	if err != nil {
		return nil, err
	}
	return prices, nil
}

// GetTopOrders returns best buy/sell market orders for a given item code.
func (s *TradingService) GetTopOrders(ctx context.Context, itemCode string, limit int) ([]TopOrder, error) {
	var orders []TopOrder
	err := s.client.trpc.Call(ctx, "tradingOrder.getTopOrders", getTopOrdersInput{
		ItemCode: itemCode,
		Limit:    limit,
	}, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
