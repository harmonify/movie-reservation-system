package kafka_driver

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	watermill_kafka "github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	watermill_pkg "github.com/harmonify/movie-reservation-system/pkg/kafka/watermill"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	"go.uber.org/fx"
)

var (
	registeredTopics = []string{
		shared.PublicUserRegisteredV1.String(),
	}

	KafkaConsumerModule = fx.Module(
		"kafka-driver",
		fx.Provide(
			func(logger logger.Logger) watermill.LoggerAdapter {
				return watermill_pkg.NewLogger(logger.GetZapLogger())
			},
			watermill_pkg.AsRoute(NewUserRegisteredRoute),
		),
		fx.Invoke(
			watermill_pkg.AsRouter(BootstrapWatermill),
		),
	)
)

func BootstrapWatermill(routes []watermill_pkg.Route, lc fx.Lifecycle, wl watermill.LoggerAdapter, cfg *config.UserServiceConfig) {
	lc.Append(fx.StartHook(func(ctx context.Context) error {
		router, err := message.NewRouter(message.RouterConfig{}, wl)
		if err != nil {
			wl.Error("Failed to initiate router", err, watermill.LogFields{})
		}

		kafkaConfig, err := kafka.BuildKafkaConfig(&kafka.KafkaConfig{
			KafkaVersion: cfg.KafkaVersion,
		})
		if err != nil {
			wl.Error("Failed to build kafka consumer config", err, watermill.LogFields{})
		}

		subscriber, err := watermill_kafka.NewSubscriber(
			watermill_kafka.SubscriberConfig{
				Brokers:               []string{cfg.KafkaBrokers},
				ConsumerGroup:         cfg.KafkaConsumerGroup,
				OverwriteSaramaConfig: kafkaConfig,
				Unmarshaler:           watermill_pkg.NewWatermillMarshaler(),
				Tracer:                watermill_kafka.NewOTELSaramaTracer(),
			},
			wl,
		)

		// promRegistry, closeMetricsServer := metrics.CreateRegistryAndServeHTTP(metricsAddr)
		// defer closeMetricsServer()

		// metricsBuilder := metrics.NewPrometheusMetricsBuilder(promRegistry, "demo", "hello")
		// metricsBuilder.AddPrometheusRouterMetrics(router)

		// SignalsHandler gracefully shutdowns Router when receiving SIGTERM
		router.AddPlugin(plugin.SignalsHandler)

		// Router level middleware are executed for every message sent to the router
		router.AddMiddleware(
			// TODO add poison queue middleware in handler level
			// Recoverer handles panics from handlers
			middleware.Recoverer,
		)

		// Registering routes
		for _, route := range routes {
			wl.Debug("Registering route", watermill.LogFields{"route": route.Identifier()})
			err := route.Register(router, subscriber)
			if err != nil {
				wl.Error("Failed to register route", err, watermill.LogFields{"route": route.Identifier()})
			}
		}

		// Run is blocking while the router is running.
		if err := router.Run(ctx); err != nil {
			return err
		}

		return nil
	}))
}
