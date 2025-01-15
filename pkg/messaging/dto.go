package messaging

import (
	"context"

	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/harmonify/movie-reservation-system/pkg/util"
	"go.uber.org/fx"
)

type Messager interface {
	Send(ctx context.Context, message Message) (id string, err error)
}

type MessagerParam struct {
	fx.In

	Config *config.Config
	Logger logger.Logger
	Tracer tracer.Tracer
	Util   *util.Util
}

type MessagerResult struct {
	fx.Out

	Messager Messager
}

type Message struct {
	To   string // phone number
	Body string
}
