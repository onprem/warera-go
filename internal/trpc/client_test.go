package trpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type dummyInput struct {
	ID string `json:"id"`
}

type dummyOutput struct {
	Name string `json:"name"`
}

func TestClient_CallDirect(t *testing.T) {
	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/trpc/country.getById" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		input := r.URL.Query().Get("input")
		if input != `{"id":"test-1"}` {
			t.Errorf("unexpected input query: %s", input)
		}

		resp := ResponseEnvelope{
			Result: &ResultEnvelope{
				Data: json.RawMessage(`{"name":"Utopia"}`),
			},
		}
		data, _ := json.Marshal(resp)
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewReader(data)),
		}, nil
	})

	httpClient := &http.Client{Transport: mockTransport}

	client := NewClient(ClientConfig{
		BaseURL:            "https://mock.warera.io/trpc",
		HTTPClient:         httpClient,
		BatchWindow:        0, // direct
		RateLimitPerMinute: 0,
	})

	var out dummyOutput
	err := client.CallDirect(context.Background(), "country.getById", dummyInput{ID: "test-1"}, &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Name != "Utopia" {
		t.Errorf("expected Utopia, got %s", out.Name)
	}
}

func TestClient_CallDirect_Error(t *testing.T) {
	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		resp := ResponseEnvelope{
			Error: &TRPCError{
				Message: "Country not found",
				Code:    -32004,
				Data: &TRPCErrorData{
					Code:       "NOT_FOUND",
					HTTPStatus: 404,
				},
			},
		}
		data, _ := json.Marshal(resp)
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "404 Not Found",
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewReader(data)),
		}, nil
	})

	httpClient := &http.Client{Transport: mockTransport}

	client := NewClient(ClientConfig{
		BaseURL:            "https://mock.warera.io/trpc",
		HTTPClient:         httpClient,
		BatchWindow:        0,
		RateLimitPerMinute: 0,
		MaxRetries:         0,
	})

	var out dummyOutput
	err := client.CallDirect(context.Background(), "country.getById", dummyInput{ID: "missing"}, &out)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	trpcErr, ok := err.(*TRPCError)
	if !ok {
		t.Fatalf("expected *TRPCError, got %T: %v", err, err)
	}

	if trpcErr.Data.HTTPStatus != 404 {
		t.Errorf("expected 404, got %d", trpcErr.Data.HTTPStatus)
	}
}

func TestClient_CallBatch_Explicit(t *testing.T) {
	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/trpc/country.getById,company.getById" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("batch") != "1" {
			t.Errorf("expected batch=1")
		}

		resp := []ResponseEnvelope{
			{
				Result: &ResultEnvelope{
					Data: json.RawMessage(`{"name":"CountryA"}`),
				},
			},
			{
				Result: &ResultEnvelope{
					Data: json.RawMessage(`{"name":"CompanyB"}`),
				},
			},
		}
		data, _ := json.Marshal(resp)
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewReader(data)),
		}, nil
	})

	httpClient := &http.Client{Transport: mockTransport}

	client := NewClient(ClientConfig{
		BaseURL:            "https://mock.warera.io/trpc",
		HTTPClient:         httpClient,
		BatchWindow:        0,
		RateLimitPerMinute: 0,
	})

	var out1, out2 dummyOutput
	calls := []BatchCall{
		{Method: "country.getById", Input: dummyInput{ID: "1"}, Output: &out1},
		{Method: "company.getById", Input: dummyInput{ID: "2"}, Output: &out2},
	}

	err := client.CallBatch(context.Background(), calls)
	if err != nil {
		t.Fatalf("unexpected batch error: %v", err)
	}

	if out1.Name != "CountryA" || out2.Name != "CompanyB" {
		t.Errorf("unexpected outputs: %+v, %+v", out1, out2)
	}
}

