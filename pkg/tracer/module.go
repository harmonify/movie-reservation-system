package tracer

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

var (
	TracerModule = fx.Module("tracer", fx.Provide(NewJaegerTracer))
)

type Tracer interface {
	Start(ctx context.Context, spanName string) (context.Context, trace.Span)
	StartSpanWithCaller(ctx context.Context) (context.Context, trace.Span)
	Shutdown(ctx context.Context) error
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
