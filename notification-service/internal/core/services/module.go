package services

import (
	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/templates"
	"go.uber.org/fx"
)

var ServiceModule = fx.Module(
	"service",
	templates.TemplateModule,
	fx.Provide(
		NewEmailTemplateService,
		NewEmailService,
		NewSmsService,
	),
)
