package shared

const (
	EmailTypeHtml  EmailType = "html"
	EmailTypePlain EmailType = "plain"
)

type EmailType string

func (p EmailType) String() string {
	return string(p)
}
