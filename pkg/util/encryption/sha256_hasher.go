package encryption

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"

	generator_util "github.com/harmonify/movie-reservation-system/pkg/util/generator"
	"go.uber.org/fx"
)

var (
	ErrSHA256InvalidHash = errors.New("sha256: hash is not in the correct format")
)

type SHA256Hasher interface {
	Hash(value string) (string, error)
	Compare(hash, value string) (bool, error)
}

type SHA256HasherParam struct {
	fx.In

	GeneratorUtil generator_util.GeneratorUtil
}

type SHA256HasherResult struct {
	fx.Out

	SHA256Hasher SHA256Hasher
}

type sha256HasherImpl struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func NewSHA256Hasher(p SHA256HasherParam) SHA256HasherResult {
	return SHA256HasherResult{
		SHA256Hasher: &sha256HasherImpl{},
	}
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
