package http_driver

import (
	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/driven/config"
	health_rest "github.com/harmonify/movie-reservation-system/theater-service/internal/driver/http/health_check"
	http_driver_shared "github.com/harmonify/movie-reservation-system/theater-service/internal/driver/http/shared"
	showtime_rest "github.com/harmonify/movie-reservation-system/theater-service/internal/driver/http/showtime"
	theater_rest "github.com/harmonify/movie-reservation-system/theater-service/internal/driver/http/theater"
	"go.uber.org/fx"
)

type BootstrapHttpServerParam struct {
	fx.In
	fx.Lifecycle
	Config     *config.TheaterServiceConfig
	HttpServer *HttpServer
}

var (
	HttpModule = fx.Module(
		"http-driver",
		http_pkg.HttpModule,
		http_driver_shared.HttpMiddlewareModule,
		fx.Provide(
			health_rest.NewHealthCheckRestHandler,
			theater_rest.NewAdminTheaterRestHandler,
			showtime_rest.NewAdminShowtimeRestHandler,
			func(p HttpServerParam, cfg *config.TheaterServiceConfig) (HttpServerResult, error) {
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
	// Disable http server in test environment
	if p.Config.Env == config_pkg.EnvironmentTest {
		return
	}
	p.Lifecycle.Append(fx.StartStopHook(p.HttpServer.Start, p.HttpServer.Shutdown))
}
