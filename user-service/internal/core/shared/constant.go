package shared

type EmailTemplateId string

func (p EmailTemplateId) String() string {
	return string(p)
}

const (
	SignupEmailTemplateId       EmailTemplateId = "signup-email"
	VerificationEmailTemplateId EmailTemplateId = "verification-email"
)
