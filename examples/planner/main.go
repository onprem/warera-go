// planner is a read-only production planner built on the analytics packages.
// It ranks what a production point earns on each item at current prices,
// shows what buying a batch of each item's inputs would cost on the order
// book, and, with an API key, values the key owner's workforce per hour.
// Run:
//
//	WARERA_API_KEY=... go run ./examples/planner
//
// Without a key it prints only the market ranking.
package main

import (
	"context"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"time"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/company"
	"github.com/onprem/warera-go/analytics/market"
	"github.com/onprem/warera-go/analytics/production"
)

// batchUnits is how many units of each item the input costs are quoted for.
const batchUnits = 100

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	key := os.Getenv("WARERA_API_KEY")
	client := warera.NewClient(warera.WithAPIKey(key))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// The analytics packages compute from records; they make no requests.
	cfg, err := client.GameConfig.Get(ctx)
	if err != nil {
		return fmt.Errorf("load game config: %w", err)
	}
	prices, err := client.Trading.GetPrices(ctx)
	if err != nil {
		return fmt.Errorf("load prices: %w", err)
	}
	values := production.Values(cfg, prices)

	// One request returns the order books of every input the ranked items
	// need, so the input quotes below reflect what is on sale now.
	var inputs []string
	for _, v := range values {
		for input := range cfg.Items[v.ItemCode].ProductionNeeds {
			if !slices.Contains(inputs, input) {
				inputs = append(inputs, input)
			}
		}
	}
	books, err := client.Trading.GetTopOrdersPerItem(ctx, inputs, 10)
	if err != nil {
		return fmt.Errorf("load order books: %w", err)
	}

	fmt.Println("Money per production point at current prices, best first.")
	fmt.Println("'buying inputs' assumes market inputs; 'from raw' makes every step in-house.")
	fmt.Printf("\n%-11s %10s %14s %10s   %s\n", "item", "price", "buying inputs", "from raw", fmt.Sprintf("inputs for %d units on the order book", batchUnits))
	for _, v := range values {
		fmt.Printf("%-11s %10.4f %14.4f %10.4f   %s\n", v.ItemCode, v.Price, v.MoneyPerPoint, v.MoneyPerPointFromRaw, inputQuote(cfg, books, v.ItemCode))
	}

	if key == "" {
		fmt.Println("\nSet WARERA_API_KEY to also value your own workforce.")
		return nil
	}
	summary, err := client.Me.GetWorkforceSummary(ctx)
	if err != nil {
		return fmt.Errorf("load workforce: %w", err)
	}
	profit, err := company.ProfitPerHour(cfg, prices, summary)
	if err != nil {
		return err
	}
	fmt.Printf("\nYour workforce: %d employees, %d managers, %.2f production points per hour.\n",
		summary.EmployeeCount, summary.ManagerCount, summary.Production)
	for _, code := range slices.Sorted(maps.Keys(profit.ByItem)) {
		out := profit.ByItem[code]
		fmt.Printf("  %-10s %7.2f points -> %7.2f units, worth %8.2f (inputs %.2f)\n", code, out.ProductionPoints, out.Units, out.Revenue, out.InputCost)
	}
	fmt.Printf("Per hour: revenue %.2f - inputs %.2f - wages %.2f - managers %.2f = profit %.2f\n",
		profit.Revenue, profit.InputCost, profit.Wages, profit.ManagerCost, profit.Profit)
	fmt.Println("Inputs are valued at market prices even when your own companies make them.")
	return nil
}

// inputQuote describes what buying the inputs for batchUnits units costs on
// the order book, or "-" for items without inputs.
func inputQuote(cfg *warera.GameConfig, books map[string]warera.TopOrders, itemCode string) string {
	needs := cfg.Items[itemCode].ProductionNeeds
	if len(needs) == 0 {
		return "-"
	}
	total, complete := 0.0, true
	for input, quantity := range needs {
		book := books[input]
		fill := market.BuyCost(&book, int(quantity*batchUnits))
		total += fill.Total
		complete = complete && fill.Complete
	}
	if !complete {
		return fmt.Sprintf("%.2f (the listed orders do not cover it)", total)
	}
	return fmt.Sprintf("%.2f", total)
}
