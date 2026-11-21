// trader prints a read-only snapshot of a player's companies and current
// market context. It prompts for a player handle. Set WARERA_API_KEY
// when the API requires authentication, then run: go run ./examples/trader.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
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
	// The SDK client owns the HTTP and batching resources used by its services.
	// Close it when this short-lived example finishes.
	client := warera.NewClient(warera.WithAPIKey(os.Getenv("WARERA_API_KEY")))
	defer client.Close()

	// Keep both interactive input and API requests interruptible, and cap this
	// exploration so a stalled endpoint cannot leave the example running forever.
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(rootCtx, 2*time.Minute)
	defer cancel()
	prompt := console.New(os.Stdin, os.Stdout)
	player, err := console.SelectPlayer(ctx, client, prompt)
	if err != nil {
		return err
	}
	userID := player.ID
	fmt.Printf("\nLoading a read-only company snapshot for %s [%s].\n\n", player.Username, player.ID)
	// These services are queried one after another. The timestamps below make it
	// clear that the resulting records are a collection window, not one atomic snapshot.
	captureStartedAt := time.Now().UTC()

	// All streams through the SDK's paginator and yield one company ID at a time.
	// This keeps listing separate from hydration: collect the IDs, then fetch the
	// richer company records in a bounded bulk request.
	companyIDs := make([]string, 0)
	for companyID, err := range client.Companies.All(ctx, warera.CompanyFilter{UserID: userID}) {
		if err != nil {
			return fmt.Errorf("list company IDs: %w", err)
		}
		companyIDs = append(companyIDs, companyID)
	}
	if len(companyIDs) == 0 {
		fmt.Println("No companies found for this user.")
		return nil
	}

	companies, err := client.Companies.GetMany(ctx, companyIDs)
	if err != nil {
		return fmt.Errorf("load company details: %w", err)
	}
	// Prices are fetched once as a lookup map, so each company can be paired with
	// the current reference price for its product code.
	prices, err := client.Trading.GetPrices(ctx)
	if err != nil {
		return fmt.Errorf("load reference prices: %w", err)
	}
	captureFinishedAt := time.Now().UTC()

	fmt.Printf("Capture window: %s to %s (UTC); records may change during collection.\n",
		captureStartedAt.Format(time.RFC3339), captureFinishedAt.Format(time.RFC3339))
	fmt.Printf("Companies for user %s (%d)\n", userID, len(companies))
	// The price map provides current reference prices by item code; it is not a
	// realized sale price and does not include this company's costs or taxes.
	for _, company := range companies {
		fmt.Printf("- %s [%s] product=%s production=%g region=%s\n",
			company.Name, company.ID, company.ItemCode, company.Production, company.Region)
		if price, ok := prices[company.ItemCode]; ok {
			fmt.Printf("  reference price for %s: %.4f\n", company.ItemCode, price)
		}
	}

	itemCode := strings.TrimSpace(os.Getenv("WARERA_ITEM_CODE"))
	if itemCode == "" {
		itemCode = companies[0].ItemCode
	}
	if itemCode != "" {
		// Top orders provide a small view of visible buy/sell demand for one item.
		// Choose WARERA_ITEM_CODE to inspect another product without editing code.
		orders, err := client.Trading.GetTopOrders(ctx, itemCode, 5)
		if err != nil {
			return fmt.Errorf("load top orders for %q: %w", itemCode, err)
		}
		fmt.Printf("\nTop visible orders for %s\n", itemCode)
		printOrders("buy", orders.BuyOrders)
		printOrders("sell", orders.SellOrders)
	}

	fmt.Println("\nAll values above are API records. This example does not estimate profit or rank products.")
	return nil
}

func printOrders(label string, orders []warera.TopOrder) {
	// Both sides use the same order shape; the label tells the reader whether the
	// displayed price is a buy offer or a sell offer.
	for _, order := range orders {
		fmt.Printf("- %s price=%.4f quantity=%d offerAt=%s\n", label, order.Price, order.Quantity, order.OfferAt)
	}
}
