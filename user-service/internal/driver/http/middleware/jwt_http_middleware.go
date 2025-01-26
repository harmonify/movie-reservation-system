package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
)

type (
	JwtHttpMiddleware interface {
		AuthenticateUser(c *gin.Context)
		OptAuthenticateUser(c *gin.Context)
	}

	jwtHttpMiddlewareImpl struct {
		tracer   tracer.Tracer
		response http_pkg.HttpResponse
		util     *util.Util
	}
)

func NewJwtHttpMiddleware(
	tracer tracer.Tracer,
	response http_pkg.HttpResponse,
	util *util.Util,
) JwtHttpMiddleware {
	return &jwtHttpMiddlewareImpl{
		tracer:   tracer,
		response: response,
		util:     util,
	}
}

func (h *jwtHttpMiddlewareImpl) AuthenticateUser(c *gin.Context) {
	var (
		err error
	)

	_, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	accessToken := c.Request.Header.Get("Authorization")
	splitAccessToken := strings.Split(accessToken, " ")
	if len(splitAccessToken) == 2 {
		if splitAccessToken[1] == "" {
			h.response.Send(c, nil, error_pkg.UnauthorizedError)
			c.Abort()
			return
		}
		accessToken = splitAccessToken[1]
	} else {
		h.response.Send(c, nil, error_pkg.UnauthorizedError)
		c.Abort()
		return
	}

	payload, err := h.util.JWTUtil.JWTVerify(accessToken)
	if err != nil {
		h.response.Send(c, nil, err)
		c.Abort()
		return
	}

	ctx := context.WithValue(c.Request.Context(), UserInfoKey, payload)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func (h *jwtHttpMiddlewareImpl) OptAuthenticateUser(c *gin.Context) {
	_, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	accessToken := c.Request.Header.Get("Authorization")

	if accessToken != "" {
		h.AuthenticateUser(c)
	}
}
