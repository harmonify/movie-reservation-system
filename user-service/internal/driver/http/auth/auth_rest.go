package auth_rest

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/middleware"
	http_constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/response"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/validator"
	http_validator "github.com/harmonify/movie-reservation-system/user-service/lib/http/validator"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"go.uber.org/zap"
)

type AuthRestHandler interface {
	Register(g *gin.RouterGroup)
	Version() string
}

type authRestHandlerImpl struct {
	authService auth_service.AuthService
	validator   http_validator.HttpValidator
	response    response.HttpResponse
	logger      logger.Logger
	tracer      tracer.Tracer
	util        *util.Util
	middleware  *middleware.HttpMiddleware
}

func NewAuthRestHandler(
	authService auth_service.AuthService,
	logger logger.Logger,
	tracer tracer.Tracer,
	response response.HttpResponse,
	validator validator.HttpValidator,
	util *util.Util,
	middleware *middleware.HttpMiddleware,
) AuthRestHandler {
	return &authRestHandlerImpl{
		authService: authService,
		logger:      logger,
		tracer:      tracer,
		response:    response,
		validator:   validator,
		util:        util,
		middleware:  middleware,
	}
}

func (h *authRestHandlerImpl) Register(g *gin.RouterGroup) {
	g.POST("/register", h.PostRegister)
	g.POST("/register/verify", h.PostVerifyEmail)
	g.POST("/login", h.PostLogin)
	g.POST("/logout", h.PostLogout)
	g.GET("/token", h.GetToken)
}

func (h *authRestHandlerImpl) Version() string {
	return "1"
}

func (h *authRestHandlerImpl) PostRegister(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		body PostRegisterReq
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	err = h.authService.Register(ctx, auth_service.RegisterParam{
		Username:    body.Username,
		Password:    body.Password,
		Email:       body.Email,
		PhoneNumber: body.PhoneNumber,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
	})

	h.response.Send(c, nil, err)
}

func (h *authRestHandlerImpl) PostVerifyEmail(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		body PostVerifyEmailReq
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	err = h.authService.VerifyEmail(ctx, auth_service.VerifyEmailParam{
		Email: body.Email,
		Token: body.Token,
	})
	if err != nil {
		h.response.Send(c, nil, err)
		return
	}

	h.response.Send(c, nil, nil)
}

func (h *authRestHandlerImpl) PostLogin(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		params PostLoginReq
		err    error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.ValidateRequestBody(c, &params); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	data, err := h.authService.Login(ctx, auth_service.LoginParam{
		Username:  params.Username,
		Password:  params.Password,
		UserAgent: c.Request.UserAgent(),
		IpAddress: h.util.HttpUtil.GetUserIP(c.Request),
	})
	if err != nil {
		h.response.Send(c, nil, err)
		return
	}

	// Set refresh token cookies
	cookieName := http_constant.HttpCookiePrefix + "token"
	cookieMaxAge := int(time.Until(data.RefreshTokenExpiredAt).Seconds())
	cookieValue := data.RefreshToken
	cookieDomain := "localhost"
	cookiePath := "/user/token"
	h.logger.Debug("Set refresh token cookie", zap.String("cookieName", cookieName), zap.Int("cookieMaxAge", cookieMaxAge), zap.String("cookieDomain", cookieDomain), zap.String("cookiePath", cookiePath), zap.String("username", params.Username))
	c.SetCookie(cookieName, cookieValue, cookieMaxAge, cookiePath, cookieDomain, true, true)

	h.response.Send(c, PostLoginRes{
		AccessToken:         data.AccessToken,
		AccessTokenDuration: data.AccessTokenDuration,
	}, err)
}

func (h *authRestHandlerImpl) GetToken(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cookieName := http_constant.HttpCookiePrefix + "token"
	refreshToken, err := c.Cookie(cookieName)
	if err != nil {
		err = h.response.BuildError(auth_service.InvalidRefreshToken, err)
		h.response.Send(c, nil, err)
		return
	}

	data, err := h.authService.GetToken(ctx, auth_service.GetTokenParam{
		RefreshToken: refreshToken,
	})
	if err != nil {
		h.response.Send(c, nil, err)
		return
	}

	h.response.Send(c, GetTokenRes{
		AccessToken:         data.AccessToken,
		AccessTokenDuration: data.AccessTokenDuration,
	}, err)
}

func (h *authRestHandlerImpl) PostLogout(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	cookieName := http_constant.HttpCookiePrefix + "token"
	refreshToken, err := c.Cookie(cookieName)
	if err != nil {
		h.response.Send(c, nil, auth_service.ErrRefreshTokenAlreadyExpired)
		return
	}

	err = h.authService.Logout(ctx, auth_service.LogoutParam{
		RefreshToken: refreshToken,
	})

	// Delete refresh token cookie
	cookieValue := ""
	cookieMaxAge := -1
	cookieDomain := "localhost"
	cookiePath := "/user/token"
	c.SetCookie(cookieName, cookieValue, cookieMaxAge, cookiePath, cookieDomain, true, true)

	if err != nil && !errors.Is(err, auth_service.ErrRefreshTokenAlreadyExpired) {
		h.response.Send(c, nil, err)
	}

	h.response.Send(c, nil, nil)
}
