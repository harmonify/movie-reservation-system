package http_pkg

const (
	HttpCookiePrefix = "mrs_"
)

type HttpContextKey string

const (
	// UserInfoKey should be used to get/set user info in the request context
	UserInfoKey HttpContextKey = "userInfo"
	// ParsedBodyKey should be used to get/set parsed body in the request context
	ParsedBodyKey HttpContextKey = "parsedBody"
)
