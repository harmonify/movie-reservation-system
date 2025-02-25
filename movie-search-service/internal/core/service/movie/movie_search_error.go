package movie_service

import (
	"net/http"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"google.golang.org/grpc/codes"
)

var (
	InvalidCursorError = &error_pkg.ErrorWithDetails{
		Code:     error_pkg.ErrorCode("INVALID_CURSOR_ERROR"),
		HttpCode: http.StatusBadRequest,
		GrpcCode: codes.InvalidArgument,
		Message: "Invalid cursor provided",
	}
)
