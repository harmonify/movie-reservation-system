package shared

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

const (
	EmailTopicV1_0_0 = `notifications_email_v1.0.0`
	SmsTopicV1_0_0   = `notifications_sms_v1.0.0`
)

var (
	templatesDirPath                                = path.Join(path.Dir(getCurrentFilePath()), "..", "templates")
	EmailVerificationTemplatePath EmailTemplatePath = EmailTemplatePath(path.Join(templatesDirPath, "email-verification.gohtml"))
)

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
