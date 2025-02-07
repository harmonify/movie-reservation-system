package encryption

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

var (
	ErrSHA256InvalidHash = errors.New("sha256: hash is not in the correct format")
)

type SHA256HasherConfig struct {
	AppSecret string `validate:"required"`
}

type SHA256Hasher interface {
	Hash(value string) (string, error)
	Compare(hash, value string) (bool, error)
}

type sha256HasherImpl struct {
	cfg *SHA256HasherConfig
}

func NewSHA256Hasher(cfg *SHA256HasherConfig) SHA256Hasher {
	return &sha256HasherImpl{
		cfg: cfg,
	}
}

func (h *sha256HasherImpl) Hash(value string) (string, error) {
	hash := hmac.New(sha256.New, []byte(h.cfg.AppSecret))
	hash.Write([]byte(value))
	return base64.RawStdEncoding.EncodeToString(hash.Sum(nil)), nil
}

func (h *sha256HasherImpl) Compare(hash, value string) (match bool, err error) {
	newHash, err := h.Hash(value)
	if err != nil {
		return false, err
	}
	return hash == newHash, nil
}
