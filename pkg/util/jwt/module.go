package jwt_util

import "go.uber.org/fx"

var (
	JWTUtilModule = fx.Provide(NewJWTUtil)
)
