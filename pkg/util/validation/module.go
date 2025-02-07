package validation

import "go.uber.org/fx"

var (
	ValidationUtilModule = fx.Module(
		"validation-util",
		fx.Provide(
			NewValidator,
			NewStructValidator,
		),
	)
)
