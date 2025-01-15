package response_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	error_constant "github.com/harmonify/movie-reservation-system/pkg/error/constant"
	"github.com/harmonify/movie-reservation-system/pkg/http/response"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	struct_util "github.com/harmonify/movie-reservation-system/pkg/util/struct"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"

	"errors"
)

func TestHttpResponse(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(ResponseTestSuite))
}

type ResponseTestSuite struct {
	suite.Suite
	app      *fx.App
	response response.HttpResponse
}

type testConfig struct {
	HttpCode int
	Data     string
	Error    error
}

type testExpectation struct {
	Success  bool
	HttpCode int
	Result   string
	Error    error
}

func (s *ResponseTestSuite) SetupSuite() {
	s.app = fx.New(
		logger.LoggerModule,
		tracer.TracerModule,
		fx.Provide(
			func() *config.Config {
				return &config.Config{
					Env:      "test",
					LogType:  "nop",
					LogLevel: "debug",
				}
			},
			func() *error_constant.CustomErrorMap {
				return &error_constant.DefaultCustomErrorMap
			},
			struct_util.NewStructUtil,
			response.NewHttpResponse,
		),
		fx.Invoke(func(response response.HttpResponse) {
			s.response = response
		}),

		fx.NopLogger,
	)
}

func (s *ResponseTestSuite) TestHttpResponse_Build() {
	testCases := []test_interface.TestCase[testConfig, testExpectation]{
		{
			Description: "Should build success response",
			Config: testConfig{
				HttpCode: 200,
				Data:     "Test data",
				Error:    nil,
			},
			Expectation: testExpectation{
				Success:  true,
				HttpCode: 200,
				Result:   "Test data",
				Error:    nil,
			},
		},
		{
			Description: "Should build error response",
			Config: testConfig{
				HttpCode: 400,
				Data:     "Test data",
				Error:    errors.New("Test error"),
			},
			Expectation: testExpectation{
				Success:  false,
				HttpCode: 500, // unknown error http code should be overwritten to 500
				Result:   "Test data",
				Error:    errors.New("Test error"),
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			httpCode, httpResponse, httpResponseError := s.response.Build(
				context.Background(),
				testCase.Config.HttpCode,
				testCase.Config.Data,
				testCase.Config.Error,
			)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}

			s.Assert().Equal(testCase.Expectation.Success, httpResponse.Success)
			s.Assert().Equal(testCase.Expectation.HttpCode, httpCode)
			s.Assert().Equal(testCase.Expectation.Result, httpResponse.Result)
			if testCase.Expectation.Error != nil {
				s.Require().IsType(httpResponseError, &response.HttpErrorHandlerImpl{})
				s.Assert().Equal(
					httpResponseError.(*response.HttpErrorHandlerImpl).Code,
					testCase.Expectation.Error.Error(),
				)
			}
		})
	}
}

func (s *ResponseTestSuite) TestHttpResponse_Build_ResponseCode() {
	testCases := []test_interface.TestCase[testConfig, testExpectation]{}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			httpCode, response, responseError := s.response.Build(context.Background(), testCase.Config.HttpCode, testCase.Config.Data, nil)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}

			s.Assert().True(testCase.Expectation.Success, response.Success)
			s.Assert().Equal(testCase.Expectation.HttpCode, httpCode)
			s.Assert().Equal(testCase.Expectation.Result, response.Result)
			s.Assert().Equal(testCase.Expectation.Error, responseError)
		})
	}
}

func (s *ResponseTestSuite) TestHttpResponse_Send() {
	s.T().Run("Should Return Http Success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create a new request.
		req, _ := http.NewRequest("GET", "/", nil)

		// Add the request to the gin context.
		c.Request = req

		// Act
		s.response.Send(c, nil, nil)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	s.T().Run("Should Return Http Error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create a new request.
		req, _ := http.NewRequest("GET", "/", nil)

		// Add the request to the gin context.
		c.Request = req

		// Act
		err := s.response.BuildError(error_constant.InternalServerError, errors.New("test error"))
		s.response.Send(c, nil, err)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})
}

func (s *ResponseTestSuite) TestHttpResponse_BuildError() {
	s.T().Run("Should Build Error Correctly", func(t *testing.T) {
		err := errors.New("test error")
		handler := s.response.BuildError("test_code", err)

		if handler.Code != "test_code" {
			t.Errorf("Expected code 'test_code', got '%s'", handler.Code)
		}
	})
}

func (s *ResponseTestSuite) TestHttpResponse_BuildValidationError() {
	s.T().Run("Should build error correctly", func(t *testing.T) {
		err := errors.New("test error")
		handler := s.response.BuildValidationError("test_code", err, []validation.BaseValidationErrorSchema{
			{
				Field:   "error",
				Message: "test error",
			},
		})

		if handler.Code != "test_code" {
			t.Errorf("Expected code 'test_code', got '%s'", handler.Code)
		}
	})
}
