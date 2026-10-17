package warera

import "context"

// RankingService handles global game rankings.
type RankingService struct {
	client *Client
}

type getRankingInput struct {
	RankingType string `json:"rankingType"`
}

// GetRanking returns the ranking for the requested ranking type.
func (s *RankingService) GetRanking(ctx context.Context, rankingType string) (*Ranking, error) {
	var ranking Ranking
	if err := s.client.trpc.Call(ctx, "ranking.getRanking", getRankingInput{RankingType: rankingType}, &ranking); err != nil {
		return nil, err
	}
	return &ranking, nil
}
