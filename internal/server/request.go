package server

import (
	"encoding/json"
	"io"

	"github.com/gofiber/fiber/v2"

	"github.com/yourorg/go-fiber-template/internal/shared"
)

// Bind parses the JSON request body into dst and validates it.
// Returns an *AppError or *ValidationError that MapError understands.
func (val *Validator) Bind(c *fiber.Ctx, dst interface{}) error {
	if err := c.BodyParser(dst); err != nil {
		if err == io.EOF {
			return shared.NewError(shared.CodeBadRequest, "empty request body")
		}
		return shared.NewError(shared.CodeBadRequest, "invalid request body", err)
	}
	if val != nil {
		if err := val.ValidateStruct(dst); err != nil {
			return err
		}
	}
	return nil
}

// Query parses the request query string into dst. Non-existent keys are left
// untouched so structs with defaults work naturally.
func Query(c *fiber.Ctx, dst interface{}) error {
	if err := c.QueryParser(dst); err != nil {
		return shared.NewError(shared.CodeBadRequest, "invalid query parameters", err)
	}
	return nil
}

// PathParams populates dst from URL path parameters.
func PathParams(c *fiber.Ctx, dst interface{}) error {
	if err := c.ParamsParser(dst); err != nil {
		return shared.NewError(shared.CodeBadRequest, "invalid path parameters", err)
	}
	return nil
}

// ParseBody is a thin wrapper for handlers that only need body parsing.
func ParseBody(c *fiber.Ctx, dst interface{}) error {
	if err := json.Unmarshal(c.Body(), dst); err != nil {
		return shared.NewError(shared.CodeBadRequest, "invalid request body", err)
	}
	return nil
}
