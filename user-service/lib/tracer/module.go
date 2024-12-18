package tracer

import (
	"context"
	"fmt"
	"runtime"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

var (
	TracerModule = fx.Module("tracer", fx.Provide(InitTracer))
)

type Tracer interface {
	Start(ctx context.Context, spanName string) (context.Context, trace.Span)
	StartSpanWithCaller(ctx context.Context) (context.Context, trace.Span)
	Shutdown(ctx context.Context) error
}

type TracerImpl struct {
	Exporter *otlptrace.Exporter
	Config   *config.Config
}

func InitTracer(cfg *config.Config) Tracer {
	secureOption := otlptracegrpc.WithInsecure()

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(cfg.OtelHost),
		),
	)

	if err != nil {
		fmt.Println("Failed to connect to Jaeger Open Telemetry", err.Error())
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.AppName),
			attribute.String("service.environment", cfg.Env),
		),
	)
	if err != nil {
		fmt.Println("Could not set resources", err.Error())
	}

	tracer := sdkTrace.NewTracerProvider(
		sdkTrace.WithSampler(sdkTrace.AlwaysSample()),
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(resources),
	)

	otel.SetTracerProvider(tracer)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Set an error handler for trace exports
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		fmt.Println("OpenTelemetry trace export error", err.Error())
	}))

	return &TracerImpl{
		Exporter: exporter,
		Config:   cfg,
	}
}

func (t *TracerImpl) Start(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.GetTracerProvider().Tracer(t.Config.AppName).Start(ctx, spanName)
}

func (s *TracerImpl) StartSpanWithCaller(ctx context.Context) (context.Context, trace.Span) {
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()

	// segments := strings.Split(callerName, ".")
	// spanName := segments[len(segments)-1]

	ctx, span := s.Start(ctx, callerName)
	return ctx, span
}

func (t *TracerImpl) Shutdown(ctx context.Context) error {
	if err := t.Exporter.Shutdown(ctx); err != nil {
		fmt.Println("Failed to shutdown OpenTelemetry exporter")
		return err
	}
	return nil
}
