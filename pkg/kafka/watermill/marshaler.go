package watermill_pkg

import (
	"context"

	"github.com/IBM/sarama"
	watermill_kafka "github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/dnwe/otelsarama"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

type WatermillMarshaler struct {
	watermill_kafka.DefaultMarshaler
}

func NewWatermillMarshaler() watermill_kafka.MarshalerUnmarshaler {
	return &WatermillMarshaler{}
}

func (j *WatermillMarshaler) Marshal(topic string, msg *message.Message) (*sarama.ProducerMessage, error) {
	kafkaMsg, err := j.DefaultMarshaler.Marshal(topic, msg)
	if err != nil {
		return nil, err
	}

	key, err := j.GeneratePartitionKey(topic, msg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot generate partition key")
	}
	if key != "" {
		kafkaMsg.Key = sarama.ByteEncoder(key)
	}

	// inject tracing span context
	ctx := msg.Context()
	if ctx != nil {
		otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(kafkaMsg))
	}

	return kafkaMsg, nil
}

func (*WatermillMarshaler) Unmarshal(kafkaMsg *sarama.ConsumerMessage) (*message.Message, error) {
	var messageID string
	metadata := make(message.Metadata, len(kafkaMsg.Headers))

	for _, header := range kafkaMsg.Headers {
		if string(header.Key) == watermill_kafka.UUIDHeaderKey {
			messageID = string(header.Value)
		} else {
			metadata.Set(string(header.Key), string(header.Value))
		}
	}

	msg := message.NewMessage(messageID, kafkaMsg.Value)
	msg.Metadata = metadata

	// extract tracing span context
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(kafkaMsg))
	msg.SetContext(ctx)

	return msg, nil
}

const PartitionKeyMetadataKey = "_watermill_partition_key"

func (*WatermillMarshaler) GeneratePartitionKey(topic string, msg *message.Message) (string, error) {
	return msg.Metadata.Get(PartitionKeyMetadataKey), nil
}
