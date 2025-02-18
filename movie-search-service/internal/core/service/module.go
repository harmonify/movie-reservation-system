package service

import (
	movie_service "github.com/harmonify/movie-reservation-system/movie-search-service/internal/core/service/movie"
	"go.uber.org/fx"
)

var (
	ServiceModule = fx.Module(
		"service",
		fx.Provide(
			movie_service.NewMovieSearchService,
		),
	)
)
