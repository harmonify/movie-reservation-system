package util

import (
	"github.com/harmonify/movie-reservation-system/pkg/util/encryption"
	"github.com/harmonify/movie-reservation-system/pkg/util/formatter"
	generator_util "github.com/harmonify/movie-reservation-system/pkg/util/generator"
	http_util "github.com/harmonify/movie-reservation-system/pkg/util/http"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	struct_util "github.com/harmonify/movie-reservation-system/pkg/util/struct"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"
	"go.uber.org/fx"
)

type Util struct {
	EncryptionUtil *encryption.Encryption
	FormatterUtil  formatter.FormatterUtil
	GeneratorUtil  generator_util.GeneratorUtil
	HttpUtil       http_util.HttpUtil
	JWTUtil        jwt_util.JwtUtil
	StructUtil     struct_util.StructUtil
	ValidationUtil validation.Validator
}

func NewUtil(
	encryptionUtil *encryption.Encryption,
	formatterUtil formatter.FormatterUtil,
	generatorUtil generator_util.GeneratorUtil,
	httpUtil http_util.HttpUtil,
	jwtUtil jwt_util.JwtUtil,
	structUtil struct_util.StructUtil,
	validationUtil validation.Validator,
) *Util {
	return &Util{
		EncryptionUtil: encryptionUtil,
		FormatterUtil:  formatterUtil,
		GeneratorUtil:  generatorUtil,
		HttpUtil:       httpUtil,
		JWTUtil:        jwtUtil,
		StructUtil:     structUtil,
		ValidationUtil: validationUtil,
	}
}

var (
	UtilModule = fx.Module(
		"util",
		generator_util.GeneratorUtilModule,
		encryption.EncryptionModule,
		formatter.FormatterModule,
		http_util.HttpUtilModule,
		jwt_util.JWTUtilModule,
		struct_util.StructUtilModule,
		validation.ValidationUtilModule,
		fx.Provide(NewUtil),
	)
)
