package jwt_util

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	error_constant "github.com/harmonify/movie-reservation-system/pkg/error/constant"
	"github.com/harmonify/movie-reservation-system/pkg/util/encryption"
)

type (
	JWTCustomClaims struct {
		Data JWTBodyPayload `json:"data"`
		jwt.RegisteredClaims
	}

	JWTHeaderPayload struct {
		Typ string `json:"typ"` // JWT
		Alg string `json:"alg"` // RS256
		Kid string `json:"kid"` // user public key
	}

	JWTBodyPayload struct {
		UUID string `json:"uuid"`
	}

	JWTSignParam struct {
		ExpInSeconds int
		PrivateKey   []byte // in PEM format
		BodyPayload  JWTBodyPayload
	}
)

type JwtUtil interface {
	JWTSign(payload JWTSignParam) (string, error)
	JWTVerify(token string) (*JWTBodyPayload, error)
}

type jwtUtilImpl struct {
	encryption *encryption.Encryption
	config     *JwtUtilConfig
}

type JwtUtilConfig struct {
	AppJwtAudiences string
	ServiceBaseUrl  string
}

func NewJwtUtil(
	encryption *encryption.Encryption,
	config *config.Config,
) (JwtUtil, error) {
	if config.AppJwtAudiences == "" {
		return nil, fmt.Errorf("AppJwtAudiences is empty")
	}
	if config.ServiceBaseUrl == "" {
		return nil, fmt.Errorf("ServiceBaseUrl is empty")
	}

	return &jwtUtilImpl{
		encryption: encryption,
		config: &JwtUtilConfig{
			AppJwtAudiences: config.AppJwtAudiences,
			ServiceBaseUrl:  config.ServiceBaseUrl,
		},
	}, nil
}

func (i *jwtUtilImpl) JWTSign(payload JWTSignParam) (string, error) {
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(payload.PrivateKey)
	if err != nil {
		return "", err
	}

	// Define time expiration
	now := time.Now()

	// Claim Property
	claims := JWTCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    i.config.ServiceBaseUrl,
			Subject:   payload.BodyPayload.UUID,
			Audience:  strings.Split(i.config.AppJwtAudiences, ","),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(payload.ExpInSeconds))),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			// ID: "", // not needed, apparently to prevent replay attack but it will cause the token to be one-time use, open an issue if this is false
		},
		Data: payload.BodyPayload,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = string(i.encryption.RSAEncryption.EncodePublicKey(&privKey.PublicKey))

	// Sign the JWT
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (i *jwtUtilImpl) JWTVerify(tokenString string) (*JWTBodyPayload, error) {
	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		&JWTCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			publicKey, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("failed to get public key from token header")
			}
			if publicKey == "" {
				return nil, fmt.Errorf("public key is empty")
			}

			// Decode PEM block
			block, _ := pem.Decode([]byte(publicKey))
			if block == nil {
				return nil, fmt.Errorf("failed to decode PEM: no PEM data is found.")
			}

			// Parse RSA public key
			rsaPublicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				return nil, err
			}

			return rsaPublicKey, nil
		},
		jwt.WithAudience(i.config.ServiceBaseUrl),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(i.config.ServiceBaseUrl),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, error_constant.ErrInvalidJwt
	}

	claims, ok := parsedToken.Claims.(*JWTCustomClaims)
	if !ok {
		return nil, error_constant.ErrInvalidJwtClaims
	}

	return &claims.Data, nil
}
