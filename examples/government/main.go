// government interactively selects a country and collects a read-only briefing
// from several SDK services. Set WARERA_API_KEY when authentication is required,
// then run: go run ./examples/government. It scans all military units because
// the current listing filter does not include a country ID.
package main

import (
	"context"
	"fmt"
	"iter"
	"log"
	"os"
	"os/signal"
	"sort"
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
		for _, id := range []string{government.President, government.VicePresident} {
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
		// collect keeps partial results on an error.
		users, err := collect(client.Users.All(ctx, warera.UserFilter{CountryID: countryID, Limit: 100}))
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
		events, err := collect(client.Events.All(ctx, warera.EventFilter{CountryID: countryID, Limit: 100}))
		snapshot.events = events
		return err
	})
	active := true
	launch("active battles", func(ctx context.Context) error {
		// A pointer bool lets the filter distinguish "active only" from an unset
		// IsActive filter, which would return both active and completed battles.
		battles, err := collect(client.Battles.All(ctx, warera.BattleFilter{
			CountryID: countryID,
			IsActive:  &active,
			Limit:     100,
		}))
		snapshot.battles = battles
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

func collect[T any](items iter.Seq2[T, error]) ([]T, error) {
	// Iterator errors arrive alongside values. Return what was already collected
	// with the error so a partial briefing can still show useful sections.
	collected := make([]T, 0)
	for item, err := range items {
		if err != nil {
			return collected, err
		}
		collected = append(collected, item)
	}
	return collected, nil
}

func printSnapshot(snapshot countrySnapshot, countryID string) {
	fmt.Printf("Country briefing: %s\n", countryID)
	fmt.Printf("Capture window: %s to %s (UTC)\n", snapshot.startedAt.Format(time.RFC3339), snapshot.finishedAt.Format(time.RFC3339))
	fmt.Println("Sections are fetched concurrently and may reflect slightly different server times.")

	if country := snapshot.country; country != nil {
		fmt.Printf("Country: %s (%s), government type=%s, reported population=%d, wars=%d, allies=%d\n",
			country.Name, country.Code, country.GovType, country.CurrentPopulation, len(country.WarsWith), len(country.Allies))
		if income := country.Income; income != nil {
			fmt.Printf("Daily income estimate: total=%.2f (development=%.2f, employee income tax=%.2f, manager tax=%.2f, strategic resources=%.2f)\n",
				income.Total, income.Development, income.EmployeeIncomeTax, income.ManagerIncomeTax, income.StrategicResourcesBonus)
		}
	}
	if government := snapshot.government; government != nil {
		fmt.Printf("Government: president=%s, vice-president=%s, congress members=%d\n",
			resolveUser(government.President, snapshot.officials), resolveUser(government.VicePresident, snapshot.officials), len(government.CongressMembers))
		for _, minister := range government.Ministers {
			title := minister.Title
			if minister.CustomTitle != "" {
				title += " (" + minister.CustomTitle + ")"
			}
			fmt.Printf("  minister: %s, %s, rights=%v\n", resolveUser(minister.UserID, snapshot.officials), title, minister.Rights)
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

	sort.Slice(snapshot.events, func(i, j int) bool { return snapshot.events[i].CreatedAt > snapshot.events[j].CreatedAt })
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
