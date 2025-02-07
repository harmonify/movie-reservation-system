package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewConsoleLogger() Logger {
	// The bundled Config struct only supports the most common configuration
	// options. More complex needs, like splitting logs between multiple files
	// or writing to non-file outputs, require use of the zapcore package.
	//
	// In this example, imagine we're both sending our logs to Kafka and writing
	// them to the console. We'd like to encode the console output and the Kafka
	// topics differently, and we'd also like special treatment for
	// high-priority logs.

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the console output for human.
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core)
	defer logger.Sync()

	return &ConsoleLoggerImpl{logger, nil}
}

func (l *ConsoleLoggerImpl) GetZapLogger() *zap.Logger {
	return l.Logger
}

func (l *ConsoleLoggerImpl) With(fields ...zap.Field) Logger {
	return &ConsoleLoggerImpl{
		Logger: l.Logger.With(fields...),
		span:   l.span,
	}
}

func (l *ConsoleLoggerImpl) WithCtx(ctx context.Context) Logger {
	var log *zap.Logger = l.Logger
	span := trace.SpanFromContext(ctx)

	// Extract the trace ID from the span's context
	spanContext := span.SpanContext()

	if spanContext.TraceID().IsValid() {
		traceID := spanContext.TraceID().String()
		traceId := zap.String("traceID", traceID)
		log = l.Logger.With(traceId)
	}

	return &ConsoleLoggerImpl{
		Logger: log,
		span:   span,
	}
}

func (l *ConsoleLoggerImpl) Error(msg string, fields ...zap.Field) {
	// Obtain caller information
	_, file, line, _ := runtime.Caller(1)

	callerInfo := fmt.Sprintf("%s:%d", file, line)

	if l.span != nil {
		l.span.RecordError(fmt.Errorf("%s", msg), trace.WithAttributes(attribute.String("caller", callerInfo)))
		l.span.SetStatus(codes.Error, msg)
	}

	fields = append(fields, zap.String("caller", fmt.Sprintf("%s:%d", file, line)))
	l.Logger.Error(msg, fields...)
}
