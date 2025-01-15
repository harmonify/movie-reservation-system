package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

// type RecordHeader struct {
// 	Key   []byte
// 	Value []byte
// }

type Event struct {
	Headers   []*sarama.RecordHeader `json:"headers"`
	Timestamp time.Time              `json:"timestamp"`
	TraceID   string                 `json:"trace_id"`
	Key       string                 `json:"key"`
	Value     interface{}            `json:"value"`
	Topic     string                 `json:"topic"`
}
