package warera

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onprem/warera-go/internal/trpc"
)

func TestResponseFixturesDecode(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		decode   func(json.RawMessage) error
	}{
		{"article detail", "article.getArticleById.json", decodeFixtureAs[Article]},
		{"article list", "article.getArticlesPaginated.json", decodeFixtureAs[CursorPage[Article]]},
		{"article lite", "article.getArticleLiteById.json", decodeFixtureAs[*Article]},
		{"battle list", "battle.getBattles.json", decodeFixtureAs[CursorPage[Battle]]},
		{"battle detail", "battle.getById.json", decodeFixtureAs[Battle]},
		{"live battle", "battle.getLiveBattleData.json", decodeFixtureAs[LiveBattleData]},
		{"battle loot summary", "battleLootSummary.getByBattleAndUser.json", decodeFixtureAs[BattleLootSummary]},
		{"battle orders", "battleOrder.getByBattle.json", decodeFixtureAs[[]BattleOrder]},
		{"battle ranking", "battleRanking.getRanking.json", decodeFixtureAs[BattleRanking]},
		{"company detail", "company.getById.json", decodeFixtureAs[Company]},
		{"company IDs", "company.getCompanies.json", decodeFixtureAs[CursorPage[string]]},
		{"all countries", "country.getAllCountries.json", decodeFixtureAs[[]Country]},
		{"country detail", "country.getCountryById.json", decodeFixtureAs[Country]},
		{"events", "event.getEventsPaginated.json", decodeFixtureAs[CursorPage[Event]]},
		{"game dates", "gameConfig.getDates.json", decodeFixtureAs[GameDates]},
		{"game config", "gameConfig.getGameConfig.json", decodeFixtureAs[GameConfig]},
		{"government", "government.getByCountryId.json", decodeFixtureAs[Government]},
		{"current equipment", "inventory.fetchCurrentEquipment.json", decodeFixtureAs[CurrentEquipment]},
		{"item prices", "itemTrading.getPrices.json", decodeFixtureAs[ItemPrices]},
		{"military unit detail", "mu.getById.json", decodeFixtureAs[MilitaryUnit]},
		{"military unit list", "mu.getManyPaginated.json", decodeFixtureAs[CursorPage[MilitaryUnit]]},
		{"mercenary auctions", "mercenaryContractAuction.getPaginatedAuctions.json", decodeFixtureAs[CursorPage[MercenaryContractAuction]]},
		{"ranking", "ranking.getRanking.json", decodeFixtureAs[Ranking]},
		{"region detail", "region.getById.json", decodeFixtureAs[Region]},
		{"region map", "region.getRegionsObject.json", decodeFixtureAs[map[string]Region]},
		{"round detail", "round.getById.json", decodeFixtureAs[Round]},
		{"round last hits", "round.getLastHits.json", decodeFixtureAs[RoundLastHits]},
		{"search results", "search.searchAnything.json", decodeFixtureAs[SearchResult]},
		{"top orders", "tradingOrder.getTopOrders.json", decodeFixtureAs[TopOrders]},
		{"upgrade", "upgrade.getUpgradeByTypeAndEntity.json", decodeFixtureAs[Upgrade]},
		{"user detail", "user.getUserById.json", decodeFixtureAs[User]},
		{"user lite", "user.getUserLite.json", decodeFixtureAs[UserLite]},
		{"users by country", "user.getUsersByCountry.json", decodeFixtureAs[CursorPage[User]]},
		{"work offer", "workOffer.getById.json", decodeFixtureAs[WorkOffer]},
		{"work offer by company", "workOffer.getWorkOfferByCompanyId.json", decodeFixtureAs[WorkOffer]},
		{"work offers", "workOffer.getWorkOffersPaginated.json", decodeFixtureAs[CursorPage[WorkOffer]]},
	}

	covered := make(map[string]struct{}, len(tests))
	for _, tt := range tests {
		covered[tt.filename] = struct{}{}
	}
	fixtures, err := filepath.Glob(filepath.Join("spec", "fixtures", "*.json"))
	if err != nil {
		t.Fatalf("list response fixtures: %v", err)
	}
	for _, fixture := range fixtures {
		if _, ok := covered[filepath.Base(fixture)]; !ok {
			t.Errorf("fixture %s has no typed decode check", fixture)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := loadFixture(t, tt.filename)
			if err := tt.decode(data); err != nil {
				t.Fatalf("fixture %s does not decode: %v", tt.filename, err)
			}
		})
	}
}

