package shared

const (
	SignupEmailTemplateId       EmailTemplateId = "signup-email"
	VerificationEmailTemplateId EmailTemplateId = "verification-email"
)

type EmailTemplateId string

func (p EmailTemplateId) String() string {
	return string(p)
}
