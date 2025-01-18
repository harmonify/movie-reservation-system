package shared

import (
	"context"
)

type (
	EmailProvider interface {
		Send(ctx context.Context, msg EmailMessage) (email_message string, email_id string, err error)
	}

	SmsProvider interface {
		Send(ctx context.Context, message SmsMessage) (sms_id string, err error)
		BulkSend(ctx context.Context, message BulkSmsMessage) (sms_ids []string, err error)
	}
)
