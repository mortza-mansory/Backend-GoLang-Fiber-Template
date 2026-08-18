package auth

import (
	"github.com/gofiber/fiber/v2"

	"github.com/yourorg/go-fiber-template/internal/server"
)

// Handler is the HTTP layer for the auth module.
//
// It is the ONLY place in this module that touches Fiber. It delegates all
// work to the Service and maps results to the standardized server responses.
type Handler struct {
	service  *Service
	validate *server.Validator
}

// NewHandler wires an auth Handler.
func NewHandler(service *Service, validate *server.Validator) *Handler {
	return &Handler{service: service, validate: validate}
}

// RegisterRoutes mounts the auth routes. This is called from the app layer.
//
// TODO(auth): implement /auth/login, /auth/register, /auth/refresh, etc.
// below using the real service methods. Each handler should:
//  1. bind + validate the request via h.validate.Bind(c, &dto)
//  2. call the service
//  3. return server.OK(c, data) on success
func (h *Handler) RegisterRoutes(group fiber.Router) {
	auth := group.Group("/auth")

	auth.Post("/login", h.Login)
	auth.Post("/register", h.Register)
	auth.Get("/ping", h.Ping)
}

// Login handles POST /auth/login.
//
//	@Summary		Login
//	@Description	Authenticate a user and return tokens.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	LoginRequest	true	"Credentials"
//	@Success		200		{object}	server.SuccessEnvelope
//	@Failure		401		{object}	server.ErrorEnvelope
//	@Router			/auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := h.validate.Bind(c, &req); err != nil {
		return err
	}
	resp, err := h.service.Login(c.Context(), req)
	if err != nil {
		return err
	}
	return server.OK(c, resp, nil)
}

// Register handles POST /auth/register.
//
//	@Summary		Register
//	@Description	Create a new user account.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	RegisterRequest	true	"Account details"
//	@Success		201		{object}	server.SuccessEnvelope
//	@Failure		409		{object}	server.ErrorEnvelope
//	@Router			/auth/register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := h.validate.Bind(c, &req); err != nil {
		return err
	}
	resp, err := h.service.Register(c.Context(), req)
	if err != nil {
		return err
	}
	return server.Created(c, resp)
}

// Ping handles GET /auth/ping.
//
//	@Summary		Auth ping
//	@Description	Verify auth module dependencies are reachable.
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	server.SuccessEnvelope
//	@Router			/auth/ping [get]
func (h *Handler) Ping(c *fiber.Ctx) error {
	if err := h.service.Ping(c.Context()); err != nil {
		return err
	}
	return server.OK(c, fiber.Map{"status": "ok"}, nil)
}
