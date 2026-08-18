// Package database prepares PostgreSQL infrastructure: connection, pooling,
// health checks, graceful close and migration integration.
//
// It does NOT contain any business tables or SQL.
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/yourorg/go-fiber-template/internal/config"
)

// DB wraps a pgxpool.Pool to keep the driver detail contained.
type DB struct {
	pool *pgxpool.Pool
}

// New connects to PostgreSQL using the supplied configuration.
func New(ctx context.Context, cfg config.Database) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("database: parse config: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.MaxOpenConns)
	poolCfg.MinConns = int32(cfg.MaxIdleConns)
	poolCfg.MaxConnLifetime = cfg.ConnMaxLifetime
	poolCfg.MaxConnIdleTime = cfg.ConnMaxIdleTime

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("database: create pool: %w", err)
	}

	// Verify connectivity before returning.
	pingCtx, cancelPing := context.WithTimeout(ctx, 10*time.Second)
	defer cancelPing()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database: ping: %w", err)
	}

	return &DB{pool: pool}, nil
}

// Pool exposes the underlying pgxpool.Pool for queries.
func (db *DB) Pool() *pgxpool.Pool { return db.pool }

// Ping verifies the connection is alive.
func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Close gracefully closes the connection pool.
func (db *DB) Close() {
	db.pool.Close()
}
