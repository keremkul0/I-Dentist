package tracingMiddleware

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Middleware creates a Fiber middleware for OpenTelemetry tracing
func Middleware(serviceName string) fiber.Handler {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(c *fiber.Ctx) error {
		// Extract context from incoming request headers
		ctx := propagator.Extract(c.Context(), &FiberCarrier{c: c})

		// Start a new span
		spanName := c.Method() + " " + c.Route().Path
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithAttributes(
				attribute.String("http.method", c.Method()),
				attribute.String("http.url", c.OriginalURL()),
				attribute.String("http.target", c.Path()),
				attribute.String("http.host", c.Hostname()),
				attribute.String("http.scheme", c.Protocol()),
				attribute.String("http.user_agent", c.Get("User-Agent")),
				attribute.String("http.route", c.Route().Path),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		// Store the context in Fiber context
		c.SetUserContext(ctx)

		// Continue with the request
		err := c.Next()

		// Record response status
		span.SetAttributes(
			attribute.Int("http.status_code", c.Response().StatusCode()),
		)

		// Record error if any
		if err != nil {
			span.RecordError(err)
		}

		return err
	}
}

// FiberCarrier adapts Fiber context to propagation.TextMapCarrier
type FiberCarrier struct {
	c *fiber.Ctx
}

func (fc *FiberCarrier) Get(key string) string {
	return fc.c.Get(key)
}

func (fc *FiberCarrier) Set(key string, value string) {
	fc.c.Set(key, value)
}

func (fc *FiberCarrier) Keys() []string {
	keys := make([]string, 0)
	fc.c.Request().Header.VisitAll(func(key, value []byte) {
		keys = append(keys, string(key))
	})
	return keys
}
