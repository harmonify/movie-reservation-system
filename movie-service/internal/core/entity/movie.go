package entity

import (
	"database/sql"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Movie struct {
	// Administration data
	MovieID   string    `json:"movie_id" bson:"_id"`
	TraceID   string    `json:"trace_id" bson:"trace_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	// Movie details
	Title              string              `json:"title" bson:"title"`
	Description        string              `json:"description" bson:"description"`
	Genres             []string            `json:"genres" bson:"genres"`
	PosterImageURL     string              `json:"poster_image_url" bson:"poster_image_url"`
	PhotoURLs          []string            `json:"photo_urls" bson:"photo_urls"`
	TrailerURL         string              `json:"trailer_url" bson:"trailer_url"`
	Runtime            time.Duration       `json:"runtime" bson:"runtime"`
	ReleaseDate        time.Time           `json:"release_date" bson:"release_date"`
	ParentalGuidances  []*ParentalGuidance `json:"parental_guide" bson:"parental_guide"`
	Dub                Language            `json:"dub" bson:"dub"`
	AvailableSubtitles []Language          `json:"available_subtitles" bson:"available_subtitles"`
	// Crew details
	Cast              []*People `json:"cast" bson:"cast"`
	Director          *People   `json:"director" bson:"director"`
	Writer            *People   `json:"writer" bson:"writer"`
	ProductionCompany string    `json:"production_company" bson:"production_company"`
}

type Language struct {
	Name string `json:"name" bson:"name" validate:"required"`
	Code string `json:"code" bson:"code" validate:"required"`
}

func (l *Language) String() string {
	return l.Code
}

type ParentalGuidance struct {
	Code string `json:"code" bson:"code" validate:"required"` // Parental guidance code
	CountryCode string `json:"country_code" bson:"country_code" validate:"required"` // ISO 3166-1 alpha-2 country code
}

func (pg *ParentalGuidance) String() string {
	return pg.CountryCode + ":" + pg.Code
}

type SaveMovie struct {
	// Movie details
	Title              string              `json:"title" bson:"title" validate:"required"`
	Description        string              `json:"description" bson:"description" validate:"required"`
	Genres             []string            `json:"genres" bson:"genres" validate:"required"`
	PosterImageURL     string              `json:"poster_image_url" bson:"poster_image_url" validate:"required"`
	PhotoURLs          []string            `json:"photo_urls" bson:"photo_urls" validate:"required"`
	TrailerURL         string              `json:"trailer_url" bson:"trailer_url" validate:"required"`
	Runtime            time.Duration       `json:"runtime" bson:"runtime" validate:"required"`
	ReleaseDate        time.Time           `json:"release_date" bson:"release_date" validate:"required"`
	ParentalGuidances  []*ParentalGuidance `json:"parental_guide" bson:"parental_guide" validate:"required"`
	Dub                Language            `json:"dub" bson:"dub" validate:"required"`
	AvailableSubtitles []Language          `json:"available_subtitles" bson:"available_subtitles" validate:"required"`
	// Crew details
	Cast              []*People `json:"cast" bson:"cast" validate:"required"`
	Director          *People   `json:"director" bson:"director" validate:"required"`
	Writer            *People   `json:"writer" bson:"writer" validate:"required"`
	ProductionCompany string    `json:"production_company" bson:"production_company" validate:"required"`
}

type UpdateMovie struct {
	// Movie details
	Title              string              `json:"title" bson:"title" validate:"required"`
	Description        string              `json:"description" bson:"description" validate:"required"`
	Genres             []string            `json:"genres" bson:"genres" validate:"required"`
	PosterImageURL     string              `json:"poster_image_url" bson:"poster_image_url" validate:"required"`
	PhotoURLs          []string            `json:"photo_urls" bson:"photo_urls" validate:"required"`
	TrailerURL         string              `json:"trailer_url" bson:"trailer_url" validate:"required"`
	Runtime            time.Duration       `json:"runtime" bson:"runtime" validate:"required"`
	ReleaseDate        time.Time           `json:"release_date" bson:"release_date" validate:"required"`
	ParentalGuidances  []*ParentalGuidance `json:"parental_guide" bson:"parental_guide" validate:"required"`
	Dub                Language            `json:"dub" bson:"dub" validate:"required"`
	AvailableSubtitles []Language          `json:"available_subtitles" bson:"available_subtitles" validate:"required"`
	// Crew details
	Cast              []*People `json:"cast" bson:"cast" validate:"required"`
	Director          *People   `json:"director" bson:"director" validate:"required"`
	Writer            *People   `json:"writer" bson:"writer" validate:"required"`
	ProductionCompany string    `json:"production_company" bson:"production_company" validate:"required"`
}

// SearchMovie is a struct to search for movies
// Note: This struct differs from the one used in the [movie-search-service](github.com/harmonify/movie-reservation-system/movie-search-service/internal/core/entity/movie.go)
type SearchMovie struct {
	MovieIDs        []string       // movie IDs to include in the search
	Genre           sql.NullString // genre to search for in the movie genres
	Keyword         sql.NullString // keyword to search for in the movie text fields
	ReleaseDateFrom sql.NullTime
	ReleaseDateTo   sql.NullTime
	SortBy          MovieSortBy // sorting order of the movies. If keyword is provided, this setting will be ignored and the result will be sorted by relevance
	Page            int64       // page number
	PageSize        int64       // number of movies to return per page
}

type MovieSortBy string

const (
	MovieSortByRelevance       MovieSortBy = "relevance"
	MovieSortByReleaseDateDesc MovieSortBy = "release_date_desc"
	MovieSortByTitleAsc        MovieSortBy = "title_asc"
	MovieSortByTitleDesc       MovieSortBy = "title_desc"
)

type SearchMovieResult struct {
	// Administration data
	MovieID   bson.ObjectID `json:"movie_id" bson:"_id"`
	TraceID   string        `json:"trace_id" bson:"trace_id"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
	// Movie details
	Title              string              `json:"title" bson:"title"`
	Description        string              `json:"description" bson:"description"`
	Genres             []string            `json:"genres" bson:"genres"`
	PosterImageURL     string              `json:"poster_image_url" bson:"poster_image_url"`
	PhotoURLs          []string            `json:"photo_urls" bson:"photo_urls"`
	TrailerURL         string              `json:"trailer_url" bson:"trailer_url"`
	Runtime            time.Duration       `json:"runtime" bson:"runtime"`
	ReleaseDate        time.Time           `json:"release_date" bson:"release_date"`
	ParentalGuidances  []*ParentalGuidance `json:"parental_guide" bson:"parental_guide"`
	Dub                Language            `json:"dub" bson:"dub"`
	AvailableSubtitles []Language          `json:"available_subtitles" bson:"available_subtitles"`
	// Crew details
	Cast              []*People `json:"cast" bson:"cast"`
	Director          *People   `json:"director" bson:"director"`
	Writer            *People   `json:"writer" bson:"writer"`
	ProductionCompany string    `json:"production_company" bson:"production_company"`
	// Search result details
	Score float64 `json:"score" bson:"score"`
}
