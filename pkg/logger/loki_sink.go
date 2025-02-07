package logger

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
)

const lokiSinkKey = "loki"

type lokiSink zap.Sink

func newLokiSink(lp *lokiPusherImpl) lokiSink {
	return sink{
		lokiPusher: lp,
	}
}

type sink struct {
	lokiPusher *lokiPusherImpl
}

func (s sink) Sync() error {
	// No need to check the batch in sink; batching is handled locally in lokiPusherImpl
	return nil
}
func (s sink) Close() error {
	s.lokiPusher.stop()
	return nil
}

func (s sink) Write(p []byte) (int, error) {
	var entry logEntry
	err := json.Unmarshal(p, &entry)
	if err != nil {
		return 0, err
	}
	entry.raw = string(p)

	retries := 3
	for retries > 0 {
		select {
		case s.lokiPusher.entry <- entry:
			return len(p), nil
		default:
			retries--
			time.Sleep(20 * time.Millisecond) // Small backoff
		}
	}

	return 0, fmt.Errorf("failed to write log after retries: channel is full")
}
