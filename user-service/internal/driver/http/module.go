package http_driver

import (
	"context"

	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	auth_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/auth"
	health_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/health_check"
	http_driver_shared "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/shared"
	user_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/user"
	"go.uber.org/fx"
)

type BootstrapHttpServerParam struct {
	fx.In
	fx.Lifecycle
	Routes     []http_pkg.RestHandler `group:"http_routes"`
	HttpServer *HttpServer
}

var (
	HttpModule = fx.Module(
		"http-driver",
		http_pkg.HttpModule,
		http_driver_shared.HttpMiddlewareModule,
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

func BootstrapHttpServer(p BootstrapHttpServerParam) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := p.HttpServer.configure(p.Routes...); err != nil {
				return err
			}
			return p.HttpServer.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return p.HttpServer.Shutdown(ctx)
		},
	})
}
