package driven

import (
	"path"
	"runtime"

	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/driven/config"
	mongo_repository "github.com/harmonify/movie-reservation-system/movie-search-service/internal/driven/database/mongo"
	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/driven/rpc/grpc"
	"go.uber.org/fx"
)

var (
	DrivenModule = fx.Module(
		"driven",
		fx.Provide(
			func() (*config.MovieSearchServiceConfig, error) {
				_, filename, _, _ := runtime.Caller(0)
				configFile := path.Join(filename, "..", "..", "..", ".env")
				return config.NewMovieServiceConfig(configFile)
			},
		),
		mongo_repository.DrivenMongoModule,
		grpc.DrivenGrpcModule,
	)
)
