package auth

import (
	"context"
	"log/slog"
)

// Service contains authentication business rules.
//
// Future responsibilities may include:
//   - validating credentials
//   - generating tokens (via an injected token.Manager)
//   - login policies, rate limiting, refresh tokens
//
// Do NOT put Fiber-specific logic here. The Service depends only on the
// Repository interface and injected collaborators, so it stays testable and
// independent of the HTTP framework.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService wires an auth Service.
//
// TODO(auth): when implementing real logic, add collaborators here, e.g. a
// *token.Manager, a *shared.IDGenerator, a password-hasher interface.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// Login validates credentials and returns a token pair.
//
// This is a commented placeholder. Replace the body with real logic.
//
// TODO(auth):
//   - load the user by email via repo.FindByEmail
//   - compare the password hash
//   - sign access/refresh tokens with token.Manager
//   - return tokens or ErrInvalidCredentials
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// s.logger.Debug("auth.login", "email", req.Email)
	return nil, ErrInvalidCredentials
}

// Register creates a new user account.
//
// TODO(auth):
//   - hash the password (bcrypt via security config)
//   - persist via repo.CreateUser
//   - return ErrUserAlreadyExists on conflict
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	return nil, ErrUserAlreadyExists
}

// Ping verifies the module's dependencies are reachable.
func (s *Service) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
