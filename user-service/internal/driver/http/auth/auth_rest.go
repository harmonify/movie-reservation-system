package auth_rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
	http_interface "github.com/harmonify/movie-reservation-system/user-service/lib/http/interface"
	middleware "github.com/harmonify/movie-reservation-system/user-service/lib/http/middleware"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/response"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/validator"
	http_validator "github.com/harmonify/movie-reservation-system/user-service/lib/http/validator"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
)

type AuthRestHandler interface {
	Register(g *gin.RouterGroup)
}

type authRestHandlerImpl struct {
	authService auth_service.AuthService
	validator   http_validator.HttpValidator
	response    response.HttpResponse
	tracer      tracer.Tracer
	util        *util.Util
	middleware  middleware.JWTMiddleware
}

func NewAuthRestHandler(
	authService auth_service.AuthService,
	tracer tracer.Tracer,
	response response.HttpResponse,
	validator validator.HttpValidator,
	util *util.Util,
	middleware middleware.JWTMiddleware,
) AuthRestHandler {
	return &authRestHandlerImpl{
		authService: authService,
		tracer:      tracer,
		response:    response,
		validator:   validator,
		util:        util,
		middleware:  middleware,
	}
}

func (h *authRestHandlerImpl) Register(g *gin.RouterGroup) {
	g.POST("/register", h.PostUserRegistration)
	// TODO: email verification & connect to email service
	// g.GET("/email/verify", h.middleware.JWTMiddleware.AuthenticateUser, h.GetEmailVerificationLink)
	// g.POST("/email/verify", h.PostVerifyEmail)
	// g.GET("/email/edit/verify", h.middleware.JWTMiddleware.AuthenticateUser, h.GetChangeEmailVerificationLink)
	// g.POST("/email/edit/verify", h.PostVerifyChangeEmail)
	g.POST("/login", h.PostUserLogin)
	g.POST("/logout", h.middleware.AuthenticateUser, h.PostUserLogout)
	g.GET("/token", h.GetRefreshToken)
}

func (h *authRestHandlerImpl) PostUserRegistration(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		body PostUserRegisterReq
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.Validate(c, &body); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	headers := http_interface.HeadersExtension{
		UserAgent: c.Request.UserAgent(),
		IpAddress: h.util.HttpUtil.GetUserIP(c.Request),
	}

	res, err = h.authService.Register(ctx, &auth_service.RegisterParam{
		Username:         body.Username,
		Password:         body.Password,
		Email:            body.Email,
		PhoneNumber:      body.PhoneNumber,
		FullName:         body.FullName,
		HeadersExtension: headers,
	})

	if err != nil {
		if errors.Is(err, constant.ErrUnauthorized) {
			err = h.response.BuildError(constant.Unauthorized, err)
		} else if errors.Is(err, constant.ErrServiceUnavailable) {
			err = h.response.BuildError(constant.ServiceUnavailable, err)
		} else if errors.Is(err, auth_service.ErrDuplicateEmail) {
			err = h.response.BuildError(auth_service.DuplicateEmail, err)
		} else if errors.Is(err, auth_service.ErrDuplicatePhoneNumber) {
			err = h.response.BuildError(auth_service.DuplicatePhoneNumber, err)
		} else if errors.Is(err, auth_service.ErrInvalidPhoneNumber) {
			err = h.response.BuildError(auth_service.InvalidPhoneNumber, err)
		} else if errors.Is(err, constant.ErrRateLimitExceeded) {
			err = h.response.BuildError(constant.RateLimitExceeded, err)
		} else if errors.Is(err, constant.ErrInvalidJwt) {
			err = h.response.BuildError(constant.InvalidJwt, err)
		} else {
			err = h.response.BuildError(constant.InternalServerError, err)
		}

		h.response.Send(c, nil, err)
		return
	}

	h.response.Send(c, nil, nil)
}

