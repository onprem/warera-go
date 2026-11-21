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
		{"company production bonus", "company.getProductionBonus.json", decodeFixtureAs[CompanyProductionBonus]},
		{"company roster", "employee.getCompanyRoster.json", decodeFixtureAs[CompanyRoster]},
		{"employee", "employee.getById.json", decodeFixtureAs[Employee]},
		{"company manager", "manager.getByCompanyId.json", decodeFixtureAs[Manager]},
		{"manager offer", "managerOffer.getByCompanyId.json", decodeFixtureAs[ManagerOffer]},
		{"manager offers", "managerOffer.getManyPaginated.json", decodeFixtureAs[CursorPage[ManagerOffer]]},
		{"manager wage stats", "managerOffer.getWageStats.json", decodeFixtureAs[ManagerWageStats]},
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
		{"ranking page", "ranking.getRankingPage.json", decodeFixtureAs[RankingPage]},
		{"ranking around", "ranking.getRankingAround.json", decodeFixtureAs[[]RankingEntry]},
		{"poll", "poll.getPoll.json", decodeFixtureAs[Poll]},
		{"poll voters", "poll.getPollVoters.json", decodeFixtureAs[CursorPage[PollVote]]},
		{"transactions", "transaction.getPaginatedTransactions.json", decodeFixtureAs[CursorPage[TransactionListItem]]},
		{"region detail", "region.getById.json", decodeFixtureAs[Region]},
		{"region map", "region.getRegionsObject.json", decodeFixtureAs[map[string]Region]},
		{"round detail", "round.getById.json", decodeFixtureAs[Round]},
		{"round last hits", "round.getLastHits.json", decodeFixtureAs[RoundLastHits]},
		{"search results", "search.searchAnything.json", decodeFixtureAs[SearchResult]},
		{"top orders", "tradingOrder.getTopOrders.json", decodeFixtureAs[TopOrders]},
		{"upgrade", "upgrade.getUpgradeByTypeAndEntity.json", decodeFixtureAs[Upgrade]},
		{"user detail", "user.getUserById.json", decodeFixtureAs[User]},
		{"user lite", "user.getUserLite.json", decodeFixtureAs[UserLite]},
		{"user skills", "user.getUserSkills.json", decodeFixtureAs[UserSkillProfile]},
		{"users by country", "user.getUsersByCountry.json", decodeFixtureAs[CursorPage[UserListItem]]},
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

func TestUserSkillFixtureDistinguishesAbsentAndNumericValues(t *testing.T) {
	var user User
	if err := json.Unmarshal(loadFixture(t, "user.getUserById.json"), &user); err != nil {
		t.Fatalf("decode user fixture: %v", err)
	}
	attack := user.Skills.Attack
	if attack.Equipment != nil {
		t.Errorf("attack equipment = %v, want nil when the API omits it", *attack.Equipment)
	}
	if attack.Value != 275 || attack.Weapon == nil || *attack.Weapon != 163 {
		t.Errorf("attack value/weapon = %v/%v, want 275/163", attack.Value, attack.Weapon)
	}
	if !bytes.Equal(attack.Overflow, []byte("72")) {
		t.Errorf("attack overflow = %s, want preserved JSON 72", attack.Overflow)
	}
	armor := user.Skills.Armor
	if armor.TotalAfterSoftCap == nil || *armor.TotalAfterSoftCap != 68 {
		t.Errorf("armor total after soft cap = %v, want 68", armor.TotalAfterSoftCap)
	}
}

// decodeFixtureAs decodes a fixture strictly, so a field the API returns but the
// model lacks (or has renamed) fails the test instead of silently decoding to
// a zero value. MongoDB's "__v" document version is dropped first because it is
// not part of any model's contract. Custom UnmarshalJSON methods decode their
// own subtree leniently.
func decodeFixtureAs[T any](data json.RawMessage) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var raw any
	if err := dec.Decode(&raw); err != nil {
		return err
	}
	stripped, err := json.Marshal(withoutDocumentVersions(raw))
	if err != nil {
		return err
	}
	dec = json.NewDecoder(bytes.NewReader(stripped))
	dec.DisallowUnknownFields()
	var got T
	return dec.Decode(&got)
}

func withoutDocumentVersions(value any) any {
	switch v := value.(type) {
	case map[string]any:
		delete(v, "__v")
		for key, child := range v {
			v[key] = withoutDocumentVersions(child)
		}
	case []any:
		for i, child := range v {
			v[i] = withoutDocumentVersions(child)
		}
	}
	return value
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
	if config == nil || config.Announcement.DurationHours != 24 || config.Badge["babyBoomer"].Reward != 100 {
		t.Fatalf("unexpected decoded config: %+v", config)
	}
	if level := config.Skills["armor"].Levels["1"]; level.Cost != 1 || level.Value != 6 {
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

func TestRankingPageFixtureFillsEmbeddedCursorPage(t *testing.T) {
	var page RankingPage
	if err := json.Unmarshal(loadFixture(t, "ranking.getRankingPage.json"), &page); err != nil {
		t.Fatalf("decode ranking page fixture: %v", err)
	}
	if len(page.Items) != 5 || page.Items[0].Rank != 1 || !page.HasNext() {
		t.Fatalf("page items = %d, first rank = %v, has next = %t; want 5 entries from rank 1 and a next cursor",
			len(page.Items), page.Items, page.HasNext())
	}
	if page.Ranking == nil || page.Ranking.Count != 181 || page.Ranking.TierValues.Master == 0 {
		t.Fatalf("first-page ranking summary = %+v, want count 181 with tier values", page.Ranking)
	}
}

func TestManagersGetByCompanyIDReportsMissingManagerAsNil(t *testing.T) {
	for _, tt := range []struct {
		name        string
		data        json.RawMessage
		wantManager bool
	}{
		{"company with a manager", loadFixture(t, "manager.getByCompanyId.json"), true},
		{"company without a manager", json.RawMessage("null"), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(
				WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					return fixtureResponse(t, tt.data), nil
				})}),
				WithBatchWindow(0),
				WithRateLimit(0),
			)
			defer client.Close()
			manager, err := client.Managers.GetByCompanyID(context.Background(), "6a2e45271cbde97f9cc3422f")
			if err != nil {
				t.Fatalf("GetByCompanyID: %v", err)
			}
			if !tt.wantManager {
				if manager != nil {
					t.Fatalf("manager = %+v, want nil", manager)
				}
				return
			}
			if manager == nil || manager.UserID != "6997f39ec016842739f5dca1" || manager.ManagementLevel != 1 {
				t.Fatalf("manager = %+v, want the recorded manager 6997f39ec016842739f5dca1 at level 1", manager)
			}
		})
	}
}
