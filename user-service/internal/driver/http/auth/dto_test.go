package auth_rest_test

import (
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	auth_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/auth"
)

type (
	postRegisterTestConfig struct {
		Data auth_rest.PostRegisterReq
	}

	postRegisterTestExpectation struct {
		ResponseStatusCode                   test_interface.NullInt
		ResponseBodyStatus                   test_interface.NullBool
		ResponseBodyErrorCode                test_interface.NullString
		ResponseBodyErrorMessage             test_interface.NullString
		ResponseBodyErrorObject              []interface{}
		ResponseHeaderRefreshTokenExist      test_interface.NullBool
		ResponseBodyAccessTokenExist         test_interface.NullBool
		ResponseBodyAccessTokenDurationExist test_interface.NullBool
	}

	postVerifyEmailTestConfig func() auth_rest.PostVerifyEmailReq

	postVerifyEmailTestExpectation struct {
		ResponseStatusCode       test_interface.NullInt
		ResponseBodyStatus       test_interface.NullBool
		ResponseBodyResult       interface{}
		ResponseBodyErrorCode    test_interface.NullString
		ResponseBodyErrorMessage test_interface.NullString
		ResponseBodyErrorObject  []interface{}
		IsEmailVerified          test_interface.NullBool
	}

	postLoginTestConfig struct {
		Data auth_rest.PostRegisterReq
	}

	postLoginTestExpectation struct {
		Result auth_rest.PostLoginRes
	}

	getTokenTestExpectation struct {
		AccessTokenExist         test_interface.NullBool
		AccessTokenDurationExist test_interface.NullBool
	}
)
