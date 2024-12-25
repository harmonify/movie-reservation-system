package rbac

import (
	"go.uber.org/fx"
)

var (
	DrivenCasbinModule = fx.Module(
		"driven-casbin",
		fx.Provide(
			NewCasbin,
		),
	)
)
