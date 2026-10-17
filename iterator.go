package warera

import (
	"context"
	"iter"
)

// PaginatorFunc is a function that fetches a single page given a cursor.
type PaginatorFunc[T any] func(ctx context.Context, cursor string) (*CursorPage[T], error)

// NewCursorIterator returns a Go 1.23+ iterator (iter.Seq2[T, error]) that seamlessly
// paginates through all items until no further cursor is returned or context is cancelled.
func NewCursorIterator[T any](ctx context.Context, initialCursor string, fetchPage PaginatorFunc[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		currentCursor := initialCursor

		for {
			select {
			case <-ctx.Done():
				var zero T
				yield(zero, ctx.Err())
				return
			default:
			}

			page, err := fetchPage(ctx, currentCursor)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}

			for _, item := range page.Items {
				if !yield(item, nil) {
					return
				}
			}

			if !page.HasNext() || len(page.Items) == 0 {
				return
			}

			currentCursor = page.NextCursor
		}
	}
}
