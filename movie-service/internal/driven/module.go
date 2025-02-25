package driven

import (
	"path"
	"runtime"

	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven/config"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven/database/mongo"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven/rpc/grpc"
	"go.uber.org/fx"
)

var (
	DrivenModule = fx.Module(
		"driven",
		fx.Provide(
			func() (*config.MovieServiceConfig, error) {
				_, filename, _, _ := runtime.Caller(0)
				configFile := path.Join(path.Dir(filename), "..", "..", ".env")
				return config.NewMovieServiceConfig(configFile)
			},
		),
		mongo.DrivenMongoModule,
		grpc.DrivenGrpcModule,
	)
)
