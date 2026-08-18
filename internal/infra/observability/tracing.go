package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/yourorg/go-fiber-template/internal/config"
)

// setupTracing configures the OTLP trace exporter and a batch span processor,
// then sets the global TracerProvider and propagators.
func setupTracing(ctx context.Context, cfg config.OTel, shutdown *Shutdown) error {
	endpoint := cfg.OTLPEndpoint
	if cfg.TracesEndpoint != "" {
		endpoint = cfg.TracesEndpoint
	}

	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(endpoint),
	}
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exp, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return fmt.Errorf("observability: create trace exporter: %w", err)
	}

	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SamplingRatio))

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(defaultResource(cfg)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	shutdown.Add(func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	})

	return nil
}

// Tracer returns the application tracer. Business modules should generally NOT
// call this; infrastructure handles instrumentation.
func Tracer() trace.Tracer {
	return otel.Tracer("github.com/yourorg/go-fiber-template")
}
