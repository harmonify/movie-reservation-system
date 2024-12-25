package util

import (
	"github.com/harmonify/movie-reservation-system/user-service/lib/util/encryption"
	generator_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/generator"
	http_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/http"
	jwt_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/jwt"
	struct_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/struct"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util/validation"
	"go.uber.org/fx"
)

type Util struct {
	EncryptionUtil *encryption.Encryption
	GeneratorUtil  generator_util.GeneratorUtil
	HttpUtil       http_util.HttpUtil
	JWTUtil        jwt_util.JwtUtil
	StructUtil     struct_util.StructUtil
	ValidationUtil validation.ValidationUtil
}

func NewUtil(
	encryptionUtil *encryption.Encryption,
	generatorUtil generator_util.GeneratorUtil,
	httpUtil http_util.HttpUtil,
	jwtUtil jwt_util.JwtUtil,
	structUtil struct_util.StructUtil,
	validationUtil validation.ValidationUtil,
) *Util {
	return &Util{
		EncryptionUtil: encryptionUtil,
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
		http_util.HttpUtilModule,
		jwt_util.JWTUtilModule,
		struct_util.StructUtilModule,
		validation.ValidationUtilModule,
		fx.Provide(NewUtil),
	)
)
