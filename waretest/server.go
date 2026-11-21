// Package waretest provides a fake WarEra API server for testing tools built
// on warera-go without network access or an API key.
//
// A new server answers every procedure the SDK wraps with a response recorded
// from the live API, so a tool can run against realistic data in one line:
//
//	srv := waretest.NewServer(t)
//	client := warera.NewClient(warera.WithBaseURL(srv.URL), warera.WithRateLimit(0))
//
// WithRateLimit(0) turns off the client's own rate limit, which would
// otherwise slow a test that makes many calls.
//
// Recorded responses answer any input. Use [Server.Handle] to answer a
// procedure yourself, for example to return chosen data, a null result, or an
// API error:
//
//	srv.Handle("user.getUserLite", waretest.Fail(waretest.NotFound))
//
// The server speaks tRPC's HTTP protocol, including the batched requests the
// client sends for concurrent calls and bulk lookups.
package waretest

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

//go:embed fixtures/*.json
var fixtures embed.FS

// Fixture returns the response recorded for procedure, such as
// "gameConfig.getGameConfig", decoded into T. It lets tests and examples of
// code that computes from API records, like the analytics packages, start
// from real data without a server:
//
//	cfg, err := waretest.Fixture[*warera.GameConfig]("gameConfig.getGameConfig")
func Fixture[T any](procedure string) (T, error) {
	var v T
	data, err := fixtures.ReadFile("fixtures/" + procedure + ".json")
	if err != nil {
		return v, fmt.Errorf("waretest: no recorded response for %s", procedure)
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return v, fmt.Errorf("waretest: decode %s: %w", procedure, err)
	}
	return v, nil
}

// A Handler answers one call to a procedure. input is the call's JSON input,
// or nil for a call without input. A returned *[Error] becomes a tRPC error;
// any other error becomes INTERNAL_SERVER_ERROR. A nil result is sent as null.
type Handler func(input json.RawMessage) (any, error)

// Respond returns a Handler that answers every call with result.
func Respond(result any) Handler {
	return func(json.RawMessage) (any, error) { return result, nil }
}

// Fail returns a Handler that answers every call with err.
func Fail(err *Error) Handler {
	return func(json.RawMessage) (any, error) { return nil, err }
}

// Error is a tRPC error a Handler can return. Code is a tRPC error code such
// as "NOT_FOUND"; the server sends the matching HTTP status.
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string { return e.Code + ": " + e.Message }

// Errors the API commonly returns.
var (
	NotFound     = &Error{Code: "NOT_FOUND", Message: "Not found."}
	BadRequest   = &Error{Code: "BAD_REQUEST", Message: "Invalid input."}
	Unauthorized = &Error{Code: "UNAUTHORIZED", Message: "API key required."}
	Forbidden    = &Error{Code: "FORBIDDEN", Message: "Forbidden."}
	RateLimited  = &Error{Code: "TOO_MANY_REQUESTS", Message: "Too many requests."}
	Internal     = &Error{Code: "INTERNAL_SERVER_ERROR", Message: "Internal server error."}
)

// errorCodes maps tRPC error codes to their JSON-RPC code and HTTP status.
var errorCodes = map[string][2]int{
	"BAD_REQUEST":           {-32600, http.StatusBadRequest},
	"UNAUTHORIZED":          {-32001, http.StatusUnauthorized},
	"FORBIDDEN":             {-32003, http.StatusForbidden},
	"NOT_FOUND":             {-32004, http.StatusNotFound},
	"TIMEOUT":               {-32008, http.StatusRequestTimeout},
	"TOO_MANY_REQUESTS":     {-32029, http.StatusTooManyRequests},
	"INTERNAL_SERVER_ERROR": {-32603, http.StatusInternalServerError},
}

// Server is a fake WarEra API. Pass its URL to warera.WithBaseURL. It is safe
// for concurrent use.
type Server struct {
	*httptest.Server

	mu       sync.Mutex
	handlers map[string]Handler
	calls    map[string]int
	inputs   map[string][]json.RawMessage
}

