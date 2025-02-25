package driven

import (
	"path"
	"runtime"

	"github.com/harmonify/movie-reservation-system/theater-service/internal/driven/config"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/driven/database/mysql/repository"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/driven/grpc"
	"go.uber.org/fx"
)

var (
	DrivenModule = fx.Module(
		"driven",
		fx.Provide(
			func() (*config.TheaterServiceConfig, error) {
				_, filename, _, _ := runtime.Caller(0)
				configFile := path.Join(filename, "..", "..", "..", ".env")
				return config.NewTheaterServiceConfig(configFile)
			},
		),
		repository.DrivenMysqlRepositoryModule,
		grpc.DrivenGrpcModule,
	)
)
