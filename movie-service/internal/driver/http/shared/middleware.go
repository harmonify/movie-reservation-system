package http_driver_shared

import (
	"github.com/gin-gonic/gin"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven/config"
	http_middleware "github.com/harmonify/movie-reservation-system/pkg/http/middleware"
	"github.com/harmonify/movie-reservation-system/pkg/metrics"
	"github.com/harmonify/movie-reservation-system/pkg/ratelimiter"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	user_http_middleware "github.com/harmonify/movie-reservation-system/user-service/pkg/http/middleware"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/fx"
)

type (
	HttpMiddleware struct {
		Recovery    http_middleware.RecoveryHttpMiddleware
		Metrics     metrics.PrometheusHttpMiddleware
		Auth        http_middleware.JwtAuthenticationHttpMiddleware
		AuthV2      user_http_middleware.AuthHttpMiddleware
		Rbac        http_middleware.JwtRbacHttpMiddleware
		RateLimiter http_middleware.RateLimiterHttpMiddleware
		Trace       TraceHttpMiddleware
	}
)

var HttpMiddlewareModule = fx.Module(
	"http-middleware",
	fx.Provide(
		http_middleware.NewRecoveryHttpMiddleware,
		metrics.NewPrometheusHttpMiddleware,
		http_middleware.NewJwtHttpMiddleware,
		func(p http_middleware.JwtRbacHttpMiddlewareParam, cfg *config.MovieServiceConfig) (http_middleware.JwtRbacHttpMiddlewareResult, error) {
			return http_middleware.NewJwtRbacHttpMiddleware(p, &http_middleware.JwtHttpMiddlewareConfig{
				Domain: cfg.ServiceIdentifier,
			})
		},
		func(p ratelimiter.RateLimiterRegistryParam, cfg *config.MovieServiceConfig) (ratelimiter.RateLimiterRegistry, error) {
			return ratelimiter.NewRateLimiterRegistry(p, &ratelimiter.RateLimiterRegistryConfig{
				ServiceIdentifier: cfg.ServiceIdentifier,
			})
		},
		http_middleware.NewRateLimiterHttpMiddleware,
		user_http_middleware.NewAuthHttpMiddleware,
		NewTraceHttpMiddleware,
		NewHttpMiddleware,
	),
)

func NewHttpMiddleware(
	recovery http_middleware.RecoveryHttpMiddleware,
	metrics metrics.PrometheusHttpMiddleware,
	jwt http_middleware.JwtAuthenticationHttpMiddleware,
	jwtRbac http_middleware.JwtRbacHttpMiddleware,
	rateLimiter http_middleware.RateLimiterHttpMiddleware,
	auth user_http_middleware.AuthHttpMiddleware,
	trace TraceHttpMiddleware,
) *HttpMiddleware {
	return &HttpMiddleware{
		Recovery:    recovery,
		Metrics:     metrics,
		Auth:        jwt,
		AuthV2:      auth,
		Rbac:        jwtRbac,
		RateLimiter: rateLimiter,
		Trace:       trace,
	}
}

type TraceHttpMiddleware interface {
	ExtractTraceContext(c *gin.Context)
}

type TraceHttpMiddlewareParam struct {
	fx.In
	tracer.Tracer
}

type traceHttpMiddlewareImpl struct {
	tracer tracer.Tracer
}

func NewTraceHttpMiddleware(p TraceHttpMiddlewareParam) TraceHttpMiddleware {
	return &traceHttpMiddlewareImpl{
		tracer: p.Tracer,
	}
}

func (t *traceHttpMiddlewareImpl) ExtractTraceContext(c *gin.Context) {
	// Inject the W3C compliant trace context from the incoming HTTP request into the current trace context
	// This is useful for propagating trace context across services
	t.tracer.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
	c.Next()
}
