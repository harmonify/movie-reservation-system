package services

import (
	"go.uber.org/fx"
)

var ServiceModule = fx.Module(
	"service",
	fx.Provide(
		NewTheaterService,
	),
)
