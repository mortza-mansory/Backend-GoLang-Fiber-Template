// Package server contains reusable HTTP/server concerns: standardized
// responses, request parsing, validation, error mapping, middleware and
// request-context helpers. It contains no business rules.
package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/yourorg/go-fiber-template/internal/shared"
)

// SuccessEnvelope is the standard success response shape.
type SuccessEnvelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorEnvelope is the standard error response shape.
type ErrorEnvelope struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

// ErrorDetail holds a single machine-readable error.
type ErrorDetail struct {
	Code    shared.ErrorCode `json:"code"`
	Message string           `json:"message"`
}

// OK writes a 200 success response with optional data and meta.
func OK(c *fiber.Ctx, data interface{}, meta interface{}) error {
	return c.Status(fiber.StatusOK).JSON(SuccessEnvelope{Success: true, Data: data, Meta: meta})
}

// Created writes a 201 success response.
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(SuccessEnvelope{Success: true, Data: data})
}

// NoContent writes a 204 response with no body.
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Fail writes an error response using the mapped HTTP status for the code.
func Fail(c *fiber.Ctx, code shared.ErrorCode, message string) error {
	status := StatusForCode(code)
	return c.Status(status).JSON(ErrorEnvelope{
		Success: false,
		Error:   ErrorDetail{Code: code, Message: message},
	})
}

// FailStatus writes an error response with an explicit HTTP status.
func FailStatus(c *fiber.Ctx, status int, code shared.ErrorCode, message string) error {
	return c.Status(status).JSON(ErrorEnvelope{
		Success: false,
		Error:   ErrorDetail{Code: code, Message: message},
	})
}

// StatusForCode maps a shared.ErrorCode to an HTTP status code.
func StatusForCode(code shared.ErrorCode) int {
	switch code {
	case shared.CodeBadRequest:
		return fiber.StatusBadRequest
	case shared.CodeUnauthorized:
		return fiber.StatusUnauthorized
	case shared.CodeForbidden:
		return fiber.StatusForbidden
	case shared.CodeNotFound:
		return fiber.StatusNotFound
	case shared.CodeConflict:
		return fiber.StatusConflict
	case shared.CodeUnprocessable:
		return fiber.StatusUnprocessableEntity
	case shared.CodeTooManyRequests:
		return fiber.StatusTooManyRequests
	case shared.CodeInternal, "":
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusInternalServerError
	}
}
