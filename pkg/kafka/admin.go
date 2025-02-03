package kafka

import (
	"context"
	"strings"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"go.uber.org/fx"
)

// KafkaAdmin wraps a Sarama ClusterAdmin.
type KafkaAdmin struct {
	Client sarama.ClusterAdmin

	logger logger.Logger
}

type KafkaAdminConfig struct {
	*KafkaConfig
	KafkaBrokers string `validate:"required"`
}

// NewKafkaAdmin initializes the Kafka admin.
func NewKafkaAdmin(lc fx.Lifecycle, cfg *KafkaAdminConfig, logger logger.Logger) (*KafkaAdmin, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return nil, err
	}

	kafkaConfig, err := BuildKafkaConfig(cfg.KafkaConfig)
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewClusterAdmin(strings.Split(cfg.KafkaBrokers, ","), kafkaConfig)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Kafka admin")
			return client.Close()
		},
	})

	return &KafkaAdmin{
		Client: client,
		logger: logger,
	}, nil
}
