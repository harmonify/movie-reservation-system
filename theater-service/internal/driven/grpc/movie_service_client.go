package grpc

import (
	"context"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/failsafe-go/failsafe-go/timeout"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	grpc_failsafe "github.com/harmonify/movie-reservation-system/pkg/grpc/failsafe"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	failsafe_object_logger "github.com/harmonify/movie-reservation-system/pkg/logger/object/failsafe"
	circuitbreaker_object_logger "github.com/harmonify/movie-reservation-system/pkg/logger/object/failsafe/circuitbreaker"
	movie_proto "github.com/harmonify/movie-reservation-system/pkg/proto/movie"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/driven/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type movieServiceClientParam struct {
	fx.In
	grpc_pkg.GrpcClientParam
	error_pkg.ErrorMapper
	tracer.Tracer
	Logger logger.Logger
}

type movieServiceClientImpl struct {
	client      movie_proto.MovieServiceClient
	errorMapper error_pkg.ErrorMapper
	tracer      tracer.Tracer
	logger      logger.Logger
}

func NewMovieServiceClient(p movieServiceClientParam, cfg *config.TheaterServiceConfig) (movie_proto.MovieServiceClient, error) {
	executor := failsafe.NewExecutor(
		retrypolicy.Builder[*movie_proto.Movie]().
			AbortOnErrors(circuitbreaker.ErrOpen).
			ReturnLastFailure().
			WithBackoff(100*time.Millisecond, time.Second).
			WithJitterFactor(0.2).
			WithMaxAttempts(4).
			OnRetry(func(event failsafe.ExecutionEvent[*movie_proto.Movie]) {
				p.Logger.WithCtx(event.Context()).Warn("failsafe retry policy retrying", zap.Any("event", failsafe_object_logger.NewLoggableAnyExecutionEvent(event, false)))
			}).
			OnRetriesExceeded(func(event failsafe.ExecutionEvent[*movie_proto.Movie]) {
				p.Logger.WithCtx(event.Context()).Warn("failsafe retry policy retries exceeded", zap.Any("event", failsafe_object_logger.NewLoggableAnyExecutionEvent(event, false)))
			}).
			OnFailure(func(event failsafe.ExecutionEvent[*movie_proto.Movie]) {
				p.Logger.WithCtx(event.Context()).Warn("failsafe retry policy failure", zap.Any("event", failsafe_object_logger.NewLoggableAnyExecutionEvent(event, false)))
			}).
			OnSuccess(func(event failsafe.ExecutionEvent[*movie_proto.Movie]) {
				p.Logger.WithCtx(event.Context()).Debug("failsafe retry policy success", zap.Any("event", failsafe_object_logger.NewLoggableAnyExecutionEvent(event, false)))
			}).
			Build(),
		circuitbreaker.Builder[*movie_proto.Movie]().
			// Handle if the error is a ServiceUnavailableError or Unavailable.
			HandleIf(func(_ *movie_proto.Movie, err error) bool {
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
			// 4 failures in 10 attempts when the circuit is half-open will open the circuit breaker.
			WithFailureThresholdRatio(4, 10).
			// 6 successes in 10 attempts when the circuit is half-open will close the circuit breaker.
			WithSuccessThresholdRatio(6, 10).
			// The circuit will be half-open for 5 seconds before transitioning to open.
			WithDelay(5*time.Second).
			OnStateChanged(func(event circuitbreaker.StateChangedEvent) {
				p.Logger.WithCtx(event.Context()).Debug("failsafe circuit breaker policy state changed", zap.Any("state", circuitbreaker_object_logger.NewLoggableStateChangedEvent(event)))
			}).
			Build(),
		timeout.Builder[*movie_proto.Movie](10*time.Second).
			OnTimeoutExceeded(func(event failsafe.ExecutionDoneEvent[*movie_proto.Movie]) {
				p.Logger.WithCtx(event.Context()).Warn("failsafe timeout policy exceeded", zap.Any("event", failsafe_object_logger.NewLoggableAnyExecutionDoneEvent(event, false)))
			}).
			Build(),
	)

	interceptor := grpc_failsafe.NewUnaryClientInterceptorWithExecutorContext(executor, p.Tracer)

	client, err := grpc_pkg.NewGrpcClient(
		p.GrpcClientParam,
		&grpc_pkg.GrpcClientConfig{
			Address: cfg.GrpcMovieServiceUrl,
		},
		grpc.WithUnaryInterceptor(interceptor),
	)
	if err != nil {
		return nil, err
	}

	return &movieServiceClientImpl{
		client:      movie_proto.NewMovieServiceClient(client.Conn),
		errorMapper: p.ErrorMapper,
		logger:      p.Logger,
		tracer:      p.Tracer,
	}, nil
}

func (c *movieServiceClientImpl) GetMovieByID(ctx context.Context, in *movie_proto.GetMovieByIDRequest, opts ...grpc.CallOption) (*movie_proto.GetMovieByIDResponse, error) {
	ctx, span := c.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	res, err := c.client.GetMovieByID(ctx, in, opts...)
	if err != nil {
		c.logger.WithCtx(ctx).Error("failed to call MovieService.GetMovieByID gRPC method", zap.Error(err), zap.Any("input", in))
		if de, ok := c.errorMapper.FromFailsafeError(err); ok {
			return nil, de
		} else {
			de, _ := c.errorMapper.FromGrpcError(err)
			return nil, de
		}
	}

	return res, nil
}
