package mongo

import (
	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven/database/mongo/repository"
	"go.uber.org/fx"
)

var DrivenMongoModule = fx.Module(
	"driven-mongo",
	fx.Provide(
		mongo_repository.NewMovieMongoRepository,
	),
)
