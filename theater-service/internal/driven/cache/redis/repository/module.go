package redis_repository

import "go.uber.org/fx"

var DrivenRedisRepositoryModule = fx.Module(
	"driven-redis-repository",
	fx.Provide(NewMovieRedisRepository),
)
