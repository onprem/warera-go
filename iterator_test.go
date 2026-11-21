package warera

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestNewCursorIteratorFetchesPagesInOrder(t *testing.T) {
	var cursors []string
	iterator := NewCursorIterator(context.Background(), "start", func(_ context.Context, cursor string) (*CursorPage[int], error) {
		cursors = append(cursors, cursor)
		switch cursor {
		case "start":
			return &CursorPage[int]{Items: []int{1, 2}, NextCursor: "next"}, nil
		case "next":
			return &CursorPage[int]{Items: []int{3}}, nil
		default:
			return nil, errors.New("unexpected cursor: " + cursor)
		}
	})

	var got []int
	for item, err := range iterator {
		if err != nil {
			t.Fatalf("iterator returned error: %v", err)
		}
		got = append(got, item)
	}
	if !equalInts(got, []int{1, 2, 3}) {
		t.Fatalf("items = %v, want [1 2 3]", got)
	}
	if !equalStrings(cursors, []string{"start", "next"}) {
		t.Fatalf("fetched cursors = %v, want [start next]", cursors)
	}
}

func TestNewCursorIteratorStopsWhenConsumerStops(t *testing.T) {
	fetches := 0
	iterator := NewCursorIterator(context.Background(), "", func(context.Context, string) (*CursorPage[int], error) {
		fetches++
		return &CursorPage[int]{Items: []int{1, 2}, NextCursor: "next"}, nil
	})

	items := 0
	for item, err := range iterator {
		if err != nil {
			t.Fatalf("iterator returned error: %v", err)
		}
		if item != 1 {
			t.Fatalf("first yielded item = %d, want 1", item)
		}
		items++
		break
	}
	if items != 1 || fetches != 1 {
		t.Fatalf("yielded %d item(s) and fetched %d page(s), want 1 each", items, fetches)
	}
}

func TestNewCursorIteratorReturnsFetchAndNilPageErrors(t *testing.T) {
	wantErr := errors.New("fetch failed")
	tests := []struct {
		name  string
		fetch PaginatorFunc[int]
		want  string
	}{
		{
			name: "fetch error",
			fetch: func(context.Context, string) (*CursorPage[int], error) {
				return nil, wantErr
			},
			want: wantErr.Error(),
		},
		{
			name: "nil page",
			fetch: func(context.Context, string) (*CursorPage[int], error) {
				return nil, nil
			},
			want: "pagination returned a nil page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorCount := 0
			for _, err := range NewCursorIterator(context.Background(), "", tt.fetch) {
				if err != nil {
					errorCount++
					if err.Error() != tt.want {
						t.Fatalf("iterator error = %v, want %q", err, tt.want)
					}
				}
			}
			if errorCount != 1 {
				t.Fatalf("yielded %d errors, want 1", errorCount)
			}
		})
	}
}

func TestNewCursorIteratorDetectsRepeatedCursor(t *testing.T) {
	fetches := 0
	iterator := NewCursorIterator(context.Background(), "", func(_ context.Context, cursor string) (*CursorPage[int], error) {
		fetches++
		if cursor == "" {
			return &CursorPage[int]{Items: []int{1}, NextCursor: "repeat"}, nil
		}
		return &CursorPage[int]{Items: []int{2}, NextCursor: "repeat"}, nil
	})

	var items []int
	var gotErr error
	for item, err := range iterator {
		if err != nil {
			gotErr = err
			continue
		}
		items = append(items, item)
	}
	if !equalInts(items, []int{1, 2}) {
		t.Fatalf("items = %v, want [1 2]", items)
	}
	if gotErr == nil || !strings.Contains(gotErr.Error(), `pagination cursor repeated: "repeat"`) {
		t.Fatalf("iterator error = %v, want repeated cursor error", gotErr)
	}
	if fetches != 2 {
		t.Fatalf("fetch count = %d, want 2", fetches)
	}
}

func TestNewCursorIteratorHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fetches := 0
	iterator := NewCursorIterator(ctx, "", func(context.Context, string) (*CursorPage[int], error) {
		fetches++
		return &CursorPage[int]{Items: []int{1}}, nil
	})

	var gotErr error
	for _, err := range iterator {
		if err != nil {
			gotErr = err
		}
	}
	if !errors.Is(gotErr, context.Canceled) {
		t.Fatalf("iterator error = %v, want context.Canceled", gotErr)
	}
	if fetches != 0 {
		t.Fatalf("fetch count = %d, want 0", fetches)
	}
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestPagesFollowsCursorsAndResumes(t *testing.T) {
	pages := map[string]*CursorPage[Event]{
		"":   {Items: []Event{{ID: "e1"}}, NextCursor: "c2"},
		"c2": {Items: []Event{{ID: "e2"}}, NextCursor: "c3"},
		"c3": {Items: []Event{{ID: "e3"}}},
	}
	var requested []EventFilter
	list := func(_ context.Context, f EventFilter) (*CursorPage[Event], error) {
		requested = append(requested, f)
		return pages[f.Cursor], nil
	}

	var ids []string
	for page, err := range Pages(context.Background(), list, EventFilter{CountryID: "x"}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, page.Items[0].ID)
	}
	if strings.Join(ids, ",") != "e1,e2,e3" || len(requested) != 3 || requested[2].Cursor != "c3" || requested[2].CountryID != "x" {
		t.Fatalf("pages = %v, requests = %+v; want e1,e2,e3 with the filter kept and cursors advanced", ids, requested)
	}

	// A saved cursor resumes from that page.
	ids = nil
	for page, err := range Pages(context.Background(), list, EventFilter{Cursor: "c2"}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, page.Items[0].ID)
	}
	if strings.Join(ids, ",") != "e2,e3" {
		t.Fatalf("resumed pages = %v, want e2,e3", ids)
	}
}

func TestUntilStopsWithoutFetchingMorePages(t *testing.T) {
	fetches := 0
	seq := NewCursorIterator(context.Background(), "", func(_ context.Context, cursor string) (*CursorPage[string], error) {
		fetches++
		if cursor == "" {
			return &CursorPage[string]{Items: []string{"new2", "new1", "seen"}, NextCursor: "older"}, nil
		}
		return &CursorPage[string]{Items: []string{"old"}}, nil
	})
	got, err := Collect(Until(seq, func(id string) bool { return id == "seen" }), 0)
	if err != nil || strings.Join(got, ",") != "new2,new1" || fetches != 1 {
		t.Fatalf("Until = %v, %v after %d fetches; want new2,new1 after 1", got, err, fetches)
	}
}

// Every List method that takes a cursor filter must work with Pages; this
// fails to compile when a filter lacks a cursor setter.
var _ = func(c *Client, ctx context.Context) {
	_ = Pages(ctx, c.Articles.List, ArticleFilter{})
	_ = Pages(ctx, c.Battles.List, BattleFilter{})
	_ = Pages(ctx, c.BattleRankings.List, BattleRankingFilter{})
	_ = Pages(ctx, c.Companies.List, CompanyFilter{})
	_ = Pages(ctx, c.Donations.List, DonationFilter{})
	_ = Pages(ctx, c.Elections.List, ElectionFilter{})
	_ = Pages(ctx, c.Events.List, EventFilter{})
	_ = Pages(ctx, c.Managers.ListOffers, ManagerOfferFilter{})
	_ = Pages(ctx, c.MercenaryAuctions.List, MercenaryAuctionFilter{})
	_ = Pages(ctx, c.MilitaryUnits.List, MilitaryUnitFilter{})
	_ = Pages(ctx, c.MilitaryUnits.ListActiveSponsorships, MUSponsorshipFilter{})
	_ = Pages(ctx, c.Parties.List, PartyFilter{})
	_ = Pages(ctx, c.Polls.ListVoters, PollVoterFilter{})
	_ = Pages(ctx, c.Rankings.List, RankingPageFilter{})
	_ = Pages(ctx, c.StateArmy.ListLedger, StateArmyLedgerFilter{})
	_ = Pages(ctx, c.Tournaments.List, TournamentFilter{})
	_ = Pages(ctx, c.Transactions.List, TransactionFilter{})
	_ = Pages(ctx, c.Users.List, UserFilter{})
}
