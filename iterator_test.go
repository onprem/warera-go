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
