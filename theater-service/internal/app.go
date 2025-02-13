package internal

import (
	"context"
	"fmt"

	"github.com/harmonify/movie-reservation-system/pkg/cache"
	"github.com/harmonify/movie-reservation-system/pkg/database"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/metrics"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/pkg/util/encryption"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/services"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/driven"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/driven/config"
	grpc_driver "github.com/harmonify/movie-reservation-system/theater-service/internal/driver/grpc"
	http_driver "github.com/harmonify/movie-reservation-system/theater-service/internal/driver/http"
	"go.uber.org/fx"
)

func StartApp() error {
	app := NewApp()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Start(ctx); err != nil {
		fmt.Println(">> App failed to start. Error:", err)
		return err
	}

	<-app.Done()
	fmt.Println(">> App shutdown")
	return nil
}

// This is a function to initialize all services and invoke their functions.
func NewApp(p ...fx.Option) *fx.App {
	options := []fx.Option{
		driven.DrivenModule,

		// LIB
		fx.Provide(
			func(cfg *config.TheaterServiceConfig) (logger.Logger, error) {
				return logger.NewLogger(&logger.LoggerConfig{
					Env:               cfg.Env,
					ServiceIdentifier: cfg.ServiceIdentifier,
					LogType:           cfg.LogType,
					LogLevel:          cfg.LogLevel,
					LokiUrl:           cfg.LokiUrl,
				})
			},
			func(lc fx.Lifecycle, cfg *config.TheaterServiceConfig) (tracer.Tracer, error) {
				return tracer.NewTracer(lc, &tracer.TracerConfig{
					Env:               cfg.Env,
					ServiceIdentifier: cfg.ServiceIdentifier,
					Type:              cfg.TracerType,
					OtelEndpoint:      cfg.OtelEndpoint,
				})
			},
			func(cfg *config.TheaterServiceConfig) *encryption.AESEncryptionConfig {
				return &encryption.AESEncryptionConfig{
					AppSecret: cfg.AppSecret,
				}
			},
			func(cfg *config.TheaterServiceConfig) *encryption.SHA256HasherConfig {
				return &encryption.SHA256HasherConfig{
					AppSecret: cfg.AppSecret,
				}
			},
			func(cfg *config.TheaterServiceConfig) *jwt_util.JwtUtilConfig {
				return &jwt_util.JwtUtilConfig{
					ServiceIdentifier:      cfg.ServiceIdentifier,
					JwtAudienceIdentifiers: cfg.AuthJwtAudienceIdentifiers,
					JwtIssuerIdentifier:    cfg.AuthJwtIssuerIdentifier,
				}
			},
			func(p database.DatabaseParam, cfg *config.TheaterServiceConfig) (database.DatabaseResult, error) {
				return database.NewDatabase(p, &database.DatabaseConfig{
					Env:                   cfg.Env,
					DbType:                cfg.DbType,
					DbHost:                cfg.DbHost,
					DbPort:                cfg.DbPort,
					DbUser:                cfg.DbUser,
					DbPassword:            cfg.DbPassword,
					DbName:                cfg.DbName,
					DbMaxIdleConn:         cfg.DbMaxIdleConn,
					DbMaxOpenConn:         cfg.DbMaxOpenConn,
					DbMaxLifetimeInMinute: cfg.DbMaxLifetimeInMinute,
				})
			},
			func(cfg *config.TheaterServiceConfig) (*cache.Redis, error) {
				return cache.NewRedis(&cache.RedisConfig{
					RedisHost: cfg.RedisHost,
					RedisPort: cfg.RedisPort,
					RedisPass: cfg.RedisPass,
				})
			},
		),
		error_pkg.ErrorModule,
		util.UtilModule,
		metrics.MetricsModule,

		// CORE
		services.ServiceModule,

		// API (DRIVER)
		http_driver.HttpModule,
		grpc_driver.GrpcModule,
	}

	// Override dependencies
	if len(p) > 0 {
		for _, c := range p {
			options = append(options, c)
		}
	}

	return fx.New(options...)
}
