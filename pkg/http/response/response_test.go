package response_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/harmonify/movie-reservation-system/pkg/constant"
	"github.com/harmonify/movie-reservation-system/pkg/http/response"
	logger_shared "github.com/harmonify/movie-reservation-system/pkg/logger/shared"
	"github.com/harmonify/movie-reservation-system/pkg/mocks"
	"github.com/harmonify/movie-reservation-system/pkg/test"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	struct_util "github.com/harmonify/movie-reservation-system/pkg/util/struct"
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
	logger     logger_shared.Logger
	tracer     tracer.Tracer
	structUtil struct_util.StructUtil
	response   response.Response
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
	Error    *response.ErrorHandler
}

func (s *ResponseTestSuite) SetupSuite() {
	s.logger = mocks.NewLogger(s.T())
	s.tracer = mocks.NewTracer(s.T())
	s.structUtil = mocks.NewStructUtil(s.T())
	s.response = response.NewResponse(s.logger, s.tracer, s.structUtil, &constant.DefaultCustomHttpErrorMap)
}

func (s *ResponseTestSuite) TestBuild() {
	testCases := []test.TestCase[testConfig, testExpectation]{
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
		config := testCase.Config.(testConfig)

		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(config)
			}

			httpCode, response, responseError := s.response.Build(context.Background(), config.HttpCode, config.Data, config.Error)

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
	testCases := []test.TestCase[testConfig, testExpectation]{}

	for _, testCase := range testCases {
		config := testCase.Config.(testConfig)

		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(config)
			}

			httpCode, response, responseError := s.response.Build(context.Background(), config.HttpCode, config.Data, nil)

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
		httpCode, response, responseError := s.response.Build(ctx, http.StatusOK, nil, nil)
		s.response.Send(c, httpCode, response, responseError)

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
		httpCode, response, responseError := s.response.Build(ctx, http.StatusInternalServerError, nil, errors.New("test error"))
		s.response.Send(c, httpCode, response, responseError)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})
}

func (s *ResponseTestSuite) TestBuildError(t *testing.T) {
	t.Run("Should Build Error Correctly", func(t *testing.T) {
		err := errors.New("test error")
		handler := s.response.BuildError("test_code", err)

		if handlerImpl, ok := handler.(*response.ErrorHandlerImpl); ok {
			if handlerImpl.Code != "test_code" {
				t.Errorf("Expected code 'test_code', got '%s'", handlerImpl.Code)
			}
		} else {
			t.Errorf("Expected ErrorHandlerImpl, got %T", handler)
		}
	})
}

func (s *ResponseTestSuite) TestBuildErrorValidation(t *testing.T) {
	t.Run("Should Validate Buildn Error Correctly", func(t *testing.T) {
		err := errors.New("test error")
		errorFields := map[string]string{"field": "error"}
		handler := s.response.BuildValidationError("test_code", err, errorFields)

		if handlerImpl, ok := handler.(*response.ErrorHandlerImpl); ok {
			if handlerImpl.Code != "test_code" {
				t.Errorf("Expected code 'test_code', got '%s'", handlerImpl.Code)
			}
		} else {
			t.Errorf("Expected ErrorHandlerImpl, got %T", handler)
		}
	})
}
