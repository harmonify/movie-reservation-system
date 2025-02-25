package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/harmonify/movie-reservation-system/pkg/database"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/pkg/util/encryption"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	entityfactory "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/factory"
	entityseeder "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/seeder"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/seeder"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	if err := newMinimalApp().Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func newMinimalApp() *fx.App {
	return fx.New(
		fx.NopLogger,
		entityfactory.UserEntityFactoryModule,
		postgresql.DrivenPostgresqlModule,
		seeder.DrivenPostgresqlSeederModule,
		util.UtilModule,
		fx.Provide(
			func() (*config.UserServiceConfig, error) {
				_, filename, _, _ := runtime.Caller(0)
				configFile := path.Join(path.Dir(filename), "..", "..", ".env")
				return config.NewUserServiceConfig(configFile)
			},
			func(cfg *config.UserServiceConfig) (logger.Logger, error) {
				return logger.NewLogger(&logger.LoggerConfig{
					Env:               cfg.Env,
					ServiceIdentifier: cfg.ServiceIdentifier,
					LogType:           "console",
					LogLevel:          cfg.LogLevel,
					LokiUrl:           cfg.LokiUrl,
				})
			},
			func(lc fx.Lifecycle, cfg *config.UserServiceConfig) (tracer.Tracer, error) {
				return tracer.NewTracer(lc, &tracer.TracerConfig{
					Env:               cfg.Env,
					ServiceIdentifier: cfg.ServiceIdentifier,
					Type:              cfg.TracerType,
					OtelEndpoint:      cfg.OtelEndpoint,
				})
			},
			func(cfg *config.UserServiceConfig) *encryption.AESEncryptionConfig {
				return &encryption.AESEncryptionConfig{
					AppSecret: cfg.AppSecret,
				}
			},
			func(cfg *config.UserServiceConfig) *encryption.SHA256HasherConfig {
				return &encryption.SHA256HasherConfig{
					AppSecret: cfg.AppSecret,
				}
			},
			func(cfg *config.UserServiceConfig) *jwt_util.JwtUtilConfig {
				return &jwt_util.JwtUtilConfig{
					ServiceIdentifier:      cfg.ServiceIdentifier,
					JwtAudienceIdentifiers: cfg.AuthJwtAudienceIdentifiers,
					JwtIssuerIdentifier:    cfg.AuthJwtIssuerIdentifier,
				}
			},
			func(p database.DatabaseParam, cfg *config.UserServiceConfig) (database.DatabaseResult, error) {
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
		),
		fx.Invoke(func(lc fx.Lifecycle, logger logger.Logger, s entityseeder.UserSeeder) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info("Application started")
					if err := run(ctx, logger, s, os.Args...); err != nil {
						logger.Error(err.Error(), zap.Stack("stack"))
						return err
					}
					return nil
				},
				OnStop: func(context.Context) error {
					logger.Info("Application stopped")
					return nil
				},
			})
		}),
	)
}

func run(ctx context.Context, logger logger.Logger, s entityseeder.UserSeeder, args ...string) error {
	username := args[1]
	if username == "" {
		return fmt.Errorf("username is required")
	}

	u, err := s.CreateAdmin(ctx, username)
	if err != nil {
		logger.Error(err.Error(), zap.Stack("stack"))
		return err
	}
	fmt.Println("===============================================================")
	fmt.Printf("ADMIN USERNAME: %v\n", u.User.Username)
	fmt.Printf("ADMIN PASSWORD: %v\n", u.UserRaw.Password)
	fmt.Println("===============================================================")
	return nil
}
