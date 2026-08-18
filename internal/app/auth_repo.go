package app

import (
	"context"

	"github.com/yourorg/go-fiber-template/internal/infra/database"
)

// authRepo implements the auth.Repository interface using PostgreSQL.
//
// It lives in the app package (not in the module) because it is an
// infrastructure adapter wiring the module's business interface to a driver.
//
// TODO(auth): add the real methods from auth.Repository here once defined,
// e.g.:
//
//	func (r *authRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
//	    // query the users table via r.db.Pool()
//	}
type authRepository struct {
	db *database.DB
}

// Ping verifies the database connection.
func (r *authRepository) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}
