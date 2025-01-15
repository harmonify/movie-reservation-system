package kafka_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/kafka"
	"github.com/harmonify/movie-reservation-system/user-service/lib/kafka/test"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

var (
	testTopic = "test-topic"
)

func TestKafkaSuite(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(KafkaTestSuite))
}

type KafkaTestSuite struct {
	suite.Suite

	logger        logger.Logger
	app           *fx.App
	producer      *kafka.KafkaProducer
	consumerGroup *kafka.KafkaConsumerGroup
	testConsumer  *test.TestConsumer
}

func (s *KafkaTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.Provide(
			func() *config.Config {
				return &config.Config{
					KafkaBrokers:       "localhost:9092",
					KafkaVersion:       "3.9.0",
					KafkaConsumerGroup: "user-service",
				}
			},
			func() *sync.WaitGroup {
				var wg sync.WaitGroup
				return &wg
			},
			logger.NewConsoleLogger,
			kafka.NewKafkaProducer,
			kafka.NewKafkaConsumerGroup,
			test.NewTestConsumer,
		),
		fx.Invoke(func(logger logger.Logger, p *kafka.KafkaProducer, cg *kafka.KafkaConsumerGroup, c *test.TestConsumer) {
			s.logger = logger
			s.producer = p
			s.consumerGroup = cg
			s.testConsumer = c
		}),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *KafkaTestSuite) TestKafka() {
	// Set up a timeout context for consumer to be ready
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Start the consumer in a separate goroutine
	s.consumerGroup.StartConsumer(ctx, []string{testTopic}, s.testConsumer)

	// Wait for the consumer to be ready
	select {
	case <-s.testConsumer.Ready():
		s.logger.Info("Sarama consumer up and running!")
	case <-ctx.Done():
		s.T().Fatal("Consumer failed to become ready within the timeout")
	}

	// Send a message to the test topic
	expectedKey := []byte("test-key")
	expectedValue := []byte("{\"hello\":\"world\"}")

	err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: testTopic,
		Key:   sarama.ByteEncoder(expectedKey),
		Value: sarama.ByteEncoder(expectedValue),
	})
	s.Require().Nil(err, "Producer should send message successfully")

	// Set up a timeout context waiting for consumer to consume a message
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	select {
	case message := <-s.testConsumer.Messages():
		s.logger.Info(fmt.Sprintf("Received message: %s", string(message.Value)))
		s.Require().Equal(testTopic, message.Topic, "Consumer should receive the message from the correct topic")
		s.Require().Equal(expectedKey, message.Key, "Consumer should receive the message with the correct key")
		s.Require().Equal(expectedValue, message.Value, "Consumer should receive the correct message")
	case <-ctx.Done():
		s.T().Fatal("Test timed out waiting for the message to be processed")
	}
}
