package user_http_middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	auth_proto "github.com/harmonify/movie-reservation-system/pkg/proto/auth"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthHttpMiddleware interface {
	// Default returns the authentication middleware
	Default() gin.HandlerFunc
	// WithPolicy returns the authentication + authorization middleware
	WithPolicy(policy_id string) gin.HandlerFunc
}

type AuthHttpMiddlewareParam struct {
	fx.In
	tracer.Tracer
	logger.Logger
	error_pkg.ErrorMapper
	http_pkg.HttpResponseBuilder
	jwt_util.JwtUtil
	auth_proto.AuthServiceClient
}

type AuthHttpMiddlewareResult struct {
	fx.Out
	AuthHttpMiddleware AuthHttpMiddleware
}

type authHttpMiddlewareImpl struct {
	tracer            tracer.Tracer
	logger            logger.Logger
	errorMapper       error_pkg.ErrorMapper
	responseBuilder   http_pkg.HttpResponseBuilder
	jwtUtil           jwt_util.JwtUtil
	authServiceClient auth_proto.AuthServiceClient

	defaultMiddleware gin.HandlerFunc
}

func NewAuthHttpMiddleware(p AuthHttpMiddlewareParam) AuthHttpMiddlewareResult {
	m := &authHttpMiddlewareImpl{
		tracer:            p.Tracer,
		logger:            p.Logger,
		errorMapper:       p.ErrorMapper,
		responseBuilder:   p.HttpResponseBuilder,
		jwtUtil:           p.JwtUtil,
		authServiceClient: p.AuthServiceClient,
	}

	m.defaultMiddleware = m.build(nil)

	return AuthHttpMiddlewareResult{
		AuthHttpMiddleware: m,
	}
}

func (a *authHttpMiddlewareImpl) Default() gin.HandlerFunc {
	return a.defaultMiddleware
}

func (a authHttpMiddlewareImpl) WithPolicy(policyId string) gin.HandlerFunc {
	return a.build(&policyId)
}

func (a *authHttpMiddlewareImpl) build(policyId *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := a.tracer.StartSpanWithCaller(c.Request.Context())
		defer span.End()

		a.logger.WithCtx(ctx).Debug("auth header", zap.String("authorization_header", c.Request.Header.Get("Authorization")))

		accessToken, err := a.extractAuthorizationHeader(ctx, c)
		if err != nil {
			a.responseBuilder.New().WithCtx(ctx).WithError(err).Send(c)
			c.Abort()
			return
		}

		res, err := a.authServiceClient.Auth(ctx, &auth_proto.AuthRequest{
			AccessToken: accessToken,
			PolicyId:    policyId,
		})
		if err != nil {
			terr, valid := a.errorMapper.FromGrpcError(err)
			if !valid {
				a.logger.WithCtx(ctx).Error("uncatched grpc error", zap.Error(terr))
			}
			a.responseBuilder.New().WithCtx(ctx).WithError(terr).Send(c)
			c.Abort()
			return
		}

		ctxWithUser := context.WithValue(c.Request.Context(), http_pkg.UserInfoKey, res.GetUserInfo())
		c.Request = c.Request.WithContext(ctxWithUser)
		c.Next()
	}
}

func (a *authHttpMiddlewareImpl) extractAuthorizationHeader(ctx context.Context, c *gin.Context) (string, error) {
	ctx, span := a.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	authHeader := c.Request.Header.Get("Authorization")
	splitAccessToken := strings.Split(authHeader, " ")
	if len(splitAccessToken) != 2 || splitAccessToken[1] == "" {
		a.logger.WithCtx(ctx).Debug("invalid authorization header")
		return "", error_pkg.InvalidAuthorizationHeaderError
	}

	return splitAccessToken[1], nil
}
