.PHONY: help test build generate fmt update-specs record record-force

## Show available development targets.
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //'

## Run the complete Go test suite.
test:
	go test ./...

## Build the SDK and examples.
build:
	go build ./...
	go build ./examples/...

## Regenerate request and filter models from the checked-in OpenAPI spec.
generate:
	go generate ./...

## Format Go source files.
fmt:
	gofmt -s -w .

## Refresh the checked-in OpenAPI document and upstream response declarations.
update-specs:
	go run ./cmd/update-specs

## Record missing live API fixtures (requires WARERA_API_KEY).
record:
	go run ./cmd/record

## Re-record all live API fixtures (overwrites existing fixtures).
record-force:
	go run ./cmd/record --force
