package tracer

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

var (
	TracerModule = fx.Module("tracer", fx.Provide(NewJaegerTracer))
)

// Tracer provides a set of methods for working with distributed tracing using OpenTelemetry.
// It allows you to start spans, inject and extract trace context to manage propagation across services.
type Tracer interface {
	// Start creates a span and a context.Context containing the newly-created span.
	Start(ctx context.Context, spanName string) (context.Context, trace.Span)
	// Starts a span with the name of the caller function.
	StartSpanWithCaller(ctx context.Context) (context.Context, trace.Span)
	// Injects the current trace context into the provided carrier (e.g., Kafka headers, HTTP request headers).
	Inject(ctx context.Context, carrier propagation.TextMapCarrier)
	// Extracts the trace context from the provided carrier and returns a new context with the extracted trace information.
	Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context
}

type TracerParam struct {
	fx.In

	Config    *config.Config
	Lifecycle fx.Lifecycle
}

type TracerResult struct {
	fx.Out

	Tracer Tracer
}
