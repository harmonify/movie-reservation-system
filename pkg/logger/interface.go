package logger

import (
	"context"
	"time"

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

type LoggerConfig struct {
	Level      string      // log level: debug, info, warn, error, fatal
	Type       string      // log encoding: console, loki
	LokiConfig *LokiConfig // optional for LokiLogger
}

type LokiConfig struct {
	// Url of the loki server including http:// or https://
	Url string
	// BatchMaxSize is the maximum number of log lines that are sent in one request
	BatchMaxSize int
	// BatchMaxWait is the maximum time to wait before sending a request
	BatchMaxWait time.Duration
	// Labels that are added to all log lines
	Labels   map[string]string
	Username string
	Password string
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
