package ratelimiter_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/lib/cache"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/ratelimiter"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type RateLimiterRegistryTestSuite struct {
	suite.Suite
	app                 *fx.App
	rateLimiterRegistry ratelimiter.RateLimiterRegistry
}

func TestRateLimiterRegistryTestSuite(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}
	suite.Run(t, new(RateLimiterRegistryTestSuite))
}

func (s *RateLimiterRegistryTestSuite) SetupSuite() {
	app, rateLimiterRegistry, err := createNewRegistry()
	if err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
	s.app = app
	s.rateLimiterRegistry = rateLimiterRegistry
}

func (s *RateLimiterRegistryTestSuite) TearDownSuite() {
	s.app.Done()
}

func (s *RateLimiterRegistryTestSuite) TestRateLimiterRegistry_GetHttpRequestRateLimiter_Simple() {
	rl, err := s.rateLimiterRegistry.GetHttpRequestRateLimiter(ratelimiter.HttpRequestRateLimiterParam{
		IP:     "192.168.200.2",
		Method: http.MethodGet,
		Path:   "/test",
	})
	s.Require().NoError(err)

	for i := 0; i < 5; i++ {
		retryAfter, err := rl.Take(context.Background(), 1)
		if i < 2 {
			s.Require().NoError(err)
			s.Require().Equal(time.Duration(0), retryAfter)
		} else {
			s.Require().Error(err)
			s.Require().Greater(retryAfter, time.Duration(0))
		}
	}

	// Test another client
	rl, err = s.rateLimiterRegistry.GetHttpRequestRateLimiter(ratelimiter.HttpRequestRateLimiterParam{
		IP:     "192.168.200.3",
		Method: http.MethodGet,
		Path:   "/test",
	})
	s.Require().NoError(err)

	for i := 0; i < 5; i++ {
		retryAfter, err := rl.Take(context.Background(), 1)
		if i < 2 {
			s.Require().NoError(err)
			s.Require().Equal(time.Duration(0), retryAfter)
		} else {
			s.Require().Error(err)
			s.Require().Greater(retryAfter, time.Duration(0))
		}
	}
}

func (s *RateLimiterRegistryTestSuite) TestRateLimiterRegistry_GetHttpRequestRateLimiter_Refill() {
	rl, err := s.rateLimiterRegistry.GetHttpRequestRateLimiter(ratelimiter.HttpRequestRateLimiterParam{
		IP:     "192.168.200.4",
		Method: http.MethodPost,
		Path:   "/test",
	})
	s.Require().NoError(err)

	// 0ms mark
	retryAfter, err := rl.Take(context.Background(), 1)
	s.Assert().NoError(err)
	s.Assert().Equal(time.Duration(0), retryAfter)
	time.Sleep(500 * time.Millisecond)

	// 500ms mark
	retryAfter, err = rl.Take(context.Background(), 1)
	s.Assert().NoError(err)
	s.Assert().Equal(time.Duration(0), retryAfter)
	time.Sleep(500 * time.Millisecond)

	// 1000ms mark
	retryAfter, err = rl.Take(context.Background(), 1)
	s.Assert().Error(err)
	s.Assert().Equal(3*time.Second, retryAfter)
	time.Sleep(500 * time.Millisecond)

	// 1500ms mark
	retryAfter, err = rl.Take(context.Background(), 1)
	s.Assert().Error(err)
	s.Assert().Equal(3*time.Second, retryAfter)
	time.Sleep(500 * time.Millisecond)

	// 2000ms mark
	retryAfter, err = rl.Take(context.Background(), 1)
	s.Assert().Error(err)
	s.Assert().Equal(3*time.Second, retryAfter)
	time.Sleep(500 * time.Millisecond)

	// 2500ms mark
	retryAfter, err = rl.Take(context.Background(), 1)
	s.Assert().Error(err)
	s.Assert().Equal(3*time.Second, retryAfter)
	time.Sleep(600 * time.Millisecond)

	// 3100ms mark
	retryAfter, err = rl.Take(context.Background(), 1)
	s.Assert().NoError(err)
	s.Assert().Equal(time.Duration(0), retryAfter)
	time.Sleep(1900 * time.Millisecond)

	// 5000ms mark
	retryAfter, err = rl.Take(context.Background(), 1)
	s.Assert().Error(err)
	s.Assert().Equal(3*time.Second, retryAfter)
	time.Sleep(2000 * time.Millisecond)

	// 7000ms mark
	retryAfter, err = rl.Take(context.Background(), 1)
	s.Assert().NoError(err)
	s.Assert().Equal(time.Duration(0), retryAfter)
	time.Sleep(100 * time.Millisecond)

	// 7100ms mark
	retryAfter, err = rl.Take(context.Background(), 1)
	s.Assert().NoError(err)
	s.Assert().Equal(time.Duration(0), retryAfter)

	// +- 7100ms mark
	retryAfter, err = rl.Take(context.Background(), 1)
	s.Assert().Error(err)
	s.Assert().Equal(3*time.Second, retryAfter)
}

