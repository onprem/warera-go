package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var defaultFixtures = []struct {
	method string
	input  string
}{
	{"gameConfig.getDates", ""},
	{"search.searchAnything", `{"searchText":"Kautilya"}`},
	{"country.getCountryById", `{"countryId":"6813b6d446e731854c7ac7a6"}`},
	{"user.getUserLite", `{"userId":"6aa7ad06d91c79361d0b0673"}`},
	{"itemTrading.getPrices", ""},
}

func main() {
	methodFlag := flag.String("method", "", "tRPC method to record (e.g. country.getCountryById)")
	inputFlag := flag.String("input", "", "JSON input for the method")
	allFlag := flag.Bool("all", false, "Record all standard public test fixtures")
	outDirFlag := flag.String("out", "spec/fixtures", "Output directory for fixtures")
	flag.Parse()

	apiKey := os.Getenv("WARERA_API_KEY")

	if *allFlag {
		for _, f := range defaultFixtures {
			if err := recordEndpoint(f.method, f.input, apiKey, *outDirFlag); err != nil {
				fmt.Printf("Error recording %s: %v\n", f.method, err)
			}
		}
		return
	}

	if *methodFlag == "" {
		fmt.Println("Usage: go run ./cmd/record -method <method> [-input '<json>'] [-out <dir>]")
		fmt.Println("   or: go run ./cmd/record -all")
		os.Exit(1)
	}

	if err := recordEndpoint(*methodFlag, *inputFlag, apiKey, *outDirFlag); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func recordEndpoint(method, input, apiKey, outDir string) error {
	baseURL := "https://api2.warera.io/trpc"
	endpointURL := fmt.Sprintf("%s/%s", baseURL, method)

	if input != "" {
		endpointURL = fmt.Sprintf("%s?input=%s", endpointURL, url.QueryEscape(input))
	} else if method != "gameStat.getWorldDevelopment" {
		endpointURL = fmt.Sprintf("%s?input=%s", endpointURL, url.QueryEscape("{}"))
	}

	req, err := http.NewRequest("GET", endpointURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if apiKey != "" {
		req.Header.Set("x-api-key", apiKey)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var envelope struct {
		Result struct {
			Data json.RawMessage `json:"data"`
		} `json:"result"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("failed to decode tRPC envelope: %w", err)
	}

	if envelope.Error != nil {
		return fmt.Errorf("tRPC server returned error: %s", envelope.Error.Message)
	}

	// Indent the data JSON for clean reading
	var formatted any
	if err := json.Unmarshal(envelope.Result.Data, &formatted); err != nil {
		return fmt.Errorf("failed to re-parse data: %w", err)
	}

	indented, err := json.MarshalIndent(formatted, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	outPath := filepath.Join(outDir, method+".json")
	if err := os.WriteFile(outPath, indented, 0644); err != nil {
		return fmt.Errorf("failed to write fixture file: %w", err)
	}

	fmt.Printf("Recorded fixture: %s (%d bytes)\n", outPath, len(indented))
	return nil
}
