package encryption_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/util/encryption"
	generator_util "github.com/harmonify/movie-reservation-system/pkg/util/generator"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

func TestRSAEncryption(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(RSAEncryptionTestSuite))
}

type RSAEncryptionTestSuite struct {
	suite.Suite
	app           *fx.App
	cfg           *config.Config
	rsaEncryption encryption.RSAEncryption
}

func (s *RSAEncryptionTestSuite) SetupSuite() {
	s.app = fx.New(
		generator_util.GeneratorUtilModule,
		fx.Provide(
			func() *config.Config {
				return &config.Config{
					AppSecret: "1234567891123456",
				}
			},
			func() *encryption.AesGcmPbkdf2EncryptionConfig {
				return &encryption.AesGcmPbkdf2EncryptionConfig{
					PBKDF2Iterations: int(15000),
				}
			},
			encryption.NewAESEncryption,
			encryption.NewRSAEncryption,
		),
		fx.Invoke(func(rsaEncryption encryption.RSAEncryption) {
			s.rsaEncryption = rsaEncryption
		}),

		fx.NopLogger,
	)
	ctx, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()

	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *RSAEncryptionTestSuite) TestGenerate() {
	keyPair, err := s.rsaEncryption.Generate()

	s.Require().Nil(err)
	fmt.Println("PUBLIC:", string(keyPair.PublicKey))
	fmt.Println("PRIVATE:", string(keyPair.PrivateKey))
}
