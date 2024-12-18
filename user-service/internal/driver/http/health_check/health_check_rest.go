package health_rest

import (
	"github.com/gin-gonic/gin"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/response"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
)

type HealthCheckRestHandler interface {
	Register(g *gin.RouterGroup)
	HealthCheck(c *gin.Context)
}

type healthCheckRestHandlerImpl struct {
	response response.HttpResponse
	tracer   tracer.Tracer
}

type HealthCheckResponse struct {
	Ok bool `json:"ok"`
}

func NewHealthCheckRestHandler(response response.HttpResponse, tracer tracer.Tracer) HealthCheckRestHandler {
	return &healthCheckRestHandlerImpl{
		response: response,
		tracer:   tracer,
	}
}

func (h *healthCheckRestHandlerImpl) Register(g *gin.RouterGroup) {
	g.GET("/health", h.HealthCheck)
}

func (h *healthCheckRestHandlerImpl) HealthCheck(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		err  error
		data = &HealthCheckResponse{
			Ok: true,
		}
	)

	_, span := h.tracer.Start(ctx, "healthCheckRestHandlerImpl.HealthCheck")
	defer span.End()

	h.response.Send(c, data, err)
}
