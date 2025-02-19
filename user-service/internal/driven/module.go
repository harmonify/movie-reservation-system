package driven

import (
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/cache/redis"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/rpc/grpc"
	"go.uber.org/fx"
)

var (
	DrivenModule = fx.Module(
		"driven",
		redis.DrivenRedisModule,
		postgresql.DrivenPostgresqlModule,
		grpc.DrivenGrpcModule,
	)
)
