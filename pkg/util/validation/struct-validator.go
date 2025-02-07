package validation

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gobeam/stringy"
)

type StructValidator interface {
	Validate(schema interface{}) (original error, errorFields []error)
	ConstructValidationErrorFields(err error) []error
}

type structValidatorImpl struct {
	uni           *ut.UniversalTranslator
	trans         ut.Translator
	validator     *validator.Validate
	validatorUtil Validator
}

func NewStructValidator(validatorUtil Validator) (StructValidator, error) {
	structValidator := &structValidatorImpl{
		validator:     validator.New(validator.WithRequiredStructEnabled()),
		validatorUtil: validatorUtil,
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
	err := en_translations.RegisterDefaultTranslations(v.validator, v.trans)
	if err != nil {
		return err
	}

	return nil
}

func (v *structValidatorImpl) registerCustomValidations() error {
	err := v.validator.RegisterValidation(
		PhoneNumberKey.String(),
		func(fl validator.FieldLevel) bool {
			return v.validatorUtil.ValidatePhoneNumber(fl.Field().String())
		},
	)
	if err != nil {
		return err
	}

	err = v.validator.RegisterValidation(
		AlphaSpaceKey.String(),
		func(fl validator.FieldLevel) bool {
			return v.validatorUtil.ValidateAlphaSpace(fl.Field().String())
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (v *structValidatorImpl) Validate(schema interface{}) (error, []error) {
	if err := v.validator.Struct(schema); err != nil {
		errFields := v.ConstructValidationErrorFields(err)
		return err, errFields
	}
	return nil, nil
}

// ConstructValidationErrorFields constructs validation error fields
// Accepts error (will only process if the type is validator.ValidationErrors)
// Returns boolean (true if error is validator.ValidationErrors) and array of constructed error fields
func (v *structValidatorImpl) ConstructValidationErrorFields(err error) (errorFields []error) {
	var val validator.ValidationErrors
	if errors.As(err, &val) {
		errorFields = make([]error, len(val))
		for i, fe := range val {
			// Use tag name whenever possible
			fieldPath := v.extractNestedField(fe.Namespace())

			// Fall back to struct namespace if tag name is not available
			if fieldPath == "" && fe.StructNamespace() != "" {
				fieldPath = v.extractNestedField(fe.StructNamespace())
			}

			// Construct validation error fields
			errorFields[i] = &ValidationError{
				Field:   fieldPath,
				Message: fe.Translate(v.trans),
			}
		}
	} else {
		errorFields = make([]error, 0)
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
	if value == "" {
		return value
	}
	str := stringy.New(value)
	return str.SnakeCase("?", "").ToLower()
}
