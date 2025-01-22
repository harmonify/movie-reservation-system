package shared

import (
	"fmt"
)

type InvalidPhoneNumberError struct {
	PhoneNumber string
}

func (e *InvalidPhoneNumberError) Error() string {
	return fmt.Sprintf("invalid phone number: %s", e.PhoneNumber)
}

func (e *InvalidPhoneNumberError) Is(err error) bool {
	_, ok := err.(*InvalidPhoneNumberError)
	return ok
}

func (e *InvalidPhoneNumberError) As(target interface{}) bool {
	_, ok := target.(*InvalidPhoneNumberError)
	return ok
}

func (e *InvalidPhoneNumberError) Unwrap() error {
	return e
}

func NewInvalidPhoneNumberError(phoneNumber string) error {
	return &InvalidPhoneNumberError{
		PhoneNumber: phoneNumber,
	}
}

type RateLimitError struct {
	Original   error
	RetryAfter int
}

func (e *RateLimitError) Error() string {
	return e.Original.Error()
}

func (e *RateLimitError) Is(err error) bool {
	_, ok := err.(*RateLimitError)
	return ok
}

func (e *RateLimitError) As(target interface{}) bool {
	_, ok := target.(*RateLimitError)
	return ok
}

func (e *RateLimitError) Unwrap() error {
	return e
}

func NewRateLimitError(err error, retryAfter int) error {
	return &RateLimitError{
		Original:   err,
		RetryAfter: retryAfter,
	}
}
