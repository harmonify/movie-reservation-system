package encryption_test

import (
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/test"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util/encryption"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestArgon2Hasher(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(Argon2HasherTestSuite))
}

type Argon2HasherTestSuite struct {
	suite.Suite
	app          any
	cfg          *config.Config
	argon2Hasher encryption.Argon2Hasher
}

func (s *Argon2HasherTestSuite) SetupSuite() {
	s.app = test.NewTestApp(s.invoker, s.mock()...)
}

func (t *Argon2HasherTestSuite) invoker(
	cfg *config.Config,
	argon2Hasher encryption.Argon2Hasher,
) {
	t.cfg = &config.Config{
		AppName:   "RSA Encryption Tester",
		AppSecret: "1234567891123456",
	}
	t.argon2Hasher = argon2Hasher
}

func (s *Argon2HasherTestSuite) mock() []any {
	// s.mockExample = mocks.NewExample(s.T())
	return []any{
		// func() interfaces.Example { return s.mockExample },
	}
}

func (s *Argon2HasherTestSuite) TestArgon2Hasher() {
	var hashed string
	var err error

	s.T().Run("Should return hashed password", func(t *testing.T) {
		hashed, err = s.argon2Hasher.Hash("password")

		require.NoError(t, err)
	})

	s.T().Run("Should return true when password match", func(t *testing.T) {
		match, err := s.argon2Hasher.Compare(hashed, "password")

		require.NoError(t, err)
		require.True(t, match)
	})
}
