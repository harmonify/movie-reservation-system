package opa

import (
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	"go.uber.org/fx"
)

var DrivenOpaModule = fx.Module(
	"driven-opa",
	fx.Provide(
		func(p OpaClientParam, cfg *config.UserServiceConfig) OpaClientResult {
			return NewOpaClient(p, &OpaClientConfig{
				OpaServerUrl: cfg.OpaServerUrl,
			})
		},
	),
)
