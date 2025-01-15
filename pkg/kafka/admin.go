package kafka

import (
	"context"
	"strings"

	"github.com/IBM/sarama"
	"go.uber.org/fx"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
)

// KafkaAdmin wraps a Sarama ClusterAdmin.
type KafkaAdmin struct {
	Client sarama.ClusterAdmin

	logger logger.Logger
}

// NewKafkaAdmin initializes the Kafka admin.
func NewKafkaAdmin(lc fx.Lifecycle, cfg *config.Config, logger logger.Logger) (*KafkaAdmin, error) {
	kafkaConfig, err := buildKafkaConfig(cfg)
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
