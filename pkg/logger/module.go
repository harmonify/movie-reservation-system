package logger

import (
	"fmt"
	"time"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger(cfg *config.Config) (Logger, error) {
	if cfg.LogLevel == "" {
		return nil, fmt.Errorf("Log level is required")
	}
	if cfg.LogType == "" {
		return nil, fmt.Errorf("Log type is required")
	}

	zapConfig := zap.NewProductionConfig()
	// zapConfig.EncoderConfig.CallerKey = zapcore.OmitKey

	logLevel, err := zap.ParseAtomicLevel(cfg.LogLevel)
	if err == nil {
		zapConfig.Level = logLevel
	} else {
		fmt.Println("Failed to set log level")
	}

	switch cfg.LogType {
	case "nop":
		{
			return NewNopLogger(), nil
		}
	case "console":
		{
			return NewConsoleLogger(), nil
		}
	default:
		{
			if cfg.LokiUrl == "" {
				return nil, fmt.Errorf("Loki URL is required")
			}
			logger, err := NewLokiLogger(zapConfig, LokiConfig{
				Url:          cfg.LokiUrl,
				BatchMaxSize: 1000,
				BatchMaxWait: 10 * time.Second,
				Labels:       map[string]string{"app": cfg.AppName, "env": cfg.Env},
			})

			if err != nil {
				return nil, err
			}

			return logger, nil
		}
	}
}

var (
	LoggerModule = fx.Provide(NewLogger)
)
