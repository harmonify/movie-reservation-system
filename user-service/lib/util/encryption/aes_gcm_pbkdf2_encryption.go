package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	generator_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/generator"
	"go.uber.org/fx"
	"golang.org/x/crypto/pbkdf2"
)

// AES encryption with AES-GCM and PBKDF2 key derivation.
type AESEncryption interface {
	Encrypt(text string) (string, error)
	Decrypt(cipherTextCompleteBase64 string) (string, error)
}

type AESEncryptionParam struct {
	fx.In

	Config        *config.Config
	GeneratorUtil generator_util.GeneratorUtil
}

type AesGcmPbkdf2EncryptionConfig struct {
	PBKDF2Iterations int
}

type AesGcmPbkdf2EncryptionImpl struct {
	generatorUtil    generator_util.GeneratorUtil
	PBKDF2Iterations int
	secret           string
}

type AesGcmPbkdf2CipherParam struct {
	Salt             []byte // 32 bytes salt
	PBKDF2Iterations int    // using the class PBKDF2Iterations when encrypting, and the value from the cipher text when decrypting
}

func NewAESEncryption(p AESEncryptionParam, cfg *AesGcmPbkdf2EncryptionConfig) (AESEncryption, error) {
	if p.Config.AppSecret == "" {
		return nil, fmt.Errorf("app secret is required")
	}
	return &AesGcmPbkdf2EncryptionImpl{
		generatorUtil:    p.GeneratorUtil,
		secret:           p.Config.AppSecret,
		PBKDF2Iterations: cfg.PBKDF2Iterations,
	}, nil
}

func (i *AesGcmPbkdf2EncryptionImpl) buildAesGcmCipher(p AesGcmPbkdf2CipherParam) (cipher.AEAD, error) {
	// Derive key with PBKDF2
	key := []byte(pbkdf2.Key([]byte(i.secret), p.Salt, p.PBKDF2Iterations, 32, sha256.New))

	// Build block cipher wrapped in GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesGCM, nil
}

func (i *AesGcmPbkdf2EncryptionImpl) Encrypt(text string) (string, error) {
	// Generate salt
	salt, err := i.generatorUtil.GenerateRandomBytes(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt")
	}

	// Build cipher
	aesGCM, err := i.buildAesGcmCipher(AesGcmPbkdf2CipherParam{
		Salt:             salt,
		PBKDF2Iterations: i.PBKDF2Iterations,
	})
	if err != nil {
		return "", fmt.Errorf("failed to build cipher. Error: %s", err.Error())
	}

	// Generate nonce
	nonce, err := i.generatorUtil.GenerateRandomBytes(uint32(aesGCM.NonceSize()))
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce")
	}

	// Encrypt value
	cipherText := aesGCM.Seal(nil, nonce, []byte(text), nil)

	cipherTextBase64 := base64.StdEncoding.EncodeToString(cipherText)
	saltbase64 := base64.StdEncoding.EncodeToString(salt)
	nonceBase64 := base64.StdEncoding.EncodeToString(nonce)

	cipherTextCompleteBase64 := strings.Join([]string{cipherTextBase64, saltbase64, nonceBase64, strconv.Itoa(i.PBKDF2Iterations)}, ".")
	return cipherTextCompleteBase64, nil
}

func (i *AesGcmPbkdf2EncryptionImpl) Decrypt(cipherTextCompleteBase64 string) (string, error) {
	data := strings.Split(cipherTextCompleteBase64, ".")
	if len(data) < 4 {
		return "", fmt.Errorf("invalid cipher text format")
	}

	cipherText, err := base64.StdEncoding.DecodeString(data[0])
	if err != nil {
		return "", err
	}

	salt, err := base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		return "", err
	}

	nonce, err := base64.StdEncoding.DecodeString(data[2])
	if err != nil {
		return "", err
	}

	pbkdf2Iterations, err := strconv.Atoi(data[3])
	if err != nil {
		return "", err
	}

	// Build cipher
	aesGCM, err := i.buildAesGcmCipher(AesGcmPbkdf2CipherParam{
		Salt:             salt,
		PBKDF2Iterations: pbkdf2Iterations,
	})
	if err != nil {
		return "", fmt.Errorf("failed to build cipher. Error: %s", err.Error())
	}

	// Decrypt
	plaintext, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
