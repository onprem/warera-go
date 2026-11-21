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
// yielded by the iterator. [Pages] iterates whole pages so a collector can
// save its cursor and resume, [Collect] gathers an iterator's values, and
// [Until] stops at a record already seen, so polling the newest-first event
// and transaction listings fetches only what is new.
//
// GetMany methods look up many IDs in batched requests and return the records
// found keyed by ID. [Directory] caches countries, regions, the game
// configuration, and user profiles for turning IDs into names. To test a tool
// without the network, point a client at the fake server in the waretest
// package.
//
// Concurrent calls made through a client are batched automatically, and the
// API counts a batch as one request against its rate limit. The client keeps
// to the API's limit (500 requests a minute with a key, 100 without), waits
// when the server reports its quota exhausted, and retries rate-limited and
// transient failures after the delay the server asks for. Use [WithHooks] to
// log or measure requests, retries, and waits. Use
// [Client.Batch] to submit an explicit batch, or [Client.RawCall] to invoke a
// procedure without a typed service method.
//
// Most services wrap procedures from WarEra's published OpenAPI document. Some
// wrap procedures the API serves but does not list there: the Employees,
// Managers, Alliances, Wars, Parties, Elections, MarketStats, Donations,
// Tournaments, StateArmy, Contributions, and Leagues services, plus
// Users.GetSkills, Companies.GetProductionBonus, and the methods in
// politics.go, market.go, and military_units.go. They may change without
// notice and need an API key.
//
// API errors are available as [TRPCError] and transport failures as
// [HTTPError]; use [errors.As] to inspect their fields. Bulk lookup methods
// report failed entity IDs with [BatchLookupError]. [IsNotFound],
// [IsRateLimited], and [IsTemporary] classify any of these errors, so callers
// can tell a missing entity from a failure worth retrying.
package warera
