package encryption

import "go.uber.org/fx"

type Encryption struct {
	AESEncryption AESEncryption
	BcryptHash    BcryptHash
}

func NewEncryption(
	aesEncryption AESEncryption,
	bcryptHash BcryptHash,
) *Encryption {
	return &Encryption{
		AESEncryption: aesEncryption,
		BcryptHash:    bcryptHash,
	}
}

var (
	EncryptionModule = fx.Provide(NewAESEncryption, NewBcryptHash, NewEncryption)
)
