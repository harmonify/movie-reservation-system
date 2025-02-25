package shared

import (
	"context"
	"errors"

	notification_proto "github.com/harmonify/movie-reservation-system/pkg/proto/notification"
)

var (
	ErrEmptyRecipient = errors.New("empty recipient")
	ErrEmptySubject   = errors.New("empty subject")
	ErrEmptyTemplate  = errors.New("empty template")

	ErrEmptyMessage = errors.New("empty message body")
)

type (
	NotificationProvider interface {
		SendEmail(ctx context.Context, p *notification_proto.SendEmailRequest) error
		SendSms(ctx context.Context, p *notification_proto.SendSmsRequest) error
		BulkSendSms(ctx context.Context, p *notification_proto.BulkSendSmsRequest) error
	}
)
