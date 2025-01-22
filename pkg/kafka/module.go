package kafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	"go.uber.org/fx"
)

var KafkaModule = fx.Module(
	"kafka",
	fx.Provide(
		NewKafkaAdmin,
		NewKafkaProducer,
		NewKafkaConsumerGroup,
		NewKafkaDLQProducer,
		// Example register kafka route:
		// AsRoute(
		// test.NewTestRoute,
		// ),
		fx.Annotate(
			NewKafkaRouter,
			fx.ParamTags(`group:"kafka-routes"`),
		),
	),
)

func AsRoute(f any, anns ...fx.Annotation) any {
	finalAnns := []fx.Annotation{
		fx.As(new(Route)),
		fx.ResultTags(`group:"kafka-routes"`),
	}
	if len(anns) > 0 {
		finalAnns = append(finalAnns, anns...)
	}

	return fx.Annotate(
		f,
		finalAnns...,
	)
}

func buildKafkaConfig(cfg *config.Config) (*sarama.Config, error) {
	version, err := sarama.ParseKafkaVersion(cfg.KafkaVersion)
	if err != nil {
		return nil, err
	}

	c := sarama.NewConfig()

	c.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	c.Consumer.Offsets.Initial = sarama.OffsetOldest

	c.Producer.Compression = sarama.CompressionLZ4 // very fast and reasonable compression ratio
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Retry.Max = 5
	c.Producer.Retry.Backoff = 100 * time.Millisecond
	c.Producer.Return.Successes = true

	c.Version = version

	return c, nil
}
