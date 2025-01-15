package http_util

import "go.uber.org/fx"

var (
	HttpUtilModule = fx.Module(
		"http-util",
		fx.Provide(
			NewHttpUtil,
		),
	)
)
