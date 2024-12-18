package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

type AESEncryption interface {
	Encrypt(payload *AESPayload) (string, error)
	Decrypt(payload *AESPayload) (string, error)
}

type AESEncryptionImpl struct{}

type AESPayload struct {
	Secret  string
	Payload string
}

func NewAESEncryption() AESEncryption {
	return &AESEncryptionImpl{}
}

func (i *AESEncryptionImpl) Encrypt(payload *AESPayload) (string, error) {
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

func (i *AESEncryptionImpl) Decrypt(payload *AESPayload) (string, error) {
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
