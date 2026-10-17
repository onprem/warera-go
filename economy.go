package warera

import (
	"context"
	"iter"
)

// --- Company Service ---

// CompanyService handles interactions with company management and production.
type CompanyService struct {
	client *Client
}

type getCompanyByIdInput struct {
	CompanyID string `json:"companyId"`
}

// GetByID returns detailed information about a specific company.
func (s *CompanyService) GetByID(ctx context.Context, companyID string) (*Company, error) {
	var company Company
	err := s.client.trpc.Call(ctx, "company.getById", getCompanyByIdInput{CompanyID: companyID}, &company)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

// GetCompanies returns a paginated list of company IDs.
func (s *CompanyService) GetCompanies(ctx context.Context, filter CompanyFilter) (*CursorPage[string], error) {
	var page CursorPage[string]
	err := s.client.trpc.Call(ctx, "company.getCompanies", filter, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over all company IDs matching the filter across pages.
func (s *CompanyService) All(ctx context.Context, filter CompanyFilter) iter.Seq2[string, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cur string) (*CursorPage[string], error) {
		f := filter
		f.Cursor = cur
		return s.GetCompanies(c, f)
	})
}

// --- Trading Service ---

// TradingService handles market prices and exchange trading orders.
type TradingService struct {
	client *Client
}

type getTopOrdersInput struct {
	ItemCode string `json:"itemCode"`
	Limit    int    `json:"limit,omitempty"`
}

// GetPrices returns reference/current prices for tradeable items.
func (s *TradingService) GetPrices(ctx context.Context) (ItemPrices, error) {
	var prices ItemPrices
	err := s.client.trpc.Call(ctx, "itemTrading.getPrices", struct{}{}, &prices)
	if err != nil {
		return nil, err
	}
	return prices, nil
}

// GetTopOrders returns best buy/sell market orders for a given item code.
func (s *TradingService) GetTopOrders(ctx context.Context, itemCode string, limit int) (*TopOrders, error) {
	var orders TopOrders
	err := s.client.trpc.Call(ctx, "tradingOrder.getTopOrders", getTopOrdersInput{
		ItemCode: itemCode,
		Limit:    limit,
	}, &orders)
	if err != nil {
		return nil, err
	}
	return &orders, nil
}
