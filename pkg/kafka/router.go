package kafka

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"go.uber.org/zap"
)

// KafkaRouter distributes incoming messages to the correct handler
type KafkaRouter interface {
	// Ready returns a channel that signals when the router is ready
	Ready() <-chan bool
	// GetRoutes returns routes that are registered within the router
	GetRoutes() []Route

	sarama.ConsumerGroupHandler
}

// Route handle incoming messages from a Topic.
// The first generic type argument corresponds to the message value type.
type Route interface {
	// Match determines if this route should handle the message
	Match(topic string) bool
	// Decode decodes the message value
	Decode(message *sarama.ConsumerMessage) (*Event, error)
	// Handle handles the incoming message that has been decoded
	Handle(ctx context.Context, event *Event) error
}

func NewKafkaRouter(routes []Route, logger logger.Logger) KafkaRouter {
	return &kafkaRouterImpl{
		ready:  make(chan bool),
		routes: routes,
		logger: logger,
	}
}

type kafkaRouterImpl struct {
	ready  chan bool
	routes []Route

	logger logger.Logger
}

func (c *kafkaRouterImpl) Ready() <-chan bool {
	return c.ready
}

func (c *kafkaRouterImpl) GetRoutes() []Route {
	return c.routes
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *kafkaRouterImpl) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *kafkaRouterImpl) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim starts a consumer loop for the given claim's messages
func (c *kafkaRouterImpl) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			var finalErr error
			var traceId string

			if !ok {
				finalErr = errors.New("message channel was closed")
			}

			for _, route := range c.routes {
				if route.Match(message.Topic) {
					event, err := route.Decode(message)
					if err != nil {
						finalErr = errors.Join(finalErr, err)
						continue
					}
					if traceId == "" {
						traceId = event.TraceID
					}
					err = route.Handle(session.Context(), event)
					if err != nil {
						finalErr = errors.Join(finalErr, err)
					}
				}
			}

			if finalErr != nil {
				c.logger.With(zap.String("trace_id", traceId)).Warn(finalErr.Error())
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

var _ KafkaRouter = (*kafkaRouterImpl)(nil)
