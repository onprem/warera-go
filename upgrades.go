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

// GetMapUpgrades returns every region upgrade in the world, such as bunkers
// and bases, with its level and status, in one call.
//
// Unstable: not in the published OpenAPI document; may change without notice.
func (s *UpgradeService) GetMapUpgrades(ctx context.Context) ([]MapUpgrade, error) {
	var upgrades []MapUpgrade
	if err := s.client.trpc.Call(ctx, "upgrade.getMapUpgrades", nil, &upgrades); err != nil {
		return nil, err
	}
	return upgrades, nil
}
