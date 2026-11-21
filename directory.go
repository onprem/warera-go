package warera

import (
	"context"
	"sync"
	"time"
)

// Directory caches reference data that tools look up constantly and that
// changes slowly: countries, regions, the game configuration, and compact user
// profiles. Use it to turn the IDs in API records into names without fetching
// the same data again. Entries are refreshed after the directory's TTL. A
// Directory is safe for concurrent use, and concurrent callers share one fetch.
//
//	dir := warera.NewDirectory(client, 10*time.Minute)
//	country, err := dir.Country(ctx, battle.Attacker.CountryID)
type Directory struct {
	client *Client
	ttl    time.Duration

	countries cached[map[string]Country]
	regions   cached[map[string]Region]
	config    cached[*GameConfig]

	usersMu sync.Mutex
	users   map[string]cachedUser
}

type cachedUser struct {
	user      UserLite
	fetchedAt time.Time
}

// NewDirectory returns a Directory that loads data through client and keeps
// it for ttl.
func NewDirectory(client *Client, ttl time.Duration) *Directory {
	return &Directory{client: client, ttl: ttl, users: map[string]cachedUser{}}
}

// Countries returns every country keyed by ID.
func (d *Directory) Countries(ctx context.Context) (map[string]Country, error) {
	return d.countries.get(ctx, d.ttl, func(ctx context.Context) (map[string]Country, error) {
		list, err := d.client.Countries.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		byID := make(map[string]Country, len(list))
		for _, country := range list {
			byID[country.ID] = country
		}
		return byID, nil
	})
}

// Country returns the country with the given ID, or [ErrNotFound].
func (d *Directory) Country(ctx context.Context, countryID string) (Country, error) {
	return lookup(ctx, d.Countries, countryID)
}

// Regions returns every region keyed by ID.
func (d *Directory) Regions(ctx context.Context) (map[string]Region, error) {
	return d.regions.get(ctx, d.ttl, d.client.Regions.GetAll)
}

// Region returns the region with the given ID, or [ErrNotFound].
func (d *Directory) Region(ctx context.Context, regionID string) (Region, error) {
	return lookup(ctx, d.Regions, regionID)
}

// GameConfig returns the game configuration, which holds item, skill, and
// upgrade definitions.
func (d *Directory) GameConfig(ctx context.Context) (*GameConfig, error) {
	return d.config.get(ctx, d.ttl, d.client.GameConfig.Get)
}

// Users returns compact profiles for the given user IDs keyed by ID, fetching
// only those not already cached. Like [UserService.GetLiteMany], it returns
// the profiles it has together with a [BatchLookupError] when some fail.
func (d *Directory) Users(ctx context.Context, userIDs []string) (map[string]UserLite, error) {
	now := time.Now()
	found := make(map[string]UserLite, len(userIDs))
	var missing []string
	d.usersMu.Lock()
	for _, id := range uniqueIDs(userIDs) {
		if entry, ok := d.users[id]; ok && now.Sub(entry.fetchedAt) < d.ttl {
			found[id] = entry.user
		} else {
			missing = append(missing, id)
		}
	}
	d.usersMu.Unlock()
	if len(missing) == 0 {
		return found, nil
	}

	fetched, err := d.client.Users.GetLiteMany(ctx, missing)
	d.usersMu.Lock()
	for id, user := range fetched {
		d.users[id] = cachedUser{user: user, fetchedAt: now}
		found[id] = user
	}
	d.usersMu.Unlock()
	return found, err
}

// cached holds one value and when it was loaded. Its mutex is held while
// loading, so concurrent callers wait for one fetch instead of repeating it.
type cached[T any] struct {
	mu       sync.Mutex
	value    T
	loadedAt time.Time
}

func (c *cached[T]) get(ctx context.Context, ttl time.Duration, load func(context.Context) (T, error)) (T, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.loadedAt.IsZero() && time.Since(c.loadedAt) < ttl {
		return c.value, nil
	}
	value, err := load(ctx)
	if err != nil {
		var zero T
		return zero, err
	}
	c.value, c.loadedAt = value, time.Now()
	return value, nil
}

func lookup[T any](ctx context.Context, all func(context.Context) (map[string]T, error), id string) (T, error) {
	byID, err := all(ctx)
	if err != nil {
		var zero T
		return zero, err
	}
	value, ok := byID[id]
	if !ok {
		return value, ErrNotFound
	}
	return value, nil
}
