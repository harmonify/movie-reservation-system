package internal

import (
	"context"
	"fmt"

	"github.com/harmonify/movie-reservation-system/user-service/internal/driver/http"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"go.uber.org/fx"
)

func StartApp() error {
	app := NewApp(bootstrap)

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
// It accepts a generic type of invoker, so you can specify some modules to be mocked in the test file
func NewApp(invoker interface{}, overrideConstructors ...any) *fx.App {
	options := []fx.Option{
		config.ConfigModule,
		logger.LoggerModule,
		tracer.TracerModule,

		http.HttpModule,
		util.UtilModule,

		// CORE

		// INFRA (DRIVEN)

		// API (DRIVER)
		http.HttpModule,

		fx.Invoke(invoker),
	}

	// Override dependencies
	if len(overrideConstructors) > 0 {
		for _, c := range overrideConstructors {
			options = append(options, fx.Decorate(c))
		}
	}

	return fx.New(options...)
}

func bootstrap(lc fx.Lifecycle, l logger.Logger, h http.HttpServer, t tracer.Tracer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := h.Start(ctx, http.RestHandlers...)
			if err != nil {
				l.WithCtx(ctx).Error(err)
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := h.Shutdown(ctx)
			if err != nil {
				l.WithCtx(ctx).Error(err)
				return err
			}

			t.Shutdown(ctx)
			if err != nil {
				l.WithCtx(ctx).Error(err)
				return err
			}

			return nil
		},
	})
}
