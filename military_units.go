package warera

import (
	"context"
	"iter"
)

// The procedures in this file are served by the WarEra API but are not listed
// in its published OpenAPI document, so they may change without notice. They
// require an API key.

// GetMembers returns the damage and help statistics of each member of a
// military unit.
func (s *MilitaryUnitService) GetMembers(ctx context.Context, muID string) ([]MUMember, error) {
	var members []MUMember
	if err := s.client.trpc.Call(ctx, "muMember.getByMu", getMilitaryUnitByIDInput{MUID: muID}, &members); err != nil {
		return nil, err
	}
	return members, nil
}

// GetLeague returns a military unit's group in this week's MU league, with the
// group's standings and how the unit's reward would be shared among its
// members. League and Me are nil when the unit is not in a league.
func (s *MilitaryUnitService) GetLeague(ctx context.Context, muID string) (*MULeague, error) {
	var league MULeague
	if err := s.client.trpc.Call(ctx, "muLeague.getMuLeague", getMilitaryUnitByIDInput{MUID: muID}, &league); err != nil {
		return nil, err
	}
	return &league, nil
}

// GetLeagueDivisionCounts returns how many military units are in each MU
// league division this week.
func (s *MilitaryUnitService) GetLeagueDivisionCounts(ctx context.Context) (*MULeagueDivisionCounts, error) {
	var counts MULeagueDivisionCounts
	if err := s.client.trpc.Call(ctx, "muLeague.getDivisionCounts", struct{}{}, &counts); err != nil {
		return nil, err
	}
	return &counts, nil
}

// GetSponsorship returns a military unit's active or offered state
// sponsorship, or nil when it has none.
func (s *MilitaryUnitService) GetSponsorship(ctx context.Context, muID string) (*MUSponsorship, error) {
	var sponsorship *MUSponsorship
	if err := s.client.trpc.Call(ctx, "muSponsorship.getForMu", getMilitaryUnitByIDInput{MUID: muID}, &sponsorship); err != nil {
		return nil, err
	}
	return sponsorship, nil
}

type getMUSponsorshipsByCountryInput struct {
	CountryID string `json:"countryId"`
}

// GetSponsorshipsByCountry returns a country's active, offered, and past
// military unit sponsorships, with the sponsored units' names and sizes.
func (s *MilitaryUnitService) GetSponsorshipsByCountry(ctx context.Context, countryID string) (*MUSponsorshipsByCountry, error) {
	var sponsorships MUSponsorshipsByCountry
	if err := s.client.trpc.Call(ctx, "muSponsorship.getByCountry", getMUSponsorshipsByCountryInput{CountryID: countryID}, &sponsorships); err != nil {
		return nil, err
	}
	return &sponsorships, nil
}

// ListActiveSponsorships returns one page of a country's active military unit
// sponsorships.
func (s *MilitaryUnitService) ListActiveSponsorships(ctx context.Context, filter MUSponsorshipFilter) (*MUSponsorshipPage, error) {
	var page MUSponsorshipPage
	if err := s.client.trpc.Call(ctx, "muSponsorship.getActiveByCountryPaginated", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// AllActiveSponsorships returns an iterator over every active military unit
// sponsorship of a country.
func (s *MilitaryUnitService) AllActiveSponsorships(ctx context.Context, filter MUSponsorshipFilter) iter.Seq2[MUSponsorship, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[MUSponsorship], error) {
		f := filter
		f.Cursor = cursor
		page, err := s.ListActiveSponsorships(c, f)
		if err != nil {
			return nil, err
		}
		return &page.CursorPage, nil
	})
}

// GetSponsorable returns the military units based in a country that it could
// sponsor.
func (s *MilitaryUnitService) GetSponsorable(ctx context.Context, countryID string) ([]MUSponsorshipCandidate, error) {
	var units []MUSponsorshipCandidate
	if err := s.client.trpc.Call(ctx, "muSponsorship.getSponsorable", getMUSponsorshipsByCountryInput{CountryID: countryID}, &units); err != nil {
		return nil, err
	}
	return units, nil
}

// MUSponsorshipFilter specifies options for listing a country's active
// military unit sponsorships.
type MUSponsorshipFilter struct {
	CountryID string `json:"countryId"`
	Cursor    string `json:"cursor,omitempty"`
	Limit     int    `json:"limit,omitempty"`
}

// StateArmyService handles what countries spend supporting their state army:
// oil for military unit headquarters, mercenary contracts, and sponsorship
// laws.
type StateArmyService struct {
	client *Client
}

type getStateArmyOverviewInput struct {
	CountryID string `json:"countryId"`
	Days      int    `json:"days,omitempty"`
}

