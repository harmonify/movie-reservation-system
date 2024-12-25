package struct_util

import "go.uber.org/fx"

var (
	StructUtilModule = fx.Module(
		"struct-util",
		fx.Provide(
			NewStructUtil,
		),
	)
)
