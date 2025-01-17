package tracer

import (
	"context"
	"fmt"
	"runtime"

	"github.com/harmonify/movie-reservation-system/pkg/config"
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

type jaegerTracerImpl struct {
	exporter *otlptrace.Exporter
	config   *config.Config
}

func NewJaegerTracer(p TracerParam) TracerResult {
	secureOption := otlptracegrpc.WithInsecure()

	exporter := otlptrace.NewUnstarted(
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(p.Config.OtelHost),
		),
	)

	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", p.Config.AppName),
			attribute.String("service.environment", p.Config.Env),
		),
	)
	if err != nil {
		fmt.Printf("Could not set tracer resources: %v\n", err)
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

	t := &jaegerTracerImpl{
		exporter: exporter,
		config:   p.Config,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := t.exporter.Start(ctx); err != nil {
				fmt.Printf("Failed to connect to Jaeger OpenTelemetry: %v\n", err)
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err = t.Shutdown(ctx); err != nil {
				fmt.Printf("Error shutting down tracer: %v\n", err)
				return err
			}
			return nil
		},
	})

	return TracerResult{
		Tracer: t,
	}
}

func (t *jaegerTracerImpl) Start(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.GetTracerProvider().Tracer(t.config.AppName).Start(ctx, spanName)
}

func (s *jaegerTracerImpl) StartSpanWithCaller(ctx context.Context) (context.Context, trace.Span) {
	pc, _, _, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()

	// segments := strings.Split(callerName, ".")
	// spanName := segments[len(segments)-1]

	ctx, span := s.Start(ctx, callerName)
	return ctx, span
}

func (t *jaegerTracerImpl) Shutdown(ctx context.Context) error {
	if err := t.exporter.Shutdown(ctx); err != nil {
		fmt.Println("Failed to shutdown OpenTelemetry exporter")
		return err
	}
	return nil
}
