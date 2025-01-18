package carrier

import (
	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel/propagation"
)

// KafkaCarrier wraps the Kafka headers to implement the TextMapCarrier interface
type KafkaCarrier []*sarama.RecordHeader

func (kc KafkaCarrier) Get(key string) string {
	for _, h := range kc {
		if string(h.Key) == key {
			return string(h.Value)
		}
	}
	return ""
}

func (kc KafkaCarrier) Set(key, value string) {
	kc = append(kc, &sarama.RecordHeader{Key: []byte(key), Value: []byte(value)})
}

func (kc KafkaCarrier) Keys() []string {
	keys := make([]string, 0, len(kc))
	for _, h := range kc {
		keys = append(keys, string(h.Key))
	}
	return keys
}

var _ propagation.TextMapCarrier = KafkaCarrier{}
