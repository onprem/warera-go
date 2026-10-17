package warera

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/onprem/warera-go/internal/trpc"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestClient_Countries_GetAll(t *testing.T) {
	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if !strings.HasSuffix(r.URL.Path, "/country.getAllCountries") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := trpc.ResponseEnvelope{
			Result: &trpc.ResultEnvelope{
				Data: json.RawMessage(`[
					{"_id":"c1","name":"Switzerland","code":"ch","money":100.5,"taxes":{"income":1,"market":1,"selfWork":1}},
					{"_id":"c2","name":"Singapore","code":"sg","money":250.0,"taxes":{"income":5,"market":5,"selfWork":10}}
				]`),
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

	client := NewClient(
		WithHTTPClient(&http.Client{Transport: mockTransport}),
		WithBatchWindow(0), // direct
		WithRateLimit(0),
	)
	defer client.Close()

	countries, err := client.Countries.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(countries) != 2 {
		t.Fatalf("expected 2 countries, got %d", len(countries))
	}

	if countries[0].Name != "Switzerland" || countries[0].Code != "ch" {
		t.Errorf("unexpected country[0]: %+v", countries[0])
	}
	if countries[1].Name != "Singapore" || countries[1].Money != 250.0 {
		t.Errorf("unexpected country[1]: %+v", countries[1])
	}
}

func TestClient_Countries_GetByID(t *testing.T) {
	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if !strings.HasSuffix(r.URL.Path, "/country.getCountryById") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if !strings.Contains(r.URL.RawQuery, "%22countryId%22%3A%22c1%22") {
			t.Errorf("expected countryId in query, got: %s", r.URL.RawQuery)
		}

		resp := trpc.ResponseEnvelope{
			Result: &trpc.ResultEnvelope{
				Data: json.RawMessage(`{"_id":"c1","name":"Switzerland","code":"ch","money":100.5,"currentPopulation":34}`),
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

	client := NewClient(
		WithHTTPClient(&http.Client{Transport: mockTransport}),
		WithBatchWindow(0),
		WithRateLimit(0),
	)
	defer client.Close()

	country, err := client.Countries.GetByID(context.Background(), "c1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if country.ID != "c1" || country.CurrentPopulation != 34 {
		t.Errorf("unexpected country data: %+v", country)
	}
}

func TestClient_TransparentBatching_MultiService(t *testing.T) {
	var httpRequestsCount int32

	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		atomic.AddInt32(&httpRequestsCount, 1)

		if r.URL.Query().Get("batch") != "1" {
			t.Errorf("expected batch=1, got %s", r.URL.String())
		}

		methods := strings.Split(strings.TrimPrefix(r.URL.Path, "/trpc/"), ",")
		resp := make([]trpc.ResponseEnvelope, len(methods))
		for i, m := range methods {
			switch m {
			case "country.getCountryById":
				resp[i] = trpc.ResponseEnvelope{Result: &trpc.ResultEnvelope{Data: json.RawMessage(`{"_id":"c1","name":"Switzerland"}`)}}
			case "company.getById":
				resp[i] = trpc.ResponseEnvelope{Result: &trpc.ResultEnvelope{Data: json.RawMessage(`{"_id":"comp1","name":"Iron Mine Inc"}`)}}
			case "gameConfig.getDates":
				resp[i] = trpc.ResponseEnvelope{Result: &trpc.ResultEnvelope{Data: json.RawMessage(`{"nextDayAt":"2026-10-18T00:00:00Z"}`)}}
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

	client := NewClient(
		WithHTTPClient(&http.Client{Transport: mockTransport}),
		WithBatchWindow(25*time.Millisecond),
		WithRateLimit(0),
	)
	defer client.Close()

	var wg sync.WaitGroup
	var country *Country
	var company *Company
	var dates *GameDates

	wg.Add(3)
	go func() {
		defer wg.Done()
		var err error
		country, err = client.Countries.GetByID(context.Background(), "c1")
		if err != nil {
			t.Errorf("country error: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		company, err = client.Companies.GetByID(context.Background(), "comp1")
		if err != nil {
			t.Errorf("company error: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		var err error
		dates, err = client.GameConfig.GetDates(context.Background())
		if err != nil {
			t.Errorf("dates error: %v", err)
		}
	}()

	wg.Wait()

	reqs := atomic.LoadInt32(&httpRequestsCount)
	if reqs != 1 {
		t.Errorf("expected 1 batched HTTP request, got %d", reqs)
	}

	if country == nil || country.Name != "Switzerland" {
		t.Errorf("expected Switzerland, got %+v", country)
	}
	if company == nil || company.Name != "Iron Mine Inc" {
		t.Errorf("expected Iron Mine Inc, got %+v", company)
	}
	if dates == nil || dates.NextDayAt != "2026-10-18T00:00:00Z" {
		t.Errorf("expected date 2026-10-18T00:00:00Z, got %+v", dates)
	}
}

func TestClient_Articles_Iterator(t *testing.T) {
	pageCounter := 0

	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		pageCounter++
		var resp trpc.ResponseEnvelope
		if pageCounter == 1 {
			resp = trpc.ResponseEnvelope{
				Result: &trpc.ResultEnvelope{
					Data: json.RawMessage(`{
						"items": [
							{"_id":"a1","title":"Article 1","author":"Dog"},
							{"_id":"a2","title":"Article 2","author":"Cat"}
						],
						"nextCursor": "cur-page-2"
					}`),
				},
			}
		} else {
			resp = trpc.ResponseEnvelope{
				Result: &trpc.ResultEnvelope{
					Data: json.RawMessage(`{
						"items": [
							{"_id":"a3","title":"Article 3","author":"Bird"}
						],
						"nextCursor": ""
					}`),
				},
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

	client := NewClient(
		WithHTTPClient(&http.Client{Transport: mockTransport}),
		WithBatchWindow(0),
		WithRateLimit(0),
	)
	defer client.Close()

	var titles []string
	for article, err := range client.Articles.All(context.Background(), ArticleFilter{Type: "last"}) {
		if err != nil {
			t.Fatalf("unexpected iterator err: %v", err)
		}
		titles = append(titles, article.Title)
	}

	if len(titles) != 3 {
		t.Fatalf("expected 3 iterated articles, got %d (%v)", len(titles), titles)
	}

	if titles[0] != "Article 1" || titles[1] != "Article 2" || titles[2] != "Article 3" {
		t.Errorf("unexpected titles: %v", titles)
	}
}

func TestClient_SearchAnything(t *testing.T) {
	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if !strings.HasSuffix(r.URL.Path, "/search.searchAnything") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := trpc.ResponseEnvelope{
			Result: &trpc.ResultEnvelope{
				Data: json.RawMessage(`{
					"userIds": ["u123"],
					"muIds": [],
					"countryIds": [],
					"regionIds": [],
					"hasData": true
				}`),
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

	client := NewClient(
		WithHTTPClient(&http.Client{Transport: mockTransport}),
		WithBatchWindow(0),
		WithRateLimit(0),
	)
	defer client.Close()

	res, err := client.Search.SearchAnything(context.Background(), "Kautilya")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.UserIDs) != 1 || res.UserIDs[0] != "u123" {
		t.Errorf("expected [u123], got %v", res.UserIDs)
	}
	if !res.HasData {
		t.Errorf("expected hasData true")
	}
}

func TestClient_SearchUsers(t *testing.T) {
	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if strings.HasSuffix(r.URL.Path, "/search.searchAnything") {
			resp := trpc.ResponseEnvelope{
				Result: &trpc.ResultEnvelope{
					Data: json.RawMessage(`{"userIds":["u123"],"hasData":true}`),
				},
			}
			data, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(data)),
			}, nil
		}

		if strings.HasSuffix(r.URL.Path, "/user.getUserLite") {
			resp := trpc.ResponseEnvelope{
				Result: &trpc.ResultEnvelope{
					Data: json.RawMessage(`{
						"_id": "u123",
						"username": "Kautilya",
						"country": "c1",
						"avatarUrl": "https://media.warera.io/avatar.jpg",
						"leveling": {"level": 24, "xp": 20100}
					}`),
				},
			}
			data, _ := json.Marshal(resp)
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(data)),
			}, nil
		}

		return nil, nil
	})

	client := NewClient(
		WithHTTPClient(&http.Client{Transport: mockTransport}),
		WithBatchWindow(0),
		WithRateLimit(0),
	)
	defer client.Close()

	users, err := client.Search.SearchUsers(context.Background(), "Kautilya")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}

	if users[0].Username != "Kautilya" || users[0].Level() != 24 {
		t.Errorf("unexpected user data: %+v", users[0])
	}
}
