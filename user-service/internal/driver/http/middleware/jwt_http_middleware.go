package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/response"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
)

type (
	JwtHttpMiddleware interface {
		AuthenticateUser(c *gin.Context)
		OptAuthenticateUser(c *gin.Context)
	}

	jwtHttpMiddlewareImpl struct {
		tracer   tracer.Tracer
		response response.HttpResponse
		util     *util.Util
	}
)

func NewJwtHttpMiddleware(
	tracer tracer.Tracer,
	response response.HttpResponse,
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