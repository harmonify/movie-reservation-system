package entityfactory

import "go.uber.org/fx"

var UserEntityFactoryModule = fx.Module(
	"user-entity-factory",
	fx.Provide(
		NewUserFactory,
		NewUserKeyFactory,
		NewUserSessionFactory,
	),
)