// NewServer starts a server that answers every wrapped procedure with its
// recorded response, and closes it when the test ends.
func NewServer(t testing.TB) *Server {
	t.Helper()
	s := &Server{handlers: map[string]Handler{}, calls: map[string]int{}, inputs: map[string][]json.RawMessage{}}
	entries, err := fixtures.ReadDir("fixtures")
	if err != nil {
		t.Fatalf("waretest: read fixtures: %v", err)
	}
	for _, entry := range entries {
		data, err := fixtures.ReadFile("fixtures/" + entry.Name())
		if err != nil {
			t.Fatalf("waretest: read fixture %s: %v", entry.Name(), err)
		}
		s.handlers[strings.TrimSuffix(entry.Name(), ".json")] = Respond(json.RawMessage(data))
	}
	s.Server = httptest.NewServer(http.HandlerFunc(s.serveHTTP))
	t.Cleanup(s.Close)
	return s
}

// Handle answers procedure, such as "country.getCountryById", with h instead
// of its recorded response.
func (s *Server) Handle(procedure string, h Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[procedure] = h
}

// Calls reports how many calls to procedure the server has answered. A
// batched request counts once per call in it.
func (s *Server) Calls(procedure string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls[procedure]
}

// Inputs returns the JSON inputs of the calls to procedure the server has
// answered, oldest first. A call without input has a nil entry.
func (s *Server) Inputs(procedure string) []json.RawMessage {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]json.RawMessage(nil), s.inputs[procedure]...)
}

type envelope struct {
	Result *struct {
		Data any `json:"data"`
	} `json:"result,omitempty"`
	Error *errorBody `json:"error,omitempty"`
}

type errorBody struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    struct {
		Code       string `json:"code"`
		HTTPStatus int    `json:"httpStatus"`
		Path       string `json:"path"`
	} `json:"data"`
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	procedures := strings.Split(r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:], ",")
	rawInput := r.URL.Query().Get("input")
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rawInput = string(body)
	}

	if r.URL.Query().Get("batch") != "1" {
		var input json.RawMessage
		if rawInput != "" {
			input = json.RawMessage(rawInput)
		}
		env, status := s.call(procedures[0], input)
		writeJSON(w, status, env)
		return
	}

	inputs := map[string]json.RawMessage{}
	if rawInput != "" {
		if err := json.Unmarshal([]byte(rawInput), &inputs); err != nil {
			http.Error(w, "waretest: invalid batch input: "+err.Error(), http.StatusBadRequest)
			return
		}
	}
	envs := make([]envelope, len(procedures))
	statuses := map[int]bool{}
	for i, procedure := range procedures {
		var status int
		envs[i], status = s.call(procedure, inputs[fmt.Sprint(i)])
		statuses[status] = true
	}
	// Like tRPC, a batch whose calls all share a status uses it; a mixed batch
	// answers 207 Multi-Status.
	status := http.StatusMultiStatus
	if len(statuses) == 1 {
		for only := range statuses {
			status = only
		}
	}
	writeJSON(w, status, envs)
}

func (s *Server) call(procedure string, input json.RawMessage) (envelope, int) {
	s.mu.Lock()
	h, ok := s.handlers[procedure]
	s.calls[procedure]++
	s.inputs[procedure] = append(s.inputs[procedure], input)
	s.mu.Unlock()
	if !ok {
		return errorEnvelope(procedure, &Error{Code: "NOT_FOUND", Message: fmt.Sprintf("No procedure found on path %q", procedure)})
	}
	result, err := h(input)
	if err != nil {
		apiErr, ok := err.(*Error)
		if !ok {
			apiErr = &Error{Code: "INTERNAL_SERVER_ERROR", Message: err.Error()}
		}
		return errorEnvelope(procedure, apiErr)
	}
	env := envelope{Result: &struct {
		Data any `json:"data"`
	}{Data: result}}
	return env, http.StatusOK
}

func errorEnvelope(procedure string, err *Error) (envelope, int) {
	codes, ok := errorCodes[err.Code]
	if !ok {
		codes = errorCodes["INTERNAL_SERVER_ERROR"]
	}
	body := &errorBody{Message: err.Message, Code: codes[0]}
	body.Data.Code = err.Code
	body.Data.HTTPStatus = codes[1]
	body.Data.Path = procedure
	return envelope{Error: body}, codes[1]
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
