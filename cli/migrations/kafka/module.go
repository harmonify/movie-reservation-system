package kafka

import (
	v_1_0_0 "github.com/harmonify/movie-reservation-system/cli/migrations/kafka/v1.0.0"
	"github.com/harmonify/movie-reservation-system/cli/shared"

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
