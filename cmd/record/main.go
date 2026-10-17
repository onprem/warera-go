package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// OpenAPISpec represents the minimal OpenAPI paths structure for endpoint discovery.
type OpenAPISpec struct {
	Paths map[string]struct {
		Post struct {
			OperationID string `json:"operationId"`
			Description string `json:"description"`
			RequestBody *struct {
				Content struct {
					JSON struct {
						Schema struct {
							Type       string                     `json:"type"`
							Properties map[string]OpenAPIProperty `json:"properties"`
							Required   []string                   `json:"required"`
						} `json:"schema"`
					} `json:"application/json"`
				} `json:"content"`
			} `json:"requestBody"`
		} `json:"post"`
	} `json:"paths"`
}

type OpenAPIProperty struct {
	Type        string           `json:"type"`
	Description string           `json:"description"`
	Enum        []string         `json:"enum"`
	Items       *OpenAPIProperty `json:"items"`
}

func main() {
	methodFlag := flag.String("method", "", "Targeted tRPC method to record (e.g. country.getCountryById). Overwrites local fixture by default.")
	inputFlag := flag.String("input", "", "Optional custom JSON input for the targeted method")
	forceFlag := flag.Bool("force", false, "Re-record all fixtures cleanly, overwriting existing files on disk")
	outDirFlag := flag.String("out", "spec/fixtures", "Output directory for fixture files")
	specFlag := flag.String("spec", "spec/openapi.json", "Path to WarEra OpenAPI specification")
	flag.Parse()

	apiKey := os.Getenv("WARERA_API_KEY")

	// Ensure output directory exists
	if err := os.MkdirAll(*outDirFlag, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// 1. Load OpenAPI spec to infer available endpoints and request schemas
	spec, err := loadOpenAPISpec(*specFlag)
	if err != nil {
		fmt.Printf("Error loading OpenAPI spec from %s: %v\n", *specFlag, err)
		os.Exit(1)
	}

	// 2. Discover seed parameters from disk and fallback defaults
	seeds := loadSeedValues(*outDirFlag)

	// Mode A: Targeted single-endpoint recording (replaces local fixture by default)
	if *methodFlag != "" {
		method := *methodFlag
		input := *inputFlag
		if input == "" {
			// Auto-construct input using OpenAPI property schema and seeds
			props := getEndpointProperties(spec, method)
			input = buildInputPayload(props, seeds)
		}

		fmt.Printf("Recording targeted endpoint: %s (input: %s)\n", method, input)
		size, err := recordEndpoint(method, input, apiKey, *outDirFlag)
		if err != nil {
			fmt.Printf("Error recording %s: %v\n", method, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully recorded %s (%d bytes) to %s\n", method, size, filepath.Join(*outDirFlag, method+".json"))
		return
	}

	// Mode B: Bulk discovery mode across all OpenAPI endpoints
	// By default: records only missing fixtures. Use --force to re-record existing ones.
	endpoints := discoverEndpoints(spec)
	fmt.Printf("Discovered %d endpoints from %s\n", len(endpoints), *specFlag)
	if *forceFlag {
		fmt.Println("Running in --force mode: all fixtures will be re-recorded cleanly.")
	} else {
		fmt.Println("Default mode: skipping fixtures that already exist on disk (use --force to overwrite).")
	}
	fmt.Println()

	var recordedCount, skippedCount, failedCount int

	for _, ep := range endpoints {
		fixturePath := filepath.Join(*outDirFlag, ep.method+".json")
		exists := fileExists(fixturePath)

		if exists && !*forceFlag {
			skippedCount++
			fmt.Printf("  [Skip] %s (fixture already exists)\n", ep.method)
			continue
		}

		input := buildInputPayload(ep.properties, seeds)
		if hasUnresolvedRequiredProps(ep.properties, ep.required, seeds) {
			failedCount++
			fmt.Printf("  [Skip] %s (missing required seed parameters)\n", ep.method)
			continue
		}

		size, err := recordEndpoint(ep.method, input, apiKey, *outDirFlag)
		if err != nil {
			failedCount++
			fmt.Printf("  [Fail] %s: %v\n", ep.method, err)
			continue
		}

		recordedCount++
		action := "Recorded"
		if exists && *forceFlag {
			action = "Re-recorded"
		}
		fmt.Printf("  [%s] %s (%d bytes)\n", action, ep.method, size)
	}

	fmt.Println()
	fmt.Printf("Summary: %d recorded, %d skipped, %d failed/unresolved.\n", recordedCount, skippedCount, failedCount)
	if failedCount > 0 {
		os.Exit(1)
	}
}

type endpointInfo struct {
	method     string
	properties map[string]OpenAPIProperty
	required   []string
}

func discoverEndpoints(spec *OpenAPISpec) []endpointInfo {
	var list []endpointInfo
	for pathKey, pathItem := range spec.Paths {
		method := strings.TrimPrefix(pathKey, "/")
		var props map[string]OpenAPIProperty
		var req []string
		if rb := pathItem.Post.RequestBody; rb != nil {
			props = rb.Content.JSON.Schema.Properties
			req = rb.Content.JSON.Schema.Required
		}
		list = append(list, endpointInfo{
			method:     method,
			properties: props,
			required:   req,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].method < list[j].method
	})
	return list
}

func getEndpointProperties(spec *OpenAPISpec, method string) map[string]OpenAPIProperty {
	pathKey := "/" + method
	if item, ok := spec.Paths[pathKey]; ok {
		if rb := item.Post.RequestBody; rb != nil {
			return rb.Content.JSON.Schema.Properties
		}
	}
	return nil
}

func loadOpenAPISpec(specPath string) (*OpenAPISpec, error) {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return nil, err
	}
	var spec OpenAPISpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

// loadSeedValues returns sample domain parameters, augmenting defaults with local fixtures.
func loadSeedValues(fixturesDir string) map[string]any {
	seeds := map[string]any{
		"countryId":        "6813b6d446e731854c7ac7a6",
		"userId":           "6aa7ad06d91c79361d0b0673",
		"muId":             "6aa9311d8af2bd246aff6189",
		"battleId":         "6ad3b89d7b158f464a343664",
		"roundId":          "6ad3b89d7b158f464a343666",
		"warId":            "6ad3b89d7b158f464a343661",
		"regionId":         "6813b7059403bc4170a5d79b",
		"defenderRegionId": "6813b7059403bc4170a5d79b",
		"articleId":        "6ad3b816bf1676fa11bff564",
		"itemCode":         "ammo",
		"searchText":       "Kautilya",
		"limit":            5,
		"perPage":          10,
		"type":             "last",
		"rankingType":      "countryDamages",
	}

	// Refine IDs from existing local fixtures if available
	if data, err := os.ReadFile(filepath.Join(fixturesDir, "country.getCountryById.json")); err == nil {
		var c struct {
			ID string `json:"_id"`
		}
		if json.Unmarshal(data, &c) == nil && c.ID != "" {
			seeds["countryId"] = c.ID
		}
	}
	if data, err := os.ReadFile(filepath.Join(fixturesDir, "user.getUserLite.json")); err == nil {
		var u struct {
			ID string `json:"_id"`
			MU string `json:"mu"`
		}
		if json.Unmarshal(data, &u) == nil {
			if u.ID != "" {
				seeds["userId"] = u.ID
			}
			if u.MU != "" {
				seeds["muId"] = u.MU
			}
		}
	}
	if data, err := os.ReadFile(filepath.Join(fixturesDir, "battle.getBattles.json")); err == nil {
		var b struct {
			Items []struct {
				ID           string `json:"_id"`
				War          string `json:"war"`
				CurrentRound struct {
					ID string `json:"_id"`
				} `json:"currentRound"`
				Defender struct {
					Region string `json:"region"`
				} `json:"defender"`
			} `json:"items"`
		}
		if json.Unmarshal(data, &b) == nil && len(b.Items) > 0 {
			if b.Items[0].ID != "" {
				seeds["battleId"] = b.Items[0].ID
			}
			if b.Items[0].War != "" {
				seeds["warId"] = b.Items[0].War
			}
			if b.Items[0].CurrentRound.ID != "" {
				seeds["roundId"] = b.Items[0].CurrentRound.ID
			}
			if b.Items[0].Defender.Region != "" {
				seeds["regionId"] = b.Items[0].Defender.Region
				seeds["defenderRegionId"] = b.Items[0].Defender.Region
			}
		}
	}
	if data, err := os.ReadFile(filepath.Join(fixturesDir, "article.getArticlesPaginated.json")); err == nil {
		var a struct {
			Items []struct {
				ID string `json:"_id"`
			} `json:"items"`
		}
		if json.Unmarshal(data, &a) == nil && len(a.Items) > 0 && a.Items[0].ID != "" {
			seeds["articleId"] = a.Items[0].ID
		}
	}

	return seeds
}

func buildInputPayload(props map[string]OpenAPIProperty, seeds map[string]any) string {
	if len(props) == 0 {
		return ""
	}
	payload := make(map[string]any)
	for propName, prop := range props {
		if propName == "cursor" {
			continue // Omit cursor to get first page
		}
		if val, ok := seeds[propName]; ok {
			payload[propName] = val
			continue
		}
		if len(prop.Enum) > 0 {
			payload[propName] = prop.Enum[0]
			continue
		}
		if prop.Type == "number" || prop.Type == "integer" {
			if propName == "limit" || propName == "perPage" {
				payload[propName] = 5
			}
		}
	}
	if len(payload) == 0 {
		return "{}"
	}
	b, _ := json.Marshal(payload)
	return string(b)
}

func hasUnresolvedRequiredProps(props map[string]OpenAPIProperty, required []string, seeds map[string]any) bool {
	for _, req := range required {
		if req == "cursor" {
			continue
		}
		if _, hasSeed := seeds[req]; hasSeed {
			continue
		}
		prop, exists := props[req]
		if exists && len(prop.Enum) > 0 {
			continue
		}
		return true
	}
	return false
}

func recordEndpoint(method, input, apiKey, outDir string) (int, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	return recordEndpointWithClient(method, input, apiKey, outDir, "https://api2.warera.io/trpc", client)
}

func recordEndpointWithClient(method, input, apiKey, outDir, baseURL string, client *http.Client) (int, error) {
	endpointURL := fmt.Sprintf("%s/%s", baseURL, method)

	if input != "" {
		endpointURL = fmt.Sprintf("%s?input=%s", endpointURL, url.QueryEscape(input))
	} else if method != "gameStat.getWorldDevelopment" {
		endpointURL = fmt.Sprintf("%s?input=%s", endpointURL, url.QueryEscape("{}"))
	}

	req, err := http.NewRequest("GET", endpointURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json")
	if apiKey != "" {
		req.Header.Set("x-api-key", apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
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
		return 0, fmt.Errorf("failed to decode tRPC envelope: %w", err)
	}

	if envelope.Error != nil {
		return 0, fmt.Errorf("tRPC server returned error: %s", envelope.Error.Message)
	}

	// Indent the original JSON bytes so large integers retain their exact value.
	var formatted bytes.Buffer
	if err := json.Indent(&formatted, envelope.Result.Data, "", "  "); err != nil {
		return 0, fmt.Errorf("failed to format JSON: %w", err)
	}
	indented := formatted.Bytes()

	outPath := filepath.Join(outDir, method+".json")
	if err := os.WriteFile(outPath, indented, 0644); err != nil {
		return 0, fmt.Errorf("failed to write fixture file: %w", err)
	}

	return len(indented), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
