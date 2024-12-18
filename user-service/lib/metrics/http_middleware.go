package metrics

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harmonify/movie-reservation-system/user-service/lib/http/response"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusHttpMiddleware interface {
	LogHttpMetrics(c *gin.Context)
}

type prometheusHttpMiddlewareImpl struct {
	tracer   tracer.Tracer
	response response.HttpResponse
	util     *util.Util
}

func NewPrometheusHttpMiddleware(
	tracer tracer.Tracer,
	response response.HttpResponse,
	util *util.Util,
) PrometheusHttpMiddleware {
	return &prometheusHttpMiddlewareImpl{
		tracer:   tracer,
		response: response,
		util:     util,
	}
}

func (h *prometheusHttpMiddlewareImpl) LogHttpMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := h.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	timer := prometheus.NewTimer(httpDurationCollector.WithLabelValues(
		c.Request.Method,
		c.FullPath(),
	))

	c.Next()

	responseStatusCollector.WithLabelValues(c.Request.Method, c.FullPath(), strconv.Itoa(c.Writer.Status())).Inc()

	totalRequestsCollector.WithLabelValues(c.Request.Method, c.FullPath()).Inc()

	timer.ObserveDuration()
}
