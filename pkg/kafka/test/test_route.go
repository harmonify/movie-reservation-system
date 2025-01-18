package test

import (
	"context"
	"fmt"

	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	test_proto "github.com/harmonify/movie-reservation-system/pkg/test/proto"
	"google.golang.org/protobuf/proto"
)

// TestRoute represents a Kafka route handler
type TestRoute struct {
	logger    logger.Logger
	listeners []kafka.EventListener
}

func (r *TestRoute) Decode(ctx context.Context, value []byte) (interface{}, error) {
	val := &test_proto.Test{}
	if err := proto.Unmarshal(value, val); err != nil {
		return nil, kafka.ErrDecodeFailed
	}
	return val, nil
}

func (r *TestRoute) Match(ctx context.Context, event *kafka.Event) (bool, error) {
	return event.Topic == TestRouterTopic, nil
}

func (r *TestRoute) Handle(ctx context.Context, event *kafka.Event) error {
	// Notify listeners
	for _, listener := range r.listeners {
		listener.OnEvent(ctx, event)
	}

	// Production handling logic
	val, ok := event.Value.(*test_proto.Test)
	if !ok {
		return kafka.ErrInvalidValueType
	}
	r.logger.Debug(
		fmt.Sprintf(
			"Message claimed: topic = %s, timestamp = %v, trace_id = %s, key = %s, value = %s",
			event.Topic,
			event.Timestamp,
			event.TraceID,
			event.Key,
			val,
		),
	)

	return nil
}

func (r *TestRoute) AddEventListener(listener kafka.EventListener) {
	r.listeners = append(r.listeners, listener)
}

// NewTestRoute initializes a new TestRoute
func NewTestRoute(logger logger.Logger) *TestRoute {
	return &TestRoute{
		logger:    logger,
		listeners: []kafka.EventListener{},
	}
}
