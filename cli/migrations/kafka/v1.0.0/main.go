package v_1_0_0

import (
	"context"
	"log"
	"time"

	"github.com/harmonify/movie-reservation-system/cli/shared"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var NewOrderTopic = kafka.TopicSpecification{
	Topic:             shared.NewOrderTopic.String(),
	NumPartitions:     3,
	ReplicationFactor: 3,
	// Topic config reference: <https://kafka.apache.org/documentation/#topicconfigs>
	// Config: map[string]string{
	// 	"retention.ms": strconv.Itoa(7 * 24 * 60 * 60), // 7 days
	// },
}

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
		kafka.NewTopicCollectionOfTopicNames([]string{NewOrderTopic.Topic}),
		kafka.SetAdminRequestTimeout(time.Second*15),
	)
	if err != nil {
		m.logger.Printf("Failed to describe topic: \"%s\"\n", NewOrderTopic.Topic)
		return err
	}
	if len(topics.TopicDescriptions) < 1 {
		_, err := m.client.CreateTopics(
			ctx,
			[]kafka.TopicSpecification{NewOrderTopic},
			kafka.SetAdminRequestTimeout(time.Second*15),
			kafka.SetAdminOperationTimeout(time.Minute),
		)
		if err != nil {
			m.logger.Printf("Failed to create topic: \"%s\"\n", NewOrderTopic.Topic)
			return err
		}
	}

	return nil
}

func (m *CreateNewOrderTopicMigration) Down(ctx context.Context) error {
	_, err := m.client.DeleteTopics(
		ctx,
		[]string{NewOrderTopic.Topic},
		kafka.SetAdminRequestTimeout(time.Second*15),
		kafka.SetAdminOperationTimeout(time.Minute),
	)
	if err != nil {
		m.logger.Printf("Failed to delete topic: \"%s\"\n", NewOrderTopic.Topic)
		return err
	}
	return nil
}
