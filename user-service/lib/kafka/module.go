package kafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"go.uber.org/fx"
)

var KafkaModule = fx.Module(
	"kafka",
	fx.Provide(
		NewKafkaConsumerGroup,
		NewKafkaProducer,
	),
)

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
