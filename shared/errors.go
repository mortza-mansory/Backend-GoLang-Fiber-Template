// Package shared holds only genuinely cross-module concepts: base error
// types, pagination helpers, role constants and small shared types.
//
// It is NOT a dumping ground for module-specific code. If a type only makes
// sense inside one feature, keep it there.
package shared

// ErrorCode is a stable machine-readable error identifier returned in API
// responses. Consumers can switch on it without parsing messages.
type ErrorCode string

const (
	// CodeBadRequest indicates a client input error.
	CodeBadRequest ErrorCode = "BAD_REQUEST"
	// CodeUnauthorized indicates missing/invalid credentials.
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	// CodeForbidden indicates insufficient permissions.
	CodeForbidden ErrorCode = "FORBIDDEN"
	// CodeNotFound indicates the requested resource does not exist.
	CodeNotFound ErrorCode = "NOT_FOUND"
	// CodeConflict indicates a state conflict with existing data.
	CodeConflict ErrorCode = "CONFLICT"
	// CodeUnprocessable indicates valid but semantically invalid input.
	CodeUnprocessable ErrorCode = "UNPROCESSABLE_ENTITY"
	// CodeInternal indicates an unexpected server error.
	CodeInternal ErrorCode = "INTERNAL_ERROR"
	// CodeTooManyRequests indicates rate limiting.
	CodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
)

// AppError is the canonical error type for expected, client-visible failures.
//
// Unexpected/internal errors should generally be wrapped by the server layer
// rather than constructed as AppError at the call site.
type AppError struct {
	Code    ErrorCode
	Message string
	// Err optionally holds the underlying cause. It is never exposed to clients.
	Err error
}

// NewError builds an AppError with an optional underlying cause.
func NewError(code ErrorCode, message string, cause ...error) *AppError {
	e := &AppError{Code: code, Message: message}
	if len(cause) > 0 && cause[0] != nil {
		e.Err = cause[0]
	}
	return e
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap exposes the underlying cause for errors.Is/errors.As.
func (e *AppError) Unwrap() error { return e.Err }
