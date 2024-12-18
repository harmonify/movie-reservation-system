package logger

import (
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

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

type ConsoleLoggerImpl struct {
	*zap.Logger
	span trace.Span
}

type LokiLoggerImpl struct {
	*zap.Logger
	span trace.Span
}
