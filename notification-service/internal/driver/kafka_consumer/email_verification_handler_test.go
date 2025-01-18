package kafka_consumer_test

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	notification_proto "github.com/harmonify/movie-reservation-system/notification-service/internal/driven/proto/notification"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/driver/kafka_consumer"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/kafka/test"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	runId                      = uuid.New().String()
	TestEmailVerificationTopic = "test_email-verification_" + runId
)

func TestEmailVerificationHandlerSuite(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(KafkaTestSuite))
}

type KafkaTestSuite struct {
	suite.Suite

	app           *fx.App
	tracer        tracer.Tracer
	admin         *kafka.KafkaAdmin
	producer      *kafka.KafkaProducer
	consumerGroup *kafka.KafkaConsumerGroup
	router        kafka.KafkaRouter
}

func (s *KafkaTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.Provide(
			func() *config.Config {
				return &config.Config{
					KafkaBrokers:       "localhost:9092",
					KafkaVersion:       "3.9.0",
					KafkaConsumerGroup: "notification-service",
				}
			},
			logger.NewConsoleLogger,
			tracer.NewConsoleTracer,
			kafka.NewKafkaAdmin,
			kafka.NewKafkaProducer,
			kafka.NewKafkaConsumerGroup,
			kafka.AsRoute(
				kafka_consumer.NewEmailVerificationRoute,
			),
			fx.Annotate(
				kafka.NewKafkaRouter,
				fx.ParamTags(`group:"kafka-routes"`),
			),
		),
		fx.Invoke(func(t tracer.Tracer, a *kafka.KafkaAdmin, p *kafka.KafkaProducer, cg *kafka.KafkaConsumerGroup, r kafka.KafkaRouter) {
			s.tracer = t
			s.admin = a
			s.producer = p
			s.consumerGroup = cg
			s.router = r
		}),

		fx.NopLogger,
	)

	err := s.admin.Client.CreateTopic(
		TestEmailVerificationTopic,
		&sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		false,
	)
	s.Require().Nil(err, "Admin should successfully create test topic for setup process")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *KafkaTestSuite) TearDownSuite() {
	err := s.admin.Client.DeleteTopic(TestEmailVerificationTopic)
	s.Require().Nil(err, "Admin should successfully delete test topic for teardown process")
}

func (s *KafkaTestSuite) TestEmailVerificationHandlerSuite_Handle() {
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
	expectedValue := &notification_proto.Email{
		Recipients:   []string{"john_doe@example.com"},
		Subject:      "Email verification",
		TemplateId:   "",
		TemplateData: &anypb.Any{},
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
		s.Require().Equal(consumerSpanCtx.TraceID().String(), ce.Event.TraceID, "Router should correctly parse the message header to the event trace id property")
		s.Require().Equal(producerSpan.SpanContext().TraceID().String(), ce.Event.TraceID, "Email verification handler should receive the correct trace id")

		s.Require().Equal(test.TestRouterTopic, ce.Event.Topic, "Email verification handler should receive the event from the correct topic")

		s.Require().Equal(string(expectedKey), string(ce.Event.Key), "Email verification handler should receive the event with the correct key")

		val, ok := ce.Event.Value.(*notification_proto.Email)
		s.Require().True(ok, "Email verification handler should receive the correct event value type of %s, but got: %s", reflect.TypeFor[*notification_proto.Email](), reflect.TypeOf(val).Name())
		s.Require().Equal(expectedValue.GetRecipients(), val.GetRecipients(), "Email verification handler should receive the correct recipients")
		s.Require().Equal(expectedValue.GetSubject(), val.GetSubject(), "Email verification handler should receive the correct subject")
		s.Require().Equal(expectedValue.GetTemplateId(), val.GetTemplateId(), "Email verification handler should receive the correct template id")
		s.Require().Equal(expectedValue.GetTemplateData(), val.GetTemplateData(), "Email verification handler should receive the correct template id")
	case <-consumerCtx.Done():
		s.T().Fatal("Test timed out waiting for the event to be processed")
	}
}
