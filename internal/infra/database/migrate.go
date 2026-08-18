package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/yourorg/go-fiber-template/internal/config"
)

// Migrator runs SQL migrations against the database.
type Migrator struct {
	logger *slog.Logger
}

// NewMigrator creates a Migrator.
func NewMigrator(logger *slog.Logger) *Migrator {
	return &Migrator{logger: logger}
}

// Run executes pending migrations. Files live in ./migrations and are loaded
// via golang-migrate's file source.
func (m *Migrator) Run(ctx context.Context, cfg config.Database) error {
	sourceURL := cfg.MigrationsSource

	databaseURL := buildDatabaseURL(cfg)

	instance, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("migrations: create instance: %w", err)
	}
	defer instance.Close()

	if err := instance.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("migrations: no pending migrations")
			return nil
		}
		return fmt.Errorf("migrations: up: %w", err)
	}

	m.logger.Info("migrations: applied successfully")
	return nil
}

// buildDatabaseURL converts config into a pgx5:// URL expected by golang-migrate.
func buildDatabaseURL(cfg config.Database) string {
	q := url.Values{}
	q.Set("sslmode", cfg.SSLMode)

	u := url.URL{
		Scheme:   "pgx5",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host + ":" + cfg.Port,
		Path:     "/" + cfg.Name,
		RawQuery: q.Encode(),
	}
	return u.String()
}
