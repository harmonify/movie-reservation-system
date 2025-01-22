package kafka

import (
	"context"
	"strings"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/fx"
)

// KafkaProducer wraps a sarama.SyncProducer and OTel instrumentation logic.
type KafkaProducer struct {
	Client sarama.SyncProducer

	logger logger.Logger
	tracer tracer.Tracer
}

type KafkaProducerParam struct {
	Lifecycle fx.Lifecycle
	Config    *config.Config
	Logger    logger.Logger
	Tracer    tracer.Tracer
}

// NewKafkaProducer initializes the Kafka producer.
func NewKafkaProducer(p KafkaProducerParam) (*KafkaProducer, error) {
	kafkaConfig, err := buildKafkaConfig(p.Config)
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewSyncProducer(strings.Split(p.Config.KafkaBrokers, ","), kafkaConfig)
	if err != nil {
		return nil, err
	}

	client = otelsarama.WrapSyncProducer(kafkaConfig, client)

	p.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Closing Kafka producer")
			return client.Close()
		},
	})

	return &KafkaProducer{
		Client: client,
		logger: p.Logger,
		tracer: p.Tracer,
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
