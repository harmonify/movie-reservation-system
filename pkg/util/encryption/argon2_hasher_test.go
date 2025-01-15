package encryption_test

import (
	"context"
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/util/encryption"
	generator_util "github.com/harmonify/movie-reservation-system/pkg/util/generator"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

func TestArgon2Hasher(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(Argon2HasherTestSuite))
}

type Argon2HasherTestSuite struct {
	suite.Suite
	app          *fx.App
	argon2Hasher encryption.Argon2Hasher
}

func (s *Argon2HasherTestSuite) SetupSuite() {
	s.app = fx.New(
		generator_util.GeneratorUtilModule,
		fx.Provide(
			func() *config.Config {
				return &config.Config{
					AppSecret: "1234567891123456",
				}
			},
			func() encryption.Argon2HasherConfig {
				return *encryption.Argon2HasherDefaultConfig
			},
			encryption.NewArgon2Hasher,
		),
		fx.Invoke(func(argon2Hasher encryption.Argon2Hasher) {
			s.argon2Hasher = argon2Hasher
		}),

		fx.NopLogger,
	)
	ctx, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()

	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
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
