package warera

import (
	"context"
	"iter"
)

// The procedures in this file are served by the WarEra API but are not listed
// in its published OpenAPI document, so they may change without notice. They
// require an API key.

// --- Alliance Service ---

// AllianceService handles alliances of countries. Alliance IDs appear on
// countries as Country.AllianceID.
type AllianceService struct {
	client *Client
}

type getAllianceInput struct {
	AllianceID string `json:"allianceId"`
}

type getAlliancesInput struct {
	IDs []string `json:"ids"`
}

// GetByID returns an alliance by ID.
func (s *AllianceService) GetByID(ctx context.Context, allianceID string) (*Alliance, error) {
	var alliance Alliance
	if err := s.client.trpc.Call(ctx, "alliance.getById", getAllianceInput{AllianceID: allianceID}, &alliance); err != nil {
		return nil, err
	}
	return &alliance, nil
}

// GetMany returns the alliances with the given IDs in one request. The API
// accepts at most 100 IDs and omits IDs it does not find, so the result may be
// shorter than allianceIDs and is not guaranteed to follow its order.
func (s *AllianceService) GetMany(ctx context.Context, allianceIDs []string) ([]Alliance, error) {
	if len(allianceIDs) == 0 {
		return []Alliance{}, nil
	}
	var alliances []Alliance
	if err := s.client.trpc.Call(ctx, "alliance.getByIds", getAlliancesInput{IDs: allianceIDs}, &alliances); err != nil {
		return nil, err
	}
	return alliances, nil
}

// GetStrategicResourceBonuses returns the strategic resources an alliance's
// countries control and the bonuses they grant.
func (s *AllianceService) GetStrategicResourceBonuses(ctx context.Context, allianceID string) (*StrategicResourceBonuses, error) {
	var bonuses StrategicResourceBonuses
	if err := s.client.trpc.Call(ctx, "alliance.getStrategicResourceBonuses", getAllianceInput{AllianceID: allianceID}, &bonuses); err != nil {
		return nil, err
	}
	return &bonuses, nil
}

// --- War Service ---

// WarService handles wars between countries. War IDs appear on battles as
// Battle.WarID.
type WarService struct {
	client *Client
}

type getWarInput struct {
	WarID string `json:"warId"`
}

// GetByID returns a war by ID.
func (s *WarService) GetByID(ctx context.Context, warID string) (*War, error) {
	var war War
	if err := s.client.trpc.Call(ctx, "war.getById", getWarInput{WarID: warID}, &war); err != nil {
		return nil, err
	}
	return &war, nil
}

// GetMany returns wars by ID, in the same order as warIDs, using one batched
// request.
func (s *WarService) GetMany(ctx context.Context, warIDs []string) ([]War, error) {
	return getManyByID[War](ctx, s.client, "wars", "war.getById", warIDs, func(id string) any {
		return getWarInput{WarID: id}
	})
}

// --- Party Service ---

// PartyService handles political parties. Party IDs appear on countries as
// Country.RulingPartyID and in party listings.
type PartyService struct {
	client *Client
}

type getPartyInput struct {
	PartyID string `json:"partyId"`
}

// GetByID returns a party by ID.
func (s *PartyService) GetByID(ctx context.Context, partyID string) (*Party, error) {
	var party Party
	if err := s.client.trpc.Call(ctx, "party.getById", getPartyInput{PartyID: partyID}, &party); err != nil {
		return nil, err
	}
	return &party, nil
}

// GetMany returns parties by ID, in the same order as partyIDs, using one
// batched request.
func (s *PartyService) GetMany(ctx context.Context, partyIDs []string) ([]Party, error) {
	return getManyByID[Party](ctx, s.client, "parties", "party.getById", partyIDs, func(id string) any {
		return getPartyInput{PartyID: id}
	})
}

// List returns one page of parties matching the filter, newest first.
func (s *PartyService) List(ctx context.Context, filter PartyFilter) (*CursorPage[Party], error) {
	var page CursorPage[Party]
	if err := s.client.trpc.Call(ctx, "party.getManyPaginated", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over every party matching the filter.
func (s *PartyService) All(ctx context.Context, filter PartyFilter) iter.Seq2[Party, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[Party], error) {
		f := filter
		f.Cursor = cursor
		return s.List(c, f)
	})
}

