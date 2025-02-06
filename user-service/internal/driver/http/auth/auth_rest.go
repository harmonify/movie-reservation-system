package auth_rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/ratelimiter"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	auth_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/auth"
	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/config"
	http_driver_shared "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/shared"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthRestHandlerParam struct {
	fx.In

	AuthService auth_service.AuthService
	Validator   http_pkg.HttpValidator
	Response    http_pkg.HttpResponse
	Config      *config.UserServiceConfig
	Logger      logger.Logger
	Tracer      tracer.Tracer
	Util        *util.Util
	Middleware  *http_driver_shared.HttpMiddleware
}

type AuthRestHandlerResult struct {
	fx.Out

	AuthRestHandler http_pkg.RestHandler `group:"http_routes"`
}

type authRestHandlerImpl struct {
	authService auth_service.AuthService
	validator   http_pkg.HttpValidator
	response    http_pkg.HttpResponse
	config      *config.UserServiceConfig
	logger      logger.Logger
	tracer      tracer.Tracer
	util        *util.Util
	middleware  *http_driver_shared.HttpMiddleware
}

func NewAuthRestHandler(p AuthRestHandlerParam) AuthRestHandlerResult {
	return AuthRestHandlerResult{
		AuthRestHandler: &authRestHandlerImpl{
			authService: p.AuthService,
			config:      p.Config,
			logger:      p.Logger,
			tracer:      p.Tracer,
			response:    p.Response,
			validator:   p.Validator,
			util:        p.Util,
			middleware:  p.Middleware,
		},
	}
}

func (h *authRestHandlerImpl) Register(g *gin.RouterGroup) error {
	var registerIPCap int64 = 10
	if h.config.Env == config_pkg.EnvironmentDevelopment {
		registerIPCap = 100
	}

	var loginCap int64 = 10
	if h.config.Env == config_pkg.EnvironmentDevelopment {
		loginCap = 100
	}

	var getTokenCap int64 = 10
	if h.config.Env == config_pkg.EnvironmentDevelopment {
		getTokenCap = 100
	}

	g.POST(
		"/register",
		h.middleware.RateLimiter.LimitByIP(&ratelimiter.RateLimiterConfig{
			Capacity:   registerIPCap,
			RefillRate: time.Minute,
		}),
		h.postRegister,
	)
	g.POST(
		"/login",
		h.middleware.RateLimiter.LimitBy(
			&ratelimiter.RateLimiterConfig{
				Capacity:   loginCap,
				RefillRate: time.Minute,
			},
			func(c *gin.Context) string {
				var params PostLoginReq
				if err := h.validator.ValidateRequestBody(c, &params); err != nil {
					h.response.Send(c, nil, err)
					c.Abort()
					return ""
				}
				ctxWithBody := context.WithValue(c.Request.Context(), http_pkg.ParsedBodyKey, params)
				c.Request = c.Request.WithContext(ctxWithBody)
				return params.Username
			},
		),
		h.postLogin,
	)
	g.GET(
		"/token",
		h.middleware.RateLimiter.LimitByIP(&ratelimiter.RateLimiterConfig{
			Capacity:   getTokenCap,
			RefillRate: time.Minute,
		}),
		h.getToken,
	)
	g.POST("/logout", h.postLogout)

	return nil
}

func (h *authRestHandlerImpl) Version() string {
	return "1"
}

func (h *authRestHandlerImpl) postRegister(c *gin.Context) {
	var (
		body PostRegisterReq
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
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

func (h *authRestHandlerImpl) postLogin(c *gin.Context) {
	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	params, ok := c.Request.Context().Value(http_pkg.ParsedBodyKey).(PostLoginReq)
	if !ok {
		h.logger.WithCtx(ctx).Error("Failed to assert parsed body")
		h.response.Send(c, nil, error_pkg.InvalidRequestBodyError)
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
	cookieName := http_pkg.HttpCookiePrefix + "token"
	cookieMaxAge := int(time.Until(data.RefreshTokenExpiredAt).Seconds())
	cookieValue := data.RefreshToken
	cookieDomain := "localhost"
	cookiePath := "/user/token"
	h.logger.WithCtx(ctx).Debug("Set refresh token cookie", zap.String("cookieName", cookieName), zap.Int("cookieMaxAge", cookieMaxAge), zap.String("cookieDomain", cookieDomain), zap.String("cookiePath", cookiePath), zap.String("username", params.Username))
	c.SetCookie(cookieName, cookieValue, cookieMaxAge, cookiePath, cookieDomain, true, true)

	h.response.Send(c, &PostLoginRes{
		AccessToken:         data.AccessToken,
		AccessTokenDuration: data.AccessTokenDuration,
	}, nil)
}

func (h *authRestHandlerImpl) getToken(c *gin.Context) {
	var (
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	cookieName := http_pkg.HttpCookiePrefix + "token"
	refreshToken, err := c.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			h.response.Send(c, nil, auth_service.RefreshTokenExpiredError)
			return
		}
		h.logger.WithCtx(ctx).Error("Failed to get refresh token cookie", zap.Error(err))
		h.response.Send(c, nil, err)
		return
	}

	data, err := h.authService.GetToken(ctx, auth_service.GetTokenParam{
		RefreshToken: refreshToken,
	})

	h.response.Send(c, &GetTokenRes{
		AccessToken:         data.AccessToken,
		AccessTokenDuration: data.AccessTokenDuration,
	}, err)
}

func (h *authRestHandlerImpl) postLogout(c *gin.Context) {
	var (
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	cookieName := http_pkg.HttpCookiePrefix + "token"
	refreshToken, err := c.Cookie(cookieName)
	if errors.Is(err, http.ErrNoCookie) {
		h.response.Send(c, nil, nil)
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

	h.response.Send(c, nil, err)
}
