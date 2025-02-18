package movie_service

import (
	"github.com/harmonify/movie-reservation-system/movie-search-service/internal/core/entity"
)

// SearchMovieParam is a struct to search for movies
type SearchMovieParam struct {
	// filters
	TheaterID       string // theater ID to search for
	IncludeUpcoming bool   // include upcoming movies in the search
	Genre           string // genre to search for in the movie genres
	Keyword         string // keyword to search for in the movie text fields
	// pagination
	SortBy            entity.MovieSortBy
	Limit             int64
	LastSeenSortValue interface{}
	LastSeenID        string
}

type SearchMovieResult struct {
	Data []*entity.SearchMovieResult `json:"data"`
	Meta *SearchMovieMeta            `json:"meta"`
}

type SearchMovieMeta struct {
	TotalCount  int64  `json:"total_count"`
	HasNextPage bool   `json:"has_next_page"`
	NextCursor  string `json:"next_cursor,omitempty"`
}

type SearchMovieCursor struct {
	Cursor            string             `json:"cursor" form:"cursor" validate:"omitempty,alphanum"`
	SortBy            entity.MovieSortBy `json:"sort_by" form:"sort_by" validate:"required,oneof=relevance release_date_desc title_asc title_desc"`
	Limit             int64              `json:"limit" form:"limit" validate:"required,numeric"`
	TheaterID         string             `json:"theater_id" form:"theater_id" validate:"required,alphanumunicode"`
	IncludeUpcoming   bool               `json:"include_upcoming" form:"include_upcoming" validate:"boolean"`
	Genre             string             `json:"genre" form:"genre"`
	Keyword           string             `json:"keyword" form:"keyword"`
	LastSeenSortValue interface{}        `json:"last_seen_sort_value" validate:"required"`
	LastSeenID        string             `json:"last_seen_id" validate:"required"`
}
