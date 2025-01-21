package logger

import (
	"time"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"go.uber.org/fx"
)

func NewLogger(cfg *config.Config) (Logger, error) {
	switch cfg.LogType {
	case "nop":
		{
			return NewNopLogger(), nil
		}
	case "loki":
		{
			logger, err := NewLokiZapLogger(&LokiZapConfig{
				LogLevel:     cfg.LogLevel,
				Url:          cfg.LokiUrl,
				BatchMaxSize: 1000,
				BatchMaxWait: 10 * time.Second,
				// https://grafana.com/docs/loki/latest/get-started/labels/
				Labels: map[string]string{
					"env":          cfg.Env,
					"service_name": cfg.ServiceIdentifier,
				},
			})

			if err != nil {
				return nil, err
			}

			return logger, nil
		}
	default:
		{
			return NewConsoleLogger(), nil
		}
	}
}

var (
	LoggerModule = fx.Provide(NewLogger)
)
