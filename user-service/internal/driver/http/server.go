package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	config "github.com/harmonify/movie-reservation-system/pkg/config"
	constant "github.com/harmonify/movie-reservation-system/pkg/config/constant"
	http_interface "github.com/harmonify/movie-reservation-system/pkg/http/interface"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/metrics"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type HttpServer struct {
	Server      *http.Server
	Gin         *gin.Engine
	cfg         *config.Config
	logger      logger.Logger
	middlewares *httpServerMiddlewares
}

type HttpServerParam struct {
	fx.In

	Lifecycle         fx.Lifecycle
	Config            *config.Config
	Logger            logger.Logger
	MetricsMiddleware metrics.PrometheusHttpMiddleware
	Routes            []http_interface.RestHandler `group:"http_routes"`
}

type HttpServerResult struct {
	fx.Out

	HttpServer *HttpServer
}

type httpServerMiddlewares struct {
	metrics metrics.PrometheusHttpMiddleware
}

type httpMethodPath struct {
	Method string
	Path   string
}

func NewHttpServer(p HttpServerParam) (HttpServerResult, error) {
	gin := gin.New()

	readTimeout, err := time.ParseDuration(p.Config.ServiceHttpReadTimeOut)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("HTTP: Failed to parse HTTP read timeout. Error: %v", err))
		return HttpServerResult{}, err
	}

	writeTimeout, err := time.ParseDuration(p.Config.ServiceHttpWriteTimeOut)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("HTTP: Failed to parse HTTP write timeout. Error: %v", err))
		return HttpServerResult{}, err
	}

	h := &HttpServer{
		Gin: gin,
		Server: &http.Server{
			Addr:         ":" + p.Config.ServicePort,
			Handler:      gin,
			ReadTimeout:  time.Second * readTimeout,
			WriteTimeout: time.Second * writeTimeout,
		},
		cfg:    p.Config,
		logger: p.Logger,
		middlewares: &httpServerMiddlewares{
			metrics: p.MetricsMiddleware,
		},
	}

	h.configure(p.Routes...)

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := h.Start(ctx); err != nil {
				h.logger.WithCtx(ctx).Error("Failed to start HTTP server", zap.Error(err))
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := h.Shutdown(ctx); err != nil {
				h.logger.WithCtx(ctx).Error("Failed to shutdown HTTP server", zap.Error(err))
				return err
			}
			return nil
		},
	})

	return HttpServerResult{
		HttpServer: h,
	}, nil
}

func (h *HttpServer) Start(ctx context.Context) error {
	h.logger.WithCtx(ctx).Info(">> HTTP server run on port: " + h.cfg.ServicePort)
	var err error
	if err = h.Server.ListenAndServe(); err == nil {
		h.logger.WithCtx(ctx).Info(">> HTTP server started on port " + h.cfg.ServicePort)
	} else {
		h.logger.WithCtx(ctx).Error(">> HTTP server failed to start. Error: " + err.Error())
	}
	return err
}

func (h *HttpServer) Shutdown(ctx context.Context) error {
	var err error
	if err = h.Server.Shutdown(ctx); err == nil {
		h.logger.WithCtx(ctx).Info(">> HTTP server shutdown")
	} else {
		h.logger.WithCtx(ctx).Warn(">> HTTP server failed to shutdown")
	}
	return err
}

func (h *HttpServer) configure(handlers ...http_interface.RestHandler) {
	h.configureMiddlewares()

	if h.cfg.Env == constant.EnvironmentProduction {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		h.Gin.TrustedPlatform = gin.PlatformCloudflare
	}

	h.registerRoutes(handlers...)
}

func (h *HttpServer) configureMiddlewares() {
	h.Gin.Use(h.configureCorsMiddleware)
	h.Gin.Use(otelgin.Middleware(h.cfg.AppName))
	h.Gin.Use(ginzap.RecoveryWithZap(h.logger.GetZapLogger(), true))
	h.Gin.Use(ginzap.GinzapWithConfig(h.logger.GetZapLogger(), &ginzap.Config{
		TimeFormat: time.RFC3339Nano,
		UTC:        true,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			// log request ID
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			// log trace and span ID
			if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
				fields = append(fields, zap.String("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()))
				fields = append(fields, zap.String("span_id", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()))
			}

			// log request body
			var body []byte
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)
			fields = append(fields, zap.String("body", string(body)))

			return fields
		}),
		Skipper: func(c *gin.Context) bool {
			skip_list := []httpMethodPath{
				{
					Method: "GET",
					Path:   "/health",
				},
				{
					Method: "GET",
					Path:   "/ping",
				},
				{
					Method: "GET",
					Path:   "/metrics",
				},
			}

			for _, el := range skip_list {
				if c.Request.Method == el.Path && c.Request.URL.Path == el.Path {
					return true
				}
			}

			return false
		},
	}))
	h.Gin.Use(h.middlewares.metrics.LogHttpMetrics)
}

func (h *HttpServer) configureCorsMiddleware(c *gin.Context) {
	if h.cfg.ServiceEnableCors {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
	}

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func (h *HttpServer) registerRoutes(handlers ...http_interface.RestHandler) {
	baseGroup := h.Gin.Group(h.cfg.ServiceBasePath)
	groupMap := map[string]*gin.RouterGroup{}
	for _, handler := range handlers {
		version := "v" + handler.Version()
		if _, found := groupMap[version]; !found {
			groupMap[version] = baseGroup.Group(version)
		}
		handler.Register(groupMap[version])
	}
}
