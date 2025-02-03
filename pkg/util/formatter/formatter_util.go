package formatter

import (
	"regexp"
)

type FormatterUtil interface {
	FormatPhoneNumberToE164(phoneNumber, countryCode string) string
}

type formatterUtilImpl struct{}

func NewFormatterUtil() FormatterUtil {
	return &formatterUtilImpl{}
}

// FormatPhoneNumberToE164 formats a phone number into E.164 format
func (*formatterUtilImpl) FormatPhoneNumberToE164(phoneNumber, countryCode string) string {
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
