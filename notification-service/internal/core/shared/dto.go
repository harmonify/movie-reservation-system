package shared

type (
	EmailMessage struct {
		Recipients []string // emails
		Subject    string
		Body       string
	}

	SmsMessage struct {
		Recipient string // phone number
		Body      string
	}

	BulkSmsMessage struct {
		Recipients []string // phone numbers
		Body       string
	}
)
