package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/onprem/warera-go"
)

func main() {
	apiKey := os.Getenv("WARERA_API_KEY")

	// Create a new WarEra client
	client := warera.NewClient(
		warera.WithAPIKey(apiKey),
		warera.WithBatchWindow(25*time.Millisecond), // automatically coalesces concurrent requests
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Fetch game server dates and rollover clock
	dates, err := client.GameConfig.GetDates(ctx)
	if err != nil {
		log.Fatalf("failed to fetch game dates: %v", err)
	}
	fmt.Printf("WarEra Server Next Day: %s, Next Energy Regen: %s\n", dates.NextDayAt, dates.NextRegenAt)

	// 2. Fetch all countries
	countries, err := client.Countries.GetAll(ctx)
	if err != nil {
		log.Fatalf("failed to fetch countries: %v", err)
	}
	fmt.Printf("Loaded %d countries\n", len(countries))
	if len(countries) > 0 {
		top := countries[0]
		fmt.Printf("Featured Country: %s (Code: %s, Population: %d, Money: %.2f)\n",
			top.Name, top.Code, top.CurrentPopulation, top.Money)
	}

	// 3. Search the universe
	results, err := client.Search.SearchAnything(ctx, "Dog")
	if err != nil {
		log.Fatalf("failed to search: %v", err)
	}
	fmt.Printf("Search Results for 'Dog': %d players found\n", len(results.Users))
}
