package formatter

import (
	"regexp"

	"go.uber.org/fx"
)

type FormatterUtil interface {
	FormatPhoneNumberToE164(phoneNumber, countryCode string) string
}

type FormatterUtilParam struct {
	fx.In
}

type FormatterUtilResult struct {
	fx.Out

	FormatterUtil FormatterUtil
}

type formatterUtilImpl struct {
}

func NewFormatterUtil(p FormatterUtilParam) FormatterUtilResult {
	return FormatterUtilResult{
		FormatterUtil: &formatterUtilImpl{},
	}
}

// FormatPhoneNumberToE164 formats a phone number into E.164 format
func (u *formatterUtilImpl) FormatPhoneNumberToE164(phoneNumber, countryCode string) string {
	// Remove any spaces, parenthesis or other punctuation.
	e164Number := regexp.MustCompile(`\D`).ReplaceAllString(phoneNumber, "")

	// Remove leading zero and prepend country code
	if e164Number[0] == '0' {
		e164Number = "62" + e164Number[1:]
	}

	// Prepend '+' to match E.164 format
	if e164Number[0] != '+' {
		e164Number = "+" + e164Number
	}

	return e164Number
}
