package movie_rest

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/core/entity"
	movie_service "github.com/harmonify/movie-reservation-system/movie-service/internal/core/service/movie"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/driven/config"
	http_driver_shared "github.com/harmonify/movie-reservation-system/movie-service/internal/driver/http/shared"
	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/ratelimiter"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AdminMovieRestHandlerParam struct {
	fx.In

	Config            *config.MovieServiceConfig
	Logger            logger.Logger
	Tracer            tracer.Tracer
	Util              *util.Util
	Middleware        *http_driver_shared.HttpMiddleware
	Validator         http_pkg.HttpValidator
	ResponseBuilder   http_pkg.HttpResponseBuilder
	AdminMovieService movie_service.AdminMovieService
}

type AdminMovieRestHandlerResult struct {
	fx.Out

	AdminMovieRestHandler http_pkg.RestHandler `group:"http_routes"`
}

type adminMovieRestHandlerImpl struct {
	config            *config.MovieServiceConfig
	logger            logger.Logger
	tracer            tracer.Tracer
	util              *util.Util
	middleware        *http_driver_shared.HttpMiddleware
	validator         http_pkg.HttpValidator
	responseBuilder   http_pkg.HttpResponseBuilder
	adminMovieService movie_service.AdminMovieService
}

func NewAdminMovieRestHandler(p AdminMovieRestHandlerParam) AdminMovieRestHandlerResult {
	return AdminMovieRestHandlerResult{
		AdminMovieRestHandler: &adminMovieRestHandlerImpl{
			config:            p.Config,
			logger:            p.Logger,
			tracer:            p.Tracer,
			util:              p.Util,
			middleware:        p.Middleware,
			validator:         p.Validator,
			responseBuilder:   p.ResponseBuilder,
			adminMovieService: p.AdminMovieService,
		},
	}
}

func (h *adminMovieRestHandlerImpl) Register(g *gin.RouterGroup) error {
	var getMovieCap int64 = 10
	var modifyMovieCap int64 = 2
	if h.config.Env == config_pkg.EnvironmentDevelopment || h.config.Env == config_pkg.EnvironmentTest {
		getMovieCap = 100
		modifyMovieCap = 100
	}

	amg := g.Group("/admin/movies")

	amg.GET(
		"",
		h.middleware.Trace.ExtractTraceContext,
		h.middleware.AuthV2.WithPolicy("policies.movie.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getMovieCap,
			RefillRate: 3 * time.Second,
		}),
		h.searchMovies,
	)
	amg.POST(
		"",
		h.middleware.AuthV2.WithPolicy("policies.movie.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   modifyMovieCap,
			RefillRate: time.Second * 3,
		}),
		h.postMovie,
	)
	amg.GET(
		":movieId",
		h.middleware.AuthV2.WithPolicy("policies.movie.manage.allow"),
		h.getMovieByID,
	)
	amg.PUT(
		":movieId",
		h.middleware.AuthV2.WithPolicy("policies.movie.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   modifyMovieCap,
			RefillRate: time.Second * 3,
		}),
		h.putMovie,
	)
	amg.DELETE(
		":movieId",
		h.middleware.AuthV2.WithPolicy("policies.movie.manage.allow"),
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   modifyMovieCap,
			RefillRate: time.Second * 3,
		}),
		h.deleteMovie,
	)

	return nil
}

func (h *adminMovieRestHandlerImpl) Version() string {
	return "1"
}

func (h *adminMovieRestHandlerImpl) searchMovies(c *gin.Context) {
	var (
		err   error
		query AdminSearchMovieRequestQuery
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	if err := h.validator.ValidateRequestQuery(c, &query); err != nil {
		response.WithCtx(ctx).WithError(err).Send(c)
		return
	}

	h.logger.WithCtx(ctx).Debug("searchMovies query", zap.Any("query", query))

	data, err := h.adminMovieService.SearchMovies(ctx, &movie_service.SearchMovieParam{
		TheaterID:       query.TheaterID,
		IncludeUpcoming: query.IncludeUpcoming,
		Genre:           sql.NullString{String: query.Genre, Valid: query.Genre != ""},
		Keyword:         sql.NullString{String: query.Keyword, Valid: query.Keyword != ""},
		ReleaseDateFrom: sql.NullTime{Time: query.ReleaseDateFrom, Valid: !query.ReleaseDateFrom.IsZero()},
		ReleaseDateTo:   sql.NullTime{Time: query.ReleaseDateTo, Valid: !query.ReleaseDateTo.IsZero()},
		SortBy:          entity.MovieSortBy(query.SortBy),
		Page:            query.Page,
		PageSize:        query.PageSize,
	})

	if err == nil {
		response = response.WithResult(data.Data).WithMetadataFromStruct(data.Meta)
	} else {
		response = response.WithError(err)
	}

	response.Send(c)
}

func (h *adminMovieRestHandlerImpl) getMovieByID(c *gin.Context) {
	var (
		err   error
		query AdminGetMovieByIDRequestQuery
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	movieId := c.Param("movieId")
	if movieId == "" {
		response.WithError(error_pkg.InvalidRequestPathError).Send(c)
		return
	}

	if err := h.validator.ValidateRequestQuery(c, &query); err != nil {
		response.WithError(err).Send(c)
		return
	}

	data, err := h.adminMovieService.GetMovieByID(ctx, &movie_service.GetMovieByIDParam{
		TheaterID: query.TheaterID,
		MovieID:   movieId,
	})

	if err == nil {
		response = response.WithResult(data)
	} else {
		response = response.WithError(err)
	}

	response.Send(c)
}

func (h *adminMovieRestHandlerImpl) postMovie(c *gin.Context) {
	var (
		body entity.SaveMovie
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		response.WithError(err).Send(c)
		return
	}

	id, err := h.adminMovieService.SaveMovie(ctx, &body)

	response.WithError(err).WithResult(&AdminPostMovieResponse{
		MovieID: id,
	}).Send(c)
}

func (h *adminMovieRestHandlerImpl) putMovie(c *gin.Context) {
	var (
		body entity.UpdateMovie
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		response.WithError(err).Send(c)
		return
	}

	movieId := c.Param("movieId")
	if movieId == "" {
		response.WithError(error_pkg.InvalidRequestPathError).Send(c)
		return
	}

	err = h.adminMovieService.UpdateMovie(ctx, movieId, &body)

	response.WithError(err).Send(c)
}

func (h *adminMovieRestHandlerImpl) deleteMovie(c *gin.Context) {
	var (
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	response := h.responseBuilder.New().WithCtx(ctx)

	movieId := c.Param("movieId")
	if movieId == "" {
		response.WithError(error_pkg.InvalidRequestPathError).Send(c)
		return
	}

	err = h.adminMovieService.SoftDeleteMovie(ctx, movieId)

	response.WithError(err).Send(c)
}
