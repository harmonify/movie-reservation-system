package internal

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/harmonify/movie-reservation-system/pkg/cache"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/database"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/mail"
	"github.com/harmonify/movie-reservation-system/pkg/messaging"
	"github.com/harmonify/movie-reservation-system/pkg/metrics"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/service"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven"
	http_driver "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http"
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
		error_pkg.ErrorModule,
		util.UtilModule,
		fx.Provide(
			kafka.NewKafkaProducer,
		),

		// CORE
		service.ServiceModule,

		// INFRA (DRIVEN)
		database.DatabaseModule,
		cache.RedisModule,
		mail.MailerModule,
		messaging.MessagingModule,
		driven.DrivenModule,

		// API (DRIVER)
		http_pkg.HttpModule,
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
