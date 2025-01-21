package kafka_producer

import (
	"go.uber.org/fx"
)

var (
	DrivenKafkaModule = fx.Module(
		"driven-kafka",
		fx.Provide(
			NewKafkaEmailProvider,
			NewKafkaSmsProvider,
		),
	)
)
