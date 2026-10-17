package warera

import (
	"context"
	"iter"
)

// CompanyService handles interactions with company endpoints.
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

// GetCompanies returns a paginated list of companies.
func (s *CompanyService) GetCompanies(ctx context.Context, filter CompanyFilter) (*CursorPage[Company], error) {
	var page CursorPage[Company]
	err := s.client.trpc.Call(ctx, "company.getCompanies", filter, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over all companies matching the filter across pages.
func (s *CompanyService) All(ctx context.Context, filter CompanyFilter) iter.Seq2[Company, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cur string) (*CursorPage[Company], error) {
		f := filter
		f.Cursor = cur
		return s.GetCompanies(c, f)
	})
}
