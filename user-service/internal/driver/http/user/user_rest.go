package user_rest

import (
	"errors"

	"github.com/gin-gonic/gin"
	constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/response"
	http_validator "github.com/harmonify/movie-reservation-system/user-service/lib/http/validator"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	jwt_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/jwt"

	user_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/user"
	middleware "github.com/harmonify/movie-reservation-system/user-service/lib/http/middleware"
)

type UserRestHandler interface {
	Register(g *gin.RouterGroup)
	GetUser(c *gin.Context)
	PatchUser(c *gin.Context)
}

type userRestHandlerImpl struct {
	userService user_service.UserService
	validator   http_validator.HttpValidator
	response    response.HttpResponse
	tracer      tracer.Tracer
	jwtUtil     jwt_util.JWTUtil
	middleware  *middleware.JWTMiddleware
}

func NewUserRestHandler(
	userService user_service.UserService,
	tracer tracer.Tracer,
	response response.HttpResponse,
	validator http_validator.HttpValidator,
	jwtUtil jwt_util.JWTUtil,
	middleware middleware.JWTMiddleware,
) UserRestHandler {
	return &userRestHandlerImpl{
		userService: userService,
		tracer:      tracer,
		response:    response,
		validator:   validator,
		jwtUtil:     jwtUtil,
		middleware:  middleware,
	}
}

func (h *userRestHandlerImpl) GetUser(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		params GetQueryInfoReq
		err    error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	var fUserInfo *jwt_util.JWTBodyPayload
	userInfo := c.Request.Context().Value(middleware.UserInfoKey)
	if userInfo != nil {
		fUserInfo = userInfo.(*jwt_util.JWTBodyPayload)
	}

	if err = h.validator.ValidateQueryParams(c, &params); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	data, err := h.userService.GetUser(ctx, fUserInfo, &params)
	if err != nil {
		if errors.Is(err, constant.ErrRateLimitExceeded) {
			err = h.response.BuildError(constant.RateLimitExceeded, err)
		} else if errors.Is(err, constant.ErrInvalidJwt) {
			err = h.response.BuildError(constant.ErrInvalidJwt, err)
		} else if errors.Is(err, constant.ErrNotFound) {
			err = h.response.BuildError(constant.ErrUnregisteredAccount, err)
		} else {
			err = h.response.BuildError(constant.InternalServerError, err)
		}

		h.response.Send(c, nil, err)
		return
	}

	h.response.Send(c, data, err)
}

func (h *userRestHandlerImpl) PatchUser(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		body PatchUserReq
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.Validate(c, &body); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	userInfo := c.Request.Context().Value(middleware.UserInfoKey).(*jwt_util.JWTBodyPayload)

	data, err := h.userService.EditUser(ctx, &body, userInfo)
	if err != nil {
		if errors.Is(err, constant.InvalidCredentialError) {
			err = h.response.BuildError(constant.InvalidCredential, err)
		} else if errors.Is(err, constant.ChangeLimitEmailError) {
			err = h.response.BuildError(constant.ChangeLimitEmail, err)
		} else if errors.Is(err, constant.ChangeLimitPhoneNumberError) {
			err = h.response.BuildError(constant.ChangeLimitPhoneNumber, err)
		} else {
			err = h.response.BuildError(constant.InternalServerError, err)
		}
		h.response.Send(ctx, nil, err)
		return
	}

	h.response.Send(ctx, data, err)
}
