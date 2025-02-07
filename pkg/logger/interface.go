package logger

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	GetZapLogger() *zap.Logger
	With(fields ...zap.Field) Logger
	WithCtx(ctx context.Context) Logger
	Level() zapcore.Level
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Log(debugLevel zapcore.Level, msg string, fields ...zap.Field)
}

type NopLoggerImpl struct {
	*zap.Logger
	span trace.Span
}

type ConsoleLoggerImpl struct {
	*zap.Logger
	span trace.Span
}

type LokiLoggerImpl struct {
	*zap.Logger
	span trace.Span
}
