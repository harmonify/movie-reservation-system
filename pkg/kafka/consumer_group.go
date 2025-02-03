package kafka

import (
	"context"
	"errors"
	"strings"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/go-playground/validator/v10"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// KafkaConsumerGroup wraps a Sarama ConsumerGroup.
type KafkaConsumerGroup struct {
	Client sarama.ConsumerGroup

	logger logger.Logger
}

type KafkaConsumerGroupConfig struct {
	*KafkaConfig
	KafkaBrokers       string `validate:"required"`
	KafkaConsumerGroup string `validate:"required"`
}

// NewKafkaConsumerGroup initializes the Kafka consumer.
func NewKafkaConsumerGroup(lc fx.Lifecycle, cfg *KafkaConsumerGroupConfig, logger logger.Logger) (*KafkaConsumerGroup, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return nil, err
	}

	kafkaConfig, err := BuildKafkaConfig(cfg.KafkaConfig)
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewConsumerGroup(strings.Split(cfg.KafkaBrokers, ","), cfg.KafkaConsumerGroup, kafkaConfig)
	if err != nil {
		return nil, err
	}

	kc := &KafkaConsumerGroup{
		Client: client,
		logger: logger,
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Kafka consumer")
			return client.Close()
		},
	})

	return kc, nil
}

// StartConsumer starts consuming messages from the given topic.
// StartConsumer function do not need to be called inside a goroutine.
func (kc *KafkaConsumerGroup) StartConsumer(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) {
	go func() {
		wrappedHandler := otelsarama.WrapConsumerGroupHandler(handler)
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := kc.Client.Consume(ctx, topics, wrappedHandler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				kc.logger.WithCtx(ctx).Error("Consumer error", zap.Error(err))
			}
		}
	}()
}
