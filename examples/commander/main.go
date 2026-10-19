// commander interactively selects a military unit and shows its listed roster
// and orders for one active battle. Set WARERA_API_KEY when the API requires
// authentication, then run: go run ./examples/commander.
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
	fmt.Printf("\nMilitary unit: %s [%s], country=%s, region=%s\n", unit.Name, unit.ID, unit.CountryID, unit.Region)

	// The unit record contains member IDs; hydrate them in one bounded bulk
	// lookup so the example demonstrates the common ID-to-profile join.
	if len(unit.Members) == 0 {
		fmt.Println("The API returned no member IDs for this unit.")
	} else {
		members, err := client.Users.GetLiteMany(ctx, unit.Members)
		if err != nil {
			return fmt.Errorf("load unit member profiles: %w", err)
		}
		fmt.Printf("Listed roster: %d members\n", len(members))
		modeCounts := map[skillmix.Mode]int{}
		for _, member := range members {
			lastSeen := member.Dates.LastConnectionAt
			if lastSeen == "" {
				lastSeen = "not reported"
			}
			mix := skillmix.Analyze(member.Skills)
			modeCounts[mix.Mode()]++
			fmt.Printf("- %s | level=%d militaryRank=%d activeFlag=%t lastConnection=%s\n",
				member.Username, member.Level(), member.MilitaryRank, member.IsActive, lastSeen)
			fmt.Printf("  skill mix: %s leaning; war %.1f%% (%d pts), eco %.1f%% (%d pts); prestige levels war=%d eco=%d\n",
				mix.Mode(), mix.WarPercent, mix.WarPoints, mix.EcoPercent, mix.EcoPoints,
				mix.WarPrestigeLevels, mix.EcoPrestigeLevels)
			if mix.ClassifiedPoints != member.Leveling.SpentSkillPoints {
				fmt.Printf("  note: categorized regular points=%d, profile reports %d spent points\n",
					mix.ClassifiedPoints, member.Leveling.SpentSkillPoints)
			}
		}
		fmt.Printf("Roster by larger regular skill-point allocation: war=%d eco=%d balanced=%d unknown=%d\n",
			modeCounts[skillmix.ModeWar], modeCounts[skillmix.ModeEco],
			modeCounts[skillmix.ModeBalanced], modeCounts[skillmix.ModeUnknown])
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
	side := ""
	if battle.Attacker.CountryID == unit.CountryID {
		side = "attacker"
	} else if battle.Defender.CountryID == unit.CountryID {
		side = "defender"
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
