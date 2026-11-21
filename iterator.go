package warera

import (
	"context"
	"fmt"
	"iter"
)

// PaginatorFunc is a function that fetches a single page given a cursor.
type PaginatorFunc[T any] func(ctx context.Context, cursor string) (*CursorPage[T], error)

// NewCursorIterator returns a Go 1.23+ iterator (iter.Seq2[T, error]) that seamlessly
// paginates through all items until no further cursor is returned or context is cancelled.
func NewCursorIterator[T any](ctx context.Context, initialCursor string, fetchPage PaginatorFunc[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		currentCursor := initialCursor
		seenCursors := map[string]struct{}{}
		if currentCursor != "" {
			seenCursors[currentCursor] = struct{}{}
		}

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
			if page == nil {
				var zero T
				yield(zero, fmt.Errorf("pagination returned a nil page"))
				return
			}

			for _, item := range page.Items {
				if !yield(item, nil) {
					return
				}
			}

			if !page.HasNext() {
				return
			}
			if _, exists := seenCursors[page.NextCursor]; exists {
				var zero T
				yield(zero, fmt.Errorf("pagination cursor repeated: %q", page.NextCursor))
				return
			}

			currentCursor = page.NextCursor
			seenCursors[currentCursor] = struct{}{}
		}
	}
}

// cursorFilter is a filter whose page cursor [Pages] can advance.
type cursorFilter[F any] interface {
	withCursor(cursor string) F
}

// cursorPage is a page that reports the cursor of the page after it.
type cursorPage interface {
	nextCursor() string
}

func (p *CursorPage[T]) nextCursor() string {
	if !p.HasNext() {
		return ""
	}
	return p.NextCursor
}

// Pages returns an iterator over whole pages from a List method, starting at
// filter's cursor. Use it instead of a service's All method when a tool must
// resume where it stopped: save each page's NextCursor after processing the
// page, and pass it back as the filter's Cursor on the next run.
//
//	for page, err := range warera.Pages(ctx, client.Events.List, warera.EventFilter{CountryID: id}) {
//		if err != nil {
//			return err
//		}
//		process(page.Items)
//		saveCursor(page.NextCursor)
//	}
//
// Pages works with every List method that takes a filter with a Cursor field.
// It stops after the last page, at the first error, and if the API repeats a
// cursor.
func Pages[F cursorFilter[F], P cursorPage](ctx context.Context, list func(context.Context, F) (P, error), filter F) iter.Seq2[P, error] {
	return func(yield func(P, error) bool) {
		var zero P
		seen := map[string]bool{}
		for {
			if err := ctx.Err(); err != nil {
				yield(zero, err)
				return
			}
			page, err := list(ctx, filter)
			if err != nil {
				yield(zero, err)
				return
			}
			if !yield(page, nil) {
				return
			}
			next := page.nextCursor()
			if next == "" {
				return
			}
			if seen[next] {
				yield(zero, fmt.Errorf("pagination cursor repeated: %q", next))
				return
			}
			seen[next] = true
			filter = filter.withCursor(next)
		}
	}
}

// Collect gathers the values from an iterator such as a service's All method,
// stopping after limit values when limit is positive. On an error it returns
// the values collected so far together with the error.
func Collect[T any](seq iter.Seq2[T, error], limit int) ([]T, error) {
	values := make([]T, 0)
	for value, err := range seq {
		if err != nil {
			return values, err
		}
		values = append(values, value)
		if limit > 0 && len(values) >= limit {
			break
		}
	}
	return values, nil
}

// Until passes on values from seq until stop reports true for one, which it
// does not pass on, and then stops without fetching further pages. Listings
// such as events and transactions come newest first, so Until fetches only
// what is new since a tool's last poll:
//
//	newTrades := warera.Until(client.Transactions.All(ctx, filter), func(tx warera.TransactionListItem) bool {
//		return tx.ID == lastSeenID
//	})
func Until[T any](seq iter.Seq2[T, error], stop func(T) bool) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for value, err := range seq {
			if err == nil && stop(value) {
				return
			}
			if !yield(value, err) || err != nil {
				return
			}
		}
	}
}

// Cursor setters for the hand-written filters, used by [Pages].

func (f PartyFilter) withCursor(cursor string) PartyFilter               { f.Cursor = cursor; return f }
func (f ElectionFilter) withCursor(cursor string) ElectionFilter         { f.Cursor = cursor; return f }
func (f DonationFilter) withCursor(cursor string) DonationFilter         { f.Cursor = cursor; return f }
func (f ManagerOfferFilter) withCursor(cursor string) ManagerOfferFilter { f.Cursor = cursor; return f }
func (f MUSponsorshipFilter) withCursor(cursor string) MUSponsorshipFilter {
	f.Cursor = cursor
	return f
}
func (f StateArmyLedgerFilter) withCursor(cursor string) StateArmyLedgerFilter {
	f.Cursor = cursor
	return f
}
func (f TournamentFilter) withCursor(cursor string) TournamentFilter { f.Cursor = cursor; return f }
