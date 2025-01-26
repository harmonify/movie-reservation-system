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

	"github.com/harmonify/movie-reservation-system/pkg/database"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"
	"github.com/harmonify/movie-reservation-system/user-service/internal"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/factory"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/seeder"
	http_driver "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http"
	auth_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/auth"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/auth/test"
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

type AuthRestTestSuite struct {
	suite.Suite
	app                    *fx.App
	database               *database.Database
	httpServer             *http_driver.HttpServer
	testUser               *model.User
	testUserHashedPassword *model.User
	userSessionFactory     factory.UserSessionFactory
	userSeeder             seeder.UserSeeder
	userSessionSeeder      seeder.UserSessionSeeder
	userStorage            shared.UserStorage
	otpStorage             shared.OtpStorage
}

func (s *AuthRestTestSuite) SetupSuite() {
	s.app = internal.NewApp(
		seeder.DrivenPostgresqlSeederModule,
		factory.DrivenPostgresqlFactoryModule,
		fx.Invoke(func(
			database *database.Database,
			httpServer *http_driver.HttpServer,
			authService auth_service.AuthService,
			userFactory factory.UserFactory,
			userSessionFactory factory.UserSessionFactory,
			userSeeder seeder.UserSeeder,
			userSessionSeeder seeder.UserSessionSeeder,
			userStorage shared.UserStorage,
			otpStorage shared.OtpStorage,
		) {
			s.database = database
			s.httpServer = httpServer
			s.testUser = userFactory.CreateTestUser(factory.CreateTestUserParam{HashPassword: false})
			s.testUserHashedPassword = userFactory.CreateTestUser(factory.CreateTestUserParam{HashPassword: true})
			s.userSeeder = userSeeder
			s.userSessionFactory = userSessionFactory
			s.userSessionSeeder = userSessionSeeder
			s.userStorage = userStorage
			s.otpStorage = otpStorage
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

	testCases := []test_interface.HttpTestCase[auth_rest.PostRegisterReq, interface{}]{
		{
			Description: "It should return a 200 OK response",
			Config: test_interface.Request[auth_rest.PostRegisterReq]{
				RequestBody: auth_rest.PostRegisterReq{
					Username:    s.testUser.Username,
					Password:    s.testUser.Password,
					Email:       s.testUser.Email,
					PhoneNumber: s.testUser.PhoneNumber,
					FirstName:   s.testUser.FirstName,
					LastName:    s.testUser.LastName,
				},
			},
			Expectation: test_interface.ResponseExpectation[interface{}]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
			},
			BeforeCall: func(req *http.Request) {
				if err := s.userSeeder.DeleteUser(s.testUser.Username); err != nil {
					s.T().Log("Failed to delete test user before call")
				}
				if _, err := s.otpStorage.DeleteEmailVerificationToken(context.Background(), s.testUser.Email); err != nil {
					s.T().Log("Failed to delete test user email verification token before call")
				}
			},
			AfterCall: func(w *httptest.ResponseRecorder) {
				if err := s.userSeeder.DeleteUser(s.testUser.Username); err != nil {
					s.T().Log("Failed to delete test user after call")
				}
				if _, err := s.otpStorage.DeleteEmailVerificationToken(context.Background(), s.testUser.Email); err != nil {
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

func (s *AuthRestTestSuite) TestAuthRest_PostVerifyEmail() {
	var (
		PATH   = "/v1/register/verify"
		METHOD = "POST"
	)

	testCases := []test_interface.TestCase[test.PostVerifyEmailTestConfig, test.PostVerifyEmailTestExpectation]{
		{
			Description: "It should return a 200 OK response",
			Config: func() auth_rest.PostVerifyEmailReq {
				token := "123456"
				err := s.otpStorage.SaveEmailVerificationToken(context.Background(), shared.SaveEmailVerificationTokenParam{
					Email: s.testUser.Email,
					Token: token,
					TTL:   time.Minute * 5,
				})
				s.Require().NoError(err)
				return auth_rest.PostVerifyEmailReq{
					Email: s.testUser.Email,
					Token: token,
				}
			},
			Expectation: test.PostVerifyEmailTestExpectation{
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
				err := s.otpStorage.SaveEmailVerificationToken(context.Background(), shared.SaveEmailVerificationTokenParam{
					Email: s.testUser.Email,
					Token: token,
					TTL:   time.Minute * 5,
				})
				s.Require().NoError(err)
				return auth_rest.PostVerifyEmailReq{
					Email: s.testUser.Email,
					Token: "INCORRECT",
				}
			},
			Expectation: test.PostVerifyEmailTestExpectation{
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
			if _, err := s.userSeeder.SaveUser(*s.testUserHashedPassword); err != nil {
				s.T().Log("Failed to create test user before call")
			}
			defer func() {
				if err := s.userSeeder.DeleteUser(s.testUser.Username); err != nil {
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
					Email: sql.NullString{String: s.testUser.Email, Valid: true},
				})
				s.Assert().NoError(err)
				s.Assert().Equal(testCase.Expectation.IsEmailVerified.Bool, user.IsEmailVerified)
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

	testCases := []test_interface.TestCase[auth_rest.PostLoginReq, test.PostRegisterTestExpectation]{
		{
			Description: "User exist and correct password should return a 200 OK response",
			Config: auth_rest.PostLoginReq{
				Username: s.testUser.Username,
				Password: s.testUser.Password,
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
				Username: s.testUser.Username,
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
			if _, err := s.userSeeder.SaveUser(*s.testUserHashedPassword); err != nil {
				s.T().Log("Failed to create test user before call")
			}
			defer func() {
				if err := s.userSeeder.DeleteUser(s.testUser.Username); err != nil {
					s.T().Log("Failed to delete test user before call")
				}
			}()

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
				session, hashedRefreshToken := s.userSessionFactory.CreateUserSession(factory.CreateUserSessionParam{
					UserUUID:         s.testUser.UUID.String(),
					HashRefreshToken: false,
				})
				s.Require().Less(time.Now(), session.ExpiredAt)
				unhashedRefreshToken := session.RefreshToken
				session.RefreshToken = hashedRefreshToken
				session, err := s.userSessionSeeder.SaveUserSession(*session)
				s.Require().NoError(err)
				req.AddCookie(&http.Cookie{
					Name:     refreshTokenCookieName,
					Value:    unhashedRefreshToken,
					Path:     "/user/token",
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
				ResponseBodyErrorCode:    test_interface.NullString{String: "INVALID_REFRESH_TOKEN", Valid: true},
				ResponseBodyErrorMessage: test_interface.NullString{String: "Your session is expired. Please login again.", Valid: true},
				ResponseBodyErrorObject:  make([]interface{}, 0),
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			if _, err := s.userSeeder.SaveUser(*s.testUserHashedPassword); err != nil {
				s.T().Log("Failed to create test user before call")
			}
			defer func() {
				if err := s.userSeeder.DeleteUser(s.testUser.Username); err != nil {
					s.T().Log("Failed to delete test user before call")
				}
			}()

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

	testCases := []test_interface.HttpTestCase[any, any]{
		{
			Description: "Refresh token exist should return a 200 OK response",
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
			},
			BeforeCall: func(req *http.Request) {
				session, hashedRefreshToken := s.userSessionFactory.CreateUserSession(factory.CreateUserSessionParam{
					UserUUID:         s.testUser.UUID.String(),
					HashRefreshToken: false,
				})
				unhashedRefreshToken := session.RefreshToken
				session.RefreshToken = hashedRefreshToken
				session, err := s.userSessionSeeder.SaveUserSession(*session)
				s.Require().NoError(err)
				req.AddCookie(&http.Cookie{
					Name:     refreshTokenCookieName,
					Value:    unhashedRefreshToken,
					Path:     "/user/token",
					Domain:   "localhost",
					MaxAge:   86400,
					HttpOnly: true,
					Secure:   true,
				})
			},
		},
		{
			Description: "Refresh token not exist should return a 400 Bad Request response",
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode:       test_interface.NullInt{Int: http.StatusBadRequest, Valid: true},
				ResponseBodyStatus:       test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyErrorCode:    test_interface.NullString{String: "REFRESH_TOKEN_ALREADY_EXPIRED", Valid: true},
				ResponseBodyErrorMessage: test_interface.NullString{String: "Your session is already expired.", Valid: true},
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			if _, err := s.userSeeder.SaveUser(*s.testUserHashedPassword); err != nil {
				s.T().Log("Failed to create test user before call")
			}
			defer func() {
				if err := s.userSeeder.DeleteUser(s.testUser.Username); err != nil {
					s.T().Log("Failed to delete test user before call")
				}
			}()

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
