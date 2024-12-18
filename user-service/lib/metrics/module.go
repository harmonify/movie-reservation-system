package metrics

import "go.uber.org/fx"

var (
	MetricsModule = fx.Module("metrics", fx.Provide(NewPrometheusHttpMiddleware))
)
