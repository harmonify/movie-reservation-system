package validation

import (
	"regexp"
)

type ValidationUtil interface {
	ValidatePhoneNumber(value string) bool
	// Validate if the phone number is in E.164 format
	ValidateE164PhoneNumber(value string) bool
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

func (v *ValidationUtilImpl) ValidateE164PhoneNumber(value string) bool {
	// E.164 regex: starts with a '+' followed by 1 to 15 digits
	e164Regex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return e164Regex.MatchString(value)
}
