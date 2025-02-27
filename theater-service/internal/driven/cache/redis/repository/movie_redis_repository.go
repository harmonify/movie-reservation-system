package redis_repository

import (
	"context"
	"time"

	"github.com/harmonify/movie-reservation-system/pkg/cache"
	movie_proto "github.com/harmonify/movie-reservation-system/pkg/proto/movie"
	"github.com/harmonify/movie-reservation-system/theater-service/internal/core/shared"
	"go.uber.org/fx"
)

var _ shared.MovieCache = (*MovieRedisRepository)(nil)

type MovieRedisRepositoryParam struct {
	fx.In
	Redis *cache.Redis
}

type MovieRedisRepository struct {
	redis *cache.Redis
}

func NewMovieRedisRepository(p MovieRedisRepositoryParam) *MovieRedisRepository {
	return &MovieRedisRepository{
		redis: p.Redis,
	}
}

func (r *MovieRedisRepository) constructMovieKey(movieID string) string {
	return "movie:" + movieID
}

func (r *MovieRedisRepository) Set(ctx context.Context, movie *movie_proto.Movie, ttl time.Duration) error {
	return r.redis.Client.Set(ctx, r.constructMovieKey(movie.MovieId), movie, ttl).Err()
}

func (r *MovieRedisRepository) Get(ctx context.Context, movieID string) (*movie_proto.Movie, error) {
	res := r.redis.Client.Get(ctx, r.constructMovieKey(movieID))
	if res.Err() != nil {
		return nil, res.Err()
	}
	var movie movie_proto.Movie
	err := res.Scan(&movie)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRedisRepository) Delete(ctx context.Context, movieID string) error {
	return r.redis.Client.Del(ctx, r.constructMovieKey(movieID)).Err()
}
