// trader prints a read-only snapshot of a player's companies, their NPC
// employees and managers, and market context for their products: recent price
// history, order books, and the money a production point earns. It prompts for
// a player handle. Set WARERA_API_KEY, which the employee procedures require, then run:
// go run ./examples/trader.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
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
	// Companies now run on NPC employees. Each roster also carries the game's
	// own estimate of hourly output and costs, and a company can have one
	// player manager. These calls are issued concurrently and the client
	// batches them into a few HTTP requests.
	rosters := make([]*warera.CompanyRoster, len(companies))
	managers := make([]*warera.Manager, len(companies))
	errs := make([]error, 2*len(companies))
	var wg sync.WaitGroup
	for i, company := range companies {
		wg.Add(2)
		go func() {
			defer wg.Done()
			rosters[i], errs[2*i] = client.Employees.GetCompanyRoster(ctx, company.ID)
		}()
		go func() {
			defer wg.Done()
			managers[i], errs[2*i+1] = client.Managers.GetByCompanyID(ctx, company.ID)
		}()
	}
	wg.Wait()
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("load employees and managers: %w", err)
	}
	captureFinishedAt := time.Now().UTC()

	fmt.Printf("Capture window: %s to %s (UTC); records may change during collection.\n",
		captureStartedAt.Format(time.RFC3339), captureFinishedAt.Format(time.RFC3339))
	fmt.Printf("Companies for user %s (%d)\n", userID, len(companies))
	// The price map provides current reference prices by item code; it is not a
	// realized sale price and does not include this company's costs or taxes.
	for i, company := range companies {
		fmt.Printf("- %s [%s] product=%s production=%.2f region=%s\n",
			company.Name, company.ID, company.ItemCode, company.Production, company.RegionID)
		if price, ok := prices[company.ItemCode]; ok {
			fmt.Printf("  reference price for %s: %.4f\n", company.ItemCode, price)
		}
		roster, hourly := rosters[i], rosters[i].Hourly
		fmt.Printf("  employees %d/%d; hourly estimate: production=%.2f wages=%.4f manager=%.4f\n",
			len(roster.Employees), hourly.Cap, hourly.Production, hourly.WageCost, hourly.ManagerCost)
		if enabled := hourly.IfEnabled; enabled != nil {
			fmt.Printf("  company is disabled; once enabled: production=%.2f wages=%.4f manager=%.4f\n",
				enabled.Production, enabled.WageCost, enabled.ManagerCost)
		}
		for _, employee := range roster.Employees {
			fmt.Printf("    %s (age %d) production=%g luck=%g wage=%.4f happiness=%.1f bonus=%.1f%%\n",
				employee.Name, employee.Age,
				employee.Base.Production+employee.Trained.Production, employee.Base.Luck+employee.Trained.Luck,
				employee.Wage, employee.HappinessBreakdown.Value, employee.ProductionBonus.Total)
		}
		if manager := managers[i]; manager != nil {
			fmt.Printf("  manager %s: management level %d, +%g%% production, %.4f/hour\n",
				manager.UserID, manager.ManagementLevel, manager.ProductionBonusPercent, manager.HourlyWage)
		}
	}

	// Market context covers every product these companies make, or the
	// comma-separated item codes in WARERA_ITEM_CODES.
	itemCodes := []string{}
	if env := strings.TrimSpace(os.Getenv("WARERA_ITEM_CODES")); env != "" {
		itemCodes = strings.Split(env, ",")
	} else {
		seen := map[string]bool{}
		for _, company := range companies {
			if company.ItemCode != "" && !seen[company.ItemCode] {
				seen[company.ItemCode] = true
				itemCodes = append(itemCodes, company.ItemCode)
			}
		}
	}
	if err := printMarket(ctx, client, itemCodes); err != nil {
		return err
	}

	fmt.Println("\nAll values above are API records and the game's own hourly estimates. This example does not estimate profit or rank products.")
	return nil
}

func printMarket(ctx context.Context, client *warera.Client, itemCodes []string) error {
	// One request returns the top of the order book for every item.
	books, err := client.Trading.GetTopOrdersPerItem(ctx, itemCodes, 3)
	if err != nil {
		return fmt.Errorf("load order books: %w", err)
	}
	// The game's own estimate of what a production point earns when spent on
	// each item at current prices.
	ppValues, err := client.MarketStats.GetPPMoneyValueByItem(ctx)
	if err != nil {
		return fmt.Errorf("load production point values: %w", err)
	}
	// Trading history is one call per item; issued together, they share a
	// batched request.
	histories := make([]*warera.ItemTrading, len(itemCodes))
	errs := make([]error, len(itemCodes))
	var wg sync.WaitGroup
	for i, itemCode := range itemCodes {
		wg.Add(1)
		go func() {
			defer wg.Done()
			histories[i], errs[i] = client.Trading.GetItemTrading(ctx, itemCode)
		}()
	}
	wg.Wait()
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("load trading history: %w", err)
	}

	for i, itemCode := range itemCodes {
		fmt.Printf("\nMarket for %s: money per production point=%.4f\n", itemCode, ppValues[itemCode])
		// The last day in the history is today and still incomplete, so
		// compare the last complete day with the one a week before it.
		if days := histories[i].Values; len(days) >= 9 {
			last, weekAgo := days[len(days)-2], days[len(days)-9]
			fmt.Printf("  average price %s: %.4f (%+.1f%% vs %s), volume=%d in %d trades\n",
				last.ValueAt, last.AvgValue, 100*(last.AvgValue/weekAgo.AvgValue-1), weekAgo.ValueAt,
				last.TotalQuantity, last.TransactionsCount)
		}
		book := books[itemCode]
		printOrders("buy", book.BuyOrders)
		printOrders("sell", book.SellOrders)
	}
	return nil
}

func printOrders(label string, orders []warera.TopOrder) {
	// Both sides use the same order shape; the label tells the reader whether the
	// displayed price is a buy offer or a sell offer.
	for _, order := range orders {
		fmt.Printf("  %s price=%.4f quantity=%d offerAt=%s\n", label, order.Price, order.Quantity, order.OfferAt)
	}
}