// GetInactiveMemberIDs returns the user IDs of a party's inactive members.
func (s *PartyService) GetInactiveMemberIDs(ctx context.Context, partyID string) ([]string, error) {
	var ids []string
	if err := s.client.trpc.Call(ctx, "party.getInactiveMemberIds", getPartyInput{PartyID: partyID}, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// PartyFilter specifies options for listing parties.
type PartyFilter struct {
	CountryID string `json:"countryId,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
	// Limit is the page size, at most 100.
	Limit    int    `json:"limit,omitempty"`
	RegionID string `json:"regionId,omitempty"`
	// Search matches party names.
	Search string `json:"search,omitempty"`
}

// --- Election Service ---

// ElectionService handles country and party elections.
type ElectionService struct {
	client *Client
}

type getElectionInput struct {
	ElectionID string `json:"electionId"`
}

// GetByID returns an election by ID.
func (s *ElectionService) GetByID(ctx context.Context, electionID string) (*Election, error) {
	var election Election
	if err := s.client.trpc.Call(ctx, "election.getElection", getElectionInput{ElectionID: electionID}, &election); err != nil {
		return nil, err
	}
	return &election, nil
}

// GetMany returns elections by ID, in the same order as electionIDs, using one
// batched request.
func (s *ElectionService) GetMany(ctx context.Context, electionIDs []string) ([]Election, error) {
	return getManyByID[Election](ctx, s.client, "elections", "election.getElection", electionIDs, func(id string) any {
		return getElectionInput{ElectionID: id}
	})
}

// List returns one page of a country's or a party's elections, newest first.
// The filter must set CountryID or PartyID.
func (s *ElectionService) List(ctx context.Context, filter ElectionFilter) (*CursorPage[Election], error) {
	var page CursorPage[Election]
	if err := s.client.trpc.Call(ctx, "election.getElections", filter, &page); err != nil {
		return nil, err
	}
	return &page, nil
}

// All returns an iterator over every election matching the filter. The filter
// must set CountryID or PartyID.
func (s *ElectionService) All(ctx context.Context, filter ElectionFilter) iter.Seq2[Election, error] {
	return NewCursorIterator(ctx, filter.Cursor, func(c context.Context, cursor string) (*CursorPage[Election], error) {
		f := filter
		f.Cursor = cursor
		return s.List(c, f)
	})
}

// ElectionFilter specifies options for listing elections. CountryID selects
// a country's presidential and congress elections, and PartyID a party's
// internal elections.
type ElectionFilter struct {
	CountryID string `json:"countryId,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
	// Limit is the page size, at most 100.
	Limit   int    `json:"limit,omitempty"`
	PartyID string `json:"partyId,omitempty"`
}

// --- Country and government additions ---

// GetDiplomacy returns a country's sworn enemy and defensive pacts with the
// bonuses they give.
func (s *CountryService) GetDiplomacy(ctx context.Context, countryID string) (*Diplomacy, error) {
	var diplomacy Diplomacy
	if err := s.client.trpc.Call(ctx, "countryDiplomacy.getByCountry", getCountryByIdInput{CountryID: countryID}, &diplomacy); err != nil {
		return nil, err
	}
	return &diplomacy, nil
}

// GetUnrest returns a country's unrest bar and revolution cooldown.
func (s *CountryService) GetUnrest(ctx context.Context, countryID string) (*UnrestStatus, error) {
	var unrest UnrestStatus
	if err := s.client.trpc.Call(ctx, "country.getUnrestData", getCountryByIdInput{CountryID: countryID}, &unrest); err != nil {
		return nil, err
	}
	return &unrest, nil
}

// GetStrategicResourceBonuses returns the strategic resources available to a
// country, through its alliance when it has one, and the bonuses they grant.
// It returns nil and no error when the country has none.
func (s *CountryService) GetStrategicResourceBonuses(ctx context.Context, countryID string) (*StrategicResourceBonuses, error) {
	var bonuses *StrategicResourceBonuses
	if err := s.client.trpc.Call(ctx, "country.getStrategicResourceBonuses", getCountryByIdInput{CountryID: countryID}, &bonuses); err != nil {
		return nil, err
	}
	return bonuses, nil
}

// GetCongressSeats returns a country's congress members and their parties.
func (s *GovernmentService) GetCongressSeats(ctx context.Context, countryID string) ([]CongressSeat, error) {
	var seats []CongressSeat
	if err := s.client.trpc.Call(ctx, "government.getCongressSeats", getGovernmentByCountryIdInput{CountryID: countryID}, &seats); err != nil {
		return nil, err
	}
	return seats, nil
}
