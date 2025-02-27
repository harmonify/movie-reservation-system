package service

import (
	"go.uber.org/fx"
)

var ServiceModule = fx.Module(
	"service",
	fx.Provide(
		NewAdminTheaterService,
		NewAdminShowtimeService,
		NewTheaterService,
		NewShowtimeService,
		NewSeatService,
	),
)
