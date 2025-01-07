package shared

import "go.uber.org/fx"

var (
	SharedModule = fx.Module(
		"shared",
		fx.Provide(
			NewConfig,
			NewKafkaAdminClient,
			NewMigrationStorage,
		),
	)
)
