package shared

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/fx"
)

func NewKafkaAdminClient(cfg *Config, lc fx.Lifecycle) (*kafka.AdminClient, error) {
	client, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": cfg.KafkaServerUrl,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			client.Close()
			return nil
		},
	})

	return client, nil
}
