package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

func NewAesCfbEncryption() *aesCFBEncryptionImpl {
	return &aesCFBEncryptionImpl{}
}

// Deprecated: Use AES-GCM with PBKDF2 instead.
type aesCFBEncryptionImpl struct{}

// Deprecated: Use AES-GCM with PBKDF2 instead.
func (i *aesCFBEncryptionImpl) Encrypt(payload *AESPayload) (string, error) {
	key, err := hex.DecodeString(payload.Secret)
	if err != nil {
		fmt.Println("Invalid AES key:", err)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	strBytes := []byte(payload.Payload)
	cfb := cipher.NewCFBEncrypter(block, key[:aes.BlockSize])
	cipherText := make([]byte, len(strBytes))
	cfb.XORKeyStream(cipherText, strBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Deprecated: Use AES-GCM with PBKDF2 instead.
func (i *aesCFBEncryptionImpl) Decrypt(payload *AESPayload) (string, error) {
	key, err := hex.DecodeString(payload.Secret)
	if err != nil {
		fmt.Println("Invalid AES key:", err)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	strBytes, err := base64.StdEncoding.DecodeString(payload.Payload)
	if err != nil {
		return "", errors.New("INVALID_BASE64")
	}

	cfb := cipher.NewCFBDecrypter(block, key[:aes.BlockSize])
	cipherText := make([]byte, len(strBytes))
	cfb.XORKeyStream(cipherText, strBytes)

	return string(cipherText), nil

}
