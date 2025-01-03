package validation

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gobeam/stringy"
	http_constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
)

type StructValidator interface {
	Validate(schema interface{}) (original error, errorFields []BaseValidationErrorSchema)
	ConstructValidationErrorFields(err error) (processed bool, errorFields []BaseValidationErrorSchema)
}

type structValidatorImpl struct {
	uni       *ut.UniversalTranslator
	trans     ut.Translator
	validate  *validator.Validate
	validator Validator
}

func NewStructValidator(validationUtil Validator) (StructValidator, error) {
	validate := validator.New()

	structValidator := &structValidatorImpl{
		validate: validate,
	}

	err := structValidator.registerTranslations()
	if err != nil {
		return nil, err
	}

	err = structValidator.registerCustomValidations()
	if err != nil {
		return nil, err
	}

	return structValidator, nil
}

func (v *structValidatorImpl) registerTranslations() error {
	en := en.New()
	v.uni = ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	v.trans, _ = v.uni.GetTranslator("en")
	err := en_translations.RegisterDefaultTranslations(v.validate, v.trans)
	if err != nil {
		return err
	}

	return nil
}

func (v *structValidatorImpl) registerCustomValidations() error {
	err := v.validate.RegisterValidation(
		string(http_constant.PhoneNumberKey),
		func(fl validator.FieldLevel) bool {
			return v.validator.ValidatePhoneNumber(fl.Field().String())
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (v *structValidatorImpl) Validate(schema interface{}) (original error, errorFields []BaseValidationErrorSchema) {
	if err := v.validate.Struct(schema); err != nil {
		_, errFields := v.ConstructValidationErrorFields(err)
		return err, errFields
	}
	return nil, nil
}

// ConstructValidationErrorFields constructs validation error fields
// Accepts error (will only process if the type is validator.ValidationErrors)
// Returns boolean (true if error is validator.ValidationErrors) and array of constructed error fields
func (v *structValidatorImpl) ConstructValidationErrorFields(err error) (processed bool, errorFields []BaseValidationErrorSchema) {
	var val validator.ValidationErrors

	processed = errors.As(err, &val)
	if processed {
		errorFields = make([]BaseValidationErrorSchema, len(val))
		for i, fe := range val {
			// Use tag name whenever possible
			fieldPath := fe.Field()

			// Fallback to struct namespace if tag name is not available
			if fieldPath == "" && fe.StructNamespace() != "" {
				fieldPath = v.extractNestedField(fe.StructNamespace())
			}

			// Construct validation error fields
			errorFields[i] = BaseValidationErrorSchema{
				Field:   fieldPath,
				Message: fe.Translate(v.trans),
			}
		}
	}

	return
}

func (v *structValidatorImpl) extractNestedField(fieldPath string) string {
	parts := strings.Split(fieldPath, ".")

	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}

	if len(parts) == 2 {
		return v.convertToSnakeCase(parts[1])
	}

	return v.convertToSnakeCase(parts[0])
}

func (v *structValidatorImpl) convertToSnakeCase(value string) string {
	str := stringy.New(value)
	return str.SnakeCase("?", "").ToLower()
}
