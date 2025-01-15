package generator_util

import "go.uber.org/fx"

var (
	GeneratorUtilModule = fx.Module(
		"generator-util",
		fx.Provide(
			NewGeneratorUtil,
		),
	)
)
