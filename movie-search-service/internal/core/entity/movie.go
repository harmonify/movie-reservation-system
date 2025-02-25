package entity

import (
	"database/sql"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Movie struct {
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
}

type Language struct {
	Name string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
}

func (l *Language) String() string {
	return l.Code
}

type ParentalGuidance struct {
	Code        string `json:"code" bson:"code"`
	CountryCode string `json:"country_code" bson:"country_code"` // ISO 3166-1 alpha-2 country code
}

func (pg *ParentalGuidance) String() string {
	return pg.CountryCode + ":" + pg.Code
}

// SearchMovieParam is a struct to search for movies
// Note: This struct differs from the one used in the [movie-service](github.com/harmonify/movie-reservation-system/movie-service/internal/core/entity/movie.go)
type SearchMovieParam struct {
	MovieIDs          []string       // movie IDs to include in the search
	Genre             sql.NullString // genre to search for in the movie genres
	Keyword           sql.NullString // keyword to search for in the movie text fields
	SortBy            MovieSortBy    // sorting order of the movies. If keyword is provided, this setting will be ignored and the result will be sorted by relevance
	Limit             int64          // number of movies to return
	LastSeenSortValue interface{}    // last seen sort value for pagination
	LastSeenID        string         // last seen movie ID for pagination
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
