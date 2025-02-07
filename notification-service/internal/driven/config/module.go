package config

import "go.uber.org/fx"

var NotificationServiceConfigModule = fx.Module(
	"notification-service-config",
	fx.Provide(NewNotificationServiceConfig),
)
