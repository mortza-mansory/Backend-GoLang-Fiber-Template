// Package app composes the application: dependency wiring, route registration
// and lifecycle. It contains no business logic.
package app

import (
	"context"
	"log/slog"

	"github.com/yourorg/go-fiber-template/internal/config"
	"github.com/yourorg/go-fiber-template/internal/infra/cache"
	"github.com/yourorg/go-fiber-template/internal/infra/database"
	"github.com/yourorg/go-fiber-template/internal/infra/email"
	"github.com/yourorg/go-fiber-template/internal/infra/logger"
	"github.com/yourorg/go-fiber-template/internal/infra/observability"
	"github.com/yourorg/go-fiber-template/internal/infra/sms"
	"github.com/yourorg/go-fiber-template/internal/infra/storage"
	"github.com/yourorg/go-fiber-template/internal/infra/token"
	"github.com/yourorg/go-fiber-template/internal/server"
)

// Deps bundles every dependency the application needs. It is the dependency
// graph root built once at startup.
type Deps struct {
	Cfg           *config.Config
	Logger        *slog.Logger
	DB            *database.DB
	Redis         *cache.Cache
	Migrator      *database.Migrator
	Validator     *server.Validator
	Middleware    *server.Middleware
	Token         *token.Manager
	Email         email.Sender
	SMS           sms.Sender
	Storage       storage.Store
	Observability *observability.Shutdown
	Metrics       *observability.HTTPMetrics
	Modules       *Modules
}

// Build constructs the full dependency graph. Callers must Close it.
func Build(ctx context.Context, cfg *config.Config) (*Deps, error) {
	log := logger.New(cfg.Logger)

	otel, err := observability.Init(ctx, cfg.OTel)
	if err != nil {
		log.Error("failed to init observability", "error", err)
		return nil, err
	}

	db, err := database.New(ctx, cfg.Database)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		_ = otel.Close(ctx)
		return nil, err
	}

	rdb, err := cache.New(ctx, cfg.Redis)
	if err != nil {
		log.Error("failed to connect to redis", "error", err)
		_ = otel.Close(ctx)
		return nil, err
	}

	tok, err := token.NewManager(cfg.Security)
	if err != nil {
		log.Error("failed to init token manager", "error", err)
		return nil, err
	}

	validator := server.NewValidator()
	middleware := server.NewMiddleware(log)

	// External services: choose implementation based on config.
	var mail email.Sender = email.NewNoopSender()
	if cfg.External.Email.Host != "" {
		mail = email.NewSMTPClient(cfg.External.Email)
	}
	smsSender := sms.NewNoopSender()
	store, err := storage.NewLocalStore("storage")
	if err != nil {
		log.Error("failed to init storage", "error", err)
		return nil, err
	}

	deps := &Deps{
		Cfg:           cfg,
		Logger:        log,
		DB:            db,
		Redis:         rdb,
		Migrator:      database.NewMigrator(log),
		Validator:     validator,
		Middleware:    middleware,
		Token:         tok,
		Email:         mail,
		SMS:           smsSender,
		Storage:       store,
		Observability: otel,
		Metrics:       observability.NewHTTPMetrics(),
	}

	deps.Modules = newModules(deps)
	return deps, nil
}

// Close gracefully releases all resources in the correct order.
func (d *Deps) Close(ctx context.Context) {
	// Redis, then DB, then OpenTelemetry providers/exporters.
	if d.Redis != nil {
		_ = d.Redis.Close()
	}
	if d.DB != nil {
		d.DB.Close()
	}
	if d.Observability != nil {
		_ = d.Observability.Close(ctx)
	}
}
