package logger

import (
	"fmt"
	"time"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	console_logger "github.com/harmonify/movie-reservation-system/pkg/logger/console"
	loki_logger "github.com/harmonify/movie-reservation-system/pkg/logger/loki"
	logger_shared "github.com/harmonify/movie-reservation-system/pkg/logger/shared"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger(cfg config.Config) logger_shared.Logger {
	zapConfig := zap.NewProductionConfig()
	// zapConfig.EncoderConfig.CallerKey = zapcore.OmitKey

	logLevel, err := zap.ParseAtomicLevel(cfg.LogLevel)
	if err == nil {
		zapConfig.Level = logLevel
	} else {
		fmt.Println("Failed to set log level")
	}

	switch cfg.LogType {
	case "console":
		{
			return console_logger.NewConsoleLogger()
		}
	default:
		{
			logger, err := loki_logger.NewLokiLogger(zapConfig, logger_shared.LokiConfig{
				Url:          cfg.LokiUrl,
				BatchMaxSize: 1000,
				BatchMaxWait: 10 * time.Second,
				Labels:       map[string]string{"app": cfg.ServiceName, "env": cfg.Env},
			})

			if err != nil {
				fmt.Printf("%s", err)
			}

			return logger
		}
	}
}

var (
	LoggerModule = fx.Provide(NewLogger)
)
