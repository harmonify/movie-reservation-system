package jwt_util_test

import (
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/test"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	"github.com/stretchr/testify/suite"
)

func TestJWTUtil(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(JWTUtilTestSuite))
}

type JWTUtilTestSuite struct {
	suite.Suite
	app     any
	cfg     *config.Config
	jwtUtil jwt_util.JWTUtil
}

func (s *JWTUtilTestSuite) SetupSuite() {
	// TODO: fix
	// &config.Config{
	// 	AppSecret: "1234567891123456",
	// }
	s.app = test.NewTestApp(s.invoker, s.mock()...)
}

func (t *JWTUtilTestSuite) invoker(
	cfg *config.Config,
) {
	t.cfg = cfg
}

func (s *JWTUtilTestSuite) mock() []any {
	// s.mockExample = mocks.NewExample(s.T())
	return []any{
		// func() interfaces.Example { return s.mockExample },
	}
}

func (s *JWTUtilTestSuite) TestJWTSign() {
	s.T().Run("Should return error when signing JWT", func(t *testing.T) {
		// Define a JWTPayload
		payload := &jwt_util.JWTPayload{
			SecretKey:    "1234567891123456",
			ExpInMinutes: 60,
			PrivateKey:   "privateKey",
			PublicKey:    "publicKey",
			BodyPayload: jwt_util.JWTBodyPayload{
				Email:  "test@example.com",
				UserID: "testUserID",
			},
		}

		// Call the JWTSign function
		_, err := s.jwtUtil.JWTSign(payload)

		// Assert that there is error
		s.Assert().Error(err)
	})
}
