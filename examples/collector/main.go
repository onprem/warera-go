// collector is a long-running, read-only data collector: every interval it
// appends market prices, new market trades, and new world events to JSON Lines
// files, remembering what it has seen so a restart picks up where it stopped.
// It shows the SDK features such tools rely on: Until for fetching only new
// records, Directory for resolving IDs to names, and hooks for logging rate
// limits. Set WARERA_API_KEY, then run:
//
//	go run ./examples/collector -out ./warera-data -interval 15m
//
// Use -once to collect a single round and exit.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"syscall"
	"time"

	"github.com/onprem/warera-go"
)

// maxPerRound caps how many new trades or events one round reads, so a first
// run or a restart after a long gap does not crawl the whole history.
const maxPerRound = 1000

// state is what the collector has already stored, saved between runs.
type state struct {
	LastTradeID string `json:"lastTradeId"`
	LastEventID string `json:"lastEventId"`
}

type collector struct {
	client *warera.Client
	dir    *warera.Directory
	out    string
	state  state
}

func main() {
	out := flag.String("out", "warera-data", "directory for the JSON Lines files and saved state")
	interval := flag.Duration("interval", 15*time.Minute, "time between collection rounds")
	once := flag.Bool("once", false, "collect one round and exit")
	flag.Parse()
	if err := run(*out, *interval, *once); err != nil {
		slog.Error("collector stopped", "err", err)
		os.Exit(1)
	}
}

func run(out string, interval time.Duration, once bool) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client := warera.NewClient(
		warera.WithAPIKey(os.Getenv("WARERA_API_KEY")),
		warera.WithUserAgent("warera-go-collector-example"),
		// A long-running tool should say when it is being slowed down.
		warera.WithHooks(warera.Hooks{OnWait: func(w warera.WaitInfo) {
			if w.Reason != "retry" || w.Duration > time.Second {
				slog.Info("waiting", "reason", w.Reason, "duration", w.Duration.Round(time.Millisecond), "procedures", w.Procedures)
			}
		}}),
	)
	defer client.Close()

	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	c := &collector{client: client, dir: warera.NewDirectory(client, time.Hour), out: out}
	if err := c.loadState(); err != nil {
		return err
	}

	for {
		if err := c.collect(ctx); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			// Keep running through temporary failures; stop on anything else,
			// such as a rejected API key.
			if !warera.IsTemporary(err) {
				return err
			}
			slog.Warn("round failed; will retry next round", "err", err)
		}
		if once {
			return nil
		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(interval):
		}
	}
}

func (c *collector) collect(ctx context.Context) error {
	now := time.Now().UTC()
	prices, err := c.client.Trading.GetPrices(ctx)
	if err != nil {
		return fmt.Errorf("prices: %w", err)
	}
	if err := c.appendLines("prices.jsonl", []any{map[string]any{"at": now, "prices": prices}}); err != nil {
		return err
	}

	// Listings come newest first, so Until stops at the newest record stored
	// last time and fetches only what is new.
	trades, err := warera.Collect(warera.Until(
		c.client.Transactions.All(ctx, warera.TransactionFilter{
			TransactionTypes: []warera.TransactionType{warera.TransactionTypeTrading},
			Limit:            100,
		}),
		func(tx warera.TransactionListItem) bool { return tx.ID == c.state.LastTradeID },
	), maxPerRound)
	if err != nil {
		return fmt.Errorf("trades: %w", err)
	}
	events, err := warera.Collect(warera.Until(
		c.client.Events.All(ctx, warera.EventFilter{Limit: 100}),
		func(event warera.Event) bool { return event.ID == c.state.LastEventID },
	), maxPerRound)
	if err != nil {
		return fmt.Errorf("events: %w", err)
	}

	// Store records oldest first so the files read chronologically.
	slices.Reverse(trades)
	slices.Reverse(events)
	if err := c.appendLines("trades.jsonl", toAny(trades)); err != nil {
		return err
	}
	if err := c.appendLines("events.jsonl", toAny(events)); err != nil {
		return err
	}
	if len(trades) > 0 {
		c.state.LastTradeID = trades[len(trades)-1].ID
	}
	if len(events) > 0 {
		c.state.LastEventID = events[len(events)-1].ID
	}
	if err := c.saveState(); err != nil {
		return err
	}

	slog.Info("collected", "items priced", len(prices), "new trades", len(trades), "new events", len(events))
	for _, event := range events {
		c.logEvent(ctx, event)
	}
	return nil
}

// logEvent prints an event with its countries named, using the Directory so
// the country list is fetched once an hour rather than per event.
func (c *collector) logEvent(ctx context.Context, event warera.Event) {
	names := make([]string, 0, len(event.CountryIDs))
	for _, id := range event.CountryIDs {
		if country, err := c.dir.Country(ctx, id); err == nil {
			names = append(names, country.Name)
		} else {
			names = append(names, id)
		}
	}
	slog.Info("event", "type", event.Data.Type, "countries", names, "at", event.CreatedAt)
}

func (c *collector) appendLines(name string, records []any) error {
	if len(records) == 0 {
		return nil
	}
	file, err := os.OpenFile(filepath.Join(c.out, name), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	for _, record := range records {
		if err := encoder.Encode(record); err != nil {
			file.Close()
			return err
		}
	}
	return file.Close()
}

func (c *collector) loadState() error {
	data, err := os.ReadFile(filepath.Join(c.out, "state.json"))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.state)
}

func (c *collector) saveState() error {
	data, err := json.MarshalIndent(c.state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(c.out, "state.json"), data, 0o644)
}

func toAny[T any](records []T) []any {
	out := make([]any, len(records))
	for i, record := range records {
		out[i] = record
	}
	return out
}
