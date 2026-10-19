package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/onprem/warera-go"
)

func main() {
	client := warera.NewClient(
		warera.WithAPIKey(os.Getenv("WARERA_API_KEY")),
		// Calls submitted close together have a chance to share a batch request.
		// The window trades a tiny amount of latency for request coalescing.
		warera.WithBatchWindow(30*time.Millisecond),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("=== Example 1: Transparent Request Coalescing ===")
	// Start independent reads concurrently. The client can combine calls that
	// reach its batcher during the configured window; batching is transparent to
	// each service method, which still returns its own typed result.
	var wg sync.WaitGroup
	var countries []warera.Country
	var dates *warera.GameDates
	var prices warera.ItemPrices

	wg.Add(3)
	go func() {
		defer wg.Done()
		var err error
		countries, err = client.Countries.GetAll(ctx)
		if err != nil {
			log.Printf("error fetching countries: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		dates, err = client.GameConfig.GetDates(ctx)
		if err != nil {
			log.Printf("error fetching dates: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		prices, err = client.Trading.GetPrices(ctx)
		if err != nil {
			log.Printf("error fetching prices: %v", err)
		}
	}()

	wg.Wait()
	fmt.Printf("Coalesced response received: %d countries, dates=%v, %d prices\n\n",
		len(countries), dates != nil, len(prices))

	fmt.Println("=== Example 2: Go 1.23+ Range Iterator Auto-Pagination ===")
	// All follows cursor pages and yields one value at a time. Breaking early
	// stops iteration, which is useful when a caller only needs a sample.
	count := 0
	for article, err := range client.Articles.All(ctx, warera.ArticleFilter{Type: "last", Limit: 10}) {
		if err != nil {
			log.Fatalf("iterator error: %v", err)
		}
		fmt.Printf("[%d] %s (by %s)\n", count+1, article.Title, article.Author)
		count++
		if count >= 5 { // Stop after five records even if more pages are available.
			break
		}
	}
}
