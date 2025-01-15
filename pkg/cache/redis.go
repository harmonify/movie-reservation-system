package cache

import (
	"context"
	"fmt"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/redis/go-redis/v9"
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
	if p.Config.RedisHost == "" {
		return RedisResult{}, fmt.Errorf("redis host is required")
	}
	if p.Config.RedisPort == "" {
		return RedisResult{}, fmt.Errorf("redis port is required")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     p.Config.RedisHost + ":" + p.Config.RedisPort,
		Password: p.Config.RedisPass,
	})

	_, err := client.Ping(context.TODO()).Result()

	return RedisResult{
		Redis: &Redis{
			Client: client,
		},
	}, err
}
