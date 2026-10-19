package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	openAPIURL = "https://api2.warera.io/openapi.json"
	userAgent  = "warera-go-spec-updater"
)

type openAPIDocument struct {
	OpenAPI string `json:"openapi"`
	Info    struct {
		Version string `json:"version"`
	} `json:"info"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "update-specs: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	specDir := filepath.Join("spec")

	openAPIBytes, err := get(openAPIURL)
	if err != nil {
		return err
	}
	var openAPI openAPIDocument
	if err := json.Unmarshal(openAPIBytes, &openAPI); err != nil {
		return fmt.Errorf("decode OpenAPI document: %w", err)
	}
	if len(openAPI.OpenAPI) < 2 || openAPI.OpenAPI[:2] != "3." || openAPI.Info.Version == "" {
		return fmt.Errorf("%s did not return a versioned OpenAPI 3 document", openAPIURL)
	}
	var formattedOpenAPI bytes.Buffer
	if err := json.Indent(&formattedOpenAPI, openAPIBytes, "", "  "); err != nil {
		return fmt.Errorf("format OpenAPI document: %w", err)
	}
	prettyOpenAPI := append(bytes.TrimSpace(formattedOpenAPI.Bytes()), '\n')

	outputs := []fileOutput{{path: filepath.Join(specDir, "openapi.json"), data: prettyOpenAPI}}
	changed, err := stageAndReplace(specDir, outputs)
	if err != nil {
		return err
	}

	fmt.Printf("Updated %d spec file(s).\n", changed)
	fmt.Printf("OpenAPI: %s\n", openAPI.Info.Version)
	return nil
}

func get(url string) ([]byte, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	return getWithClient(url, client)
}

func getWithClient(url string, client *http.Client) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request for %s: %w", url, err)
	}
	request.Header.Set("User-Agent", userAgent)

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", url, err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("download %s: HTTP %s", url, response.Status)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", url, err)
	}
	return data, nil
}

type fileOutput struct {
	path string
	data []byte
}

func stageAndReplace(specDir string, outputs []fileOutput) (int, error) {
	tempDir, err := os.MkdirTemp(specDir, ".update-specs-")
	if err != nil {
		return 0, fmt.Errorf("create staging directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	type stagedFile struct {
		temporary   string
		destination string
	}
	var changed []stagedFile
	for index, output := range outputs {
		temporary := filepath.Join(tempDir, fmt.Sprintf("%d", index))
		if err := os.WriteFile(temporary, output.data, 0o644); err != nil {
			return 0, fmt.Errorf("stage %s: %w", output.path, err)
		}
		current, err := os.ReadFile(output.path)
		if err == nil && bytes.Equal(current, output.data) {
			continue
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return 0, fmt.Errorf("read %s: %w", output.path, err)
		}
		changed = append(changed, stagedFile{temporary: temporary, destination: output.path})
	}

	for _, file := range changed {
		if err := os.Rename(file.temporary, file.destination); err != nil {
			return 0, fmt.Errorf("replace %s: %w", file.destination, err)
		}
	}
	return len(changed), nil
}
