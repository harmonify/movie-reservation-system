package mail

import (
	"context"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/mailgun/mailgun-go"
)

func NewMailgunMailer(p MailerParam) MailerResult {
	mg := mailgun.NewMailgun(p.Config.MailgunDomain, p.Config.MailgunApiKey)

	return MailerResult{
		Mailer: &mailgunMailerImpl{
			cfg: p.Config,
			mg:  mg,
		},
	}
}

type mailgunMailerImpl struct {
	cfg *config.Config
	// https://github.com/mailgun/mailgun-go/blob/master/examples/examples.go
	mg mailgun.Mailgun
}

func (m *mailgunMailerImpl) Send(ctx context.Context, message Message) (id string, err error) {
	msg := m.mg.NewMessage(
		m.cfg.MailgunDefaultSender,
		message.Subject,
		message.Body,
		message.To...,
	)

	_, id, err = m.mg.Send(msg)
	return
}
