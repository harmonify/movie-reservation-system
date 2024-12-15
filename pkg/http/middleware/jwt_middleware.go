package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	constant "github.com/harmonify/movie-reservation-system/pkg/http/constant"
	"github.com/harmonify/movie-reservation-system/pkg/http/response"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
)

type JWTMiddleware interface {
	AuthenticateUser(c *gin.Context)
	OptAuthenticateUser(c *gin.Context)
}

type jwtHttpMiddlewareImpl struct {
	tracer   tracer.Tracer
	response response.HttpResponse
	util     *util.Util
}

// suppress golangci-lint
// SA1029: should not use built-in type string as key for value; define your own type to avoid collisions
type contextKey string

// UserInfoKey should be used when access userInfo
const UserInfoKey contextKey = "userInfo"

func NewJWTMiddleware(
	tracer tracer.Tracer,
	response response.HttpResponse,
	util *util.Util,
) JWTMiddleware {
	return &jwtHttpMiddlewareImpl{
		tracer:   tracer,
		response: response,
		util:     util,
	}
}

func (h *jwtHttpMiddlewareImpl) AuthenticateUser(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	_, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	accessToken := c.Request.Header.Get("Authorization")
	splitAccessToken := strings.Split(accessToken, " ")
	if len(splitAccessToken) == 2 {
		if splitAccessToken[1] == "" {
			err := h.response.BuildError(constant.Unauthorized, nil)
			h.response.Send(c, nil, err)
			c.Abort()
			return
		}
		accessToken = splitAccessToken[1]
	} else {
		err := h.response.BuildError(constant.Unauthorized, nil)
		h.response.Send(c, nil, err)
		c.Abort()
		return
	}

	payload, err := h.util.JWTUtil.JWTVerify(accessToken)
	if err != nil {
		if errors.Is(err, constant.ErrUnauthorized) {
			err = h.response.BuildError(constant.Unauthorized, err)
		} else {
			err = h.response.BuildError(constant.InternalServerError, err)
		}

		h.response.Send(c, nil, err)
		c.Abort()
		return
	}

	ctx = context.WithValue(c.Request.Context(), UserInfoKey, payload)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func (h *jwtHttpMiddlewareImpl) OptAuthenticateUser(c *gin.Context) {
	var (
		ctx = c.Request.Context()
	)

	_, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	accessToken := c.Request.Header.Get("Authorization")

	if accessToken != "" {
		h.AuthenticateUser(c)
	}
}
