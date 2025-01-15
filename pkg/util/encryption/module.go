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
		"encryption",
		fx.Provide(
			func() *AesGcmPbkdf2EncryptionConfig {
				return &AesGcmPbkdf2EncryptionConfig{
					PBKDF2Iterations: int(15000),
				}
			},
			NewAESEncryption,

			func() Argon2HasherConfig {
				return *Argon2HasherDefaultConfig
			},
			NewArgon2Hasher,

			NewBcryptHasher,

			NewRSAEncryption,

			NewSHA256Hasher,

			NewEncryption,
		),
	)
)
