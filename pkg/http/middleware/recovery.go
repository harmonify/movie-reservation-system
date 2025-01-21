package http_middleware

import (
	"net"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	error_constant "github.com/harmonify/movie-reservation-system/pkg/error/constant"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"go.uber.org/zap"
)

// NewRecoveryHttpMiddleware returns a gin.HandlerFunc (middleware) that recovers from any panics and
// logs requests with OTel + Loki + Zap. stack means whether to output the stack info.
// This middleware MUST be registered only after the otelgin middleware is registered.
// Otelgin: go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin
func NewRecoveryHttpMiddleware(response http_pkg.HttpResponse, logger logger.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx := c.Request.Context()

				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.WithCtx(ctx).Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) //nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.WithCtx(ctx).Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.WithCtx(ctx).Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				response.Send(c, nil, error_constant.ErrInternalServerError)
			}
		}()
		c.Next()
	}
}
