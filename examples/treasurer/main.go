// treasurer interactively selects a country and prints its current treasury,
// tax fields, strategic-resource bonuses, state army spending, donations, and a
// bounded page of country transactions. Set WARERA_API_KEY, which the
// strategic-resource, state army, and donation procedures require, then run:
// go run ./examples/treasurer.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
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
	client := warera.NewClient(warera.WithAPIKey(os.Getenv("WARERA_API_KEY")))
	defer client.Close()

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(rootCtx, 2*time.Minute)
	defer cancel()
	prompt := console.New(os.Stdin, os.Stdout)
	countryChoice, err := console.SelectCountry(ctx, client, prompt)
	if err != nil {
		return err
	}

	capturedAt := time.Now().UTC()
	country, err := client.Countries.GetByID(ctx, countryChoice.ID)
	if err != nil {
		return fmt.Errorf("load country details: %w", err)
	}
	transactions, err := client.Transactions.List(ctx, warera.TransactionFilter{
		CountryID: country.ID,
		Limit:     25,
	})
	if err != nil {
		return fmt.Errorf("load recent country transactions: %w", err)
	}

	fmt.Printf("\nCountry economy snapshot: %s (%s) [%s]\n", country.Name, country.Code, country.ID)
	fmt.Printf("Captured at: %s (UTC)\n", capturedAt.Format(time.RFC3339))
	fmt.Printf("Reported population=%d government=%s treasuryMoney=%.4f taxes(income=%.4f market=%.4f manager=%.4f)\n",
		country.CurrentPopulation, country.GovType, country.Money,
		country.Taxes.Income, country.Taxes.Market, country.Taxes.Manager)
	// The income breakdown is the game's own daily estimate by source.
	if income := country.Income; income != nil {
		fmt.Printf("Daily income estimate: total=%.4f development=%.4f employeeIncomeTax=%.4f managerIncomeTax=%.4f strategicResources=%.4f governmentType=%.4f\n",
			income.Total, income.Development, income.EmployeeIncomeTax, income.ManagerIncomeTax,
			income.StrategicResourcesBonus, income.GovType)
	} else {
		fmt.Println("No income breakdown was present in the country response.")
	}

	// Strategic resources are pooled across an alliance, so the bonuses can
	// come from regions held by other member countries.
	bonuses, err := client.Countries.GetStrategicResourceBonuses(ctx, country.ID)
	if err != nil {
		return fmt.Errorf("load strategic resource bonuses: %w", err)
	}
	if bonuses == nil {
		fmt.Println("\nStrategic resources: none")
	} else {
		resources := make([]string, 0, len(bonuses.RegionIDs))
		for resource, regions := range bonuses.RegionIDs {
			resources = append(resources, fmt.Sprintf("%s=%d", resource, len(regions)))
		}
		sort.Strings(resources)
		fmt.Printf("\nStrategic resources (regions per resource, shared by %d countries): %s\n", len(bonuses.CountryIDs), strings.Join(resources, " "))
		skills := make([]string, 0, len(bonuses.Bonuses.Skills))
		for skill, bonus := range bonuses.Bonuses.Skills {
			skills = append(skills, fmt.Sprintf("%s=%+g", skill, bonus))
		}
		sort.Strings(skills)
		fmt.Printf("Bonuses: development=%+.1f%% production=%+.1f%% skills: %s\n",
			bonuses.Bonuses.DevelopmentPercent, bonuses.Bonuses.ProductionPercent, strings.Join(skills, " "))
	}

	// The state army is the country's spending on military units: oil for
	// headquarters, mercenary contracts, and items given under sponsorship laws.
	const days = 7
	army, err := client.StateArmy.GetOverview(ctx, country.ID, days)
	if err != nil {
		return fmt.Errorf("load state army overview: %w", err)
	}
	fmt.Printf("\nState army spending over %d days:", days)
	if len(army.ByKind) == 0 {
		fmt.Print(" none")
	}
	fmt.Println()
	for _, total := range army.ByKind {
		fmt.Printf("- %s money=%.2f items=%d\n", total.Kind, total.Money, total.Quantity)
	}

	// Donations are cumulative per donor, so the first page lists the largest
	// donors. Their user IDs resolve to names in one batched lookup.
	donated, err := client.Donations.GetTotal(ctx, warera.DonationTotalFilter{CountryID: country.ID})
	if err != nil {
		return fmt.Errorf("load donation total: %w", err)
	}
	donors, err := client.Donations.List(ctx, warera.DonationFilter{CountryID: country.ID, Limit: 5})
	if err != nil {
		return fmt.Errorf("load top donors: %w", err)
	}
	fmt.Printf("\nDonations: %.2f from %d donors\n", donated.TotalAmount, donated.DonorCount)
	donorIDs := make([]string, len(donors.Items))
	for i, donor := range donors.Items {
		donorIDs[i] = donor.UserID
	}
	profiles, err := client.Users.GetLiteMany(ctx, donorIDs)
	if err != nil {
		return fmt.Errorf("load donor profiles: %w", err)
	}
	for _, donor := range donors.Items {
		fmt.Printf("- %s donated %.2f, latest %s\n", profiles[donor.UserID].Username, donor.Amount, donor.UpdatedAt)
	}

	fmt.Printf("\nTransaction page: %d records", len(transactions.Items))
	if transactions.NextCursor != "" {
		fmt.Print(" (more pages available)")
	}
	fmt.Println()
	for _, tx := range transactions.Items {
		itemCode := tx.ItemCode
		if itemCode == "" && tx.Item != nil {
			itemCode = tx.Item.Code
		}
		fmt.Printf("- %s type=%s item=%s quantity=%d money=%.4f buyer=%s seller=%s sellerCountry=%s sellerMU=%s\n",
			tx.CreatedAt, tx.TransactionType, itemCode, tx.Quantity, tx.Money,
			tx.BuyerID, tx.SellerID, tx.SellerCountryID, tx.SellerMUID)
	}

	fmt.Println("These are current country fields and one transaction page. The sample is not a complete ledger or a forecast of budget sustainability.")
	return nil
}
