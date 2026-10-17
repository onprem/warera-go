.PHONY: update-specs

## Refresh the checked-in OpenAPI document and upstream response declarations.
update-specs:
	go run ./cmd/update-specs
