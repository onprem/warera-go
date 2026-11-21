package warera

import (
	"context"
	"iter"
)

// The procedures in this file are served by the WarEra API but are not listed
// in its published OpenAPI document, so they may change without notice. They
// require an API key.

// EmployeeService handles NPC employees in player companies.
type EmployeeService struct {
	client *Client
}

type getEmployeeInput struct {
	EmployeeID string `json:"employeeId"`
}

// GetCompanyRoster returns a company's employees and its estimated hourly
// output and costs.
func (s *EmployeeService) GetCompanyRoster(ctx context.Context, companyID string) (*CompanyRoster, error) {
	var roster CompanyRoster
	if err := s.client.trpc.Call(ctx, "employee.getCompanyRoster", getCompanyByIdInput{CompanyID: companyID}, &roster); err != nil {
		return nil, err
	}
	return &roster, nil
}

// GetByID returns an employee by ID.
func (s *EmployeeService) GetByID(ctx context.Context, employeeID string) (*Employee, error) {
	var employee Employee
	if err := s.client.trpc.Call(ctx, "employee.getById", getEmployeeInput{EmployeeID: employeeID}, &employee); err != nil {
		return nil, err
	}
	return &employee, nil
}

// ManagerService handles company managers and manager job offers.
type ManagerService struct {
	client *Client
}

// GetByCompanyID returns a company's manager, or nil when it has none.
func (s *ManagerService) GetByCompanyID(ctx context.Context, companyID string) (*Manager, error) {
	var manager *Manager
	if err := s.client.trpc.Call(ctx, "manager.getByCompanyId", getCompanyByIdInput{CompanyID: companyID}, &manager); err != nil {
		return nil, err
	}
	return manager, nil
}

// GetOfferByCompanyID returns a company's open manager job, or nil when it
// has none.
func (s *ManagerService) GetOfferByCompanyID(ctx context.Context, companyID string) (*ManagerOffer, error) {
	var offer *ManagerOffer
	if err := s.client.trpc.Call(ctx, "managerOffer.getByCompanyId", getCompanyByIdInput{CompanyID: companyID}, &offer); err != nil {
		return nil, err
	}
	return offer, nil
}

// ListOffers returns one page of open manager jobs matching the filter.
func (s *ManagerService) ListOffers(ctx context.Context, filter ManagerOfferFilter) (*CursorPage[ManagerOffer], error) {
	var page CursorPage[ManagerOffer]
	if err := s.client.trpc.Call(ctx, "managerOffer.getManyPaginated", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// AllOffers returns an iterator over every open manager job matching the
// filter.
func (s *ManagerService) AllOffers(ctx context.Context, filter ManagerOfferFilter) iter.Seq2[ManagerOffer, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[ManagerOffer], error) {
		f := filter
		f.Cursor = cursor
		return s.ListOffers(c, f)
	})
}

// GetWageStats returns the range manager pay rates must fall within.
func (s *ManagerService) GetWageStats(ctx context.Context) (*ManagerWageStats, error) {
	var stats ManagerWageStats
	if err := s.client.trpc.Call(ctx, "managerOffer.getWageStats", struct{}{}, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
