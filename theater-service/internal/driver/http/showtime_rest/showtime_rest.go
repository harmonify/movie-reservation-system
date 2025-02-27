package showtime_rest

import (
	"time"

	"github.com/gin-gonic/gin"
	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/ratelimiter"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/service"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/driven/config"
	http_driver_shared "github.com/harmonify/movie-reservation-system/theater-service/internal/driver/http/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
)

type ShowtimeRestHandlerParam struct {
	fx.In

	Config          *config.TheaterServiceConfig
	Logger          logger.Logger
	Tracer          tracer.Tracer
	Util            *util.Util
	Middleware      *http_driver_shared.HttpMiddleware
	Validator       http_pkg.HttpValidator
	ResponseBuilder http_pkg.HttpResponseBuilder
	ShowtimeService service.ShowtimeService
}

type ShowtimeRestHandlerResult struct {
	fx.Out

	ShowtimeRestHandler http_pkg.RestHandler `group:"http_routes"`
}

type showtimeRestHandlerImpl struct {
	config          *config.TheaterServiceConfig
	logger          logger.Logger
	tracer          tracer.Tracer
	util            *util.Util
	middleware      *http_driver_shared.HttpMiddleware
	validator       http_pkg.HttpValidator
	responseBuilder http_pkg.HttpResponseBuilder
	showtimeService service.ShowtimeService
}

func NewShowtimeRestHandler(p ShowtimeRestHandlerParam) ShowtimeRestHandlerResult {
	return ShowtimeRestHandlerResult{
		ShowtimeRestHandler: &showtimeRestHandlerImpl{
			config:          p.Config,
			logger:          p.Logger,
			tracer:          p.Tracer,
			util:            p.Util,
			middleware:      p.Middleware,
			validator:       p.Validator,
			responseBuilder: p.ResponseBuilder,
			showtimeService: p.ShowtimeService,
		},
	}
}

func (h *showtimeRestHandlerImpl) Register(g *gin.RouterGroup) error {
	var getTheaterShowtimeCap int64 = 5
	if h.config.Env == config_pkg.EnvironmentDevelopment || h.config.Env == config_pkg.EnvironmentTest {
		getTheaterShowtimeCap = 100
	}

	sg := g.Group("/showtimes")

	sg.GET(
		":showtimeId",
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getTheaterShowtimeCap,
			RefillRate: time.Second * 3,
		}),
		h.getShowtimeDetail,
	)

	return nil
}

func (h *showtimeRestHandlerImpl) Version() string {
	return "1"
}

func (h *showtimeRestHandlerImpl) getShowtimeDetail(c *gin.Context) {
	var (
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	showtimeId := c.Param("showtimeId")
	if showtimeId == "" {
		response.WithError(error_pkg.InvalidRequestPathError).Send(c)
		return
	}

	span.SetAttributes(attribute.String("showtime_id", showtimeId))

	data, err := h.showtimeService.GetShowtimeDetail(ctx, showtimeId)

	if err == nil {
		response = response.WithResult(data)
	} else {
		response = response.WithError(err)
	}

	response.Send(c)
}
