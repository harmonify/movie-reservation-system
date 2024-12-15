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
	"github.com/harmonify/movie-reservation-system/pkg/config"
	constant "github.com/harmonify/movie-reservation-system/pkg/http/constant"
	"github.com/harmonify/movie-reservation-system/pkg/util/encryption"
)

type JWTUtil interface {
	JWTSign(payload *JWTPayload) (string, error)
	JWTVerify(accessToken string) (*JWTBodyPayload, error)
}

type jwtUtilImpl struct {
	encryption *encryption.Encryption
	cfg        *config.Config
}

type JWTPayload struct {
	ExpInMinutes int // expiration in minutes
	SecretKey    string
	PrivateKey   string
	PublicKey    string
	BodyPayload  JWTBodyPayload
}

type JWTBodyPayload struct {
	Email       string `json:"email"`
	UserID      string `json:"userId"` // user UUID
	PhoneNumber string `json:"phoneNumber"`
}

type JWTCustomClaims struct {
	Aud  string         `json:"aud"`
	Sub  string         `json:"sub"`
	Exp  int64          `json:"exp"`
	Data JWTBodyPayload `json:"data"`
}

func NewJWTUtil(
	encryption *encryption.Encryption,
	cfg *config.Config,
) JWTUtil {
	return &jwtUtilImpl{
		encryption: encryption,
		cfg:        cfg,
	}
}

func (i *jwtUtilImpl) JWTSign(payload *JWTPayload) (string, error) {
	decPrivKey, err := i.encryption.AESEncryption.Decrypt(&encryption.AESPayload{
		Secret:  payload.SecretKey,
		Payload: payload.PrivateKey,
	})
	if err != nil {
		return "", err
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(decPrivKey))
	if err != nil {
		return "", err
	}

	decPubKey, err := i.encryption.AESEncryption.Decrypt(&encryption.AESPayload{
		Secret:  payload.SecretKey,
		Payload: payload.PublicKey,
	})
	if err != nil {
		return "", err
	}

	encPubKey, err := i.encryption.AESEncryption.Encrypt(&encryption.AESPayload{
		Secret:  i.cfg.AppSecret,
		Payload: decPubKey,
	})
	if err != nil {
		return "", err
	}

	// Define time expiration
	timeNow := time.Now()
	timeSubtract := time.Duration(payload.ExpInMinutes)
	expDate := timeNow.Add(time.Minute * timeSubtract).Unix()

	// Claim Property
	var claimsProperty JWTCustomClaims
	claimsProperty.Aud = encPubKey
	claimsProperty.Sub = payload.BodyPayload.UserID
	claimsProperty.Exp = expDate
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud":  claimsProperty.Aud,
		"sub":  claimsProperty.Sub,
		"exp":  claimsProperty.Exp,
		"data": payload.BodyPayload,
	})

	// Sign the JWT
	tokenString, err := token.SignedString(privKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (i *jwtUtilImpl) JWTVerify(accessToken string) (*JWTBodyPayload, error) {
	payload, err := i.parser(accessToken)
	if err != nil {
		return nil, err
	}

	decodePubKey, err := i.encryption.AESEncryption.Decrypt(&encryption.AESPayload{
		Secret:  i.cfg.AppSecret,
		Payload: payload.Aud,
	},
	)
	if err != nil {
		return nil, err
	}

	// Decode PEM block
	block, _ := pem.Decode([]byte(decodePubKey))
	if block == nil {
		return nil, err
	}

	// Parse RSA public key
	rsaPublicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, constant.ErrInvalidJwtSigningMethod
		}
		return rsaPublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, constant.ErrInvalidJwt
	}

	return &payload.Data, nil
}

func (i *jwtUtilImpl) parser(accessToken string) (*JWTCustomClaims, error) {
	var claims JWTCustomClaims

	splittedString := strings.Split(accessToken, ".")
	if len(splittedString) < 2 {
		return &claims, constant.ErrInvalidJwtFormat
	}

	encPayload := splittedString[1]

	decPayload, err := base64.RawStdEncoding.DecodeString(encPayload)
	if err != nil {
		return &claims, err
	}

	err = json.Unmarshal([]byte(decPayload), &claims)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return &claims, err
	}

	return &claims, nil
}
