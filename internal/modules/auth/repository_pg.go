package auth

import (
	"context"

	"github.com/yourorg/go-fiber-template/internal/infra/database"
)

// PostgresRepository is a concrete implementation of the Repository interface
// backed by PostgreSQL. It lives here (in the module) rather than in the app
// layer so the app folder only wires dependencies.
//
// This is the infrastructure adapter: it depends on the database driver so the
// business/interface code in this package stays driver-agnostic.
type PostgresRepository struct {
	db *database.DB
}

// NewPostgresRepository creates a PostgreSQL-backed Repository.
func NewPostgresRepository(db *database.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Ping verifies the database connection.
func (r *PostgresRepository) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}

// TODO(auth): implement the real persistence methods here once the business
// logic is written. For example, after adding to the Repository interface:
//
//	func (r *PostgresRepository) Create(ctx context.Context, u *User) error {
//		_, err := r.db.Pool().Exec(ctx,
//			`INSERT INTO users (email, password_hash) VALUES ($1, $2)
//			 RETURNING id, created_at, updated_at`,
//			u.Email, u.PasswordHash,
//		)
//		return err
//	}
//
//	func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
//		var u User
//		err := r.db.Pool().QueryRow(ctx,
//			`SELECT id, email, password_hash, created_at, updated_at
//			 FROM users WHERE email = $1`, email,
//		).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
//		if err == pgx.ErrNoRows {
//			return nil, nil
//		}
//		if err != nil {
//			return nil, err
//		}
//		return &u, nil
//	}
