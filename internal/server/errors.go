package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"

	"github.com/yourorg/go-fiber-template/internal/shared"
)

// HTTPError is an error carrying an explicit HTTP status. It is the bridge
// between the server layer and the HTTP framework.
type HTTPError struct {
	Status  int
	Code    shared.ErrorCode
	Message string
	Cause   error
}

// Error implements error.
func (e *HTTPError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap exposes the cause.
func (e *HTTPError) Unwrap() error { return e.Cause }

// NewHTTPError builds an HTTPError.
func NewHTTPError(status int, code shared.ErrorCode, message string, cause ...error) *HTTPError {
	e := &HTTPError{Status: status, Code: code, Message: message}
	if len(cause) > 0 {
		e.Cause = cause[0]
	}
	return e
}

// MapError converts any error into an HTTPError suitable for a response.
//
//   - *AppError -> mapped via its code
//   - *HTTPError -> passed through
//   - validation errors -> 400
//   - everything else -> 500 (internal), hidden in production
func MapError(err error, isProduction bool) *HTTPError {
	if err == nil {
		return nil
	}

	var he *HTTPError
	if errors.As(err, &he) {
		return he
	}

	var appErr *shared.AppError
	if errors.As(err, &appErr) {
		return NewHTTPError(StatusForCode(appErr.Code), appErr.Code, appErr.Message, appErr.Err)
	}

	var valErr *ValidationError
	if errors.As(err, &valErr) {
		return NewHTTPError(fiber.StatusBadRequest, shared.CodeBadRequest, valErr.Error())
	}

	// fasthttp/fiber sentinel errors.
	if errors.Is(err, fasthttp.ErrBodyTooLarge) {
		return NewHTTPError(fiber.StatusRequestEntityTooLarge, shared.CodeBadRequest, "request body too large", err)
	}
	if errors.Is(err, fiber.ErrUnprocessableEntity) {
		return NewHTTPError(fiber.StatusUnprocessableEntity, shared.CodeUnprocessable, "unprocessable entity", err)
	}

	// Unknown error: never leak internals to the client in production.
	if isProduction {
		return NewHTTPError(fiber.StatusInternalServerError, shared.CodeInternal, "internal server error", err)
	}
	return NewHTTPError(fiber.StatusInternalServerError, shared.CodeInternal, err.Error(), err)
}
