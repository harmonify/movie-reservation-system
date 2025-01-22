package shared

import (
	"context"
)

type (
	EmailProvider interface {
		// Send sends an email. Send returns email_id string and err error.
		// If err is not nil, email_id should be empty.
		Send(ctx context.Context, msg EmailMessage) (email_id string, err error)
	}

	SmsProvider interface {
		Send(ctx context.Context, message SmsMessage) (sms_id string, err error)
		BulkSend(ctx context.Context, message BulkSmsMessage) (sms_ids []string, err error)
	}
)
