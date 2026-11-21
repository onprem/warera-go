package warera_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/onprem/warera-go"
)

type errorResponseTransport func(*http.Request) (*http.Response, error)

func (f errorResponseTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestPublicTRPCErrorCanBeInspected(t *testing.T) {
	client := warera.NewClient(
		warera.WithBaseURL("https://example.test/trpc"),
		warera.WithHTTPClient(&http.Client{Transport: errorResponseTransport(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     "400 Bad Request",
				Body: io.NopCloser(strings.NewReader(
					`{"error":{"message":"country not found","code":-32004,"data":{"code":"NOT_FOUND","httpStatus":400}}}`,
				)),
			}, nil
		})}),
		warera.WithRateLimit(0),
		warera.WithMaxRetries(0),
	)
	defer client.Close()

	err := client.RawCall(context.Background(), "country.getCountryById", map[string]string{"countryId": "missing"}, &struct{}{})
	var apiErr *warera.TRPCError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %T (%v), want *warera.TRPCError", err, err)
	}
	if apiErr.Data == nil || apiErr.Data.Code != "NOT_FOUND" || apiErr.Data.HTTPStatus != http.StatusBadRequest {
		t.Fatalf("TRPC error data = %#v, want NOT_FOUND / 400", apiErr.Data)
	}
}

func TestPublicHTTPErrorCanBeInspected(t *testing.T) {
	client := warera.NewClient(
		warera.WithBaseURL("https://example.test/trpc"),
		warera.WithHTTPClient(&http.Client{Transport: errorResponseTransport(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     "400 Bad Request",
				Body:       io.NopCloser(strings.NewReader("gateway rejected request")),
			}, nil
		})}),
		warera.WithRateLimit(0),
		warera.WithMaxRetries(0),
	)
	defer client.Close()

	err := client.RawCall(context.Background(), "country.getCountryById", map[string]string{"countryId": "missing"}, &struct{}{})
	var httpErr *warera.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("error = %T (%v), want *warera.HTTPError", err, err)
	}
	if httpErr.StatusCode != http.StatusBadRequest || string(httpErr.Body) != "gateway rejected request" {
		t.Fatalf("HTTP error = %#v, want status 400 and response body", httpErr)
	}
}

func clientReturning(status int, body string) *warera.Client {
	return warera.NewClient(
		warera.WithBaseURL("https://example.test/trpc"),
		warera.WithHTTPClient(&http.Client{Transport: errorResponseTransport(func(*http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: status, Status: http.StatusText(status), Body: io.NopCloser(strings.NewReader(body))}, nil
		})}),
		warera.WithRateLimit(0),
		warera.WithMaxRetries(0),
	)
}

// When every call in a batch fails the same way, tRPC answers with that HTTP
// status and one error per call. Each caller must still get its own tRPC error.
func TestBatchWhereEveryCallFailsReturnsPerCallTRPCErrors(t *testing.T) {
	client := clientReturning(http.StatusNotFound, `[`+
		`{"error":{"message":"user a missing","code":-32004,"data":{"code":"NOT_FOUND","httpStatus":404}}},`+
		`{"error":{"message":"user b missing","code":-32004,"data":{"code":"NOT_FOUND","httpStatus":404}}}]`)
	defer client.Close()

	_, err := client.Users.GetLiteMany(context.Background(), []string{"a", "b"})
	var lookupErr *warera.BatchLookupError
	if !errors.As(err, &lookupErr) || len(lookupErr.Failures) != 2 {
		t.Fatalf("error = %v, want BatchLookupError with two failures", err)
	}
	for i, id := range []string{"a", "b"} {
		failure := lookupErr.Failures[i]
		var apiErr *warera.TRPCError
		if failure.ID != id || !errors.As(failure.Err, &apiErr) || apiErr.Message != "user "+id+" missing" {
			t.Errorf("failure %d = %s: %v, want %s with its own tRPC error", i, failure.ID, failure.Err, id)
		}
	}
	if !warera.IsNotFound(err) {
		t.Error("IsNotFound = false, want true")
	}
}

