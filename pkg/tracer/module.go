package tracer

import (
	"context"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
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

type TracerConfig struct {
	Env               string `validate:"required,oneof=dev test prod"`
	ServiceIdentifier string `validate:"required"`
	Type              string `validate:"required,oneof=jaeger console nop"`
	OtelEndpoint      string `validate:"required_if=Type jaeger"`
}

func NewTracer(lc fx.Lifecycle, cfg *TracerConfig) (Tracer, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return nil, err
	}

	switch cfg.Type {
	case "jaeger":
		return NewJaegerTracer(cfg, lc)
	case "console":
		return NewConsoleTracer(cfg)
	default:
		return NewNopTracer(cfg), nil
	}
}
