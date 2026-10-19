// soldier interactively chooses a player and a battle involving their country,
// then prints a read-only profile, loadout, and combat snapshot. Set
// WARERA_API_KEY when authentication is required, then run: go run ./examples/soldier.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/examples/internal/console"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	// API services share one client, so the same configuration and connection
	// resources are reused for profile, inventory, battle, and hit requests.
	client := warera.NewClient(warera.WithAPIKey(os.Getenv("WARERA_API_KEY")))
	defer client.Close()

	// Ctrl-C cancels both prompts and network calls. The overall deadline is a
	// guardrail for a multi-request exploration rather than a gameplay timer.
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(rootCtx, 2*time.Minute)
	defer cancel()
	prompt := console.New(os.Stdin, os.Stdout)

	player, err := console.SelectPlayer(ctx, client, prompt)
	if err != nil {
		return err
	}
	// Search returns a lightweight player record. Fetch the full profile for
	// combat skills and country, and query inventory separately for equipped items.
	captureStartedAt := time.Now().UTC()
	profile, err := client.Users.GetByID(ctx, player.ID)
	if err != nil {
		return fmt.Errorf("load full profile for %s: %w", player.Username, err)
	}
	equipment, err := client.Inventory.GetCurrentEquipment(ctx, player.ID)
	if err != nil {
		return fmt.Errorf("load current equipment: %w", err)
	}

	active := true
	// Start with active battles involving the player's country. This is only the
	// first page (up to 50) to keep the interactive choice list manageable.
	battlePage, err := client.Battles.List(ctx, warera.BattleFilter{
		CountryID: profile.CountryID,
		IsActive:  &active,
		Limit:     50,
	})
	if err != nil {
		return fmt.Errorf("list active battles for the player's country: %w", err)
	}
	battles := battlePage.Items
	if len(battles) == 0 {
		// Do not ask the developer to know or look up a battle ID. If no active
		// battle is listed, offer the latest page of this country's battle history.
		fmt.Println("No active battle was listed; loading recent battles instead.")
		recentPage, err := client.Battles.List(ctx, warera.BattleFilter{
			CountryID: profile.CountryID,
			Limit:     50,
		})
		if err != nil {
			return fmt.Errorf("list recent battles for the player's country: %w", err)
		}
		battles = recentPage.Items
		// The listing response is paginated by the API. Sort this page by its
		// creation timestamp so the most recently created records appear first.
		sort.SliceStable(battles, func(i, j int) bool {
			return battles[i].CreatedAt > battles[j].CreatedAt
		})
	}
	if len(battles) == 0 {
		return fmt.Errorf("no active or recent battles were found for country %s", profile.CountryID)
	}
	battleID, roundNumber, side, err := chooseBattle(ctx, prompt, profile.CountryID, battles)
	if err != nil {
		return err
	}
	// Live battle data is round-specific. The selected round supplies the round
	// record ID needed by the separate endpoint for recent hit details.
	live, err := client.Battles.GetLiveBattleData(ctx, battleID, roundNumber)
	if err != nil {
		return fmt.Errorf("load live battle data: %w", err)
	}
	if live.Round.RoundID == "" {
		return fmt.Errorf("battle %s has no round ID for round %d", battleID, roundNumber)
	}
	hits, err := client.Rounds.GetLastHits(ctx, live.Round.RoundID)
	if err != nil {
		return fmt.Errorf("load round hits: %w", err)
	}
	// A user's loot summary is useful for a compact per-battle contribution view.
	// Keep it optional: if this endpoint fails, the profile, loadout, and hit data
	// can still teach users how the other service calls fit together.
	contribution, contributionErr := client.LootSummaries.GetByBattleAndUser(ctx, battleID, player.ID)
	captureFinishedAt := time.Now().UTC()

	weapon := "none"
	if equipment.Weapon != nil {
		weapon = equipment.Weapon.Code
	}
	fmt.Printf("\nCapture window: %s to %s (UTC)\n", captureStartedAt.Format(time.RFC3339), captureFinishedAt.Format(time.RFC3339))
	fmt.Printf("Player: %s (level %d, military rank %d)\n", profile.Username, profile.Level(), profile.MilitaryRank)
	fmt.Printf("Combat skills: attack=%d criticalChance=%d criticalDamage=%d precision=%d armor=%d\n",
		profile.Skills.Attack.Total, profile.Skills.CriticalChance.Total,
		profile.Skills.CriticalDamages.Total, profile.Skills.Precision.Total, profile.Skills.Armor.Total)
	fmt.Printf("Current equipment: weapon=%s helmet=%s chest=%s gloves=%s pants=%s ammo=%s\n",
		weapon, equipment.Helmet.Code, equipment.Chest.Code, equipment.Gloves.Code, equipment.Pants.Code, equipment.Ammo)
	fmt.Printf("Battle %s, round %d (%s): active=%t nextTick=%s attackerDamage=%d defenderDamage=%d\n",
		battleID, roundNumber, side, live.Round.IsActive, live.Round.NextTickAt,
		live.Round.AttackerDamages, live.Round.DefenderDamages)
	if contributionErr != nil {
		fmt.Fprintf(os.Stderr, "Battle contribution summary unavailable: %v\n", contributionErr)
	} else {
		fmt.Printf("User battle summary: hits=%d damage=%d bounty=%.4f contract=%.4f\n",
			contribution.Hits, contribution.TotalDamage, contribution.TotalMoneyFromBounty, contribution.TotalMoneyFromContract)
	}

	var sideHits []warera.RoundLastHit
	// The hit endpoint separates attacker and defender hits. Use the side selected
	// above so the output is scoped to the player's chosen battle perspective.
	if side == "attacker" {
		sideHits = hits.Attacker
	} else {
		sideHits = hits.Defender
	}
	fmt.Printf("Recent hits returned for %s side\n", side)
	foundHit := false
	for _, hit := range sideHits {
		if hit.UserID != player.ID {
			continue
		}
		foundHit = true
		weaponCode := "unknown"
		if hit.Weapon != nil {
			weaponCode = hit.Weapon.Code
		}
		fmt.Printf("- %s damage=%d critical=%t missed=%t weapon=%s\n",
			hit.HitAt, hit.Damages, hit.IsCriticalHit, hit.IsMissed, weaponCode)
	}
	if !foundHit {
		fmt.Println("No hits for this player were present in the API's recent-hit response.")
	}
	fmt.Println("These are current API records and a per-battle summary, not a complete historical damage series.")
	fmt.Println("No build score is calculated; damage optimization needs explicit formula and cost assumptions.")
	return nil
}

