package shared

type EmailTemplateId string

func (p EmailTemplateId) String() string {
	return string(p)
}

const (
	EmailVerificationTemplateId EmailTemplateId = "email-verification"
)

type MessageBrokerTopic string

func (p MessageBrokerTopic) String() string {
	return string(p)
}

const (
	PublicUserRegisteredV1 MessageBrokerTopic = "public.user.registered.v1"
)
