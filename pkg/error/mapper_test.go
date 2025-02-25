package error_pkg_test

import (
	"errors"
	"os"
	"testing"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

func TestHttpResponse(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(ResponseTestSuite))
}

type ResponseTestSuite struct {
	suite.Suite
	app         *fx.App
	errorMapper error_pkg.ErrorMapper
}

type testConfig struct {
	SuccessHttpCode int
	Data            string
	Error           *error_pkg.ErrorWithDetails
}

type testExpectation struct {
	Error *error_pkg.ErrorWithDetails
	Valid bool
}

func (s *ResponseTestSuite) SetupSuite() {
	s.app = fx.New(
		error_pkg.ErrorModule,
		fx.Invoke(func(errorMapper error_pkg.ErrorMapper) {
			s.errorMapper = errorMapper
		}),

		fx.NopLogger,
	)
}

func (s *ResponseTestSuite) TestErrorMapper_FromError() {
	tests := []struct {
		name        string
		errFactory  func() error
		expectation testExpectation
	}{
		{
			name: "nil",
			errFactory: func() error {
				return nil
			},
			expectation: testExpectation{
				Error: nil,
				Valid: true,
			},
		},
		{
			name: "unknown error",
			errFactory: func() error {
				return errors.New("unknown error")
			},
			expectation: testExpectation{
				Error: error_pkg.DefaultError,
				Valid: false,
			},
		},
		{
			name: "valid error",
			errFactory: func() error {
				return error_pkg.InternalServerError
			},
			expectation: testExpectation{
				Error: error_pkg.InternalServerError,
				Valid: true,
			},
		},
		{
			name: "error with stack",
			errFactory: func() error {
				err := errors.New("db crashed")
				return error_pkg.NewErrorWithStack(err, 1)
			},
			expectation: testExpectation{
				Error: error_pkg.ServiceUnavailableError,
				Valid: true,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err, valid := s.errorMapper.FromError(tt.errFactory())
			assert.Equal(s.T(), tt.expectation.Error, err)
			assert.Equal(s.T(), tt.expectation.Valid, valid)
		})
	}
}
