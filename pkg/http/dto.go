package http_pkg

type (
	BaseResponseSchema struct {
		Success  bool        `json:"success"`
		TraceId  string      `json:"traceId"`
		Error    interface{} `json:"error"`    // BaseErrorResponseSchema
		Metadata interface{} `json:"metadata"` // could be PaginationMetadataSchema
		Result   interface{} `json:"result"`
	}

	BaseErrorResponseSchema struct {
		Original error       `json:"-"`
		Code     string      `json:"code"`
		Message  string      `json:"message"`
		Errors   interface{} `json:"errors"` // BaseValidationErrorSchema
	}

	BaseValidationErrorSchema struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	PaginationMetadataSchema struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"totalPages"`
	}
)
