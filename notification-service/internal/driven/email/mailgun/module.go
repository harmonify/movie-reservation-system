package mailgun

import (
	"go.uber.org/fx"
)

var (
	MailgunMailerModule = fx.Module(
		"driven-mail-mailgun",
		fx.Provide(
			NewMailgunEmailProvider,
		),
	)
)
