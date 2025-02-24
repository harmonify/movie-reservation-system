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
	auth_proto "github.com/harmonify/movie-reservation-system/pkg/proto/auth"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type AuthServiceClientParam struct {
	fx.In
	grpc_pkg.GrpcClientParam
	error_pkg.ErrorMapper
}

type authServiceClientImpl struct {
	client      auth_proto.AuthServiceClient
	errorMapper error_pkg.ErrorMapper
}

func NewAuthServiceClient(p AuthServiceClientParam, cfg *config.MovieServiceConfig) (auth_proto.AuthServiceClient, error) {
	interceptor := failsafegrpc.NewUnaryClientInterceptor(
		retrypolicy.Builder[*auth_proto.AuthResponse]().
			WithBackoff(100*time.Millisecond, time.Second).
			WithJitterFactor(0.2).
			WithMaxRetries(4).
			Build(),
		circuitbreaker.Builder[*auth_proto.AuthResponse]().
			HandleErrorTypes(&error_pkg.ErrorWithDetails{}).
			HandleIf(func(_ *auth_proto.AuthResponse, err error) bool {
				ed := err.(*error_pkg.ErrorWithDetails)
				return ed.Code == error_pkg.BadGatewayError.Code || ed.GrpcCode == codes.Unavailable
			}).
			// 4 failures in 10 attempts when the circuit is half-open will open the circuit breaker.
			WithFailureThresholdRatio(4, 10).
			// 6 successes in 10 attempts when the circuit is half-open will close the circuit breaker.
			WithSuccessThresholdRatio(6, 10).
			WithDelay(5*time.Second).
			Build(),
		timeout.With[*auth_proto.AuthResponse](10*time.Second),
	)

	client, err := grpc_pkg.NewGrpcClient(
		p.GrpcClientParam,
		&grpc_pkg.GrpcClientConfig{
			Address: cfg.GrpcAuthServiceUrl,
		},
		grpc.WithUnaryInterceptor(interceptor),
	)
	if err != nil {
		return nil, err
	}

	return &authServiceClientImpl{
		client:      auth_proto.NewAuthServiceClient(client.Conn),
		errorMapper: p.ErrorMapper,
	}, nil
}

func (c *authServiceClientImpl) Auth(ctx context.Context, in *auth_proto.AuthRequest, opts ...grpc.CallOption) (*auth_proto.AuthResponse, error) {
	res, err := c.client.Auth(ctx, in, opts...)
	if err != nil {
		de, _ := c.errorMapper.FromGrpcError(err)
		return nil, de
	}
	return res, nil
}
