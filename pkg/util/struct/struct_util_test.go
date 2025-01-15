package struct_util_test

import (
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	struct_util "github.com/harmonify/movie-reservation-system/pkg/util/struct"
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
	s.app = fx.New(
		fx.Provide(
			logger.NewConsoleLogger,
			struct_util.NewStructUtil,
		),
		fx.Invoke(func(structUtil struct_util.StructUtil) {
			s.structUtil = structUtil
		}),
		// fx.Decorate(func() interfaces.ExampleService {
		// s.exampleService = mocks.NewExampleService(s.T())
		// 	return s.exampleService
		// }),

		fx.NopLogger,
	)
}

func (s *StructUtilTestSuite) TestSetValueIfNotEmpty() {
	var data interface{} = struct{}{}

	testCases := []test_interface.TestCase[testConfig, testExpectation]{
		{
			Description: "Should handle nil-variable",
			Config:      testConfig{Data: data},
			Expectation: testExpectation{Result: map[string]any{}},
		},
		{
			Description: "Should handle nil",
			Config:      testConfig{Data: nil},
			Expectation: testExpectation{Result: struct{}{}},
		},
		{
			Description: "Should handle empty string",
			Config:      testConfig{Data: ""},
			Expectation: testExpectation{Result: []any{}},
		},
		{
			Description: "Should handle empty struct",
			Config:      testConfig{Data: struct{ Name string }{}},
			Expectation: testExpectation{Result: map[string]any{}},
		},
		{
			Description: "Should handle non-empty string",
			Config:      testConfig{Data: "hello"},
			Expectation: testExpectation{Result: "hello"},
		},
		{
			Description: "Should handle non-empty struct",
			Config:      testConfig{Data: struct{ Hello string }{Hello: "world"}},
			Expectation: testExpectation{Result: struct{ Hello string }{Hello: "world"}},
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
