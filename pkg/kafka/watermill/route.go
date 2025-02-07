package watermill_pkg

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/fx"
)

type Route interface {
	Identifier() string
	Register(router *message.Router, subscriber message.Subscriber) error
}

func AsRoute(f any, anns ...fx.Annotation) any {
	finalAnns := []fx.Annotation{
		fx.As(new(Route)),
		fx.ResultTags(`group:"kafka-routes"`),
	}
	if len(anns) > 0 {
		finalAnns = append(finalAnns, anns...)
	}

	return fx.Annotate(
		f,
		finalAnns...,
	)
}
func AsRouter(f any, anns ...fx.Annotation) any {
	finalAnns := []fx.Annotation{
		fx.ParamTags(`group:"kafka-routes"`),
	}
	if len(anns) > 0 {
		finalAnns = append(finalAnns, anns...)
	}

	return fx.Annotate(
		f,
		finalAnns...,
	)
}
