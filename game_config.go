package warera

import "context"

// GameConfigService handles game clock, timers, and configuration parameters.
type GameConfigService struct {
	client *Client
}

// GetDates returns the official game clock and upcoming server event times
// (day rollover, energy regeneration, election schedules, etc.).
func (s *GameConfigService) GetDates(ctx context.Context) (*GameDates, error) {
	var dates GameDates
	err := s.client.trpc.Call(ctx, "gameConfig.getDates", struct{}{}, &dates)
	if err != nil {
		return nil, err
	}
	return &dates, nil
}

// GetGameConfig returns the general game config settings.
func (s *GameConfigService) GetGameConfig(ctx context.Context) (map[string]any, error) {
	var config map[string]any
	err := s.client.trpc.Call(ctx, "gameConfig.getGameConfig", struct{}{}, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
