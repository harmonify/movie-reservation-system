package kafka_consumer

import (
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"go.uber.org/fx"
)

var (
	KafkaConsumerModule = fx.Module(
		"driver-kafka-consumer",
		fx.Provide(
			kafka.AsRoute(NewEmailVerificationRoute),
			kafka.AsRoute(NewSmsRoute),
			fx.Annotate(
				kafka.NewKafkaRouter,
				fx.ParamTags(`group:"kafka-routes"`),
			),
		),
	)
)
