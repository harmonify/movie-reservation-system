package config

import "go.uber.org/fx"

var TheaterServiceConfigModule = fx.Module(
	"theater-service-config",
	fx.Provide(NewTheaterServiceConfig),
)
