package test

import (
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	user_rest "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/user"
)

type (
	GetUserTestExpectation struct {
		ResponseStatusCode       test_interface.NullInt
		ResponseBodyStatus       test_interface.NullBool
		ResponseBodyResult       user_rest.GetUserRes
		ResponseBodyErrorCode    test_interface.NullString
		ResponseBodyErrorMessage test_interface.NullString
		ResponseBodyErrorObject  []interface{}
	}

	PatchUserTestConfig func() user_rest.PatchUserReq

	PatchUserTestExpectation struct {
		ResponseStatusCode       test_interface.NullInt
		ResponseBodyStatus       test_interface.NullBool
		ResponseBodyResult       user_rest.PatchUserRes
		ResponseBodyErrorCode    test_interface.NullString
		ResponseBodyErrorMessage test_interface.NullString
		ResponseBodyErrorObject  []interface{}
	}

	SendVerificationEmailTestConfig struct {
		// Data user_rest.SendVerificationEmailReq
	}

	VerifyEmailTestConfig struct {
		Data user_rest.VerifyEmailReq
	}

	SendPhoneNumberVerificationTestConfig struct {
		// Data user_rest.SendPhoneNumberVerificationReq
	}

	VerifyPhoneNumberTestConfig struct {
		Data user_rest.VerifyPhoneNumberReq
	}
)
