package kafka

import (
	"context"
	"strings"

	"github.com/IBM/sarama"
	"go.uber.org/fx"

	"github.com/dnwe/otelsarama"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	// "github.com/harmonify/movie-reservation-system/pkg/tracer/carrier"
)

// KafkaProducer wraps a Sarama AsyncProducer.
type KafkaProducer struct {
	Client sarama.SyncProducer

	logger logger.Logger
	tracer tracer.Tracer
}

// NewKafkaProducer initializes the Kafka producer.
func NewKafkaProducer(lc fx.Lifecycle, cfg *config.Config, logger logger.Logger, tracer tracer.Tracer) (*KafkaProducer, error) {
	kafkaConfig, err := buildKafkaConfig(cfg)
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewSyncProducer(strings.Split(cfg.KafkaBrokers, ","), kafkaConfig)
	if err != nil {
		return nil, err
	}

	client = otelsarama.WrapSyncProducer(kafkaConfig, client)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Kafka producer")
			return client.Close()
		},
	})

	return &KafkaProducer{
		Client: client,
		logger: logger,
		tracer: tracer,
	}, nil
}

func (kp *KafkaProducer) SendMessage(ctx context.Context, msg *sarama.ProducerMessage) error {
	ctx, span := kp.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if msg == nil {
		return ErrNilMessage
	}

	if msg.Headers == nil {
		msg.Headers = []sarama.RecordHeader{}
	}

	kp.tracer.Inject(ctx, otelsarama.NewProducerMessageCarrier(msg))

	_, _, err := kp.Client.SendMessage(msg)
	return err
}
