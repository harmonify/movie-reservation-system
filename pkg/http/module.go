package http

import (
	"github.com/harmonify/movie-reservation-system/pkg/http/middleware"
	"github.com/harmonify/movie-reservation-system/pkg/http/response"
	"github.com/harmonify/movie-reservation-system/pkg/http/validator"
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
