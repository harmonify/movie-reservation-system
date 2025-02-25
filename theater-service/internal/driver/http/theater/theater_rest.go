package theater_rest

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
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

type AdminTheaterRestHandlerParam struct {
	fx.In

	Config              *config.TheaterServiceConfig
	Logger              logger.Logger
	Tracer              tracer.Tracer
	Util                *util.Util
	Middleware          *http_driver_shared.HttpMiddleware
	Validator           http_pkg.HttpValidator
	ResponseBuilder     http_pkg.HttpResponseBuilder
	AdminTheaterService service.AdminTheaterService
}

type AdminTheaterRestHandlerResult struct {
	fx.Out

	AdminTheaterRestHandler http_pkg.RestHandler `group:"http_routes"`
}

type adminTheaterRestHandlerImpl struct {
	config              *config.TheaterServiceConfig
	logger              logger.Logger
	tracer              tracer.Tracer
	util                *util.Util
	middleware          *http_driver_shared.HttpMiddleware
	validator           http_pkg.HttpValidator
	responseBuilder     http_pkg.HttpResponseBuilder
	adminTheaterService service.AdminTheaterService
}

func NewAdminTheaterRestHandler(p AdminTheaterRestHandlerParam) AdminTheaterRestHandlerResult {
	return AdminTheaterRestHandlerResult{
		AdminTheaterRestHandler: &adminTheaterRestHandlerImpl{
			config:              p.Config,
			logger:              p.Logger,
			tracer:              p.Tracer,
			util:                p.Util,
			middleware:          p.Middleware,
			validator:           p.Validator,
			responseBuilder:     p.ResponseBuilder,
			adminTheaterService: p.AdminTheaterService,
		},
	}
}

func (h *adminTheaterRestHandlerImpl) Register(g *gin.RouterGroup) error {
	var getTheaterCap int64 = 10
	var modifyTheaterCap int64 = 2
	if h.config.Env == config_pkg.EnvironmentDevelopment || h.config.Env == config_pkg.EnvironmentTest {
		getTheaterCap = 100
		modifyTheaterCap = 100
	}

	amg := g.Group("/admin/theaters")

	amg.GET(
		"",
		h.middleware.Trace.ExtractTraceContext,
		h.middleware.AuthV2.WithPolicy("policies.theater.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getTheaterCap,
			RefillRate: 3 * time.Second,
		}),
		h.searchTheaters,
	)
	amg.POST(
		"",
		h.middleware.AuthV2.WithPolicy("policies.theater.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   modifyTheaterCap,
			RefillRate: time.Second * 3,
		}),
		h.postTheater,
	)
	amg.GET(
		":theaterId",
		h.middleware.AuthV2.WithPolicy("policies.theater.manage.allow"),
		h.getTheaterByID,
	)
	amg.GET(
		":theaterId/seats",
		h.middleware.AuthV2.WithPolicy("policies.theater.manage.allow"),
		h.getTheaterByID,
	)
	amg.PUT(
		":theaterId",
		h.middleware.AuthV2.WithPolicy("policies.theater.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   modifyTheaterCap,
			RefillRate: time.Second * 3,
		}),
		h.putTheater,
	)
	amg.DELETE(
		":theaterId",
		h.middleware.AuthV2.WithPolicy("policies.theater.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   modifyTheaterCap,
			RefillRate: time.Second * 3,
		}),
		h.deleteTheater,
	)

	return nil
}

func (h *adminTheaterRestHandlerImpl) Version() string {
	return "1"
}

func (h *adminTheaterRestHandlerImpl) searchTheaters(c *gin.Context) {
	var (
		err   error
		query AdminSearchTheaterRequestQuery
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

	data, err := h.adminTheaterService.SearchTheaters(ctx, &entity.FindManyTheaters{
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

func (h *adminTheaterRestHandlerImpl) getTheaterByID(c *gin.Context) {
	var (
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	theaterId := c.Param("theaterId")
	if theaterId == "" {
		response.WithError(error_pkg.InvalidRequestPathError).Send(c)
		return
	}

	span.SetAttributes(attribute.String("theater_id", theaterId))

	data, err := h.adminTheaterService.GetTheaterByID(ctx, &entity.FindOneTheater{
		TheaterID: sql.NullString{String: theaterId, Valid: true},
	})

	if err == nil {
		response = response.WithResult(data)
	} else {
		response = response.WithError(err)
	}

	response.Send(c)
}

func (h *adminTheaterRestHandlerImpl) postTheater(c *gin.Context) {
	var (
		body entity.SaveTheater
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		response.WithError(err).Send(c)
		return
	}

	res, err := h.adminTheaterService.SaveTheater(ctx, &body)
	if err != nil {
		response.WithError(err).Send(c)
	} else {
		response.WithResult(&AdminPostTheaterResponse{
			TheaterID: res.TheaterID,
		}).Send(c)
	}

}

func (h *adminTheaterRestHandlerImpl) putTheater(c *gin.Context) {
	var (
		body entity.UpdateTheater
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		response.WithError(err).Send(c)
		return
	}

	theaterId := c.Param("theaterId")
	if theaterId == "" {
		response.WithError(error_pkg.InvalidRequestPathError).Send(c)
		return
	}

	span.SetAttributes(attribute.String("theater_id", theaterId))

	err = h.adminTheaterService.UpdateTheater(
		ctx,
		&entity.FindOneTheater{
			TheaterID: sql.NullString{String: theaterId, Valid: true},
		},
		&body,
	)

	response.WithError(err).Send(c)
}

func (h *adminTheaterRestHandlerImpl) deleteTheater(c *gin.Context) {
	var (
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	theaterId := c.Param("theaterId")
	if theaterId == "" {
		response.WithError(error_pkg.InvalidRequestPathError).Send(c)
		return
	}

	span.SetAttributes(attribute.String("theater_id", theaterId))

	err = h.adminTheaterService.SoftDeleteTheater(ctx, &entity.FindOneTheater{
		TheaterID: sql.NullString{String: theaterId, Valid: true},
	})

	response.WithError(err).Send(c)
}
