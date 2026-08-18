package server

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/yourorg/go-fiber-template/internal/shared"
)

// ContextKey is a type-safe key for request-scoped values.
type ContextKey string

const (
	// CtxRequestID is the injected request ID.
	CtxRequestID ContextKey = "request_id"
	// CtxPrincipal is the authenticated principal (set by auth middleware).
	CtxPrincipal ContextKey = "principal"
	// CtxTraceID is the OpenTelemetry trace ID when present.
	CtxTraceID ContextKey = "trace_id"
)

// Principal describes an authenticated caller. Modules may extend this via
// their own auth middleware.
type Principal struct {
	ID    shared.ID
	Roles []shared.Role
}

// RequestID returns the request ID or an empty string.
func RequestID(c *fiber.Ctx) string {
	if v, ok := c.Locals(string(CtxRequestID)).(string); ok {
		return v
	}
	return ""
}

// TraceID returns the OpenTelemetry trace ID or an empty string.
func TraceID(c *fiber.Ctx) string {
	if v, ok := c.Locals(string(CtxTraceID)).(string); ok {
		return v
	}
	return ""
}

// PrincipalFrom returns the authenticated principal, or nil.
func PrincipalFrom(c *fiber.Ctx) *Principal {
	if v, ok := c.Locals(string(CtxPrincipal)).(*Principal); ok {
		return v
	}
	return nil
}

// SetPrincipal stores the principal on the request context.
func SetPrincipal(c *fiber.Ctx, p *Principal) {
	c.Locals(string(CtxPrincipal), p)
}

// WithRequestValues returns a context.Context carrying the request ID and
// trace ID for use by infrastructure (DB, Redis, external calls).
func WithRequestValues(c *fiber.Ctx) context.Context {
	ctx := context.Background()
	if id := RequestID(c); id != "" {
		ctx = context.WithValue(ctx, CtxRequestID, id)
	}
	if id := TraceID(c); id != "" {
		ctx = context.WithValue(ctx, CtxTraceID, id)
	}
	return ctx
}
