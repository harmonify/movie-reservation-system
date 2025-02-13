package http_pkg

import (
	"go.uber.org/fx"
)

var (
	HttpModule = fx.Provide(
		NewHttpResponseBuilder,
		NewHttpResponse,
		NewHttpValidator,
	)
)
