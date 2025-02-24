package error_pkg

import (
	"context"
	"errors"
	"sync"

	shared_proto "github.com/harmonify/movie-reservation-system/pkg/proto/shared"
	"go.uber.org/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	mu sync.RWMutex
)

type ErrorMapper interface {
	// RegisterErrors registers an error with a custom error.
	// If the error code already exists, it will be ignored.
	RegisterErrors(errs ...*ErrorWithDetails)
	// FromError asserts the provided error to a *DetailedError.
	// If the provided error is nil, it will return nil.
	// If the provided error concrete type of provided error is a *DetailedError,
	// it will return the provided error, otherwise it will return the default error.
	FromError(original error) (err *ErrorWithDetails, valid bool)
	// FromCode gets a detailed error from the provided error code.
	// If the error code does not exist, it will return the default error.
	FromCode(code ErrorCode) (err *ErrorWithDetails, found bool)
	// ToGrpcError converts the provided error to a grpc error.
	// If the provided error is nil, it will return nil.
	// If the provided error is not a detailed error, it will return the default error.
	ToGrpcError(err error) error
	// FromGrpcError converts the provided grpc error to a detailed error.
	// If the provided grpc error is nil, it will return nil.
	// If the provided grpc error is not a detailed error, it will return the default error.
	FromGrpcError(original error) (err *ErrorWithDetails, valid bool)
}

type ErrorMapperParam struct {
	fx.In
	fx.Lifecycle
}

type errorMapperImpl struct {
	errorMap map[ErrorCode]*ErrorWithDetails
}

func NewErrorMapper(p ErrorMapperParam) ErrorMapper {
	m := &errorMapperImpl{
		errorMap: make(map[ErrorCode]*ErrorWithDetails),
	}

	p.Lifecycle.Append(fx.StartHook(func(ctx context.Context) error {
		m.RegisterErrors(
			DefaultError,
			InternalServerError,
			InvalidRequestBodyError,
			UnauthorizedError,
			InvalidJwtError,
			InvalidJwtClaimsError,
			InvalidJwtSigningMethodError,
			ForbiddenError,
			NotFoundError,
			UnprocessableEntityError,
			RateLimitExceededError,
			BadGatewayError,
			ServiceUnavailableError,
		)
		return nil
	}))

	return m
}

func (e *errorMapperImpl) RegisterErrors(errs ...*ErrorWithDetails) {
	mu.Lock()
	defer mu.Unlock()
	for _, err := range errs {
		if err.Code == "" {
			continue
		}
		if _, ok := e.errorMap[err.Code]; ok {
			continue
		}
		e.errorMap[err.Code] = err
	}
}

func (e *errorMapperImpl) FromError(original error) (*ErrorWithDetails, bool) {
	if original == nil {
		return nil, true
	}

	var target *ErrorWithDetails
	if errors.As(original, &target) {
		return target, true
	}
	return DefaultError, false
}

func (e *errorMapperImpl) FromCode(code ErrorCode) (*ErrorWithDetails, bool) {
	mu.RLock()
	defer mu.RUnlock()
	if errorDetails, ok := e.errorMap[code]; ok {
		return errorDetails, true
	} else {
		return DefaultError, false
	}
}

func (e *errorMapperImpl) ToGrpcError(err error) error {
	if err == nil {
		return nil
	}
	if e, valid := e.FromError(err); e != nil && valid {
		st := status.New(e.GrpcCode, e.Message)
		ds, err := st.WithDetails(
			&shared_proto.Error{
				Code:     e.Code.String(),
				Message:  e.Message,
				HttpCode: int32(e.HttpCode),
				GrpcCode: int32(e.GrpcCode),
			},
		)
		if err != nil {
			return st.Err()
		}
		return ds.Err()
	}
	return status.Error(DefaultError.GrpcCode, err.Error())
}

func (e *errorMapperImpl) FromGrpcError(original error) (*ErrorWithDetails, bool) {
	if original == nil {
		return nil, true
	}
	var derr *ErrorWithDetails
	if errors.As(original, &derr) {
		return derr, true
	}

	st := status.Convert(original)
	for _, d := range st.Details() {
		switch info := d.(type) {
		case *shared_proto.Error:
			return &ErrorWithDetails{
				Code:     ErrorCode(info.Code),
				Message:  info.Message,
				HttpCode: int(info.HttpCode),
				GrpcCode: codes.Code(info.GrpcCode),
			}, true
		default:
			continue
		}
	}
	return &ErrorWithDetails{
		Code:     DefaultError.Code,
		Message:  DefaultError.Message,
		HttpCode: DefaultError.HttpCode,
		GrpcCode: st.Code(),
	}, false
}
