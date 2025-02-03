package logger

import (
	"strconv"
	"time"
)

type (
	LokiZapConfig struct {
		// Log level
		LogLevel string `validate:"required,oneof=debug info warn error"`
		// Url of the loki server including http:// or https://
		Url string `validate:"required,http_url"`
		// BatchMaxSize is the maximum number of log lines that are sent in one request
		BatchMaxSize int `validate:"required,min=1"`
		// BatchMaxWait is the maximum time to wait before sending a request
		BatchMaxWait time.Duration `validate:"required,min=1s"`
		// Labels that are added to all log lines
		Labels map[string]string
		// Basic auth username
		Username string
		// Basic auth password
		Password string
	}

	lokiPushRequest struct {
		Streams []stream `json:"streams"`
	}

	stream struct {
		Stream map[string]string `json:"stream"`
		Values []streamValue     `json:"values"`
	}

	streamValue [2]string

	logEntry struct {
		Level     string  `json:"level"`
		Timestamp float64 `json:"ts"`
		Message   string  `json:"msg"`
		Caller    string  `json:"caller"`
		TraceID   string  `json:"traceId"`
		Stack     string  `json:"stack"`
		raw       string
	}
)

func newLog(entry logEntry) streamValue {
	ts := time.Unix(int64(entry.Timestamp), 0)
	return streamValue{strconv.FormatInt(ts.UnixNano(), 10), entry.raw}
}
