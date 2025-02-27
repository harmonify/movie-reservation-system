package service

import (
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
)

var (
	TheaterNotFoundError = &error_pkg.ErrorWithDetails{
		Code:     "THEATER_NOT_FOUND",
		Message:  "theater not found",
		HttpCode: 404,
		GrpcCode: 5,
	}

	TheaterExistsError = &error_pkg.ErrorWithDetails{
		Code:     "THEATER_EXISTS",
		Message:  "theater already exists",
		HttpCode: 409,
		GrpcCode: 6,
	}

	TheaterIDRequiredError = &error_pkg.ErrorWithDetails{
		Code:     "THEATER_ID_REQUIRED",
		Message:  "theater id is required",
		HttpCode: 400,
		GrpcCode: 3,
	}

	RoomNotFoundError = &error_pkg.ErrorWithDetails{
		Code:     "ROOM_NOT_FOUND",
		Message:  "room not found",
		HttpCode: 404,
		GrpcCode: 5,
	}

	MovieIDRequiredError = &error_pkg.ErrorWithDetails{
		Code:     "MOVIE_ID_REQUIRED",
		Message:  "movie id is required",
		HttpCode: 400,
		GrpcCode: 3,
	}

	ShowtimeNotFoundError = &error_pkg.ErrorWithDetails{
		Code:     "SHOWTIME_NOT_FOUND",
		Message:  "showtime not found",
		HttpCode: 404,
		GrpcCode: 5,
	}

	ShowtimeIDRequiredError = &error_pkg.ErrorWithDetails{
		Code:     "SHOWTIME_ID_REQUIRED",
		Message:  "showtime id is required",
		HttpCode: 400,
		GrpcCode: 3,
	}

	ShowtimeTicketsNotFoundError = &error_pkg.ErrorWithDetails{
		Code:     "SHOWTIME_TICKETS_NOT_FOUND",
		Message:  "showtime tickets not found",
		HttpCode: 404,
		GrpcCode: 5,
	}
)
