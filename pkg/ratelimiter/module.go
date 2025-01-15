package ratelimiter

import "go.uber.org/fx"

var RateLimiterModule = fx.Option(
	fx.Provide(
		NewRateLimiterRegistry,
	),
)
