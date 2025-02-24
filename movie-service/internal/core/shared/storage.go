package shared

import (
	"context"

	"github.com/harmonify/movie-reservation-system/movie-service/internal/core/entity"
)

type (
	MovieStorage interface {
		SearchMovies(ctx context.Context, searchModel *entity.SearchMovie) (*SearchMovieResult, error)
		SaveMovie(ctx context.Context, saveModel *entity.SaveMovie) (id string, err error)
		UpdateMovie(ctx context.Context, movieId string, updateModel *entity.UpdateMovie) error
		SoftDeleteMovie(ctx context.Context, movieId string) error
	}

	SearchMovieResult struct {
		Data []*entity.SearchMovieResult
		Meta *SearchMovieMetadata
	}

	SearchMovieMetadata struct {
		TotalCount  int64
		HasNextPage bool
	}
)
