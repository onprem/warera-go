package waretest_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/waretest"
)

func TestServerAnswersWithRecordedResponses(t *testing.T) {
	srv := waretest.NewServer(t)
	client := warera.NewClient(warera.WithBaseURL(srv.URL))
	defer client.Close()
	ctx := context.Background()

	countries, err := client.Countries.GetAll(ctx)
	if err != nil || len(countries) < 100 {
		t.Fatalf("GetAll = %d countries, %v; want the recorded list", len(countries), err)
	}
	// Void-input procedures and procedures outside the OpenAPI document are
	// recorded too.
	if value, err := client.MarketStats.GetPPMoneyValue(ctx); err != nil || value <= 0 {
		t.Fatalf("GetPPMoneyValue = %v, %v; want the recorded value", value, err)
	}
	if srv.Calls("country.getAllCountries") != 1 {
		t.Errorf("Calls = %d, want 1", srv.Calls("country.getAllCountries"))
	}
}

func TestServerHandlersShapeBatchesAndErrors(t *testing.T) {
	srv := waretest.NewServer(t)
	srv.Handle("user.getUserLite", func(input json.RawMessage) (any, error) {
		var in struct {
			UserID string `json:"userId"`
		}
		if err := json.Unmarshal(input, &in); err != nil {
			return nil, err
		}
		if strings.HasPrefix(in.UserID, "gone") {
			return nil, waretest.NotFound
		}
		return map[string]string{"_id": in.UserID, "username": "user-" + in.UserID}, nil
	})
	srv.Handle("tournament.getById", waretest.Respond(nil))
	client := warera.NewClient(warera.WithBaseURL(srv.URL), warera.WithMaxRetries(0))
	defer client.Close()
	ctx := context.Background()

	// A batch with one success and one failure: the client keeps the success.
	users, err := client.Users.GetLiteMany(ctx, []string{"a", "gone"})
	if users["a"].Username != "user-a" || len(users) != 1 || !warera.IsNotFound(err) {
		t.Fatalf("GetLiteMany = %+v, %v; want user a and a not-found failure", users, err)
	}
	// A batch where every call fails gets per-call tRPC errors.
	_, err = client.Users.GetLiteMany(ctx, []string{"gone", "gone-too"})
	var lookupErr *warera.BatchLookupError
	if !errors.As(err, &lookupErr) || len(lookupErr.Failures) != 2 {
		t.Fatalf("all-failing GetLiteMany error = %v, want two failures", err)
	}

	if _, err := client.Tournaments.GetByID(ctx, "x"); !errors.Is(err, warera.ErrNotFound) {
		t.Errorf("null result error = %v, want ErrNotFound", err)
	}
	if err := client.RawCall(ctx, "nope.unknown", nil, nil); !warera.IsNotFound(err) {
		t.Errorf("unknown procedure error = %v, want NOT_FOUND", err)
	}
	if srv.Calls("user.getUserLite") != 4 {
		t.Errorf("Calls = %d, want 4", srv.Calls("user.getUserLite"))
	}
}
