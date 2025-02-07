package user_rest

import (
	"time"

	"github.com/gin-gonic/gin"
	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/ratelimiter"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	otp_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/otp"
	user_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/user"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	http_driver_shared "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/shared"
	"go.uber.org/fx"
)

type UserRestHandlerParam struct {
	fx.In

	Config      *config.UserServiceConfig
	Logger      logger.Logger
	Tracer      tracer.Tracer
	Util        *util.Util
	Middleware  *http_driver_shared.HttpMiddleware
	Validator   http_pkg.HttpValidator
	Response    http_pkg.HttpResponse
	UserService user_service.UserService
	OtpService  otp_service.OtpService
}

type UserRestHandlerResult struct {
	fx.Out

	UserRestHandler http_pkg.RestHandler `group:"http_routes"`
}

type userRestHandlerImpl struct {
	config      *config.UserServiceConfig
	logger      logger.Logger
	tracer      tracer.Tracer
	util        *util.Util
	middleware  *http_driver_shared.HttpMiddleware
	validator   http_pkg.HttpValidator
	response    http_pkg.HttpResponse
	userService user_service.UserService
	otpService  otp_service.OtpService
}

func NewUserRestHandler(p UserRestHandlerParam) UserRestHandlerResult {
	return UserRestHandlerResult{
		UserRestHandler: &userRestHandlerImpl{
			config:      p.Config,
			logger:      p.Logger,
			tracer:      p.Tracer,
			util:        p.Util,
			middleware:  p.Middleware,
			validator:   p.Validator,
			response:    p.Response,
			userService: p.UserService,
			otpService:  p.OtpService,
		},
	}
}

func (h *userRestHandlerImpl) Register(g *gin.RouterGroup) error {
	var getOrUpdateUserCap int64 = 5
	var getOrVerifyEmailCodeCap int64 = 1
	var getOrVerifyOtpCap int64 = 1
	if h.config.Env == config_pkg.EnvironmentDevelopment || h.config.Env == config_pkg.EnvironmentTest {
		getOrUpdateUserCap = 100
		getOrVerifyEmailCodeCap = 100
		getOrVerifyOtpCap = 100
	}

	ug := g.Group("/profile")

	ug.GET(
		"",
		h.middleware.Auth.AuthenticateUser,
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getOrUpdateUserCap,
			RefillRate: 3 * time.Second,
		}),
		h.getUser,
	)
	ug.PATCH(
		"",
		h.middleware.Auth.AuthenticateUser,
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getOrUpdateUserCap,
			RefillRate: time.Minute,
		}),
		h.patchUser,
	)

	ug.GET(
		"/email/verification",
		h.middleware.Auth.AuthenticateUser,
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getOrVerifyEmailCodeCap,
			RefillRate: time.Minute,
		}),
		h.sendVerificationEmail,
	)
	ug.POST(
		"/email/verification",
		h.middleware.Auth.AuthenticateUser,
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getOrVerifyEmailCodeCap,
			RefillRate: time.Minute,
		}),
		h.verifyEmail,
	)

	ug.GET(
		"/phone/verification",
		h.middleware.Auth.AuthenticateUser,
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getOrVerifyOtpCap,
			RefillRate: time.Minute,
		}),
		h.sendPhoneNumberVerification,
	)
	ug.POST(
		"/phone/verification",
		h.middleware.Auth.AuthenticateUser,
		h.middleware.RateLimiter.LimitByUUID(&ratelimiter.RateLimiterConfig{
			Capacity:   getOrVerifyOtpCap,
			RefillRate: time.Minute,
		}),
		h.verifyPhoneNumber,
	)

	return nil
}

func (h *userRestHandlerImpl) Version() string {
	return "1"
}

func (h *userRestHandlerImpl) getUser(c *gin.Context) {
	var (
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	userInfo, err := h.util.HttpUtil.GetUserInfo(c.Request)
	if err != nil {
		h.response.Send(c, nil, err)
		return
	}

	data, err := h.userService.GetUser(ctx, user_service.GetUserParam{
		UUID: userInfo.UUID,
	})

	h.response.Send(c, &GetUserRes{
		UUID:                  data.UUID,
		Username:              data.Username,
		Email:                 data.Email,
		PhoneNumber:           data.PhoneNumber,
		FirstName:             data.FirstName,
		LastName:              data.LastName,
		IsEmailVerified:       data.IsEmailVerified,
		IsPhoneNumberVerified: data.IsPhoneNumberVerified,
		CreatedAt:             data.CreatedAt,
		UpdatedAt:             data.UpdatedAt,
		DeletedAt:             data.DeletedAt,
	}, err)
}

func (h *userRestHandlerImpl) patchUser(c *gin.Context) {
	var (
		body PatchUserReq
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	userInfo, err := h.util.HttpUtil.GetUserInfo(c.Request)
	if err != nil {
		h.response.Send(c, nil, err)
		return
	}

	data, err := h.userService.UpdateUser(ctx, user_service.UpdateUserParam{
		UUID:        userInfo.UUID,
		Username:    body.Username,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		PhoneNumber: body.PhoneNumber,
	})

	h.response.Send(c, &PatchUserRes{
		UUID:                  data.UUID,
		Username:              data.Username,
		Email:                 data.Email,
		PhoneNumber:           data.PhoneNumber,
		FirstName:             data.FirstName,
		LastName:              data.LastName,
		IsEmailVerified:       data.IsEmailVerified,
		IsPhoneNumberVerified: data.IsPhoneNumberVerified,
		CreatedAt:             data.CreatedAt,
		UpdatedAt:             data.UpdatedAt,
		DeletedAt:             data.DeletedAt,
	}, err)
}

func (h *userRestHandlerImpl) sendVerificationEmail(c *gin.Context) {
	var (
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	userInfo, err := h.util.HttpUtil.GetUserInfo(c.Request)
	if err != nil {
		h.response.Send(c, nil, err)
		return
	}

	err = h.otpService.SendVerificationEmail(ctx, otp_service.SendVerificationEmailParam{
		UUID: userInfo.UUID,
	})

	h.response.Send(c, nil, err)
}

func (h *userRestHandlerImpl) verifyEmail(c *gin.Context) {
	var (
		body VerifyEmailReq
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	userInfo, err := h.util.HttpUtil.GetUserInfo(c.Request)
	if err != nil {
		h.response.Send(c, nil, err)
		return
	}

	err = h.otpService.VerifyEmail(ctx, otp_service.VerifyEmailParam{
		UUID: userInfo.UUID,
		Code: body.Code,
	})

	h.response.Send(c, nil, err)
}

func (h *userRestHandlerImpl) sendPhoneNumberVerification(c *gin.Context) {
	var (
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	userInfo, err := h.util.HttpUtil.GetUserInfo(c.Request)
	if err != nil {
		h.response.Send(c, nil, err)
		return
	}

	err = h.otpService.SendPhoneNumberVerificationOtp(ctx, otp_service.SendPhoneNumberVerificationOtpParam{
		UUID: userInfo.UUID,
	})

	h.response.Send(c, nil, err)
}

func (h *userRestHandlerImpl) verifyPhoneNumber(c *gin.Context) {
	var (
		body VerifyPhoneNumberReq
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	userInfo, err := h.util.HttpUtil.GetUserInfo(c.Request)
	if err != nil {
		h.response.Send(c, nil, err)
		return
	}

	err = h.otpService.VerifyPhoneNumber(ctx, otp_service.VerifyPhoneNumberParam{
		UUID: userInfo.UUID,
		Otp:  body.Otp,
	})

	h.response.Send(c, nil, err)
}
