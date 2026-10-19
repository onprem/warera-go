package warera

import (
	"fmt"
	"strings"

	"github.com/onprem/warera-go/internal/trpc"
)

// TRPCError is a structured error returned by the WarEra tRPC API.
// Use errors.As to inspect its message, code, and optional data.
type TRPCError = trpc.TRPCError

// TRPCErrorData contains structured metadata attached to a TRPCError.
type TRPCErrorData = trpc.TRPCErrorData

// HTTPError represents an HTTP failure that did not contain a tRPC error payload.
// Its status and response body are available for callers to inspect.
type HTTPError = trpc.HTTPError

// BatchLookupFailure describes one failed entity lookup in a bulk retrieval.
type BatchLookupFailure struct {
	// ID is the entity ID whose lookup failed.
	ID string
	// Err is the underlying API, transport, or decoding error.
	Err error
}

// BatchLookupError reports lookup failures from a bulk retrieval method such
// as [CompanyService.GetMany]. Such methods return no records when this error
// is returned. Use [errors.As] to inspect the failed IDs, and [errors.Is] or
// [errors.As] to inspect the underlying errors.
type BatchLookupError struct {
	// Entity is the plural entity name, such as "companies".
	Entity string
	// Failures contains one entry for each ID whose lookup failed.
	Failures []BatchLookupFailure
}

func (e *BatchLookupError) Error() string {
	if e == nil {
		return "warera: batch lookup failed"
	}
	parts := make([]string, 0, len(e.Failures))
	for _, failure := range e.Failures {
		if failure.Err != nil {
			parts = append(parts, fmt.Sprintf("%s: %v", failure.ID, failure.Err))
		} else {
			parts = append(parts, failure.ID)
		}
	}
	entity := e.Entity
	if entity == "" {
		entity = "records"
	}
	if len(parts) == 0 {
		return fmt.Sprintf("warera: batch lookup for %s failed", entity)
	}
	return fmt.Sprintf("warera: batch lookup for %s failed (%s)", entity, strings.Join(parts, "; "))
}

// Unwrap exposes all underlying lookup errors for errors.Is and errors.As.
func (e *BatchLookupError) Unwrap() []error {
	if e == nil {
		return nil
	}
	errs := make([]error, 0, len(e.Failures))
	for _, failure := range e.Failures {
		if failure.Err != nil {
			errs = append(errs, failure.Err)
		}
	}
	return errs
}
