package theater_rest

type (
	AdminSearchTheaterRequestQuery struct {
		Keyword   string  `json:"keyword" form:"keyword"`
		Latitude  float32 `json:"latitude" form:"latitude" validate:"gte=-90,lte=90"`     // user's latitude
		Longitude float32 `json:"longitude" form:"longitude" validate:"gte=-180,lte=180"` // user's longitude
		Radius    float32 `json:"radius" form:"radius" validate:"gte=0"`                  // search radius in meters
		Page      uint32  `json:"page" form:"page" validate:"gte=1"`
		PageSize  uint32  `json:"page_size" form:"page_size" validate:"gte=1"`
		SortBy    string  `json:"sort_by" form:"sort_by" validate:"oneof=newest nearest"`
	}

	AdminPostTheaterResponse struct {
		TheaterID string `json:"theater_id"`
	}

	// AdminPutTheaterResponse entity.Theater

	// AdminDeleteTheaterResponse entity.Theater
)
