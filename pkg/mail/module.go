package mail

import (
	"go.uber.org/fx"
)

var (
	MailerModule = fx.Provide(NewMailer)
)
