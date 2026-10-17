package warera

import (
	"context"
	"iter"
)

// UserService handles player profiles and user data.
type UserService struct {
	client *Client
}

type getByIdInput struct {
	UserID string `json:"userId"`
}

// GetByID returns the full user profile.
func (s *UserService) GetByID(ctx context.Context, userID string) (*User, error) {
	var user User
	err := s.client.trpc.Call(ctx, "user.getUserById", getByIdInput{UserID: userID}, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetLite returns a compact user profile.
func (s *UserService) GetLite(ctx context.Context, userID string) (*UserLite, error) {
	var user UserLite
	err := s.client.trpc.Call(ctx, "user.getUserLite", getByIdInput{UserID: userID}, &user)
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
