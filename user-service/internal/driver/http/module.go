package http

import (
	auth_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/auth"
	health_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/health_check"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/middleware"

	// user_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/user"
	http_interface "github.com/harmonify/movie-reservation-system/pkg/http/interface"
	"go.uber.org/fx"
)

type RestHandlers = []http_interface.RestHandler

var (
	HttpModule = fx.Module(
		"http",
		middleware.HttpMiddlewareModule,
		fx.Provide(
			health_rest.NewHealthCheckRestHandler,
			auth_rest.NewAuthRestHandler,
			// user_rest.NewUserRestHandler,

			func(
				h health_rest.HealthCheckRestHandler,
				a auth_rest.AuthRestHandler,
				// u user_rest.UserRestHandler,
			) RestHandlers {
				return RestHandlers{
					h, a,
					// u
				}
			},
			NewHttpServer,
		),
	)
)
