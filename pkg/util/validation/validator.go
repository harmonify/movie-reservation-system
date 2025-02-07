package validation

import (
	"regexp"
)

var (
	// This regex allows for country codes, spaces, dashes, parentheses, and extensions
	phoneRegex = regexp.MustCompile(`^(\+?[1-9]\d{0,2})?[-.●\s]?\(?\d{1,4}\)?[-.●\s]?\d{1,4}[-.●\s]?\d{1,9}(?:\s?(ext|x|extension)\s?\d{1,5})?$`)
	// E.164 regex: starts with a '+' followed by 1 to 15 digits
	e164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	// Regular expression to match alphabetic characters and spaces
	alphaSpaceRegex = regexp.MustCompile(`^[a-zA-Z\s]+$`)
)

type Validator interface {
	ValidatePhoneNumber(value string) bool
	ValidateE164PhoneNumber(value string) bool
	ValidateAlphaSpace(value string) bool
}

type validatorImpl struct{}

func NewValidator() Validator {
	return &validatorImpl{}
}

func (v *validatorImpl) ValidatePhoneNumber(value string) bool {
	return phoneRegex.MatchString(value)
}

func (v *validatorImpl) ValidateE164PhoneNumber(value string) bool {
	return e164Regex.MatchString(value)
}

func (v *validatorImpl) ValidateAlphaSpace(value string) bool {
	return alphaSpaceRegex.MatchString(value)
}
