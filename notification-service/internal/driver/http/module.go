package http_driver

import (
	"context"

	"github.com/harmonify/movie-reservation-system/notification-service/internal/driven/config"
	health_rest "github.com/harmonify/movie-reservation-system/notification-service/internal/driver/http/health_check"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
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
		fx.Provide(
			health_rest.NewHealthCheckRestHandler,
			func(p HttpServerParam, cfg *config.NotificationServiceConfig) (HttpServerResult, error) {
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
