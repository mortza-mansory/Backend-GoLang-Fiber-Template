// Package observability prepares production-ready OpenTelemetry integration.
//
// The design goal: observability should observe the application, not reshape
// the business code. Instrumentation is wired at the infrastructure/framework
// layer (HTTP middleware, DB/Redis drivers) so business functions stay clean.
//
// All of it is optional: when OTEL_ENABLED=false the application runs with a
// no-op tracer/metric provider and requires no external collector.
package observability

import (
	"context"

	"github.com/yourorg/go-fiber-template/internal/config"
)

// Shutdown is a set of cleanup functions run in reverse order at shutdown.
type Shutdown struct {
	funcs []func(context.Context) error
}

// Add registers a cleanup function.
func (s *Shutdown) Add(f func(context.Context) error) {
	s.funcs = append(s.funcs, f)
}

// Close runs all registered cleanup functions in reverse registration order.
func (s *Shutdown) Close(ctx context.Context) error {
	for i := len(s.funcs) - 1; i >= 0; i-- {
		if err := s.funcs[i](ctx); err != nil {
			return err
		}
	}
	return nil
}

// Init wires up tracing and metrics based on config. It always returns a
// Shutdown that must be closed on graceful shutdown.
//
// When OTel is disabled it returns the standard no-op providers, so the rest
// of the application needs no conditional branching.
func Init(ctx context.Context, cfg config.OTel) (*Shutdown, error) {
	shutdown := &Shutdown{}

	if !cfg.Enabled {
		// Register no-op providers so the app works without a collector.
		RegisterNoopProviders()
		return shutdown, nil
	}

	if err := setupTracing(ctx, cfg, shutdown); err != nil {
		return nil, err
	}
	if err := setupMetrics(ctx, cfg, shutdown); err != nil {
		return nil, err
	}

	return shutdown, nil
}
