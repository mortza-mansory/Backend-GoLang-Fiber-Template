package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/yourorg/go-fiber-template/internal/infra/observability"
)

// NewFiber builds the base Fiber application with global middleware registered.
// Route registration happens separately in RegisterRoutes.
func (d *Deps) NewFiber() *fiber.App {
	cfg := d.Cfg

	app := fiber.New(fiber.Config{
		AppName:           cfg.App.Name,
		BodyLimit:         cfg.Server.BodyLimit,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
		ErrorHandler:      d.Middleware.ErrorHandler(cfg.App.Environment == "production"),
		EnablePrintRoutes: false,
	})

	// Recovery must be near the top so panics become 500s handled by the
	// centralized error handler.
	app.Use(recover.New())

	// Request ID first so downstream middleware/logs can reference it.
	app.Use(d.Middleware.RequestID())

	// CORS.
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Security.AllowedOrigin,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Compression.
	app.Use(compress.New())

	// Request logging (infrastructure-level, not business logic).
	app.Use(d.Middleware.RequestLogger())

	// OpenTelemetry HTTP tracing + metrics middleware. No-op when disabled.
	app.Use(observability.Middleware(d.Metrics))

	// Optional global rate limiting.
	if cfg.RateLimit.Enabled {
		app.Use(limiter.New(limiter.Config{
			Max:        cfg.RateLimit.Max,
			Expiration: cfg.RateLimit.Expiration,
		}))
	}

	return app
}
