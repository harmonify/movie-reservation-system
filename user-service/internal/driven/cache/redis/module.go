package redis

import "go.uber.org/fx"

var (
	DrivenRedisModule = fx.Module(
		"driven-redis",
		fx.Provide(
			NewOtpRedisRepository,
		),
	)
)
