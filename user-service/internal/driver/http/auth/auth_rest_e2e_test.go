package auth_rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/internal"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	http_driver "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http"
	auth_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/auth"
	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	"github.com/harmonify/movie-reservation-system/user-service/lib/test/interface"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
)

func TestAuthRest(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(AuthRestTestSuite))
}

type postUserRegisterTestConfig struct {
	Data auth_rest.PostUserRegisterReq
}

type postUserLoginTestConfig struct {
	Data auth_rest.PostUserRegisterReq
}

type postUserLoginTestExpectation struct {
	Result auth_rest.PostUserLoginRes
}

type AuthRestTestSuite struct {
	suite.Suite
	app        *fx.App
	database   *database.Database
	httpServer *http_driver.HttpServer
	authRest   auth_rest.AuthRestHandler
}

func (s *AuthRestTestSuite) SetupSuite() {
	s.app = internal.NewApp(
		fx.Invoke(func(
			authRest auth_rest.AuthRestHandler,
			database *database.Database,
			httpServer *http_driver.HttpServer,
			handlers http_driver.RestHandlers,
		) {
			s.authRest = authRest
			s.database = database
			s.httpServer = httpServer
			s.httpServer.Configure(handlers...)
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*105)
	defer cancel()

	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *AuthRestTestSuite) TestPostUserRegister() {
	testCases := []test_interface.HttpTestCase[auth_rest.PostUserRegisterReq, any]{
		{
			Description: "It should return a 200 OK response",
			Config: test_interface.Request[auth_rest.PostUserRegisterReq]{
				RequestBody: auth_rest.PostUserRegisterReq{
					Username:    "user1234",
					Password:    "user1234",
					Email:       "user1234@example.com",
					PhoneNumber: "+6281234567890",
					FirstName:   "Example",
					LastName:    "User",
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode: http.StatusOK,
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
			},
			BeforeCall: func(req *http.Request) {
				user := &model.User{}
				s.database.DB.First(user, &model.User{Username: "user1234"})
				s.database.DB.Exec("DELETE FROM user_keys WHERE user_uuid = ?", user.UUID.String())
				s.database.DB.Exec("DELETE FROM users WHERE uuid = ?", user.UUID.String())
			},
			AfterCall: func(w *httptest.ResponseRecorder) {
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			jsonPayload, err := json.Marshal(testCase.Config.RequestBody)
			s.Require().NoError(err)

			req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonPayload))
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
			status := body.Get("success").Bool()
			responseError := body.Get("error")
			resultBody := body.Get("result")

			if testCase.Expectation.ResponseStatusCode != 0 {
				s.Require().Equal(testCase.Expectation.ResponseStatusCode, w.Result().StatusCode)
			}
			if testCase.Expectation.ResponseBodyStatus.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyStatus.Bool, status)
			}
			if testCase.Expectation.ResponseBodyResult != nil {
				expected, err := json.Marshal(testCase.Expectation.ResponseBodyResult)
				s.Require().NoError(err)
				s.Require().JSONEq(string(expected), resultBody.Raw)
			}
			if testCase.Expectation.ResponseBodyErrorCode != "" {
				s.Require().Equal(testCase.Expectation.ResponseBodyErrorCode, responseError.Get("code").String())
			}
			if testCase.Expectation.ResponseBodyErrorMessage != "" {
				s.Require().Equal(testCase.Expectation.ResponseBodyErrorMessage, responseError.Get("message").String())
			}
			if testCase.Expectation.ResponseBodyErrorObject != nil {
				s.Require().True(responseError.Get("errors").IsArray(), "Expected 'errors' to be an array")
				for i, errData := range testCase.Expectation.ResponseBodyErrorObject {
					s.Equal(errData.Field, responseError.Get("errors").Array()[i].Get("field").String())
					s.Equal(errData.Message, responseError.Get("errors").Array()[i].Get("message").String())
				}
			}
		})
	}
}
