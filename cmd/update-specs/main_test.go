package main

import (
	"errors"
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

func TestGetWithClientSendsUserAgentAndReadsBody(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.Header.Get("User-Agent") != userAgent {
			t.Errorf("User-Agent = %q, want %q", r.Header.Get("User-Agent"), userAgent)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(strings.NewReader("response bytes")),
		}, nil
	})}

	got, err := getWithClient("https://example.test/resource", client)
	if err != nil {
		t.Fatalf("getWithClient: %v", err)
	}
	if string(got) != "response bytes" {
		t.Fatalf("body = %q, want response bytes", got)
	}
}

func TestGetWithClientRejectsHTTPFailure(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadGateway,
			Status:     "502 Bad Gateway",
			Body:       io.NopCloser(strings.NewReader("upstream unavailable")),
		}, nil
	})}
	_, err := getWithClient("https://example.test/resource", client)
	if err == nil || !strings.Contains(err.Error(), "HTTP 502 Bad Gateway") {
		t.Fatalf("getWithClient error = %v, want HTTP status error", err)
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, errors.New("read failed") }
func (errorReader) Close() error             { return nil }

func TestGetWithClientReportsBodyReadFailure(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       errorReader{},
		}, nil
	})}
	_, err := getWithClient("https://example.test/resource", client)
	if err == nil || !strings.Contains(err.Error(), "read failed") {
		t.Fatalf("getWithClient error = %v, want body read error", err)
	}
}

func TestStageAndReplaceWritesChangesAndSkipsIdenticalFiles(t *testing.T) {
	specDir := t.TempDir()
	firstPath := filepath.Join(specDir, "openapi.json")
	secondPath := filepath.Join(specDir, "second.json")
	outputs := []fileOutput{
		{path: firstPath, data: []byte("openapi")},
		{path: secondPath, data: []byte("responses")},
	}

	changed, err := stageAndReplace(specDir, outputs)
	if err != nil {
		t.Fatalf("stageAndReplace: %v", err)
	}
	if changed != 2 {
		t.Fatalf("changed count = %d, want 2", changed)
	}
	for _, output := range outputs {
		got, err := os.ReadFile(output.path)
		if err != nil {
			t.Fatalf("read %s: %v", output.path, err)
		}
		if string(got) != string(output.data) {
			t.Errorf("%s = %q, want %q", output.path, got, output.data)
		}
	}

	changed, err = stageAndReplace(specDir, outputs)
	if err != nil {
		t.Fatalf("stageAndReplace identical files: %v", err)
	}
	if changed != 0 {
		t.Fatalf("identical file changed count = %d, want 0", changed)
	}
}

func TestStageAndReplaceDoesNotReplaceBeforeAllOutputsAreStaged(t *testing.T) {
	specDir := t.TempDir()
	firstPath := filepath.Join(specDir, "first.json")
	if err := os.WriteFile(firstPath, []byte("original"), 0o600); err != nil {
		t.Fatalf("write first destination: %v", err)
	}
	blockedPath := filepath.Join(specDir, "blocked")
	if err := os.Mkdir(blockedPath, 0o700); err != nil {
		t.Fatalf("create blocking directory: %v", err)
	}

	_, err := stageAndReplace(specDir, []fileOutput{
		{path: firstPath, data: []byte("replacement")},
		{path: blockedPath, data: []byte("cannot replace a directory")},
	})
	if err == nil {
		t.Fatal("expected staging error for directory destination")
	}
	got, readErr := os.ReadFile(firstPath)
	if readErr != nil {
		t.Fatalf("read first destination: %v", readErr)
	}
	if string(got) != "original" {
		t.Fatalf("first destination changed before staging completed: %q", got)
	}
}
