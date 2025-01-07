package kafka

import (
	v_1_0_0 "kafka-playground/migrations/kafka/v1.0.0"
	"kafka-playground/shared"

	"go.uber.org/fx"
)

var (
	MigrationModule = fx.Module(
		"migrations",
		fx.Provide(
			AsMigration(v_1_0_0.NewCreateNewOrderTopicMigration),
			// fx.Annotate(
			// ,
			// fx.ParamTags(`group:"migrations"`)
			// ),
		),
	)
)

func AsMigration(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(shared.KafkaMigration)),
		fx.ResultTags(`group:"migrations"`),
	)
}
