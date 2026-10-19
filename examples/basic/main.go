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
		// Concurrent calls that arrive during this small window can share a
		// transport batch; calls made later remain separate requests.
		warera.WithBatchWindow(25*time.Millisecond),
	)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Fetch game server dates and rollover clock. These timestamps come from
	// the game API, so clients can use them instead of assuming a local timezone.
	dates, err := client.GameConfig.GetDates(ctx)
	if err != nil {
		log.Fatalf("failed to fetch game dates: %v", err)
	}
	fmt.Printf("WarEra Server Next Day: %s, Next Energy Regen: %s\n", dates.NextDayAt, dates.NextRegenAt)

	// 2. Fetch all countries. GetAll hides the endpoint's pagination and returns
	// the complete list when successful.
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

	// 3. Search for a user by a sample name. Search returns IDs and lightweight
	// matches; use a user lookup to retrieve a player record for a chosen ID.
	results, err := client.Search.SearchAnything(ctx, "Kautilya")
	if err != nil {
		log.Fatalf("failed to search: %v", err)
	}
	fmt.Printf("Search Results for 'Kautilya': %d user IDs found (%v)\n", len(results.UserIDs), results.UserIDs)

	// Fetch a profile for the first match to demonstrate the ID-to-record flow.
	if len(results.UserIDs) > 0 {
		user, err := client.Users.GetLiteByID(ctx, results.UserIDs[0])
		if err != nil {
			log.Fatalf("failed to get user profile: %v", err)
		}
		fmt.Printf("Player Found: %s | Level: %d | Avatar: %s\n",
			user.Username, user.Level(), user.AvatarURL)
	}
}
