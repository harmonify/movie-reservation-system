package seeder

import (
	"go.uber.org/fx"
)

var (
	DrivenPostgresqlSeederModule = fx.Module(
		"driven-postgresql-seeder",
		fx.Provide(
			NewUserSeeder,
			NewUserKeySeeder,
			NewUserSessionSeeder,
		),
	)
)
