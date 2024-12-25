package response

type (
	BaseResponseSchema struct {
		Success  bool        `json:"success"`
		TraceId  string      `json:"traceId"`
		Error    interface{} `json:"error"`
		Metadata interface{} `json:"meta"`
		Result   interface{} `json:"result"`
	}

	BaseErrorResponseSchema struct {
		Code     string      `json:"code"`
		Message  string      `json:"message"`
		Errors   interface{} `json:"errors"`
		Original error       `json:"-"`
	}

	BaseErrorValidationSchema struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)
