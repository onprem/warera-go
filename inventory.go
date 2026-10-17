package warera

import "context"

// InventoryService handles a user's equipped items.
type InventoryService struct {
	client *Client
}

type fetchCurrentEquipmentInput struct {
	UserID string `json:"userId"`
}

// GetCurrentEquipment returns the currently equipped items for a user.
func (s *InventoryService) GetCurrentEquipment(ctx context.Context, userID string) (*CurrentEquipment, error) {
	var equipment CurrentEquipment
	if err := s.client.trpc.Call(ctx, "inventory.fetchCurrentEquipment", fetchCurrentEquipmentInput{UserID: userID}, &equipment); err != nil {
		return nil, err
	}
	return &equipment, nil
}
