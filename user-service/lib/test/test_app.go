package test

import (
	"path"
	"runtime"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"go.uber.org/fx"
)

// This is a function to initialize all components of the library.
func NewTestApp(p ...fx.Option) *fx.App {
	options := []fx.Option{
		fx.Provide(
			func() *config.ConfigFile {
				_, filename, _, _ := runtime.Caller(0)
				return &config.ConfigFile{
					Path: path.Join(filename, "..", "..", "..", ".env.test"),
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
	}

	// Override dependencies
	if len(p) > 0 {
		for _, c := range p {
			options = append(options, c)
		}
	}

	return fx.New(options...)
}
