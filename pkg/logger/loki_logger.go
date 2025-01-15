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

func (l *LokiLoggerImpl) GetZapLogger() *zap.Logger {
	return l.Logger
}

func (l *LokiLoggerImpl) WithCtx(ctx context.Context) Logger {
	var log *zap.Logger = l.Logger
	span := trace.SpanFromContext(ctx)

	// Extract the trace ID from the span's context
	spanContext := span.SpanContext()

	if spanContext.TraceID().IsValid() {
		traceID := spanContext.TraceID().String()
		traceId := zap.String("traceID", traceID)
		log = l.Logger.With(traceId)
	}

	return &LokiLoggerImpl{
		Logger: log,
		span:   span,
	}
}

func (l *LokiLoggerImpl) Error(msg string, fields ...zap.Field) {
	// Obtain caller information
	_, file, line, _ := runtime.Caller(1)

	callerInfo := fmt.Sprintf("%s:%d", file, line)

	eventOpt := trace.EventOption(trace.WithAttributes(attribute.String("caller", callerInfo)))
	if l.span != nil {
		l.span.RecordError(fmt.Errorf("%s", msg), eventOpt)
		l.span.SetStatus(codes.Error, msg)
	}

	fields = append(fields, zap.String("caller", fmt.Sprintf("%s:%d", file, line)))
	l.Logger.Error(msg, fields...)
}
