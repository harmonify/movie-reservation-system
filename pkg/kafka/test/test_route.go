package test

import (
	"context"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	test_proto "github.com/harmonify/movie-reservation-system/pkg/test/proto"
	"google.golang.org/protobuf/proto"
)

// TestRoute represents a Kafka route handler
type TestRoute struct {
	events chan *sarama.ConsumerMessage
	logger logger.Logger
}

// Messages returns the channel for consumed event messages.
// This method serve as a utility for testing purposes.
func (c *TestRoute) Messages() <-chan *sarama.ConsumerMessage {
	return c.events
}

// Match implements kafka.Route.
func (r *TestRoute) Match(message *sarama.ConsumerMessage) bool {
	return message.Topic == TestBasicTopic
}

// Handle implements kafka.Route.
func (r *TestRoute) Handle(ctx context.Context, message *sarama.ConsumerMessage) (err error) {
	val := &test_proto.Test{}
	err = proto.Unmarshal(message.Value, val)
	if err != nil {
		return errors.New("invalid event value type")
	}

	r.logger.Debug(
		fmt.Sprintf(
			"Message claimed: topic = %s, timestamp = %v, trace_id = %s, key = %s, value = %s, headers = %v",
			message.Topic,
			message.Timestamp,
			val.GetTraceId(),
			message.Key,
			proto.Message(val),
			message.Headers,
		),
	)
	r.events <- message
	return nil
}

// NewTestRoute initializes a new TestRoute
func NewTestRoute(logger logger.Logger) *TestRoute {
	return &TestRoute{
		logger: logger,
		events: make(chan *sarama.ConsumerMessage, 100),
	}
}

var _ kafka.Route = &TestRoute{}
