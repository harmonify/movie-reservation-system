package encryption

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHash interface {
	Hash(value string) (string, error)
	Compare(hashedValue, currValue string) (bool, error)
}

type bcryptHashImpl struct{}

func NewBcryptHash() BcryptHash {
	return &bcryptHashImpl{}
}

func (h *bcryptHashImpl) Hash(value string) (string, error) {
	var passwordBytes = []byte(value)

	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)

	return string(hashedPasswordBytes), err
}

func (h *bcryptHashImpl) Compare(hashedValue, currValue string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedValue), []byte(currValue))

	return err == nil, err
}
