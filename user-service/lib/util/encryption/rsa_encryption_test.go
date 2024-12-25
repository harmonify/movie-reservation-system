package encryption_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/test"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util/encryption"
	"github.com/stretchr/testify/suite"
)

func TestRSAEncryption(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(RSAEncryptionTestSuite))
}

type RSAEncryptionTestSuite struct {
	suite.Suite
	app           any
	cfg           *config.Config
	rsaEncryption encryption.RSAEncryption
}

func (s *RSAEncryptionTestSuite) SetupSuite() {
	s.app = test.NewTestApp(s.invoker, s.mock()...)
}

func (t *RSAEncryptionTestSuite) invoker(
	cfg *config.Config,
	rsaEncryption encryption.RSAEncryption,
) {
	t.cfg = &config.Config{
		AppName:   "RSA Encryption Tester",
		AppSecret: "1234567891123456",
	}
	t.rsaEncryption = rsaEncryption
}

func (s *RSAEncryptionTestSuite) mock() []any {
	// s.mockExample = mocks.NewExample(s.T())
	return []any{
		// func() interfaces.Example { return s.mockExample },
	}
}

func (s *RSAEncryptionTestSuite) TestGenerate() {
	keyPair, err := s.rsaEncryption.Generate()

	s.Require().Nil(err)
	fmt.Println("PUBLIC:", string(keyPair.PublicKey))
	fmt.Println("PRIVATE:", string(keyPair.PrivateKey))
}
