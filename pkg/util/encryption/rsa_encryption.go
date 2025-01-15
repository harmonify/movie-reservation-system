package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"

	generator_util "github.com/harmonify/movie-reservation-system/pkg/util/generator"
	"go.uber.org/fx"
)

type RSAEncryption interface {
	Generate() (RSAKeyPair, error)
	EncodePrivateKey(privateKey *rsa.PrivateKey) []byte
	EncodePublicKey(publicKey *rsa.PublicKey) []byte
}

type RSAEncryptionParam struct {
	fx.In

	AESEncryption AESEncryption
	GeneratorUtil generator_util.GeneratorUtil
}

type RSAEncryptionResult struct {
	fx.Out

	RSAEncryption RSAEncryption
}

type rsaEncryptionImpl struct {
	aesEncryption AESEncryption
	generatorUtil generator_util.GeneratorUtil
}

func NewRSAEncryption(p RSAEncryptionParam) RSAEncryptionResult {
	return RSAEncryptionResult{
		RSAEncryption: &rsaEncryptionImpl{
			aesEncryption: p.AESEncryption,
			generatorUtil: p.GeneratorUtil,
		},
	}
}

type RSAKeyPair struct {
	// Encoded in PKCS#1 ASN.1 PEM format
	PrivateKey []byte
	// Encoded in PKCS#1 ASN.1 PEM format
	PublicKey []byte
}

type CustomReader struct {
	io.Reader
}

// Generate RSA key pair
func (i *rsaEncryptionImpl) Generate() (key RSAKeyPair, err error) {
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return key, err
	}

	key.PrivateKey = i.EncodePrivateKey(privateKey)
	key.PublicKey = i.EncodePublicKey(&privateKey.PublicKey)

	return key, nil
}

// EncodePrivateKey convert RSA private key to PKCS #1, ASN.1 DER form,
// then encode it in PEM blocks of type "RSA PRIVATE KEY"
func (i *rsaEncryptionImpl) EncodePrivateKey(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
}

// EncodePrivateKey convert RSA private key to PKCS #1, ASN.1 DER form,
// then encode it in PEM blocks of type "RSA PUBLIC KEY"
func (i *rsaEncryptionImpl) EncodePublicKey(publicKey *rsa.PublicKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(publicKey),
		},
	)
}
