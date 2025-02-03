package shared

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"slices"
)

const (
	EmailTypeHtml  EmailType = "html"
	EmailTypePlain EmailType = "plain"

	EmailVerificationTemplateId EmailTemplateId = "email-verification"
)

var (
	templatesDirPath                                = path.Join(path.Dir(getCurrentFilePath()), "..", "templates")
	EmailVerificationTemplatePath EmailTemplatePath = EmailTemplatePath(path.Join(templatesDirPath, "email-verification.gohtml"))
)

func ValidateEmailTemplateId(id string) bool {
	return slices.Contains([]string{
		EmailVerificationTemplateId.String(),
	}, id)
}

func MapEmailTemplateIdToPath(id EmailTemplateId) EmailTemplatePath {
	switch id {
	case EmailVerificationTemplateId:
		return EmailVerificationTemplatePath
	default:
		return ""
	}
}

func init() {
	if _, err := os.Stat(templatesDirPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("templates directory not found: %s", templatesDirPath))
	}
	if _, err := os.Stat(EmailVerificationTemplatePath.String()); os.IsNotExist(err) {
		panic(fmt.Sprintf("email verification template file not found: %s", EmailVerificationTemplatePath))
	}
}

func getCurrentFilePath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve correct path"))
	}
	return file
}

type EmailType string

func (p EmailType) String() string {
	return string(p)
}

type EmailTemplateId string

func (p EmailTemplateId) String() string {
	return string(p)
}

type EmailTemplatePath string

func (p EmailTemplatePath) String() string {
	return string(p)
}
