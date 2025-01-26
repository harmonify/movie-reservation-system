package http_pkg_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	struct_util "github.com/harmonify/movie-reservation-system/pkg/util/struct"
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
	app          *fx.App
	httpResponse http_pkg.HttpResponse
}

type testConfig struct {
	SuccessHttpCode int
	Data            string
	Error           *error_pkg.ErrorWithDetails
}

type testExpectation struct {
	Success  bool
	HttpCode int
	Result   string
	Error    *error_pkg.ErrorWithDetails
}

func (s *ResponseTestSuite) SetupSuite() {
	s.app = fx.New(
		logger.LoggerModule,
		tracer.TracerModule,
		error_pkg.ErrorModule,
		fx.Provide(
			func() *config.Config {
				return &config.Config{
					Env:      "test",
					LogType:  "nop",
					LogLevel: "debug",
				}
			},
			struct_util.NewStructUtil,
			http_pkg.NewHttpResponse,
		),
		fx.Invoke(func(response http_pkg.HttpResponse) {
			s.httpResponse = response
		}),

		fx.NopLogger,
	)
}

func (s *ResponseTestSuite) TestHttpResponse_Build() {
	testCases := []test_interface.TestCase[testConfig, testExpectation]{
		{
			Description: "Should build success response",
			Config: testConfig{
				SuccessHttpCode: 200,
				Data:            "Test data",
				Error:           nil,
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
				Data:  "Test data",
				Error: error_pkg.InvalidRequestBodyError,
			},
			Expectation: testExpectation{
				Success:  false,
				HttpCode: 400,
				Result:   "Test data",
				Error:    error_pkg.InvalidRequestBodyError,
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			httpCode, httpResponse, httpResponseError := s.httpResponse.Build(
				context.Background(),
				testCase.Config.SuccessHttpCode,
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
				s.Require().IsType(httpResponseError, &http_pkg.HttpError{})
				s.Assert().Equal(
					httpResponseError.Code,
					testCase.Expectation.Error.Code,
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

			httpCode, response, responseError := s.httpResponse.Build(context.Background(), testCase.Config.SuccessHttpCode, testCase.Config.Data, nil)

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
		s.httpResponse.Send(c, nil, nil)

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
		s.httpResponse.Send(c, nil, error_pkg.InternalServerError)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "error")
	})
}
