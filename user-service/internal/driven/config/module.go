package config

import "go.uber.org/fx"

var UserServiceConfigModule = fx.Provide(NewUserServiceConfig)
