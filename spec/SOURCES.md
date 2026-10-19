# Specification sources and response model policy

- [`openapi.json`](openapi.json) is WarEra's published OpenAPI document: <https://api2.warera.io/openapi.json>. Its `info.version` identifies the specification version. Refresh it with `make update-specs`.
- Response models in `models_response.go` and `models_additional.go` are maintained as Go code in this repository. They are informed by observed API behavior, recorded fixtures, and the [WarEraProjects/TRPC response declarations](https://github.com/WarEraProjects/TRPC/blob/e7b7765d2487b43e01d3b13654bedf9ee5edb16f/src/api/Responses.d.ts) as a reference. The declaration file is not vendored or required by this SDK or its generation workflow.

The OpenAPI document is used to generate request and filter models. Recorded fixtures are examples for decoding and tests; they do not generate response model types. When endpoint responses change, update the relevant Go model and add or update fixture-backed decoding assertions for the fields the SDK relies on.