func TestNullResultIsNotFoundUnlessNoneIsDocumented(t *testing.T) {
	client := clientReturning(http.StatusOK, `{"result":{"data":null}}`)
	defer client.Close()

	tournament, err := client.Tournaments.GetByID(context.Background(), "missing")
	if tournament != nil || !errors.Is(err, warera.ErrNotFound) || !warera.IsNotFound(err) {
		t.Fatalf("GetByID = %v, %v; want nil, ErrNotFound", tournament, err)
	}
	// Managers.GetByCompanyID documents null as "no manager".
	manager, err := client.Managers.GetByCompanyID(context.Background(), "company")
	if manager != nil || err != nil {
		t.Fatalf("GetByCompanyID = %v, %v; want nil, nil", manager, err)
	}
}

func TestErrorClassification(t *testing.T) {
	trpcErr := func(code string) error {
		return &warera.TRPCError{Message: code, Data: &warera.TRPCErrorData{Code: code}}
	}
	tests := []struct {
		name                             string
		err                              error
		notFound, rateLimited, temporary bool
	}{
		{"tRPC not found", trpcErr("NOT_FOUND"), true, false, false},
		{"HTTP 404", &warera.HTTPError{StatusCode: 404}, true, false, false},
		{"tRPC rate limit", trpcErr("TOO_MANY_REQUESTS"), false, true, true},
		{"HTTP 429 after retries", fmt.Errorf("request failed after 4 attempts: %w", &warera.HTTPError{StatusCode: 429}), false, true, true},
		{"HTTP 503", &warera.HTTPError{StatusCode: 503}, false, false, true},
		{"tRPC validation", trpcErr("BAD_REQUEST"), false, false, false},
		{"HTTP 401", &warera.HTTPError{StatusCode: 401}, false, false, false},
		{"network", &net.OpError{Op: "dial", Err: errors.New("connection refused")}, false, false, true},
	}
	for _, tt := range tests {
		if got := warera.IsNotFound(tt.err); got != tt.notFound {
			t.Errorf("%s: IsNotFound = %t", tt.name, got)
		}
		if got := warera.IsRateLimited(tt.err); got != tt.rateLimited {
			t.Errorf("%s: IsRateLimited = %t", tt.name, got)
		}
		if got := warera.IsTemporary(tt.err); got != tt.temporary {
			t.Errorf("%s: IsTemporary = %t", tt.name, got)
		}
	}
}

func TestGetManyReturnsFoundRecordsAndSendsEachIDOnce(t *testing.T) {
	var path string
	client := warera.NewClient(
		warera.WithBaseURL("https://example.test/trpc"),
		warera.WithHTTPClient(&http.Client{Transport: errorResponseTransport(func(r *http.Request) (*http.Response, error) {
			path = r.URL.Path
			return &http.Response{StatusCode: http.StatusMultiStatus, Body: io.NopCloser(strings.NewReader(`[` +
				`{"result":{"data":{"_id":"a","username":"alice"}}},` +
				`{"error":{"message":"Users not found.","code":-32004,"data":{"code":"NOT_FOUND","httpStatus":404}}}]`))}, nil
		})}),
		warera.WithRateLimit(0),
		warera.WithMaxRetries(0),
	)
	defer client.Close()

	users, err := client.Users.GetLiteMany(context.Background(), []string{"a", "gone", "a", ""})
	if path != "/trpc/user.getUserLite,user.getUserLite" {
		t.Errorf("request path = %s, want one call per distinct ID", path)
	}
	if len(users) != 1 || users["a"].Username != "alice" {
		t.Errorf("users = %+v, want alice keyed by her ID", users)
	}
	var lookupErr *warera.BatchLookupError
	if !errors.As(err, &lookupErr) || len(lookupErr.Failures) != 1 || lookupErr.Failures[0].ID != "gone" || !warera.IsNotFound(err) {
		t.Fatalf("error = %v, want a BatchLookupError naming only the missing ID", err)
	}
}
