package http_interface

import (
	"time"

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

	UserToken struct {
		AccessToken  string    `json:"accessToken,omitempty"`
		RefreshToken string    `json:"refreshToken,omitempty"`
		ExpiresAt    time.Time `json:"expiresAt"`
	}

	HeadersExtension struct {
		UserAgent string
		IpAddress string
	}

	VerifierHeader struct {
		TokenVerifier string `json:"token_verifier"`
		HeadersExtension
	}
)
