package health_check_rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"
	"github.com/harmonify/movie-reservation-system/user-service/internal"
	http_driver "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http"
	health_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/health_check"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
)

func TestHealthCheckRest(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(HealthCheckRestTestSuite))
}

type HealthCheckRestTestSuite struct {
	suite.Suite
	app        *fx.App
	httpServer *http_driver.HttpServer
}

func (s *HealthCheckRestTestSuite) SetupSuite() {
	s.app = internal.NewApp(
		fx.Invoke(func(
			httpServer *http_driver.HttpServer,
		) {
			s.httpServer = httpServer
		}),
		fx.NopLogger,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*105)
	defer cancel()

	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *HealthCheckRestTestSuite) TestHealthCheckRest_GetHealthCheck() {
	var (
		PATH   = "/v1/health"
		METHOD = "GET"
	)

	testCases := []test_interface.HttpTestCase[interface{}, *health_rest.HealthCheckResponse]{
		{
			Description: "It should return a 200 OK response",
			Expectation: test_interface.ResponseExpectation[*health_rest.HealthCheckResponse]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
				ResponseBodyResult: &health_rest.HealthCheckResponse{
					Ok: true,
				},
				ResponseBodyErrorCode:    test_interface.NullString{String: "", Valid: false},
				ResponseBodyErrorMessage: test_interface.NullString{String: "", Valid: false},
				ResponseBodyErrorObject:  nil,
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			jsonPayload, err := json.Marshal(testCase.Config.RequestBody)
			s.Require().NoError(err)

			req, err := http.NewRequest(METHOD, PATH, bytes.NewBuffer(jsonPayload))
			s.Require().NoError(err)

			req.Header.Set("Content-Type", "application/json")
			if testCase.Config.RequestHeader != nil && len(testCase.Config.RequestHeader) > 0 {
				for _, rh := range testCase.Config.RequestHeader {
					req.Header.Set(rh.Key, rh.Value)
				}
			}

			if testCase.Config.RequestQuery != nil && len(testCase.Config.RequestQuery) > 0 {
				q := req.URL.Query()
				for _, rq := range testCase.Config.RequestQuery {
					q.Set(rq.Key, rq.Value)
				}
				req.URL.RawQuery = q.Encode()
			}

			if testCase.BeforeCall != nil {
				testCase.BeforeCall(req)
			}

			w := httptest.NewRecorder()
			s.httpServer.Gin.ServeHTTP(w, req)

			if testCase.AfterCall != nil {
				testCase.AfterCall(w)
			}

			bodyString := w.Body.String()

			s.Require().True(
				gjson.Valid(bodyString),
				fmt.Sprintf("response body should be a valid JSON, but got %s", bodyString),
			)
			body := gjson.Parse(bodyString)
			s.T().Log(body)
			status := body.Get("success").Bool()
			responseError := body.Get("error")
			resultBody := body.Get("result")

			if testCase.Expectation.ResponseStatusCode.Valid {
				s.Require().Equal(testCase.Expectation.ResponseStatusCode.Int, w.Result().StatusCode)
			}
			if testCase.Expectation.ResponseBodyStatus.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyStatus.Bool, status)
			}
			if testCase.Expectation.ResponseBodyResult != nil {
				expected, err := json.Marshal(testCase.Expectation.ResponseBodyResult)
				s.Require().NoError(err)
				s.Require().JSONEq(string(expected), resultBody.Raw)
			}
			if testCase.Expectation.ResponseBodyErrorCode.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyErrorCode.String, responseError.Get("code").String())
			}
			if testCase.Expectation.ResponseBodyErrorMessage.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyErrorMessage.String, responseError.Get("message").String())
			}
			if testCase.Expectation.ResponseBodyErrorObject != nil {
				s.Require().True(responseError.Get("errors").IsArray(), "Expected 'errors' to be an array")
				for i, errData := range testCase.Expectation.ResponseBodyErrorObject {
					if expectedErrorObject, ok := errData.(validation.ValidationError); ok {
						s.Equal(expectedErrorObject.Field, responseError.Get("errors").Array()[i].Get("field").String())
						s.Equal(expectedErrorObject.Message, responseError.Get("errors").Array()[i].Get("message").String())
					} else {
						s.T().Fatalf("Expected error object to be %s, but got %s", reflect.TypeFor[validation.ValidationError]().Name(), reflect.TypeOf(errData).Name())
					}
				}
			}
		})
	}
}
