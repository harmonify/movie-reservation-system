package movie_service

import (
	"context"

	"github.com/harmonify/movie-reservation-system/movie-service/internal/core/entity"
	"github.com/harmonify/movie-reservation-system/movie-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	theater_proto "github.com/harmonify/movie-reservation-system/pkg/proto/theater"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/fx"
)

type MovieService interface {
	GetMovieByID(ctx context.Context, movieId string) (*entity.Movie, error)
}

type MovieServiceParam struct {
	fx.In
	logger.Logger
	tracer.Tracer
	shared.MovieStorage
	theater_proto.TheaterServiceClient
}

type MovieServiceImpl struct {
	logger       logger.Logger
	tracer       tracer.Tracer
	movieStorage shared.MovieStorage
}

func NewMovieService(p MovieServiceParam) MovieService {
	return &MovieServiceImpl{
		logger:       p.Logger,
		tracer:       p.Tracer,
		movieStorage: p.MovieStorage,
	}
}

func (s *MovieServiceImpl) GetMovieByID(ctx context.Context, movieId string) (*entity.Movie, error) {
	ctx, span := s.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	movie, err := s.movieStorage.GetMovieByID(ctx, movieId)
	if err != nil {
		return nil, err
	}

	return movie, nil
}
