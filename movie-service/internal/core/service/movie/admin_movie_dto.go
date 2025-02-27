package movie_service

import (
	"database/sql"

	"github.com/harmonify/movie-reservation-system/movie-service/internal/core/entity"
	theater_proto "github.com/harmonify/movie-reservation-system/pkg/proto/theater"
)

// SearchMovieParam is a struct to search for movies
type SearchMovieParam struct {
	// filters
	TheaterID       string         `json:"theater_id" validate:"required"` // theater ID to search for
	IncludeUpcoming bool           `json:"include_upcoming"`               // include upcoming movies in the search
	Genre           sql.NullString `json:"genre"`                          // genre to search for in the movie genres
	Keyword         sql.NullString `json:"keyword"`                        // keyword to search for in the movie text fields
	ReleaseDateFrom sql.NullTime   `json:"release_date_from"`
	ReleaseDateTo   sql.NullTime   `json:"release_date_to"`
	// pagination
	SortBy   entity.MovieSortBy `json:"sort_by"`   // sorting order of the movies. If keyword is provided, this setting will be ignored and the result will be sorted by relevance
	Page     int64              `json:"page"`      // page number
	PageSize int64              `json:"page_size"` // number of movies to return per page
}

type SearchMovieResult struct {
	Data []*entity.SearchMovieResult `json:"data"`
	Meta *SearchMovieMeta            `json:"meta"`
}

type SearchMovieMeta struct {
	TotalCount  int64 `json:"total_count"`
	HasNextPage bool  `json:"has_next_page"`
}

type GetMovieByIDParam struct {
	TheaterID string `json:"theater_id" validate:"required"` // theater ID to search for
	MovieID   string `json:"movie_id" validate:"required"`   // movie ID to search for
}

type GetMovieByIDResult struct {
	Movie     *entity.Movie                                        `json:"movie"`
	Showtimes []*theater_proto.GetActiveShowtimesResponse_Showtime `json:"showtimes"`
}
