package internal

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/services"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/driven/config"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/driven/email/mailgun"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/driven/sms/twilio"
	grpc_driver "github.com/harmonify/movie-reservation-system/notification-service/internal/driver/grpc"
	http_driver "github.com/harmonify/movie-reservation-system/notification-service/internal/driver/http"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/metrics"
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
			func() (*config.NotificationServiceConfig, error) {
				_, filename, _, _ := runtime.Caller(0)
				configFile := path.Join(filename, "..", "..", ".env")
				return config.NewNotificationServiceConfig(configFile)
			},
			func(cfg *config.NotificationServiceConfig) (logger.Logger, error) {
				return logger.NewLogger(&logger.LoggerConfig{
					Env:               cfg.Env,
					ServiceIdentifier: cfg.ServiceIdentifier,
					LogType:           cfg.LogType,
					LogLevel:          cfg.LogLevel,
					LokiUrl:           cfg.LokiUrl,
				})
			},
			func(lc fx.Lifecycle, cfg *config.NotificationServiceConfig) (tracer.Tracer, error) {
				return tracer.NewTracer(lc, &tracer.TracerConfig{
					Env:               cfg.Env,
					ServiceIdentifier: cfg.ServiceIdentifier,
					Type:              cfg.TracerType,
					OtelEndpoint:      cfg.OtelEndpoint,
				})
			},
			func(cfg *config.NotificationServiceConfig) *encryption.AESEncryptionConfig {
				return &encryption.AESEncryptionConfig{
					AppSecret: cfg.AppSecret,
				}
			},
			func(cfg *config.NotificationServiceConfig) *jwt_util.JwtUtilConfig {
				return &jwt_util.JwtUtilConfig{
					ServiceIdentifier:      cfg.ServiceIdentifier,
					JwtAudienceIdentifiers: cfg.AuthJwtAudienceIdentifiers,
					JwtIssuerIdentifier:    cfg.AuthJwtIssuerIdentifier,
				}
			},
		),
		error_pkg.ErrorModule,
		util.UtilModule,
		metrics.MetricsModule,

		// CORE
		services.ServiceModule,

		// INFRA (DRIVEN)
		mailgun.MailgunMailerModule,
		twilio.TwilioSmsModule,

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
