package middleware

import "go.uber.org/fx"

var HttpMiddlewareModule = fx.Module(
	"http-middleware",
	fx.Provide(
		NewJwtHttpMiddleware,
		NewRbacHttpMiddleware,
		NewHttpMiddleware,
	),
)

type (
	HttpMiddleware struct {
		JwtHttpMiddleware  JwtHttpMiddleware
		RbacHttpMiddleware RbacHttpMiddleware
	}

	HttpMiddlewareParam struct {
		fx.In

		JwtHttpMiddleware  JwtHttpMiddleware
		RbacHttpMiddleware RbacHttpMiddleware
	}

	HttpMiddlewareResult struct {
		fx.Out

		HttpMiddleware *HttpMiddleware
	}
)

func NewHttpMiddleware(p HttpMiddlewareParam) HttpMiddlewareResult {
	return HttpMiddlewareResult{
		HttpMiddleware: &HttpMiddleware{
			JwtHttpMiddleware:  p.JwtHttpMiddleware,
			RbacHttpMiddleware: p.RbacHttpMiddleware,
		},
	}
}
