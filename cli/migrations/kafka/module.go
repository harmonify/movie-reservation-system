package kafka_migration

import (
	"github.com/harmonify/movie-reservation-system/cli/shared"
	"go.uber.org/fx"
)

var (
	MigrationModule = fx.Module(
		"migrations",
		fx.Provide(
			AsMigration(NewCreatePublicUserRegisteredV1TopicMigration),
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
