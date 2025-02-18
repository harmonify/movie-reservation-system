package movie_rest

import (
	"github.com/gin-gonic/gin"
	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/core/entity"
	movie_service "github.com/harmonify/movie-reservation-system/movie-search-service/internal/core/service/movie"
	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/driven/config"
	http_driver_shared "github.com/harmonify/movie-reservation-system/movie-search-service/internal/driver/http/shared"
	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"

	// "github.com/harmonify/movie-reservation-system/pkg/ratelimiter"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MovieRestHandlerParam struct {
	fx.In

	Config             *config.MovieSearchServiceConfig
	Logger             logger.Logger
	Tracer             tracer.Tracer
	Util               *util.Util
	Middleware         *http_driver_shared.HttpMiddleware
	Validator          http_pkg.HttpValidator
	ResponseBuilder    http_pkg.HttpResponseBuilder
	MovieSearchService movie_service.MovieSearchService
	StructValidator    validation.StructValidator
}

type MovieRestHandlerResult struct {
	fx.Out

	MovieRestHandler http_pkg.RestHandler `group:"http_routes"`
}

type movieRestHandlerImpl struct {
	config             *config.MovieSearchServiceConfig
	logger             logger.Logger
	tracer             tracer.Tracer
	util               *util.Util
	middleware         *http_driver_shared.HttpMiddleware
	validator          http_pkg.HttpValidator
	responseBuilder    http_pkg.HttpResponseBuilder
	movieSearchService movie_service.MovieSearchService
	structValidator    validation.StructValidator
}

func NewMovieRestHandler(p MovieRestHandlerParam) MovieRestHandlerResult {
	return MovieRestHandlerResult{
		MovieRestHandler: &movieRestHandlerImpl{
			config:             p.Config,
			logger:             p.Logger,
			tracer:             p.Tracer,
			util:               p.Util,
			middleware:         p.Middleware,
			validator:          p.Validator,
			responseBuilder:    p.ResponseBuilder,
			movieSearchService: p.MovieSearchService,
			structValidator:    p.StructValidator,
		},
	}
}

func (h *movieRestHandlerImpl) Register(g *gin.RouterGroup) error {
	var searchMovieCap int64 = 5
	if h.config.Env == config_pkg.EnvironmentDevelopment || h.config.Env == config_pkg.EnvironmentTest {
		searchMovieCap = 100
	}
	h.logger.Debug("searchMovieCap", zap.Int("searchMovieCap", int(searchMovieCap)))

	g.GET(
		"/movie",
		// h.middleware.Auth.AuthenticateUser,
		// h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
		// 	Capacity:   searchMovieCap,
		// 	RefillRate: 3 * time.Second,
		// }),
		h.searchMovies,
	)

	return nil
}

func (h *movieRestHandlerImpl) Version() string {
	return "1"
}

func (h *movieRestHandlerImpl) searchMovies(c *gin.Context) {
	var (
		err   error
		query CustomerGetMovieRequestQuery
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	if err := c.ShouldBindQuery(&query); err != nil {
		h.responseBuilder.New().
			WithCtx(ctx).
			WithError(error_pkg.InvalidRequestQueryError.WithErrors(h.structValidator.ConstructValidationErrorFields(err)...)).
			Send(c)
		return
	}

	h.logger.WithCtx(ctx).Debug("searchMovies query", zap.Any("query", query))

	var data *movie_service.SearchMovieResult
	if query.Cursor != "" {
		data, err = h.movieSearchService.SearchMoviesWithCursor(ctx, query.Cursor)
	} else {
		if err, validationErrs := h.structValidator.Validate(query); err != nil {
			h.logger.WithCtx(ctx).Error("invalid search movie param", zap.Error(err))
			h.responseBuilder.New().
				WithCtx(ctx).
				WithError(error_pkg.InvalidRequestQueryError.WithErrors(validationErrs...)).
				Send(c)
			return
		}

		data, err = h.movieSearchService.SearchMovies(ctx, &movie_service.SearchMovieParam{
			TheaterID:         query.TheaterID,
			IncludeUpcoming:   query.IncludeUpcoming,
			Genre:             query.Genre,
			Keyword:           query.Keyword,
			SortBy:            entity.MovieSortBy(query.SortBy),
			Limit:             query.Limit,
			LastSeenSortValue: nil,
			LastSeenID:        "",
		})
	}

	response := h.responseBuilder.New().WithCtx(ctx)

	if err == nil {
		response = response.WithResult(data.Data).WithMetadataFromStruct(data.Meta)
	} else {
		response = response.WithError(err)
	}

	response.Send(c)
}
