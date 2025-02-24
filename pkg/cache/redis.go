package cache

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var RedisModule = fx.Module("redis", fx.Provide(NewRedis))

type Redis struct {
	Client *redis.Client
}

type RedisConfig struct {
	RedisHost string `validate:"required"`
	RedisPort string `validate:"required,min=1,max=65535"`
	RedisPass string `validate:"required"`
	RedisDB   int    `validate:"min=0,max=15"`
}

func NewRedis(cfg *RedisConfig) (*Redis, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		return nil, err
	}
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		return nil, err
	}

	_, err := rdb.Ping(context.TODO()).Result()

	return &Redis{
		Client: rdb,
	}, err
}
