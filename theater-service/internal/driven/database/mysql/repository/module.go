package repository

import "go.uber.org/fx"

var DrivenMysqlRepositoryModule = fx.Module(
	"driven-mysql-repository",
	fx.Provide(
		NewTheaterRepository,
		NewRoomRepository,
		NewSeatRepository,
		NewShowtimeRepository,
		NewTicketRepository,
	),
)
