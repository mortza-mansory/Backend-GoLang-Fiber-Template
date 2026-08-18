// Package swagger wires the OpenAPI/Swagger UI into the HTTP server.
//
// The spec is generated from annotations using the swag CLI:
//
//	swag init -g cmd/api/main.go -o docs
//
// and served at /swagger. Toggle via SWAGGER_ENABLED.
package swagger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "github.com/yourorg/go-fiber-template/docs" // generated spec
)

// Enabled reports whether swagger UI should be mounted.
func Enabled() bool {
	return getEnvBool("SWAGGER_ENABLED", true)
}

// Register mounts the Swagger UI at /swagger and the raw JSON at /swagger/doc.json.
func Register(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("docs/swagger.json")
	})
}
