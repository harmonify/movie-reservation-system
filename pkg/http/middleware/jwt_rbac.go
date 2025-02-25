package http_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type (
	JwtRbacHttpMiddleware interface {
		CheckPermissions(func(userInfo *jwt_util.JWTBodyPayload) (authorized bool)) gin.HandlerFunc
	}

	JwtRbacHttpMiddlewareParam struct {
		fx.In

		Logger   logger.Logger
		Tracer   tracer.Tracer
		Response http_pkg.HttpResponse
	}

	JwtHttpMiddlewareConfig struct {
		Domain string `validate:"required"` // THe service domain, used for RBAC
	}

	JwtRbacHttpMiddlewareResult struct {
		fx.Out

		JwtRbacHttpMiddleware JwtRbacHttpMiddleware
	}

	jwtRbacHttpMiddlewareImpl struct {
		logger   logger.Logger
		tracer   tracer.Tracer
		response http_pkg.HttpResponse
		config   *JwtHttpMiddlewareConfig
	}
)

func NewJwtRbacHttpMiddleware(p JwtRbacHttpMiddlewareParam, cfg *JwtHttpMiddlewareConfig) (JwtRbacHttpMiddlewareResult, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return JwtRbacHttpMiddlewareResult{}, err
	}
	return JwtRbacHttpMiddlewareResult{
		JwtRbacHttpMiddleware: &jwtRbacHttpMiddlewareImpl{
			logger:   p.Logger,
			tracer:   p.Tracer,
			response: p.Response,
			config:   cfg,
		},
	}, nil
}

func (h *jwtRbacHttpMiddlewareImpl) CheckPermissions(fn func(*jwt_util.JWTBodyPayload) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.StartSpanWithCaller(c.Request.Context())
		defer span.End()

		_userInfo := c.Request.Context().Value(http_pkg.UserInfoKey)
		if _userInfo == nil {
			h.logger.WithCtx(ctx).Error("failed to get user info from context")
			h.response.Send(c, nil, error_pkg.InternalServerError)
			c.Abort()
			return
		}

		userInfo, ok := _userInfo.(*jwt_util.JWTBodyPayload)
		if !ok {
			h.logger.WithCtx(ctx).Error("invalid user info from context")
			h.response.Send(c, nil, error_pkg.InternalServerError)
			c.Abort()
			return
		}

		authorized := fn(userInfo)
		if !authorized {
			h.logger.WithCtx(ctx).Debug(
				"user is forbidden to access this resource",
				zap.String("user_uuid", userInfo.UUID),
				zap.String("domain", h.config.Domain),
				zap.String("resource", c.Request.URL.Path),
				zap.String("action", c.Request.Method),
			)
			h.response.Send(c, nil, error_pkg.ForbiddenError)
		}

		c.Next()
	}
}
