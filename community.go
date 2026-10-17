package warera

import (
	"context"
	"iter"
)

// EventService handles paginated world and country events.
type EventService struct {
	client *Client
}

// GetEventsPaginated returns one page of events matching the filter.
func (s *EventService) GetEventsPaginated(ctx context.Context, filter EventFilter) (*CursorPage[Event], error) {
	var page CursorPage[Event]
	if err := s.client.trpc.Call(ctx, "event.getEventsPaginated", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over all events matching the filter.
func (s *EventService) All(ctx context.Context, filter EventFilter) iter.Seq2[Event, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[Event], error) {
		f := filter
		f.Cursor = cursor
		return s.GetEventsPaginated(c, f)
	})
}

// --- Article Service ---

// ArticleService handles interactions with in-game news and articles.
type ArticleService struct {
	client *Client
}

type getArticleByIdInput struct {
	ArticleID string `json:"articleId"`
}

// GetByID returns full article information by ID.
func (s *ArticleService) GetByID(ctx context.Context, articleID string) (*Article, error) {
	var article Article
	err := s.client.trpc.Call(ctx, "article.getArticleById", getArticleByIdInput{ArticleID: articleID}, &article)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetLiteByID returns lightweight article information by ID.
func (s *ArticleService) GetLiteByID(ctx context.Context, articleID string) (*Article, error) {
	var article Article
	err := s.client.trpc.Call(ctx, "article.getArticleLiteById", getArticleByIdInput{ArticleID: articleID}, &article)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetArticlesPaginated returns a single paginated page of articles matching the filter.
func (s *ArticleService) GetArticlesPaginated(ctx context.Context, filter ArticleFilter) (*CursorPage[Article], error) {
	var page CursorPage[Article]
	err := s.client.trpc.Call(ctx, "article.getArticlesPaginated", filter, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over all articles across pages matching the filter.
func (s *ArticleService) All(ctx context.Context, filter ArticleFilter) iter.Seq2[Article, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cur string) (*CursorPage[Article], error) {
		f := filter
		f.Cursor = cur
		return s.GetArticlesPaginated(c, f)
	})
}

// --- User Service ---

// UserService handles player profiles and user data.
type UserService struct {
	client *Client
}

type getUserByIdInput struct {
	UserID string `json:"userId"`
}

// GetByID returns the full user profile.
func (s *UserService) GetByID(ctx context.Context, userID string) (*User, error) {
	var user User
	err := s.client.trpc.Call(ctx, "user.getUserById", getUserByIdInput{UserID: userID}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetLite returns a compact user profile.
func (s *UserService) GetLite(ctx context.Context, userID string) (*UserLite, error) {
	var user UserLite
	err := s.client.trpc.Call(ctx, "user.getUserLite", getUserByIdInput{UserID: userID}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsersByCountry returns a page of users belonging to a specific country.
func (s *UserService) GetUsersByCountry(ctx context.Context, filter UserFilter) (*CursorPage[User], error) {
	var page CursorPage[User]
	err := s.client.trpc.Call(ctx, "user.getUsersByCountry", filter, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// AllByCountry returns an iterator over all citizens of a country across pages.
func (s *UserService) AllByCountry(ctx context.Context, filter UserFilter) iter.Seq2[User, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cur string) (*CursorPage[User], error) {
		f := filter
		f.Cursor = cur
		return s.GetUsersByCountry(c, f)
	})
}

// --- Search Service ---

// SearchService handles global searches across the WarEra game universe.
type SearchService struct {
	client *Client
}

type searchAnythingInput struct {
	SearchText string `json:"searchText"`
}

// SearchAnything searches for players, companies, countries, regions, and military units,
// returning the IDs of matching entities.
func (s *SearchService) SearchAnything(ctx context.Context, searchText string) (*SearchResult, error) {
	var result SearchResult
	err := s.client.trpc.Call(ctx, "search.searchAnything", searchAnythingInput{SearchText: searchText}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchUsers searches for users matching the query and fetches their profiles.
func (s *SearchService) SearchUsers(ctx context.Context, searchText string) ([]UserLite, error) {
	result, err := s.SearchAnything(ctx, searchText)
	if err != nil {
		return nil, err
	}
	if len(result.UserIDs) == 0 {
		return nil, nil
	}

	users := make([]UserLite, len(result.UserIDs))
	calls := make([]BatchCall, len(result.UserIDs))
	for i, id := range result.UserIDs {
		calls[i] = BatchCall{
			Method: "user.getUserLite",
			Input:  getUserByIdInput{UserID: id},
			Output: &users[i],
		}
	}
	if err := s.client.Batch(ctx, calls); err != nil {
		return nil, err
	}

	return users, nil
}
