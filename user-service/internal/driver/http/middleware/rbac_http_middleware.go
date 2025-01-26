package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	"github.com/harmonify/movie-reservation-system/user-service/internal/core/shared"
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
		Response    http_pkg.HttpResponse
		RbacStorage shared.RbacStorage
		Config      *config.Config
	}

	RbacHttpMiddlewareResult struct {
		fx.Out

		RbacHttpMiddleware RbacHttpMiddleware
	}

	rbacHttpMiddlewareImpl struct {
		logger      logger.Logger
		tracer      tracer.Tracer
		response    http_pkg.HttpResponse
		rbacStorage shared.RbacStorage
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
	_, span := m.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

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
	ctx, span := m.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	var userInfo *jwt_util.JWTBodyPayload
	_userInfo := c.Request.Context().Value(UserInfoKey)
	if _userInfo == nil {
		return false, error_pkg.UnauthorizedError
	} else {
		userInfo = _userInfo.(*jwt_util.JWTBodyPayload)
	}

	allowed, err := m.rbacStorage.CheckPermission(ctx, shared.CheckPermissionParam{
		UUID:     userInfo.UUID,
		Domain:   m.config.ServiceIdentifier,
		Resource: c.Request.URL.Path,
		Action:   shared.Action(c.Request.Method),
	})

	if err != nil {
		return false, err
	}

	return allowed, nil
}
