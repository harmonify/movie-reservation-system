package http_util

import "net/http"

type HttpUtil interface {
	GetUserIP(r *http.Request) string
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
