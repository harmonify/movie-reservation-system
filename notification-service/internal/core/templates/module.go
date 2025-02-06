package templates

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"go.uber.org/fx"
)

const (
	SignupEmailTemplateId       EmailTemplateId = "signup-email"
	VerificationEmailTemplateId EmailTemplateId = "verification-email"
)

var (
	templatesDirPath              = path.Dir(getCurrentFilePath())
	signupEmailTemplatePath       = EmailTemplatePath(path.Join(templatesDirPath, "signup.gohtml"))
	verificationEmailTemplatePath = EmailTemplatePath(path.Join(templatesDirPath, "email-verification.gohtml"))
)

var TemplateModule = fx.Module(
	"templates",
	fx.Provide(
		AsTemplate(signupEmailTemplatePath),
		AsTemplate(verificationEmailTemplatePath),
	),
)

func getCurrentFilePath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve correct path"))
	}
	return file
}

func AsTemplate(p EmailTemplatePath) any {
	if _, err := os.Stat(p.String()); os.IsNotExist(err) {
		panic(fmt.Sprintf("template file not found: %s", p))
	}

	return fx.Annotate(
		func() EmailTemplatePath {
			return p
		},
		fx.ResultTags(`group:"email-template-paths"`),
	)
}

type EmailTemplateId string

func (p EmailTemplateId) String() string {
	return string(p)
}

type EmailTemplatePath string

func (p EmailTemplatePath) String() string {
	return string(p)
}

func MapEmailTemplateIdToPath(id string) EmailTemplatePath {
	switch id {
	case SignupEmailTemplateId.String():
		return signupEmailTemplatePath
	case VerificationEmailTemplateId.String():
		return verificationEmailTemplatePath
	default:
		return ""
	}
}
