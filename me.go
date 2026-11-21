package warera

import "context"

// MeService reads the account that owns the client's API key: its NPC
// employees, their workforce totals, and their story events. Every method
// needs an API key and answers for that key's owner only, so a tool serving
// several players creates one client per player's key.
//
// The API forbids API keys from most other caller-scoped procedures, such as
// the owner's profile, inventory, and missions; MeService wraps the ones it
// allows. Like the other procedures missing from the published OpenAPI
// document, they may change without notice.
type MeService struct {
	client *Client
}

type employeeIDInput struct {
	EmployeeID string `json:"employeeId"`
}

// GetEmployees returns all of the owner's employees across their companies.
//
// Unstable: not in the published OpenAPI document; may change without notice.
// Requires an API key.
func (s *MeService) GetEmployees(ctx context.Context) ([]Employee, error) {
	var employees []Employee
	if err := s.client.trpc.Call(ctx, "employee.getMyEmployees", nil, &employees); err != nil {
		return nil, err
	}
	return employees, nil
}

// GetWorkforceSummary totals the owner's employees and their hourly output
// and costs.
//
// Unstable: not in the published OpenAPI document; may change without notice.
// Requires an API key.
func (s *MeService) GetWorkforceSummary(ctx context.Context) (*WorkforceSummary, error) {
	var summary WorkforceSummary
	if err := s.client.trpc.Call(ctx, "employee.getMyWorkforceSummary", nil, &summary); err != nil {
		return nil, err
	}
	return &summary, nil
}

// GetPendingEmployeeEvents returns the owner's employee events that wait for
// a choice. Its responses have only been seen empty; the events are assumed
// to have the shape [MeService.GetEmployeeEventHistory] returns.
//
// Unstable: not in the published OpenAPI document; may change without notice.
// Requires an API key.
func (s *MeService) GetPendingEmployeeEvents(ctx context.Context) ([]EmployeeEvent, error) {
	var events []EmployeeEvent
	if err := s.client.trpc.Call(ctx, "employeeEvent.getMine", nil, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// GetPendingEmployeeEventCount returns how many employee events wait for the
// owner's choice.
//
// Unstable: not in the published OpenAPI document; may change without notice.
// Requires an API key.
func (s *MeService) GetPendingEmployeeEventCount(ctx context.Context) (int, error) {
	var count int
	if err := s.client.trpc.Call(ctx, "employeeEvent.getPendingCount", nil, &count); err != nil {
		return 0, err
	}
	return count, nil
}

// GetEmployeeEventHistory returns the events one of the owner's employees
// took part in, most recently resolved first. For another owner's employee
// the API returns an empty list.
//
// Unstable: not in the published OpenAPI document; may change without notice.
// Requires an API key.
func (s *MeService) GetEmployeeEventHistory(ctx context.Context, employeeID string) ([]EmployeeEvent, error) {
	var events []EmployeeEvent
	if err := s.client.trpc.Call(ctx, "employeeEvent.getHistory", employeeIDInput{EmployeeID: employeeID}, &events); err != nil {
		return nil, err
	}
	return events, nil
}
