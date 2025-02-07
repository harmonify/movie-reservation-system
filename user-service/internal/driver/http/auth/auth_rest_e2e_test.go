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
	"reflect"
	"testing"
	"time"

	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/database"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"
	"github.com/harmonify/movie-reservation-system/user-service/internal"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	entityfactory "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/factory"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/seeder"
	http_driver "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http"
	auth_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/auth"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/auth/test"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func TestAuthRest(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}
	os.Setenv("ENV", config_pkg.EnvironmentTest)
	suite.Run(t, new(AuthRestTestSuite))
}

type AuthRestTestSuite struct {
	suite.Suite
	app         *fx.App
	httpServer  *http_driver.HttpServer
	otpCache    shared.OtpCache
	database    *database.Database
	userStorage shared.UserStorage
	userFactory entityfactory.UserFactory
	userSeeder  seeder.UserSeeder
}

func (s *AuthRestTestSuite) SetupSuite() {
	s.app = internal.NewApp(
		fx.Decorate(func(l logger.Logger) logger.Logger {
			return &logger.ConsoleLoggerImpl{
				Logger: l.GetZapLogger().WithOptions(zap.IncreaseLevel(zap.InfoLevel)),
			}
		}),
		seeder.DrivenPostgresqlSeederModule,
		entityfactory.UserEntityFactoryModule,
		fx.Invoke(func(
			httpServer *http_driver.HttpServer,
			otpCache shared.OtpCache,
			database *database.Database,
			userStorage shared.UserStorage,
			userFactory entityfactory.UserFactory,
			userSeeder seeder.UserSeeder,
		) {
			s.httpServer = httpServer
			s.otpCache = otpCache
			s.database = database
			s.userStorage = userStorage
			s.userFactory = userFactory
			s.userSeeder = userSeeder
		}),
		fx.NopLogger,
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

	testCases := []func() test_interface.HttpTestCase[auth_rest.PostRegisterReq, interface{}]{
		func() test_interface.HttpTestCase[auth_rest.PostRegisterReq, interface{}] {
			testUser, _, err := s.userFactory.GenerateUser()
			s.Require().NoError(err)
			return test_interface.HttpTestCase[auth_rest.PostRegisterReq, interface{}]{
				Description: "It should return a 200 OK response",
				Config: test_interface.Request[auth_rest.PostRegisterReq]{
					RequestBody: auth_rest.PostRegisterReq{
						Username:    testUser.Username,
						Password:    testUser.Password,
						Email:       testUser.Email,
						PhoneNumber: testUser.PhoneNumber,
						FirstName:   testUser.FirstName,
						LastName:    testUser.LastName,
					},
				},
				Expectation: test_interface.ResponseExpectation[interface{}]{
					ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
					ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
				},
			}
		},
	}

	for _, tc := range testCases {
		ctx := context.Background()

		testCase := tc()
		defer func() {
			user, err := s.userStorage.FindUser(ctx, entity.FindUser{Username: sql.NullString{String: testCase.Config.RequestBody.Username, Valid: true}})
			if err != nil {
				s.T().Log("Failed to find test user after call")
				return
			}
			if err := s.userSeeder.DeleteUser(ctx, entity.FindUser{Username: sql.NullString{String: testCase.Config.RequestBody.Username, Valid: true}}); err != nil {
				s.T().Log("Failed to delete test user after call")
				return
			}
			if _, err := s.otpCache.DeleteEmailVerificationCode(ctx, user.UUID); err != nil {
				s.T().Log("Failed to delete otp cache after call")
			}
		}()

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

func (s *AuthRestTestSuite) TestAuthRest_PostLogin() {
	var (
		PATH                   = "/v1/login"
		METHOD                 = "POST"
		refreshTokenCookieName = http_pkg.HttpCookiePrefix + "token"
	)

	ctx := context.Background()
	testUser, err := s.userSeeder.CreateUser(ctx)
	s.Require().NoError(err)
	defer func() {
		if err := s.userSeeder.DeleteUser(ctx, entity.FindUser{UUID: sql.NullString{String: testUser.User.UUID, Valid: true}}); err != nil {
			s.T().Log("Failed to delete test user before call")
		}
	}()

	testCases := []test_interface.TestCase[auth_rest.PostLoginReq, test.PostRegisterTestExpectation]{
		{
			Description: "User exist and correct password should return a 200 OK response",
			Config: auth_rest.PostLoginReq{
				Username: testUser.User.Username,
				Password: testUser.UserRaw.Password,
			},
			Expectation: test.PostRegisterTestExpectation{
				ResponseStatusCode:                   test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus:                   test_interface.NullBool{Bool: true, Valid: true},
				ResponseHeaderRefreshTokenExist:      test_interface.NullBool{Bool: true, Valid: true},
				ResponseBodyAccessTokenExist:         test_interface.NullBool{Bool: true, Valid: true},
				ResponseBodyAccessTokenDurationExist: test_interface.NullBool{Bool: true, Valid: true},
			},
		},
		{
			Description: "User not exist should return a 404 Not Found response",
			Config: auth_rest.PostLoginReq{
				Username: "nonexistentuser@example.com",
				Password: "password",
			},
			Expectation: test.PostRegisterTestExpectation{
				ResponseStatusCode:                   test_interface.NullInt{Int: http.StatusNotFound, Valid: true},
				ResponseBodyStatus:                   test_interface.NullBool{Bool: false, Valid: true},
				ResponseHeaderRefreshTokenExist:      test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyAccessTokenExist:         test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyAccessTokenDurationExist: test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyErrorCode:                test_interface.NullString{String: "ACCOUNT_NOT_FOUND", Valid: true},
				ResponseBodyErrorMessage:             test_interface.NullString{String: "The account you're trying to access is not found. Please register an account or check the username you've entered.", Valid: true},
				ResponseBodyErrorObject:              make([]interface{}, 0),
			},
		},
		{
			Description: "User exist and incorrect password should return a 403 Forbidden response",
			Config: auth_rest.PostLoginReq{
				Username: testUser.User.Username,
				Password: "incorrect_password",
			},
			Expectation: test.PostRegisterTestExpectation{
				ResponseStatusCode:                   test_interface.NullInt{Int: http.StatusForbidden, Valid: true},
				ResponseBodyStatus:                   test_interface.NullBool{Bool: false, Valid: true},
				ResponseHeaderRefreshTokenExist:      test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyAccessTokenExist:         test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyAccessTokenDurationExist: test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyErrorCode:                test_interface.NullString{String: "INCORRECT_PASSWORD", Valid: true},
				ResponseBodyErrorMessage:             test_interface.NullString{String: "The password you've entered is incorrect. Please try again.", Valid: true},
				ResponseBodyErrorObject:              make([]interface{}, 0),
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			jsonPayload, err := json.Marshal(testCase.Config)
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

			if testCase.Expectation.ResponseHeaderRefreshTokenExist.Valid {
				cookies := w.Result().Cookies()

				if testCase.Expectation.ResponseHeaderRefreshTokenExist.Bool {
					s.Require().NotEmpty(cookies)
					var cookie *http.Cookie
					for _, c := range cookies {
						if c.Name == refreshTokenCookieName {
							cookie = c
							break
						}
					}
					// MaxAge   int
					// Domain   string
					// Path     string
					// Secure   bool
					// HttpOnly bool
					s.Require().NotNil(cookie)
					s.Require().NotEmpty(cookie.Value)
				} else {
					s.Require().Empty(cookies)
				}
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
				s.Require().Equal(testCase.Expectation.ResponseStatusCode.Int, w.Result().StatusCode)
			}
			if testCase.Expectation.ResponseBodyStatus.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyStatus.Bool, status)
			}
			if testCase.Expectation.ResponseBodyAccessTokenExist.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyAccessTokenDurationExist.Bool, resultBody.Get("access_token").Exists())
			}
			if testCase.Expectation.ResponseBodyAccessTokenDurationExist.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyAccessTokenDurationExist.Bool, resultBody.Get("access_token_duration").Exists())
			}
			if testCase.Expectation.ResponseBodyErrorCode.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyErrorCode.String, responseError.Get("code").String())
			}
			if testCase.Expectation.ResponseBodyErrorMessage.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyErrorMessage.String, responseError.Get("message").String())
			}
			if testCase.Expectation.ResponseBodyErrorObject != nil {
				s.Require().True(responseError.Get("errors").IsArray(), "Expected 'errors' to be an array")
			}
		})
	}
}

