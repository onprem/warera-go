package warera

import "context"

func getManyByID[T any](
	ctx context.Context,
	client *Client,
	entity string,
	procedure string,
	ids []string,
	inputForID func(string) any,
) ([]T, error) {
	results := make([]T, len(ids))
	if len(ids) == 0 {
		return results, nil
	}

	calls := make([]BatchCall, len(ids))
	for i, id := range ids {
		calls[i] = BatchCall{
			Method: procedure,
			Input:  inputForID(id),
			Output: &results[i],
		}
	}

	batchErr := client.Batch(ctx, calls)
	failures := make([]BatchLookupFailure, 0)
	for i, call := range calls {
		if call.Err != nil {
			failures = append(failures, BatchLookupFailure{ID: ids[i], Err: call.Err})
		}
	}
	if batchErr != nil && len(failures) == 0 {
		// Preserve a useful ID list if a transport-level failure was not copied
		// onto individual calls by the batch implementation.
		for _, id := range ids {
			failures = append(failures, BatchLookupFailure{ID: id, Err: batchErr})
		}
	}
	if len(failures) > 0 {
		return nil, &BatchLookupError{Entity: entity, Failures: failures}
	}
	return results, nil
}
