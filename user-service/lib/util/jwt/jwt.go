package jwt_util

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util/encryption"
)

type (
	JWTSignParam struct {
		ExpInSeconds int
		SecretKey    string
		PrivateKey   []byte // in PEM format
		BodyPayload  JWTBodyPayload
	}

	JWTCustomClaims struct {
		*jwt.RegisteredClaims
		Data      JWTBodyPayload `json:"data"`
		PublicKey string         `json:"publicKey"` // in base64 format
	}

	JWTBodyPayload struct {
		UUID        string `json:"uuid"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
	}
)

type JwtUtil interface {
	JWTSign(payload JWTSignParam) (string, error)
	JWTVerify(token string) (*JWTBodyPayload, error)
}

type jwtUtilImpl struct {
	encryption *encryption.Encryption
	cfg        *config.Config
}

func NewJwtUtil(
	encryption *encryption.Encryption,
	cfg *config.Config,
) JwtUtil {
	return &jwtUtilImpl{
		encryption: encryption,
		cfg:        cfg,
	}
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
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer: i.cfg.AppName, // TODO: auth server URI
			// Audience:  jwt.ClaimStrings{}, // TODO: resource servers URI
			Subject:   payload.BodyPayload.UUID,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(payload.ExpInSeconds))),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			// ID: "", // TODO: secure random value
		},
		Data:      payload.BodyPayload,
		PublicKey: string(i.encryption.RSAEncryption.EncodePublicKey(&privKey.PublicKey)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the JWT
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (i *jwtUtilImpl) JWTVerify(token string) (*JWTBodyPayload, error) {
	claims, err := i.decodeClaims(token)
	if err != nil {
		return nil, err
	}

	// Decode PEM block
	block, _ := pem.Decode([]byte(claims.PublicKey))
	if block == nil {
		return nil, fmt.Errorf("Failed to decode PEM: no PEM data is found.")
	}

	// Parse RSA public key
	rsaPublicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	parsedToken, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, constant.ErrInvalidJwtSigningMethod
		}
		return rsaPublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, constant.ErrInvalidJwt
	}

	return &claims.Data, nil
}

func (i *jwtUtilImpl) decodeClaims(token string) (claims *JWTCustomClaims, err error) {
	splittedString := strings.Split(token, ".")
	if len(splittedString) < 2 {
		return claims, constant.ErrInvalidJwtFormat
	}

	// header := splittedString[0]
	encodedClaims := splittedString[1]

	rawClaims, err := base64.RawStdEncoding.DecodeString(encodedClaims)
	if err != nil {
		return claims, err
	}

	err = json.Unmarshal([]byte(rawClaims), &claims)
	if err != nil {
		return claims, err
	}

	return claims, nil
}
