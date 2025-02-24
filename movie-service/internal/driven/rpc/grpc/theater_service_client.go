package grpc

import (
	"context"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven/config"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	grpc_failsafe "github.com/harmonify/movie-reservation-system/pkg/grpc/failsafe"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	failsafe_object_logger "github.com/harmonify/movie-reservation-system/pkg/logger/object/failsafe"
	circuitbreaker_object_logger "github.com/harmonify/movie-reservation-system/pkg/logger/object/failsafe/circuitbreaker"
	theater_proto "github.com/harmonify/movie-reservation-system/pkg/proto/theater"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type TheaterServiceClientParam struct {
	fx.In
	grpc_pkg.GrpcClientParam
	error_pkg.ErrorMapper
	logger.Logger
	tracer.Tracer
}

type theaterServiceClientImpl struct {
	errorMapper error_pkg.ErrorMapper
	client      theater_proto.TheaterServiceClient
	logger      logger.Logger
	tracer      tracer.Tracer
}

func NewTheaterServiceClient(p TheaterServiceClientParam, cfg *config.MovieServiceConfig) (theater_proto.TheaterServiceClient, error) {
	executor := failsafe.NewExecutor(
		retrypolicy.Builder[any]().
			WithBackoff(100*time.Millisecond, time.Second).
			WithJitterFactor(0.2).
			WithMaxAttempts(4).
			OnRetry(func(event failsafe.ExecutionEvent[any]) {
				p.Logger.WithCtx(event.Context()).Warn("failsafe retry policy retrying", zap.Any("event", failsafe_object_logger.NewLoggableAnyExecutionEvent(event, true)))
			}).
			OnRetriesExceeded(func(event failsafe.ExecutionEvent[any]) {
				p.Logger.WithCtx(event.Context()).Warn("failsafe retry policy retries exceeded", zap.Any("event", failsafe_object_logger.NewLoggableAnyExecutionEvent(event, true)))
			}).
			OnFailure(func(event failsafe.ExecutionEvent[any]) {
				p.Logger.WithCtx(event.Context()).Warn("failsafe retry policy failure", zap.Any("event", failsafe_object_logger.NewLoggableAnyExecutionEvent(event, true)))
			}).
			OnSuccess(func(event failsafe.ExecutionEvent[any]) {
				p.Logger.WithCtx(event.Context()).Debug("failsafe retry policy success", zap.Any("event", failsafe_object_logger.NewLoggableAnyExecutionEvent(event, true)))
			}).
			ReturnLastFailure().
			Build(),
		circuitbreaker.Builder[any]().
			HandleIf(func(_ any, err error) bool {
				if err == nil {
					return false
				}
				if ed, ok := p.ErrorMapper.FromGrpcError(err); ok && ed != nil {
					return (ed.Code == error_pkg.ServiceUnavailableError.Code ||
						ed.GrpcCode == codes.Unavailable ||
						ed.GrpcCode == codes.DeadlineExceeded ||
						ed.GrpcCode == codes.ResourceExhausted)
				}
				return false
			}).
			WithFailureThresholdRatio(4, 10).
			WithSuccessThresholdRatio(6, 10).
			WithDelay(5*time.Second).
			OnStateChanged(func(event circuitbreaker.StateChangedEvent) {
				p.Logger.WithCtx(event.Context()).Debug("failsafe circuit breaker policy state changed", zap.Any("event", circuitbreaker_object_logger.NewLoggableStateChangedEvent(event)))
			}).
			Build(),
		timeout.Builder[any](10*time.Second).
			OnTimeoutExceeded(func(event failsafe.ExecutionDoneEvent[any]) {
				p.Logger.WithCtx(event.Context()).Warn("failsafe timeout policy exceeded", zap.Any("event", failsafe_object_logger.NewLoggableAnyExecutionDoneEvent(event, true)))
			}).
			Build(),
	)

	interceptor := grpc_failsafe.NewUnaryClientInterceptorWithExecutorContext(executor, p.Tracer)

	client, err := grpc_pkg.NewGrpcClient(
		p.GrpcClientParam,
		&grpc_pkg.GrpcClientConfig{
			Address: cfg.GrpcTheaterServiceUrl,
		},
		grpc.WithUnaryInterceptor(interceptor),
	)
	if err != nil {
		return nil, err
	}

	return &theaterServiceClientImpl{
		client:      theater_proto.NewTheaterServiceClient(client.Conn),
		errorMapper: p.ErrorMapper,
		logger:      p.Logger,
		tracer:      p.Tracer,
	}, nil
}

// Get movies with active showtimes
func (c *theaterServiceClientImpl) GetActiveMovies(ctx context.Context, in *theater_proto.GetActiveMoviesRequest, opts ...grpc.CallOption) (*theater_proto.GetActiveMoviesResponse, error) {
	ctx, span := c.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := c.client.GetActiveMovies(ctx, in, opts...)
	if err != nil {
		c.logger.WithCtx(ctx).Error("failed to get active movies", zap.Error(err), zap.Any("in", in))
		if de, ok := c.errorMapper.FromFailsafeError(err); ok {
			return nil, de
		} else {
			de, _ := c.errorMapper.FromGrpcError(err)
			return nil, de
		}
	}

	return res, nil
}

// Get active showtimes for a movie
func (c *theaterServiceClientImpl) GetActiveShowtimes(ctx context.Context, in *theater_proto.GetActiveShowtimesRequest, opts ...grpc.CallOption) (*theater_proto.GetActiveShowtimesResponse, error) {
	ctx, span := c.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := c.client.GetActiveShowtimes(ctx, in, opts...)
	if err != nil {
		c.logger.WithCtx(ctx).Error("failed to get active showtimes", zap.Error(err), zap.Any("input", in))
		if de, ok := c.errorMapper.FromFailsafeError(err); ok {
			return nil, de
		} else {
			de, _ := c.errorMapper.FromGrpcError(err)
			return nil, de
		}
	}

	return res, nil
}

// Get available seats for a showtime
func (c *theaterServiceClientImpl) GetAvailableSeats(ctx context.Context, in *theater_proto.GetAvailableSeatsRequest, opts ...grpc.CallOption) (*theater_proto.GetAvailableSeatsResponse, error) {
	ctx, span := c.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := c.client.GetAvailableSeats(ctx, in, opts...)
	if err != nil {
		c.logger.WithCtx(ctx).Error("failed to get available seats", zap.Error(err), zap.Any("input", in))
		if de, ok := c.errorMapper.FromFailsafeError(err); ok {
			return nil, de
		} else {
			de, _ := c.errorMapper.FromGrpcError(err)
			return nil, de
		}
	}

	return res, nil
}
