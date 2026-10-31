package imazeo

import (
	"fmt"
	"net/http"
	"strings"
)

// errorBody holds shared API error fields. Embedded under a name other than
// Error so the Error() method promotes cleanly on typed wrappers.
type errorBody struct {
	Status    int
	Code      string
	Message   string
	RequestID string
	Details   any
	cause     error
}

func (e *errorBody) Error() string {
	if e == nil {
		return "imazeo: <nil>"
	}
	parts := []string{}
	if e.Status != 0 {
		parts = append(parts, fmt.Sprintf("status=%d", e.Status))
	}
	if e.Code != "" {
		parts = append(parts, "code="+e.Code)
	}
	if e.RequestID != "" {
		parts = append(parts, "requestId="+e.RequestID)
	}
	msg := e.Message
	if msg == "" {
		msg = "request failed"
	}
	if len(parts) == 0 {
		return "imazeo: " + msg
	}
	return fmt.Sprintf("imazeo: %s (%s)", msg, strings.Join(parts, " "))
}

func (e *errorBody) Unwrap() error { return e.cause }

// Error is the base typed SDK error mapped from Imazeo API envelopes.
type Error = errorBody

// AuthenticationError is returned for HTTP 401.
type AuthenticationError struct{ *errorBody }

func (e *AuthenticationError) Unwrap() error { return e.errorBody }

// AuthorizationError is returned for HTTP 403.
type AuthorizationError struct{ *errorBody }

func (e *AuthorizationError) Unwrap() error { return e.errorBody }

// ValidationError is returned for HTTP 400 / 413 / 415.
type ValidationError struct{ *errorBody }

func (e *ValidationError) Unwrap() error { return e.errorBody }

// NotFoundError is returned for HTTP 404.
type NotFoundError struct{ *errorBody }

func (e *NotFoundError) Unwrap() error { return e.errorBody }

// ConflictError is returned for HTTP 409.
type ConflictError struct{ *errorBody }

func (e *ConflictError) Unwrap() error { return e.errorBody }

// RateLimitError is returned for HTTP 429.
type RateLimitError struct{ *errorBody }

func (e *RateLimitError) Unwrap() error { return e.errorBody }

// APIError is returned for other 4xx/5xx responses with an Imazeo envelope.
type APIError struct{ *errorBody }

func (e *APIError) Unwrap() error { return e.errorBody }

// NetworkError wraps transport failures and timeouts.
type NetworkError struct{ *errorBody }

func (e *NetworkError) Unwrap() error {
	if e.errorBody == nil {
		return nil
	}
	if e.errorBody.cause != nil {
		return e.errorBody.cause
	}
	return e.errorBody
}

func mapHTTPError(status int, code, message, requestID string, details any) error {
	base := &errorBody{
		Status:    status,
		Code:      code,
		Message:   message,
		RequestID: requestID,
		Details:   details,
	}
	if message == "" {
		base.Message = "request failed"
	}
	switch status {
	case http.StatusUnauthorized:
		return &AuthenticationError{errorBody: base}
	case http.StatusForbidden:
		return &AuthorizationError{errorBody: base}
	case http.StatusNotFound:
		return &NotFoundError{errorBody: base}
	case http.StatusConflict:
		return &ConflictError{errorBody: base}
	case http.StatusTooManyRequests:
		return &RateLimitError{errorBody: base}
	case http.StatusBadRequest, http.StatusRequestEntityTooLarge, http.StatusUnsupportedMediaType:
		return &ValidationError{errorBody: base}
	default:
		return &APIError{errorBody: base}
	}
}

func newNetworkError(message string, cause error, requestID string) error {
	return &NetworkError{errorBody: &errorBody{
		Message:   message,
		RequestID: requestID,
		cause:     cause,
	}}
}
