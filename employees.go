package warera

import (
	"context"
	"iter"
)

// EmployeeService handles NPC employees in player companies.
type EmployeeService struct {
	client *Client
}

type getEmployeeInput struct {
	EmployeeID string `json:"employeeId"`
}

// GetCompanyRoster returns a company's employees and its estimated hourly
// output and costs.
//
// Unstable: not in the published OpenAPI document; may change without notice.
// Requires an API key.
func (s *EmployeeService) GetCompanyRoster(ctx context.Context, companyID string) (*CompanyRoster, error) {
	var roster CompanyRoster
	if err := s.client.trpc.Call(ctx, "employee.getCompanyRoster", getCompanyByIdInput{CompanyID: companyID}, &roster); err != nil {
		return nil, err
	}
	return &roster, nil
}

// GetByID returns an employee by ID.
//
// Unstable: not in the published OpenAPI document; may change without notice.
// Requires an API key.
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
//
// Unstable: not in the published OpenAPI document; may change without notice.
func (s *ManagerService) GetByCompanyID(ctx context.Context, companyID string) (*Manager, error) {
	var manager *Manager
	if err := s.client.trpc.Call(ctx, "manager.getByCompanyId", getCompanyByIdInput{CompanyID: companyID}, &manager); err != nil {
		return nil, err
	}
	return manager, nil
}

// GetOfferByCompanyID returns a company's open manager job, or nil when it
// has none.
//
// Unstable: not in the published OpenAPI document; may change without notice.
func (s *ManagerService) GetOfferByCompanyID(ctx context.Context, companyID string) (*ManagerOffer, error) {
	var offer *ManagerOffer
	if err := s.client.trpc.Call(ctx, "managerOffer.getByCompanyId", getCompanyByIdInput{CompanyID: companyID}, &offer); err != nil {
		return nil, err
	}
	return offer, nil
}

// ListOffers returns one page of open manager jobs matching the filter.
//
// Unstable: not in the published OpenAPI document; may change without notice.
func (s *ManagerService) ListOffers(ctx context.Context, filter ManagerOfferFilter) (*CursorPage[ManagerOffer], error) {
	var page CursorPage[ManagerOffer]
	if err := s.client.trpc.Call(ctx, "managerOffer.getManyPaginated", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// AllOffers returns an iterator over every open manager job matching the
// filter.
//
// Unstable: not in the published OpenAPI document; may change without notice.
func (s *ManagerService) AllOffers(ctx context.Context, filter ManagerOfferFilter) iter.Seq2[ManagerOffer, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[ManagerOffer], error) {
		f := filter
		f.Cursor = cursor
		return s.ListOffers(c, f)
	})
}

// GetPointValue returns how much a company produces per manager point.
//
// Unstable: not in the published OpenAPI document; may change without notice.
// Requires an API key.
func (s *ManagerService) GetPointValue(ctx context.Context, companyID string) (*ManagerPointValue, error) {
	var value ManagerPointValue
	if err := s.client.trpc.Call(ctx, "employee.getManagerPointValue", getCompanyByIdInput{CompanyID: companyID}, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

// GetWageStats returns the range manager pay rates must fall within.
//
// Unstable: not in the published OpenAPI document; may change without notice.
func (s *ManagerService) GetWageStats(ctx context.Context) (*ManagerWageStats, error) {
	var stats ManagerWageStats
	if err := s.client.trpc.Call(ctx, "managerOffer.getWageStats", struct{}{}, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
