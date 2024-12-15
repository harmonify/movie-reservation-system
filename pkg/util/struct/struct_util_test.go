package struct_util_test

import (
	"os"
	"testing"

	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	"github.com/harmonify/movie-reservation-system/pkg/util/struct"
	"github.com/stretchr/testify/suite"
)

func TestStructUtil(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(StructUtilTestSuite))
}

type StructUtilTestSuite struct {
	suite.Suite
	structUtil struct_util.StructUtil
}

type testConfig struct {
	Data any
}

type testExpectation struct {
	Result any
}

func (s *StructUtilTestSuite) SetupSuite() {
	s.structUtil = struct_util.NewStructUtil()
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
		config := testCase.Config.(testConfig)

		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(config)
			}

			result := s.structUtil.SetValueIfNotEmpty(config.Data)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}

			s.Assert().Equal(testCase.Expectation.Result, result)
		})
	}
}
