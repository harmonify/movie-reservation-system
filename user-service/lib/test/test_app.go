package test

import (
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"go.uber.org/fx"
)

// This is a function to initialize all services and invoke their functions.
// It accepts a generic type of invoker, so you can specify some modules to be mocked in the test file
func NewTestApp(invoker interface{}, overrideConstructors ...any) *fx.App {
	options := []fx.Option{
		fx.Provide(
			func() *config.ConfigFile {
				return &config.ConfigFile{
					Path: ".env",
				}
			},
		),
		config.ConfigModule,
		logger.LoggerModule,
		tracer.TracerModule,

		http.HttpModule,
		util.UtilModule,

		// Invoke the function
		fx.Invoke(logger.NewLogger),
		fx.Invoke(tracer.InitTracer),
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
