// Package warera provides a Go client for the WarEra online game API.
//
// Create a client with [NewClient] and close it when finished. Supply
// [WithAPIKey] when an endpoint requires authentication:
//
//	client := warera.NewClient(warera.WithAPIKey(apiKey))
//	defer client.Close()
//
// Service methods accept a [context.Context] and return typed response models.
// Pass a context with a deadline or cancellation signal when the request should
// be bounded by the caller.
// Cursor-paginated services provide List to fetch a single page and All to
// lazily iterate across pages. All returns an [iter.Seq2], so check each error
// yielded by the iterator.
//
// Concurrent calls made through a client are batched automatically. Use
// [Client.Batch] to submit an explicit batch, or [Client.RawCall] to invoke a
// procedure without a typed service method.
//
// Most services wrap procedures from WarEra's published OpenAPI document. The
// Employees and Managers services, Users.GetSkills, and
// Companies.GetProductionBonus wrap procedures the API serves but does not
// list there; they may change without notice and need an API key.
//
// API errors are available as [TRPCError] and transport failures as
// [HTTPError]; use [errors.As] to inspect their fields. Bulk lookup methods
// report failed entity IDs with [BatchLookupError].
package warera
