package http_driver

import (
	"github.com/harmonify/movie-reservation-system/notification-service/internal/driven/config"
	health_rest "github.com/harmonify/movie-reservation-system/notification-service/internal/driver/http/health_check"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"go.uber.org/fx"
)

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

func BootstrapHttpServer(h *HttpServer) {}
