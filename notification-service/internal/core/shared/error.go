package shared

import (
	"fmt"
	"net/http"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"google.golang.org/grpc/codes"
)

var (
	EmptyRecipientError = &error_pkg.ErrorWithDetails{
		Code:     "EMPTY_RECIPIENT",
		Message:  "recipient is empty",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
	}

	EmptySubjectError = &error_pkg.ErrorWithDetails{
		Code:     "EMPTY_SUBJECT",
		Message:  "subject is empty",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
	}

	EmptyTemplateError = &error_pkg.ErrorWithDetails{
		Code:     "EMPTY_TEMPLATE",
		Message:  "template is empty",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
	}

	EmptyBodyError = &error_pkg.ErrorWithDetails{
		Code:     "EMPTY_BODY",
		Message:  "body is empty",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
	}

	InvalidTemplateIdError = &error_pkg.ErrorWithDetails{
		Code:     "INVALID_TEMPLATE_ID",
		Message:  "invalid template id",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
	}

	InvalidTemplateDataError = &error_pkg.ErrorWithDetails{
		Code:     "INVALID_TEMPLATE_DATA",
		Message:  "invalid template data",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
	}

	InvalidPhoneNumberError = &error_pkg.ErrorWithDetails{
		Code:     "INVALID_PHONE_NUMBER",
		Message:  "invalid phone number",
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
	}
)

type InvalidPhoneNumberErrorData struct {
	PhoneNumber string `json:"phone_number"`
}

func NewInvalidPhoneNumberError(phoneNumber string) *error_pkg.ErrorWithDetails {
	return &error_pkg.ErrorWithDetails{
		Code:     InvalidPhoneNumberError.Code,
		Message:  fmt.Sprintf(InvalidPhoneNumberError.Message+": %s", phoneNumber),
		HttpCode: InvalidPhoneNumberError.HttpCode,
		GrpcCode: InvalidPhoneNumberError.GrpcCode,
		Data:     &InvalidPhoneNumberErrorData{PhoneNumber: phoneNumber},
	}
}
