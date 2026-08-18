package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/yourorg/go-fiber-template/internal/config"
)

// defaultResource builds the OTel resource describing this service.
func defaultResource(cfg config.OTel) *resource.Resource {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(cfg.ServiceName),
		semconv.ServiceVersion(cfg.ServiceVersion),
		semconv.DeploymentEnvironment(cfg.Environment),
	}
	return resource.NewWithAttributes(semconv.SchemaURL, attrs...)
}

// setupMetrics configures the OTLP metric exporter and registers runtime +
// process metrics.
func setupMetrics(ctx context.Context, cfg config.OTel, shutdown *Shutdown) error {
	endpoint := cfg.OTLPEndpoint
	if cfg.MetricsEndpoint != "" {
		endpoint = cfg.MetricsEndpoint
	}

	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(endpoint),
	}
	if cfg.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exp, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return fmt.Errorf("observability: create metric exporter: %w", err)
	}

	reader := sdkmetric.NewPeriodicReader(exp)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(defaultResource(cfg)),
	)

	otel.SetMeterProvider(mp)

	shutdown.Add(func(ctx context.Context) error {
		return mp.Shutdown(ctx)
	})

	return nil
}

// HTTPMetrics exposes the infrastructure metric instruments used by the Fiber
// middleware. Handlers never need to touch these directly.
type HTTPMetrics struct {
	RequestCount   metric.Int64Counter
	RequestLatency metric.Float64Histogram
	ErrorCount     metric.Int64Counter
}

// NewHTTPMetrics builds and registers the HTTP metric instruments.
func NewHTTPMetrics() *HTTPMetrics {
	meter := otel.Meter("github.com/yourorg/go-fiber-template/http")

	count, _ := meter.Int64Counter(
		"http.server.request.count",
		metric.WithDescription("Total number of HTTP requests received"),
	)
	latency, _ := meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
	)
	errors, _ := meter.Int64Counter(
		"http.server.error.count",
		metric.WithDescription("Total number of HTTP errors"),
	)

	return &HTTPMetrics{
		RequestCount:   count,
		RequestLatency: latency,
		ErrorCount:     errors,
	}
}
