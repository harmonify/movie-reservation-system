package grpc

import (
	"context"
	"time"

	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/failsafegrpc"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven/config"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	grpc_pkg "github.com/harmonify/movie-reservation-system/pkg/grpc"
	theater_proto "github.com/harmonify/movie-reservation-system/pkg/proto/theater"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type TheaterServiceClientParam struct {
	fx.In
	grpc_pkg.GrpcClientParam
	error_pkg.ErrorMapper
}

type theaterServiceClientImpl struct {
	errorMapper error_pkg.ErrorMapper
	client      theater_proto.TheaterServiceClient
}

func NewTheaterServiceClient(p TheaterServiceClientParam, cfg *config.MovieServiceConfig) (theater_proto.TheaterServiceClient, error) {
	interceptor := failsafegrpc.NewUnaryClientInterceptor(
		retrypolicy.Builder[any]().
			WithBackoff(100*time.Millisecond, time.Second).
			WithJitterFactor(0.2).
			WithMaxRetries(4).
			Build(),
		circuitbreaker.Builder[any]().
			HandleErrorTypes(&error_pkg.ErrorWithDetails{}).
			HandleIf(func(_ any, err error) bool {
				ed := err.(*error_pkg.ErrorWithDetails)
				return ed.Code == error_pkg.BadGatewayError.Code || ed.GrpcCode == codes.Unavailable
			}).
			// 4 failures in 10 attempts when the circuit is half-open will open the circuit breaker.
			WithFailureThresholdRatio(4, 10).
			// 6 successes in 10 attempts when the circuit is half-open will close the circuit breaker.
			WithSuccessThresholdRatio(6, 10).
			WithDelay(5*time.Second).
			Build(),
		timeout.With[any](10*time.Second),
	)

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
		errorMapper: p.ErrorMapper,
		client:      theater_proto.NewTheaterServiceClient(client.Conn),
	}, nil
}

// Get movies with active showtimes
func (c *theaterServiceClientImpl) GetActiveMovies(ctx context.Context, in *theater_proto.GetActiveMoviesRequest, opts ...grpc.CallOption) (*theater_proto.GetActiveMoviesResponse, error) {
	res, err := c.client.GetActiveMovies(ctx, in, opts...)
	if err != nil {
		de, _ := c.errorMapper.FromGrpcError(err)
		return nil, de
	}
	return res, nil
}

// Get active showtimes for a movie
func (c *theaterServiceClientImpl) GetActiveShowtimes(ctx context.Context, in *theater_proto.GetActiveShowtimesRequest, opts ...grpc.CallOption) (*theater_proto.GetActiveShowtimesResponse, error) {
	res, err := c.client.GetActiveShowtimes(ctx, in, opts...)
	if err != nil {
		de, _ := c.errorMapper.FromGrpcError(err)
		return nil, de
	}
	return res, nil
}

// Get available seats for a showtime
func (c *theaterServiceClientImpl) GetAvailableSeats(ctx context.Context, in *theater_proto.GetAvailableSeatsRequest, opts ...grpc.CallOption) (*theater_proto.GetAvailableSeatsResponse, error) {
	res, err := c.client.GetAvailableSeats(ctx, in, opts...)
	if err != nil {
		de, _ := c.errorMapper.FromGrpcError(err)
		return nil, de
	}
	return res, nil
}
