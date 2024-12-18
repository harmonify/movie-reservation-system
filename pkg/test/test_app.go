package test

import (
	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
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
