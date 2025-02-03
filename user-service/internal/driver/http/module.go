package http_driver

import (
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	auth_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/auth"
	health_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/health_check"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/middleware"
	user_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/user"
	"go.uber.org/fx"
)

var (
	HttpModule = fx.Module(
		"http-driver",
		http_pkg.HttpModule,
		middleware.HttpMiddlewareModule,
		fx.Provide(
			health_rest.NewHealthCheckRestHandler,
			auth_rest.NewAuthRestHandler,
			user_rest.NewUserRestHandler,
			func(p HttpServerParam, cfg *config.UserServiceConfig) (HttpServerResult, error) {
				return NewHttpServer(p, &HttpServerConfig{
					Env:                     cfg.Env,
					ServiceIdentifier:       cfg.ServiceIdentifier,
					ServiceHttpPort:         cfg.ServiceHttpPort,
					ServiceHttpBaseUrl:      cfg.ServiceHttpBaseUrl,
					ServiceHttpBasePath:     cfg.ServiceHttpBasePath,
					ServiceHttpReadTimeOut:  cfg.ServiceHttpReadTimeOut,
					ServiceHttpWriteTimeOut: cfg.ServiceHttpWriteTimeOut,
					ServiceHttpEnableCors:   cfg.ServiceHttpEnableCors,
				})
			},
		),
		fx.Invoke(BootstrapHttpServer),
	)
)

func BootstrapHttpServer(h *HttpServer) {}
