package kafkaconsumer

import (
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"go.uber.org/fx"
)

var (
	KafkaConsumerModule = fx.Module(
		"driver-kafka-consumer",
		fx.Provide(
			kafka.AsRoute(NewEmailRoute),
			kafka.AsRoute(NewSmsRoute),
		),
	)
)
