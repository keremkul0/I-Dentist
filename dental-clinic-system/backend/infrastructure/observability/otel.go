package observability

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// TracerShutdown is a function type for shutting down the tracer
type TracerShutdown func(context.Context) error

type TracingConfig struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	Enabled        bool
	OTLPEndpoint   string
}

// InitTracer initializes the OpenTelemetry tracer with configuration
func InitTracer(ctx context.Context, config TracingConfig) (TracerShutdown, error) {
	if !config.Enabled {
		log.Info().Msg("Tracing is disabled")
		return func(ctx context.Context) error { return nil }, nil
	}

	log.Info().Str("service", config.ServiceName).Msg("Initializing OpenTelemetry tracer...")

	// Create OTLP HTTP exporter
	opts := []otlptracehttp.Option{
		otlptracehttp.WithInsecure(), // Use HTTP instead of HTTPS for local development
	}

	if config.OTLPEndpoint != "" {
		opts = append(opts, otlptracehttp.WithEndpoint(config.OTLPEndpoint))
	}

	exporter, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(config.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	// Set global propagator to tracecontext (W3C Trace Context)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	log.Info().Msg("OpenTelemetry tracer initialized successfully")

	// Return shutdown function
	return func(ctx context.Context) error {
		log.Info().Msg("Shutting down tracer provider...")
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer provider: %w", err)
		}
		log.Info().Msg("Tracer provider shutdown successfully")
		return nil
	}, nil
}
