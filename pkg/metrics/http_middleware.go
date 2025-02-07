package metrics

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusHttpMiddleware interface {
	LogHttpMetrics(c *gin.Context)
}

type prometheusHttpMiddlewareImpl struct {
	tracer   tracer.Tracer
	response http_pkg.HttpResponse
	util     *util.Util
}

func NewPrometheusHttpMiddleware(
	tracer tracer.Tracer,
	response http_pkg.HttpResponse,
	util *util.Util,
) PrometheusHttpMiddleware {
	return &prometheusHttpMiddlewareImpl{
		tracer:   tracer,
		response: response,
		util:     util,
	}
}

func (h *prometheusHttpMiddlewareImpl) LogHttpMetrics(c *gin.Context) {
	_, span := h.tracer.StartSpanWithCaller(c.Request.Context())
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
