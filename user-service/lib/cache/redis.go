package cache

import (
	"github.com/go-redis/redis"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"go.uber.org/fx"
)

var RedisModule = fx.Module("redis", fx.Provide(NewRedis))

type Redis struct {
	Client *redis.Client
}

type RedisParam struct {
	fx.In

	Config *config.Config
}

type RedisResult struct {
	fx.Out

	Redis *Redis
}

func NewRedis(p RedisParam) (RedisResult, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     p.Config.RedisHost + ":" + p.Config.RedisPort,
		Password: p.Config.RedisPass,
	})

	_, err := client.Ping().Result()

	return RedisResult{
		Redis: &Redis{
			Client: client,
		},
	}, err
}
