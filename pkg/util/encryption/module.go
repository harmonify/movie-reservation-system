package encryption

import "go.uber.org/fx"

type Encryption struct {
	AESEncryption AESEncryption
	Argon2Hasher  Argon2Hasher
	BcryptHasher  BcryptHasher
	RSAEncryption RSAEncryption
	SHA256Hasher  SHA256Hasher
}

func NewEncryption(
	aesEncryption AESEncryption,
	argon2Hasher Argon2Hasher,
	bcryptHash BcryptHasher,
	rsaEncryption RSAEncryption,
	sha256Hasher SHA256Hasher,
) *Encryption {
	return &Encryption{
		AESEncryption: aesEncryption,
		Argon2Hasher:  argon2Hasher,
		BcryptHasher:  bcryptHash,
		RSAEncryption: rsaEncryption,
		SHA256Hasher:  sha256Hasher,
	}
}

var (
	EncryptionModule = fx.Module(
		"encryption-util",
		fx.Provide(
			NewAESEncryption,

			func(p Argon2HasherParam) (Argon2HasherResult, error) {
				return NewArgon2Hasher(p, *Argon2HasherDefaultConfig)
			},

			NewBcryptHasher,

			NewRSAEncryption,

			NewSHA256Hasher,

			NewEncryption,
		),
	)
)
