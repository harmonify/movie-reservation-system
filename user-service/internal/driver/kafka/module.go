package kafka_driver

import (
	"context"
	"errors"

	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"go.uber.org/fx"
)

var (
	registeredTopics = []string{
		shared.PublicUserRegisteredV1.String(),
	}

	KafkaConsumerModule = fx.Module(
		"kafka-driver",
		fx.Provide(
			kafka.NewKafkaConsumerGroup,
			kafka.NewKafkaDLQProducer,
			fx.Annotate(
				kafka.NewKafkaRouter,
				fx.ParamTags(`group:"kafka-routes"`),
			),
			kafka.AsRoute(NewOutboxRoute),
		),
		fx.Invoke(BootstrapKafkaConsumer),
	)
)

func BootstrapKafkaConsumer(lc fx.Lifecycle, l logger.Logger, cg *kafka.KafkaConsumerGroup, r kafka.KafkaRouter) {
	lc.Append(fx.StartHook(func(ctx context.Context) error {
		cg.StartConsumer(ctx, registeredTopics, r)
		select {
		case <-r.Ready():
			l.WithCtx(ctx).Info("Sarama consumer up and running!")
			return nil
		case <-ctx.Done():
			err := errors.New("consumer failed to become ready within the timeout")
			l.WithCtx(ctx).Error(err.Error())
			return err
		}
	}))
}
