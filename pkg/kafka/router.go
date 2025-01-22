package kafka

import (
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
)

// KafkaRouter distributes incoming messages to the correct handler.
// KafkaRouter also handles DLQ logic.
type KafkaRouter interface {
	// Ready returns a channel that signals when the router is ready
	Ready() <-chan bool
	// GetRoutes returns routes that are registered within the router
	GetRoutes() []Route
	// Implement the underlying interface
	sarama.ConsumerGroupHandler
}

func NewKafkaRouter(routes []Route, logger logger.Logger, tracer tracer.Tracer, dlq *KafkaDLQProducer) KafkaRouter {
	return &kafkaRouterImpl{
		ready:  make(chan bool),
		routes: routes,
		logger: logger,
		tracer: tracer,
		dlq:    dlq,
	}
}

type kafkaRouterImpl struct {
	ready  chan bool
	routes []Route

	logger logger.Logger
	tracer tracer.Tracer
	dlq    *KafkaDLQProducer
}

func (r *kafkaRouterImpl) Ready() <-chan bool {
	return r.ready
}

func (r *kafkaRouterImpl) GetRoutes() []Route {
	return r.routes
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (r *kafkaRouterImpl) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(r.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (r *kafkaRouterImpl) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim starts a consumer loop for the given claim's messages.
// Once the Messages() channel is closed, the Handler must finish its processing loop and exit.
func (r *kafkaRouterImpl) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			ctx := r.tracer.Extract(session.Context(), otelsarama.NewConsumerMessageCarrier(message))
			ctx, span := r.tracer.StartSpanWithCaller(ctx)
			defer span.End()

			if !ok {
				r.logger.WithCtx(ctx).Info("Kafka consumer message channel is closed.")
				return nil
			}

			var errs []DLQError

			for _, route := range r.routes {
				val, err := route.Decode(ctx, message.Value)
				if err != nil {
					errs = append(errs, DLQError{Error: err, RouteID: route.Identifier()})
					continue
				}

				event := &Event{
					Headers:   message.Headers,
					Timestamp: message.Timestamp,
					Key:       string(message.Key),
					Value:     val,
					Topic:     message.Topic,
				}

				match, err := route.Match(ctx, event)
				if err != nil {
					errs = append(errs, DLQError{Error: err, RouteID: route.Identifier()})
					continue
				}

				if match {
					err = route.Handle(ctx, event)
					if err != nil {
						errs = append(errs, DLQError{Error: err, RouteID: route.Identifier()})
					}
				}
			}

			if len(errs) > 0 {
				r.dlq.MoveMessageToDLQ(ctx, message, errs)
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

var _ KafkaRouter = (*kafkaRouterImpl)(nil)
