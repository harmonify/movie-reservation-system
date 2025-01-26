package error_pkg

import (
	"go.uber.org/fx"
)

var (
	ErrorModule = fx.Module(
		"error_pkg",
		fx.Provide(NewErrorMapper),
	)
)
