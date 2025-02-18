package movie_service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"

	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/core/shared"
	theater_proto "github.com/harmonify/movie-reservation-system/movie-search-service/internal/driven/proto/theater"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MovieSearchService interface {
	SearchMovies(ctx context.Context, p *SearchMovieParam) (*SearchMovieResult, error)
	SearchMoviesWithCursor(ctx context.Context, cursor string) (*SearchMovieResult, error)
}

type MovieSearchServiceParam struct {
	fx.In
	logger.Logger
	tracer.Tracer
	shared.MovieStorage
	theater_proto.TheaterServiceClient
	validation.StructValidator
}

type movieSearchServiceImpl struct {
	logger          logger.Logger
	tracer          tracer.Tracer
	movieStorage    shared.MovieStorage
	theaterService  theater_proto.TheaterServiceClient
	structValidator validation.StructValidator
}

func NewMovieSearchService(p MovieSearchServiceParam) MovieSearchService {
	return &movieSearchServiceImpl{
		logger:          p.Logger,
		tracer:          p.Tracer,
		movieStorage:    p.MovieStorage,
		theaterService:  p.TheaterServiceClient,
		structValidator: p.StructValidator,
	}
}

func (s *movieSearchServiceImpl) SearchMovies(ctx context.Context, p *SearchMovieParam) (*SearchMovieResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	activeMovies, err := s.theaterService.GetActiveMovies(ctx, &theater_proto.GetActiveMoviesRequest{
		TheaterId:       p.TheaterID,
		IncludeUpcoming: p.IncludeUpcoming,
	})
	if err != nil {
		s.logger.WithCtx(ctx).Error("failed to get active movies", zap.Error(err), zap.String("theater_id", p.TheaterID), zap.Bool("include_upcoming", p.IncludeUpcoming))
		return nil, err
	}

	s.logger.WithCtx(ctx).Debug("active movies", zap.Any("movies", activeMovies.Movies))

	activeMovieIds := make([]string, 0, len(activeMovies.Movies))
	for _, movie := range activeMovies.Movies {
		activeMovieIds = append(activeMovieIds, movie.MovieId)
	}

	movieResults, err := s.movieStorage.SearchMovies(ctx,
		&entity.SearchMovieParam{
			MovieIDs:          activeMovieIds,
			Keyword:           sql.NullString{String: p.Keyword, Valid: p.Keyword != ""},
			Genre:             sql.NullString{String: p.Genre, Valid: p.Genre != ""},
			SortBy:            p.SortBy,
			Limit:             p.Limit,
			LastSeenSortValue: p.LastSeenSortValue,
			LastSeenID:        p.LastSeenID,
		},
	)
	if err != nil {
		s.logger.WithCtx(ctx).Error("failed to search movies", zap.Error(err))
		return nil, err
	}

	var nextCursor string
	if movieResults.Meta.HasNextPage {
		nextCursor, err = s.generateCursor(&SearchMovieCursor{
			TheaterID:         p.TheaterID,
			IncludeUpcoming:   p.IncludeUpcoming,
			Genre:             p.Genre,
			Keyword:           p.Keyword,
			SortBy:            p.SortBy,
			Limit:             p.Limit,
			LastSeenSortValue: movieResults.Meta.LastSeenSortValue,
			LastSeenID:        movieResults.Meta.LastSeenID,
		})
		if err != nil {
			s.logger.WithCtx(ctx).Error("failed to generate cursor", zap.Error(err))
			return nil, err
		}
	}

	return &SearchMovieResult{
		Data: movieResults.Data,
		Meta: &SearchMovieMeta{
			TotalCount:  movieResults.Meta.TotalCount,
			HasNextPage: movieResults.Meta.HasNextPage,
			NextCursor:  nextCursor,
		},
	}, nil
}

func (s *movieSearchServiceImpl) SearchMoviesWithCursor(ctx context.Context, cursor string) (*SearchMovieResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cursorData, err := s.parseCursor(cursor)
	if err != nil {
		s.logger.WithCtx(ctx).Error("failed to decode cursor", zap.Error(err), zap.Any("cursor_data", cursorData), zap.String("cursor", cursor))
		return nil, err
	}

	s.logger.WithCtx(ctx).Debug("searching movies with cursor", zap.Any("pagination", cursorData))

	return s.SearchMovies(ctx, &SearchMovieParam{
		TheaterID:         cursorData.TheaterID,
		IncludeUpcoming:   cursorData.IncludeUpcoming,
		Genre:             cursorData.Genre,
		Keyword:           cursorData.Keyword,
		SortBy:            cursorData.SortBy,
		Limit:             cursorData.Limit,
		LastSeenSortValue: cursorData.LastSeenSortValue,
		LastSeenID:        cursorData.LastSeenID,
	})
}

func (s *movieSearchServiceImpl) parseCursor(cursor string) (*SearchMovieCursor, error) {
	jsonPagination, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var p SearchMovieCursor
	if err := json.Unmarshal(jsonPagination, &p); err != nil {
		return nil, err
	}

	if _, validationErrs := s.structValidator.Validate(p); len(validationErrs) > 0 {
		return &p, InvalidCursorError.WithErrors(validationErrs...)
	}

	return &p, nil
}

func (s *movieSearchServiceImpl) generateCursor(c *SearchMovieCursor) (string, error) {
	jsonPagination, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(jsonPagination), nil
}
