package test

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
)

// TestConsumer represents a Sarama consumer handler
type TestConsumer struct {
	ready     chan bool
	logger    logger.Logger
	tracer    tracer.Tracer
	listeners []kafka.MessageListener
}

// NewTestConsumer initializes a new TestConsumer
func NewTestConsumer(logger logger.Logger, tracer tracer.Tracer) *TestConsumer {
	return &TestConsumer{
		ready:     make(chan bool),
		logger:    logger,
		tracer:    tracer,
		listeners: []kafka.MessageListener{},
	}
}

// Ready returns a channel that signals when the consumer is ready.
// This method serve as a utility for testing purposes.
func (c *TestConsumer) Ready() <-chan bool {
	return c.ready
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *TestConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *TestConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim starts a consumer loop for the given claim's messages
func (c *TestConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				c.logger.Info("message channel was closed")
				return nil
			}

			ctx := c.tracer.Extract(session.Context(), otelsarama.NewConsumerMessageCarrier(message))
			ctx, span := c.tracer.StartSpanWithCaller(ctx)
			defer span.End()

			// Notify listeners
			for _, listener := range c.listeners {
				listener.OnMessage(ctx, message)
			}

			c.logger.Info(fmt.Sprintf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic))
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (r *TestConsumer) AddMessageListener(listener kafka.MessageListener) {
	r.listeners = append(r.listeners, listener)
}
