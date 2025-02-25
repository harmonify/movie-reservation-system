package shared

import (
	"context"

	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/core/entity"
)

type (
	MovieStorage interface {
		SearchMovies(ctx context.Context, searchModel *entity.SearchMovieParam) (*SearchMovieResult, error)
	}

	SearchMovieResult struct {
		Data []*entity.SearchMovieResult
		Meta *SearchMovieMetadata
	}

	SearchMovieMetadata struct {
		TotalCount        int64
		HasNextPage       bool
		LastSeenSortValue interface{}
		LastSeenID        string
	}
)
