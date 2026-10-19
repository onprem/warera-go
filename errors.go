package warera

import "github.com/onprem/warera-go/internal/trpc"

// TRPCError is a structured error returned by the WarEra tRPC API.
// Use errors.As to inspect its message, code, and optional data.
type TRPCError = trpc.TRPCError

// TRPCErrorData contains structured metadata attached to a TRPCError.
type TRPCErrorData = trpc.TRPCErrorData

// HTTPError represents an HTTP failure that did not contain a tRPC error payload.
// Its status and response body are available for callers to inspect.
type HTTPError = trpc.HTTPError
