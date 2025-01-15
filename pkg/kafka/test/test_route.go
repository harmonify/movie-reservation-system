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
	events chan *kafka.Event
	logger logger.Logger
}

// Events returns the channel for consumed event messages.
// This method serve as a utility for testing purposes.
func (c *TestRoute) Events() <-chan *kafka.Event {
	return c.events
}

// Match implements kafka.Route.
func (r *TestRoute) Match(topic string) bool {
	return topic == TestBasicTopic
}

// Decode implements kafka.Route.
func (r *TestRoute) Decode(message *sarama.ConsumerMessage) (*kafka.Event, error) {
	val := &test_proto.Test{}
	err := proto.Unmarshal(message.Value, val)
	if err != nil {
		return nil, err
	}

	return &kafka.Event{
		Headers:   message.Headers,
		Timestamp: message.Timestamp,
		TraceID:   val.GetTraceId(),
		Key:       string(message.Key),
		Value:     val,
		Topic:     message.Topic,
	}, nil
}

// Handle implements kafka.Route.
func (r *TestRoute) Handle(ctx context.Context, event *kafka.Event) error {
	val, ok := event.Value.(*test_proto.Test)
	if !ok {
		return errors.New("invalid event value type")
	}
	r.logger.Debug(
		fmt.Sprintf(
			"Message claimed: topic = %s, timestamp = %v, trace_id = %s, key = %s, value = %s, headers = %v",
			event.Topic,
			event.Timestamp,
			event.TraceID,
			event.Key,
			proto.Message(val),
			event.Headers,
		),
	)
	r.events <- event
	return nil
}

// NewTestRoute initializes a new TestRoute
func NewTestRoute(logger logger.Logger) *TestRoute {
	return &TestRoute{
		logger: logger,
		events: make(chan *kafka.Event, 100),
	}
}
