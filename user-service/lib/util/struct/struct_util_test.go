package struct_util_test

import (
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/user-service/internal"
	test_interface "github.com/harmonify/movie-reservation-system/user-service/lib/test/interface"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util/struct"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

func TestStructUtil(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(StructUtilTestSuite))
}

type StructUtilTestSuite struct {
	suite.Suite
	app        *fx.App
	structUtil struct_util.StructUtil
}

type testConfig struct {
	Data any
}

type testExpectation struct {
	Result any
}

func (s *StructUtilTestSuite) SetupSuite() {
	s.app = internal.NewApp(
		fx.Invoke(func(structUtil struct_util.StructUtil) {
			s.structUtil = structUtil
		}),
		// fx.Decorate(func() interfaces.ExampleService {
		// s.exampleService = mocks.NewExampleService(s.T())
		// 	return s.exampleService
		// }),
	)
}

func (s *StructUtilTestSuite) TestSetValueIfNotEmpty() {
	var data interface{}

	testCases := []test_interface.TestCase[testConfig, testExpectation]{
		{
			Description: "Should send correct response",
			Config: testConfig{
				Data: data,
			},
			Expectation: testExpectation{
				Result: "Test data",
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			result := s.structUtil.SetValueIfNotEmpty(testCase.Config.Data)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}

			s.Assert().Equal(testCase.Expectation.Result, result)
		})
	}
}
