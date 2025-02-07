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
		vErr := error_pkg.InvalidRequestBodyError
		var data []error
		data = v.structValidator.ConstructValidationErrorFields(err)
		data = append(data, &validation.ValidationError{
			Field:   "",
			Message: err.Error(),
		})
		vErr.Errors = data
		return vErr
	}

	if err, validationErrs := v.structValidator.Validate(schema); err != nil && len(validationErrs) > 0 {
		vErr := error_pkg.InvalidRequestBodyError
		vErr.Errors = validationErrs
		return vErr
	}

	return nil
}

func (v *HttpValidatorImpl) ValidateRequestQuery(c *gin.Context, schema interface{}) error {
	if err := c.ShouldBindQuery(schema); err != nil {
		vErr := error_pkg.InvalidRequestBodyError
		var data []error
		data = v.structValidator.ConstructValidationErrorFields(err)
		data = append(data, &validation.ValidationError{
			Field:   "",
			Message: err.Error(),
		})
		vErr.Errors = data
		return vErr
	}

	if err, validationErrs := v.structValidator.Validate(schema); err != nil && len(validationErrs) > 0 {
		vErr := error_pkg.InvalidRequestBodyError
		vErr.Errors = validationErrs
		return vErr
	}

	return nil
}
