package health_check_rest

import (
	"github.com/gin-gonic/gin"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
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

	Response http_pkg.HttpResponse
	Tracer   tracer.Tracer
}

type HealthCheckRestHandlerResult struct {
	fx.Out

	HealthCheckRestHandler http_pkg.RestHandler `group:"http_routes"`
}

type healthCheckRestHandlerImpl struct {
	response http_pkg.HttpResponse
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
		err  error
		data = &HealthCheckResponse{
			Ok: true,
		}
	)

	_, span := h.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	h.response.Send(c, data, err)
}
