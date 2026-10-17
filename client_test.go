package warera

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func loadFixture(t *testing.T, filename string) json.RawMessage {
	t.Helper()
	path := filepath.Join("spec", "fixtures", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		path = filepath.Join("..", "spec", "fixtures", filename)
		data, err = os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to load fixture %s: %v", filename, err)
		}
	}
	return json.RawMessage(data)
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

func TestClient_Countries_GetByID_Fixture(t *testing.T) {
	fixtureData := loadFixture(t, "country.getCountryById.json")

	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if !strings.HasSuffix(r.URL.Path, "/country.getCountryById") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if !strings.Contains(r.URL.RawQuery, "6813b6d446e731854c7ac7a6") {
			t.Errorf("expected countryId in query, got: %s", r.URL.RawQuery)
		}

		resp := trpc.ResponseEnvelope{
			Result: &trpc.ResultEnvelope{
				Data: fixtureData,
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

	country, err := client.Countries.GetByID(context.Background(), "6813b6d446e731854c7ac7a6")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if country.ID != "6813b6d446e731854c7ac7a6" {
		t.Errorf("expected ID 6813b6d446e731854c7ac7a6, got %s", country.ID)
	}
	if country.Name != "Switzerland" {
		t.Errorf("expected Switzerland, got %s", country.Name)
	}
	if country.Code != "ch" {
		t.Errorf("expected ch, got %s", country.Code)
	}
	if country.CurrentPopulation != 34 {
		t.Errorf("expected population 34, got %d", country.CurrentPopulation)
	}
	if country.Money < 3.0 {
		t.Errorf("expected money > 3.0, got %f", country.Money)
	}
	if len(country.DefensivePacts) != 6 {
		t.Errorf("expected 6 defensive pacts, got %d", len(country.DefensivePacts))
	}
	if popRank, ok := country.Rankings["countryActivePopulation"]; !ok || popRank.Rank != 91 {
		t.Errorf("expected countryActivePopulation rank 91, got %+v", popRank)
	}
}

func TestClient_TransparentBatching_MultiService(t *testing.T) {
	var httpRequestsCount int32
	countryFixture := loadFixture(t, "country.getCountryById.json")
	datesFixture := loadFixture(t, "gameConfig.getDates.json")
	pricesFixture := loadFixture(t, "itemTrading.getPrices.json")

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
				resp[i] = trpc.ResponseEnvelope{Result: &trpc.ResultEnvelope{Data: countryFixture}}
			case "gameConfig.getDates":
				resp[i] = trpc.ResponseEnvelope{Result: &trpc.ResultEnvelope{Data: datesFixture}}
			case "itemTrading.getPrices":
				resp[i] = trpc.ResponseEnvelope{Result: &trpc.ResultEnvelope{Data: pricesFixture}}
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
	var dates *GameDates
	var prices ItemPrices

	wg.Add(3)
	go func() {
		defer wg.Done()
		var err error
		country, err = client.Countries.GetByID(context.Background(), "6813b6d446e731854c7ac7a6")
		if err != nil {
			t.Errorf("country error: %v", err)
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
	go func() {
		defer wg.Done()
		var err error
		prices, err = client.Trading.GetPrices(context.Background())
		if err != nil {
			t.Errorf("prices error: %v", err)
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
	if dates == nil || dates.NextDayAt == "" {
		t.Errorf("expected dates populated, got %+v", dates)
	}
	if _, err := dates.NextDayTime(); err != nil {
		t.Errorf("failed to parse NextDayTime: %v", err)
	}
	if prices == nil || prices["ammo"] <= 0 {
		t.Errorf("expected ammo price > 0, got %v", prices)
	}
}

func TestClient_Articles_Iterator_Fixture(t *testing.T) {
	articleFixture := loadFixture(t, "article.getArticlesPaginated.json")
	pageCounter := 0

	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		pageCounter++
		var resp trpc.ResponseEnvelope
		if pageCounter == 1 {
			resp = trpc.ResponseEnvelope{
				Result: &trpc.ResultEnvelope{
					Data: articleFixture,
				},
			}
		} else {
			resp = trpc.ResponseEnvelope{
				Result: &trpc.ResultEnvelope{
					Data: json.RawMessage(`{"items": [], "nextCursor": ""}`),
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

	var articles []Article
	for article, err := range client.Articles.All(context.Background(), ArticleFilter{Type: "last"}) {
		if err != nil {
			t.Fatalf("unexpected iterator err: %v", err)
		}
		articles = append(articles, article)
	}

	if len(articles) != 5 {
		t.Fatalf("expected 5 iterated articles from fixture, got %d", len(articles))
	}

	if articles[0].Title != "Two Revolts One Nation" {
		t.Errorf("unexpected article 0 title: %s", articles[0].Title)
	}
	if articles[0].Stats.Likes != 5 {
		t.Errorf("expected 5 likes on article 0, got %d", articles[0].Stats.Likes)
	}
	if articles[1].Language != "tr" {
		t.Errorf("expected language tr on article 1, got %s", articles[1].Language)
	}
}

func TestClient_SearchAnything_Fixture(t *testing.T) {
	fixtureData := loadFixture(t, "search.searchAnything.json")

	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if !strings.HasSuffix(r.URL.Path, "/search.searchAnything") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := trpc.ResponseEnvelope{
			Result: &trpc.ResultEnvelope{
				Data: fixtureData,
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

	if len(res.UserIDs) != 1 || res.UserIDs[0] != "6aa7ad06d91c79361d0b0673" {
		t.Errorf("expected [6aa7ad06d91c79361d0b0673], got %v", res.UserIDs)
	}
	if !res.HasData {
		t.Errorf("expected hasData true")
	}
}

func TestClient_SearchUsers_Fixture(t *testing.T) {
	searchFixture := loadFixture(t, "search.searchAnything.json")
	userFixture := loadFixture(t, "user.getUserLite.json")

	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if strings.HasSuffix(r.URL.Path, "/search.searchAnything") {
			resp := trpc.ResponseEnvelope{
				Result: &trpc.ResultEnvelope{
					Data: searchFixture,
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
					Data: userFixture,
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

	u := users[0]
	if u.ID != "6aa7ad06d91c79361d0b0673" {
		t.Errorf("expected ID 6aa7ad06d91c79361d0b0673, got %s", u.ID)
	}
	if u.Username != "Kautilya" {
		t.Errorf("expected Username Kautilya, got %s", u.Username)
	}
	if u.Level() != 24 {
		t.Errorf("expected Level 24, got %d", u.Level())
	}
	if u.MilitaryRank != 31 {
		t.Errorf("expected military rank 31, got %d", u.MilitaryRank)
	}
	if u.MUID != "6aa9311d8af2bd246aff6189" {
		t.Errorf("expected MU 6aa9311d8af2bd246aff6189, got %s", u.MUID)
	}
}

func TestClient_Battles_GetBattles_Fixture(t *testing.T) {
	battleFixture := loadFixture(t, "battle.getBattles.json")

	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if !strings.HasSuffix(r.URL.Path, "/battle.getBattles") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := trpc.ResponseEnvelope{
			Result: &trpc.ResultEnvelope{
				Data: battleFixture,
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

	page, err := client.Battles.GetBattles(context.Background(), BattleFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(page.Items) != 5 {
		t.Fatalf("expected 5 battles, got %d", len(page.Items))
	}
	if !page.HasNext() {
		t.Errorf("expected HasNext to be true")
	}

	b := page.Items[0]
	if b.ID != "6ad3b89d7b158f464a343664" {
		t.Errorf("expected battle ID 6ad3b89d7b158f464a343664, got %s", b.ID)
	}
	if b.AttackerCountry() != "6813b6d446e731854c7ac7ba" {
		t.Errorf("expected attacker country 6813b6d446e731854c7ac7ba, got %s", b.AttackerCountry())
	}
	if b.DefenderCountry() != "6813b6d446e731854c7ac7bc" {
		t.Errorf("expected defender country 6813b6d446e731854c7ac7bc, got %s", b.DefenderCountry())
	}
	if b.DefenderRegion() != "6813b7059403bc4170a5d79b" {
		t.Errorf("expected defender region 6813b7059403bc4170a5d79b, got %s", b.DefenderRegion())
	}
	if b.CurrentRound == nil || b.CurrentRound.Number != 1 {
		t.Errorf("expected CurrentRound number 1, got %+v", b.CurrentRound)
	}
	if b.Stats.HitCount != 0 {
		t.Errorf("expected hit count 0, got %d", b.Stats.HitCount)
	}
	if b.RoundsToWin != 2 {
		t.Errorf("expected roundsToWin 2, got %d", b.RoundsToWin)
	}
}

func TestClient_Market_GetPrices_Fixture(t *testing.T) {
	pricesFixture := loadFixture(t, "itemTrading.getPrices.json")

	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if !strings.HasSuffix(r.URL.Path, "/itemTrading.getPrices") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := trpc.ResponseEnvelope{
			Result: &trpc.ResultEnvelope{
				Data: pricesFixture,
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

	prices, err := client.Trading.GetPrices(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if prices["ammo"] <= 0 {
		t.Errorf("expected ammo price > 0, got %f", prices["ammo"])
	}
	if prices["bread"] <= 0 {
		t.Errorf("expected bread price > 0, got %f", prices["bread"])
	}

	list := prices.List()
	if len(list) != len(prices) {
		t.Errorf("expected list length %d, got %d", len(prices), len(list))
	}
	// Verify sorted order
	for i := 1; i < len(list); i++ {
		if list[i-1].ItemCode > list[i].ItemCode {
			t.Errorf("expected list sorted, found %s after %s", list[i].ItemCode, list[i-1].ItemCode)
		}
	}
}

func TestClient_GameConfig_GetDates_Fixture(t *testing.T) {
	datesFixture := loadFixture(t, "gameConfig.getDates.json")

	mockTransport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if !strings.HasSuffix(r.URL.Path, "/gameConfig.getDates") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		resp := trpc.ResponseEnvelope{
			Result: &trpc.ResultEnvelope{
				Data: datesFixture,
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

	dates, err := client.GameConfig.GetDates(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	nextDay, err := dates.NextDayTime()
	if err != nil {
		t.Fatalf("unexpected error parsing NextDayTime: %v", err)
	}
	if nextDay.IsZero() {
		t.Errorf("expected valid non-zero time")
	}

	nextRegen, err := dates.NextRegenTime()
	if err != nil {
		t.Fatalf("unexpected error parsing NextRegenTime: %v", err)
	}
	if nextRegen.IsZero() {
		t.Errorf("expected valid non-zero time")
	}
}
