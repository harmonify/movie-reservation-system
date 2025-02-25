package http_middleware

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/ratelimiter"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/mennanov/limiters"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// RateLimiterHttpMiddleware is the interface for rate limiter middleware
// It provides functions to limit the request by IP address or UUID
// The middleware will return 429 Too Many Requests if the request is limited
// Currently, the middleware will still limit the request even if the request is failed (e.g. internal server error)
type RateLimiterHttpMiddleware interface {
	// LimitByIP limits the request by IP address. It currently panics if the config is invalid
	LimitByIP(cfg *ratelimiter.RateLimiterConfig) gin.HandlerFunc
	// LimitByUUID limits the request by UUID. It currently panics if the config is invalid
	LimitByUUID(cfg *ratelimiter.RateLimiterConfig) gin.HandlerFunc
	// LimitBy is a generic function to limit the request. It currently panics if the provided config is invalid
	// The function f is used to get the ID to limit the request. If it returns an empty string, the middleware will not limit the request and call the next handler (only if the request context is not aborted).
	LimitBy(cfg *ratelimiter.RateLimiterConfig, f func(*gin.Context) string) gin.HandlerFunc
}

type RateLimiterHttpMiddlewareParam struct {
	fx.In

	Logger              logger.Logger
	Tracer              tracer.Tracer
	Util                *util.Util
	RateLimiterRegistry ratelimiter.RateLimiterRegistry
	Response            http_pkg.HttpResponse
}

type RateLimiterHttpMiddlewareResult struct {
	fx.Out

	RateLimiterHttpMiddleware RateLimiterHttpMiddleware
}

type rateLimiterHttpMiddlewareImpl struct {
	logger              logger.Logger
	tracer              tracer.Tracer
	util                *util.Util
	rateLimiterRegistry ratelimiter.RateLimiterRegistry
	response            http_pkg.HttpResponse
}

func NewRateLimiterHttpMiddleware(p RateLimiterHttpMiddlewareParam) RateLimiterHttpMiddlewareResult {
	return RateLimiterHttpMiddlewareResult{
		RateLimiterHttpMiddleware: &rateLimiterHttpMiddlewareImpl{
			logger:              p.Logger,
			tracer:              p.Tracer,
			util:                p.Util,
			rateLimiterRegistry: p.RateLimiterRegistry,
			response:            p.Response,
		},
	}
}

func (m *rateLimiterHttpMiddlewareImpl) LimitByIP(cfg *ratelimiter.RateLimiterConfig) gin.HandlerFunc {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		// TODO: Currently panic at startup if the config is invalid
		panic(err)
	}

	return func(c *gin.Context) {
		ctx, span := m.tracer.StartSpanWithCaller(c.Request.Context())
		defer span.End()
		m.limitById(ctx, c, c.ClientIP(), cfg)
	}
}

func (m *rateLimiterHttpMiddlewareImpl) LimitByUUID(cfg *ratelimiter.RateLimiterConfig) gin.HandlerFunc {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		// TODO: Currently panic at startup if the config is invalid
		panic(err)
	}

	return func(c *gin.Context) {
		ctx, span := m.tracer.StartSpanWithCaller(c.Request.Context())
		defer span.End()

		userInfo, err := m.util.HttpUtil.GetUserInfo(c.Request)
		if err != nil {
			m.response.Send(c, nil, err)
			c.Abort()
			return
		}

		m.limitById(ctx, c, userInfo.UUID, cfg)
	}
}

func (m *rateLimiterHttpMiddlewareImpl) LimitBy(cfg *ratelimiter.RateLimiterConfig, f func(*gin.Context) string) gin.HandlerFunc {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		// TODO: Currently panic at startup if the config is invalid
		panic(err)
	}

	return func(c *gin.Context) {
		ctx, span := m.tracer.StartSpanWithCaller(c.Request.Context())
		defer span.End()

		id := f(c)
		m.logger.WithCtx(ctx).Debug("limit by", zap.String("id", id), zap.Any("config", cfg))

		if c.IsAborted() {
			return
		}

		if id == "" {
			c.Next()
			return
		}

		m.limitById(ctx, c, id, cfg)
	}
}

// Generic flow to limit the request by ID
func (m *rateLimiterHttpMiddlewareImpl) limitById(ctx context.Context, c *gin.Context, id string, cfg *ratelimiter.RateLimiterConfig) {
	rl, err := m.rateLimiterRegistry.GetHttpRequestRateLimiter(
		&ratelimiter.HttpRequestRateLimiterParam{
			ID:     id,
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
		},
		cfg,
	)
	if err != nil {
		m.logger.WithCtx(ctx).Error("failed to get rate limiter", zap.Error(err))
		m.response.Send(c, nil, error_pkg.InternalServerError)
		c.Abort()
		return
	}

	retryAfter, err := rl.Limit(ctx)
	if err != nil {
		if errors.Is(err, limiters.ErrLimitExhausted) {
			m.logger.WithCtx(ctx).Debug("limit exhausted", zap.Error(err))
			if retryAfter > 0 {
				c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
				m.response.Send(c, nil, error_pkg.RateLimitExceededError)
				c.Abort()
				return
			}
		} else {
			m.logger.WithCtx(ctx).Error("failed to limit request", zap.Error(err))
			m.response.Send(c, nil, error_pkg.InternalServerError)
			c.Abort()
			return
		}
	}

	c.Next()
}
