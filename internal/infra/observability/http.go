package observability

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"

	"github.com/yourorg/go-fiber-template/internal/server"
)

// Middleware instruments HTTP requests with OpenTelemetry spans and metrics.
//
// It lives at the infrastructure layer so business handlers stay free of
// telemetry code. When OTel is disabled the no-op providers make this cheap.
func Middleware(metrics *HTTPMetrics) fiber.Handler {
	tracer := Tracer()
	return func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(c.Context(), "HTTP "+c.Method()+" "+c.Route().Path)
		span.SetAttributes(
			semconv.HTTPRequestMethodKey.String(c.Method()),
			semconv.URLPathKey.String(c.Path()),
			attribute.String("http.request_id", server.RequestID(c)),
		)

		start := time.Now()
		err := c.Next()

		status := c.Response().StatusCode()
		span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(status))
		if status >= 500 {
			span.SetStatus(codes.Error, "server error")
		}

		attrs := []attribute.KeyValue{
			semconv.HTTPRequestMethodKey.String(c.Method()),
			semconv.URLPathKey.String(c.Route().Path),
			semconv.HTTPResponseStatusCodeKey.Int(status),
		}
		metrics.RequestCount.Add(ctx, 1, metric.WithAttributes(attrs...))
		metrics.RequestLatency.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(attrs...))
		if status >= 500 {
			metrics.ErrorCount.Add(ctx, 1, metric.WithAttributes(attrs...))
		}

		span.End()
		return err
	}
}
