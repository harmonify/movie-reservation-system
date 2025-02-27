package grpc_repository

import (
	"context"
	"errors"
	"time"

	movie_proto "github.com/harmonify/movie-reservation-system/pkg/proto/movie"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/shared"
	redis_repository "github.com/harmonify/movie-reservation-system/theater-service/internal/driven/cache/redis/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type MovieGrpcRepositoryParam struct {
	fx.In
	*redis_repository.MovieRedisRepository
	movie_proto.MovieServiceClient
}

type movieGrpcRepositoryImpl struct {
	movieRedisRepository   *redis_repository.MovieRedisRepository
	movieServiceGrpcClient movie_proto.MovieServiceClient
}

func NewMovieGrpcRepository(p MovieGrpcRepositoryParam) shared.MovieCache {
	return &movieGrpcRepositoryImpl{
		movieRedisRepository:   p.MovieRedisRepository,
		movieServiceGrpcClient: p.MovieServiceClient,
	}
}

func (r *movieGrpcRepositoryImpl) Set(ctx context.Context, movie *movie_proto.Movie, ttl time.Duration) error {
	return r.movieRedisRepository.Set(ctx, movie, ttl)
}

func (r *movieGrpcRepositoryImpl) Get(ctx context.Context, movieID string) (*movie_proto.Movie, error) {
	movie, err := r.movieRedisRepository.Get(ctx, movieID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			movieRes, err := r.movieServiceGrpcClient.GetMovieByID(ctx, &movie_proto.GetMovieByIDRequest{MovieId: movieID})
			if err != nil {
				return nil, err
			}
			movie = movieRes.Movie
			go r.Set(ctx, movie, time.Hour)
		}
	}
	return movie, nil
}

func (r *movieGrpcRepositoryImpl) Delete(ctx context.Context, movieID string) error {
	return r.movieRedisRepository.Delete(ctx, movieID)
}
