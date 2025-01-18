package twilio

import (
	"go.uber.org/fx"
)

var TwilioSmsModule = fx.Module(
	"driven-sms-twilio",
	fx.Provide(
		NewTwilioSmsProvider,
	),
)
