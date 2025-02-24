package movie_rest

import (
	"time"

	"github.com/harmonify/movie-reservation-system/movie-service/internal/core/entity"
)

type (
	AdminGetMovieRequestQuery struct {
		TheaterID       string    `json:"theater_id" form:"theater_id" validate:"required,alphanumunicode"`
		IncludeUpcoming bool      `json:"include_upcoming" form:"include_upcoming" validate:"boolean"`
		Genre           string    `json:"genre"`
		Keyword         string    `json:"keyword"`
		ReleaseDateFrom time.Time `json:"release_date_from"`
		ReleaseDateTo   time.Time `json:"release_date_to"`
		SortBy          string    `json:"sort_by" form:"sort_by" validate:"required,oneof=relevance release_date_desc title_asc title_desc"`
		Page            int64     `json:"page" form:"page" validate:"required,numeric"`
		PageSize        int64     `json:"page_size" form:"page_size" validate:"required,numeric"`
	}

	AdminGetMovieResponse []*entity.SearchMovieResult

	AdminPostMovieResponse struct {
		MovieID string `json:"movie_id"`
	}

	// AdminPutMovieResponse entity.Movie

	// AdminDeleteMovieResponse entity.Movie
)
