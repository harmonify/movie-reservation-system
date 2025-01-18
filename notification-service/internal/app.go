package internal

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/services"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/driven/email/mailgun"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/driven/sms/twilio"
	http_driver "github.com/harmonify/movie-reservation-system/notification-service/internal/driver/http"
	kafkaconsumer "github.com/harmonify/movie-reservation-system/notification-service/internal/driver/kafka_consumer"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/metrics"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
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
		util.UtilModule,

		// CORE
		services.ServiceModule,

		// INFRA (DRIVEN)
		mailgun.MailgunMailerModule,
		twilio.TwilioSmsModule,

		// API (DRIVER)
		http_driver.HttpModule,
		kafkaconsumer.KafkaConsumerModule,
		fx.Provide(
			fx.Annotate(
				kafka.NewKafkaRouter,
				fx.ParamTags(`group:"kafka-routes"`),
			),
		),
	}

	// Override dependencies
	if len(p) > 0 {
		for _, c := range p {
			options = append(options, c)
		}
	}

	return fx.New(options...)
}
