package encryption

import (
	"golang.org/x/crypto/bcrypt"
)

// Deprecated: Use Argon2Hasher instead.
type BcryptHasher interface {
	Hash(value string) (string, error)
	Compare(hashedValue, currValue string) (bool, error)
}

type bcryptHasherImpl struct{}

func NewBcryptHasher() BcryptHasher {
	return &bcryptHasherImpl{}
}

// Deprecated: Use Argon2Hasher instead.
func (h *bcryptHasherImpl) Hash(value string) (string, error) {
	var passwordBytes = []byte(value)

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)

	return string(hashedPasswordBytes), err
}

// Deprecated: Use Argon2Hasher instead.
func (h *bcryptHasherImpl) Compare(hash, value string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(value),
	)

	return err == nil, err
}
