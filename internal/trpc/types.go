package trpc

import (
	"encoding/json"
	"fmt"
)

// ResponseEnvelope represents a standard tRPC response object.
type ResponseEnvelope struct {
	Result *ResultEnvelope `json:"result,omitempty"`
	Error  *TRPCError      `json:"error,omitempty"`
}

// ResultEnvelope contains the output data of a successful tRPC procedure.
type ResultEnvelope struct {
	Data json.RawMessage `json:"data"`
}

// TRPCError represents an error returned by the tRPC server.
type TRPCError struct {
	Message string         `json:"message"`
	Code    int            `json:"code"`
	Data    *TRPCErrorData `json:"data,omitempty"`
}

// TRPCErrorData contains detailed error metadata.
type TRPCErrorData struct {
	Code       string `json:"code"`
	HTTPStatus int    `json:"httpStatus"`
	Path       string `json:"path,omitempty"`
}

func (e *TRPCError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Data != nil && e.Data.Code != "" {
		return fmt.Sprintf("trpc: %s (code: %s, httpStatus: %d)", e.Message, e.Data.Code, e.Data.HTTPStatus)
	}
	return fmt.Sprintf("trpc: %s (code: %d)", e.Message, e.Code)
}

// HTTPError represents an unexpected HTTP status from the server.
type HTTPError struct {
	StatusCode int
	StatusText string
	Body       []byte
}

func (e *HTTPError) Error() string {
	if len(e.Body) > 0 {
		return fmt.Sprintf("http %d %s: %s", e.StatusCode, e.StatusText, string(e.Body))
	}
	return fmt.Sprintf("http %d %s", e.StatusCode, e.StatusText)
}
