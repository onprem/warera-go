package warera_test

import (
	"context"
	"errors"
	"io"
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
