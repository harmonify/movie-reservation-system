package http_pkg

type (
	Response struct {
		Success  bool        `json:"success"`
		TraceId  string      `json:"traceId"`
		Error    interface{} `json:"error"`    // ErrorResponse
		Metadata interface{} `json:"metadata"` // could be PaginationMetadataSchema
		Result   interface{} `json:"result"`
	}

	ErrorResponse struct {
		Original error       `json:"-"`
		Code     string      `json:"code"`
		Message  string      `json:"message"`
		Errors   interface{} `json:"errors"` // validation.ValidationError
	}

	PaginationMetadataSchema struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"totalPages"`
	}
)
