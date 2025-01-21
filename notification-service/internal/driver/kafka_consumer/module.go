package kafka_consumer

import (
	"context"
	"errors"

	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"go.uber.org/fx"
)

var (
	topics = []string{}

	KafkaConsumerModule = fx.Module(
		"driver-kafka-consumer",
		fx.Provide(
			kafka.NewKafkaConsumerGroup,
			fx.Annotate(
				kafka.NewKafkaRouter,
				fx.ParamTags(`group:"kafka-routes"`),
			),
			kafka.AsRoute(NewEmailVerificationRoute),
			kafka.AsRoute(NewSmsRoute),
		),
		fx.Invoke(BootstrapKafkaConsumer),
	)
)

func BootstrapKafkaConsumer(lc fx.Lifecycle, l logger.Logger, cg *kafka.KafkaConsumerGroup, r kafka.KafkaRouter) {
	lc.Append(fx.StartHook(func(ctx context.Context) error {
		cg.StartConsumer(ctx, shared.RegisteredTopics, r)
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
