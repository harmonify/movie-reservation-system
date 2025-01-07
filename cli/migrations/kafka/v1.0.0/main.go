package v_1_0_0

import (
	"context"
	"kafka-playground/shared"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type CreateNewOrderTopicMigration struct {
	client *kafka.AdminClient
	logger *log.Logger
}

func NewCreateNewOrderTopicMigration(client *kafka.AdminClient, logger *log.Logger) *CreateNewOrderTopicMigration {
	return &CreateNewOrderTopicMigration{
		client: client,
		logger: logger,
	}
}

func (m *CreateNewOrderTopicMigration) GetIdentifier() string {
	return "20250107234012_create-new-order-topic"
}

func (m *CreateNewOrderTopicMigration) Up(ctx context.Context) error {
	topics, err := m.client.DescribeTopics(
		ctx,
		kafka.NewTopicCollectionOfTopicNames([]string{shared.NewOrderTopic.String()}),
		kafka.SetAdminRequestTimeout(time.Second*15),
	)
	if err != nil {
		m.logger.Printf("Failed to describe topic: \"%s\"\n", shared.NewOrderTopic.String())
		return err
	}
	if len(topics.TopicDescriptions) < 1 {
		_, err := m.client.CreateTopics(
			ctx,
			[]kafka.TopicSpecification{
				{
					Topic:         shared.NewOrderTopic.String(),
					NumPartitions: 3,
				},
			},
			kafka.SetAdminRequestTimeout(time.Second*15),
			kafka.SetAdminOperationTimeout(time.Minute),
		)
		if err != nil {
			m.logger.Printf("Failed to create topic: \"%s\"\n", shared.NewOrderTopic.String())
			return err
		}
	}

	return nil
}

func (m *CreateNewOrderTopicMigration) Down(ctx context.Context) error {
	_, err := m.client.DeleteTopics(
		ctx,
		[]string{shared.NewOrderTopic.String()},
		kafka.SetAdminRequestTimeout(time.Second*15),
		kafka.SetAdminOperationTimeout(time.Minute),
	)
	if err != nil {
		m.logger.Printf("Failed to delete topic: \"%s\"\n", shared.NewOrderTopic.String())
		return err
	}
	return nil
}
