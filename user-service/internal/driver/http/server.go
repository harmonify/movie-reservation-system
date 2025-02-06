package http_driver

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	http_driver_shared "github.com/harmonify/movie-reservation-system/user-service/internal/driver/http/shared"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type HttpServer struct {
	started bool
	mu      sync.RWMutex

	Server         *http.Server
	Gin            *gin.Engine
	cfg            *HttpServerConfig
	logger         logger.Logger
	httpMiddleware *http_driver_shared.HttpMiddleware
}

type HttpServerConfig struct {
	Env                     string `validate:"required,oneof=dev test prod"`
	ServiceIdentifier       string `validate:"required"`
	ServiceHttpPort         string `validate:"required,numeric"`
	ServiceHttpBaseUrl      string `validate:"required"`
	ServiceHttpBasePath     string `validate:"required"`
	ServiceHttpReadTimeOut  string `validate:"required"`
	ServiceHttpWriteTimeOut string `validate:"required"`
	ServiceHttpEnableCors   bool   `validate:"boolean"`
}

type HttpServerParam struct {
	fx.In

	Logger         logger.Logger
	HttpMiddleware *http_driver_shared.HttpMiddleware
}

type HttpServerResult struct {
	fx.Out

	HttpServer *HttpServer
}

type httpMethodPath struct {
	Method string
	Path   string
}

func NewHttpServer(p HttpServerParam, cfg *HttpServerConfig) (HttpServerResult, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return HttpServerResult{}, err
	}

	gin := gin.New()

	readTimeout, err := time.ParseDuration(cfg.ServiceHttpReadTimeOut)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("HTTP: Failed to parse HTTP read timeout. Error: %v", err))
		return HttpServerResult{}, err
	}

	writeTimeout, err := time.ParseDuration(cfg.ServiceHttpWriteTimeOut)
	if err != nil {
		p.Logger.Error(fmt.Sprintf("HTTP: Failed to parse HTTP write timeout. Error: %v", err))
		return HttpServerResult{}, err
	}

	h := &HttpServer{
		Gin: gin,
		Server: &http.Server{
			Addr:         ":" + cfg.ServiceHttpPort,
			Handler:      gin,
			ReadTimeout:  time.Second * readTimeout,
			WriteTimeout: time.Second * writeTimeout,
		},
		cfg:            cfg,
		logger:         p.Logger,
		httpMiddleware: p.HttpMiddleware,
	}

	return HttpServerResult{
		HttpServer: h,
	}, nil
}

func (h *HttpServer) Start(ctx context.Context) error {
	go func() {
		h.setStarted(true)
		if err := h.Server.ListenAndServe(); err != nil {
			h.setStarted(false)
			h.logger.WithCtx(ctx).Error(fmt.Sprintf(">> HTTP server failed to shutdown gracefully. error: %s", err.Error()))
		}
	}()

	time.Sleep(1 * time.Second)
	if h.getStarted() {
		h.logger.WithCtx(ctx).Info(">> HTTP server started on port " + h.cfg.ServiceHttpPort)
		return nil
	} else {
		err := fmt.Errorf("HTTP server failed to start on port: %s", h.cfg.ServiceHttpPort)
		h.logger.WithCtx(ctx).Error(err.Error())
		return err
	}
}

func (h *HttpServer) Shutdown(ctx context.Context) error {
	var err error
	if err = h.Server.Shutdown(ctx); err == nil {
		h.logger.WithCtx(ctx).Info(">> HTTP server shutdown")
	} else {
		h.logger.WithCtx(ctx).Warn(">> HTTP server failed to shutdown: " + err.Error())
	}
	return err
}

func (h *HttpServer) configure(handlers ...http_pkg.RestHandler) error {
	h.configureMiddlewares()

	if h.cfg.Env == config.EnvironmentProduction {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		h.Gin.TrustedPlatform = gin.PlatformCloudflare
	}

	return h.registerRoutes(handlers...)
}

func (h *HttpServer) configureMiddlewares() {
	h.Gin.Use(h.configureCorsMiddleware)
	h.Gin.Use(otelgin.Middleware(h.cfg.ServiceIdentifier))
	h.Gin.Use(h.httpMiddleware.Recovery.WithStack(true))
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
	h.Gin.Use(h.httpMiddleware.Metrics.LogHttpMetrics)
}

func (h *HttpServer) configureCorsMiddleware(c *gin.Context) {
	if h.cfg.ServiceHttpEnableCors {
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

func (h *HttpServer) registerRoutes(handlers ...http_pkg.RestHandler) error {
	baseGroup := h.Gin.Group(h.cfg.ServiceHttpBasePath)
	groupMap := map[string]*gin.RouterGroup{}
	for _, handler := range handlers {
		version := "v" + handler.Version()
		if _, found := groupMap[version]; !found {
			groupMap[version] = baseGroup.Group(version)
		}
		err := handler.Register(groupMap[version])
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HttpServer) getStarted() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.started
}

func (h *HttpServer) setStarted(started bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.started = started
}
