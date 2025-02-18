package mongo_repository

import "go.uber.org/fx"

var DrivenMongoModule = fx.Module(
	"driven-mongo",
	fx.Provide(
		NewMovieMongoRepository,
	),
)
