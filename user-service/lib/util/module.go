package util

import (
	"github.com/harmonify/movie-reservation-system/user-service/lib/util/encryption"
	http_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/http"
	jwt_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/jwt"
	struct_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/struct"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util/validation"
	"go.uber.org/fx"
)

type Util struct {
	EncryptionUtil encryption.Encryption
	HttpUtil       http_util.HttpUtil
	JWTUtil        jwt_util.JWTUtil
	StructUtil     struct_util.StructUtil
	ValidationUtil validation.ValidationUtil
}

func NewUtil(
	encryptionUtil encryption.Encryption,
	httpUtil http_util.HttpUtil,
	jwtUtil jwt_util.JWTUtil,
	structUtil struct_util.StructUtil,
	validationUtil validation.ValidationUtil,
) *Util {
	return &Util{
		EncryptionUtil: encryptionUtil,
		HttpUtil:       httpUtil,
		JWTUtil:        jwtUtil,
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
