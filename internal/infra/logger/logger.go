// Package logger provides structured logging built on log/slog.
//
// It supports levels, structured fields, JSON output (production) and
// human-readable text output (development). Request correlation and
// trace/span IDs are added by middleware in the server layer.
package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/yourorg/go-fiber-template/internal/config"
)

// New creates an *slog.Logger based on configuration.
//   - JSON output when Logger.JSON is true (recommended for production).
//   - Text output otherwise (nice for local development).
func New(cfg config.Logger) *slog.Logger {
	level := parseLevel(cfg.Level)

	var handler slog.Handler
	if cfg.JSON {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	return slog.New(handler)
}

// parseLevel maps a string level to a slog.Level.
func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "info", "":
		return slog.LevelInfo
	default:
		return slog.LevelInfo
	}
}
