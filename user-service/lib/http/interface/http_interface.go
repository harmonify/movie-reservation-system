package http_interface

import (
	"time"

	"github.com/gin-gonic/gin"
)

type (
	RestHandler interface {
		Register(g *gin.RouterGroup)
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
