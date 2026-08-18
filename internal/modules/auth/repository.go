package auth

import "context"

// Repository is the persistence abstraction for the auth module.
//
// The business layer depends on this interface, not on PostgreSQL. The actual
// implementation (e.g. using pgxpool) lives elsewhere and is injected.
//
// TODO(auth): add the methods your real auth flow needs, for example:
//
//	CreateUser(ctx, *User) error
//	FindByEmail(ctx, string) (*User, error)
//	UpdatePassword(ctx, string, string) error
type Repository interface {
	// Ping is a placeholder proving the boundary wires up.
	Ping(ctx context.Context) error
}
