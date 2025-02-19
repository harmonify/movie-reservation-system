package kafka_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/kafka/test"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	test_proto "github.com/harmonify/movie-reservation-system/pkg/test/proto"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
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
	// Deps
	tracer        tracer.Tracer
	app           *fx.App
	admin         *kafka.KafkaAdmin
	producer      *kafka.KafkaProducer
	consumerGroup *kafka.KafkaConsumerGroup
	// Basic consumer
	basicConsumer *test.TestConsumer
	// Router consumer
	router kafka.KafkaRouter
}

func (s *KafkaTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.Provide(
			func() *tracer.TracerConfig {
				return &tracer.TracerConfig{
					Env:               "test",
					ServiceIdentifier: "test",
					Type:              "console",
					OtelEndpoint:      "",
				}
			},
			func() *kafka.KafkaConfig {
				return &kafka.KafkaConfig{
					KafkaVersion: "3.9.2",
				}
			},
			func() *logger.LoggerConfig {
				return &logger.LoggerConfig{
					Env:               "test",
					ServiceIdentifier: "test",
					LogType:           "console",
					LogLevel:          "debug",
				}
			},
			func(kc *kafka.KafkaConfig) *kafka.KafkaAdminConfig {
				return &kafka.KafkaAdminConfig{
					KafkaConfig:  kc,
					KafkaBrokers: "localhost:9092",
				}
			},
			func(kc *kafka.KafkaConfig) *kafka.KafkaProducerConfig {
				return &kafka.KafkaProducerConfig{
					KafkaConfig:  kc,
					KafkaBrokers: "localhost:9092",
				}
			},
			func(kc *kafka.KafkaConfig) *kafka.KafkaConsumerGroupConfig {
				return &kafka.KafkaConsumerGroupConfig{
					KafkaConfig:        kc,
					KafkaBrokers:       "localhost:9092",
					KafkaConsumerGroup: "test-group",
				}
			},
			logger.NewConsoleLogger,
			tracer.NewConsoleTracer,
			kafka.NewKafkaAdmin,
			kafka.NewKafkaProducer,
			kafka.NewKafkaConsumerGroup,
			kafka.NewKafkaDLQProducer,
			test.NewTestConsumer,
			kafka.AsRoute(
				test.NewTestRoute,
			),
			fx.Annotate(
				kafka.NewKafkaRouter,
				fx.ParamTags(`group:"kafka-routes"`),
			),
		),
		fx.Invoke(func(t tracer.Tracer, a *kafka.KafkaAdmin, p *kafka.KafkaProducer, cg *kafka.KafkaConsumerGroup, c *test.TestConsumer, r kafka.KafkaRouter) {
			s.tracer = t
			s.admin = a
			s.producer = p
			s.consumerGroup = cg
			s.basicConsumer = c
			s.router = r
		}),
		fx.Invoke(func(a *kafka.KafkaAdmin) {
			err := a.Client.CreateTopic(
				test.TestBasicTopic,
				&sarama.TopicDetail{
					NumPartitions:     1,
					ReplicationFactor: 1,
				},
				false,
			)
			s.Require().Nil(err, "Admin should successfully create test topic for setup process, but got: %s", err)
			err = a.Client.CreateTopic(
				test.TestRouterTopic,
				&sarama.TopicDetail{
					NumPartitions:     1,
					ReplicationFactor: 1,
				},
				false,
			)
			s.Require().Nil(err, "Admin should successfully create 2nd test topic for setup process, but got: %s", err)
		}),
		fx.NopLogger,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *KafkaTestSuite) TearDownSuite() {
	err := s.admin.Client.DeleteTopic(test.TestBasicTopic)
	s.Require().Nil(err, "Admin should successfully delete test topic for teardown process, but got: %s", err)
	err = s.admin.Client.DeleteTopic(test.TestRouterTopic)
	s.Require().Nil(err, "Admin should successfully delete 2nd test topic for teardown process, but got: %s", err)
}

