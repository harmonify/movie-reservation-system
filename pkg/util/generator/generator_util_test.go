package generator_util_test

import (
	"os"
	"testing"

	generator_util "github.com/harmonify/movie-reservation-system/pkg/util/generator"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

func TestGeneratorUtil(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(GeneratorUtilTestSuite))
}

type GeneratorUtilTestSuite struct {
	suite.Suite
	app           any
	generatorUtil generator_util.GeneratorUtil
}

func (s *GeneratorUtilTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.Provide(generator_util.NewGeneratorUtil),
		fx.Invoke(func(generatorUtil generator_util.GeneratorUtil) {
			s.generatorUtil = generatorUtil
		}),

		fx.NopLogger,
	)
}

func (s *GeneratorUtilTestSuite) TestGenerateRandomBytes() {
	bt, err := s.generatorUtil.GenerateRandomBytes(32)

	s.Require().Nil(err)
	s.Require().Len(bt, 32, "Should return bytes with length of 32")
}

func (s *GeneratorUtilTestSuite) TestGenerateRandomHex() {
	hx, err := s.generatorUtil.GenerateRandomHex(32)

	s.Require().Nil(err)
	s.Require().Len(hx, 64, "Should return hex string with length of 64")
}

func (s *GeneratorUtilTestSuite) TestGenerateRandomBase64() {
	b64, err := s.generatorUtil.GenerateRandomBase64(32)

	s.Require().Nil(err)
	s.Require().Len(b64, 43, "Should return base-64 string with length of 43")
}

func (s *GeneratorUtilTestSuite) TestGenerateRandomNumber() {
	number, err := s.generatorUtil.GenerateRandomNumber(6)

	s.Require().Nil(err)
	s.Require().Len(number, 6, "Should return number with length of 6")
}
