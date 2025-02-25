package struct_util_test

import (
	"context"
	"os"
	"testing"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	test_proto "github.com/harmonify/movie-reservation-system/pkg/test/proto"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	struct_util "github.com/harmonify/movie-reservation-system/pkg/util/struct"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
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

type setNonPrimitiveDefaultValueTestConfig struct {
	Data interface{}
}

type setNonPrimitiveDefaultValueTestExpectation struct {
	Result interface{}
}

type convertProtoToMapTestConfig struct {
	Data proto.Message
}

type convertProtoToMapTestExpectation struct {
	Result map[string]interface{}
}

func (s *StructUtilTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.Provide(
			func(lc fx.Lifecycle) (tracer.Tracer, error) {
				return tracer.NewTracer(lc, &tracer.TracerConfig{
					Env:               "test",
					ServiceIdentifier: "test",
					Type:              "console",
				})
			},
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

	if err := s.app.Start(context.Background()); err != nil {
		s.T().Fatal(err)
	}
}

func (s *StructUtilTestSuite) TestStructUtilSuite_SetOrDefault() {
	var data interface{}
	nonEmptyPrimitive := "hello"
	emptyPrimitive := ""

	testCases := []test_interface.TestCase[setNonPrimitiveDefaultValueTestConfig, setNonPrimitiveDefaultValueTestExpectation]{
		{
			Description: "Should handle nil",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: data},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle nil pointer",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: &data},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle non-empty struct",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: struct{ Hello string }{Hello: "world"}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: struct{ Hello string }{Hello: "world"}},
		},
		{
			Description: "Should handle non-empty struct pointer",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: &struct{ Hello string }{Hello: "world"}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: &struct{ Hello string }{Hello: "world"}},
		},
		{
			Description: "Should handle empty struct",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: struct{ Name string }{}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle empty struct pointer",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: &struct{ Name string }{}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle non-empty slice",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: []string{"hello"}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: []string{"hello"}},
		},
		{
			Description: "Should handle non-empty slice pointer",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: &[]string{"hello"}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: &[]string{"hello"}},
		},
		{
			Description: "Should handle non-empty slice of pointers",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: []*string{&nonEmptyPrimitive}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: []*string{&nonEmptyPrimitive}},
		},
		{
			Description: "Should handle empty slice",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: []string{}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: []string{}},
		},
		{
			Description: "Should handle empty slice pointer",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: &[]string{}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: &[]string{}},
		},
		{
			Description: "Should handle non-empty array",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: [1]string{"hello"}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: [1]string{"hello"}},
		},
		{
			Description: "Should handle empty array",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: [1]string{}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: []interface{}{}},
		},
		{
			Description: "Should handle non-empty map",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: map[string]string{"hello": "world"}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: map[string]string{"hello": "world"}},
		},
		{
			Description: "Should handle non-empty map pointer",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: &map[string]string{"hello": "world"}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: &map[string]string{"hello": "world"}},
		},
		{
			Description: "Should handle empty map",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: map[string]string{}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: map[string]string{}},
		},
		{
			Description: "Should handle empty map pointer",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: &map[string]string{}},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: &map[string]string{}},
		},
		{
			Description: "Should handle empty interface",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: interface{}(nil)},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle non-empty interface",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: interface{}(struct{ Hello string }{Hello: "world"})},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: struct{ Hello string }{Hello: "world"}},
		},
		{
			Description: "Should handle non-empty interface pointer",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: interface{}(&struct{ Hello string }{Hello: "world"})},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: &struct{ Hello string }{Hello: "world"}},
		},
		{
			Description: "Should handle empty interface pointer",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: interface{}(&struct{ Hello string }{})},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle non-empty primitives",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: "hello"},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: "hello"},
		},
		{
			Description: "Should handle non-empty primitives pointer",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: &nonEmptyPrimitive},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: &nonEmptyPrimitive},
		},
		{
			Description: "Should handle empty primitives",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: ""},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: nil},
		},
		{
			Description: "Should handle empty primitives pointer",
			Config:      setNonPrimitiveDefaultValueTestConfig{Data: &emptyPrimitive},
			Expectation: setNonPrimitiveDefaultValueTestExpectation{Result: nil},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()

		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			result := s.structUtil.SetOrDefault(ctx, testCase.Config.Data)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}

			s.Assert().Equal(testCase.Expectation.Result, result)
		})
	}
}

func (s *StructUtilTestSuite) TestStructUtilSuite_ConvertProtoToMap() {
	var data proto.Message
	var emptyProto = test_proto.Test{}
	var nonEmptyProto = test_proto.Test{Message: "Hello world"}
	var complexProto = test_proto.ComplexTestMessage{
		Message:         "Hello message",
		RepeatedMessage: []string{"Hello", "repeated", "message"},
		NestedTestMessage: &test_proto.NestedTestMessage{
			NestedMessage: "Hello nested message",
		},
		RepeatedNestedTestMessage: []*test_proto.NestedTestMessage{
			{NestedMessage: "Hello repeated nested message 1"},
			{NestedMessage: "Hello repeated nested message 2"},
		},
	}

	testCases := []test_interface.TestCase[convertProtoToMapTestConfig, convertProtoToMapTestExpectation]{
		{
			Description: "Should handle nil",
			Config:      convertProtoToMapTestConfig{Data: data},
			Expectation: convertProtoToMapTestExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle empty proto",
			Config:      convertProtoToMapTestConfig{Data: &emptyProto},
			Expectation: convertProtoToMapTestExpectation{Result: map[string]interface{}{}},
		},
		{
			Description: "Should handle non-empty proto",
			Config:      convertProtoToMapTestConfig{Data: &nonEmptyProto},
			Expectation: convertProtoToMapTestExpectation{Result: map[string]interface{}{"message": "Hello world"}},
		},
		{
			Description: "Should handle complex proto",
			Config:      convertProtoToMapTestConfig{Data: &complexProto},
			Expectation: convertProtoToMapTestExpectation{Result: map[string]interface{}{
				"message":         "Hello message",
				"repeatedMessage": []interface{}{"Hello", "repeated", "message"},
				"nestedTestMessage": map[string]interface{}{
					"nestedMessage": "Hello nested message",
				},
				"repeatedNestedTestMessage": []interface{}{
					map[string]interface{}{"nestedMessage": "Hello repeated nested message 1"},
					map[string]interface{}{"nestedMessage": "Hello repeated nested message 2"},
				},
			}},
		},
	}

	for _, testCase := range testCases {
		ctx := context.Background()

		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			result, err := s.structUtil.ConvertProtoToMap(ctx, testCase.Config.Data)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}

			s.Assert().Nil(err)
			s.Assert().Equal(testCase.Expectation.Result, result)
		})
	}
}
