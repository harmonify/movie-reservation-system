package service

import (
	"go.uber.org/fx"
)

var ServiceModule = fx.Module(
	"service",
	fx.Provide(
		NewAdminShowtimeService,
		NewAdminTheaterService,
		NewShowtimeService,
		NewSeatService,
	),
)
