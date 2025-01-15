package validation

type (
	BaseValidationErrorSchema struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)
