package kafka

import (
	"context"
	"log"
	"strings"

	"github.com/IBM/sarama"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// KafkaConsumerGroup wraps a Sarama ConsumerGroup.
type KafkaConsumerGroup struct {
	Client sarama.ConsumerGroup

	logger logger.Logger
}

// NewKafkaConsumerGroup initializes the Kafka consumer.
func NewKafkaConsumerGroup(lc fx.Lifecycle, cfg *config.Config, logger logger.Logger) (*KafkaConsumerGroup, error) {
	kafkaConfig, err := buildKafkaConfig(cfg)
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewConsumerGroup(strings.Split(cfg.KafkaBrokers, ","), cfg.KafkaConsumerGroup, kafkaConfig)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing Kafka consumer")
			return client.Close()
		},
	})

	return &KafkaConsumerGroup{
		Client: client,
		logger: logger,
	}, nil
}

// StartConsumer starts consuming messages from the given topic.
// StartConsumer function do not need to be called inside a goroutine.
func (kc *KafkaConsumerGroup) StartConsumer(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) {
	go func() {
		for {
			if err := kc.Client.Consume(ctx, topics, handler); err != nil {
				kc.logger.WithCtx(ctx).Error("Error consuming messages", zap.Error(err))
			}
		}
	}()
}