// GetOverview returns a country's state army spending over the last days days
// (1 to 365), totaled by kind, by military unit, and per day. With days 0 the
// API uses its default period and leaves out the per-day totals.
func (s *StateArmyService) GetOverview(ctx context.Context, countryID string, days int) (*StateArmyOverview, error) {
	var overview StateArmyOverview
	if err := s.client.trpc.Call(ctx, "stateArmyExpense.getOverview", getStateArmyOverviewInput{CountryID: countryID, Days: days}, &overview); err != nil {
		return nil, err
	}
	return &overview, nil
}

// ListLedger returns one page of a country's state army expenses, newest
// first.
func (s *StateArmyService) ListLedger(ctx context.Context, filter StateArmyLedgerFilter) (*CursorPage[StateArmyExpense], error) {
	var page CursorPage[StateArmyExpense]
	if err := s.client.trpc.Call(ctx, "stateArmyExpense.getLedger", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// AllLedger returns an iterator over every state army expense of a country
// matching the filter, newest first.
func (s *StateArmyService) AllLedger(ctx context.Context, filter StateArmyLedgerFilter) iter.Seq2[StateArmyExpense, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[StateArmyExpense], error) {
		f := filter
		f.Cursor = cursor
		return s.ListLedger(c, f)
	})
}

// StateArmyLedgerFilter specifies options for listing a country's state army
// expenses.
type StateArmyLedgerFilter struct {
	CountryID string               `json:"countryId"`
	Cursor    string               `json:"cursor,omitempty"`
	Kind      StateArmyExpenseKind `json:"kind,omitempty"`
	Limit     int                  `json:"limit,omitempty"`
	MUID      string               `json:"muId,omitempty"`
}

// StateArmyExpenseKind is what a state army expense paid for.
type StateArmyExpenseKind string

const (
	// StateArmyExpenseHQOil is oil sent to a military unit's headquarters.
	StateArmyExpenseHQOil StateArmyExpenseKind = "hqOil"
	// StateArmyExpenseContract is a mercenary contract payment.
	StateArmyExpenseContract StateArmyExpenseKind = "contract"
	// StateArmyExpenseSponsorLaw is items given under a sponsorship law.
	StateArmyExpenseSponsorLaw StateArmyExpenseKind = "sponsorLaw"
	// StateArmyExpenseRevokeLaw is a sponsorship law being revoked.
	StateArmyExpenseRevokeLaw StateArmyExpenseKind = "revokeLaw"
)

// TournamentService handles MU and country tournaments and their teams.
type TournamentService struct {
	client *Client
}

type getTournamentInput struct {
	TournamentID string `json:"tournamentId"`
}

type getTournamentTeamInput struct {
	TournamentTeamID string `json:"tournamentTeamId"`
}

// GetByID returns a tournament, including its bracket, by ID.
func (s *TournamentService) GetByID(ctx context.Context, tournamentID string) (*Tournament, error) {
	var tournament Tournament
	if err := s.client.trpc.Call(ctx, "tournament.getById", getTournamentInput{TournamentID: tournamentID}, &tournament); err != nil {
		return nil, err
	}
	return &tournament, nil
}

// List returns one page of tournaments, newest first.
func (s *TournamentService) List(ctx context.Context, filter TournamentFilter) (*CursorPage[Tournament], error) {
	var page CursorPage[Tournament]
	if err := s.client.trpc.Call(ctx, "tournament.getManyPaginated", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over every tournament, newest first.
func (s *TournamentService) All(ctx context.Context, filter TournamentFilter) iter.Seq2[Tournament, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[Tournament], error) {
		f := filter
		f.Cursor = cursor
		return s.List(c, f)
	})
}

// GetTeam returns a tournament team by ID.
func (s *TournamentService) GetTeam(ctx context.Context, tournamentTeamID string) (*TournamentTeam, error) {
	var team TournamentTeam
	if err := s.client.trpc.Call(ctx, "tournamentTeam.getById", getTournamentTeamInput{TournamentTeamID: tournamentTeamID}, &team); err != nil {
		return nil, err
	}
	return &team, nil
}

// GetTeams returns every team of a tournament.
func (s *TournamentService) GetTeams(ctx context.Context, tournamentID string) ([]TournamentTeam, error) {
	var teams []TournamentTeam
	if err := s.client.trpc.Call(ctx, "tournamentTeam.getByTournamentId", getTournamentInput{TournamentID: tournamentID}, &teams); err != nil {
		return nil, err
	}
	return teams, nil
}

// TournamentFilter specifies options for listing tournaments.
type TournamentFilter struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type getActiveBattlesByCountryInput struct {
	CountryID string `json:"countryId"`
}

