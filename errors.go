package warera

import (
	"errors"
	"fmt"
	"net"
	"net/http"
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

// ErrNotFound is returned when an entity lookup's result is null, such as
// [TournamentService.GetByID] with an unknown ID. Many procedures report
// unknown IDs with a NOT_FOUND [TRPCError] instead; use [IsNotFound] to detect
// both.
var ErrNotFound = trpc.ErrNotFound

// IsNotFound reports whether err means the requested entity does not exist:
// [ErrNotFound], a NOT_FOUND [TRPCError], or an HTTP 404. For a
// [BatchLookupError] it reports whether any lookup was not found.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) || hasTRPCCode(err, "NOT_FOUND") || hasHTTPStatus(err, http.StatusNotFound)
}

// IsRateLimited reports whether err is the API rejecting a request for
// exceeding its rate limit, after the client's retries.
func IsRateLimited(err error) bool {
	return hasTRPCCode(err, "TOO_MANY_REQUESTS") || hasHTTPStatus(err, http.StatusTooManyRequests)
}

// IsTemporary reports whether err is likely to succeed if retried later: a
// rate limit, timeout, server error, or network failure. Validation,
// authorization, and not-found errors are not temporary.
func IsTemporary(err error) bool {
	if IsRateLimited(err) || hasTRPCCode(err, "TIMEOUT", "INTERNAL_SERVER_ERROR", "SERVICE_UNAVAILABLE", "GATEWAY_TIMEOUT", "BAD_GATEWAY") {
		return true
	}
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		switch code := httpErr.StatusCode; {
		case code == http.StatusRequestTimeout, code == http.StatusTooEarly, code >= 500:
			return true
		}
	}
	var netErr net.Error
	return errors.As(err, &netErr)
}

func hasTRPCCode(err error, codes ...string) bool {
	var trpcErr *TRPCError
	if !errors.As(err, &trpcErr) || trpcErr.Data == nil {
		return false
	}
	for _, code := range codes {
		if trpcErr.Data.Code == code {
			return true
		}
	}
	return false
}

func hasHTTPStatus(err error, status int) bool {
	var httpErr *HTTPError
	return errors.As(err, &httpErr) && httpErr.StatusCode == status
}

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
