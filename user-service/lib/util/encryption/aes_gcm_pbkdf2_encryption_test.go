package encryption_test

import (
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/test"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util/encryption"
	generator_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/generator"
	"github.com/stretchr/testify/suite"
)

func TestAesGcmPbkdf2Encryption(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(AesGcmPbkdf2EncryptionTestSuite))
}

type AesGcmPbkdf2EncryptionTestSuite struct {
	suite.Suite
	app                    any
	config                 *config.Config
	generatorUtil          generator_util.GeneratorUtil
	aesGcmPbkdf2Encryption *encryption.AesGcmPbkdf2EncryptionImpl
}

func (s *AesGcmPbkdf2EncryptionTestSuite) SetupSuite() {
	s.app = test.NewTestApp(s.invoker, s.mock()...)

	s.aesGcmPbkdf2Encryption = encryption.NewAesGcmPbkdf2Encryption(encryption.AesGcmPbkdf2EncryptionParam{
		GeneratorUtil:    s.generatorUtil,
		PBKDF2Iterations: int(15000),
		Secret:           s.config.AppSecret,
	})
}

func (s *AesGcmPbkdf2EncryptionTestSuite) invoker(
	cfg *config.Config,
	generatorUtil generator_util.GeneratorUtil,
) {
	s.config = &config.Config{
		AppName:   "RSA Encryption Tester",
		AppSecret: "1234567891123456",
	}
	s.generatorUtil = generatorUtil
}

func (s *AesGcmPbkdf2EncryptionTestSuite) mock() []any {
	// s.mockExample = mocks.NewExample(s.T())
	return []any{
		// func() interfaces.Example { return s.mockExample },
	}
}

func (s *AesGcmPbkdf2EncryptionTestSuite) TestEncryption() {
	text := "The quick brown fox jumps over the lazy dog"

	cipherText, err := s.aesGcmPbkdf2Encryption.Encrypt(text)
	s.Require().Nil(err)

	plainText, err := s.aesGcmPbkdf2Encryption.Decrypt(cipherText)
	s.Require().Nil(err)

	s.Require().Equal(text, plainText)
}