// GetActiveIDsByCountry returns the IDs of the active battles a country is
// fighting in.
func (s *BattleService) GetActiveIDsByCountry(ctx context.Context, countryID string) ([]string, error) {
	var ids []string
	if err := s.client.trpc.Call(ctx, "battle.getActiveBattlesByCountry", getActiveBattlesByCountryInput{CountryID: countryID}, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

type getPrivateMercenaryDamageInput struct {
	BattleID string `json:"battleId"`
	RoundID  string `json:"roundId"`
}

// GetPrivateMercenaryDamage returns the damage that mercenaries under private
// contracts dealt for each side of a battle round.
func (s *BattleRankingService) GetPrivateMercenaryDamage(ctx context.Context, battleID, roundID string) (*BattleMercenaryDamage, error) {
	var damage BattleMercenaryDamage
	if err := s.client.trpc.Call(ctx, "battleRanking.getPrivateMercenaryDamage", getPrivateMercenaryDamageInput{BattleID: battleID, RoundID: roundID}, &damage); err != nil {
		return nil, err
	}
	return &damage, nil
}

// ContributionService handles the health and energy players spend supporting
// or opposing country unrest and region resistance.
type ContributionService struct {
	client *Client
}

type getCountryContributionsInput struct {
	CountryID string `json:"countryId"`
	Page      int    `json:"page,omitempty"`
}

type getRegionContributionsInput struct {
	RegionID string `json:"regionId"`
	Page     int    `json:"page,omitempty"`
}

// ListCountryUnrest returns one page of contributions to a country's unrest,
// newest first. Pages are numbered from 1, and 0 requests the first page.
func (s *ContributionService) ListCountryUnrest(ctx context.Context, countryID string, page int) (*ContributionPage, error) {
	var result ContributionPage
	if err := s.client.trpc.Call(ctx, "contribution.getCountryUnrestContributions", getCountryContributionsInput{CountryID: countryID, Page: page}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListRegion returns one page of contributions to a region's resistance,
// newest first. Pages are numbered from 1, and 0 requests the first page.
func (s *ContributionService) ListRegion(ctx context.Context, regionID string, page int) (*ContributionPage, error) {
	var result ContributionPage
	if err := s.client.trpc.Call(ctx, "contribution.getRegionContributions", getRegionContributionsInput{RegionID: regionID, Page: page}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// LeagueService handles the weekly player leagues.
type LeagueService struct {
	client *Client
}

// GetUserLeague returns a user's group in this week's player league and the
// group's standings. League and Me are nil when the user is not in a league.
func (s *LeagueService) GetUserLeague(ctx context.Context, userID string) (*League, error) {
	var league League
	if err := s.client.trpc.Call(ctx, "league.getUserLeague", getUserByIdInput{UserID: userID}, &league); err != nil {
		return nil, err
	}
	return &league, nil
}

// GetHistory returns an entity's daily rank and value in a ranking, oldest
// first. Days without a snapshot are left out, and an entity that does not
// belong in the ranking, such as a user in a country ranking, has no history.
func (s *RankingService) GetHistory(ctx context.Context, filter RankingHistoryFilter) ([]RankingSnapshot, error) {
	var snapshots []RankingSnapshot
	if err := s.client.trpc.Call(ctx, "rankingSnapshot.getHistory", filter, &snapshots); err != nil {
		return nil, err
	}
	return snapshots, nil
}

type getEntityRecentRanksInput struct {
	EntityID string `json:"entityId"`
}

// GetRecentRanks returns a user's, country's, military unit's, or alliance's
// daily rank in every ranking it belongs to over roughly the last ten days.
func (s *RankingService) GetRecentRanks(ctx context.Context, entityID string) (*RankingSnapshotRecentRanks, error) {
	var ranks RankingSnapshotRecentRanks
	if err := s.client.trpc.Call(ctx, "rankingSnapshot.getEntityRecentRanks", getEntityRecentRanksInput{EntityID: entityID}, &ranks); err != nil {
		return nil, err
	}
	return &ranks, nil
}

// RankingHistoryFilter selects an entity's history in one ranking.
type RankingHistoryFilter struct {
	// EntityID is the ID of the user, country, military unit, or alliance.
	EntityID    string              `json:"entityId"`
	RankingType RankingType         `json:"rankingType"`
	Range       RankingHistoryRange `json:"range"`
}

// RankingHistoryRange is the period [RankingService.GetHistory] returns.
type RankingHistoryRange string

const (
	RankingHistoryRange7Days  RankingHistoryRange = "7d"
	RankingHistoryRange30Days RankingHistoryRange = "30d"
	RankingHistoryRange90Days RankingHistoryRange = "90d"
	RankingHistoryRange1Year  RankingHistoryRange = "1y"
)
