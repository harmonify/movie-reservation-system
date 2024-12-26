package database

import (
	"go.uber.org/fx"
)

var DatabaseModule = fx.Module(
	"database",
	fx.Provide(
		NewDatabase,
		NewPostgresqlErrorTranslator,
	),
)

func NewDatabase(p DatabaseParam) (DatabaseResult, error) {
	return newPostgresqlDatabase(p)
}
