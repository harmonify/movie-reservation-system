package internal

import (
	"context"
	"fmt"

	"github.com/harmonify/movie-reservation-system/movie-service/internal/core/service"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven/config"
	http_driver "github.com/harmonify/movie-reservation-system/movie-service/internal/driver/http"
	"github.com/harmonify/movie-reservation-system/pkg/cache"
	"github.com/harmonify/movie-reservation-system/pkg/database/mongo"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/pkg/util/encryption"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
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
		// LIB
		fx.Provide(
			func(cfg *config.MovieServiceConfig) (logger.Logger, error) {
				return logger.NewLogger(&logger.LoggerConfig{
					Env:               cfg.Env,
					ServiceIdentifier: cfg.ServiceIdentifier,
					LogType:           cfg.LogType,
					LogLevel:          cfg.LogLevel,
					LokiUrl:           cfg.LokiUrl,
				})
			},
			func(lc fx.Lifecycle, cfg *config.MovieServiceConfig) (tracer.Tracer, error) {
				return tracer.NewTracer(lc, &tracer.TracerConfig{
					Env:               cfg.Env,
					ServiceIdentifier: cfg.ServiceIdentifier,
					Type:              cfg.TracerType,
					OtelEndpoint:      cfg.OtelEndpoint,
				})
			},
			func(cfg *config.MovieServiceConfig) *encryption.AESEncryptionConfig {
				return &encryption.AESEncryptionConfig{
					AppSecret: cfg.AppSecret,
				}
			},
			func(cfg *config.MovieServiceConfig) *encryption.SHA256HasherConfig {
				return &encryption.SHA256HasherConfig{
					AppSecret: cfg.AppSecret,
				}
			},
			func(cfg *config.MovieServiceConfig) *jwt_util.JwtUtilConfig {
				return &jwt_util.JwtUtilConfig{
					ServiceIdentifier:      cfg.ServiceIdentifier,
					JwtAudienceIdentifiers: cfg.AuthJwtAudienceIdentifiers,
					JwtIssuerIdentifier:    cfg.AuthJwtIssuerIdentifier,
				}
			},
			func(p mongo.MongoClientParam, cfg *config.MovieServiceConfig) (*mongo.MongoClient, error) {
				return mongo.NewMongoClient(p, &mongo.MongoClientConfig{
					URI:        cfg.MongoUri,
					ReplicaSet: cfg.MongoReplicaSet,
				})
			},
			func(cfg *config.MovieServiceConfig) (*cache.Redis, error) {
				return cache.NewRedis(&cache.RedisConfig{
					RedisHost: cfg.RedisHost,
					RedisPort: cfg.RedisPort,
					RedisPass: cfg.RedisPass,
				})
			},
		),
		error_pkg.ErrorModule,
		util.UtilModule,

		// CORE
		service.ServiceModule,

		// DRIVEN
		driven.DrivenModule,

		// API (DRIVER)
		http_driver.HttpModule,
	}

	// Override dependencies
	if len(p) > 0 {
		for _, c := range p {
			options = append(options, c)
		}
	}

	return fx.New(options...)
}
