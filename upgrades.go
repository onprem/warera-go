package warera

import "context"

// UpgradeService handles region, company, and military unit upgrades.
type UpgradeService struct {
	client *Client
}

// GetByTypeAndEntity returns an upgrade matching the type and supplied entity ID.
func (s *UpgradeService) GetByTypeAndEntity(ctx context.Context, filter UpgradeFilter) (*Upgrade, error) {
	var upgrade Upgrade
	if err := s.client.trpc.Call(ctx, "upgrade.getUpgradeByTypeAndEntity", filter, &upgrade); err != nil {
		return nil, err
	}
	return &upgrade, nil
}
