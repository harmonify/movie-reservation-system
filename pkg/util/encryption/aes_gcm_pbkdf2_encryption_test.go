package encryption_test

import (
	"context"
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/pkg/util/encryption"
	generator_util "github.com/harmonify/movie-reservation-system/pkg/util/generator"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

func TestAESEncryption(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(AESEncryptionTestSuite))
}

type AESEncryptionTestSuite struct {
	suite.Suite
	app           *fx.App
	aesEncryption encryption.AESEncryption
}

func (s *AESEncryptionTestSuite) SetupSuite() {
	s.app = fx.New(
		generator_util.GeneratorUtilModule,
		fx.Provide(
			func() *encryption.AESEncryptionConfig {
				return &encryption.AESEncryptionConfig{
					AppSecret: "test",
				}
			},
			encryption.NewAESEncryption,
		),
		fx.Invoke(func(
			aesEncryption encryption.AESEncryption,
		) {
			s.aesEncryption = aesEncryption
		}),

		fx.NopLogger,
	)
	ctx, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()

	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *AESEncryptionTestSuite) TestEncryption() {
	text := "The quick brown fox jumps over the lazy dog"

	cipherText, err := s.aesEncryption.Encrypt(text)
	s.Require().Nil(err)

	plainText, err := s.aesEncryption.Decrypt(cipherText)
	s.Require().Nil(err)

	s.Require().Equal(text, plainText)
}
