package warera

import (
	"context"
	"iter"
)

// BattleService handles interactions with battles and military conflicts.
type BattleService struct {
	client *Client
}

type getBattleByIdInput struct {
	BattleID string `json:"battleId"`
}

type getLiveBattleDataInput struct {
	BattleID    string `json:"battleId"`
	RoundNumber int    `json:"roundNumber,omitempty"`
}

// GetByID returns battle information by ID.
func (s *BattleService) GetByID(ctx context.Context, battleID string) (*Battle, error) {
	var battle Battle
	err := s.client.trpc.Call(ctx, "battle.getById", getBattleByIdInput{BattleID: battleID}, &battle)
	if err != nil {
		return nil, err
	}
	return &battle, nil
}

// GetBattles returns a page of battles matching the filter.
func (s *BattleService) GetBattles(ctx context.Context, filter BattleFilter) (*CursorPage[Battle], error) {
	var page CursorPage[Battle]
	err := s.client.trpc.Call(ctx, "battle.getBattles", filter, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over all battles matching the filter across pages.
func (s *BattleService) All(ctx context.Context, filter BattleFilter) iter.Seq2[Battle, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cur string) (*CursorPage[Battle], error) {
		f := filter
		f.Cursor = cur
		return s.GetBattles(c, f)
	})
}

// GetLiveBattleData returns real-time damage and combat state for a battle and round.
func (s *BattleService) GetLiveBattleData(ctx context.Context, battleID string, roundNumber int) (*LiveBattleData, error) {
	var live LiveBattleData
	err := s.client.trpc.Call(ctx, "battle.getLiveBattleData", getLiveBattleDataInput{
		BattleID:    battleID,
		RoundNumber: roundNumber,
	}, &live)
	if err != nil {
		return nil, err
	}
	return &live, nil
}
