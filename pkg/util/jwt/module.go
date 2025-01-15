package jwt_util

import "go.uber.org/fx"

var (
	JWTUtilModule = fx.Module(
		"jwt-util",
		fx.Provide(
			NewJwtUtil,
		),
	)
)
