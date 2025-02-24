package http_util

import (
	"net/http"

	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	http_pkg "github.com/harmonify/movie-reservation-system/pkg/http"
	auth_proto "github.com/harmonify/movie-reservation-system/pkg/proto/auth"
	jwt_util "github.com/harmonify/movie-reservation-system/pkg/util/jwt"
)

type HttpUtil interface {
	GetUserIP(r *http.Request) string
	GetUserInfo(r *http.Request) (*jwt_util.JWTBodyPayload, error)
}

type httpUtilImpl struct{}

func NewHttpUtil() HttpUtil {
	return &httpUtilImpl{}
}

func (h *httpUtilImpl) GetUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func (h *httpUtilImpl) GetUserInfo(r *http.Request) (*jwt_util.JWTBodyPayload, error) {
	userInfo := r.Context().Value(http_pkg.UserInfoKey)
	if userInfo == nil {
		return nil, error_pkg.UnauthorizedError
	}

	fUserInfo, ok := userInfo.(*jwt_util.JWTBodyPayload)
	if !ok {
		// Forward compability
		pUserInfo, ok := userInfo.(*auth_proto.UserInfo)
		if !ok || pUserInfo == nil {
			return nil, error_pkg.UnauthorizedError
		}
		fUserInfo = &jwt_util.JWTBodyPayload{
			UUID: pUserInfo.GetUuid(),
		}
	}

	return fUserInfo, nil
}
