package http_pkg

import (
	"github.com/gin-gonic/gin"
)

type (
	RestHandler interface {
		// Register will be invoked when starting HTTP server
		// this function can be used to register REST API routes
		Register(g *gin.RouterGroup)
		// Version return API version for this REST handler
		Version() string
	}
)