func (s *AuthRestTestSuite) TestAuthRest_GetToken() {
	var (
		PATH                   = "/v1/token"
		METHOD                 = "GET"
		refreshTokenCookieName = http_pkg.HttpCookiePrefix + "token"
	)

	ctx := context.Background()
	testUser, err := s.userSeeder.CreateUser(ctx)
	s.Require().NoError(err)
	s.Require().Less(time.Now(), testUser.UserSessions[0].ExpiredAt)
	defer func() {
		if err := s.userSeeder.DeleteUser(ctx, entity.FindUser{UUID: sql.NullString{String: testUser.User.UUID, Valid: true}}); err != nil {
			s.T().Log("Failed to delete test user before call")
		}
	}()

	testCases := []test_interface.HttpTestCase[any, test.GetTokenTestExpectation]{
		{
			Description: "Refresh token exist should return a 200 OK response",
			Expectation: test_interface.ResponseExpectation[test.GetTokenTestExpectation]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
				ResponseBodyResult: test.GetTokenTestExpectation{
					AccessTokenExist:         test_interface.NullBool{Bool: true, Valid: true},
					AccessTokenDurationExist: test_interface.NullBool{Bool: true, Valid: true},
				},
			},
			BeforeCall: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:     refreshTokenCookieName,
					Value:    testUser.UserSessionRaws[0].RefreshToken,
					Path:     "/token",
					Domain:   "localhost",
					MaxAge:   2592000,
					HttpOnly: true,
					Secure:   true,
				})
			},
		},
		{
			Description: "Refresh token not exist should return a 401 Unauthorized response",
			Expectation: test_interface.ResponseExpectation[test.GetTokenTestExpectation]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusUnauthorized, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult: test.GetTokenTestExpectation{
					AccessTokenExist:         test_interface.NullBool{Bool: false, Valid: true},
					AccessTokenDurationExist: test_interface.NullBool{Bool: false, Valid: true},
				},
				ResponseBodyErrorCode:    test_interface.NullString{String: "REFRESH_TOKEN_EXPIRED", Valid: true},
				ResponseBodyErrorMessage: test_interface.NullString{String: "Your session has expired. Please login again.", Valid: true},
				ResponseBodyErrorObject:  make([]interface{}, 0),
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			jsonPayload, err := json.Marshal(testCase.Config)
			s.Require().NoError(err)

			req, err := http.NewRequest(METHOD, PATH, bytes.NewBuffer(jsonPayload))
			s.Require().NoError(err)

			req.Header.Set("Content-Type", "application/json")

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

			if testCase.Expectation.ResponseStatusCode.Valid {
				s.Require().Equal(testCase.Expectation.ResponseStatusCode.Int, w.Result().StatusCode)
			}
			if testCase.Expectation.ResponseBodyStatus.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyStatus.Bool, status)
			}
			if testCase.Expectation.ResponseBodyResult.AccessTokenExist.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyResult.AccessTokenExist.Bool, resultBody.Get("access_token").Exists())
			}
			if testCase.Expectation.ResponseBodyResult.AccessTokenDurationExist.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyResult.AccessTokenDurationExist.Bool, resultBody.Get("access_token_duration").Exists())
			}
			if testCase.Expectation.ResponseBodyErrorCode.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyErrorCode.String, responseError.Get("code").String())
			}
			if testCase.Expectation.ResponseBodyErrorMessage.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyErrorMessage.String, responseError.Get("message").String())
			}
			if testCase.Expectation.ResponseBodyErrorObject != nil {
				s.Require().True(responseError.Get("errors").IsArray(), "Expected 'errors' to be an array")
			}
		})
	}
}

