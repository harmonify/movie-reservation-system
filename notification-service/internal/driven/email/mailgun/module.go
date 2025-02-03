package mailgun

import (
	"go.uber.org/fx"
)

var (
	MailgunMailerModule = fx.Module(
		"driven-email-mailgun",
		fx.Provide(
			NewMailgunEmailProvider,
		),
	)
)
