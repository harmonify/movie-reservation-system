package http

import (
	auth_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/auth"
	health_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/health_check"
	user_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/user"
	http_interface "github.com/harmonify/movie-reservation-system/user-service/lib/http/interface"
	"go.uber.org/fx"
)

var (
	HttpModule = fx.Module(
		"http",
		fx.Provide(NewHttpServer),
	)
	RestHandlers = []http_interface.RestHandler{
		health_rest.NewHealthCheckRestHandler,
		auth_rest.NewAuthRestHandler,
		user_rest.NewUserRestHandler,
	}
)
