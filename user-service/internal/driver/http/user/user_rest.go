package user_rest

import (
	"github.com/gin-gonic/gin"
	http_interface "github.com/harmonify/movie-reservation-system/pkg/http/interface"
	"github.com/harmonify/movie-reservation-system/pkg/http/response"
	http_validator "github.com/harmonify/movie-reservation-system/pkg/http/validator"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
	"go.uber.org/fx"

	user_service "github.com/harmonify/movie-reservation-system/user-service/internal/core/service/user"
	middleware "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/middleware"
)

type UserRestHandler interface {
	Register(g *gin.RouterGroup)
}

type UserRestHandlerParam struct {
	fx.In

	Logger      logger.Logger
	Tracer      tracer.Tracer
	Middleware  middleware.HttpMiddleware
	Validator   http_validator.HttpValidator
	Response    response.HttpResponse
	UserService user_service.UserService
}

type UserRestHandlerResult struct {
	fx.Out

	UserRestHandler http_interface.RestHandler `group:"http_routes"`
}

type userRestHandlerImpl struct {
	logger      logger.Logger
	tracer      tracer.Tracer
	middleware  middleware.HttpMiddleware
	validator   http_validator.HttpValidator
	response    response.HttpResponse
	userService user_service.UserService
}

func NewUserRestHandler(p UserRestHandlerParam) UserRestHandlerResult {
	return UserRestHandlerResult{
		UserRestHandler: &userRestHandlerImpl{
			logger:      p.Logger,
			tracer:      p.Tracer,
			middleware:  p.Middleware,
			validator:   p.Validator,
			response:    p.Response,
			userService: p.UserService,
		},
	}
}

func (h *userRestHandlerImpl) Register(g *gin.RouterGroup) {
	g.GET("/profile", h.middleware.JwtHttpMiddleware.AuthenticateUser, h.GetUser)
	g.PATCH("/profile", h.middleware.JwtHttpMiddleware.AuthenticateUser, h.PatchUser)
	// TODO
	// g.GET("/profile/email/verify", h.middleware.JwtHttpMiddleware.AuthenticateUser, h.GetVerifyUpdateEmail)
	// g.POST("/profile/email/verify", h.PostVerifyUpdateEmail)
	// g.GET("/profile/phone/verify", h.middleware.JwtHttpMiddleware.AuthenticateUser, h.GetVerifyUpdatePhoneNumber)
	// g.POST("/profile/phone/verify", h.PostVerifyUpdatePhoneNumber)
}

func (h *userRestHandlerImpl) Version() string {
	return "1"
}

func (h *userRestHandlerImpl) GetUser(c *gin.Context) {
	var (
		ctx = c.Request.Context()
		err error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	var fUserInfo *jwt_util.JWTBodyPayload
	userInfo := c.Request.Context().Value(middleware.UserInfoKey)
	if userInfo != nil {
		fUserInfo = userInfo.(*jwt_util.JWTBodyPayload)
	}

	data, err := h.userService.GetUser(ctx, user_service.GetUserParam{
		UUID: fUserInfo.UUID,
	})

	if err != nil {
		h.response.Send(c, nil, err)
	} else {
		h.response.Send(c, data, nil)
	}
}

func (h *userRestHandlerImpl) PatchUser(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		body PatchUserReq
		err  error
	)

	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	if err = h.validator.ValidateRequestBody(c, &body); err != nil {
		h.response.Send(c, nil, err)
		return
	}

	userInfo := c.Request.Context().Value(middleware.UserInfoKey).(*jwt_util.JWTBodyPayload)

	data, err := h.userService.UpdateUser(ctx, user_service.UpdateUserParam{
		UUID:      userInfo.UUID,
		Username:  body.Username,
		FirstName: body.FirstName,
		LastName:  body.LastName,
	})

	if err != nil {
		h.response.Send(c, nil, err)
	} else {
		h.response.Send(c, data, nil)
	}
}
