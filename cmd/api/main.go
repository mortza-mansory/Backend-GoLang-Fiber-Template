// Command api is the HTTP entry point. It only performs bootstrapping;
// business logic lives in internal/.
//
// Boot sequence:
//
//	Load config -> logger -> observability -> DB -> migrations -> Redis
//	-> build deps -> Fiber app -> middleware -> routes -> serve -> shutdown
//
//	@title			Go Fiber Template API
//	@version		0.1.0
//	@description	Production-ready Go + Fiber modular monolith template.
//	@host			localhost:8080
//	@BasePath		/api
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/yourorg/go-fiber-template/internal/app"
	"github.com/yourorg/go-fiber-template/internal/config"
)

func main() {
	// Best-effort load of .env for local development. Missing file is fine.
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		fatal("load config", err)
	}

	log := slog.Default()
	log.Info("starting application",
		"name", cfg.App.Name,
		"env", cfg.App.Environment,
		"version", cfg.App.Version,
	)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	deps, err := app.Build(ctx, cfg)
	if err != nil {
		fatal("build application", err)
	}

	// Run migrations before serving (when enabled).
	if cfg.Database.RunMigrations {
		if err := deps.Migrator.Run(ctx, cfg.Database); err != nil {
			log.Error("failed to run migrations", "error", err)
		}
	}

	fiberApp := deps.NewFiber()
	deps.RegisterRoutes(fiberApp)

	go func() {
		addr := cfg.Server.Address()
		log.Info("http server listening", "addr", addr)
		if err := fiberApp.Listen(addr); err != nil {
			// During shutdown Listen returns nil; anything else is fatal.
			log.Error("http server error", "error", err)
			stop()
		}
	}()

	// Wait for shutdown signal.
	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// HTTP server first, then background resources via deps.Close.
	_ = fiberApp.ShutdownWithContext(shutdownCtx)
	deps.Close(shutdownCtx)

	log.Info("shutdown complete")
}

func fatal(msg string, err error) {
	slog.Default().Error(msg, "error", err)
	os.Exit(1)
}