func TestClient_CallBatch_PreservesPerCallErrors(t *testing.T) {
	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		resp := []ResponseEnvelope{
			{Result: &ResultEnvelope{Data: json.RawMessage(`{"name":"CountryA"}`)}},
			{Error: &TRPCError{Message: "company unavailable", Code: -32004, Data: &TRPCErrorData{Code: "NOT_FOUND", HTTPStatus: 404}}},
		}
		data, _ := json.Marshal(resp)
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewReader(data)),
		}, nil
	})
	client := NewClient(ClientConfig{
		BaseURL:            "https://mock.warera.io/trpc",
		HTTPClient:         &http.Client{Transport: mockTransport},
		RateLimitPerMinute: 0,
	})

	var country, company dummyOutput
	calls := []BatchCall{
		{Method: "country.getById", Input: dummyInput{ID: "1"}, Output: &country},
		{Method: "company.getById", Input: dummyInput{ID: "2"}, Output: &company},
	}
	err := client.CallBatch(context.Background(), calls)
	if err == nil {
		t.Fatal("expected batch to report its failed call")
	}
	if country.Name != "CountryA" {
		t.Errorf("successful call output = %q, want CountryA", country.Name)
	}
	if calls[0].Err != nil {
		t.Errorf("successful call error = %v, want nil", calls[0].Err)
	}
	callErr, ok := calls[1].Err.(*TRPCError)
	if !ok || callErr.Message != "company unavailable" {
		t.Errorf("failed call error = %#v, want company unavailable tRPC error", calls[1].Err)
	}
}

func TestClient_CallBatch_RejectsResponseCountMismatch(t *testing.T) {
	mockTransport := roundTripFunc(func(*http.Request) (*http.Response, error) {
		body := `[{"result":{"data":{"name":"CountryA"}}}]`
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})
	client := NewClient(ClientConfig{
		BaseURL:            "https://mock.warera.io/trpc",
		HTTPClient:         &http.Client{Transport: mockTransport},
		RateLimitPerMinute: 0,
	})

	var first, second dummyOutput
	calls := []BatchCall{
		{Method: "country.getById", Input: dummyInput{ID: "1"}, Output: &first},
		{Method: "company.getById", Input: dummyInput{ID: "2"}, Output: &second},
	}
	err := client.CallBatch(context.Background(), calls)
	if err == nil || !strings.Contains(err.Error(), "mismatched batch response count") {
		t.Fatalf("CallBatch error = %v, want response count mismatch", err)
	}
	for i, call := range calls {
		if call.Err == nil || call.Err.Error() != err.Error() {
			t.Errorf("call %d error = %v, want batch mismatch error %v", i, call.Err, err)
		}
	}
}

func TestClient_CallDirect_RetriesTemporaryHTTPFailure(t *testing.T) {
	requests := 0
	mockTransport := roundTripFunc(func(*http.Request) (*http.Response, error) {
		requests++
		status := http.StatusServiceUnavailable
		statusText := "503 Service Unavailable"
		body := `busy`
		if requests > 1 {
			status = http.StatusOK
			statusText = "200 OK"
			body = `{"result":{"data":{"name":"Recovered"}}}`
		}
		return &http.Response{
			StatusCode: status,
			Status:     statusText,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})
	client := NewClient(ClientConfig{
		BaseURL:            "https://mock.warera.io/trpc",
		HTTPClient:         &http.Client{Transport: mockTransport},
		RateLimitPerMinute: 0,
		MaxRetries:         1,
	})

	var out dummyOutput
	if err := client.CallDirect(context.Background(), "country.getById", dummyInput{ID: "1"}, &out); err != nil {
		t.Fatalf("CallDirect after temporary failure: %v", err)
	}
	if requests != 2 {
		t.Fatalf("request count = %d, want 2", requests)
	}
	if out.Name != "Recovered" {
		t.Fatalf("decoded result = %q, want Recovered", out.Name)
	}
}

func TestClient_CallDirect_StopsAfterMaxRetries(t *testing.T) {
	requests := 0
	mockTransport := roundTripFunc(func(*http.Request) (*http.Response, error) {
		requests++
		return &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Status:     "503 Service Unavailable",
			Header:     http.Header{"Content-Type": []string{"text/plain"}},
			Body:       io.NopCloser(strings.NewReader("still unavailable")),
		}, nil
	})
	client := NewClient(ClientConfig{
		BaseURL:            "https://mock.warera.io/trpc",
		HTTPClient:         &http.Client{Transport: mockTransport},
		RateLimitPerMinute: 0,
		MaxRetries:         1,
	})

	err := client.CallDirect(context.Background(), "country.getById", dummyInput{ID: "1"}, &dummyOutput{})
	if err == nil || !strings.Contains(err.Error(), "request failed after 2 attempts") {
		t.Fatalf("CallDirect error = %v, want exhausted retry error", err)
	}
	if requests != 2 {
		t.Fatalf("request count = %d, want 2", requests)
	}
}