func (s *KafkaTestSuite) TestKafkaSuite_Basic() {
	// ARRANGE
	// Start the consumer in a separate goroutine
	consumerStartupCtx, consumerStartupCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer consumerStartupCancel()
	testMsgListener := test.NewTestMessageListener()
	s.basicConsumer.AddMessageListener(testMsgListener)
	s.consumerGroup.StartConsumer(consumerStartupCtx, []string{test.TestBasicTopic}, s.basicConsumer)
	// Wait for the consumer to be ready
	select {
	case <-s.basicConsumer.Ready():
		s.T().Log("Sarama consumer up and running!")
	case <-consumerStartupCtx.Done():
		s.T().Fatal("Consumer failed to become ready within the timeout")
	}

	// Construct message key and value
	expectedKey := []byte("test-key")
	expectedValue := &test_proto.Test{
		Message: "hello world",
	}

	// ACT
	// Send a message to the test topic
	producerCtx, producerCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer producerCancel()
	producerCtx, producerSpan := s.tracer.StartSpanWithCaller(producerCtx)
	defer producerSpan.End()
	err := s.producer.SendMessage(producerCtx, &sarama.ProducerMessage{
		Topic: test.TestBasicTopic,
		Key:   sarama.ByteEncoder(expectedKey),
		Value: kafka.ProtoEncoder(expectedValue),
	})
	s.Require().True(producerSpan.SpanContext().HasTraceID(), "Producer span context must has trace id")
	s.Require().True(producerSpan.SpanContext().TraceID().IsValid(), "Producer span context must has valid trace id")
	s.Require().Nil(err, "Producer should send message successfully")

	// Consume message
	consumerCtx, consumerCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer consumerCancel()
	select {
	case cm := <-testMsgListener.Messages():
		// ASSERT
		consumerSpanCtx := trace.SpanContextFromContext(cm.Context)
		s.Require().True(consumerSpanCtx.HasTraceID(), "Consumer span context must has trace id")
		s.Require().True(consumerSpanCtx.TraceID().IsValid(), "Consumer span context must has valid trace id")

		s.Require().Equal(
			producerSpan.SpanContext().TraceID().String(),
			consumerSpanCtx.TraceID().String(),
			"Consumer route should receive the correct trace id",
		)

		val := &test_proto.Test{}
		err := proto.Unmarshal(cm.Message.Value, val)
		s.Require().Nil(err, "Consumer should successfully unmarshal the message")
		s.Require().Equal(expectedValue.Message, val.GetMessage(), "Consumer should receive the correct message")
	case <-consumerCtx.Done():
		s.T().Fatal("Test timed out waiting for the message to be processed")
	}
}

func (s *KafkaTestSuite) TestKafkaSuite_Router() {
	// ARRANGE
	// Set up event listener
	testListener := test.NewTestEventListener()
	routes := s.router.GetRoutes()
	s.Require().Len(routes, 1, "Router must have exactly 1 route registered")
	routes[0].AddEventListener(testListener)

	// Start the consumer in a separate goroutine
	consumerStartupCtx, consumerStartupCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer consumerStartupCancel()
	s.consumerGroup.StartConsumer(consumerStartupCtx, []string{test.TestRouterTopic}, s.router)
	select {
	case <-s.router.Ready():
		s.T().Log("Sarama consumer up and running!")
	case <-consumerStartupCtx.Done():
		s.T().Fatal("Consumer failed to become ready within the timeout")
	}

	// Construct message key and value
	expectedKey := []byte("test-key")
	expectedValue := &test_proto.Test{
		Message: "hello world",
	}

	// ACT
	// Send a message to the test topic
	producerCtx, producerCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer producerCancel()
	producerCtx, producerSpan := s.tracer.StartSpanWithCaller(producerCtx)
	defer producerSpan.End()
	err := s.producer.SendMessage(producerCtx, &sarama.ProducerMessage{
		Topic: test.TestRouterTopic,
		Key:   sarama.ByteEncoder(expectedKey),
		Value: kafka.ProtoEncoder(expectedValue),
	})
	s.T().Logf("Producer context: topic=%s spanid=%s traceid=%s", test.TestRouterTopic, producerSpan.SpanContext().SpanID(), producerSpan.SpanContext().TraceID())
	s.Require().True(producerSpan.SpanContext().HasTraceID(), "Producer span context must has trace id")
	s.Require().True(producerSpan.SpanContext().TraceID().IsValid(), "Producer span context must has valid trace id")
	s.Require().Nil(err, "Producer should send message successfully")

	// Consume message
	consumerCtx, consumerCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer consumerCancel()
	select {
	case ce := <-testListener.Events():
		consumerSpanCtx := trace.SpanContextFromContext(ce.Context)
		s.T().Logf("Consumer context: topic=%s spanid=%s traceid=%s", ce.Event.Topic, consumerSpanCtx.SpanID(), consumerSpanCtx.TraceID())
		// ASSERT
		s.Require().Equal(producerSpan.SpanContext().TraceID().String(), consumerSpanCtx.TraceID().String(), "Consumer route should receive the correct trace id")

		s.Require().Equal(test.TestRouterTopic, ce.Event.Topic, "Consumer route should receive the event from the correct topic")

		s.Require().Equal(string(expectedKey), string(ce.Event.Key), "Consumer route should receive the event with the correct key")

		val := &test_proto.Test{}
		if err := proto.Unmarshal(ce.Event.Value, val); err != nil {
			s.T().Fatal("Consumer should successfully unmarshal the event value")
		}
		s.Require().Equal(expectedValue.GetMessage(), val.GetMessage(), "Consumer route should receive the correct message")
	case <-consumerCtx.Done():
		s.T().Fatal("Test timed out waiting for the event to be processed")
	}
}
