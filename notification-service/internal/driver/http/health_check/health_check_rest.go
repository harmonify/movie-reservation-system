package health_check_rest

import (
	"github.com/gin-gonic/gin"
	http_interface "github.com/harmonify/movie-reservation-system/pkg/http/interface"
	"github.com/harmonify/movie-reservation-system/pkg/http/response"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.uber.org/fx"
)

type HealthCheckRestHandler interface {
	Register(g *gin.RouterGroup)
	Version() string
	GetHealthCheck(c *gin.Context)
}

type HealthCheckRestHandlerParam struct {
	fx.In

	Response response.HttpResponse
	Tracer   tracer.Tracer
}

type HealthCheckRestHandlerResult struct {
	fx.Out

	HealthCheckRestHandler http_interface.RestHandler `group:"http_routes"`
}

type healthCheckRestHandlerImpl struct {
	response response.HttpResponse
	tracer   tracer.Tracer
}

type HealthCheckResponse struct {
	Ok bool `json:"ok"`
}

func NewHealthCheckRestHandler(p HealthCheckRestHandlerParam) HealthCheckRestHandlerResult {
	return HealthCheckRestHandlerResult{
		HealthCheckRestHandler: &healthCheckRestHandlerImpl{
			response: p.Response,
			tracer:   p.Tracer,
		},
	}
}

func (h *healthCheckRestHandlerImpl) Register(g *gin.RouterGroup) {
	g.GET("/health", h.GetHealthCheck)
}

func (h *healthCheckRestHandlerImpl) Version() string {
	return "1"
}

func (h *healthCheckRestHandlerImpl) GetHealthCheck(c *gin.Context) {
	var (
		ctx  = c.Request.Context()
		err  error
		data = &HealthCheckResponse{
			Ok: true,
		}
	)

	_, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	h.response.Send(c, data, err)
}
