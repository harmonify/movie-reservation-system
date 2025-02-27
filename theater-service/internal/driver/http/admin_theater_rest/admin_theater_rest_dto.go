package admin_theater_rest

type (
	AdminSearchTheaterRequestQuery struct {
		Keyword   string  `json:"keyword" form:"keyword"`
		Page      uint32  `json:"page" form:"page" validate:"gte=1"`
		PageSize  uint32  `json:"page_size" form:"page_size" validate:"gte=1"`
		SortBy    string  `json:"sort_by" form:"sort_by" validate:"oneof=newest nearest"`
	}

	AdminPostTheaterResponse struct {
		TheaterID string `json:"theater_id"`
	}
)