func (s *AuthRestTestSuite) TestAuthRest_PostLogout() {
	var (
		PATH                   = "/v1/logout"
		METHOD                 = "POST"
		refreshTokenCookieName = http_pkg.HttpCookiePrefix + "token"
	)

	ctx := context.Background()
	testUser, err := s.userSeeder.CreateUser(ctx)
	s.Require().NoError(err)
	s.Require().Less(time.Now(), testUser.UserSessions[0].ExpiredAt)
	defer func() {
		if err := s.userSeeder.DeleteUser(ctx, entity.FindUser{UUID: sql.NullString{String: testUser.User.UUID, Valid: true}}); err != nil {
			s.T().Log("Failed to delete test user before call")
		}
	}()

	testCases := []test_interface.HttpTestCase[any, any]{
		{
			Description: "Refresh token exist should return a 200 OK response",
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
			},
			BeforeCall: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:     refreshTokenCookieName,
					Value:    testUser.UserSessionRaws[0].RefreshToken,
					Path:     "/token",
					Domain:   "localhost",
					MaxAge:   86400,
					HttpOnly: true,
					Secure:   true,
				})
			},
		},
		{
			Description: "Refresh token not exist should return a 200 OK response",
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			jsonPayload, err := json.Marshal(testCase.Config)
			s.Require().NoError(err)

			req, err := http.NewRequest(METHOD, PATH, bytes.NewBuffer(jsonPayload))
			s.Require().NoError(err)

			req.Header.Set("Content-Type", "application/json")

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

			if testCase.Expectation.ResponseStatusCode.Valid {
				s.Require().Equal(testCase.Expectation.ResponseStatusCode.Int, w.Result().StatusCode)
			}
			if testCase.Expectation.ResponseBodyStatus.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyStatus.Bool, status)
			}
			if testCase.Expectation.ResponseBodyErrorCode.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyErrorCode.String, responseError.Get("code").String())
			}
			if testCase.Expectation.ResponseBodyErrorMessage.Valid {
				s.Require().Equal(testCase.Expectation.ResponseBodyErrorMessage.String, responseError.Get("message").String())
			}
			if testCase.Expectation.ResponseBodyErrorObject != nil {
				s.Require().True(responseError.Get("errors").IsArray(), "Expected 'errors' to be an array")
			}
		})
	}
}
