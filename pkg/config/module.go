package config

import "go.uber.org/fx"

var ConfigModule = fx.Provide(NewConfig)
