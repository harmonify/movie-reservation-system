package service

import (
	movie_service "github.com/harmonify/movie-reservation-system/movie-service/internal/core/service/movie"
	"go.uber.org/fx"
)

var (
	AdminMovieServiceModule = fx.Module(
		"movie-service",
		fx.Provide(
			movie_service.NewAdminMovieService,
		),
	)

	ServiceModule = fx.Module(
		"service",
		AdminMovieServiceModule,
	)
)
