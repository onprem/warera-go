package warera

import (
	"context"
	"iter"
)

// EventService handles paginated world and country events.
type EventService struct {
	client *Client
}

// List returns one page of events matching the filter.
func (s *EventService) List(ctx context.Context, filter EventFilter) (*CursorPage[Event], error) {
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
		return s.List(c, f)
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

// List returns one page of articles matching the filter.
func (s *ArticleService) List(ctx context.Context, filter ArticleFilter) (*CursorPage[Article], error) {
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
		return s.List(c, f)
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

// GetMany retrieves full user profiles for the supplied IDs. Results preserve
// the input order and duplicate IDs. An empty input returns an empty slice
// without making a request. If any lookup fails, GetMany returns no profiles
// and a [BatchLookupError] identifying the failed IDs. Large inputs are split
// according to the client's configured batch-size limit.
func (s *UserService) GetMany(ctx context.Context, userIDs []string) ([]User, error) {
	return getManyByID[User](ctx, s.client, "users", "user.getUserById", userIDs, func(id string) any {
		return getUserByIdInput{UserID: id}
	})
}

// GetLiteByID returns a compact user profile by ID.
func (s *UserService) GetLiteByID(ctx context.Context, userID string) (*UserLite, error) {
	var user UserLite
	err := s.client.trpc.Call(ctx, "user.getUserLite", getUserByIdInput{UserID: userID}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetLiteMany retrieves compact user profiles for the supplied IDs. Results
// preserve the input order and duplicate IDs. An empty input returns an empty
// slice without making a request. If any lookup fails, GetLiteMany returns no
// profiles and a [BatchLookupError] identifying the failed IDs. Large inputs
// are split according to the client's configured batch-size limit.
func (s *UserService) GetLiteMany(ctx context.Context, userIDs []string) ([]UserLite, error) {
	return getManyByID[UserLite](ctx, s.client, "users", "user.getUserLite", userIDs, func(id string) any {
		return getUserByIdInput{UserID: id}
	})
}

// GetSkills returns a player's skills. User.getUserLite no longer includes
// skills; this calls user.getUserSkills, which the published API description
// names but does not list.
func (s *UserService) GetSkills(ctx context.Context, userID string) (*UserSkillProfile, error) {
	var profile UserSkillProfile
	if err := s.client.trpc.Call(ctx, "user.getUserSkills", getUserByIdInput{UserID: userID}, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// GetSkillsMany retrieves the skills of the supplied players. Results preserve
// the input order and duplicate IDs. If any lookup fails, GetSkillsMany returns
// no records and a [BatchLookupError] identifying the failed IDs.
func (s *UserService) GetSkillsMany(ctx context.Context, userIDs []string) ([]UserSkillProfile, error) {
	return getManyByID[UserSkillProfile](ctx, s.client, "user skills", "user.getUserSkills", userIDs, func(id string) any {
		return getUserByIdInput{UserID: id}
	})
}

// List returns one page of the citizens of a country, newest first. Entries
// carry only the user ID and join date; load profiles with [UserService.GetLiteMany].
func (s *UserService) List(ctx context.Context, filter UserFilter) (*CursorPage[UserListItem], error) {
	var page CursorPage[UserListItem]
	err := s.client.trpc.Call(ctx, "user.getUsersByCountry", filter, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over all citizens of a country across pages.
func (s *UserService) All(ctx context.Context, filter UserFilter) iter.Seq2[UserListItem, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cur string) (*CursorPage[UserListItem], error) {
		f := filter
		f.Cursor = cur
		return s.List(c, f)
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
	return s.client.Users.GetLiteMany(ctx, result.UserIDs)
}

// --- Poll Service ---

// PollService handles article polls and their voters.
type PollService struct {
	client *Client
}

type getPollInput struct {
	PollID string `json:"pollId"`
}

// GetByID returns a poll with its results. When the client is authenticated,
// the poll also reports the caller's vote and whether they can vote.
func (s *PollService) GetByID(ctx context.Context, pollID string) (*Poll, error) {
	var poll Poll
	if err := s.client.trpc.Call(ctx, "poll.getPoll", getPollInput{PollID: pollID}, &poll); err != nil {
		return nil, err
	}
	return &poll, nil
}

// ListVoters returns one page of the votes cast for one poll option.
func (s *PollService) ListVoters(ctx context.Context, filter PollVoterFilter) (*CursorPage[PollVote], error) {
	var page CursorPage[PollVote]
	if err := s.client.trpc.Call(ctx, "poll.getPollVoters", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// AllVoters returns an iterator over every vote cast for one poll option.
func (s *PollService) AllVoters(ctx context.Context, filter PollVoterFilter) iter.Seq2[PollVote, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[PollVote], error) {
		f := filter
		f.Cursor = cursor
		return s.ListVoters(c, f)
	})
}
