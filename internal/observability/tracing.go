package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "victorialogs-mcp"
)

var (
	tracer trace.Tracer
)

// TracingConfig tracing configuration
type TracingConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServiceName string `mapstructure:"service_name"`
	Endpoint    string `mapstructure:"endpoint"`
	Sampler     string `mapstructure:"sampler"` // always, never, ratio
}

// InitTracing initializes OpenTelemetry tracing
func InitTracing(cfg TracingConfig) (func(context.Context) error, error) {
	if !cfg.Enabled {
		tracer = otel.Tracer(tracerName)
		return func(context.Context) error { return nil }, nil
	}

	serviceName := cfg.ServiceName
	if serviceName == "" {
		serviceName = "victorialogs-mcp"
	}

	// Create stdout exporter for development
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	// Create resource
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(getSampler(cfg.Sampler)),
	)

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer(tracerName)

	return tp.Shutdown, nil
}

// getSampler returns sampler based on config
func getSampler(sampler string) sdktrace.Sampler {
	switch sampler {
	case "always":
		return sdktrace.AlwaysSample()
	case "never":
		return sdktrace.NeverSample()
	case "ratio":
		return sdktrace.TraceIDRatioBased(0.1) // 10% sampling
	default:
		return sdktrace.AlwaysSample()
	}
}

// GetTracer returns the tracer
func GetTracer() trace.Tracer {
	if tracer == nil {
		tracer = otel.Tracer(tracerName)
	}
	return tracer
}

// StartSpan starts a new span
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return GetTracer().Start(ctx, name, opts...)
}

// SpanFromContext returns span from context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddToolAttributes adds tool-related attributes to span
func AddToolAttributes(span trace.Span, tool string, params map[string]interface{}) {
	span.SetAttributes(
		attribute.String("mcp.tool", tool),
	)

	if query, ok := params["query"].(string); ok {
		// Truncate query for safety
		if len(query) > 100 {
			query = query[:100] + "..."
		}
		span.SetAttributes(attribute.String("mcp.query_preview", query))
	}

	if limit, ok := params["limit"].(float64); ok {
		span.SetAttributes(attribute.Int("mcp.limit", int(limit)))
	}
}

// RecordError records error in span
func RecordError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
	}
}
