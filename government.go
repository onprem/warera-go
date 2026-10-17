package warera

import "context"

// GovernmentService handles country government hierarchies.
type GovernmentService struct {
	client *Client
}

type getGovernmentByCountryIdInput struct {
	CountryID string `json:"countryId"`
}

// GetByCountryID returns government positions (president, ministers, congress) for a country.
func (s *GovernmentService) GetByCountryID(ctx context.Context, countryID string) (*Government, error) {
	var gov Government
	err := s.client.trpc.Call(ctx, "government.getByCountryId", getGovernmentByCountryIdInput{CountryID: countryID}, &gov)
	if err != nil {
		return nil, err
	}
	return &gov, nil
}
