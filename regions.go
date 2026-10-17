package warera

import "context"

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

// GetAll returns a map or slice of all world regions.
func (s *RegionService) GetAll(ctx context.Context) (map[string]Region, error) {
	var regions map[string]Region
	err := s.client.trpc.Call(ctx, "region.getRegionsObject", struct{}{}, &regions)
	if err != nil {
		return nil, err
	}
	return regions, nil
}
