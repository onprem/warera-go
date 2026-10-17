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
	err := s.client.trpc.Call(ctx, "battle.getLiveBattleData", BattleLiveBattleDataFilter{
		BattleID:    battleID,
		RoundNumber: roundNumber,
	}, &live)
	if err != nil {
		return nil, err
	}
	return &live, nil
}

// BattleOrderService handles battle side orders.
type BattleOrderService struct {
	client *Client
}

// GetByBattle returns orders for a battle side.
func (s *BattleOrderService) GetByBattle(ctx context.Context, battleID, side string) ([]BattleOrder, error) {
	var orders []BattleOrder
	if err := s.client.trpc.Call(ctx, "battleOrder.getByBattle", BattleOrderByBattleFilter{
		BattleID: battleID,
		Side:     side,
	}, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

// MilitaryUnitService handles military unit profiles and listings.
type MilitaryUnitService struct {
	client *Client
}

type getMilitaryUnitByIDInput struct {
	MUID string `json:"muId"`
}

// GetByID returns a military unit by ID.
func (s *MilitaryUnitService) GetByID(ctx context.Context, muID string) (*MilitaryUnit, error) {
	var unit MilitaryUnit
	if err := s.client.trpc.Call(ctx, "mu.getById", getMilitaryUnitByIDInput{MUID: muID}, &unit); err != nil {
		return nil, err
	}
	return &unit, nil
}

// GetManyPaginated returns one page of military units matching the filter.
func (s *MilitaryUnitService) GetManyPaginated(ctx context.Context, filter MUFilter) (*CursorPage[MilitaryUnit], error) {
	var page CursorPage[MilitaryUnit]
	if err := s.client.trpc.Call(ctx, "mu.getManyPaginated", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over all military units matching the filter.
func (s *MilitaryUnitService) All(ctx context.Context, filter MUFilter) iter.Seq2[MilitaryUnit, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[MilitaryUnit], error) {
		f := filter
		f.Cursor = cursor
		return s.GetManyPaginated(c, f)
	})
}

// RoundService handles battle round data and hit histories.
type RoundService struct {
	client *Client
}

type getRoundByIDInput struct {
	RoundID string `json:"roundId"`
}

// GetByID returns a round by ID.
func (s *RoundService) GetByID(ctx context.Context, roundID string) (*Round, error) {
	var round Round
	if err := s.client.trpc.Call(ctx, "round.getById", getRoundByIDInput{RoundID: roundID}, &round); err != nil {
		return nil, err
	}
	return &round, nil
}

// GetLastHits returns recent hits grouped by battle side for a round.
func (s *RoundService) GetLastHits(ctx context.Context, roundID string) (*RoundLastHits, error) {
	var hits RoundLastHits
	if err := s.client.trpc.Call(ctx, "round.getLastHits", getRoundByIDInput{RoundID: roundID}, &hits); err != nil {
		return nil, err
	}
	return &hits, nil
}

// BattleLootSummaryService handles player loot summaries for battles.
type BattleLootSummaryService struct {
	client *Client
}

// GetByBattleAndUser returns a user's loot summary for a battle.
func (s *BattleLootSummaryService) GetByBattleAndUser(ctx context.Context, battleID, userID string) (*BattleLootSummary, error) {
	var summary BattleLootSummary
	input := BattleLootSummaryByBattleAndUserFilter{BattleID: battleID, UserID: userID}
	if err := s.client.trpc.Call(ctx, "battleLootSummary.getByBattleAndUser", input, &summary); err != nil {
		return nil, err
	}
	return &summary, nil
}

// BattleRankingService handles rankings within battles and wars.
type BattleRankingService struct {
	client *Client
}

// GetRanking returns ranking entries matching the supplied filter.
func (s *BattleRankingService) GetRanking(ctx context.Context, filter BattleRankingRankingFilter) (*BattleRanking, error) {
	var ranking BattleRanking
	if err := s.client.trpc.Call(ctx, "battleRanking.getRanking", filter, &ranking); err != nil {
		return nil, err
	}
	return &ranking, nil
}
