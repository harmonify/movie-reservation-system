package kafka_driver

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	watermill_kafka "github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/kafka"
	watermill_pkg "github.com/harmonify/movie-reservation-system/pkg/kafka/watermill"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	"go.uber.org/fx"
)

var (
	KafkaConsumerModule = fx.Module(
		"kafka-driver",
		fx.Provide(
			func(cfg *config.UserServiceConfig, logger logger.Logger) watermill.LoggerAdapter {
				zapLogger := logger.GetZapLogger()
				return watermill_pkg.NewLogger(zapLogger)
			},
			watermill_pkg.AsRoute(NewUserRegisteredRoute),
			watermill_pkg.AsRouter(NewRouter),
		),
		fx.Invoke(BootstrapWatermill),
	)
)

func BootstrapWatermill(cfg *config.UserServiceConfig, r *Router, lc fx.Lifecycle) {
	// Disable kafka client in test environment
	if cfg.Env == config_pkg.EnvironmentTest {
		return
	}
	lc.Append(fx.StartStopHook(r.Start, r.Close))
}

type Router struct {
	router *message.Router
	logger logger.Logger
}

func NewRouter(routes []watermill_pkg.Route, wl watermill.LoggerAdapter, cfg *config.UserServiceConfig, l logger.Logger) *Router {
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

	r := &Router{
		router: router,
		logger: l,
	}

	return r
}

func (r *Router) Start(ctx context.Context) error {
	var err error

	go func() {
		err = r.router.Run(ctx)
	}()

	time.Sleep(1 * time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (r *Router) Close(ctx context.Context) error {
	r.logger.WithCtx(ctx).Info("Closing router")
	return r.router.Close()
}
