package postgresql

import (
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/repository"
	"go.uber.org/fx"
)

var (
	DrivenPostgresqlModule = fx.Module(
		"driven-postgresql",
		fx.Provide(
			repository.NewUserRepository,
			repository.NewUserKeyRepository,
			repository.NewUserSessionRepository,
		),
	)
)