func (s *RateLimiterRegistryTestSuite) TestRateLimiter_Len() {
	app, rateLimiterRegistry, err := createNewRegistry()
	defer app.Done()

	s.Require().NoError(err)
	s.Require().Equal(0, rateLimiterRegistry.Len())

	rl, err := rateLimiterRegistry.GetHttpRequestRateLimiter(ratelimiter.HttpRequestRateLimiterParam{
		IP:     "192.168.100.5",
		Method: http.MethodPost,
		Path:   "/test",
	})
	s.Require().NoError(err)
	s.Require().Equal(1, rateLimiterRegistry.Len())

	rl.Take(context.Background(), 1)
	s.Require().Equal(1, rateLimiterRegistry.Len())

	rl.Take(context.Background(), 1)
	s.Require().Equal(1, rateLimiterRegistry.Len())
	time.Sleep(1 * time.Second)

	rl2, err := rateLimiterRegistry.GetHttpRequestRateLimiter(ratelimiter.HttpRequestRateLimiterParam{
		IP:     "192.168.100.6",
		Method: http.MethodPost,
		Path:   "/test",
	})
	s.Require().NoError(err)

	rl2.Take(context.Background(), 2)
	s.Require().Equal(2, rateLimiterRegistry.Len())

	time.Sleep(3 * time.Second)
	s.Require().Equal(1, rateLimiterRegistry.Len())
	
	rl2.Take(context.Background(), 1)
	s.Require().Equal(1, rateLimiterRegistry.Len())

	time.Sleep(5 * time.Second)
	s.Require().Equal(0, rateLimiterRegistry.Len())
}

// func (s *RateLimiterRegistryTestSuite) TestGRPCRateLimiting() {
// 	interceptor := s.rateLimiterRegistry.LimitGRPCUnaryInterceptor()
// 	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
// 		return "success", nil
// 	}

// 	info := &grpc.UnaryServerInfo{
// 		FullMethod: "/test.Service/Method",
// 	}

// 	for i := 0; i < 10; i++ {
// 		_, err := interceptor(context.Background(), nil, info, handler)
// 		if i < 5 {
// 			s.Require().NoError(err)
// 		} else {
// 			s.Require().Equal(codes.ResourceExhausted, status.Code(err))
// 		}
// 	}
// }

// func (s *RateLimiterRegistryTestSuite) TestGRPCStreamRateLimiting() {
// 	interceptor := s.rateLimiterRegistry.LimitGRPCStreamInterceptor()
// 	handler := func(srv interface{}, stream grpc.ServerStream) error {
// 		return nil
// 	}

// 	info := &grpc.StreamServerInfo{
// 		FullMethod: "/test.Service/StreamMethod",
// 	}

// 	mockStream := &MockServerStream{}
// 	for i := 0; i < 10; i++ {
// 		err := interceptor(nil, mockStream, info, handler)
// 		if i < 5 {
// 			s.Require().NoError(err)
// 		} else {
// 			s.Require().Equal(codes.ResourceExhausted, status.Code(err))
// 		}
// 	}
// }

// type MockServerStream struct {
// 	grpc.ServerStream
// }

// func (m *MockServerStream) Context() context.Context {
// 	return context.Background()
// }

func createNewRegistry() (*fx.App, ratelimiter.RateLimiterRegistry, error) {
	var registry ratelimiter.RateLimiterRegistry

	app := fx.New(
		fx.Provide(
			func() *config.Config {
				return &config.Config{
					Env:               "test",
					AppSecret:         "1233334556905407",
					ServiceIdentifier: "user-service",
					ServiceBaseUrl:    "http://localhost:8080",
					LogLevel:          "debug",
					LogType:           "console",
					RedisHost:         "localhost",
					RedisPort:         "6379",
					RedisPass:         "secret",
				}
			},
			func() *ratelimiter.RateLimiterConfig {
				return &ratelimiter.RateLimiterConfig{
					Capacity:   2,
					RefillRate: 3 * time.Second,
				}
			},
		),
		logger.LoggerModule,
		cache.RedisModule,
		ratelimiter.RateLimiterModule,
		fx.Invoke(
			func(rateLimiterRegistry ratelimiter.RateLimiterRegistry) {
				registry = rateLimiterRegistry
			},
			func(lc fx.Lifecycle, redis *cache.Redis) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return redis.Client.Close()
					},
				})
			},
		),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := app.Start(ctx)

	return app, registry, err
}
