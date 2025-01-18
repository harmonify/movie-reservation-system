package mailgun

import (
	"context"

	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/mailgun/mailgun-go"
	"go.uber.org/fx"
)

type (
	MailgunEmailProviderParam struct {
		fx.In

		Config *config.Config
		Logger logger.Logger
		Tracer tracer.Tracer
	}

	MailgunEmailProviderResult struct {
		fx.Out

		EmailProvider shared.EmailProvider
	}

	mailgunEmailProviderImpl struct {
		cfg    *config.Config
		logger logger.Logger
		tracer tracer.Tracer
		mg     mailgun.Mailgun // https://github.com/mailgun/mailgun-go/blob/master/examples/examples.go
	}
)

func NewMailgunEmailProvider(p MailgunEmailProviderParam) MailgunEmailProviderResult {
	mg := mailgun.NewMailgun(p.Config.MailgunDomain, p.Config.MailgunApiKey)
	return MailgunEmailProviderResult{
		EmailProvider: &mailgunEmailProviderImpl{
			cfg:    p.Config,
			mg:     mg,
			logger: p.Logger,
			tracer: p.Tracer,
		},
	}
}

func (m *mailgunEmailProviderImpl) Send(ctx context.Context, message shared.EmailMessage) (msg string, id string, err error) {
	ctx, span := m.tracer.StartSpanWithCaller(ctx)
	defer span.End()
	return m.mg.Send(m.mg.NewMessage(
		m.cfg.MailgunDefaultSender,
		message.Subject,
		message.Body,
		message.Recipients...,
	))
}
