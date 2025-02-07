package kafka_migration

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/harmonify/movie-reservation-system/cli/shared"
)

var PublicUserRegisteredV1Topic = kafka.TopicSpecification{
	Topic:             shared.PublicUserRegisteredV1Topic.String(),
	NumPartitions:     1,
	ReplicationFactor: 1,
	// Topic config reference: <https://kafka.apache.org/documentation/#topicconfigs>
	Config: map[string]string{
		"retention.ms": strconv.Itoa(-1),
	},
}

type CreatePublicUserRegisteredV1TopicMigration struct {
	client *kafka.AdminClient
}

func NewCreatePublicUserRegisteredV1TopicMigration(client *kafka.AdminClient) *CreatePublicUserRegisteredV1TopicMigration {
	return &CreatePublicUserRegisteredV1TopicMigration{
		client: client,
	}
}

func (m *CreatePublicUserRegisteredV1TopicMigration) GetIdentifier() string {
	return "20250126074050_create_public.user.registered.v1_topic"
}

func (m *CreatePublicUserRegisteredV1TopicMigration) Up(ctx context.Context) error {
	describeTopicResult, err := m.client.DescribeTopics(
		ctx,
		kafka.NewTopicCollectionOfTopicNames([]string{PublicUserRegisteredV1Topic.Topic}),
		kafka.SetAdminRequestTimeout(time.Second*15),
	)
	if err != nil {
		return err
	}

	topic := describeTopicResult.TopicDescriptions[0]
	if topic.Error.Code() == kafka.ErrNoError {
		err := fmt.Errorf("topic \"%s\" already exists", PublicUserRegisteredV1Topic.Topic)
		return err
	}
	if topic.Error.Code() != kafka.ErrUnknownTopicOrPart {
		return topic.Error
	}

	createTopicResult, err := m.client.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{PublicUserRegisteredV1Topic},
		kafka.SetAdminRequestTimeout(time.Second*15),
		kafka.SetAdminOperationTimeout(time.Minute),
	)
	if err != nil {
		return err
	}
	if createTopicResult[0].Error.Code() != kafka.ErrNoError {
		return createTopicResult[0].Error
	}

	return nil
}

func (m *CreatePublicUserRegisteredV1TopicMigration) Down(ctx context.Context) error {
	_, err := m.client.DeleteTopics(
		ctx,
		[]string{PublicUserRegisteredV1Topic.Topic},
		kafka.SetAdminRequestTimeout(time.Second*15),
		kafka.SetAdminOperationTimeout(time.Minute),
	)
	if err != nil {
		return err
	}
	return nil
}
