package showtime_rest

type (
	AdminSearchShowtimeRequestQuery struct {
		TheaterID        string `json:"theater_id" form:"theater_id" validate:"required"`
		RoomID           string `json:"room_id" form:"room_id" validate:"required"`
		MovieID          string `json:"movie_id" form:"movie_id" validate:"required"`
		StartTimeGteUnix int64  `json:"start_time_gte_unix" form:"start_time_gte_unix" validate:"required,gte=0"`
		StartTimeLteUnix int64  `json:"start_time_lte_unix" form:"start_time_lte_unix" validate:"required,gte=0"`
		SortBy           string `json:"sort_by" form:"sort_by" validate:"oneof=latest oldest"`
		Page             uint32 `json:"page" form:"page" validate:"gte=1"`
		PageSize         uint32 `json:"page_size" form:"page_size" validate:"gte=1"`
	}

	AdminPostShowtimeResponse struct {
		ShowtimeID string `json:"showtime_id"`
	}

	// AdminPutShowtimeResponse entity.Showtime

	// AdminDeleteShowtimeResponse entity.Showtime
)
