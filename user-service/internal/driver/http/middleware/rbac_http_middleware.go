package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	shared_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/shared"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	http_constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/response"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	jwt_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/jwt"
	"go.uber.org/fx"
)

type (
	RbacHttpMiddleware interface {
		CheckPermission(*gin.Context)
	}

	RbacHttpMiddlewareParam struct {
		fx.In

		Logger      logger.Logger
		Tracer      tracer.Tracer
		Response    response.HttpResponse
		RbacStorage shared_service.RbacStorage
		Config      *config.Config
	}

	RbacHttpMiddlewareResult struct {
		fx.Out

		RbacHttpMiddleware RbacHttpMiddleware
	}

	rbacHttpMiddlewareImpl struct {
		logger      logger.Logger
		tracer      tracer.Tracer
		response    response.HttpResponse
		rbacStorage shared_service.RbacStorage
		config      *config.Config
	}
)

func NewRbacHttpMiddleware(p RbacHttpMiddlewareParam) RbacHttpMiddlewareResult {
	return RbacHttpMiddlewareResult{
		RbacHttpMiddleware: &rbacHttpMiddlewareImpl{
			logger:      p.Logger,
			tracer:      p.Tracer,
			response:    p.Response,
			rbacStorage: p.RbacStorage,
			config:      p.Config,
		},
	}
}

// CheckPermission checks the user/domain/method/path combination from the request.
func (m *rbacHttpMiddlewareImpl) CheckPermission(c *gin.Context) {
	authorized, err := m.checkPermission(c)

	if err != nil {
		m.response.Send(c, nil, err)
	}

	if !authorized {
		m.response.SendWithResponseCode(c, http.StatusForbidden, nil, nil)
	}

	c.Next()
}

func (m *rbacHttpMiddlewareImpl) checkPermission(c *gin.Context) (bool, error) {
	var (
		ctx = c.Request.Context()
	)

	_, span := m.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	var userInfo *jwt_util.JWTBodyPayload
	_userInfo := c.Request.Context().Value(UserInfoKey)
	if _userInfo == nil {
		return false, http_constant.ErrUnauthorized
	} else {
		userInfo = _userInfo.(*jwt_util.JWTBodyPayload)
	}

	allowed, err := m.rbacStorage.CheckPermission(ctx, shared_service.CheckPermissionParam{
		UUID:     userInfo.UUID,
		Domain:   m.config.ServiceIdentifier,
		Resource: c.Request.URL.Path,
		Action:   shared_service.Action(c.Request.Method),
	})

	if err != nil {
		return false, err
	}

	return allowed, nil
}
