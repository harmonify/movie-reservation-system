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
