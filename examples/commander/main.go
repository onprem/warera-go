// commander interactively selects a military unit and shows its roster with
// each member's damage, the unit's league standing and sponsorship, and its
// orders for one active battle. Set WARERA_API_KEY, which the member, league,
// and sponsorship procedures require, then run: go run ./examples/commander.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/examples/internal/console"
	"github.com/onprem/warera-go/examples/internal/skillmix"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	client := warera.NewClient(warera.WithAPIKey(os.Getenv("WARERA_API_KEY")))
	defer client.Close()

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(rootCtx, 3*time.Minute)
	defer cancel()
	prompt := console.New(os.Stdin, os.Stdout)

	unit, err := console.SelectMilitaryUnit(ctx, client, prompt)
	if err != nil {
		return err
	}
	capturedAt := time.Now().UTC()
	fmt.Printf("\nMilitary unit: %s [%s], country=%s, region=%s\n", unit.Name, unit.ID, unit.CountryID, unit.RegionID)

	// The unit record contains member IDs; hydrate them in one bounded bulk
	// lookup so the example demonstrates the common ID-to-profile join.
	if len(unit.MemberIDs) == 0 {
		fmt.Println("The API returned no member IDs for this unit.")
	} else {
		members, err := client.Users.GetLiteMany(ctx, unit.MemberIDs)
		if err != nil {
			return fmt.Errorf("load unit member profiles: %w", err)
		}
		// Compact profiles no longer include skills, so load them separately in
		// one batched request. Results come back in the same order as the IDs.
		skills, err := client.Users.GetSkillsMany(ctx, unit.MemberIDs)
		if err != nil {
			return fmt.Errorf("load unit member skills: %w", err)
		}
		// Member records carry the unit's own weekly, monthly, and all-time
		// damage and help counts for each member, keyed by user ID.
		memberStats, err := client.MilitaryUnits.GetMembers(ctx, unit.ID)
		if err != nil {
			return fmt.Errorf("load unit member stats: %w", err)
		}
		statsByUser := map[string]warera.MUMember{}
		for _, stats := range memberStats {
			statsByUser[stats.UserID] = stats
		}
		fmt.Printf("Listed roster: %d members\n", len(members))
		modeCounts := map[skillmix.Mode]int{}
		for i, member := range members {
			lastSeen := "not reported"
			if at := member.Dates.LastConnectionAt; !at.IsZero() {
				lastSeen = at.Format(time.RFC3339)
			}
			mix := skillmix.Analyze(skills[i].Skills)
			modeCounts[mix.Mode()]++
			fmt.Printf("- %s | level=%d militaryRank=%d activeFlag=%t lastConnection=%s\n",
				member.Username, member.Level(), member.MilitaryRank, member.IsActive, lastSeen)
			fmt.Printf("  skill mix: %s leaning; war %.1f%% (%d pts), eco %.1f%% (%d pts); prestige levels war=%d eco=%d\n",
				mix.Mode(), mix.WarPercent, mix.WarPoints, mix.EcoPercent, mix.EcoPoints,
				mix.WarPrestigeLevels, mix.EcoPrestigeLevels)
			if stats, ok := statsByUser[member.ID]; ok {
				fmt.Printf("  damage: week=%.0f month=%.0f total=%.0f; help: week=%d month=%d total=%d\n",
					stats.WeeklyDamagesCount, stats.MonthlyDamagesCount, stats.TotalDamagesCount,
					stats.WeeklyHelpCount, stats.MonthlyHelpCount, stats.TotalHelpCount)
			}
			if mix.ClassifiedPoints != member.Leveling.SpentSkillPoints {
				fmt.Printf("  note: categorized regular points=%d, profile reports %d spent points\n",
					mix.ClassifiedPoints, member.Leveling.SpentSkillPoints)
			}
		}
		fmt.Printf("Roster by larger regular skill-point allocation: war=%d eco=%d balanced=%d unknown=%d\n",
			modeCounts[skillmix.ModeWar], modeCounts[skillmix.ModeEco],
			modeCounts[skillmix.ModeBalanced], modeCounts[skillmix.ModeUnknown])
	}

	// The MU league ranks units weekly within divisions; League and Me are nil
	// when the unit has not joined this week.
	league, err := client.MilitaryUnits.GetLeague(ctx, unit.ID)
	if err != nil {
		return fmt.Errorf("load unit league: %w", err)
	}
	if me := league.Me; me != nil {
		fmt.Printf("\nLeague week %s: division %d, rank %d of %d, points=%.0f, next week=%s (ends %s)\n",
			league.YearWeek, me.Division, me.Rank, len(league.Rows), me.Points, me.Move, league.WeekEndsAt)
	} else {
		fmt.Printf("\nLeague week %s: the unit is not in a league.\n", league.YearWeek)
	}
	// A country can sponsor a unit through a law. GetSponsorship returns nil
	// when the unit has no active or offered sponsorship.
	sponsorship, err := client.MilitaryUnits.GetSponsorship(ctx, unit.ID)
	if err != nil {
		return fmt.Errorf("load unit sponsorship: %w", err)
	}
	if sponsorship != nil {
		fmt.Printf("Sponsorship: %s by country %s since %s\n", sponsorship.Status, sponsorship.CountryID, sponsorship.AcceptedAt)
	} else {
		fmt.Println("Sponsorship: none")
	}

	active := true
	battles, err := client.Battles.List(ctx, warera.BattleFilter{
		CountryID: unit.CountryID,
		IsActive:  &active,
		Limit:     30,
	})
	if err != nil {
		return fmt.Errorf("list active battles for the unit's country: %w", err)
	}
	if len(battles.Items) == 0 {
		fmt.Println("No active battle was listed for this country; no battle orders to inspect.")
		fmt.Printf("Roster capture time: %s (UTC)\n", capturedAt.Format(time.RFC3339))
		return nil
	}

	options := make([]string, len(battles.Items))
	for i, battle := range battles.Items {
		currentRound := "round unknown"
		if battle.CurrentRound != nil && battle.CurrentRound.Number > 0 {
			currentRound = fmt.Sprintf("round %d", battle.CurrentRound.Number)
		}
		options[i] = fmt.Sprintf("%s | %s vs %s | %s | ID %s",
			battle.Type, battle.Attacker.CountryID, battle.Defender.CountryID, currentRound, battle.ID)
	}
	selected, err := prompt.Choose(ctx, "Choose an active battle to inspect:", options)
	if err != nil {
		return fmt.Errorf("select active battle: %w", err)
	}
	battle := battles.Items[selected]
	var side warera.BattleSide
	if battle.Attacker.CountryID == unit.CountryID {
		side = warera.BattleSideAttacker
	} else if battle.Defender.CountryID == unit.CountryID {
		side = warera.BattleSideDefender
	} else {
		return fmt.Errorf("selected battle does not include unit country %s", unit.CountryID)
	}

	orders, err := client.BattleOrders.GetByBattle(ctx, battle.ID, side)
	if err != nil {
		return fmt.Errorf("load %s-side battle orders: %w", side, err)
	}
	fmt.Printf("\nOrders for %s side of battle %s (%s)\n", side, battle.ID, battle.Type)
	shown := 0
	for _, order := range orders {
		// Country orders apply to all units from that country. Unit-specific
		// orders are included only for the selected unit.
		isUnitOrder := order.MUID == unit.ID
		isCountryOrder := order.MUID == "" && order.CountryID == unit.CountryID
		if !isUnitOrder && !isCountryOrder {
			continue
		}
		target := "country"
		if isUnitOrder {
			target = "unit"
		}
		fmt.Printf("- target=%s active=%t priority=%s rank=%d: %s\n",
			target, order.IsActive, order.Priority, order.Rank, order.Text)
		shown++
	}
	if shown == 0 {
		fmt.Println("No country-level or unit-specific orders were returned for this side.")
	}
	fmt.Printf("\nCapture time: %s (UTC)\n", time.Now().UTC().Format(time.RFC3339))
	fmt.Println("The roster and active flag are API fields; they do not establish a member's combat readiness or deployment availability.")
	return nil
}
