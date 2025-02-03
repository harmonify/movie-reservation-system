package jwt_util

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
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
	AppJwtAudiences    string `validate:"required"`
	ServiceHttpBaseUrl string `validate:"required"`
}

func NewJwtUtil(
	encryption *encryption.Encryption,
	cfg *JwtUtilConfig,
) (JwtUtil, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return nil, err
	}
	return &jwtUtilImpl{
		encryption: encryption,
		config: &JwtUtilConfig{
			AppJwtAudiences:    cfg.AppJwtAudiences,
			ServiceHttpBaseUrl: cfg.ServiceHttpBaseUrl,
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
			Issuer:    i.config.ServiceHttpBaseUrl,
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
		jwt.WithAudience(i.config.ServiceHttpBaseUrl),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithIssuer(i.config.ServiceHttpBaseUrl),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, error_pkg.InvalidJwtError
	}

	claims, ok := parsedToken.Claims.(*JWTCustomClaims)
	if !ok {
		return nil, error_pkg.InvalidJwtClaimsError
	}

	return &claims.Data, nil
}
