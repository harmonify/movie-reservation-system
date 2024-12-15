package config

import "go.uber.org/fx"

var ConfigModule = fx.Options(
	fx.Provide(LoadConfig),
)
