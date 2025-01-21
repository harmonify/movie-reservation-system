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

type KafkaEmailProviderParam struct {
	fx.In
	KafkaProducer *kafka.KafkaProducer
	Tracer        tracer.Tracer
}

type KafkaEmailProviderResult struct {
	fx.Out
	EmailProvider shared.EmailProvider
}

type kafkaEmailProvider struct {
	client *kafka.KafkaProducer
	tracer tracer.Tracer
}

func NewKafkaEmailProvider(p KafkaEmailProviderParam) KafkaEmailProviderResult {
	return KafkaEmailProviderResult{
		EmailProvider: &kafkaEmailProvider{
			client: p.KafkaProducer,
			tracer: p.Tracer,
		},
	}
}

func (n *kafkaEmailProvider) Send(ctx context.Context, p *notification_proto.Email) error {
	ctx, span := n.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if len(p.Recipients) <= 0 {
		return shared.ErrEmptyRecipient
	}
	if p.Subject == "" {
		return shared.ErrEmptySubject
	}
	if p.TemplateId == "" {
		return shared.ErrEmptyTemplate
	}

	return n.client.SendMessage(ctx, &sarama.ProducerMessage{
		Topic: EmailTopicV1_0_0,
		Value: kafka.ProtoEncoder(p),
	})
}
