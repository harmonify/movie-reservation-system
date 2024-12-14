package logger

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/user-servicelogger-utility/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LoggerModule = fx.Module("logger", fx.Provide(NewLogger))

type Logger interface {
	GetZapLogger() *zap.Logger
	WithCtx(ctx context.Context) Logger
	Error(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Log(debugLevel zapcore.Level, msg string, fields ...zap.Field)
}

type LoggerImpl struct {
	*zap.Logger
	span trace.Span
}

func NewLogger(cfg *config.Config) (Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.CallerKey = zapcore.OmitKey

	loki := logger.New(context.Background(), logger.Config{
		Url:          cfg.LokiURL,
		BatchMaxSize: 1000,
		BatchMaxWait: 10 * time.Second,
		Labels:       map[string]string{"app": cfg.AppName, "env": cfg.Env},
	})

	zapLogger, err := loki.WithCreateLogger(zapConfig)

	return &LoggerImpl{zapLogger, nil}, err
}

func (w *LoggerImpl) GetZapLogger() *zap.Logger {
	return w.Logger
}

func (w *LoggerImpl) WithCtx(ctx context.Context) Logger {
	var log *zap.Logger = w.Logger
	span := trace.SpanFromContext(ctx)

	// Extract the trace ID from the span's context
	spanContext := span.SpanContext()

	if spanContext.TraceID().IsValid() {
		traceID := spanContext.TraceID().String()
		traceId := zap.String("traceID", traceID)
		log = w.With(traceId)
	}

	return &LoggerImpl{
		Logger: log,
		span:   span,
	}
}

func (w *LoggerImpl) Error(msg string, fields ...zap.Field) {
	// Obtain caller information
	_, file, line, _ := runtime.Caller(1)

	callerInfo := fmt.Sprintf("%s:%d", file, line)

	eventOpt := trace.EventOption(trace.WithAttributes(attribute.String("caller", callerInfo)))
	if w.span != nil {
		w.span.RecordError(fmt.Errorf(msg), eventOpt)
		w.span.SetStatus(codes.Error, msg)
	}

	fields = append(fields, zap.String("caller", fmt.Sprintf("%s:%d", file, line)))
	w.Logger.Error(msg, fields...)
}

func (w *LoggerImpl) Info(msg string, fields ...zap.Field) {
	w.Logger.Info(msg, fields...)
}

func (w *LoggerImpl) Warn(msg string, fields ...zap.Field) {
	w.Logger.Warn(msg, fields...)
}

func (w *LoggerImpl) Log(debugLevel zapcore.Level, msg string, fields ...zap.Field) {
	w.Logger.Log(debugLevel, msg, fields...)
}
