package warera

import "context"

// getManyByID looks up each distinct ID with one procedure call per ID, sent
// in as few batch requests as the client's batch size allows. It returns the
// records it found keyed by ID. When some lookups fail it returns those
// records together with a *BatchLookupError naming only the failed IDs.
func getManyByID[T any](
	ctx context.Context,
	client *Client,
	entity string,
	procedure string,
	ids []string,
	inputForID func(string) any,
) (map[string]T, error) {
	unique := uniqueIDs(ids)
	results := make([]T, len(unique))
	calls := make([]BatchCall, len(unique))
	for i, id := range unique {
		calls[i] = BatchCall{Method: procedure, Input: inputForID(id), Output: &results[i]}
	}
	// Batch records every failure, including transport failures, on the
	// individual calls, so its combined error adds nothing here.
	_ = client.Batch(ctx, calls)

	found := make(map[string]T, len(unique))
	var failures []BatchLookupFailure
	for i, call := range calls {
		if call.Err != nil {
			failures = append(failures, BatchLookupFailure{ID: unique[i], Err: call.Err})
		} else {
			found[unique[i]] = results[i]
		}
	}
	if len(failures) > 0 {
		return found, &BatchLookupError{Entity: entity, Failures: failures}
	}
	return found, nil
}

// uniqueIDs returns ids without empty strings and duplicates, in first-seen
// order.
func uniqueIDs(ids []string) []string {
	seen := make(map[string]bool, len(ids))
	unique := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != "" && !seen[id] {
			seen[id] = true
			unique = append(unique, id)
		}
	}
	return unique
}
