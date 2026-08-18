package app

import (
	"github.com/gofiber/fiber/v2"

	"github.com/yourorg/go-fiber-template/internal/infra/swagger"
	"github.com/yourorg/go-fiber-template/internal/server"
	"github.com/yourorg/go-fiber-template/internal/shared"
)

// RegisterRoutes mounts all application routes: health, swagger and modules.
func (d *Deps) RegisterRoutes(app *fiber.App) {
	// Liveness: the process is alive. No dependency checks.
	app.Get("/health", func(c *fiber.Ctx) error {
		return server.OK(c, fiber.Map{"status": "alive"}, nil)
	})

	// Readiness: dependencies are reachable.
	app.Get("/ready", func(c *fiber.Ctx) error {
		if err := d.DB.Ping(c.Context()); err != nil {
			return server.Fail(c, shared.CodeInternal, "database not ready")
		}
		if err := d.Redis.Ping(c.Context()); err != nil {
			return server.Fail(c, shared.CodeInternal, "cache not ready")
		}
		return server.OK(c, fiber.Map{"status": "ready"}, nil)
	})

	// API group. All feature modules mount under /api.
	api := app.Group("/api")
	d.registerModuleRoutes(api)

	// Swagger UI and spec (see internal/infra/swagger).
	registerSwagger(app)

	// 404 fallback.
	app.Use(d.Middleware.NotFound())
}

// registerModuleRoutes wires every feature module's routes.
func (d *Deps) registerModuleRoutes(api fiber.Router) {
	d.Modules.Auth.Handler.RegisterRoutes(api)
}

// registerSwagger mounts the OpenAPI UI when enabled.
func registerSwagger(app *fiber.App) {
	if swagger.Enabled() {
		swagger.Register(app)
	}
}
