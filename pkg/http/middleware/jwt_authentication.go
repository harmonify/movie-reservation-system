package http_middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	JwtAuthenticationHttpMiddleware interface {
		AuthenticateUser(c *gin.Context)    // Authenticate user
		OptAuthenticateUser(c *gin.Context) // Optionally authenticate user
	}

	JwtHttpMiddlewareParam struct {
		fx.In
		Logger   logger.Logger
		Tracer   tracer.Tracer
		Response http_pkg.HttpResponse
		Util     *util.Util
	}

	JwtHttpMiddlewareResult struct {
		fx.Out
		JwtHttpMiddleware JwtAuthenticationHttpMiddleware
	}

	jwtHttpMiddlewareImpl struct {
		logger   logger.Logger
		tracer   tracer.Tracer
		response http_pkg.HttpResponse
		util     *util.Util
	}
)

func NewJwtHttpMiddleware(
	p JwtHttpMiddlewareParam,
) JwtHttpMiddlewareResult {
	return JwtHttpMiddlewareResult{
		JwtHttpMiddleware: &jwtHttpMiddlewareImpl{
			logger:   p.Logger,
			tracer:   p.Tracer,
			response: p.Response,
			util:     p.Util,
		},
	}
}

func (h *jwtHttpMiddlewareImpl) AuthenticateUser(c *gin.Context) {
	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	h.logger.WithCtx(ctx).Debug("auth header", zap.String("authorization_header", c.Request.Header.Get("Authorization")))

	payload, err := h.verifyAuthorizationHeader(ctx, c)
	if err != nil {
		h.response.Send(c, nil, err)
		c.Abort()
		return
	}

	ctxWithUser := context.WithValue(c.Request.Context(), http_pkg.UserInfoKey, payload)
	c.Request = c.Request.WithContext(ctxWithUser)
	c.Next()
}

func (h *jwtHttpMiddlewareImpl) verifyAuthorizationHeader(ctx context.Context, c *gin.Context) (*jwt_util.JWTBodyPayload, error) {
	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	authHeader := c.Request.Header.Get("Authorization")
	splitAccessToken := strings.Split(authHeader, " ")
	if len(splitAccessToken) != 2 || splitAccessToken[1] == "" {
		h.logger.WithCtx(ctx).Debug("invalid authorization header")
		return nil, error_pkg.InvalidAuthorizationHeaderError
	}

	accessToken := splitAccessToken[1]

	payload, err := h.util.JWTUtil.JWTVerify(ctx, accessToken)
	if err != nil {
		h.logger.WithCtx(ctx).Debug("failed to verify jwt", zap.Error(err))
		return nil, err
	}

	return payload, nil
}

func (h *jwtHttpMiddlewareImpl) OptAuthenticateUser(c *gin.Context) {
	ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	authHeader := c.Request.Header.Get("Authorization")

	h.logger.WithCtx(ctx).Debug("auth header", zap.String("authorization_header", authHeader))

	if authHeader != "" {
		payload, err := h.verifyAuthorizationHeader(ctx, c)
		if err != nil {
			h.response.Send(c, nil, err)
			c.Abort()
			return
		}
		ctxWithUser := context.WithValue(c.Request.Context(), http_pkg.UserInfoKey, payload)
		c.Request = c.Request.WithContext(ctxWithUser)
	}

	c.Next()
}
