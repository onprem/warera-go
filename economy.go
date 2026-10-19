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

// GetMany retrieves company details for the supplied IDs. Results preserve the
// input order and duplicate IDs. An empty input returns an empty slice without
// making a request. If any lookup fails, GetMany returns no records and a
// [BatchLookupError] identifying the failed IDs. Large inputs are split into
// requests according to the client's configured batch-size limit.
func (s *CompanyService) GetMany(ctx context.Context, companyIDs []string) ([]Company, error) {
	return getManyByID[Company](ctx, s.client, "companies", "company.getById", companyIDs, func(id string) any {
		return getCompanyByIdInput{CompanyID: id}
	})
}

// List returns one page of company IDs matching the filter.
func (s *CompanyService) List(ctx context.Context, filter CompanyFilter) (*CursorPage[string], error) {
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
		return s.List(c, f)
	})
}

// --- Trading Service ---

// TradingService handles market prices and exchange trading orders.
type TradingService struct {
	client *Client
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
	err := s.client.trpc.Call(ctx, "tradingOrder.getTopOrders", TradingOrderTopOrdersFilter{
		ItemCode: itemCode,
		Limit:    limit,
	}, &orders)
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

// ItemOfferService handles individual marketplace offers.
type ItemOfferService struct {
	client *Client
}

type getItemOfferByIDInput struct {
	ItemOfferID string `json:"itemOfferId"`
}

// GetByID returns an item offer by ID.
func (s *ItemOfferService) GetByID(ctx context.Context, itemOfferID string) (*ItemOffer, error) {
	var offer ItemOffer
	if err := s.client.trpc.Call(ctx, "itemOffer.getById", getItemOfferByIDInput{ItemOfferID: itemOfferID}, &offer); err != nil {
		return nil, err
	}
	return &offer, nil
}

// MercenaryAuctionService handles paginated mercenary contract auctions.
type MercenaryAuctionService struct {
	client *Client
}

// List returns one page of auctions matching the filter.
func (s *MercenaryAuctionService) List(ctx context.Context, filter MercenaryAuctionFilter) (*CursorPage[MercenaryContractAuction], error) {
	var page CursorPage[MercenaryContractAuction]
	if err := s.client.trpc.Call(ctx, "mercenaryContractAuction.getPaginatedAuctions", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over all auctions matching the filter.
func (s *MercenaryAuctionService) All(ctx context.Context, filter MercenaryAuctionFilter) iter.Seq2[MercenaryContractAuction, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[MercenaryContractAuction], error) {
		f := filter
		f.Cursor = cursor
		return s.List(c, f)
	})
}

// TransactionService handles paginated transaction history.
type TransactionService struct {
	client *Client
}

// List returns one page of transactions matching the filter.
func (s *TransactionService) List(ctx context.Context, filter TransactionFilter) (*CursorPage[TransactionListItem], error) {
	var page CursorPage[TransactionListItem]
	if err := s.client.trpc.Call(ctx, "transaction.getPaginatedTransactions", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over all transactions matching the filter.
func (s *TransactionService) All(ctx context.Context, filter TransactionFilter) iter.Seq2[TransactionListItem, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[TransactionListItem], error) {
		f := filter
		f.Cursor = cursor
		return s.List(c, f)
	})
}
