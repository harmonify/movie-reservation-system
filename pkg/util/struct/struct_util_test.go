package struct_util_test

import (
	"context"
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	struct_util "github.com/harmonify/movie-reservation-system/pkg/util/struct"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

func TestStructUtilSuite(t *testing.T) {
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
	Data interface{}
}

type testExpectation struct {
	Result interface{}
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

func (s *StructUtilTestSuite) TestStructUtilSuite_SetNonPrimitiveDefaultValue() {
	var data interface{}
	nonEmptyPrimitive := "hello"
	emptyPrimitive := ""

	testCases := []test_interface.TestCase[testConfig, testExpectation]{
		{
			Description: "Should handle nil",
			Config:      testConfig{Data: data},
			Expectation: testExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle nil pointer",
			Config:      testConfig{Data: &data},
			Expectation: testExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle non-empty struct",
			Config:      testConfig{Data: struct{ Hello string }{Hello: "world"}},
			Expectation: testExpectation{Result: struct{ Hello string }{Hello: "world"}},
		},
		{
			Description: "Should handle non-empty struct pointer",
			Config:      testConfig{Data: &struct{ Hello string }{Hello: "world"}},
			Expectation: testExpectation{Result: &struct{ Hello string }{Hello: "world"}},
		},
		{
			Description: "Should handle empty struct",
			Config:      testConfig{Data: struct{ Name string }{}},
			Expectation: testExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle empty struct pointer",
			Config:      testConfig{Data: &struct{ Name string }{}},
			Expectation: testExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle non-empty slice",
			Config:      testConfig{Data: []string{"hello"}},
			Expectation: testExpectation{Result: []string{"hello"}},
		},
		{
			Description: "Should handle non-empty slice pointer",
			Config:      testConfig{Data: &[]string{"hello"}},
			Expectation: testExpectation{Result: &[]string{"hello"}},
		},
		{
			Description: "Should handle non-empty slice of pointers",
			Config:      testConfig{Data: []*string{&nonEmptyPrimitive}},
			Expectation: testExpectation{Result: []*string{&nonEmptyPrimitive}},
		},
		{
			Description: "Should handle empty slice",
			Config:      testConfig{Data: []string{}},
			Expectation: testExpectation{Result: []string{}},
		},
		{
			Description: "Should handle empty slice pointer",
			Config:      testConfig{Data: &[]string{}},
			Expectation: testExpectation{Result: &[]string{}},
		},
		{
			Description: "Should handle non-empty array",
			Config:      testConfig{Data: [1]string{"hello"}},
			Expectation: testExpectation{Result: [1]string{"hello"}},
		},
		{
			Description: "Should handle empty array",
			Config:      testConfig{Data: [1]string{}},
			Expectation: testExpectation{Result: []interface{}{}},
		},
		{
			Description: "Should handle non-empty map",
			Config:      testConfig{Data: map[string]string{"hello": "world"}},
			Expectation: testExpectation{Result: map[string]string{"hello": "world"}},
		},
		{
			Description: "Should handle non-empty map pointer",
			Config:      testConfig{Data: &map[string]string{"hello": "world"}},
			Expectation: testExpectation{Result: &map[string]string{"hello": "world"}},
		},
		{
			Description: "Should handle empty map",
			Config:      testConfig{Data: map[string]string{}},
			Expectation: testExpectation{Result: map[string]string{}},
		},
		{
			Description: "Should handle empty map pointer",
			Config:      testConfig{Data: &map[string]string{}},
			Expectation: testExpectation{Result: &map[string]string{}},
		},
		{
			Description: "Should handle empty interface",
			Config:      testConfig{Data: interface{}(nil)},
			Expectation: testExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle non-empty interface",
			Config:      testConfig{Data: interface{}(struct{ Hello string }{Hello: "world"})},
			Expectation: testExpectation{Result: struct{ Hello string }{Hello: "world"}},
		},
		{
			Description: "Should handle non-empty interface pointer",
			Config:      testConfig{Data: interface{}(&struct{ Hello string }{Hello: "world"})},
			Expectation: testExpectation{Result: &struct{ Hello string }{Hello: "world"}},
		},
		{
			Description: "Should handle empty interface pointer",
			Config:      testConfig{Data: interface{}(&struct{ Hello string }{})},
			Expectation: testExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle non-empty primitives",
			Config:      testConfig{Data: "hello"},
			Expectation: testExpectation{Result: "hello"},
		},
		{
			Description: "Should handle non-empty primitives pointer",
			Config:      testConfig{Data: &nonEmptyPrimitive},
			Expectation: testExpectation{Result: &nonEmptyPrimitive},
		},
		{
			Description: "Should handle empty primitives",
			Config:      testConfig{Data: ""},
			Expectation: testExpectation{Result: nil},
		},
		{
			Description: "Should handle empty primitives pointer",
			Config:      testConfig{Data: &emptyPrimitive},
			Expectation: testExpectation{Result: nil},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()

		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			result := s.structUtil.SetNonPrimitiveDefaultValue(ctx, testCase.Config.Data)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}

			s.Assert().Equal(testCase.Expectation.Result, result)
		})
	}
}
