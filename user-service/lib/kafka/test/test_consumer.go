package test

import (
	"github.com/IBM/sarama"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
)

// TestConsumer represents a Sarama consumer handler
type TestConsumer struct {
	messages chan *sarama.ConsumerMessage
	ready    chan bool
	logger   logger.Logger
}

// NewTestConsumer initializes a new TestConsumer
func NewTestConsumer(logger logger.Logger) *TestConsumer {
	return &TestConsumer{
		messages: make(chan *sarama.ConsumerMessage, 100), // Buffered channel for messages
		ready:    make(chan bool),
		logger:   logger,
	}
}

// Messages returns the channel for consumed messages.
// This method serve as a utility for testing purposes.
func (c *TestConsumer) Messages() <-chan *sarama.ConsumerMessage {
	return c.messages
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
	// Close the messages channel when cleanup is complete
	close(c.messages)
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
			// c.logger.Info(fmt.Sprintf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic))
			session.MarkMessage(message, "")
			// Pass the message to the channel
			c.messages <- message
		case <-session.Context().Done():
			return nil
		}
	}
}
