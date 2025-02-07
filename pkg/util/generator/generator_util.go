package generator_util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
)

type GeneratorUtil interface {
	GenerateRandomBytes(n uint32) ([]byte, error)
	GenerateRandomHex(n uint32) (string, error)
	GenerateRandomBase64(n uint32) (string, error)
	GenerateRandomNumber(length uint32) (string, error)
}

type generatorUtilImpl struct{}

func NewGeneratorUtil() GeneratorUtil {
	return &generatorUtilImpl{}
}

func (s *generatorUtilImpl) GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *generatorUtilImpl) GenerateRandomHex(n uint32) (string, error) {
	bytes, err := s.GenerateRandomBytes(n)
	hexString := hex.EncodeToString(bytes)
	return hexString, err
}

func (s *generatorUtilImpl) GenerateRandomBase64(n uint32) (string, error) {
	bytes, err := s.GenerateRandomBytes(n)
	base64String := base64.RawStdEncoding.EncodeToString(bytes)
	return base64String, err
}

func (s *generatorUtilImpl) GenerateRandomNumber(length uint32) (string, error) {
	// Digits allowed in the random number
	const DIGITS = "0123456789"

	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(DIGITS))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %v", err)
		}
		result[i] = DIGITS[num.Int64()]
	}

	return string(result), nil
}
