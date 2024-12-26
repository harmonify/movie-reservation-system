package messaging

import (
	"go.uber.org/fx"
)

var MessagingModule = fx.Module(
	"messaging",
	fx.Provide(
		NewTwilioMessager,
	),
)

func NewMessager(p MessagerParam) MessagerResult {
	return NewTwilioMessager(p)
}
