package user_rest_test

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
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"
	"github.com/harmonify/movie-reservation-system/user-service/internal"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/entity"
	entityfactory "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/factory"
	entityseeder "github.com/harmonify/movie-reservation-system/user-service/internal/core/entity/seeder"
	token_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/token"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/seeder"
	http_driver "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http"
	user_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/user"
	"github.com/harmonify/movie-reservation-system/user-service/internal/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"go.uber.org/fx"
)

func TestUserRest(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}
	os.Setenv("ENV", config_pkg.EnvironmentTest)
	suite.Run(t, new(UserRestTestSuite))
}

type UserRestTestSuite struct {
	suite.Suite
	app                      *fx.App
	database                 *database.Database
	util                     *util.Util
	httpServer               *http_driver.HttpServer
	tokenService             token_service.TokenService
	otpCacheV2               shared.OtpCacheV2
	userStorage              shared.UserStorage
	userKeyStorage           shared.UserKeyStorage
	userFactory              entityfactory.UserFactory
	userSeeder               entityseeder.UserSeeder
	notificationProviderMock *mocks.NotificationProvider
}

func (s *UserRestTestSuite) SetupSuite() {
	s.notificationProviderMock = mocks.NewNotificationProvider(s.T())

	s.app = internal.NewApp(
		fx.Decorate(func(np shared.NotificationProvider) shared.NotificationProvider {
			return s.notificationProviderMock
		}),
		seeder.DrivenPostgresqlSeederModule,
		entityfactory.UserEntityFactoryModule,
		fx.Invoke(func(
			database *database.Database,
			util *util.Util,
			httpServer *http_driver.HttpServer,
			tokenService token_service.TokenService,
			otpCacheV2 shared.OtpCacheV2,
			userStorage shared.UserStorage,
			userKeyStorage shared.UserKeyStorage,
			userSeeder entityseeder.UserSeeder,
			userFactory entityfactory.UserFactory,
		) {
			s.database = database
			s.util = util
			s.httpServer = httpServer
			s.tokenService = tokenService
			s.otpCacheV2 = otpCacheV2
			s.userStorage = userStorage
			s.userKeyStorage = userKeyStorage
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

func (s *UserRestTestSuite) TestUserRest_GetUser() {
	var (
		PATH   = "/v1/profile"
		METHOD = "GET"
	)

	ctx := context.Background()

	testUser, err := s.userSeeder.CreateUser(ctx)
	s.Require().Nil(err, "Failed to create test user before call")
	defer func() {
		if err := s.userSeeder.DeleteUser(ctx, entity.GetUser{UUID: sql.NullString{String: testUser.User.UUID, Valid: true}}); err != nil {
			s.T().Fatal("Failed to delete test user after call")
		}
	}()

	testCases := []test_interface.HttpTestCase[any, *user_rest.GetUserRes]{
		{
			Description: "It should return a 200 OK response",
			Config: test_interface.Request[any]{
				RequestHeader: []test_interface.RequestHeaderConfig{
					{Key: "Authorization", Value: fmt.Sprintf("Bearer %s", s.generateAccessToken(ctx, testUser.User.UUID))},
				},
			},
			Expectation: test_interface.ResponseExpectation[*user_rest.GetUserRes]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
				ResponseBodyResult: &user_rest.GetUserRes{
					UUID:                  testUser.User.UUID,
					Username:              testUser.User.Username,
					Email:                 testUser.User.Email,
					PhoneNumber:           testUser.User.PhoneNumber,
					FirstName:             testUser.User.FirstName,
					LastName:              testUser.User.LastName,
					IsEmailVerified:       testUser.User.IsEmailVerified,
					IsPhoneNumberVerified: testUser.User.IsPhoneNumberVerified,
					CreatedAt:             testUser.User.CreatedAt,
					UpdatedAt:             testUser.User.UpdatedAt,
					DeletedAt:             nil,
				},
			},
		},
		{
			Description: "It should return a 401 Unauthorized response",
			Config: test_interface.Request[any]{
				RequestHeader: []test_interface.RequestHeaderConfig{
					{Key: "Authorization", Value: "Bearer invalidtoken"},
				},
			},
			Expectation: test_interface.ResponseExpectation[*user_rest.GetUserRes]{
				ResponseStatusCode:      test_interface.NullInt{Int: http.StatusUnauthorized, Valid: true},
				ResponseBodyStatus:      test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:      nil,
				ResponseBodyErrorCode:   test_interface.NullString{String: "INVALID_JWT_ERROR", Valid: true},
				ResponseBodyErrorObject: nil,
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			var err error

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
			// s.T().Log(body)
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

func (s *UserRestTestSuite) TestUserRest_PatchUser() {
	var (
		PATH   = "/v1/profile"
		METHOD = "PATCH"
	)

	ctx := context.Background()

	testUser, err := s.userSeeder.CreateUser(ctx)
	s.Require().Nil(err, "Failed to create test user before call")
	defer func() {
		if err := s.userSeeder.DeleteUser(ctx, entity.GetUser{UUID: sql.NullString{String: testUser.User.UUID, Valid: true}}); err != nil {
			s.T().Fatal("Failed to delete test user after call")
		}
	}()

	accessToken := s.generateAccessToken(ctx, testUser.User.UUID)

	userUpdate, _, err := s.userFactory.GenerateUser()
	s.Require().Nil(err, "Failed to generate user update before call")

	testCases := []test_interface.HttpTestCase[*user_rest.PatchUserReq, func(actualBody string) *user_rest.PatchUserRes]{
		{
			Description: "It should return a 200 OK response",
			Config: test_interface.Request[*user_rest.PatchUserReq]{
				RequestBody: &user_rest.PatchUserReq{
					Username:    userUpdate.Username,
					Email:       userUpdate.Email,
					PhoneNumber: userUpdate.PhoneNumber,
					FirstName:   userUpdate.FirstName,
					LastName:    userUpdate.LastName,
				},
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[func(actualBody string) *user_rest.PatchUserRes]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
				ResponseBodyResult: func(actualBody string) *user_rest.PatchUserRes {
					return &user_rest.PatchUserRes{
						UUID:                  testUser.User.UUID,
						Username:              userUpdate.Username,
						Email:                 userUpdate.Email,
						PhoneNumber:           userUpdate.PhoneNumber,
						FirstName:             userUpdate.FirstName,
						LastName:              userUpdate.LastName,
						IsEmailVerified:       false,
						IsPhoneNumberVerified: false,
						CreatedAt:             testUser.User.CreatedAt,
						UpdatedAt:             gjson.Parse(actualBody).Get("result.updated_at").Time(),
						DeletedAt:             nil,
					}
				},
			},
		},
		{
			Description: "It should return a 400 Bad Request response",
			Config: test_interface.Request[*user_rest.PatchUserReq]{
				RequestBody: &user_rest.PatchUserReq{
					Email:       "invalidemail",
					PhoneNumber: "invalidphonenumber",
					Username:    "in",
					FirstName:   "John3",
					LastName:    "Doe3",
				},
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[func(actualBody string) *user_rest.PatchUserRes]{
				ResponseStatusCode:    test_interface.NullInt{Int: http.StatusBadRequest, Valid: true},
				ResponseBodyStatus:    test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:    nil,
				ResponseBodyErrorCode: test_interface.NullString{String: "INVALID_REQUEST_BODY_ERROR", Valid: true},
				ResponseBodyErrorObject: []interface{}{
					validation.ValidationError{
						Field:   "email",
						Message: "Email must be a valid email address",
					},
					validation.ValidationError{
						Field:   "phone_number",
						Message: "PhoneNumber must be a valid E.164 formatted phone number",
					},
					validation.ValidationError{
						Field:   "username",
						Message: "Username must be at least 3 characters in length",
					},
					validation.ValidationError{
						Field:   "first_name",
						Message: "FirstName can only contain alphabetic characters",
					},
					validation.ValidationError{
						Field:   "last_name",
						Message: "LastName can only contain alphabetic characters",
					},
				},
			},
		},
		{
			Description: "It should return a 401 Unauthorized response",
			Config: test_interface.Request[*user_rest.PatchUserReq]{
				RequestBody: &user_rest.PatchUserReq{
					Email:       userUpdate.Email,
					PhoneNumber: userUpdate.PhoneNumber,
					Username:    userUpdate.Username,
					FirstName:   userUpdate.FirstName,
					LastName:    userUpdate.LastName,
				},
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: "Bearer invalidtoken",
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[func(actualBody string) *user_rest.PatchUserRes]{
				ResponseStatusCode:      test_interface.NullInt{Int: http.StatusUnauthorized, Valid: true},
				ResponseBodyStatus:      test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:      nil,
				ResponseBodyErrorCode:   test_interface.NullString{String: "INVALID_JWT_ERROR", Valid: true},
				ResponseBodyErrorObject: nil,
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			var err error

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
			// s.T().Log(body)
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
				expected := testCase.Expectation.ResponseBodyResult(bodyString)
				expectedBytes, err := json.Marshal(expected)
				s.Assert().NoError(err)
				s.Assert().JSONEq(string(expectedBytes), resultBody.Raw)
			}
			if testCase.Expectation.ResponseBodyErrorCode.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorCode.String, responseError.Get("code").String())
			}
			if testCase.Expectation.ResponseBodyErrorMessage.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorMessage.String, responseError.Get("message").String())
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

func (s *UserRestTestSuite) TestUserRest_SendVerificationEmail() {
	var (
		PATH   = "/v1/profile/email/verification"
		METHOD = "GET"
	)

	ctx := context.Background()

	testUser, err := s.userSeeder.CreateUser(ctx)
	s.Require().Nil(err, "Failed to create test user before call")
	defer func() {
		if err := s.userSeeder.DeleteUser(ctx, entity.GetUser{UUID: sql.NullString{String: testUser.User.UUID, Valid: true}}); err != nil {
			s.T().Fatal("Failed to delete test user after call")
		}
	}()
	s.T().Log(testUser)

	accessToken := s.generateAccessToken(ctx, testUser.User.UUID)

	testCases := []test_interface.HttpTestCase[any, any]{
		{
			Description: "It should return a 200 OK response",
			Config: test_interface.Request[any]{
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
				ResponseBodyResult: nil,
			},
			BeforeCall: func(req *http.Request) {
				s.notificationProviderMock.EXPECT().SendEmail(mock.Anything, mock.Anything).Return(nil).Once()
			},
			AfterCall: func(w *httptest.ResponseRecorder) {
				// Ensure to delete the email verification code after call, the current function is async
				_, err := s.otpCacheV2.DeleteOtp(ctx, testUser.User.UUID, shared.EmailVerificationOtpType)
				s.Require().Nil(err, "Failed to delete email verification code after call")
			},
		},
		{
			Description: "It should return a 401 Unauthorized response",
			Config: test_interface.Request[any]{
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: "Bearer invalidtoken",
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode:      test_interface.NullInt{Int: http.StatusUnauthorized, Valid: true},
				ResponseBodyStatus:      test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:      nil,
				ResponseBodyErrorCode:   test_interface.NullString{String: "INVALID_JWT_ERROR", Valid: true},
				ResponseBodyErrorObject: nil,
			},
		},
		{
			Description: "It should return a 502 Bad Gateway response",
			Config: test_interface.Request[any]{
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode:       test_interface.NullInt{Int: http.StatusBadGateway, Valid: true},
				ResponseBodyStatus:       test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:       nil,
				ResponseBodyErrorCode:    test_interface.NullString{String: "SEND_VERIFICATION_LINK_FAILED", Valid: true},
				ResponseBodyErrorMessage: test_interface.NullString{String: "Failed to send a verification link to your email. If issue persists, please contact our technical support and try again later", Valid: true},
				ResponseBodyErrorObject:  nil,
			},
			BeforeCall: func(req *http.Request) {
				s.notificationProviderMock.EXPECT().SendEmail(mock.Anything, mock.Anything).Return(fmt.Errorf("failed to send email verification link")).Once()
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			var err error

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
			// s.T().Log(body)
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
				expected := testCase.Expectation.ResponseBodyResult
				expectedBytes, err := json.Marshal(expected)
				s.Assert().NoError(err)
				s.Assert().JSONEq(string(expectedBytes), resultBody.Raw)
			}
			if testCase.Expectation.ResponseBodyErrorCode.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorCode.String, responseError.Get("code").String())
			}
			if testCase.Expectation.ResponseBodyErrorMessage.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorMessage.String, responseError.Get("message").String())
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

func (s *UserRestTestSuite) TestUserRest_VerifyEmail() {
	var (
		PATH   = "/v1/profile/email/verification"
		METHOD = "POST"
	)

	ctx := context.Background()

	testUser, err := s.userSeeder.CreateUser(ctx)
	s.Require().Nil(err, "Failed to create test user before call")
	defer func() {
		if err := s.userSeeder.DeleteUser(ctx, entity.GetUser{UUID: sql.NullString{String: testUser.User.UUID, Valid: true}}); err != nil {
			s.T().Fatal("Failed to delete test user after call")
		}
	}()

	accessToken := s.generateAccessToken(ctx, testUser.User.UUID)
	code := s.generateEmailVerificationCode(ctx, testUser.User.UUID)

	testCases := []test_interface.HttpTestCase[user_rest.VerifyEmailReq, any]{
		{
			Description: "It should return a 200 OK response",
			Config: test_interface.Request[user_rest.VerifyEmailReq]{
				RequestBody: user_rest.VerifyEmailReq{
					Code: code,
				},
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
				ResponseBodyResult: nil,
			},
		},
		{
			Description: "It should return a 400 Bad Request response",
			Config: test_interface.Request[user_rest.VerifyEmailReq]{
				RequestBody: user_rest.VerifyEmailReq{
					Code: "",
				},
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode:       test_interface.NullInt{Int: http.StatusBadRequest, Valid: true},
				ResponseBodyStatus:       test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:       nil,
				ResponseBodyErrorCode:    test_interface.NullString{String: "INVALID_REQUEST_BODY_ERROR", Valid: true},
				ResponseBodyErrorMessage: test_interface.NullString{String: "Please ensure you have filled all the required information correctly and try again. If the problem persists, please contact our technical support.", Valid: true},
				ResponseBodyErrorObject: []interface{}{
					validation.ValidationError{
						Field:   "code",
						Message: "Code is a required field",
					},
				},
			},
		},
		{
			Description: "It should return a 401 Unauthorized response",
			Config: test_interface.Request[user_rest.VerifyEmailReq]{
				RequestBody: user_rest.VerifyEmailReq{
					Code: code,
				},
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: "Bearer invalidtoken",
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode:      test_interface.NullInt{Int: http.StatusUnauthorized, Valid: true},
				ResponseBodyStatus:      test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:      nil,
				ResponseBodyErrorCode:   test_interface.NullString{String: "INVALID_JWT_ERROR", Valid: true},
				ResponseBodyErrorObject: nil,
			},
		},
		{
			Description: "It should return a 404 Not Found response",
			Config: test_interface.Request[user_rest.VerifyEmailReq]{
				RequestBody: user_rest.VerifyEmailReq{
					Code: "invalidcode",
				},
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode:       test_interface.NullInt{Int: http.StatusNotFound, Valid: true},
				ResponseBodyStatus:       test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:       nil,
				ResponseBodyErrorCode:    test_interface.NullString{String: "VERIFICATION_TOKEN_NOT_FOUND", Valid: true},
				ResponseBodyErrorMessage: test_interface.NullString{String: "Your verification link may be expired. Please try to request a new verification link.", Valid: true},
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			var err error

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
			// s.T().Log(body)
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
				expected := testCase.Expectation.ResponseBodyResult
				expectedBytes, err := json.Marshal(expected)
				s.Assert().NoError(err)
				s.Assert().JSONEq(string(expectedBytes), resultBody.Raw)
			}
			if testCase.Expectation.ResponseBodyErrorCode.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorCode.String, responseError.Get("code").String())
			}
			if testCase.Expectation.ResponseBodyErrorMessage.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorMessage.String, responseError.Get("message").String())
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

func (s *UserRestTestSuite) TestUserRest_SendPhoneNumberOtp() {
	var (
		PATH   = "/v1/profile/phone/verification"
		METHOD = "GET"
	)

	ctx := context.Background()

	testUser, err := s.userSeeder.CreateUser(ctx)
	s.Require().Nil(err, "Failed to create test user before call")
	defer func() {
		if err := s.userSeeder.DeleteUser(ctx, entity.GetUser{UUID: sql.NullString{String: testUser.User.UUID, Valid: true}}); err != nil {
			s.T().Fatal("Failed to delete test user after call")
		}
	}()

	accessToken := s.generateAccessToken(ctx, testUser.User.UUID)

	testCases := []test_interface.HttpTestCase[any, any]{
		{
			Description: "It should return a 200 OK response",
			Config: test_interface.Request[any]{
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
				ResponseBodyResult: nil,
			},
			BeforeCall: func(req *http.Request) {
				s.notificationProviderMock.EXPECT().SendSms(mock.Anything, mock.Anything).Return(nil).Once()
			},
			AfterCall: func(w *httptest.ResponseRecorder) {
				// Ensure to delete the otp verification attempt and code after call, the current function is async
				_, err := s.otpCacheV2.DeleteOtp(ctx, testUser.User.UUID, shared.PhoneNumberVerificationOtpType)
				s.Require().Nil(err, "Failed to delete phone otp verification attempt after call")
			},
		},
		{
			Description: "It should return a 401 Unauthorized response",
			Config: test_interface.Request[any]{
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: "Bearer invalidtoken",
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode:      test_interface.NullInt{Int: http.StatusUnauthorized, Valid: true},
				ResponseBodyStatus:      test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:      nil,
				ResponseBodyErrorCode:   test_interface.NullString{String: "INVALID_JWT_ERROR", Valid: true},
				ResponseBodyErrorObject: nil,
			},
		},
		{
			Description: "It should return a 502 Bad Gateway response",
			Config: test_interface.Request[any]{
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			BeforeCall: func(req *http.Request) {
				s.notificationProviderMock.EXPECT().SendSms(mock.Anything, mock.Anything).Return(fmt.Errorf("failed to send otp")).Once()
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode:       test_interface.NullInt{Int: http.StatusBadGateway, Valid: true},
				ResponseBodyStatus:       test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:       nil,
				ResponseBodyErrorCode:    test_interface.NullString{String: "SEND_OTP_FAILED", Valid: true},
				ResponseBodyErrorMessage: test_interface.NullString{String: "Failed to send an OTP to your phone number. If issue persists, please contact our technical support and try again later", Valid: true},
				ResponseBodyErrorObject:  nil,
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			var err error

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

			s.Require().True(gjson.Valid(bodyString), fmt.Sprintf("response body should be a valid JSON, but got %s", bodyString))
			body := gjson.Parse(bodyString)
			// s.T().Log(body)

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
				expected := testCase.Expectation.ResponseBodyResult
				expectedBytes, err := json.Marshal(expected)
				s.Assert().NoError(err)
				s.Assert().JSONEq(string(expectedBytes), resultBody.Raw)
			}
			if testCase.Expectation.ResponseBodyErrorCode.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorCode.String, responseError.Get("code").String())
			}
			if testCase.Expectation.ResponseBodyErrorMessage.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorMessage.String, responseError.Get("message").String())
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

func (s *UserRestTestSuite) TestUserRest_VerifyPhoneNumber() {
	var (
		PATH   = "/v1/profile/phone/verification"
		METHOD = "POST"
	)

	ctx := context.Background()

	testUser, err := s.userSeeder.CreateUser(ctx)
	s.Require().Nil(err, "Failed to create test user before call")
	defer func() {
		if err := s.userSeeder.DeleteUser(ctx, entity.GetUser{UUID: sql.NullString{String: testUser.User.UUID, Valid: true}}); err != nil {
			s.T().Fatal("Failed to delete test user after call")
		}
	}()

	accessToken := s.generateAccessToken(ctx, testUser.User.UUID)

	otp := s.generatePhoneNumberVerificationOtp(ctx, testUser.User.UUID)
	var nonexistentOtp string
	for otp == nonexistentOtp || nonexistentOtp == "" && err == nil {
		nonexistentOtp, err = s.util.GeneratorUtil.GenerateRandomNumber(6)
		s.Require().Nil(err, "Failed to generate random number before call")
	}

	testCases := []test_interface.HttpTestCase[user_rest.VerifyPhoneNumberReq, any]{
		{
			Description: "It should return a 200 OK response",
			Config: test_interface.Request[user_rest.VerifyPhoneNumberReq]{
				RequestBody: user_rest.VerifyPhoneNumberReq{
					Otp: otp,
				},
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode: test_interface.NullInt{Int: http.StatusOK, Valid: true},
				ResponseBodyStatus: test_interface.NullBool{Bool: true, Valid: true},
				ResponseBodyResult: nil,
			},
		},
		{
			Description: "It should return a 400 Bad Request response",
			Config: test_interface.Request[user_rest.VerifyPhoneNumberReq]{
				RequestBody: user_rest.VerifyPhoneNumberReq{
					Otp: "",
				},
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode:       test_interface.NullInt{Int: http.StatusBadRequest, Valid: true},
				ResponseBodyStatus:       test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:       nil,
				ResponseBodyErrorCode:    test_interface.NullString{String: "INVALID_REQUEST_BODY_ERROR", Valid: true},
				ResponseBodyErrorMessage: test_interface.NullString{String: "Please ensure you have filled all the required information correctly and try again. If the problem persists, please contact our technical support.", Valid: true},
				ResponseBodyErrorObject: []interface{}{
					validation.ValidationError{
						Field:   "otp",
						Message: "Otp is a required field",
					},
				},
			},
		},
		{
			Description: "It should return a 401 Unauthorized response",
			Config: test_interface.Request[user_rest.VerifyPhoneNumberReq]{
				RequestBody: user_rest.VerifyPhoneNumberReq{
					Otp: otp,
				},
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: "Bearer invalidtoken",
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode:      test_interface.NullInt{Int: http.StatusUnauthorized, Valid: true},
				ResponseBodyStatus:      test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:      nil,
				ResponseBodyErrorCode:   test_interface.NullString{String: "INVALID_JWT_ERROR", Valid: true},
				ResponseBodyErrorObject: nil,
			},
		},
		{
			Description: "It should return a 404 Not Found response",
			Config: test_interface.Request[user_rest.VerifyPhoneNumberReq]{
				RequestBody: user_rest.VerifyPhoneNumberReq{
					Otp: nonexistentOtp,
				},
				RequestHeader: []test_interface.RequestHeaderConfig{
					{
						Key:   "Authorization",
						Value: fmt.Sprintf("Bearer %s", accessToken),
					},
				},
			},
			Expectation: test_interface.ResponseExpectation[any]{
				ResponseStatusCode:       test_interface.NullInt{Int: http.StatusNotFound, Valid: true},
				ResponseBodyStatus:       test_interface.NullBool{Bool: false, Valid: true},
				ResponseBodyResult:       nil,
				ResponseBodyErrorCode:    test_interface.NullString{String: "OTP_NOT_FOUND", Valid: true},
				ResponseBodyErrorMessage: test_interface.NullString{String: "Your OTP may be expired. Please try to request a new OTP.", Valid: true},
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			var err error

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
			// s.T().Log(body)
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
				expected := testCase.Expectation.ResponseBodyResult
				expectedBytes, err := json.Marshal(expected)
				s.Assert().NoError(err)
				s.Assert().JSONEq(string(expectedBytes), resultBody.Raw)
			}
			if testCase.Expectation.ResponseBodyErrorCode.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorCode.String, responseError.Get("code").String())
			}
			if testCase.Expectation.ResponseBodyErrorMessage.Valid {
				s.Assert().Equal(testCase.Expectation.ResponseBodyErrorMessage.String, responseError.Get("message").String())
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

func (s *UserRestTestSuite) generateAccessToken(ctx context.Context, uuid string) string {
	// Get user record
	user, err := s.userStorage.GetUser(ctx, entity.GetUser{UUID: sql.NullString{String: uuid, Valid: true}})
	if err != nil {
		s.T().Fatal("Failed to Get user", err)
		return ""
	}

	// Get user key record
	userKey, err := s.userKeyStorage.GetUserKey(ctx, entity.GetUserKey{
		UserUUID: sql.NullString{String: user.UUID, Valid: true},
	})
	if err != nil {
		s.T().Fatal("Failed to Get user key", err)
		return ""
	}

	// Generate access token
	accessToken, err := s.tokenService.GenerateAccessToken(ctx, token_service.GenerateAccessTokenParam{
		UUID:        user.UUID,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		PrivateKey:  userKey.PrivateKey,
	})
	if err != nil {
		s.T().Fatal("Failed to generate access token", err)
		return ""
	}

	if accessToken.AccessToken == "" {
		s.T().Fatal("Failed to generate access token, empty access token")
		return ""
	}

	return accessToken.AccessToken
}

func (s *UserRestTestSuite) generateEmailVerificationCode(ctx context.Context, uuid string) string {
	code, err := s.util.GeneratorUtil.GenerateRandomHex(32)
	if err != nil {
		s.T().Fatal("Failed to generate random number", err)
		return ""
	}

	user, err := s.userStorage.GetUser(ctx, entity.GetUser{UUID: sql.NullString{String: uuid, Valid: true}})
	if err != nil {
		s.T().Fatal("Failed to Get user", err)
		return ""
	}

	err = s.otpCacheV2.SaveOtp(ctx, user.UUID, shared.EmailVerificationOtpType, code)
	if err != nil {
		s.T().Fatal("Failed to save email verification code", err)
		return ""
	}

	return code
}

func (s *UserRestTestSuite) generatePhoneNumberVerificationOtp(ctx context.Context, uuid string) string {
	code, err := s.util.GeneratorUtil.GenerateRandomNumber(6)
	if err != nil {
		s.T().Fatal("Failed to generate random number", err)
		return ""
	}

	user, err := s.userStorage.GetUser(ctx, entity.GetUser{UUID: sql.NullString{String: uuid, Valid: true}})
	if err != nil {
		s.T().Fatal("Failed to Get user", err)
		return ""
	}

	err = s.otpCacheV2.SaveOtp(ctx, user.UUID, shared.PhoneNumberVerificationOtpType, code)
	if err != nil {
		s.T().Fatal("Failed to save phone number verification OTP", err)
		return ""
	}

	return code
}