func TestClient_CallDirect_HonorsContextCancellation(t *testing.T) {
	started := make(chan struct{})
	var startOnce sync.Once
	requests := int32(0)
	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		atomic.AddInt32(&requests, 1)
		startOnce.Do(func() { close(started) })
		<-r.Context().Done()
		return nil, r.Context().Err()
	})
	client := NewClient(ClientConfig{
		BaseURL:            "https://mock.warera.io/trpc",
		HTTPClient:         &http.Client{Transport: mockTransport},
		RateLimitPerMinute: 0,
		MaxRetries:         1,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		result <- client.CallDirect(ctx, "country.getById", dummyInput{ID: "1"}, &dummyOutput{})
	}()
	<-started
	cancel()
	if err := <-result; !errors.Is(err, context.Canceled) {
		t.Fatalf("CallDirect error = %v, want context.Canceled", err)
	}
	if got := atomic.LoadInt32(&requests); got != 1 {
		t.Fatalf("request count = %d, want 1 after cancellation", got)
	}
}

func TestClient_TransparentBatching(t *testing.T) {
	var httpRequestsCount int32

	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		atomic.AddInt32(&httpRequestsCount, 1)

		if r.URL.Query().Get("batch") != "1" {
			t.Errorf("expected batch=1")
		}

		methods := strings.Split(strings.TrimPrefix(r.URL.Path, "/trpc/"), ",")
		resp := make([]ResponseEnvelope, len(methods))
		for i, m := range methods {
			switch m {
			case "proc.one":
				resp[i] = ResponseEnvelope{Result: &ResultEnvelope{Data: json.RawMessage(`{"name":"Result1"}`)}}
			case "proc.two":
				resp[i] = ResponseEnvelope{Result: &ResultEnvelope{Data: json.RawMessage(`{"name":"Result2"}`)}}
			case "proc.three":
				resp[i] = ResponseEnvelope{Result: &ResultEnvelope{Data: json.RawMessage(`{"name":"Result3"}`)}}
			}
		}
		data, _ := json.Marshal(resp)
		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewReader(data)),
		}, nil
	})

	httpClient := &http.Client{Transport: mockTransport}

	client := NewClient(ClientConfig{
		BaseURL:            "https://mock.warera.io/trpc",
		HTTPClient:         httpClient,
		BatchWindow:        30 * time.Millisecond,
		MaxBatchSize:       10,
		RateLimitPerMinute: 0,
	})
	defer client.Close()

	var wg sync.WaitGroup
	var out1, out2, out3 dummyOutput
	var callErrs [3]error

	wg.Add(3)
	go func() {
		defer wg.Done()
		callErrs[0] = client.Call(context.Background(), "proc.one", dummyInput{ID: "1"}, &out1)
	}()
	go func() {
		defer wg.Done()
		callErrs[1] = client.Call(context.Background(), "proc.two", dummyInput{ID: "2"}, &out2)
	}()
	go func() {
		defer wg.Done()
		callErrs[2] = client.Call(context.Background(), "proc.three", dummyInput{ID: "3"}, &out3)
	}()

	wg.Wait()
	for i, err := range callErrs {
		if err != nil {
			t.Errorf("call %d failed: %v", i+1, err)
		}
	}

	reqs := atomic.LoadInt32(&httpRequestsCount)
	if reqs != 1 {
		t.Errorf("expected 1 coalesced HTTP request, got %d", reqs)
	}

	if out1.Name != "Result1" || out2.Name != "Result2" || out3.Name != "Result3" {
		t.Errorf("coalesced results mismatch: out1=%+v out2=%+v out3=%+v", out1, out2, out3)
	}
}

func TestRateLimiterWaitHonorsCancellation(t *testing.T) {
	limiter := NewRateLimiter(1)
	if err := limiter.Wait(context.Background()); err != nil {
		t.Fatalf("consume initial token: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := limiter.Wait(ctx); err != context.Canceled {
		t.Fatalf("Wait error = %v, want context.Canceled", err)
	}
}

func TestNilRateLimiterWait(t *testing.T) {
	var limiter *RateLimiter
	if err := limiter.Wait(context.Background()); err != nil {
		t.Fatalf("nil limiter Wait error = %v, want nil", err)
	}
}
