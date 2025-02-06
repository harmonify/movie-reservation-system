package http_driver_shared

import (
	http_middleware "github.com/harmonify/movie-reservation-system/pkg/http/middleware"
	"github.com/harmonify/movie-reservation-system/pkg/metrics"
	"github.com/harmonify/movie-reservation-system/pkg/ratelimiter"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	"go.uber.org/fx"
)

type (
	HttpMiddleware struct {
		Recovery    http_middleware.RecoveryHttpMiddleware
		Metrics     metrics.PrometheusHttpMiddleware
		Auth        http_middleware.JwtAuthenticationHttpMiddleware
		Rbac        http_middleware.JwtRbacHttpMiddleware
		RateLimiter http_middleware.RateLimiterHttpMiddleware
	}
)

var HttpMiddlewareModule = fx.Module(
	"http-middleware",
	fx.Provide(
		http_middleware.NewRecoveryHttpMiddleware,
		metrics.NewPrometheusHttpMiddleware,
		http_middleware.NewJwtHttpMiddleware,
		func(p http_middleware.JwtRbacHttpMiddlewareParam, cfg *config.UserServiceConfig) (http_middleware.JwtRbacHttpMiddlewareResult, error) {
			return http_middleware.NewJwtRbacHttpMiddleware(p, &http_middleware.JwtHttpMiddlewareConfig{
				Domain: cfg.ServiceIdentifier,
			})
		},
		func(p ratelimiter.RateLimiterRegistryParam, cfg *config.UserServiceConfig) (ratelimiter.RateLimiterRegistry, error) {
			return ratelimiter.NewRateLimiterRegistry(p, &ratelimiter.RateLimiterRegistryConfig{
				ServiceIdentifier: cfg.ServiceIdentifier,
			})
		},
		http_middleware.NewRateLimiterHttpMiddleware,
		NewHttpMiddleware,
	),
)

func NewHttpMiddleware(
	recovery http_middleware.RecoveryHttpMiddleware,
	metrics metrics.PrometheusHttpMiddleware,
	jwt http_middleware.JwtAuthenticationHttpMiddleware,
	jwtRbac http_middleware.JwtRbacHttpMiddleware,
	rateLimiter http_middleware.RateLimiterHttpMiddleware,
) *HttpMiddleware {
	return &HttpMiddleware{
		Recovery:    recovery,
		Metrics:     metrics,
		Auth:        jwt,
		Rbac:        jwtRbac,
		RateLimiter: rateLimiter,
	}
}
