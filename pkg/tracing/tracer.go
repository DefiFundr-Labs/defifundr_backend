package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TracerProvider holds the OpenTelemetry TracerProvider
type TracerProvider struct {
	provider *sdktrace.TracerProvider
}

// Config holds the configuration for the OpenTelemetry tracer
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	// OTLPEndpoint is the endpoint for the OpenTelemetry collector
	// e.g., "otel-collector:4317" for gRPC or "otel-collector:4318" for HTTP
	OTLPEndpoint string
	// UseStdoutExporter determines whether to use the stdout exporter (useful for local development)
	UseStdoutExporter bool
	// DisableTracing can be used to completely disable tracing
	DisableTracing bool
	// SampleRate is the fraction of traces to sample (1.0 = 100%, 0.1 = 10%)
	SampleRate float64
}

// DefaultConfig returns a default configuration for the tracer
func DefaultConfig() Config {
	return Config{
		ServiceName:       "defifundr-service",
		ServiceVersion:    "0.1.0",
		Environment:       "development",
		OTLPEndpoint:      "localhost:4317",
		UseStdoutExporter: false,
		DisableTracing:    false,
		SampleRate:        1.0, // Sample 100% of traces by default
	}
}

// NewTracerProvider creates a new TracerProvider with the given configuration
func NewTracerProvider(ctx context.Context, cfg Config) (*TracerProvider, error) {
	if cfg.DisableTracing {
		// Return a no-op tracer provider
		return &TracerProvider{provider: sdktrace.NewTracerProvider()}, nil
	}

	// Create a resource describing the service
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	// Configure the trace exporter(s)
	var exporter sdktrace.SpanExporter
	var exporterErr error

	if cfg.UseStdoutExporter {
		// Use stdout exporter for local development
		exporter, exporterErr = stdouttrace.New(stdouttrace.WithPrettyPrint())
	} else {
		// Use OTLP exporter for production
		conn, err := grpc.NewClient(cfg.OTLPEndpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}

		exporter, exporterErr = otlptrace.New(ctx, otlptracegrpc.NewClient(
			otlptracegrpc.WithGRPCConn(conn),
		))
	}

	if exporterErr != nil {
		return nil, exporterErr
	}

	// Create the trace provider with the exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter,
			sdktrace.WithMaxExportBatchSize(512),
			sdktrace.WithBatchTimeout(5*time.Second),
		),
		sdktrace.WithSampler(sdktrace.ParentBased(
			sdktrace.TraceIDRatioBased(cfg.SampleRate),
		)),
	)

	// Set the global trace provider and propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)),
	))

	return &TracerProvider{provider: tp}, nil
}

// Tracer returns a named tracer from the global provider
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// Shutdown shuts down the tracer provider
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	if tp.provider == nil {
		return nil
	}
	return tp.provider.Shutdown(ctx)
}

// SetupOTel initializes the OpenTelemetry SDK with the given configuration
// It returns a shutdown function that should be called when the application exits
func SetupOTel(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	tp, err := NewTracerProvider(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Return a shutdown function
	shutdown := func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	}

	return shutdown, nil
}

// MustSetupOTel initializes the OpenTelemetry SDK with the given configuration
// It panics if there is an error
func MustSetupOTel(ctx context.Context, cfg Config) func(context.Context) error {
	shutdown, err := SetupOTel(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return shutdown
}

// WithSpan creates a new span and context, then executes the given function
// This is a convenience function for creating spans
func WithSpan(ctx context.Context, name string, fn func(context.Context) error) error {
	tracer := otel.Tracer("")
	ctx, span := tracer.Start(ctx, name)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

// RecordError records an error on the current span
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
}
