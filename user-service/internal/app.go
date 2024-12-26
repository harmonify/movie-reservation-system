package internal

import (
	"context"
	"fmt"
	"maps"
	"path"
	"runtime"

	"github.com/harmonify/movie-reservation-system/user-service/internal/core/service"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven"
	http_driver "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http"
	"github.com/harmonify/movie-reservation-system/user-service/lib/cache"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http"
	http_constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/mail"
	"github.com/harmonify/movie-reservation-system/user-service/lib/metrics"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"go.uber.org/fx"
)

func StartApp() error {
	app := NewApp(
		fx.Invoke(Bootstrap),
	)

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
		fx.Provide(
			func() *config.ConfigFile {
				_, filename, _, _ := runtime.Caller(0)
				return &config.ConfigFile{
					Path: path.Join(filename, "..", "..", ".env"),
				}
			},
		),
		config.ConfigModule,

		// Libraries
		logger.LoggerModule,
		tracer.TracerModule,
		metrics.MetricsModule,
		util.UtilModule,

		// CORE
		service.ServiceModule,

		// INFRA (DRIVEN)
		database.DatabaseModule,
		cache.RedisModule,
		mail.MailerModule,
		driven.DrivenModule,

		// API (DRIVER)
		fx.Provide(
			func() *http_constant.CustomHttpErrorMap {
				maps.Copy(http_constant.DefaultCustomHttpErrorMap, auth_service.AuthServiceErrorMap)
				return &http_constant.DefaultCustomHttpErrorMap
			},
		),
		http.HttpModule,
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

func Bootstrap(lc fx.Lifecycle, l logger.Logger, h *http_driver.HttpServer, t tracer.Tracer, handlers http_driver.RestHandlers) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := h.Start(ctx, handlers...)
			if err != nil {
				l.WithCtx(ctx).Error(err.Error())
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := h.Shutdown(ctx)
			if err != nil {
				l.WithCtx(ctx).Error(err.Error())
				return err
			}

			err = t.Shutdown(ctx)
			if err != nil {
				l.WithCtx(ctx).Error(err.Error())
				return err
			}

			return nil
		},
	})
}
