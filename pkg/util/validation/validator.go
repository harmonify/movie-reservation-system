package validation

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
	"github.com/gobeam/stringy"
	"github.com/harmonify/movie-reservation-system/pkg/constant"
	"github.com/harmonify/movie-reservation-system/pkg/http/response"
)

type errorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Validator interface {
	ValidateBody(c *gin.Context, schema interface{}) error
	ValidateQueryParams(c *gin.Context, schema interface{}) error
}

type ValidatorImpl struct {
	response       response.Response
	validationUtil ValidationUtil
}

func NewValidator(response response.Response, validationUtil ValidationUtil) Validator {
	return &ValidatorImpl{
		response:       response,
		validationUtil: validationUtil,
	}
}

func (v *ValidatorImpl) ValidateBody(c *gin.Context, schema interface{}) error {
	if err := c.ShouldBind(schema); err != nil {
		errMsg := constructValidationField(err)
		return v.response.BuildValidationError(constant.InvalidRequestBody, err, errMsg)
	}

	validate := validator.New()
	err := v.registerCustomValidation(validate)
	if err != nil {
		return err
	}

	if err := validate.Struct(schema); err != nil {
		errMsg := constructValidationField(err)
		return v.response.BuildValidationError(constant.InvalidRequestBody, err, errMsg)
	}

	return nil
}

func (v *ValidatorImpl) ValidateQueryParams(c *gin.Context, schema interface{}) error {
	err := c.ShouldBindQuery(schema)
	if err != nil {
		errMsg := constructValidationField(err)
		return v.response.BuildValidationError(constant.InvalidRequestBody, err, errMsg)
	}

	validate := validator.New()
	err = v.registerCustomValidation(validate)
	if err != nil {
		return err
	}

	if err := validate.Struct(schema); err != nil {
		errMsg := constructValidationField(err)
		return v.response.BuildValidationError(constant.InvalidRequestBody, err, errMsg)
	}

	return nil
}

func (v *ValidatorImpl) phoneNumberValidation(fl validator.FieldLevel) bool {
	return v.validationUtil.ValidatePhoneNumber(fl.Field().String())
}

func (v *ValidatorImpl) registerCustomValidation(validate *validator.Validate) error {
	invalid_validation_err := errors.New("INVALID_VALIDATION")

	err := validate.RegisterValidation(string(constant.PhoneNumberKey), v.phoneNumberValidation)
	if err != nil {
		return invalid_validation_err
	}

	return nil
}

func constructValidationField(err error) (errorsData []errorMsg) {
	var val validator.ValidationErrors

	if errors.As(err, &val) {
		errorsData = make([]errorMsg, len(val))
		for i, fe := range val {
			fieldPath := fe.Field()

			if fe.StructNamespace() != "" {
				fieldPath = extractNestedField(fe.StructNamespace())
			}

			errorsData[i] = errorMsg{fieldPath, getErrorMsg(fe)}
		}
	}

	return
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "alpha":
		return "Should not contain number or symbol"
	case "email":
		return "Invalid email format"
	case "alphanum":
		return "Invalid format"
	case "oneof":
		return "Unknown key on requested field"
	case "min":
		return "Not meet criteria"
	case "max":
		return "Not meet criteria"
	case "number":
		return "Page format is in number"
	case "required_without":
		return "This field is required"
	case "phone_number":
		return "Phone number is not valid"
	}

	return "Unknown error"
}

func extractNestedField(fieldPath string) string {
	parts := strings.Split(fieldPath, ".")

	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}

	if len(parts) == 2 {
		return convertToSnakeCase(parts[1])
	}

	return convertToSnakeCase(parts[0])
}

func convertToSnakeCase(value string) string {
	str := stringy.New(value)
	return str.SnakeCase("?", "").ToLower()
}
