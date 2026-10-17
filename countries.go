package warera

import "context"

// CountryService handles interactions with country endpoints.
type CountryService struct {
	client *Client
}

type getCountryByIdInput struct {
	CountryID string `json:"countryId"`
}

// GetAll returns a list of all countries in the game world.
func (s *CountryService) GetAll(ctx context.Context) ([]Country, error) {
	var countries []Country
	err := s.client.trpc.Call(ctx, "country.getAllCountries", struct{}{}, &countries)
	if err != nil {
		return nil, err
	}
	return countries, nil
}

// GetByID returns detailed information for a specific country by its ID.
func (s *CountryService) GetByID(ctx context.Context, countryID string) (*Country, error) {
	var country Country
	err := s.client.trpc.Call(ctx, "country.getCountryById", getCountryByIdInput{CountryID: countryID}, &country)
	if err != nil {
		return nil, err
	}
	return &country, nil
}
