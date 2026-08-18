package observability

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// RegisterNoopProviders installs the standard OpenTelemetry no-op providers.
// This is used when observability is disabled so callers can keep using the
// global APIs without branching on configuration.
func RegisterNoopProviders() {
	otel.SetTracerProvider(trace.NewNoopTracerProvider())
	otel.SetMeterProvider(noop.NewMeterProvider())
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
