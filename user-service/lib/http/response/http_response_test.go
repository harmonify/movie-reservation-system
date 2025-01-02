package response_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	error_constant "github.com/harmonify/movie-reservation-system/user-service/lib/error/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/response"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/mocks"
	test_interface "github.com/harmonify/movie-reservation-system/user-service/lib/test/interface"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	struct_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/struct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"

	"errors"
)

func TestResponse(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(ResponseTestSuite))
}

type ResponseTestSuite struct {
	suite.Suite
	logger     logger.Logger
	tracer     tracer.Tracer
	structUtil struct_util.StructUtil
	response   response.HttpResponse
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
	s.logger = mocks.NewLogger(s.T())
	s.tracer = mocks.NewTracer(s.T())
	s.structUtil = mocks.NewStructUtil(s.T())
	s.response = response.NewHttpResponse(s.logger, s.tracer, s.structUtil, &error_constant.DefaultCustomErrorMap)
}

func (s *ResponseTestSuite) TestBuild() {
	testCases := []test_interface.TestCase[testConfig, testExpectation]{
		{
			Description: "Should build success response",
			Config: testConfig{
				HttpCode: 200,
				Data:     "Test data",
			},
			Expectation: testExpectation{
				Success: true,
				Result:  "Test data",
				Error:   nil,
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
				HttpCode: 400,
				Success:  false,
				Result:   "Test data",
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			httpCode, response, responseError := s.response.Build(context.Background(), testCase.Config.HttpCode, testCase.Config.Data, testCase.Config.Error)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}

			s.Assert().True(testCase.Expectation.Success, response.Success)
			s.Assert().Equal(200, httpCode)
			s.Assert().Equal(testCase.Expectation.Result, response.Result)
			s.Assert().Equal(testCase.Expectation.Error, responseError)
		})
	}
}

func (s *ResponseTestSuite) TestBuildWithResponseCode(t *testing.T) {
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

func (s *ResponseTestSuite) TestSend(t *testing.T) {
	t.Run("Should Return Http Success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create a new request.
		req, _ := http.NewRequest("GET", "/", nil)

		// Create a new tracer.
		tracer := otel.Tracer("test-tracer")

		// Start a new span.
		ctx, _ := tracer.Start(context.Background(), "test-span")

		// Add the span to the request context.
		req = req.WithContext(ctx)

		// Add the request to the gin context.
		c.Request = req

		// Act
		s.response.Send(c, nil, nil)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")
	})

	t.Run("Should Return Http Error", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Create a new request.
		req, _ := http.NewRequest("GET", "/", nil)

		// Create a new tracer.
		tracer := otel.Tracer("test-tracer")

		// Start a new span.
		ctx, _ := tracer.Start(context.Background(), "test-span")

		// Add the span to the request context.
		req = req.WithContext(ctx)

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

func (s *ResponseTestSuite) TestBuildError(t *testing.T) {
	t.Run("Should Build Error Correctly", func(t *testing.T) {
		err := errors.New("test error")
		handler := s.response.BuildError("test_code", err)

		if handler.Code != "test_code" {
			t.Errorf("Expected code 'test_code', got '%s'", handler.Code)
		}
	})
}

func (s *ResponseTestSuite) TestBuildErrorValidation(t *testing.T) {
	t.Run("Should build error correctly", func(t *testing.T) {
		err := errors.New("test error")
		handler := s.response.BuildValidationError("test_code", err, []response.BaseValidationErrorSchema{
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
