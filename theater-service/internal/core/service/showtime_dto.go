package service

import "github.com/harmonify/movie-reservation-system/theater-service/internal/core/entity"

type (
	ShowtimeDetail struct {
		ShowtimeID     string                `json:"showtime_id"`
		TheaterName    string                `json:"theater_name"`
		RoomName       string                `json:"room_name"`
		MovieTitle     string                `json:"movie_title"`
		StartTime      string                `json:"start_time"`
		EndTime        string                `json:"end_time"`
		AvailableSeats []*ShowtimeSeatDetail `json:"available_seats"`
	}

	ShowtimeSeatDetail struct {
		SeatID     string            `json:"seat_iD"`
		SeatRow    string            `json:"seat_row"`
		SeatColumn string            `json:"seat_column"`
		Status     entity.SeatStatus `json:"status"`
	}
)
