// researcher selects a country, reads a bounded page of country events, and
// follows its pinned-article reference when present. Set WARERA_API_KEY when
// required, then run: go run ./examples/researcher.
package main

import (
	"context"
	"encoding/json"
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

	startedAt := time.Now().UTC()
	country, err := client.Countries.GetByID(ctx, countryChoice.ID)
	if err != nil {
		return fmt.Errorf("load country details: %w", err)
	}
	events, err := client.Events.List(ctx, warera.EventFilter{CountryID: country.ID, Limit: 25})
	if err != nil {
		return fmt.Errorf("load country events: %w", err)
	}

	var pinnedArticle *warera.Article
	var articleErr error
	if country.PinnedArticle != "" {
		pinnedArticle, articleErr = client.Articles.GetByID(ctx, country.PinnedArticle)
	}
	finishedAt := time.Now().UTC()

	fmt.Printf("\nResearch snapshot: %s (%s) [%s]\n", country.Name, country.Code, country.ID)
	fmt.Printf("Capture window: %s to %s (UTC); event data is one bounded page.\n",
		startedAt.Format(time.RFC3339), finishedAt.Format(time.RFC3339))
	fmt.Printf("Country event records: %d", len(events.Items))
	if events.NextCursor != "" {
		fmt.Print(" (more pages available)")
	}
	fmt.Println()
	for _, event := range events.Items {
		data, err := json.Marshal(event.Data)
		if err != nil {
			return fmt.Errorf("format event %s data: %w", event.ID, err)
		}
		fmt.Printf("- event=%s at=%s priority=%s data=%s\n", event.ID, event.CreatedAt, event.Priority, data)
	}

	if country.PinnedArticle == "" {
		fmt.Println("No pinned article ID was present on this country record.")
	} else if articleErr != nil {
		fmt.Fprintf(os.Stderr, "Pinned article %s could not be loaded: %v\n", country.PinnedArticle, articleErr)
	} else {
		fmt.Printf("\nPinned article: %s | category=%s language=%s author=%s published=%s slug=%s\n",
			pinnedArticle.Title, pinnedArticle.Category, pinnedArticle.Language,
			pinnedArticle.Author, pinnedArticle.PublishedAt, pinnedArticle.Slug)
		content := strings.TrimSpace(pinnedArticle.Content)
		if content == "" {
			fmt.Println("Article body is empty in the returned record.")
		} else if strings.Contains(content, "<") {
			fmt.Printf("Article body is available in Article.Content (%d bytes; content may contain HTML).\n", len(content))
		} else {
			fmt.Printf("Article excerpt: %s\n", truncate(content, 400))
		}
	}
	fmt.Println("Events and article content are source records. This example does not infer causality or maintain an archive of changes.")
	return nil
}

func truncate(value string, maxRunes int) string {
	runes := []rune(value)
	if len(runes) <= maxRunes {
		return value
	}
	return string(runes[:maxRunes]) + "…"
}
