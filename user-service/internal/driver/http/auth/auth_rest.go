package auth_rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/middleware"
	error_constant "github.com/harmonify/movie-reservation-system/user-service/lib/error/constant"
	constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/response"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/validator"
	http_validator "github.com/harmonify/movie-reservation-system/user-service/lib/http/validator"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
)

type AuthRestHandler interface {
	Register(g *gin.RouterGroup)
	Version() string
}

type authRestHandlerImpl struct {
	authService auth_service.AuthService
	validator   http_validator.HttpValidator
	response    response.HttpResponse
	tracer      tracer.Tracer
	util        *util.Util
	middleware  *middleware.HttpMiddleware
}

func NewAuthRestHandler(
	authService auth_service.AuthService,
	tracer tracer.Tracer,
	response response.HttpResponse,
	validator validator.HttpValidator,
	util *util.Util,
	middleware *middleware.HttpMiddleware,
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
	g.POST("/register", h.PostRegister)
	g.POST("/register/verify", h.PostVerifyEmail)
	g.POST("/login", h.PostLogin)
	g.POST("/logout", h.middleware.JwtHttpMiddleware.AuthenticateUser, h.PostLogout)
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

	if err = h.validator.Validate(c, &body); err != nil {
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

	if err = h.validator.Validate(c, &body); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	err = h.authService.VerifyEmail(ctx, auth_service.VerifyEmailParam{
		Email: body.Email,
	})
	if err != nil {
		h.response.Send(c, nil, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *authRestHandlerImpl) PostLogin(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		params PostLoginReq
		err    error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.Validate(c, &params); err != nil {
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
	cookieName := constant.HttpCookiePrefix + "token"
	cookieMaxAge := int(time.Until(data.RefreshTokenExpiredAt).Seconds())
	cookieValue := data.RefreshToken
	cookieDomain := "*localhost"
	cookiePath := "/user/token"
	c.SetCookie(cookieName, cookieValue, cookieMaxAge, cookiePath, cookieDomain, true, true)

	h.response.Send(c, PostLoginRes{
		AccessToken:         data.AccessToken,
		AccessTokenDuration: data.AccessTokenDuration,
	}, err)
}

func (h *authRestHandlerImpl) GetToken(c *gin.Context) {
	var (
		ctx    = c.Request.Context()
		params GetTokenReq
		err    error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.Validate(c, &params); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	data, err := h.authService.GetToken(ctx, auth_service.GetTokenParam{
		RefreshToken: params.RefreshToken,
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

	cookieName := constant.HttpCookiePrefix + "token"
	refreshToken, err := c.Cookie(cookieName)
	if err != nil {
		err = h.response.BuildError(error_constant.InvalidJwt, err)
		h.response.Send(c, nil, err)
		return
	}

	err = h.authService.Logout(ctx, auth_service.LogoutParam{
		RefreshToken: refreshToken,
	})

	// Delete refresh token cookie
	cookieValue := ""
	cookieMaxAge := -1
	cookieDomain := "*localhost"
	cookiePath := "/user/token"
	c.SetCookie(cookieName, cookieValue, cookieMaxAge, cookiePath, cookieDomain, true, true)

	h.response.Send(c, nil, err)
}
