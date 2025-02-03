package logger

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"runtime"

	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lokiSinkRegistered bool

func NewLokiZapLogger(cfg *LokiZapConfig) (Logger, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return nil, err
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.CallerKey = zapcore.OmitKey
	logLevel, err := zap.ParseAtomicLevel(cfg.LogLevel)
	if err == nil {
		zapConfig.Level = logLevel
	} else {
		fmt.Println("Failed to set log level")
	}

	if !lokiSinkRegistered {
		err := zap.RegisterSink(lokiSinkKey, func(_ *url.URL) (zap.Sink, error) {
			lp := newLokiPusher(cfg)
			sink := newLokiSink(lp)
			return sink, nil
		})
		if err != nil {
			log.Fatal(err)
		}
		lokiSinkRegistered = true
	}

	fullSinkKey := fmt.Sprintf("%s://", lokiSinkKey)
	if zapConfig.OutputPaths == nil {
		zapConfig.OutputPaths = []string{fullSinkKey}
	} else {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, fullSinkKey)
	}

	zapLogger, err := zapConfig.Build()

	return &LokiLoggerImpl{zapLogger, nil}, err
}

func (l *LokiLoggerImpl) GetZapLogger() *zap.Logger {
	return l.Logger
}

func (l *LokiLoggerImpl) With(fields ...zap.Field) Logger {
	return &LokiLoggerImpl{
		Logger: l.Logger.With(fields...),
		span:   l.span,
	}
}

func (l *LokiLoggerImpl) WithCtx(ctx context.Context) Logger {
	var logger *zap.Logger = l.Logger
	span := trace.SpanFromContext(ctx)

	// Extract the trace ID from the span's context
	spanContext := span.SpanContext()

	if spanContext.TraceID().IsValid() {
		traceId := zap.String("trace_id", spanContext.TraceID().String())
		logger = l.Logger.With(traceId)
	}

	return &LokiLoggerImpl{
		Logger: logger,
		span:   span,
	}
}

func (l *LokiLoggerImpl) Error(msg string, fields ...zap.Field) {
	// Obtain caller information
	_, file, line, _ := runtime.Caller(1)
	callerInfo := fmt.Sprintf("%s:%d", file, line)

	var err error
	for _, f := range fields {
		if f.Type == zapcore.ErrorType {
			errField, ok := f.Interface.(error)
			if ok {
				err = errField
			}
		}
	}
	if err == nil {
		err = fmt.Errorf("%s", msg)
	}

	if l.span != nil {
		l.span.RecordError(
			err,
			trace.WithAttributes(
				attribute.String("caller", callerInfo),
			),
			trace.WithStackTrace(true),
		)
		l.span.SetStatus(codes.Error, msg)
	}

	fields = append(fields, zap.String("caller", callerInfo))
	l.Logger.Error(msg, fields...)
}
