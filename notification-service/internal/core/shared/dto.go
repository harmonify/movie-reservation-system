package shared

type (
	EmailTemplatePath string

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

func (p *EmailTemplatePath) String() string {
	return string(*p)
}
