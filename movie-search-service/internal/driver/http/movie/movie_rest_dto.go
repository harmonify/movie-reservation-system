package movie_rest

import (
	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/core/entity"
)

type (
	CustomerGetMovieRequestQuery struct {
		Cursor          string `json:"cursor" form:"cursor" validate:"omitempty,alphanum"`
		SortBy          string `json:"sort_by" form:"sort_by" validate:"required,oneof=relevance release_date_desc title_asc title_desc"`
		Limit           int64  `json:"limit" form:"limit" validate:"required,numeric"`
		TheaterID       string `json:"theater_id" form:"theater_id" validate:"required,alphanumunicode"`
		IncludeUpcoming bool   `json:"include_upcoming" form:"include_upcoming" validate:"boolean"`
		Genre           string `json:"genre" form:"genre"`
		Keyword         string `json:"keyword" form:"keyword"`
	}

	CustomerGetMovieResponse []entity.Movie
)
