package validator

import (
	"github.com/gin-gonic/gin"
	error_constant "github.com/harmonify/movie-reservation-system/user-service/lib/error/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/response"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util/validation"
)

type HttpValidator interface {
	ValidateRequestBody(c *gin.Context, schema interface{}) error
	ValidateRequestQuery(c *gin.Context, schema interface{}) error
}

type HttpValidatorImpl struct {
	response        response.HttpResponse
	structValidator validation.StructValidator
}

func NewHttpValidator(
	structValidator validation.StructValidator,
	response response.HttpResponse,
) HttpValidator {

	return &HttpValidatorImpl{
		structValidator: structValidator,
		response:        response,
	}
}

func (v *HttpValidatorImpl) ValidateRequestBody(c *gin.Context, schema interface{}) error {
	if err := c.ShouldBind(schema); err != nil {
		_, errFields := v.structValidator.ConstructValidationErrorFields(err)
		return v.response.BuildValidationError(error_constant.InvalidRequestBody, err, errFields)
	}

	if err, errFields := v.structValidator.Validate(schema); len(errFields) > 0 {
		return v.response.BuildValidationError(error_constant.InvalidRequestBody, err, errFields)
	}

	return nil
}

func (v *HttpValidatorImpl) ValidateRequestQuery(c *gin.Context, schema interface{}) error {
	if err := c.ShouldBindQuery(schema); err != nil {
		_, errFields := v.structValidator.ConstructValidationErrorFields(err)
		return v.response.BuildValidationError(error_constant.InvalidRequestBody, err, errFields)
	}

	if err, errFields := v.structValidator.Validate(schema); len(errFields) > 0 {
		return v.response.BuildValidationError(error_constant.InvalidRequestBody, err, errFields)
	}

	return nil
}
