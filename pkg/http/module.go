package http

import (
	"github.com/harmonify/movie-reservation-system/pkg/http/response"
	"go.uber.org/fx"
)

var (
	HttpModule = fx.Provide(
		response.NewErrorHandler,
		response.NewResponse,
	)
)
