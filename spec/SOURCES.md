# Specification sources and attribution

The checked-in files in this directory come from separate upstream sources:

- [`openapi.json`](openapi.json) is the WarEra API's published OpenAPI document: <https://api2.warera.io/openapi.json>. Its `info.version` identifies the API specification version.
- [`Responses.d.ts`](Responses.d.ts) contains generated response declarations from [WarEraProjects/TRPC](https://github.com/WarEraProjects/TRPC). Its header records the exact upstream commit and `LICENSE` URL.

The response declarations retain their upstream generated header and link to the upstream `LICENSE`; consult that file at the pinned revision for its terms. The root `LICENSE` in this repository applies to `warera-go` code and does not replace the terms for upstream material.

At the time this attribution was added, GitHub identifies the TRPC repository as Apache-2.0 while its `package.json` declares MIT. The file header links to the exact upstream `LICENSE` at the pinned revision instead of assuming either label applies to this snapshot.

Refresh both snapshots from the repository root with:

```sh
make update-specs
```

The refresh requires Go 1.23+ and network access. It validates both downloads, pretty-prints the OpenAPI JSON, and stages the files before updating them. Review the source snapshots together before regenerating Go models with `go generate ./...`.
