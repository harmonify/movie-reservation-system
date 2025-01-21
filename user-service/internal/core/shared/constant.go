package shared

type EmailTemplateId string

func (p EmailTemplateId) String() string {
	return string(p)
}

const (
	EmailVerificationTemplateId EmailTemplateId = "email-verification"
)