func (h *authRestHandlerImpl) PostVerifyEmail(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		body PostVerifyEmailReq
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.Validate(c, &body); err != nil {
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	header := http_interface.VerifierHeader{
		TokenVerifier: c.Request.Header.Get("x-token-verifier"),
	}

	err = h.authService.VerifyEmailUser(ctx, &header, &body)
	if err != nil {
		if errors.Is(err, constant.RecordNotFoundError) {
			err = h.errorResponse.ThrowError(constant.EmailNotRegistered, err)
		} else if errors.Is(err, constant.EmailVerificationLinkExpiredError) {
			err = h.errorResponse.ThrowError(constant.EmailVerificationLinkExpired, err)
		} else if errors.Is(err, constant.InvalidEmailError) {
			err = h.errorResponse.ThrowError(constant.InvalidEmail, err)
		} else {
			err = h.errorResponse.ThrowError(constant.InternalServerError, err)
		}
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *authRestHandlerImpl) GetEmailVerificationLink(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	userInfo := c.Request.Context().Value(middleware.UserInfoKey).(*jwt.JWTBodyPayload)

	err = h.authService.GetEmailVerificationLink(ctx, userInfo.Email)
	if err != nil {
		if errors.Is(err, constant.RateLimitExceededError) {
			err = h.errorResponse.ThrowError(constant.RateLimitExceeded, err)
		} else {
			err = h.errorResponse.ThrowError(constant.InternalServerError, err)
		}
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *authRestHandlerImpl) GetChangeEmailVerificationLink(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		params PostVerifyEmailReq
		err    error
	)

	userInfo := c.Request.Context().Value(middleware.UserInfoKey).(*jwt.JWTBodyPayload)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.ValidateQueryParams(c, &params); err != nil {
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	err = h.authService.GetChangeEmailVerificationLink(ctx, userInfo.Email, params.Email)
	if err != nil {
		if errors.Is(err, constant.RateLimitExceededError) {
			err = h.errorResponse.ThrowError(constant.RateLimitExceeded, err)
		} else {
			err = h.errorResponse.ThrowError(constant.InternalServerError, err)
		}
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *authRestHandlerImpl) PostVerifyChangeEmail(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		body PostVerifyEmailReq
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.Validate(c, &body); err != nil {
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	header := http_interface.VerifierHeader{
		TokenVerifier: c.Request.Header.Get("x-token-verifier"),
	}

	err = h.authService.VerifyChangeEmailUser(ctx, &header, &body)
	if err != nil {
		if errors.Is(err, constant.RecordNotFoundError) {
			err = h.errorResponse.ThrowError(constant.EmailNotRegistered, err)
		} else if errors.Is(err, constant.EmailVerificationLinkExpiredError) {
			err = h.errorResponse.ThrowError(constant.EmailVerificationLinkExpired, err)
		} else if errors.Is(err, constant.InvalidEmailError) {
			err = h.errorResponse.ThrowError(constant.InvalidEmail, err)
		} else {
			err = h.errorResponse.ThrowError(constant.InternalServerError, err)
		}
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *authRestHandlerImpl) PostUserLogin(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		params PostUserLoginReq
		err    error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.Validate(c, &params); err != nil {
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	headers := http_interface.HeadersExtension{
		UserAgent: c.Request.UserAgent(),
		IpAddress: h.util.MiscellaneousUtil.GetUserIP(c.Request),
	}

	data, err := h.authService.UserAuthentication(ctx, &UserAuthenticationReq{
		PostUserLoginReq: params,
		Attributes:       params.Attributes,
		HeadersExtension: headers,
	})
	if err != nil {
		if errors.Is(err, constant.InvalidCredentialError) {
			err = h.errorResponse.ThrowError(constant.InvalidCredential, err)
		} else if errors.Is(err, constant.ServiceUnavailableError) {
			err = h.errorResponse.ThrowError(constant.ServiceUnavailable, err)
		} else if errors.Is(err, constant.UserExistError) {
			err = h.errorResponse.ThrowError(constant.UserExist, err)
		} else {
			err = h.errorResponse.ThrowError(constant.InternalServerError, err)
		}
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	// Set refresh token cookies
	cookiesName := constant.DefaultAppCookie + "token"
	cookiesValue := data.RefreshToken
	// TODO: limit the cookies domain
	cookiesDomain := "*localhost"
	cookiesPath := "/user/token"
	c.SetCookie(cookiesName, cookiesValue, 0, cookiesPath, cookiesDomain, true, true)

	// Set token verifier for 2FA
	if data.TokenVerifier != "" {
		c.Header("x-token-verifier", data.TokenVerifier)
	}

	mapper := mappers.MapperPostUserLogin(data)
	c.JSON(h.response.Send(ctx, mapper, err))
}

func (h *authRestHandlerImpl) PostUserLogout(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cookiesName := "refreshToken"
	refreshToken, err := c.Cookie(cookiesName)
	if err != nil {
		err = h.errorResponse.ThrowError(constant.InvalidToken, err)
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	accessToken := c.Request.Header.Get("Authorization")
	splitAccessToken := strings.Split(accessToken, " ")
	if len(splitAccessToken) == 2 {
		if splitAccessToken[1] == "" {
			err = h.errorResponse.ThrowError(constant.Unauthorized, nil)
			c.JSON(h.response.Send(ctx, nil, err))
			return
		}
		accessToken = splitAccessToken[1]
	} else {
		err = h.errorResponse.ThrowError(constant.Unauthorized, nil)
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	err = h.authService.UserLogout(ctx, &HandleRefreshTokenReq{
		UserToken: http_interface.UserToken{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}})
	if err != nil {
		if errors.Is(err, constant.InvalidCredentialError) {
			err = h.errorResponse.ThrowError(constant.InvalidCredential, err)
		} else if errors.Is(err, constant.ServiceUnavailableError) {
			err = h.errorResponse.ThrowError(constant.ServiceUnavailable, err)
		} else if errors.Is(err, constant.UserExistError) {
			err = h.errorResponse.ThrowError(constant.UserExist, err)
		} else {
			err = h.errorResponse.ThrowError(constant.InternalServerError, err)
		}
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *authRestHandlerImpl) GetRefreshToken(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		err = h.errorResponse.ThrowError(constant.InvalidToken, err)
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	accessToken := c.Request.Header.Get("Authorization")
	splitAccessToken := strings.Split(accessToken, " ")
	if len(splitAccessToken) == 2 {
		if splitAccessToken[1] == "" {
			err = h.errorResponse.ThrowError(constant.Unauthorized, nil)
			c.JSON(h.response.Send(ctx, nil, err))
			return
		}
		accessToken = splitAccessToken[1]
	} else {
		err = h.errorResponse.ThrowError(constant.Unauthorized, nil)
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	headers := http_interface.HeadersExtension{
		UserAgent: c.Request.UserAgent(),
		IpAddress: h.util.MiscellaneousUtil.GetUserIP(c.Request),
	}

	data, err := h.authService.HandleRefreshToken(ctx, &HandleRefreshTokenReq{
		UserToken: http_interface.UserToken{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		HeadersExtension: headers,
	})
	if err != nil {
		if errors.Is(err, constant.InvalidTokenError) {
			err = h.errorResponse.ThrowError(constant.InvalidToken, err)
		} else {
			err = h.errorResponse.ThrowError(constant.InternalServerError, err)
		}
		c.JSON(h.response.Send(ctx, nil, err))
		return
	}

	// Set refresh token cookies
	cookiesName := constant.DefaultAppCookie + "token"
	cookiesValue := data.RefreshToken
	cookiesDomain := "*populix.co"
	cookiesPath := "/user/token"
	c.SetCookie(cookiesName, cookiesValue, 0, cookiesPath, cookiesDomain, true, true)
	c.JSON(h.response.Send(ctx, data, err))
}
