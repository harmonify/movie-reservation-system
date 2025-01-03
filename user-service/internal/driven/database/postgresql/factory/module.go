package factory

import (
	"go.uber.org/fx"
)

var (
	DrivenPostgresqlFactoryModule = fx.Module(
		"driven-postgresql-factory",
		fx.Provide(
			NewUserFactory,
			NewUserSessionFactory,
		),
	)
)
