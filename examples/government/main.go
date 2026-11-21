// government interactively selects a country and collects a read-only briefing
// from several SDK services: government, alliance and diplomacy, wars, unrest,
// congress and parties, regions, citizens, military units, events, and battles.
// Set WARERA_API_KEY, which the politics procedures require, then run:
// go run ./examples/government. It scans all military units because the
// current listing filter does not include a country ID.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/examples/internal/console"
)

type countrySnapshot struct {
	startedAt  time.Time
	finishedAt time.Time
	country    *warera.Country
	government *warera.Government
	alliance   *warera.Alliance
	diplomacy  *warera.Diplomacy
	unrest     *warera.UnrestStatus
	seats      []warera.CongressSeat
	parties    []warera.Party
	wars       []warera.War
	names      map[string]string // country ID to name
	regions    []warera.Region
	users      []warera.UserListItem
	officials  map[string]string // user ID to username
	units      []warera.MilitaryUnit
	events     []warera.Event
	battles    []warera.Battle
}

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	// All briefing sections use this client, sharing configuration and HTTP
	// resources while remaining ordinary, independently callable SDK services.
	client := warera.NewClient(warera.WithAPIKey(os.Getenv("WARERA_API_KEY")))
	defer client.Close()

	// Country-wide pagination can take several requests. A longer timeout gives
	// the example room to finish while still allowing Ctrl-C to cancel promptly.
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(rootCtx, 5*time.Minute)
	defer cancel()
	prompt := console.New(os.Stdin, os.Stdout)
	country, err := console.SelectCountry(ctx, client, prompt)
	if err != nil {
		return err
	}
	countryID := country.ID
	fmt.Printf("\nCollecting a read-only briefing for %s (%s) [%s].\n\n", country.Name, country.Code, country.ID)

	snapshot := countrySnapshot{startedAt: time.Now().UTC()}
	// Independent sections are safe to fetch concurrently and make the briefing
	// quicker. Bound the number of active requests so this stays polite to the API.
	const maxConcurrentSections = 3
	sem := make(chan struct{}, maxConcurrentSections)
	var wg sync.WaitGroup
	var errorsMu sync.Mutex
	var sectionErrors []error
	launch := func(name string, fetch func(context.Context) error) {
		// Keep sections independent: one endpoint failure should not prevent the
		// other records from being collected and printed as a partial briefing.
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				errorsMu.Lock()
				sectionErrors = append(sectionErrors, fmt.Errorf("%s: %w", name, ctx.Err()))
				errorsMu.Unlock()
				return
			}
			defer func() { <-sem }()
			if err := fetch(ctx); err != nil {
				errorsMu.Lock()
				sectionErrors = append(sectionErrors, fmt.Errorf("%s: %w", name, err))
				errorsMu.Unlock()
			}
		}()
	}

	launch("country", func(ctx context.Context) error {
		// Re-fetch the selected record by ID so this section demonstrates the
		// normal lookup pattern used when an ID came from a search or menu.
		country, err := client.Countries.GetByID(ctx, countryID)
		snapshot.country = country
		if err != nil || country.AllianceID == "" {
			return err
		}
		// The country record holds only the alliance ID; load the alliance for
		// its name and member countries.
		snapshot.alliance, err = client.Alliances.GetByID(ctx, country.AllianceID)
		return err
	})
	launch("country names", func(ctx context.Context) error {
		// Alliances, pacts, and wars refer to other countries by ID. One
		// world listing resolves them all to names.
		countries, err := client.Countries.GetAll(ctx)
		snapshot.names = map[string]string{}
		for _, c := range countries {
			snapshot.names[c.ID] = c.Name
		}
		return err
	})
	launch("diplomacy", func(ctx context.Context) error {
		diplomacy, err := client.Countries.GetDiplomacy(ctx, countryID)
		snapshot.diplomacy = diplomacy
		return err
	})
	launch("unrest", func(ctx context.Context) error {
		unrest, err := client.Countries.GetUnrest(ctx, countryID)
		snapshot.unrest = unrest
		return err
	})
	launch("congress", func(ctx context.Context) error {
		// Seats name each member's party by ID; the party listing below
		// supplies the names.
		seats, err := client.Government.GetCongressSeats(ctx, countryID)
		snapshot.seats = seats
		return err
	})
	launch("parties", func(ctx context.Context) error {
		parties, err := warera.Collect(client.Parties.All(ctx, warera.PartyFilter{CountryID: countryID, Limit: 100}), 0)
		sort.Slice(parties, func(i, j int) bool { return len(parties[i].MemberIDs) > len(parties[j].MemberIDs) })
		snapshot.parties = parties
		return err
	})
	launch("government", func(ctx context.Context) error {
		// Government is a separate record keyed by country, so this call enriches
		// the country record with office holders and congress membership.
		government, err := client.Government.GetByCountryID(ctx, countryID)
		snapshot.government = government
		if err != nil {
			return err
		}
		// The government record holds user IDs; load compact profiles in one
		// batched request so office holders can be shown by name.
		ids := []string{}
		for _, id := range []string{government.PresidentID, government.VicePresidentID} {
			if id != "" {
				ids = append(ids, id)
			}
		}
		for _, minister := range government.Ministers {
			ids = append(ids, minister.UserID)
		}
		profiles, err := client.Users.GetLiteMany(ctx, ids)
		snapshot.officials = map[string]string{}
		for _, profile := range profiles {
			snapshot.officials[profile.ID] = profile.Username
		}
		return err
	})
	launch("regions", func(ctx context.Context) error {
		// Regions have a whole-world listing, so filter locally after loading them.
		allRegions, err := client.Regions.GetAll(ctx)
		if err != nil {
			return err
		}
		for _, region := range allRegions {
			if region.CountryID == countryID {
				snapshot.regions = append(snapshot.regions, region)
			}
		}
		sort.Slice(snapshot.regions, func(i, j int) bool { return snapshot.regions[i].Name < snapshot.regions[j].Name })
		return nil
	})
	launch("users", func(ctx context.Context) error {
		// All consumes the iterator across pages; entries carry only user IDs.
		// Collect keeps partial results on an error.
		users, err := warera.Collect(client.Users.All(ctx, warera.UserFilter{CountryID: countryID, Limit: 100}), 0)
		snapshot.users = users
		return err
	})
	launch("military units", func(ctx context.Context) error {
		// The endpoint has no country filter, so this reads the paginated world
		// listing and retains only units currently associated with this country.
		// Stream the pages so unrelated units do not accumulate in memory.
		for unit, err := range client.MilitaryUnits.All(ctx, warera.MilitaryUnitFilter{Limit: 100}) {
			if err != nil {
				return err
			}
			if unit.CountryID == countryID {
				snapshot.units = append(snapshot.units, unit)
			}
		}
		sort.Slice(snapshot.units, func(i, j int) bool { return snapshot.units[i].Name < snapshot.units[j].Name })
		return nil
	})
	launch("events", func(ctx context.Context) error {
		// Iterator-based endpoints hide cursor handling, so callers can focus on
		// the filter and the records they want to inspect.
		events, err := warera.Collect(client.Events.All(ctx, warera.EventFilter{CountryID: countryID, Limit: 100}), 0)
		snapshot.events = events
		return err
	})
	active := true
	launch("active battles", func(ctx context.Context) error {
		// A pointer bool lets the filter distinguish "active only" from an unset
		// IsActive filter, which would return both active and completed battles.
		battles, err := warera.Collect(client.Battles.All(ctx, warera.BattleFilter{
			CountryID: countryID,
			IsActive:  &active,
			Limit:     100,
		}), 0)
		snapshot.battles = battles
		if err != nil {
			return err
		}
		// Battles belong to wars. Collect the distinct war IDs and load the
		// wars in one batched request for their running score.
		seen := map[string]bool{}
		warIDs := []string{}
		for _, battle := range battles {
			if battle.WarID != "" && !seen[battle.WarID] {
				seen[battle.WarID] = true
				warIDs = append(warIDs, battle.WarID)
			}
		}
		snapshot.wars, err = client.Wars.GetMany(ctx, warIDs)
		return err
	})

	wg.Wait()
	// Record the end after every launched section has completed, including
	// sections that failed; the printed window describes the whole collection.
	snapshot.finishedAt = time.Now().UTC()
	printSnapshot(snapshot, countryID)

	if len(sectionErrors) > 0 {
		sort.Slice(sectionErrors, func(i, j int) bool { return sectionErrors[i].Error() < sectionErrors[j].Error() })
		for _, err := range sectionErrors {
			fmt.Fprintf(os.Stderr, "partial section failure: %v\n", err)
		}
		return fmt.Errorf("briefing is partial: %d section(s) failed", len(sectionErrors))
	}
	return nil
}

