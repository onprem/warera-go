package warera_test

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/waretest"
)

// ret discards a method's result so each case can return just the error.
func ret[T any](_ T, err error) error { return err }

// serviceCases calls every service method once. Each names the procedure the
// method must call and the exact JSON input it must send ("" for none), so a
// wrong procedure name or input field fails here instead of as a 400 from the
// live API. Iterators are covered through their first page.
var serviceCases = []struct {
	method string
	proc   string
	input  string
	call   func(context.Context, *warera.Client) error
}{
	// Countries, government, and regions.
	{"Countries.GetAll", "country.getAllCountries", "{}", func(ctx context.Context, c *warera.Client) error { return ret(c.Countries.GetAll(ctx)) }},
	{"Countries.GetByID", "country.getCountryById", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Countries.GetByID(ctx, "c1")) }},
	{"Countries.GetDiplomacy", "countryDiplomacy.getByCountry", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Countries.GetDiplomacy(ctx, "c1")) }},
	{"Countries.GetUnrest", "country.getUnrestData", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Countries.GetUnrest(ctx, "c1")) }},
	{"Countries.GetStrategicResourceBonuses", "country.getStrategicResourceBonuses", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Countries.GetStrategicResourceBonuses(ctx, "c1"))
	}},
	{"Government.GetByCountryID", "government.getByCountryId", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Government.GetByCountryID(ctx, "c1")) }},
	{"Government.GetCongressSeats", "government.getCongressSeats", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Government.GetCongressSeats(ctx, "c1"))
	}},
	{"Regions.GetAll", "region.getRegionsObject", "{}", func(ctx context.Context, c *warera.Client) error { return ret(c.Regions.GetAll(ctx)) }},
	{"Regions.GetByID", "region.getById", `{"regionId":"r1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Regions.GetByID(ctx, "r1")) }},
	{"Regions.GetMapData", "map.getMapData", "", func(ctx context.Context, c *warera.Client) error { return ret(c.Regions.GetMapData(ctx)) }},
	{"Upgrades.GetMapUpgrades", "upgrade.getMapUpgrades", "", func(ctx context.Context, c *warera.Client) error { return ret(c.Upgrades.GetMapUpgrades(ctx)) }},
	{"Upgrades.GetByTypeAndEntity", "upgrade.getUpgradeByTypeAndEntity", `{"upgradeType":"bunker","regionId":"r1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Upgrades.GetByTypeAndEntity(ctx, warera.UpgradeFilter{UpgradeType: warera.UpgradeTypeBunker, RegionID: "r1"}))
	}},

	// Politics.
	{"Alliances.GetByID", "alliance.getById", `{"allianceId":"a1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Alliances.GetByID(ctx, "a1")) }},
	{"Alliances.GetMany", "alliance.getByIds", `{"ids":["a1","a2"]}`, func(ctx context.Context, c *warera.Client) error {
		// The recorded response does not contain these IDs, so they are
		// reported as not found.
		if _, err := c.Alliances.GetMany(ctx, []string{"a1", "a2", "a1"}); !warera.IsNotFound(err) {
			return err
		}
		return nil
	}},
	{"Alliances.GetStrategicResourceBonuses", "alliance.getStrategicResourceBonuses", `{"allianceId":"a1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Alliances.GetStrategicResourceBonuses(ctx, "a1"))
	}},
	{"Wars.GetByID", "war.getById", `{"warId":"w1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Wars.GetByID(ctx, "w1")) }},
	{"Wars.GetMany", "war.getById", `{"warId":"w1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Wars.GetMany(ctx, []string{"w1"})) }},
	{"Parties.GetByID", "party.getById", `{"partyId":"p1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Parties.GetByID(ctx, "p1")) }},
	{"Parties.GetMany", "party.getById", `{"partyId":"p1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Parties.GetMany(ctx, []string{"p1"})) }},
	{"Parties.List", "party.getManyPaginated", `{"countryId":"c1","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Parties.List(ctx, warera.PartyFilter{CountryID: "c1", Limit: 5}))
	}},
	{"Parties.All", "party.getManyPaginated", `{"search":"x"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Parties.All(ctx, warera.PartyFilter{Search: "x"}), 1))
	}},
	{"Parties.GetInactiveMemberIDs", "party.getInactiveMemberIds", `{"partyId":"p1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Parties.GetInactiveMemberIDs(ctx, "p1"))
	}},
	{"Elections.GetByID", "election.getElection", `{"electionId":"e1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Elections.GetByID(ctx, "e1")) }},
	{"Elections.GetMany", "election.getElection", `{"electionId":"e1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Elections.GetMany(ctx, []string{"e1"}))
	}},
	{"Elections.List", "election.getElections", `{"partyId":"p1","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Elections.List(ctx, warera.ElectionFilter{PartyID: "p1", Limit: 5}))
	}},
	{"Elections.All", "election.getElections", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Elections.All(ctx, warera.ElectionFilter{CountryID: "c1"}), 1))
	}},

	// Companies, employees, and managers.
	{"Companies.GetByID", "company.getById", `{"companyId":"co1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Companies.GetByID(ctx, "co1")) }},
	{"Companies.GetMany", "company.getById", `{"companyId":"co1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Companies.GetMany(ctx, []string{"co1"}))
	}},
	{"Companies.List", "company.getCompanies", `{"userId":"u1","perPage":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Companies.List(ctx, warera.CompanyFilter{UserID: "u1", Limit: 5}))
	}},
	{"Companies.All", "company.getCompanies", `{"userId":"u1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Companies.All(ctx, warera.CompanyFilter{UserID: "u1"}), 1))
	}},
	{"Companies.GetProductionBonus", "company.getProductionBonus", `{"companyId":"co1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Companies.GetProductionBonus(ctx, "co1"))
	}},
	{"Employees.GetCompanyRoster", "employee.getCompanyRoster", `{"companyId":"co1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Employees.GetCompanyRoster(ctx, "co1"))
	}},
	{"Employees.GetByID", "employee.getById", `{"employeeId":"em1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Employees.GetByID(ctx, "em1")) }},
	{"Managers.GetByCompanyID", "manager.getByCompanyId", `{"companyId":"co1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Managers.GetByCompanyID(ctx, "co1")) }},
	{"Managers.GetOfferByCompanyID", "managerOffer.getByCompanyId", `{"companyId":"co1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Managers.GetOfferByCompanyID(ctx, "co1"))
	}},
	{"Managers.ListOffers", "managerOffer.getManyPaginated", `{"regionId":"r1","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Managers.ListOffers(ctx, warera.ManagerOfferFilter{RegionID: "r1", Limit: 5}))
	}},
	{"Managers.AllOffers", "managerOffer.getManyPaginated", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Managers.AllOffers(ctx, warera.ManagerOfferFilter{CountryID: "c1"}), 1))
	}},
	{"Managers.GetWageStats", "managerOffer.getWageStats", "{}", func(ctx context.Context, c *warera.Client) error { return ret(c.Managers.GetWageStats(ctx)) }},
	{"Me.GetEmployees", "employee.getMyEmployees", "", func(ctx context.Context, c *warera.Client) error { return ret(c.Me.GetEmployees(ctx)) }},
	{"Me.GetWorkforceSummary", "employee.getMyWorkforceSummary", "", func(ctx context.Context, c *warera.Client) error { return ret(c.Me.GetWorkforceSummary(ctx)) }},
	{"Me.GetPendingEmployeeEvents", "employeeEvent.getMine", "", func(ctx context.Context, c *warera.Client) error { return ret(c.Me.GetPendingEmployeeEvents(ctx)) }},
	{"Me.GetPendingEmployeeEventCount", "employeeEvent.getPendingCount", "", func(ctx context.Context, c *warera.Client) error {
		return ret(c.Me.GetPendingEmployeeEventCount(ctx))
	}},
	{"Me.GetEmployeeEventHistory", "employeeEvent.getHistory", `{"employeeId":"e1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Me.GetEmployeeEventHistory(ctx, "e1"))
	}},
	{"MercenaryAuctions.GetBidEligibility", "mercenaryContractAuction.getBidEligibility", `{"muId":"mu1","auctionIds":["a1","a2"]}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.MercenaryAuctions.GetBidEligibility(ctx, "mu1", []string{"a1", "a2"}))
	}},
	{"MilitaryUnits.GetReputationBuyStatus", "mu.getReputationBuyStatus", `{"muId":"mu1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.MilitaryUnits.GetReputationBuyStatus(ctx, "mu1"))
	}},
	{"Managers.GetPointValue", "employee.getManagerPointValue", `{"companyId":"co1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Managers.GetPointValue(ctx, "co1")) }},
	{"Sanctions.List", "sanction.getPaginated", `{"targetUserId":"u1","type":"BAN","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Sanctions.List(ctx, warera.SanctionFilter{TargetUserID: "u1", Type: warera.SanctionBan, Limit: 5}))
	}},
	{"Sanctions.All", "sanction.getPaginated", `{"type":"MUTE_USER"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Sanctions.All(ctx, warera.SanctionFilter{Type: warera.SanctionMuteUser}), 1))
	}},

	// Markets and economy.
	{"Trading.GetPrices", "itemTrading.getPrices", "{}", func(ctx context.Context, c *warera.Client) error { return ret(c.Trading.GetPrices(ctx)) }},
	{"Trading.GetTopOrders", "tradingOrder.getTopOrders", `{"itemCode":"steel","limit":5}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Trading.GetTopOrders(ctx, "steel", 5)) }},
	{"Trading.GetTopOrdersPerItem", "tradingOrder.getTopOrdersPerItemCode", `{"itemCodes":["steel","fish"],"limit":3}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Trading.GetTopOrdersPerItem(ctx, []string{"steel", "fish"}, 3))
	}},
	{"Trading.GetItemTrading", "itemTrading.getItemTrading", `{"itemCode":"bread"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Trading.GetItemTrading(ctx, "bread")) }},
	{"Trading.GetPublicOrdersByOwner", "tradingOrder.getPublicOrdersByOwner", `{"muId":"m1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Trading.GetPublicOrdersByOwner(ctx, warera.OwnerOrdersFilter{MUID: "m1"}))
	}},
	{"MarketStats.GetPPMoneyValue", "gameStat.getPpMoneyValue", "", func(ctx context.Context, c *warera.Client) error { return ret(c.MarketStats.GetPPMoneyValue(ctx)) }},
	{"MarketStats.GetPPMoneyValueByItem", "gameStat.getPpMoneyValueByItem", "", func(ctx context.Context, c *warera.Client) error {
		return ret(c.MarketStats.GetPPMoneyValueByItem(ctx))
	}},
	{"MarketStats.GetWorldDevelopment", "gameStat.getWorldDevelopment", "", func(ctx context.Context, c *warera.Client) error { return ret(c.MarketStats.GetWorldDevelopment(ctx)) }},
	{"MarketStats.GetEquipmentAveragePrice", "gameStat.getEquipmentAvgByCode", `{"itemCode":"rifle"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.MarketStats.GetEquipmentAveragePrice(ctx, "rifle"))
	}},
	{"GameConfig.GetMaxBountyRate", "gameConfig.getMaxBountyRate", "{}", func(ctx context.Context, c *warera.Client) error { return ret(c.GameConfig.GetMaxBountyRate(ctx)) }},
	{"Transactions.List", "transaction.getPaginatedTransactions", `{"userId":"u1","transactionType":["openCase"],"limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Transactions.List(ctx, warera.TransactionFilter{UserID: "u1", TransactionTypes: []warera.TransactionType{warera.TransactionTypeOpenCase}, Limit: 5}))
	}},
	{"Transactions.All", "transaction.getPaginatedTransactions", `{"itemCode":"steel"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Transactions.All(ctx, warera.TransactionFilter{ItemCode: "steel"}), 1))
	}},
	{"Donations.List", "donation.getManyPaginated", `{"partyId":"p1","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Donations.List(ctx, warera.DonationFilter{PartyID: "p1", Limit: 5}))
	}},
	{"Donations.All", "donation.getManyPaginated", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Donations.All(ctx, warera.DonationFilter{CountryID: "c1"}), 1))
	}},
	{"Donations.GetTotal", "donation.getTotalDonations", `{"muId":"m1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Donations.GetTotal(ctx, warera.DonationTotalFilter{MUID: "m1"}))
	}},
	{"MercenaryAuctions.List", "mercenaryContractAuction.getPaginatedAuctions", `{"status":"active","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.MercenaryAuctions.List(ctx, warera.MercenaryAuctionFilter{Status: warera.MercenaryAuctionStatusActive, Limit: 5}))
	}},
	{"MercenaryAuctions.All", "mercenaryContractAuction.getPaginatedAuctions", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.MercenaryAuctions.All(ctx, warera.MercenaryAuctionFilter{CountryID: "c1"}), 1))
	}},

	// Battles.
	{"Battles.GetByID", "battle.getById", `{"battleId":"b1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Battles.GetByID(ctx, "b1")) }},
	{"Battles.GetMany", "battle.getById", `{"battleId":"b1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Battles.GetMany(ctx, []string{"b1"})) }},
	{"Battles.List", "battle.getBattles", `{"countryId":"c1","isActive":true,"limit":5}`, func(ctx context.Context, c *warera.Client) error {
		active := true
		return ret(c.Battles.List(ctx, warera.BattleFilter{CountryID: "c1", IsActive: &active, Limit: 5}))
	}},
	{"Battles.All", "battle.getBattles", `{"allianceId":"a1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Battles.All(ctx, warera.BattleFilter{AllianceID: "a1"}), 1))
	}},
	{"Battles.GetLiveBattleData", "battle.getLiveBattleData", `{"battleId":"b1","roundNumber":3}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Battles.GetLiveBattleData(ctx, "b1", 3))
	}},
	{"Battles.GetActiveIDsByCountry", "battle.getActiveBattlesByCountry", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Battles.GetActiveIDsByCountry(ctx, "c1"))
	}},
	{"BattleOrders.GetByBattle", "battleOrder.getByBattle", `{"battleId":"b1","side":"defender"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.BattleOrders.GetByBattle(ctx, "b1", warera.BattleSideDefender))
	}},
	{"BattleRankings.List", "battleRanking.getRanking", `{"battleId":"b1","dataType":"damage","side":"attacker","type":"user","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.BattleRankings.List(ctx, warera.BattleRankingFilter{BattleID: "b1", DataType: warera.BattleRankingDataTypeDamage, Side: warera.BattleRankingSideAttacker, Type: warera.BattleRankingTypeUser, Limit: 5}))
	}},
	{"BattleRankings.All", "battleRanking.getRanking", `{"warId":"w1","dataType":"money","side":"merged","type":"country"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.BattleRankings.All(ctx, warera.BattleRankingFilter{WarID: "w1", DataType: warera.BattleRankingDataTypeMoney, Side: warera.BattleRankingSideMerged, Type: warera.BattleRankingTypeCountry}), 1))
	}},
	{"BattleRankings.GetPrivateMercenaryDamage", "battleRanking.getPrivateMercenaryDamage", `{"battleId":"b1","roundId":"rd1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.BattleRankings.GetPrivateMercenaryDamage(ctx, "b1", "rd1"))
	}},
	{"LootSummaries.GetByBattleAndUser", "battleLootSummary.getByBattleAndUser", `{"battleId":"b1","userId":"u1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.LootSummaries.GetByBattleAndUser(ctx, "b1", "u1"))
	}},
	{"Rounds.GetByID", "round.getById", `{"roundId":"rd1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Rounds.GetByID(ctx, "rd1")) }},
	{"Rounds.GetLastHits", "round.getLastHits", `{"roundId":"rd1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Rounds.GetLastHits(ctx, "rd1")) }},

	// Military units, tournaments, state armies, and contributions.
	{"MilitaryUnits.GetByID", "mu.getById", `{"muId":"m1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.MilitaryUnits.GetByID(ctx, "m1")) }},
	{"MilitaryUnits.GetMany", "mu.getById", `{"muId":"m1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.MilitaryUnits.GetMany(ctx, []string{"m1"}))
	}},
	{"MilitaryUnits.List", "mu.getManyPaginated", `{"search":"x","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.MilitaryUnits.List(ctx, warera.MilitaryUnitFilter{Search: "x", Limit: 5}))
	}},
	{"MilitaryUnits.All", "mu.getManyPaginated", `{"userId":"u1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.MilitaryUnits.All(ctx, warera.MilitaryUnitFilter{UserID: "u1"}), 1))
	}},
	{"MilitaryUnits.GetMembers", "muMember.getByMu", `{"muId":"m1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.MilitaryUnits.GetMembers(ctx, "m1")) }},
	{"MilitaryUnits.GetLeague", "muLeague.getMuLeague", `{"muId":"m1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.MilitaryUnits.GetLeague(ctx, "m1")) }},
	{"MilitaryUnits.GetLeagueDivisionCounts", "muLeague.getDivisionCounts", "{}", func(ctx context.Context, c *warera.Client) error {
		return ret(c.MilitaryUnits.GetLeagueDivisionCounts(ctx))
	}},
	{"MilitaryUnits.GetSponsorship", "muSponsorship.getForMu", `{"muId":"m1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.MilitaryUnits.GetSponsorship(ctx, "m1"))
	}},
	{"MilitaryUnits.GetSponsorshipsByCountry", "muSponsorship.getByCountry", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.MilitaryUnits.GetSponsorshipsByCountry(ctx, "c1"))
	}},
	{"MilitaryUnits.ListActiveSponsorships", "muSponsorship.getActiveByCountryPaginated", `{"countryId":"c1","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.MilitaryUnits.ListActiveSponsorships(ctx, warera.MUSponsorshipFilter{CountryID: "c1", Limit: 5}))
	}},
	{"MilitaryUnits.AllActiveSponsorships", "muSponsorship.getActiveByCountryPaginated", `{"countryId":"c2"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.MilitaryUnits.AllActiveSponsorships(ctx, warera.MUSponsorshipFilter{CountryID: "c2"}), 1))
	}},
	{"MilitaryUnits.GetSponsorable", "muSponsorship.getSponsorable", `{"countryId":"c1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.MilitaryUnits.GetSponsorable(ctx, "c1"))
	}},
	{"Tournaments.GetByID", "tournament.getById", `{"tournamentId":"t1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Tournaments.GetByID(ctx, "t1")) }},
	{"Tournaments.GetMany", "tournament.getById", `{"tournamentId":"t1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Tournaments.GetMany(ctx, []string{"t1"}))
	}},
	{"Tournaments.List", "tournament.getManyPaginated", `{"limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Tournaments.List(ctx, warera.TournamentFilter{Limit: 5}))
	}},
	{"Tournaments.All", "tournament.getManyPaginated", `{"limit":2}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Tournaments.All(ctx, warera.TournamentFilter{Limit: 2}), 1))
	}},
	{"Tournaments.GetTeam", "tournamentTeam.getById", `{"tournamentTeamId":"tt1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Tournaments.GetTeam(ctx, "tt1")) }},
	{"Tournaments.GetTeams", "tournamentTeam.getByTournamentId", `{"tournamentId":"t1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Tournaments.GetTeams(ctx, "t1")) }},
	{"StateArmy.GetOverview", "stateArmyExpense.getOverview", `{"countryId":"c1","days":7}`, func(ctx context.Context, c *warera.Client) error { return ret(c.StateArmy.GetOverview(ctx, "c1", 7)) }},
	{"StateArmy.ListLedger", "stateArmyExpense.getLedger", `{"countryId":"c1","kind":"hqOil","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.StateArmy.ListLedger(ctx, warera.StateArmyLedgerFilter{CountryID: "c1", Kind: warera.StateArmyExpenseHQOil, Limit: 5}))
	}},
	{"StateArmy.AllLedger", "stateArmyExpense.getLedger", `{"countryId":"c1","muId":"m1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.StateArmy.AllLedger(ctx, warera.StateArmyLedgerFilter{CountryID: "c1", MUID: "m1"}), 1))
	}},
	{"Contributions.ListCountryUnrest", "contribution.getCountryUnrestContributions", `{"countryId":"c1","page":2}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Contributions.ListCountryUnrest(ctx, "c1", 2))
	}},
	{"Contributions.ListRegion", "contribution.getRegionContributions", `{"regionId":"r1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Contributions.ListRegion(ctx, "r1", 0))
	}},
	{"Contributions.AllCountryUnrest", "contribution.getCountryUnrestContributions", `{"countryId":"c1","page":1}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Contributions.AllCountryUnrest(ctx, "c1"), 1))
	}},
	{"Contributions.AllRegion", "contribution.getRegionContributions", `{"regionId":"r1","page":1}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Contributions.AllRegion(ctx, "r1"), 1))
	}},

	// Players, community, and rankings.
	{"Users.GetByID", "user.getUserById", `{"userId":"u1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Users.GetByID(ctx, "u1")) }},
	{"Users.GetMany", "user.getUserById", `{"userId":"u1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Users.GetMany(ctx, []string{"u1"})) }},
	{"Users.GetLiteByID", "user.getUserLite", `{"userId":"u1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Users.GetLiteByID(ctx, "u1")) }},
	{"Users.GetLiteMany", "user.getUserLite", `{"userId":"u1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Users.GetLiteMany(ctx, []string{"u1"}))
	}},
	{"Users.GetSkills", "user.getUserSkills", `{"userId":"u1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Users.GetSkills(ctx, "u1")) }},
	{"Users.GetSkillsMany", "user.getUserSkills", `{"userId":"u1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Users.GetSkillsMany(ctx, []string{"u1"}))
	}},
	{"Users.List", "user.getUsersByCountry", `{"countryId":"c1","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Users.List(ctx, warera.UserFilter{CountryID: "c1", Limit: 5}))
	}},
	{"Users.All", "user.getUsersByCountry", `{"countryId":"c2"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Users.All(ctx, warera.UserFilter{CountryID: "c2"}), 1))
	}},
	{"Inventory.GetCurrentEquipment", "inventory.fetchCurrentEquipment", `{"userId":"u1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Inventory.GetCurrentEquipment(ctx, "u1"))
	}},
	{"Leagues.GetUserLeague", "league.getUserLeague", `{"userId":"u1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Leagues.GetUserLeague(ctx, "u1")) }},
	{"Search.SearchAnything", "search.searchAnything", `{"searchText":"bob"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Search.SearchAnything(ctx, "bob")) }},
	{"Search.SearchUsers", "search.searchAnything", `{"searchText":"ann"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Search.SearchUsers(ctx, "ann")) }},
	{"Events.List", "event.getEventsPaginated", `{"countryId":"c1","eventTypes":["warDeclared"],"limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Events.List(ctx, warera.EventFilter{CountryID: "c1", EventTypes: []warera.EventType{warera.EventTypeWarDeclared}, Limit: 5}))
	}},
	{"Events.All", "event.getEventsPaginated", `{"countryId":"c2"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Events.All(ctx, warera.EventFilter{CountryID: "c2"}), 1))
	}},
	{"Articles.GetByID", "article.getArticleById", `{"articleId":"ar1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Articles.GetByID(ctx, "ar1")) }},
	{"Articles.GetMany", "article.getArticleById", `{"articleId":"ar1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Articles.GetMany(ctx, []string{"ar1"}))
	}},
	{"Articles.GetLiteByID", "article.getArticleLiteById", `{"articleId":"ar1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Articles.GetLiteByID(ctx, "ar1")) }},
	{"Articles.List", "article.getArticlesPaginated", `{"type":"weekly","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Articles.List(ctx, warera.ArticleFilter{Type: warera.ArticleListTypeWeekly, Limit: 5}))
	}},
	{"Articles.All", "article.getArticlesPaginated", `{"type":"top","userId":"u1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Articles.All(ctx, warera.ArticleFilter{Type: warera.ArticleListTypeTop, UserID: "u1"}), 1))
	}},
	{"Polls.GetByID", "poll.getPoll", `{"pollId":"po1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Polls.GetByID(ctx, "po1")) }},
	{"Polls.ListVoters", "poll.getPollVoters", `{"pollId":"po1","optionId":2,"limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Polls.ListVoters(ctx, warera.PollVoterFilter{PollID: "po1", OptionID: 2, Limit: 5}))
	}},
	{"Polls.AllVoters", "poll.getPollVoters", `{"pollId":"po1","optionId":0}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Polls.AllVoters(ctx, warera.PollVoterFilter{PollID: "po1"}), 1))
	}},
	{"Rankings.GetRanking", "ranking.getRanking", `{"rankingType":"userWealth"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Rankings.GetRanking(ctx, warera.RankingTypeUserWealth))
	}},
	{"Rankings.List", "ranking.getRankingPage", `{"rankingType":"muDamages","countryId":"c1","limit":5}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Rankings.List(ctx, warera.RankingPageFilter{RankingType: warera.RankingTypeMUDamages, CountryID: "c1", Limit: 5}))
	}},
	{"Rankings.All", "ranking.getRankingPage", `{"rankingType":"countryWealth"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(warera.Collect(c.Rankings.All(ctx, warera.RankingPageFilter{RankingType: warera.RankingTypeCountryWealth}), 1))
	}},
	{"Rankings.GetAround", "ranking.getRankingAround", `{"rankingType":"userLevel","entityId":"u1"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Rankings.GetAround(ctx, warera.RankingAroundFilter{RankingType: warera.RankingTypeUserLevel, EntityID: "u1"}))
	}},
	{"Rankings.GetHistory", "rankingSnapshot.getHistory", `{"entityId":"u1","rankingType":"userBounty","range":"30d"}`, func(ctx context.Context, c *warera.Client) error {
		return ret(c.Rankings.GetHistory(ctx, warera.RankingHistoryFilter{EntityID: "u1", RankingType: warera.RankingTypeUserBounty, Range: warera.RankingHistoryRange30Days}))
	}},
	{"Rankings.GetRecentRanks", "rankingSnapshot.getEntityRecentRanks", `{"entityId":"u1"}`, func(ctx context.Context, c *warera.Client) error { return ret(c.Rankings.GetRecentRanks(ctx, "u1")) }},
	{"GameConfig.GetDates", "gameConfig.getDates", "{}", func(ctx context.Context, c *warera.Client) error { return ret(c.GameConfig.GetDates(ctx)) }},
	{"GameConfig.Get", "gameConfig.getGameConfig", "{}", func(ctx context.Context, c *warera.Client) error { return ret(c.GameConfig.Get(ctx)) }},
}

func TestServiceMethodsCallTheirProcedures(t *testing.T) {
	srv := waretest.NewServer(t)
	client := warera.NewClient(warera.WithBaseURL(srv.URL), warera.WithBatchWindow(0), warera.WithRateLimit(0))
	defer client.Close()
	spec := loadSpecInputs(t)

	for _, tc := range serviceCases {
		t.Run(tc.method, func(t *testing.T) {
			before := len(srv.Inputs(tc.proc))
			if err := tc.call(context.Background(), client); err != nil {
				t.Fatalf("call failed: %v", err)
			}
			inputs := srv.Inputs(tc.proc)
			if len(inputs) <= before {
				t.Fatalf("%s was not called", tc.proc)
			}
			got := inputs[before]
			if !jsonEqual(got, tc.input) {
				t.Errorf("input = %s, want %s", got, tc.input)
			}
			// For procedures in the published OpenAPI document, check the
			// input against the API's own schema too.
			if schema, ok := spec[tc.proc]; ok && len(got) > 0 {
				var fields map[string]any
				if err := json.Unmarshal(got, &fields); err != nil {
					t.Fatalf("input is not an object: %s", got)
				}
				for name := range fields {
					if _, known := schema.Properties[name]; !known {
						t.Errorf("input field %q is not in the OpenAPI schema", name)
					}
				}
				for _, name := range schema.Required {
					if _, present := fields[name]; !present {
						t.Errorf("required input field %q is missing", name)
					}
				}
			}
		})
	}
}

// TestServiceCasesCoverEveryMethod fails when a service method is added
// without a case above.
func TestServiceCasesCoverEveryMethod(t *testing.T) {
	covered := map[string]bool{}
	for _, tc := range serviceCases {
		covered[tc.method] = true
	}
	c := warera.NewClient()
	defer c.Close()
	clientValue := reflect.ValueOf(c).Elem()
	for i := range clientValue.NumField() {
		field := clientValue.Type().Field(i)
		if !field.IsExported() || field.Type.Kind() != reflect.Pointer || !strings.HasSuffix(field.Type.Elem().Name(), "Service") {
			continue
		}
		for j := range field.Type.NumMethod() {
			name := field.Name + "." + field.Type.Method(j).Name
			if !covered[name] {
				t.Errorf("no service case for %s", name)
			}
		}
	}
}

type specInput struct {
	Properties map[string]json.RawMessage `json:"properties"`
	Required   []string                   `json:"required"`
}

func loadSpecInputs(t *testing.T) map[string]specInput {
	t.Helper()
	data, err := os.ReadFile("spec/openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	var spec struct {
		Paths map[string]struct {
			Post struct {
				RequestBody *struct {
					Content struct {
						JSON struct {
							Schema specInput `json:"schema"`
						} `json:"application/json"`
					} `json:"content"`
				} `json:"requestBody"`
			} `json:"post"`
		} `json:"paths"`
	}
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatal(err)
	}
	inputs := map[string]specInput{}
	for path, item := range spec.Paths {
		if item.Post.RequestBody != nil {
			inputs[strings.TrimPrefix(path, "/")] = item.Post.RequestBody.Content.JSON.Schema
		}
	}
	return inputs
}

func jsonEqual(got json.RawMessage, want string) bool {
	if want == "" {
		return len(got) == 0
	}
	var a, b any
	if json.Unmarshal(got, &a) != nil || json.Unmarshal([]byte(want), &b) != nil {
		return false
	}
	return reflect.DeepEqual(a, b)
}
