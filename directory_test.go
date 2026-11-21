package warera_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/waretest"
)

func TestDirectoryCachesUntilTTL(t *testing.T) {
	srv := waretest.NewServer(t)
	client := warera.NewClient(warera.WithBaseURL(srv.URL))
	defer client.Close()
	ctx := context.Background()
	dir := warera.NewDirectory(client, 200*time.Millisecond)

	countries, err := dir.Countries(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var someID string
	for id := range countries {
		someID = id
		break
	}

	// Concurrent lookups within the TTL share the cached list.
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if country, err := dir.Country(ctx, someID); err != nil || country.ID != someID {
				t.Errorf("Country = %+v, %v", country, err)
			}
		}()
	}
	wg.Wait()
	if _, err := dir.Country(ctx, "missing"); !errors.Is(err, warera.ErrNotFound) {
		t.Errorf("unknown country error = %v, want ErrNotFound", err)
	}
	if calls := srv.Calls("country.getAllCountries"); calls != 1 {
		t.Fatalf("fetches within TTL = %d, want 1", calls)
	}

	time.Sleep(250 * time.Millisecond)
	if _, err := dir.Country(ctx, someID); err != nil {
		t.Fatal(err)
	}
	if calls := srv.Calls("country.getAllCountries"); calls != 2 {
		t.Fatalf("fetches after TTL = %d, want 2", calls)
	}
}

func TestDirectoryUsersFetchesOnlyMissingIDs(t *testing.T) {
	srv := waretest.NewServer(t)
	srv.Handle("user.getUserLite", func(input json.RawMessage) (any, error) {
		var in struct {
			UserID string `json:"userId"`
		}
		_ = json.Unmarshal(input, &in)
		return map[string]string{"_id": in.UserID, "username": "user-" + in.UserID}, nil
	})
	client := warera.NewClient(warera.WithBaseURL(srv.URL))
	defer client.Close()
	ctx := context.Background()
	dir := warera.NewDirectory(client, time.Hour)

	if _, err := dir.Users(ctx, []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	users, err := dir.Users(ctx, []string{"a", "c"})
	if err != nil || users["a"].Username != "user-a" || users["c"].Username != "user-c" {
		t.Fatalf("Users = %+v, %v", users, err)
	}
	if calls := srv.Calls("user.getUserLite"); calls != 3 {
		t.Errorf("lookups = %d, want 3 (a, b, then only c)", calls)
	}
}
