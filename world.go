package warera

import "context"

// --- Country Service ---

// CountryService handles interactions with country statistics and diplomacy.
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

// --- Region Service ---

// RegionService handles interactions with world regions and geography.
type RegionService struct {
	client *Client
}

type getRegionByIdInput struct {
	RegionID string `json:"regionId"`
}

// GetByID returns detailed region info by ID.
func (s *RegionService) GetByID(ctx context.Context, regionID string) (*Region, error) {
	var region Region
	err := s.client.trpc.Call(ctx, "region.getById", getRegionByIdInput{RegionID: regionID}, &region)
	if err != nil {
		return nil, err
	}
	return &region, nil
}

// GetAll returns a map of all world regions indexed by region ID.
func (s *RegionService) GetAll(ctx context.Context) (map[string]Region, error) {
	var regions map[string]Region
	err := s.client.trpc.Call(ctx, "region.getRegionsObject", struct{}{}, &regions)
	if err != nil {
		return nil, err
	}
	return regions, nil
}

// --- Government Service ---

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
