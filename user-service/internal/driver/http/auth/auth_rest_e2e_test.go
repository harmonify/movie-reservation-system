package auth_rest_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/internal"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/seeder"
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

type postRegisterTestConfig struct {
	Data auth_rest.PostRegisterReq
}

type postLoginTestConfig struct {
	Data auth_rest.PostRegisterReq
}

type postLoginTestExpectation struct {
	Result auth_rest.PostLoginRes
}

type AuthRestTestSuite struct {
	suite.Suite
	app         *fx.App
	database    *database.Database
	httpServer  *http_driver.HttpServer
	authRest    auth_rest.AuthRestHandler
	authService auth_service.AuthService
	userSeeder  seeder.UserSeeder
	userStorage shared_service.UserStorage
	otpStorage  shared_service.OtpStorage
}

func (s *AuthRestTestSuite) SetupSuite() {
	s.app = internal.NewApp(
		seeder.DrivenPostgresqlSeederModule,
		fx.Invoke(func(
			authRest auth_rest.AuthRestHandler,
			database *database.Database,
			httpServer *http_driver.HttpServer,
			handlers http_driver.RestHandlers,
			authService auth_service.AuthService,
			userSeeder seeder.UserSeeder,
			userStorage shared_service.UserStorage,
			otpStorage shared_service.OtpStorage,
		) {
			s.authRest = authRest
			s.database = database
			s.httpServer = httpServer
			s.httpServer.Configure(handlers...)
			s.authService = authService
			s.userSeeder = userSeeder
			s.userStorage = userStorage
			s.otpStorage = otpStorage
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*105)
	defer cancel()

	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *AuthRestTestSuite) TestAuthRest_PostRegister() {
	var (
		PATH   = "/v1/register"
		METHOD = "POST"
	)

	testCases := []test_interface.HttpTestCase[auth_rest.PostRegisterReq, interface{}]{
		{
			Description: "It should return a 200 OK response",
			Config: test_interface.Request[auth_rest.PostRegisterReq]{
				RequestBody: auth_rest.PostRegisterReq{
					Username:    seeder.TestUser.Username,
					Password:    seeder.TestUser.Password,
					Email:       seeder.TestUser.Email,
					PhoneNumber: seeder.TestUser.PhoneNumber,
					FirstName:   seeder.TestUser.FirstName,
					LastName:    seeder.TestUser.LastName,
				},
			},
			Expectation: test_interface.ResponseExpectation[interface{}]{
				ResponseStatusCode: http.StatusOK,
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
			},
			BeforeCall: func(req *http.Request) {
				if err := s.userSeeder.DeleteTestUser(); err != nil {
					s.T().Log("Failed to delete test user before call")
				}
				if _, err := s.otpStorage.DeleteEmailVerificationToken(context.Background(), seeder.TestUser.Email); err != nil {
					s.T().Log("Failed to delete test user email verification token before call")
				}
			},
			AfterCall: func(w *httptest.ResponseRecorder) {
				if err := s.userSeeder.DeleteTestUser(); err != nil {
					s.T().Log("Failed to delete test user after call")
				}
				if _, err := s.otpStorage.DeleteEmailVerificationToken(context.Background(), seeder.TestUser.Email); err != nil {
					s.T().Log("Failed to delete test user email verification token after call")
				}
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

func (s *AuthRestTestSuite) TestAuthRest_PostLogin() {
	var (
		PATH   = "/v1/login"
		METHOD = "POST"
	)

	testCases := []test_interface.HttpTestCase[auth_rest.PostLoginReq, interface{}]{
		{
			Description: "It should return a 200 OK response",
			Config: test_interface.Request[auth_rest.PostLoginReq]{
				RequestBody: auth_rest.PostLoginReq{
					Username: "user1234",
					Password: "user1234",
				},
			},
			Expectation: test_interface.ResponseExpectation[interface{}]{
				ResponseStatusCode: http.StatusOK,
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
			},
			BeforeCall: func(req *http.Request) {
				if _, err := s.userSeeder.CreateTestUser(); err != nil {
					s.T().Log("Failed to create test user before call")
				}
			},
			AfterCall: func(w *httptest.ResponseRecorder) {
				if err := s.userSeeder.DeleteTestUser(); err != nil {
					s.T().Log("Failed to delete test user before call")
				}
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

type postVerifyEmailTestConfig func() auth_rest.PostVerifyEmailReq
type postVerifyEmailTestExpectation struct {
	ResponseStatusCode       test_interface.NullInt
	ResponseBodyStatus       test_interface.NullBool
	ResponseBodyResult       interface{}
	ResponseBodyErrorCode    test_interface.NullString
	ResponseBodyErrorMessage test_interface.NullString
	ResponseBodyErrorObject  []interface{}
	IsEmailVerified          test_interface.NullBool
}

func (s *AuthRestTestSuite) TestAuthRest_PostVerifyEmail() {
	var (
		PATH   = "/v1/register/verify"
		METHOD = "POST"
	)

	testCases := []test_interface.TestCase[postVerifyEmailTestConfig, postVerifyEmailTestExpectation]{
		{
			Description: "It should return a 200 OK response",
			Config: func() auth_rest.PostVerifyEmailReq {
				token := "123456"
				err := s.otpStorage.SaveEmailVerificationToken(context.Background(), shared_service.SaveEmailVerificationTokenParam{
					Email: seeder.TestUser.Email,
					Token: token,
					TTL:   time.Minute * 5,
				})
				s.Require().NoError(err)
				return auth_rest.PostVerifyEmailReq{
					Email: seeder.TestUser.Email,
					Token: token,
				}
			},
			Expectation: postVerifyEmailTestExpectation{
				ResponseStatusCode:       test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus:       test_interface.NullBool{Bool: true, Valid: true},
				ResponseBodyResult:       make(map[string]interface{}, 0),
				ResponseBodyErrorCode:    test_interface.NullString{String: "", Valid: false},
				ResponseBodyErrorMessage: test_interface.NullString{String: "", Valid: false},
				IsEmailVerified: test_interface.NullBool{
					Bool:  true,
					Valid: true,
				},
			},
		},
		{
			Description: "It should return a 403 Forbidden response",
			Config: func() auth_rest.PostVerifyEmailReq {
				token := "123456"
				err := s.otpStorage.SaveEmailVerificationToken(context.Background(), shared_service.SaveEmailVerificationTokenParam{
					Email: seeder.TestUser.Email,
					Token: token,
					TTL:   time.Minute * 5,
				})
				s.Require().NoError(err)
				return auth_rest.PostVerifyEmailReq{
					Email: seeder.TestUser.Email,
					Token: "INCORRECT",
				}
			},
			Expectation: postVerifyEmailTestExpectation{
				ResponseStatusCode:       test_interface.NullInt{Int: http.StatusForbidden, Valid: true},
				ResponseBodyStatus:       test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:       make(map[string]interface{}, 0),
				ResponseBodyErrorCode:    test_interface.NullString{String: "VERIFICATION_TOKEN_INVALID", Valid: false},
				ResponseBodyErrorMessage: test_interface.NullString{String: "Failed to verify your email. Please try to request a new verification link.", Valid: false},
				ResponseBodyErrorObject:  make([]interface{}, 0),
				IsEmailVerified: test_interface.NullBool{
					Bool:  true,
					Valid: true,
				},
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			if _, err := s.userSeeder.CreateTestUser(); err != nil {
				s.T().Log("Failed to create test user before call")
			}
			defer func() {
				if err := s.userSeeder.DeleteTestUser(); err != nil {
					s.T().Log("Failed to delete test user before call")
				}
			}()

			jsonPayload, err := json.Marshal(testCase.Config())
			s.Require().NoError(err)

			req, err := http.NewRequest(METHOD, PATH, bytes.NewBuffer(jsonPayload))
			s.Require().NoError(err)

			req.Header.Set("Content-Type", "application/json")

			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			w := httptest.NewRecorder()
			s.httpServer.Gin.ServeHTTP(w, req)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
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

			if testCase.Expectation.ResponseStatusCode.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseStatusCode.Int, w.Result().StatusCode)
			}
			if testCase.Expectation.ResponseBodyStatus.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyStatus.Bool, status)
			}
			if testCase.Expectation.ResponseBodyResult != nil {
				expected, err := json.Marshal(testCase.Expectation.ResponseBodyResult)
				s.Assert().NoError(err)
				s.Assert().JSONEq(string(expected), resultBody.Raw)
			}
			if testCase.Expectation.ResponseBodyErrorCode.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorCode.String, responseError.Get("code").String())
			}
			if testCase.Expectation.ResponseBodyErrorMessage.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorMessage.String, responseError.Get("message").String())
			}
			if testCase.Expectation.ResponseBodyErrorObject != nil {
				s.Assert().True(responseError.Get("errors").IsArray(), "Expected 'errors' to be an array but got %s", responseError.Get("errors").Raw)
			}
			if testCase.Expectation.IsEmailVerified.Valid {
				user, err := s.userStorage.FindUser(context.Background(), entity.FindUser{
					Email: sql.NullString{String: seeder.TestUser.Email, Valid: true},
				})
				s.Assert().NoError(err)
				s.Assert().Equal(testCase.Expectation.IsEmailVerified.Bool, user.IsEmailVerified)
			}
		})
	}
}
