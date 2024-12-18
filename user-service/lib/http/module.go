package http

import (
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/middleware"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/response"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/validator"
	"go.uber.org/fx"
)

var (
	HttpModule = fx.Provide(
		response.NewHttpErrorHandler,
		response.NewHttpResponse,
		validator.NewHttpValidator,
		middleware.NewJWTMiddleware,
	)
)
