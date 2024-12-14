package validation

import "go.uber.org/fx"

var (
	ValidationUtilModule = fx.Provide(NewValidationUtil, NewValidator)
)
