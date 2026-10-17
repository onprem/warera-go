package main

import (
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
