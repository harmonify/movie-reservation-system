package mail

import (
	"context"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"go.uber.org/fx"
)

type (
	Mailer interface {
		Send(ctx context.Context, message Message) (id string, err error)
	}

	MailerParam struct {
		fx.In

		Config *config.Config
	}

	MailerResult struct {
		fx.Out

		Mailer Mailer
	}

	Message struct {
		Subject string
		Body    string
		To      []string
	}
)

func NewMailer(p MailerParam) MailerResult {
	return NewMailgunMailer(p)
}
