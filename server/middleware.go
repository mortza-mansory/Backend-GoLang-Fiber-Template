package server

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/yourorg/go-fiber-template/internal/shared"
)

// Middleware bundles reusable, framework-facing middleware. It is wired in the
// app layer, not inside business modules.
type Middleware struct {
	logger *slog.Logger
}

// NewMiddleware creates a Middleware.
func NewMiddleware(logger *slog.Logger) *Middleware {
	return &Middleware{logger: logger}
}

// RequestID injects a request ID header (X-Request-ID) and stores it locally.
func (m *Middleware) RequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header:     fiber.HeaderXRequestID,
		ContextKey: string(CtxRequestID),
	})
}

// RequestLogger emits one structured log line per request.
func (m *Middleware) RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		attrs := []any{
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration_ms", time.Since(start)),
			slog.String("request_id", RequestID(c)),
			slog.String("ip", c.IP()),
		}
		if tid := TraceID(c); tid != "" {
			attrs = append(attrs, slog.String("trace_id", tid))
		}
		if err != nil {
			attrs = append(attrs, slog.Any("error", err))
			m.logger.Error("http_request", attrs...)
			return err
		}
		m.logger.Info("http_request", attrs...)
		return nil
	}
}

// ErrorHandler is the centralized Fiber error handler. It maps any error to a
// consistent JSON envelope and never leaks internals in production.
func (m *Middleware) ErrorHandler(isProduction bool) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		he := MapError(err, isProduction)
		// Handlers may already have sent a response; only write a new one then.
		if c.Response().StatusCode() != fiber.StatusOK {
			return nil
		}
		return c.Status(he.Status).JSON(ErrorEnvelope{
			Success: false,
			Error:   ErrorDetail{Code: he.Code, Message: he.Message},
		})
	}
}

// NotFound handles unmatched routes with a consistent JSON response.
func (m *Middleware) NotFound() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return Fail(c, shared.CodeNotFound, "route not found")
	}
}
