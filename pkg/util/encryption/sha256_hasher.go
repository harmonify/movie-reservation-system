package encryption

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

var (
	ErrSHA256InvalidHash = errors.New("sha256: hash is not in the correct format")
)

type SHA256Hasher interface {
	Hash(value string) (string, error)
	Compare(hash, value string) (bool, error)
}

type sha256HasherImpl struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func NewSHA256Hasher() SHA256Hasher {
	return &sha256HasherImpl{}
}

func (h *sha256HasherImpl) Hash(value string) (string, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write([]byte(value))
	if err != nil {
		return "", err
	}
	msgHashSum := msgHash.Sum(nil)
	return base64.RawStdEncoding.EncodeToString(msgHashSum), nil
}

func (h *sha256HasherImpl) Compare(hash, value string) (match bool, err error) {
	newHash, err := h.Hash(value)
	if err != nil {
		return false, err
	}
	return hash == newHash, nil
}