func printSnapshot(snapshot countrySnapshot, countryID string) {
	fmt.Printf("Country briefing: %s\n", countryID)
	fmt.Printf("Capture window: %s to %s (UTC)\n", snapshot.startedAt.Format(time.RFC3339), snapshot.finishedAt.Format(time.RFC3339))
	fmt.Println("Sections are fetched concurrently and may reflect slightly different server times.")

	if country := snapshot.country; country != nil {
		fmt.Printf("Country: %s (%s), government type=%s, reported population=%d, wars=%d, allies=%d\n",
			country.Name, country.Code, country.GovType, country.CurrentPopulation, len(country.WarsWithIDs), len(country.AllyIDs))
		if income := country.Income; income != nil {
			fmt.Printf("Daily income estimate: total=%.2f (development=%.2f, employee income tax=%.2f, manager tax=%.2f, strategic resources=%.2f)\n",
				income.Total, income.Development, income.EmployeeIncomeTax, income.ManagerIncomeTax, income.StrategicResourcesBonus)
		}
	}
	if government := snapshot.government; government != nil {
		fmt.Printf("Government: president=%s, vice-president=%s, congress members=%d\n",
			resolveUser(government.PresidentID, snapshot.officials), resolveUser(government.VicePresidentID, snapshot.officials), len(government.CongressMemberIDs))
		for _, minister := range government.Ministers {
			title := minister.Title
			if minister.CustomTitle != "" {
				title += " (" + minister.CustomTitle + ")"
			}
			fmt.Printf("  minister: %s, %s, rights=%v\n", resolveUser(minister.UserID, snapshot.officials), title, minister.Rights)
		}
	}

	if alliance := snapshot.alliance; alliance != nil {
		members := make([]string, 0, len(alliance.MemberCountries))
		for _, member := range alliance.MemberCountries {
			members = append(members, countryName(member.CountryID, snapshot.names))
		}
		fmt.Printf("Alliance: %s, %d member countries: %s\n", alliance.Name, len(members), strings.Join(members, ", "))
	} else if snapshot.country != nil {
		fmt.Println("Alliance: none")
	}
	if diplomacy := snapshot.diplomacy; diplomacy != nil {
		if enemy := diplomacy.SwornEnemy; enemy != nil {
			fmt.Printf("Sworn enemy: %s (+%.0f%% damage, %.0f dealt)\n", countryName(enemy.EnemyID, snapshot.names), enemy.BonusPercent, enemy.DamagesDealt)
		}
		for _, pact := range diplomacy.DefensivePacts {
			fmt.Printf("  defensive pact: %s (+%.0f%% damage, %.0f dealt)\n", countryName(pact.PartnerID, snapshot.names), pact.BonusPercent, pact.DamagesDealt)
		}
	}
	for _, war := range snapshot.wars {
		fmt.Printf("War: %s vs %s, battles won %d-%d, rounds won %d-%d\n",
			countryName(war.Attacker.CountryID, snapshot.names), countryName(war.Defender.CountryID, snapshot.names),
			war.Attacker.WonBattlesCount, war.Defender.WonBattlesCount, war.Attacker.WonRoundsCount, war.Defender.WonRoundsCount)
	}
	if unrest := snapshot.unrest; unrest != nil {
		fmt.Printf("Unrest: %.0f/%.0f, revolution cooldown=%t\n", unrest.Bar, unrest.BarMax, unrest.IsOnRevolutionCooldown)
	}
	if len(snapshot.seats) > 0 || len(snapshot.parties) > 0 {
		seatsByParty := map[string]int{}
		for _, seat := range snapshot.seats {
			seatsByParty[seat.PartyID]++
		}
		fmt.Printf("Congress: %d seats; parties: %d\n", len(snapshot.seats), len(snapshot.parties))
		for i, party := range snapshot.parties {
			if i == 5 {
				break
			}
			fmt.Printf("  party: %s, members=%d, congress seats=%d", party.Name, len(party.MemberIDs), seatsByParty[party.ID])
			if party.FavoriteGovType != "" {
				fmt.Printf(", favourite government=%s", party.FavoriteGovType)
			}
			fmt.Println()
		}
		if n := seatsByParty[""]; n > 0 {
			fmt.Printf("  independents: %d congress seats\n", n)
		}
	}

	regionPopulation := 0
	activeRegions := 0
	// This is a sum of reported regional values. It is shown alongside the
	// country's own reported population for exploration, not as a consistency check.
	for _, region := range snapshot.regions {
		regionPopulation += region.CurrentPopulation
		if region.ActiveBattle != nil && (region.ActiveBattle.BattleID != "" || region.ActiveBattle.Details != nil) {
			activeRegions++
		}
	}
	fmt.Printf("Regions: %d (reported current population sum=%d, with active battle=%d)\n",
		len(snapshot.regions), regionPopulation, activeRegions)
	fmt.Printf("Citizens: %d; military units: %d; events: %d; active battles: %d\n",
		len(snapshot.users), len(snapshot.units), len(snapshot.events), len(snapshot.battles))

	sort.Slice(snapshot.events, func(i, j int) bool { return snapshot.events[i].CreatedAt.After(snapshot.events[j].CreatedAt.Time) })
	// Keep terminal output compact while still showing representative records.
	for i, event := range snapshot.events {
		if i == 5 {
			break
		}
		fmt.Printf("Recent event: %s type=%s at=%s\n", event.ID, event.Data.Type, event.CreatedAt)
	}
	for i, battle := range snapshot.battles {
		if i == 5 {
			break
		}
		fmt.Printf("Active battle: %s type=%s attacker=%s defender=%s\n",
			battle.ID, battle.Type, battle.Attacker.CountryID, battle.Defender.CountryID)
	}
	fmt.Println("Population growth and change alerts require comparing this capture with a later saved snapshot.")
}

func countryName(countryID string, names map[string]string) string {
	if name, ok := names[countryID]; ok {
		return name
	}
	return countryID
}

func resolveUser(userID string, usernames map[string]string) string {
	// Show office holders by name when their profile loaded, keeping the ID
	// as a fallback.
	if userID == "" {
		return "(none)"
	}
	if name, ok := usernames[userID]; ok {
		return fmt.Sprintf("%s [%s]", name, userID)
	}
	return userID
}
