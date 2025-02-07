package logger

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type LoggerConfig struct {
	Env               string `validate:"required,oneof=dev test prod"`
	ServiceIdentifier string `validate:"required"`
	LogType           string `validate:"required,oneof=nop loki console"`
	LogLevel          string `validate:"required,oneof=debug info warn error"`
	LokiUrl           string `validate:"required_if=LogType loki"`
}

func NewLogger(cfg *LoggerConfig) (Logger, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return nil, err
	}

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
