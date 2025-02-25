package http_pkg

type (
	ResponseBodySchema struct {
		Success  bool                     `json:"success"`
		TraceId  string                   `json:"traceId"`
		Error    *ResponseBodyErrorSchema `json:"error"`
		Metadata interface{}              `json:"metadata"`
		Result   interface{}              `json:"result"`
	}

	ResponseBodyErrorSchema struct {
		Original error   `json:"-"`
		Code     string  `json:"code,omitempty"`
		Message  string  `json:"message,omitempty"`
		Errors   []error `json:"errors"`
	}

	PaginationMetadataSchema struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"totalPages"`
	}

	ResponseSchema struct {
		HttpCode int
		Body     *ResponseBodySchema
	}
)
