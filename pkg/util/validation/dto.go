package validation

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *ValidationError) Error() string {
	return v.Message
}
