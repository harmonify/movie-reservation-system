package http_pkg

import (
	"github.com/gin-gonic/gin"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"
)

type HttpValidator interface {
	ValidateRequestBody(c *gin.Context, schema interface{}) error
	ValidateRequestQuery(c *gin.Context, schema interface{}) error
}

type HttpValidatorImpl struct {
	errorMapper     error_pkg.ErrorMapper
	httpResponse    HttpResponse
	structValidator validation.StructValidator
}

func NewHttpValidator(
	errorMapper error_pkg.ErrorMapper,
	structValidator validation.StructValidator,
	httpResponse HttpResponse,
) HttpValidator {
	return &HttpValidatorImpl{
		errorMapper:     errorMapper,
		structValidator: structValidator,
		httpResponse:    httpResponse,
	}
}

func (v *HttpValidatorImpl) ValidateRequestBody(c *gin.Context, schema interface{}) error {
	if err := c.ShouldBind(schema); err != nil {
		_, validationErrs := v.structValidator.ConstructValidationErrorFields(err)
		errWithStack := error_pkg.NewErrorWithStack(err, error_pkg.InvalidRequestBodyError)
		return NewHttpError(errWithStack, validationErrs...)
	}

	if err, validationErrs := v.structValidator.Validate(schema); len(validationErrs) > 0 {
		_, validationErrs := v.structValidator.ConstructValidationErrorFields(err)
		errWithStack := error_pkg.NewErrorWithStack(err, error_pkg.InvalidRequestBodyError)
		return NewHttpError(errWithStack, validationErrs...)
	}

	return nil
}

func (v *HttpValidatorImpl) ValidateRequestQuery(c *gin.Context, schema interface{}) error {
	if err := c.ShouldBindQuery(schema); err != nil {
		_, validationErrs := v.structValidator.ConstructValidationErrorFields(err)
		errWithStack := error_pkg.NewErrorWithStack(err, error_pkg.InvalidRequestBodyError)
		return NewHttpError(errWithStack, validationErrs...)
	}

	if err, validationErrs := v.structValidator.Validate(schema); len(validationErrs) > 0 {
		_, validationErrs := v.structValidator.ConstructValidationErrorFields(err)
		errWithStack := error_pkg.NewErrorWithStack(err, error_pkg.InvalidRequestBodyError)
		return NewHttpError(errWithStack, validationErrs...)
	}

	return nil
}
