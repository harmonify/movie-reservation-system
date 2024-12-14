package cmd

import (
	"context"

	"github.com/harmonify/movie-reservation-system/user-service/api/grpc"
	"github.com/harmonify/movie-reservation-system/user-service/api/rest"
	"github.com/harmonify/movie-reservation-system/user-service/cmd/config"

	"github.com/harmonify/movie-reservation-system/pkg/utility"
	"github.com/harmonify/movie-reservation-system/pkg/utility/logger"
	"github.com/harmonify/movie-reservation-system/pkg/utility/sentry"
	"github.com/harmonify/movie-reservation-system/pkg/utility/tracer"

	"github.com/harmonify/movie-reservation-system/pkg/app"
	user "github.com/harmonify/movie-reservation-system/user-service/domain/user"
	"github.com/harmonify/movie-reservation-system/user-service/infrastructure/grpc_client"
	"github.com/harmonify/movie-reservation-system/user-service/infrastructure/repository"
	"go.uber.org/fx"
)

func newServer() {
	return fx.New(
		// pkg.app
		configModule,
		tracer.TracerModule,
		logger_shared.LoggerModule,

		app.HTTPModule,
		app.GRPCModule,
		app.MongoDBModule,
		app.ElasticsearchModule,
		app.RedisModule,
		app.EnterpriseGRPCClientModule,

		utility.Module,

		grpc_client.Module,
		repository.Module,
		rest.Module,
		grpc.Module,

		// Domain Module
		user.Module,

		// Invoke the function
		fx.Invoke(config.LoadConfig),
		fx.Invoke(tracer.InitTracer),
		fx.Invoke(logger.NewLogger),
		fx.Invoke(app.NewGRPCServer),
		fx.Invoke(bootstrap),
	)
}

func bootstrap(lc fx.Lifecycle, cfg *config.Config, logger logger_shared.Logger, http *app.HTTPServer, grpc *app.GRPCServer, tracer tracer.Tracer, sentry sentry.Sentry) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			http.Configure()
			go http.Start()
			go grpc.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			defer sentry.Flush()
			tracer.Shutdown(ctx)
			http.Shutdown(ctx)
			grpc.Shutdown()
			return nil
		},
	})
}
