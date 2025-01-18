package services

import (
	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/shared"
	"go.uber.org/fx"
)

var ServiceModule = fx.Module(
	"service",
	fx.Provide(
		fx.Annotate(
			func() []shared.EmailTemplatePath {
				return []shared.EmailTemplatePath{
					shared.EmailVerificationTemplatePath,
				}
			},
			fx.ResultTags(`group:"email-template-paths"`),
		),
		NewEmailService,
		NewEmailTemplateService,
		NewSmsService,
	),
)

func AsTemplate(f func() shared.EmailTemplatePath, anns ...fx.Annotation) any {
	finalAnns := []fx.Annotation{
		fx.As(new(shared.EmailTemplatePath)),
		fx.ResultTags(`group:"email-template-paths"`),
	}
	if len(anns) > 0 {
		finalAnns = append(finalAnns, anns...)
	}

	return fx.Annotate(
		f,
		finalAnns...,
	)
}
