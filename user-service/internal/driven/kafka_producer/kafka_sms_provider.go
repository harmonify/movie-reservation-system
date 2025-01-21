package kafka_producer

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	notification_proto "github.com/harmonify/movie-reservation-system/user-service/internal/driven/proto/notification"
	"go.uber.org/fx"
)

type KafkaSmsProviderParam struct {
	fx.In
	KafkaProducer *kafka.KafkaProducer
	Tracer        tracer.Tracer
}

type KafkaSmsProviderResult struct {
	fx.Out
	SmsProvider shared.SmsProvider
}

type kafkaSmsProvider struct {
	client *kafka.KafkaProducer
	tracer tracer.Tracer
}

func NewKafkaSmsProvider(p KafkaSmsProviderParam) KafkaSmsProviderResult {
	return KafkaSmsProviderResult{
		SmsProvider: &kafkaSmsProvider{
			client: p.KafkaProducer,
			tracer: p.Tracer,
		},
	}
}

func (n *kafkaSmsProvider) Send(ctx context.Context, p *notification_proto.Sms) error {
	ctx, span := n.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if p.Recipient == "" {
		return shared.ErrEmptyRecipient
	}
	if p.Body == "" {
		return shared.ErrEmptyMessage
	}

	return n.client.SendMessage(ctx, &sarama.ProducerMessage{
		Topic: SmsTopicV1_0_0,
		Value: kafka.ProtoEncoder(p),
	})
}
