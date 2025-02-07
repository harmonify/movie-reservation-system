package shared

import (
	"context"
)

type (
	EmailProvider interface {
		// Send sends an email. Send returns emailId string and err error.
		// If err is not nil, emailId should be empty.
		Send(ctx context.Context, msg EmailMessage) (emailId string, err error)
	}

	SmsProvider interface {
		Send(ctx context.Context, message SmsMessage) (smsId string, err error)
		BulkSend(ctx context.Context, message BulkSmsMessage) (smsIds []string, err error)
	}
)
