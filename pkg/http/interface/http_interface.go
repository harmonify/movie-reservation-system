package http_interface

import "github.com/gin-gonic/gin"

type (
	RestHandler interface {
		Register(g *gin.RouterGroup)
	}
)
