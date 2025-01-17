package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type nopTracerImpl struct {
}

func NewNopTracer(p TracerParam) TracerResult {
	return TracerResult{
		Tracer: &nopTracerImpl{},
	}
}

func (t *nopTracerImpl) Start(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return ctx, trace.SpanFromContext(ctx)
}

func (t *nopTracerImpl) StartSpanWithCaller(ctx context.Context) (context.Context, trace.Span) {
	return ctx, trace.SpanFromContext(ctx)
}

func (t *nopTracerImpl) Shutdown(ctx context.Context) error {
	return nil
}
