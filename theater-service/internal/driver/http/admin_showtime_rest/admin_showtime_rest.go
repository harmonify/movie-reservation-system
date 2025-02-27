package admin_showtime_rest

import (
	"database/sql"
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

type adminShowtimeRestHandlerParam struct {
	fx.In

	Config               *config.TheaterServiceConfig
	Logger               logger.Logger
	Tracer               tracer.Tracer
	Util                 *util.Util
	Middleware           *http_driver_shared.HttpMiddleware
	Validator            http_pkg.HttpValidator
	ResponseBuilder      http_pkg.HttpResponseBuilder
	AdminShowtimeService service.AdminShowtimeService
}

type AdminShowtimeRestHandlerResult struct {
	fx.Out

	AdminShowtimeRestHandler http_pkg.RestHandler `group:"http_routes"`
}

type adminShowtimeRestHandlerImpl struct {
	config               *config.TheaterServiceConfig
	logger               logger.Logger
	tracer               tracer.Tracer
	util                 *util.Util
	middleware           *http_driver_shared.HttpMiddleware
	validator            http_pkg.HttpValidator
	responseBuilder      http_pkg.HttpResponseBuilder
	adminShowtimeService service.AdminShowtimeService
}

func NewAdminShowtimeRestHandler(p adminShowtimeRestHandlerParam) AdminShowtimeRestHandlerResult {
	return AdminShowtimeRestHandlerResult{
		AdminShowtimeRestHandler: &adminShowtimeRestHandlerImpl{
			config:               p.Config,
			logger:               p.Logger,
			tracer:               p.Tracer,
			util:                 p.Util,
			middleware:           p.Middleware,
			validator:            p.Validator,
			responseBuilder:      p.ResponseBuilder,
			adminShowtimeService: p.AdminShowtimeService,
		},
	}
}

func (h *adminShowtimeRestHandlerImpl) Register(g *gin.RouterGroup) error {
	var getTheaterShowtimeCap int64 = 10
	var modifyTheaterShowtimeCap int64 = 2
	if h.config.Env == config_pkg.EnvironmentDevelopment || h.config.Env == config_pkg.EnvironmentTest {
		getTheaterShowtimeCap = 100
		modifyTheaterShowtimeCap = 100
	}

	sg := g.Group("/admin/showtimes")

	sg.GET(
		"",
		h.middleware.AuthV2.WithPolicy("policies.theater.showtime.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getTheaterShowtimeCap,
			RefillRate: time.Second * 3,
		}),
		h.getShowtimes,
	)
	sg.POST(
		"",
		h.middleware.AuthV2.WithPolicy("policies.theater.showtime.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   modifyTheaterShowtimeCap,
			RefillRate: time.Second * 3,
		}),
		h.createShowtime,
	)
	sg.GET(
		":showtimeId",
		h.middleware.AuthV2.WithPolicy("policies.theater.showtime.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getTheaterShowtimeCap,
			RefillRate: time.Second * 3,
		}),
		h.getShowtime,
	)
	sg.PUT(
		":showtimeId",
		h.middleware.AuthV2.WithPolicy("policies.theater.showtime.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   modifyTheaterShowtimeCap,
			RefillRate: time.Second * 3,
		}),
		h.updateShowtime,
	)
	sg.DELETE(
		":showtimeId",
		h.middleware.AuthV2.WithPolicy("policies.theater.showtime.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   modifyTheaterShowtimeCap,
			RefillRate: time.Second * 3,
		}),
		h.deleteShowtime,
	)

	return nil
}

func (h *adminShowtimeRestHandlerImpl) Version() string {
	return "1"
}

func (h *adminShowtimeRestHandlerImpl) getShowtimes(c *gin.Context) {
	var (
		err   error
		query AdminSearchShowtimeRequestQuery
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	span.SetAttributes(
		attribute.String("query.theater_id", query.TheaterID),
		attribute.String("query.room_id", query.RoomID),
		attribute.String("query.movie_id", query.MovieID),
		attribute.Int64("query.start_time_gte_unix", query.StartTimeGteUnix),
		attribute.Int64("query.start_time_lte_unix", query.StartTimeLteUnix),
		attribute.String("query.sort_by", query.SortBy),
		attribute.Int("query.page", int(query.Page)),
		attribute.Int("query.page_size", int(query.PageSize)),
	)

	response := h.responseBuilder.New().WithCtx(ctx)

	if err := h.validator.ValidateRequestQuery(c, &query); err != nil {
		response.WithCtx(ctx).WithError(err).Send(c)
		return
	}

	startTimeGte := time.Unix(query.StartTimeGteUnix, 0)
	startTimeLte := time.Unix(query.StartTimeLteUnix, 0)

	data, err := h.adminShowtimeService.SearchShowtimes(ctx, &entity.FindManyShowtimes{
		TheaterID:    sql.NullString{String: query.TheaterID, Valid: query.TheaterID != ""},
		RoomID:       sql.NullString{String: query.RoomID, Valid: query.RoomID != ""},
		MovieID:      sql.NullString{String: query.MovieID, Valid: query.MovieID != ""},
		StartTimeGte: sql.NullTime{Time: startTimeGte, Valid: !startTimeGte.IsZero()},
		StartTimeLte: sql.NullTime{Time: startTimeLte, Valid: !startTimeLte.IsZero()},
		SortBy:       entity.ShowtimeSortBy(query.SortBy),
		Page:         query.Page,
		PageSize:     query.PageSize,
	})

	if err == nil {
		response = response.WithResult(data.Showtimes).WithMetadataFromStruct(data.Metadata)
	} else {
		response = response.WithError(err)
	}

	response.Send(c)
}

func (h *adminShowtimeRestHandlerImpl) getShowtime(c *gin.Context) {
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

	data, err := h.adminShowtimeService.GetShowtimeByID(ctx, &entity.FindOneShowtime{
		ShowtimeID: sql.NullString{String: showtimeId, Valid: true},
	})

	if err == nil {
		response = response.WithResult(data)
	} else {
		response = response.WithError(err)
	}

	response.Send(c)
}

func (h *adminShowtimeRestHandlerImpl) createShowtime(c *gin.Context) {
	var (
		body entity.SaveShowtime
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		response.WithError(err).Send(c)
		return
	}

	res, err := h.adminShowtimeService.SaveShowtime(ctx, &body)
	if err != nil {
		response.WithError(err).Send(c)
	} else {
		response.WithResult(&AdminPostShowtimeResponse{
			ShowtimeID: res.ShowtimeID,
		}).Send(c)
	}
}

func (h *adminShowtimeRestHandlerImpl) updateShowtime(c *gin.Context) {
	var (
		body entity.UpdateShowtime
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		response.WithError(err).Send(c)
		return
	}

	showtimeId := c.Param("showtimeId")
	if showtimeId == "" {
		response.WithError(error_pkg.InvalidRequestPathError).Send(c)
		return
	}

	span.SetAttributes(attribute.String("showtime_id", showtimeId))

	err = h.adminShowtimeService.UpdateShowtime(
		ctx,
		&entity.FindOneShowtime{
			ShowtimeID: sql.NullString{String: showtimeId, Valid: true},
		},
		&body,
	)

	response.WithError(err).Send(c)
}

func (h *adminShowtimeRestHandlerImpl) deleteShowtime(c *gin.Context) {
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

	err = h.adminShowtimeService.SoftDeleteShowtime(ctx, &entity.FindOneShowtime{
		ShowtimeID: sql.NullString{String: showtimeId, Valid: true},
	})

	response.WithError(err).Send(c)
}
