package warera

import (
	"context"
	"iter"
)

// RankingService handles global game rankings.
type RankingService struct {
	client *Client
}

type getRankingInput struct {
	RankingType RankingType `json:"rankingType"`
}

// GetRanking returns the tier cutoffs and the top 100 entries of a ranking.
// Entries with a value of 0 are not listed. Use [RankingService.List] or
// [RankingService.All] for entries beyond the top 100.
func (s *RankingService) GetRanking(ctx context.Context, rankingType RankingType) (*Ranking, error) {
	var ranking Ranking
	if err := s.client.trpc.Call(ctx, "ranking.getRanking", getRankingInput{RankingType: rankingType}, &ranking); err != nil {
		return nil, err
	}
	return &ranking, nil
}

// List returns one page of ranking entries, best first. Only the first page
// includes [RankingPage.Ranking], the ranking's tier cutoffs and entry count.
func (s *RankingService) List(ctx context.Context, filter RankingPageFilter) (*RankingPage, error) {
	var page RankingPage
	if err := s.client.trpc.Call(ctx, "ranking.getRankingPage", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over every entry of a ranking, best first.
func (s *RankingService) All(ctx context.Context, filter RankingPageFilter) iter.Seq2[RankingEntry, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[RankingEntry], error) {
		f := filter
		f.Cursor = cursor
		page, err := s.List(c, f)
		if err != nil {
			return nil, err
		}
		return &page.CursorPage, nil
	})
}

// GetAround returns one user, country, military unit, or alliance in a ranking
// with the entries just above and below it, using the same country or
// military-unit filter as [RankingService.List]. It returns an empty slice when
// the entity has no entry in that ranking.
func (s *RankingService) GetAround(ctx context.Context, filter RankingAroundFilter) ([]RankingEntry, error) {
	var entries []RankingEntry
	if err := s.client.trpc.Call(ctx, "ranking.getRankingAround", filter, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
