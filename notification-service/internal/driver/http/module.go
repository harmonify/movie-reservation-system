package http

import (
	health_rest "github.com/harmonify/movie-reservation-system/notification-service/internal/driver/http/health_check"
	"go.uber.org/fx"
)

var (
	HttpModule = fx.Module(
		"http",
		fx.Provide(
			health_rest.NewHealthCheckRestHandler,
			NewHttpServer,
		),
	)
)
