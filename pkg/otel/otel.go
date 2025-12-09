package otel

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Config holds the configuration for OpenTelemetry setup.
type Config struct {
	// ServiceName is the name of the service being instrumented.
	ServiceName string
	// ServiceVersion is the version of the service.
	ServiceVersion string
	// Environment is the deployment environment (e.g., "development", "production").
	Environment string
	// OTLPEndpoint is the endpoint for the OTLP collector (e.g., "localhost:4318").
	OTLPEndpoint string
	// EnableTracing enables trace collection.
	EnableTracing bool
	// JaegerEndpoint is the endpoint for the Jaeger collector (e.g., "localhost:14268").
	JaegerEndpoint string
	// EnableMetrics enables metric collection.
	EnableMetrics bool
	// EnableLogging enables log collection.
	EnableLogging bool
	// Insecure disables TLS for OTLP connections.
	Insecure bool
}

// DefaultConfig returns a default configuration with common settings.
func DefaultConfig(serviceName string) Config {
	return Config{
		ServiceName:    serviceName,
		ServiceVersion: "1.0.0",
		Environment:    "development",
		OTLPEndpoint:   "",
		EnableTracing:  true,
		JaegerEndpoint: "jaeger:4317",
		EnableMetrics:  false,
		EnableLogging:  false,
		Insecure:       true,
	}
}

// Provider holds the OpenTelemetry providers for graceful shutdown.
type Provider struct {
	tracerProvider *sdktrace.TracerProvider
}

// Shutdown gracefully shuts down all providers.
func (p *Provider) Shutdown(ctx context.Context) error {
	if p.tracerProvider != nil {
		return p.tracerProvider.Shutdown(ctx)
	}

	return nil
}

// Setup initializes OpenTelemetry with the given configuration.
// It returns a Provider that should be used for graceful shutdown.
func Setup(ctx context.Context, cfg Config) (*Provider, error) {
	provider := &Provider{}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
		resource.WithHost(),
		resource.WithProcess(),
		resource.WithOS(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Set up propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Setup tracing
	if cfg.EnableTracing {
		tracerProvider, err := setupTracingWithGrpcJaeger(ctx, cfg, res)
		if err != nil {
			return nil, fmt.Errorf("failed to setup tracing: %w", err)
		}
		provider.tracerProvider = tracerProvider
		otel.SetTracerProvider(tracerProvider)
	}

	return provider, nil
}

func setupTracingWithGrpcJaeger(ctx context.Context, cfg Config, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.JaegerEndpoint),
	}

	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5*time.Second),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	return tracerProvider, nil
}
