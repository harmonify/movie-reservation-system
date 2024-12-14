package validation

import (
	"regexp"
)

type ValidationUtil interface {
	ValidatePhoneNumber(value string) bool
}

type ValidationUtilImpl struct{}

func NewValidationUtil() ValidationUtil {
	return &ValidationUtilImpl{}
}

func (v *ValidationUtilImpl) ValidatePhoneNumber(value string) bool {
	// Define an ultimate regex for phone numbers
	// This regex allows for country codes, spaces, dashes, parentheses, and extensions
	phoneRegex := `^(\+?[1-9]\d{0,2})?[-.●\s]?\(?\d{1,4}\)?[-.●\s]?\d{1,4}[-.●\s]?\d{1,9}(?:\s?(ext|x|extension)\s?\d{1,5})?$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(value)
}
