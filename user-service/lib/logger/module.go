package logger

import (
	"context"

	"fmt"
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	GetZapLogger() *zap.Logger
	WithCtx(ctx context.Context) Logger
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Log(debugLevel zapcore.Level, msg string, fields ...zap.Field)
}

func NewLogger(cfg config.Config) Logger {
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
			return NewConsoleLogger()
		}
	default:
		{
			logger, err := NewLokiLogger(zapConfig, LokiConfig{
				Url:          cfg.LokiUrl,
				BatchMaxSize: 1000,
				BatchMaxWait: 10 * time.Second,
				Labels:       map[string]string{"app": cfg.AppName, "env": cfg.Env},
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
