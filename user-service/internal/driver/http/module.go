package http_driver

import (
	auth_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/auth"
	health_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/health_check"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/middleware"
	user_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/user"
	"go.uber.org/fx"
)

var (
	HttpModule = fx.Module(
		"http-driver",
		middleware.HttpMiddlewareModule,
		fx.Provide(
			health_rest.NewHealthCheckRestHandler,
			auth_rest.NewAuthRestHandler,
			user_rest.NewUserRestHandler,
			NewHttpServer,
		),
		fx.Invoke(BootstrapHttpServer),
	)
)

func BootstrapHttpServer(h *HttpServer) {}
