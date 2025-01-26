package kafka_driver

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	notification_proto "github.com/harmonify/movie-reservation-system/user-service/internal/driven/proto/notification"
	user_proto "github.com/harmonify/movie-reservation-system/user-service/internal/driven/proto/user"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type OutboxRouteParam struct {
	fx.In

	Logger logger.Logger
	Tracer tracer.Tracer
}

type outboxRouteImpl struct {
	listeners []kafka.EventListener
	logger    logger.Logger
	tracer    tracer.Tracer
}

func NewOutboxRoute(p OutboxRouteParam) kafka.Route {
	return &outboxRouteImpl{
		listeners: []kafka.EventListener{},
		logger:    p.Logger,
		tracer:    p.Tracer,
	}
}

func (r *outboxRouteImpl) Identifier() string {
	return "outbox-verification-handler"
}

func (r *outboxRouteImpl) Decode(ctx context.Context, value []byte) (interface{}, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	val := &user_proto.UserRegistered{}
	if err := proto.Unmarshal(value, val); err != nil {
		r.logger.WithCtx(ctx).Error("Failed to decode outbox proto", zap.Error(err))
		return nil, kafka.ErrMalformedMessage
	}
	return val, nil
}

func (r *outboxRouteImpl) Match(ctx context.Context, event *kafka.Event) (bool, error) {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	return event.Topic == shared.PublicUserRegisteredV1.String(), nil
}

func (r *outboxRouteImpl) Handle(ctx context.Context, event *kafka.Event) error {
	ctx, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// Notify listeners
	for _, listener := range r.listeners {
		listener.OnEvent(ctx, event)
	}

	val, ok := event.Value.(*user_proto.UserRegistered)
	if !ok {
		r.logger.WithCtx(ctx).Error(
			"Failed to assert correct event type",
			zap.String("want", reflect.TypeFor[*notification_proto.SendEmailRequest]().Name()),
			zap.String("got", reflect.TypeOf(val).Name()),
		)
		return kafka.ErrMalformedMessage
	}

	eventjson, _ := json.Marshal(event)
	valjson, _ := json.Marshal(val)
	fmt.Printf("Event: %s\n", eventjson)
	fmt.Printf("Message: %s\n", valjson)

	return nil
}

func (r *outboxRouteImpl) AddEventListener(listener kafka.EventListener) {
	r.listeners = append(r.listeners, listener)
}