func TestRegionActiveBattleFixtureShapes(t *testing.T) {
	var detail Region
	if err := json.Unmarshal(loadFixture(t, "region.getById.json"), &detail); err != nil {
		t.Fatalf("decode region detail fixture: %v", err)
	}
	if detail.ActiveBattle == nil || detail.ActiveBattle.BattleID == "" || detail.ActiveBattle.Details != nil {
		t.Fatalf("region detail active battle = %#v, want battle ID variant", detail.ActiveBattle)
	}

	var regions map[string]Region
	if err := json.Unmarshal(loadFixture(t, "region.getRegionsObject.json"), &regions); err != nil {
		t.Fatalf("decode regions fixture: %v", err)
	}
	var sawExpandedBattle, sawConstructionNull, sawConstructionBool bool
	for _, region := range regions {
		if region.ActiveBattle != nil && region.ActiveBattle.Details != nil {
			if region.ActiveBattle.Details.ID == "" {
				t.Fatalf("expanded active battle has no ID: %#v", region.ActiveBattle.Details)
			}
			sawExpandedBattle = true
		}
		for _, underConstruction := range []*bool{
			region.UpgradesV2.Upgrades.Base.IsUnderConstruction,
			region.UpgradesV2.Upgrades.Bunker.IsUnderConstruction,
		} {
			if underConstruction == nil {
				sawConstructionNull = true
			} else {
				sawConstructionBool = true
			}
		}
	}
	if !sawExpandedBattle {
		t.Fatal("regions fixture contains no expanded active battle")
	}
	if !sawConstructionNull || !sawConstructionBool {
		t.Fatalf("construction values: null=%t boolean=%t, want both fixture variants", sawConstructionNull, sawConstructionBool)
	}
}

func TestUserSkillFixtureRetainsNullAndNumericValues(t *testing.T) {
	var user User
	if err := json.Unmarshal(loadFixture(t, "user.getUserById.json"), &user); err != nil {
		t.Fatalf("decode user fixture: %v", err)
	}
	attack := user.Skills.Attack
	if attack.Equipment != nil {
		t.Errorf("attack equipment = %v, want nil for JSON null", *attack.Equipment)
	}
	if attack.Value != 100 {
		t.Errorf("attack value = %v, want 100", attack.Value)
	}
	armor := user.Skills.Armor
	if !bytes.Equal(armor.Limited, []byte("null")) {
		t.Errorf("armor limited = %s, want preserved JSON null", armor.Limited)
	}
}

func decodeFixtureAs[T any](data json.RawMessage) error {
	var got T
	return json.Unmarshal(data, &got)
}

func TestGameConfigService_GetGameConfigFixture(t *testing.T) {
	fixture := loadFixture(t, "gameConfig.getGameConfig.json")
	client := NewClient(
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if !strings.HasSuffix(r.URL.Path, "/gameConfig.getGameConfig") {
				t.Errorf("unexpected path: %s", r.URL.Path)
			}
			return fixtureResponse(t, fixture), nil
		})}),
		WithBatchWindow(0),
		WithRateLimit(0),
	)
	defer client.Close()

	config, err := client.GameConfig.GetGameConfig(context.Background())
	if err != nil {
		t.Fatalf("GetGameConfig returned error: %v", err)
	}
	if config == nil || config.Announcement.DurationHours != 24 || config.Badge.BabyBoomer.Reward != 100 {
		t.Fatalf("unexpected decoded config: %+v", config)
	}
	if level := config.Skills.Armor.Levels["1"]; level.Cost != 1 || level.Value != 6 {
		t.Fatalf("unexpected armor level: %+v", level)
	}
	if level := config.UpgradesConfig.Base.Levels["1"]; level.Level != 1 || level.Stats["attackBonus"] != 5 {
		t.Fatalf("unexpected base upgrade level: %+v", level)
	}
}

func fixtureResponse(t *testing.T, data json.RawMessage) *http.Response {
	t.Helper()
	body, err := json.Marshal(trpc.ResponseEnvelope{
		Result: &trpc.ResultEnvelope{Data: data},
	})
	if err != nil {
		t.Fatalf("marshal fixture response: %v", err)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
}
