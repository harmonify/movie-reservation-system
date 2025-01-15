package kafka

import (
	"context"
	"strings"

	"github.com/IBM/sarama"
	"go.uber.org/fx"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
)

// KafkaProducer wraps a Sarama AsyncProducer.
type KafkaProducer struct {
	Client sarama.SyncProducer

	logger logger.Logger
}

// ProvideKafkaProducer initializes the Kafka producer.
func NewKafkaProducer(lc fx.Lifecycle, cfg *config.Config, logger logger.Logger) (*KafkaProducer, error) {
	kafkaConfig, err := buildKafkaConfig(cfg)
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewSyncProducer(strings.Split(cfg.KafkaBrokers, ","), kafkaConfig)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Kafka producer")
			return client.Close()
		},
	})

	return &KafkaProducer{
		Client: client,
		logger: logger,
	}, nil
}

func (kp *KafkaProducer) SendMessage(msg *sarama.ProducerMessage) error {
	_, _, err := kp.Client.SendMessage(msg)
	return err
}
