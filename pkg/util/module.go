package util

import (
	struct_util "github.com/harmonify/movie-reservation-system/pkg/util/struct"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"
	"go.uber.org/fx"
)

type Util struct {
	StructUtil     struct_util.StructUtil
	ValidationUtil validation.ValidationUtil
}

func NewUtil(
	structUtil struct_util.StructUtil,
	validationUtil validation.ValidationUtil,
) *Util {
	return &Util{
		StructUtil:     structUtil,
		ValidationUtil: validationUtil,
	}
}

var (
	UtilModule = fx.Options(
		struct_util.StructUtilModule,
		validation.ValidationUtilModule,
	)
)
