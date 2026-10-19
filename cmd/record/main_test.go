package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestRecordEndpointPreservesLargeNumbers(t *testing.T) {
	const data = `{"large":9007199254740993,"fraction":0.12345678901234568}`
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/trpc/sample.getStats" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("input") != `{"limit":1}` {
			t.Errorf("unexpected input query: %s", r.URL.RawQuery)
		}
		if r.Header.Get("x-api-key") != "fixture-test-key" {
			t.Errorf("API key header was not sent")
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"result":{"data":` + data + `}}`)),
		}, nil
	})}

	outDir := t.TempDir()
	size, err := recordEndpointWithClient("sample.getStats", `{"limit":1}`, "fixture-test-key", outDir, "https://example.test/trpc", client)
	if err != nil {
		t.Fatalf("record endpoint: %v", err)
	}

	path := filepath.Join(outDir, "sample.getStats.json")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read recorded fixture: %v", err)
	}
	if string(got) != `{
  "large": 9007199254740993,
  "fraction": 0.12345678901234568
}` {
		t.Fatalf("recorded JSON changed numeric values: %s", got)
	}
	if size != len(got) {
		t.Fatalf("reported size %d does not match file size %d", size, len(got))
	}
}

func TestRecordEndpointDoesNotOverwriteOnServerError(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"invalid input"}}`)),
		}, nil
	})}

	outDir := t.TempDir()
	path := filepath.Join(outDir, "sample.getStats.json")
	const original = "existing fixture"
	if err := os.WriteFile(path, []byte(original), 0600); err != nil {
		t.Fatalf("write initial fixture: %v", err)
	}

	if _, err := recordEndpointWithClient("sample.getStats", "{}", "", outDir, "https://example.test/trpc", client); err == nil {
		t.Fatal("expected tRPC error")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture after failed request: %v", err)
	}
	if string(got) != original {
		t.Fatalf("failed request overwrote existing fixture: %q", got)
	}
}

func TestRecordEndpointDoesNotOverwriteOnInvalidResponses(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
	}{
		{name: "http failure", statusCode: http.StatusServiceUnavailable, body: `temporary outage`},
		{name: "malformed envelope", statusCode: http.StatusOK, body: `not-json`},
		{name: "missing data", statusCode: http.StatusOK, body: `{"result":{}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: tt.statusCode,
					Status:     http.StatusText(tt.statusCode),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(tt.body)),
				}, nil
			})}
			outDir := t.TempDir()
			path := filepath.Join(outDir, "sample.getStats.json")
			const original = "existing fixture"
			if err := os.WriteFile(path, []byte(original), 0600); err != nil {
				t.Fatalf("write initial fixture: %v", err)
			}

			if _, err := recordEndpointWithClient("sample.getStats", "{}", "", outDir, "https://example.test/trpc", client); err == nil {
				t.Fatal("expected invalid response error")
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read fixture after invalid response: %v", err)
			}
			if string(got) != original {
				t.Fatalf("invalid response overwrote existing fixture: %q", got)
			}
		})
	}
}

func TestBuildInputPayloadIncludesRequiredSeedsAndOmitsUnrelatedFilters(t *testing.T) {
	properties := map[string]OpenAPIProperty{
		"userId":    {Type: "string"},
		"mode":      {Type: "string", Enum: []string{"recent", "top"}},
		"limit":     {Type: "integer"},
		"cursor":    {Type: "string"},
		"countryId": {Type: "string"},
	}
	got := buildInputPayload(properties, []string{"userId", "mode"}, map[string]any{
		"userId":    "user-1",
		"mode":      "top",
		"limit":     10,
		"cursor":    "cursor-1",
		"countryId": "country-1",
	})

	var payload map[string]any
	if err := json.Unmarshal([]byte(got), &payload); err != nil {
		t.Fatalf("decode generated input: %v", err)
	}
	want := map[string]any{"userId": "user-1", "mode": "recent", "limit": float64(10)}
	if len(payload) != len(want) {
		t.Fatalf("payload = %s, want only required fields and limit", got)
	}
	for key, value := range want {
		if payload[key] != value {
			t.Errorf("payload[%q] = %v, want %v", key, payload[key], value)
		}
	}
}

func TestHasUnresolvedRequiredProps(t *testing.T) {
	properties := map[string]OpenAPIProperty{
		"mode": {Type: "string", Enum: []string{"recent", "top"}},
		"user": {Type: "string"},
	}
	tests := []struct {
		name     string
		required []string
		seeds    map[string]any
		want     bool
	}{
		{name: "all required values seeded", required: []string{"user"}, seeds: map[string]any{"user": "user-1"}, want: false},
		{name: "enum can resolve required value", required: []string{"mode"}, want: false},
		{name: "missing required value", required: []string{"user"}, want: true},
		{name: "cursor is optional pagination state", required: []string{"cursor"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasUnresolvedRequiredProps(properties, tt.required, tt.seeds); got != tt.want {
				t.Fatalf("hasUnresolvedRequiredProps() = %v, want %v", got, tt.want)
			}
		})
	}
}
