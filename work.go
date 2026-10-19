package warera

import (
	"context"
	"iter"
)

// WorkOfferService handles offers to work at companies.
type WorkOfferService struct {
	client *Client
}

type getWorkOfferByIDInput struct {
	WorkOfferID string `json:"workOfferId"`
}

type getWorkOfferByCompanyIDInput struct {
	CompanyID string `json:"companyId"`
}

// GetByID returns a work offer by ID.
func (s *WorkOfferService) GetByID(ctx context.Context, workOfferID string) (*WorkOffer, error) {
	var offer WorkOffer
	if err := s.client.trpc.Call(ctx, "workOffer.getById", getWorkOfferByIDInput{WorkOfferID: workOfferID}, &offer); err != nil {
		return nil, err
	}
	return &offer, nil
}

// GetByCompanyID returns the work offer for a company.
func (s *WorkOfferService) GetByCompanyID(ctx context.Context, companyID string) (*WorkOffer, error) {
	var offer WorkOffer
	if err := s.client.trpc.Call(ctx, "workOffer.getWorkOfferByCompanyId", getWorkOfferByCompanyIDInput{CompanyID: companyID}, &offer); err != nil {
		return nil, err
	}
	return &offer, nil
}

// List returns one page of work offers matching the filter.
func (s *WorkOfferService) List(ctx context.Context, filter WorkOfferFilter) (*CursorPage[WorkOfferListItem], error) {
	var page CursorPage[WorkOfferListItem]
	if err := s.client.trpc.Call(ctx, "workOffer.getWorkOffersPaginated", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over all work offers matching the filter.
func (s *WorkOfferService) All(ctx context.Context, filter WorkOfferFilter) iter.Seq2[WorkOfferListItem, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[WorkOfferListItem], error) {
		f := filter
		f.Cursor = cursor
		return s.List(c, f)
	})
}

// WorkerService handles worker listings and counts.
type WorkerService struct {
	client *Client
}

type getTotalWorkersCountInput struct {
	UserID string `json:"userId"`
}

// GetTotalWorkersCount returns the total worker count for a user.
func (s *WorkerService) GetTotalWorkersCount(ctx context.Context, userID string) (int, error) {
	var count int
	if err := s.client.trpc.Call(ctx, "worker.getTotalWorkersCount", getTotalWorkersCountInput{UserID: userID}, &count); err != nil {
		return 0, err
	}
	return count, nil
}

// GetWorkers returns workers grouped by company for the supplied filters.
func (s *WorkerService) GetWorkers(ctx context.Context, filter WorkerFilter) (*Workers, error) {
	var workers Workers
	if err := s.client.trpc.Call(ctx, "worker.getWorkers", filter, &workers); err != nil {
		return nil, err
	}
	return &workers, nil
}
