package auth_rest

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	http_constant "github.com/harmonify/movie-reservation-system/pkg/http/constant"
	http_interface "github.com/harmonify/movie-reservation-system/pkg/http/interface"
	"github.com/harmonify/movie-reservation-system/pkg/http/response"
	http_validator "github.com/harmonify/movie-reservation-system/pkg/http/validator"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthRestHandler interface {
	Register(g *gin.RouterGroup)
	Version() string
}

type AuthRestHandlerParam struct {
	fx.In

	AuthService auth_service.AuthService
	Validator   http_validator.HttpValidator
	Response    response.HttpResponse
	Logger      logger.Logger
	Tracer      tracer.Tracer
	Util        *util.Util
	Middleware  *middleware.HttpMiddleware
}

type AuthRestHandlerResult struct {
	fx.Out

	AuthRestHandler http_interface.RestHandler `group:"http_routes"`
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

func NewAuthRestHandler(p AuthRestHandlerParam) AuthRestHandlerResult {
	return AuthRestHandlerResult{
		AuthRestHandler: &authRestHandlerImpl{
			authService: p.AuthService,
			logger:      p.Logger,
			tracer:      p.Tracer,
			response:    p.Response,
			validator:   p.Validator,
			util:        p.Util,
			middleware:  p.Middleware,
		},
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
