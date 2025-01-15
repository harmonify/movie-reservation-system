package kafka_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/kafka/test"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	test_proto "github.com/harmonify/movie-reservation-system/pkg/test/proto"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
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
	admin         *kafka.KafkaAdmin
	testConsumer  *test.TestConsumer
	router        kafka.KafkaRouter
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
			kafka.NewKafkaAdmin,
			test.NewTestConsumer,
			kafka.AsRoute(
				test.NewTestRoute,
			),
			fx.Annotate(
				kafka.NewKafkaRouter,
				fx.ParamTags(`group:"kafka-routes"`),
			),
		),
		fx.Invoke(func(logger logger.Logger, p *kafka.KafkaProducer, cg *kafka.KafkaConsumerGroup, a *kafka.KafkaAdmin, c *test.TestConsumer, r kafka.KafkaRouter) {
			s.logger = logger
			s.producer = p
			s.consumerGroup = cg
			s.admin = a
			s.testConsumer = c
			s.router = r
		}),

		fx.NopLogger,
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *KafkaTestSuite) SetupTest() {
	s.admin.Client.DeleteTopic(test.TestBasicTopic)
	// s.Require().Nil(err, "Admin should successfully delete test topic for setup process")
}

func (s *KafkaTestSuite) TearDownTest() {
	s.admin.Client.DeleteTopic(test.TestBasicTopic)
	// s.Require().Nil(err, "Admin should successfully delete test topic for teardown process")
}

func (s *KafkaTestSuite) TestKafkaSuite_Producer_And_Consumer() {
	// ARRANGE
	// Set up a timeout context for consumer to be ready
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Start the consumer in a separate goroutine
	s.consumerGroup.StartConsumer(ctx, []string{test.TestBasicTopic}, s.testConsumer)
	// Wait for the consumer to be ready
	select {
	case <-s.testConsumer.Ready():
		s.logger.Info("Sarama consumer up and running!")
	case <-ctx.Done():
		s.T().Fatal("Consumer failed to become ready within the timeout")
	}
	// Construct message key and value
	expectedKey := []byte("test-key")
	expectedValue := &test_proto.Test{
		TraceId: uuid.New().String(),
		Message: "hello world",
	}

	// ACT
	// Send a message to the test topic
	err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: test.TestBasicTopic,
		Key:   sarama.ByteEncoder(expectedKey),
		Value: kafka.ProtoEncoder(expectedValue),
	})
	s.Require().Nil(err, "Producer should send message successfully")
	// Set up a timeout context waiting for consumer to consume a message
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Consume message
	select {
	case message := <-s.testConsumer.Messages():
		// ASSERT
		s.Require().Equal(test.TestBasicTopic, message.Topic, "Consumer should receive the message from the correct topic")
		s.Require().Equal(expectedKey, message.Key, "Consumer should receive the message with the correct key")

		val := &test_proto.Test{}
		err := proto.Unmarshal(message.Value, val)
		s.Require().Nil(err, "Consumer should successfully unmarshal the message")
		s.Require().Equal(expectedValue.TraceId, val.GetTraceId(), "Consumer should receive the correct trace id")
		s.Require().Equal(expectedValue.Message, val.GetMessage(), "Consumer should receive the correct message")
	case <-ctx.Done():
		s.T().Fatal("Test timed out waiting for the message to be processed")
	}
}

func (s *KafkaTestSuite) TestKafkaSuite_KafkaRouter() {
	// ARRANGE
	// Set up a timeout context for consumer to be ready
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Start the consumer in a separate goroutine
	s.consumerGroup.StartConsumer(ctx, []string{test.TestBasicTopic}, s.router)
	// Wait for the consumer to be ready
	select {
	case <-s.router.Ready():
		s.logger.Info("Sarama consumer up and running!")
	case <-ctx.Done():
		s.T().Fatal("Consumer failed to become ready within the timeout")
	}
	// Construct message key and value
	expectedKey := []byte("test-key")
	expectedValue := &test_proto.Test{
		TraceId: uuid.New().String(),
		Message: "hello world",
	}

	// ACT
	// Send a message to the test topic
	err := s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: test.TestBasicTopic,
		Key:   sarama.ByteEncoder(expectedKey),
		Value: kafka.ProtoEncoder(expectedValue),
	})
	s.Require().Nil(err, "Producer should send message successfully")
	// Set up a timeout context waiting for consumer to consume a message
	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Consume message
	tr, ok := s.router.GetRoutes()[0].(*test.TestRoute)
	s.Require().True(ok, "Consumer route should be correct type")
	select {
	case event := <-tr.Events():
		// ASSERT
		s.Require().Equal(test.TestBasicTopic, event.Topic, "Consumer route should receive the event from the correct topic")
		s.Require().Equal(string(expectedKey), event.Key, "Consumer route should receive the event with the correct key")
		val, ok := event.Value.(*test_proto.Test)
		s.Require().True(ok, "Consumer route should receive the correct event value type")
		s.Require().Equal(expectedValue.TraceId, val.TraceId, "Consumer route should receive the correct trace id")
		s.Require().Equal(expectedValue.TraceId, val.GetTraceId(), "Consumer route should receive the correct trace id")
		s.Require().Equal(expectedValue.Message, val.GetMessage(), "Consumer route should receive the correct message")
	case <-ctx.Done():
		s.T().Fatal("Test timed out waiting for the event to be processed")
	}
}
