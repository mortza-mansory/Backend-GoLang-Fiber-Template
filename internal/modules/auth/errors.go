package auth

import "github.com/yourorg/go-fiber-template/internal/shared"

// This file defines module-specific sentinel errors.
//
// Use shared.AppError for common HTTP semantics and add module-specific
// messages/codes here.

// ErrInvalidCredentials is returned when a login fails.
//
// TODO(auth): replace placeholder behavior with real credential checks.
var ErrInvalidCredentials = shared.NewError(shared.CodeUnauthorized, "invalid credentials")

// ErrUserAlreadyExists is returned when trying to register a duplicate user.
var ErrUserAlreadyExists = shared.NewError(shared.CodeConflict, "user already exists")
