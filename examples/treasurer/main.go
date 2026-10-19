// treasurer interactively selects a country and prints its current treasury,
// tax fields, strategic-resource bonuses, and a bounded page of country
// transactions. Set WARERA_API_KEY when required, then run:
// go run ./examples/treasurer.
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
	fmt.Printf("Reported population=%d treasuryMoney=%.4f taxes(income=%.4f market=%.4f selfWork=%.4f)\n",
		country.CurrentPopulation, country.Money,
		country.Taxes.Income, country.Taxes.Market, country.Taxes.SelfWork)
	if resources := country.StrategicResources; resources != nil {
		fmt.Printf("Strategic-resource bonuses: development=%.4f production=%.4f\n",
			resources.Bonuses.DevelopmentPercent, resources.Bonuses.ProductionPercent)
	} else {
		fmt.Println("No strategic-resource data was present in the country response.")
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
