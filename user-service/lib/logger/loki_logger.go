package logger

import (
	"context"
	"fmt"
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func NewLokiLogger(zapConfig zap.Config, lokiConfig LokiConfig) (Logger, error) {
	loki := NewZapLoki(context.Background(), lokiConfig)
	zapLogger, err := loki.WithCreateLogger(zapConfig)
	return &LokiLoggerImpl{zapLogger, nil}, err
}

func (c *LokiLoggerImpl) GetZapLogger() *zap.Logger {
	return c.Logger
}

func (w *LokiLoggerImpl) WithCtx(ctx context.Context) Logger {
	var log *zap.Logger = w.Logger
	span := trace.SpanFromContext(ctx)

	// Extract the trace ID from the span's context
	spanContext := span.SpanContext()

	if spanContext.TraceID().IsValid() {
		traceID := spanContext.TraceID().String()
		traceId := zap.String("traceID", traceID)
		log = w.Logger.With(traceId)
	}

	return &LokiLoggerImpl{
		Logger: log,
		span:   span,
	}
}

func (w *LokiLoggerImpl) Error(msg string, fields ...zap.Field) {
	// Obtain caller information
	_, file, line, _ := runtime.Caller(1)

	callerInfo := fmt.Sprintf("%s:%d", file, line)

	eventOpt := trace.EventOption(trace.WithAttributes(attribute.String("caller", callerInfo)))
	if w.span != nil {
		w.span.RecordError(fmt.Errorf("%s", msg), eventOpt)
		w.span.SetStatus(codes.Error, msg)
	}

	fields = append(fields, zap.String("caller", fmt.Sprintf("%s:%d", file, line)))
	w.Logger.Error(msg, fields...)
}
