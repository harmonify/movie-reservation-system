package movie_service

import (
	"context"

	"github.com/harmonify/movie-reservation-system/movie-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	theater_proto "github.com/harmonify/movie-reservation-system/pkg/proto/theater"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AdminMovieService interface {
	SearchMovies(ctx context.Context, p *SearchMovieParam) (*SearchMovieResult, error)
	GetMovieByID(ctx context.Context, p *GetMovieByIDParam) (*GetMovieByIDResult, error)
	SaveMovie(ctx context.Context, saveModel *entity.SaveMovie) (string, error)
	UpdateMovie(ctx context.Context, movieId string, updateModel *entity.UpdateMovie) error
	SoftDeleteMovie(ctx context.Context, movieId string) error
}

type AdminMovieServiceParam struct {
	fx.In
	logger.Logger
	tracer.Tracer
	shared.MovieStorage
	theater_proto.TheaterServiceClient
}

type adminMovieServiceImpl struct {
	logger         logger.Logger
	tracer         tracer.Tracer
	movieStorage   shared.MovieStorage
	theaterService theater_proto.TheaterServiceClient
}

func NewAdminMovieService(p AdminMovieServiceParam) AdminMovieService {
	return &adminMovieServiceImpl{
		logger:         p.Logger,
		tracer:         p.Tracer,
		movieStorage:   p.MovieStorage,
		theaterService: p.TheaterServiceClient,
	}
}

func (s *adminMovieServiceImpl) SearchMovies(ctx context.Context, p *SearchMovieParam) (*SearchMovieResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	activeMovies, err := s.theaterService.GetActiveMovies(ctx, &theater_proto.GetActiveMoviesRequest{
		TheaterId:       p.TheaterID,
		IncludeUpcoming: p.IncludeUpcoming,
	})
	if err != nil {
		return nil, err
	}

	s.logger.WithCtx(ctx).Debug("active movies", zap.Any("movies", activeMovies.Movies))

	activeMovieIds := make([]string, 0, len(activeMovies.Movies))
	for _, movie := range activeMovies.Movies {
		activeMovieIds = append(activeMovieIds, movie.MovieId)
	}

	movieResults, err := s.movieStorage.SearchMovies(ctx, &entity.SearchMovie{
		MovieIDs:        activeMovieIds,
		Genre:           p.Genre,
		Keyword:         p.Keyword,
		ReleaseDateFrom: p.ReleaseDateFrom,
		ReleaseDateTo:   p.ReleaseDateTo,
		SortBy:          p.SortBy,
		Page:            p.Page,
		PageSize:        p.PageSize,
	})

	if err != nil {
		s.logger.WithCtx(ctx).Error("failed to search movies", zap.Error(err))
		return nil, err
	}

	return &SearchMovieResult{
		Data: movieResults.Data,
		Meta: &SearchMovieMeta{
			TotalCount:  movieResults.Meta.TotalCount,
			HasNextPage: movieResults.Meta.HasNextPage,
		},
	}, nil
}

func (s *adminMovieServiceImpl) GetMovieByID(ctx context.Context, p *GetMovieByIDParam) (*GetMovieByIDResult, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	movie, err := s.movieStorage.GetMovieByID(ctx, p.MovieID)
	if err != nil {
		return nil, err
	}

	activeShowtimes, err := s.theaterService.GetActiveShowtimes(ctx, &theater_proto.GetActiveShowtimesRequest{
		TheaterId: p.TheaterID,
		MovieId:   movie.MovieID,
	})
	if err != nil {
		return nil, err
	}

	return &GetMovieByIDResult{
		Movie:     movie,
		Showtimes: activeShowtimes.GetShowtimes(),
	}, nil
}

func (s *adminMovieServiceImpl) SaveMovie(ctx context.Context, saveModel *entity.SaveMovie) (string, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	return s.movieStorage.SaveMovie(ctx, saveModel)
}

func (s *adminMovieServiceImpl) UpdateMovie(ctx context.Context, movieId string, updateModel *entity.UpdateMovie) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	return s.movieStorage.UpdateMovie(ctx, movieId, updateModel)
}

func (s *adminMovieServiceImpl) SoftDeleteMovie(ctx context.Context, movieId string) error {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	return s.movieStorage.SoftDeleteMovie(ctx, movieId)
}
