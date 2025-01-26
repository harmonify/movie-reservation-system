package http_pkg

import (
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"
)

type HttpError struct {
	*error_pkg.ErrorWithStack
	ValidationErrors interface{} `json:"validation_errors"` // user-friendly validation errors
}

func (e *HttpError) Error() string {
	return e.ErrorWithStack.Error()
}

func (e *HttpError) As(target interface{}) bool {
	return e.ErrorWithStack.As(target)
}

func (e *HttpError) Unwrap() error {
	return e.ErrorWithStack
}

func NewHttpError(err *error_pkg.ErrorWithStack, validationErrors ...validation.ValidationError) *HttpError {
	return &HttpError{
		ErrorWithStack:   err,
		ValidationErrors: validationErrors,
	}
}
