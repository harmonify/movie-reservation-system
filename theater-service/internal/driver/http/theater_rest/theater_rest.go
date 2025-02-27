package theater_rest

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/ratelimiter"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/service"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/driven/config"
	http_driver_shared "github.com/harmonify/movie-reservation-system/theater-service/internal/driver/http/shared"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
)

type TheaterRestHandlerParam struct {
	fx.In

	Config          *config.TheaterServiceConfig
	Logger          logger.Logger
	Tracer          tracer.Tracer
	Util            *util.Util
	Middleware      *http_driver_shared.HttpMiddleware
	Validator       http_pkg.HttpValidator
	ResponseBuilder http_pkg.HttpResponseBuilder
	TheaterService  service.TheaterService
}

type TheaterRestHandlerResult struct {
	fx.Out

	TheaterRestHandler http_pkg.RestHandler `group:"http_routes"`
}

type theaterRestHandlerImpl struct {
	config          *config.TheaterServiceConfig
	logger          logger.Logger
	tracer          tracer.Tracer
	util            *util.Util
	middleware      *http_driver_shared.HttpMiddleware
	validator       http_pkg.HttpValidator
	responseBuilder http_pkg.HttpResponseBuilder
	theaterService  service.TheaterService
}

func NewTheaterRestHandler(p TheaterRestHandlerParam) TheaterRestHandlerResult {
	return TheaterRestHandlerResult{
		TheaterRestHandler: &theaterRestHandlerImpl{
			config:          p.Config,
			logger:          p.Logger,
			tracer:          p.Tracer,
			util:            p.Util,
			middleware:      p.Middleware,
			validator:       p.Validator,
			responseBuilder: p.ResponseBuilder,
			theaterService:  p.TheaterService,
		},
	}
}

func (h *theaterRestHandlerImpl) Register(g *gin.RouterGroup) error {
	var getTheaterCap int64 = 10
	if h.config.Env == config_pkg.EnvironmentDevelopment || h.config.Env == config_pkg.EnvironmentTest {
		getTheaterCap = 100
	}

	amg := g.Group("/theaters")

	amg.GET(
		":theaterId",
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getTheaterCap,
			RefillRate: 3 * time.Second,
		}),
		h.searchTheaters,
	)

	return nil
}

func (h *theaterRestHandlerImpl) Version() string {
	return "1"
}

func (h *theaterRestHandlerImpl) searchTheaters(c *gin.Context) {
	var (
		err   error
		query SearchTheaterRequestQuery
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	span.SetAttributes(
		attribute.String("query.keyword", query.Keyword),
		attribute.Float64("query.latitude", float64(query.Latitude)),
		attribute.Float64("query.longitude", float64(query.Longitude)),
		attribute.Float64("query.radius", float64(query.Radius)),
		attribute.String("query.sort_by", query.SortBy),
		attribute.Int("query.page", int(query.Page)),
		attribute.Int("query.page_size", int(query.PageSize)),
	)

	response := h.responseBuilder.New().WithCtx(ctx)

	if err := h.validator.ValidateRequestQuery(c, &query); err != nil {
		response.WithCtx(ctx).WithError(err).Send(c)
		return
	}

	sortBy := entity.TheaterSortBy(strings.ToUpper(query.SortBy))
	if !sortBy.IsValid() {
		sortBy = entity.TheaterSortByNearest
	}

	data, err := h.theaterService.SearchTheaters(ctx, &entity.FindManyTheaters{
		Keyword: sql.NullString{String: query.Keyword, Valid: query.Keyword != ""},
		Location: &entity.FindManyTheatersLocation{
			Latitude:  query.Latitude,
			Longitude: query.Longitude,
			Radius:    query.Radius,
		},
		SortBy:   sortBy,
		Page:     query.Page,
		PageSize: query.PageSize,
	})

	if err == nil {
		response = response.WithResult(data.Theaters).WithMetadataFromStruct(data.Metadata)
	} else {
		response = response.WithError(err)
	}

	response.Send(c)
}
