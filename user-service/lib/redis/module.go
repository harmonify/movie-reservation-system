package app

import (
	"github.com/go-redis/redis"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"go.uber.org/fx"
)

var RedisModule = fx.Module("redis", fx.Provide(CreateRedisConnection))

func CreateRedisConnection(cfg *config.Config) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPass,
	})

	_, err = client.Ping().Result()
	return
}
