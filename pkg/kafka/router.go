package kafka

import (
	"context"
	"errors"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/zap"
)

// KafkaRouter distributes incoming messages to the correct handler
type KafkaRouter interface {
	// Ready returns a channel that signals when the router is ready
	Ready() <-chan bool
	// GetRoutes returns routes that are registered within the router
	GetRoutes() []Route
	// Implement the underlying interface
	sarama.ConsumerGroupHandler
}

func NewKafkaRouter(routes []Route, logger logger.Logger, tracer tracer.Tracer) KafkaRouter {
	return &kafkaRouterImpl{
		ready:  make(chan bool),
		routes: routes,
		logger: logger,
		tracer: tracer,
	}
}

type kafkaRouterImpl struct {
	ready  chan bool
	routes []Route

	logger logger.Logger
	tracer tracer.Tracer
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
			if !ok {
				finalErr = ErrMessageChannelClosed
				return finalErr
			}

			ctx := c.tracer.Extract(session.Context(), otelsarama.NewConsumerMessageCarrier(message))
			ctx, span := c.tracer.StartSpanWithCaller(ctx)
			defer span.End()

			for _, route := range c.routes {
				event, err := c.constructEventForRoute(ctx, route, message)
				if err != nil {
					finalErr = errors.Join(finalErr, err)
					continue
				}

				match, err := route.Match(ctx, event)
				if err != nil {
					finalErr = errors.Join(finalErr, err)
					continue
				}

				if match {
					err = route.Handle(ctx, event)
					if err != nil {
						finalErr = errors.Join(finalErr, err)
					}
				}
			}

			if finalErr != nil {
				c.logger.WithCtx(ctx).Error("Errors", zap.Error(finalErr))
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *kafkaRouterImpl) constructEventForRoute(ctx context.Context, route Route, message *sarama.ConsumerMessage) (*Event, error) {
	ctx, span := c.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	val, err := route.Decode(ctx, message.Value)
	if err != nil {
		return nil, err
	}

	return &Event{
		TraceID:   span.SpanContext().TraceID().String(),
		Timestamp: message.Timestamp,
		Key:       string(message.Key),
		Value:     val,
		Topic:     message.Topic,
	}, nil
}

var _ KafkaRouter = (*kafkaRouterImpl)(nil)