func chooseBattle(
	ctx context.Context,
	prompt *console.Prompt,
	countryID string,
	battles []warera.Battle,
) (battleID string, roundNumber int, side string, err error) {
	options := make([]string, len(battles))
	// Show enough context to distinguish battles without requiring users to know
	// an ID in advance; include status and round when the list response has them.
	for i, battle := range battles {
		status := "completed"
		if battle.IsActive {
			status = "active"
		}
		currentRound := "round unknown"
		if battle.CurrentRound != nil && battle.CurrentRound.Number > 0 {
			currentRound = fmt.Sprintf("round %d", battle.CurrentRound.Number)
		}
		options[i] = fmt.Sprintf("%s | %s | attacker %s vs defender %s | %s, %s | created %s",
			battle.Type, battle.ID, battle.Attacker.CountryID, battle.Defender.CountryID,
			status, currentRound, battle.CreatedAt)
	}
	selected, err := prompt.Choose(ctx, "Choose a battle to inspect:", options)
	if err != nil {
		return "", 0, "", fmt.Errorf("select battle: %w", err)
	}
	battle := battles[selected]
	battleID = battle.ID
	if battle.CurrentRound != nil {
		roundNumber = battle.CurrentRound.Number
	}
	if roundNumber < 1 {
		roundNumber, err = askRoundNumber(ctx, prompt)
		if err != nil {
			return "", 0, "", err
		}
	}
	defaultSide := ""
	// Suggest the player's country's side only when it is unambiguous. The user
	// still chooses explicitly because recent-hit output is side-specific.
	if battle.Attacker.CountryID == countryID && battle.Defender.CountryID != countryID {
		defaultSide = "attacker"
	} else if battle.Defender.CountryID == countryID && battle.Attacker.CountryID != countryID {
		defaultSide = "defender"
	}
	side, err = chooseSide(ctx, prompt, defaultSide)
	return battleID, roundNumber, side, err
}

func askRoundNumber(ctx context.Context, prompt *console.Prompt) (int, error) {
	// Live battle lookup requires a positive round number; validate here so the
	// caller gets a useful retry instead of an opaque API error.
	for {
		value, err := prompt.Ask(ctx, "Enter the round number:")
		if err != nil {
			return 0, fmt.Errorf("read round number: %w", err)
		}
		roundNumber, err := strconv.Atoi(value)
		if err == nil && roundNumber > 0 {
			return roundNumber, nil
		}
		fmt.Println("Round number must be a positive integer.")
	}
}

func chooseSide(ctx context.Context, prompt *console.Prompt, defaultSide string) (string, error) {
	// The selected side controls which half of the recent-hit response is shown;
	// it does not change the battle query or perform an in-game action.
	options := []string{"Attacker", "Defender"}
	question := "Which battle side should be used for recent hits?"
	if defaultSide == "attacker" {
		options[0] += " (player country)"
		question += " The player's country is marked."
	} else if defaultSide == "defender" {
		options[1] += " (player country)"
		question += " The player's country is marked."
	}
	selected, err := prompt.Choose(ctx, question, options)
	if err != nil {
		return "", fmt.Errorf("select battle side: %w", err)
	}
	if selected == 0 {
		return "attacker", nil
	}
	return "defender", nil
}
