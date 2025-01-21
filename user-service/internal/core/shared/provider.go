package shared

import (
	"context"
	"errors"

	notification_proto "github.com/harmonify/movie-reservation-system/user-service/internal/driven/proto/notification"
)

var (
	ErrEmptyRecipient = errors.New("empty recipient")
	ErrEmptySubject   = errors.New("empty subject")
	ErrEmptyTemplate  = errors.New("empty template")

	ErrEmptyMessage = errors.New("empty message body")
)

type (
	EmailProvider interface {
		Send(ctx context.Context, p *notification_proto.Email) error
	}

	SmsProvider interface {
		Send(ctx context.Context, p *notification_proto.Sms) error
	}
)
