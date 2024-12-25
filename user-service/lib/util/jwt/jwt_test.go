package jwt_util_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/test"
	test_interface "github.com/harmonify/movie-reservation-system/user-service/lib/test/interface"
	jwt_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/jwt"
	"github.com/stretchr/testify/suite"
)

var (
	TEST_PRIVATE_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAyotBm+v0jUBCzvh9gA356KhbrIfWiVsPqxxGh90sZchN6XXv
ZDA3OaPdFwPf8gMEig8utDb/c0T6VP84IxvBJamLLlN2Hx78zdnvP4ZJz3Xj868f
Rj5PPO4kD6rRsr3kDRUedh9aNajiJJgpG7cm6MD3AT8dWN6S+hJ3OqTMrbGDJh40
qiNSM8u9gHArUw1X6aNMGqU1QJeiYR9GDhr/8XY5TT5y3OAci5MSZsMIwSxRWkdP
M0bunaI757NkCdjL+OnIxP0xUW/g+BoyAcmlfwCBcOesIW72hn49PhNXfgBBcR5W
ZlGRs9p4IL/SMg4wTe+lzvn19aM7s3W1MJnUNwIDAQABAoIBAQCErr0IG4ZUkegy
FW6BWKaB1uhXGZVc3Z5iBV/e2PCgrJr9eRidlUhYJhRLY2ps67Upi9CYlf650FH9
JEPuG9xng619Z4dV08LgRwgHoTzw/tWZaPsf1OmrjIVrDgfZA7RFLbSKxPcfd8bN
GjCzy0Nd7irhUiszcHrv/vDEJfk/PoK+LuUGKbgcHuHPgaBlpAgKx7iCoA+HD3r6
bnGZbGDlY9uoEteh56RlaQtGNVTLwbCgUVHGUKhOLuOtApU/vKGltbyhXMihFMSJ
hhhWDLCXf4/7kAY/R7i4SeY0tQmAj8F0LRtZqU14YkNSzfqMbZ77l0ui9wFDY5bW
oeZJR3KBAoGBANF3Ja6NI2wLkQOyIHhz4BuH6VtLUN4hzcrgrDEENomx4dOoVeZh
oIrD21e8vD0gn97f3QN3IbTwbxhrlzE55Kc0EkTR0fvp1X8jUrKoHsTcnKlkdS9d
E+X4yEv8o713vOzEi1vZ9c2rDkaNK7uPteGW3+LGEjpWoB6d4/MewDgxAoGBAPeK
d80W7MfCV0/8UxYTRCW0tSyv/CXZE2+XrC7StjMehEp6PyfrF9e7F3EcLVJh8m9F
JJCe6aDRb2yvqfuzfMgXWTT1YOkjyM/MLU8lKfGFK7lbUxyTzkrGgoDZ2kbqa6z/
TOGni6QokdSSi4qcSR4TBSDs4zk7ztdgXSLE8SDnAoGAdX/rsGXV1/cJCtSyKD+A
GJF+EstF+sVlpoevr/NYEJerQUrtnMVpBE5nzWi/A184rxJO7XG3g8NX3pAECQYb
wLuR/+7fZvu92orbCgMK9411iAQlRENnNRsAaLe4tkDjxsFeF1FF9HAfGu53+Mfd
1EUJJDHN6dHMEkCprSiz1RECgYAnhwvkSvHaYBUTJ6aY0tDB+J4pmZx46rXozt5m
x6zictALGIQ3OpofD7gJjsdJ7WwKCo9xLH7/+BGD2HUbRSj6xoevJjOoZtdtHxxp
E/UjpPE7cvLNkGiTlilGrALn6gzxnf7H1bo9p2DKAfCYXKZsT/s0q78I55z61V9p
6uraJwKBgFvyqH1dBDWxywz4pvomSYxR4DtUVF4jI4yj2JuSR08R5AP/c5h8FEq5
rp5kVwKpkd0G30GB2nIr0P76ORJjk8/+f+S4aXV+aOeymOJb5Srwx9+U2F2qWmBy
p9JV7cBOiF76ml+L8Y8uPurqORvBDU0ZGNvcPrFr20JQ1U6RUIsx
-----END RSA PRIVATE KEY-----`

	TEST_PUBLIC_KEY = `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAyotBm+v0jUBCzvh9gA356KhbrIfWiVsPqxxGh90sZchN6XXvZDA3
OaPdFwPf8gMEig8utDb/c0T6VP84IxvBJamLLlN2Hx78zdnvP4ZJz3Xj868fRj5P
PO4kD6rRsr3kDRUedh9aNajiJJgpG7cm6MD3AT8dWN6S+hJ3OqTMrbGDJh40qiNS
M8u9gHArUw1X6aNMGqU1QJeiYR9GDhr/8XY5TT5y3OAci5MSZsMIwSxRWkdPM0bu
naI757NkCdjL+OnIxP0xUW/g+BoyAcmlfwCBcOesIW72hn49PhNXfgBBcR5WZlGR
s9p4IL/SMg4wTe+lzvn19aM7s3W1MJnUNwIDAQAB
-----END RSA PUBLIC KEY-----`
)

func TestJWTUtil(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(JWTUtilTestSuite))
}

type jwtSignTestConfig struct {
	Data jwt_util.JWTSignParam
}

type jwtVerifyTestConfig struct {
	Data jwt_util.JWTSignParam
}

type jwtVerifyTestExpectation struct {
	Result jwt_util.JWTBodyPayload
}

type JWTUtilTestSuite struct {
	suite.Suite
	app     any
	cfg     *config.Config
	jwtUtil jwt_util.JwtUtil
}

func (s *JWTUtilTestSuite) SetupSuite() {
	s.app = test.NewTestApp(s.invoker, s.mock()...)
}

func (t *JWTUtilTestSuite) invoker(
	cfg *config.Config,
	jwtUtil jwt_util.JwtUtil,
) {
	t.cfg = &config.Config{
		AppName:   "JWT Util Tester",
		AppSecret: "1234567891123456",
	}
	t.jwtUtil = jwtUtil
}

func (s *JWTUtilTestSuite) mock() []any {
	// s.mockExample = mocks.NewExample(s.T())
	return []any{
		// func() interfaces.Example { return s.mockExample },
	}
}

func (s *JWTUtilTestSuite) TestJWTSign() {
	testCases := []test_interface.TestCase[jwtSignTestConfig, any]{
		{
			Description: "Should return no error",
			Config: jwtSignTestConfig{
				Data: jwt_util.JWTSignParam{
					BodyPayload: jwt_util.JWTBodyPayload{
						UUID:        "testUserID",
						Username:    "test_user",
						Email:       "test@example.com",
						PhoneNumber: "081234567890",
					},
					ExpInSeconds: 15,
					PrivateKey:   []byte(TEST_PRIVATE_KEY),
					SecretKey:    s.cfg.AppSecret,
				},
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			_, err := s.jwtUtil.JWTSign(testCase.Config.Data)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}

			s.Assert().Nil(err)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		})
	}
}

func (s *JWTUtilTestSuite) TestJWTVerify() {
	testCases := []test_interface.TestCase[jwtVerifyTestConfig, jwtVerifyTestExpectation]{
		{
			Description: "Should return no error",
			Config: jwtVerifyTestConfig{
				Data: jwt_util.JWTSignParam{
					BodyPayload: jwt_util.JWTBodyPayload{
						UUID:        "testUserID",
						Username:    "test_user",
						Email:       "test@example.com",
						PhoneNumber: "081234567890",
					},
					ExpInSeconds: 15,
					PrivateKey:   []byte(TEST_PRIVATE_KEY),
					SecretKey:    s.cfg.AppSecret,
				},
			},
			Expectation: jwtVerifyTestExpectation{
				Result: jwt_util.JWTBodyPayload{
					UUID:        "testUserID",
					Username:    "test_user",
					Email:       "test@example.com",
					PhoneNumber: "081234567890",
				},
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			token, err := s.jwtUtil.JWTSign(testCase.Config.Data)

			s.Assert().Nil(err)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			payload, err := s.jwtUtil.JWTVerify(token)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}

			s.Assert().Nil(err)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			s.Require().Equal(testCase.Expectation.Result.UUID, payload.UUID)
			s.Require().Equal(testCase.Expectation.Result.Username, payload.Username)
			s.Require().Equal(testCase.Expectation.Result.Email, payload.Email)
			s.Require().Equal(testCase.Expectation.Result.PhoneNumber, payload.PhoneNumber)
		})
	}
}
