package services_test

import (
	"context"
	"html/template"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/services"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/templates"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/driven/config"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	test_interface "github.com/harmonify/movie-reservation-system/pkg/test/interface"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

func TestEmailTemplateService(t *testing.T) {
	if os.Getenv("CI") == "true" && os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping test")
	}

	suite.Run(t, new(EmailTemplateServiceTestSuite))
}

type renderTestConfig struct {
	Path templates.EmailTemplatePath
	Data map[string]interface{}
}

type renderTestExpectation struct {
	Result string
}

type EmailTemplateServiceTestSuite struct {
	suite.Suite
	app                  *fx.App
	emailTemplateService services.EmailTemplateService
}

func (s *EmailTemplateServiceTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.Provide(
			func() *config.NotificationServiceConfig {
				return &config.NotificationServiceConfig{
					AppDefaultSupportEmail: "support@example.com",
				}
			},
			func() *tracer.TracerConfig {
				return &tracer.TracerConfig{
					ServiceIdentifier: "notification-service",
				}
			},
			logger.NewConsoleLogger,
			tracer.NewNopTracer,
			services.NewEmailTemplateService,
		),
		fx.Invoke(func(t services.EmailTemplateService) {
			s.emailTemplateService = t
		}),
		fx.NopLogger,
	)
	ctx, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()
	if err := s.app.Start(ctx); err != nil {
		s.T().Fatal(">> App failed to start. Error:", err)
	}
}

func (s *EmailTemplateServiceTestSuite) TestEmailTemplateService_Render() {
	_, file, _, _ := runtime.Caller(0)

	testCases := []test_interface.TestCase[renderTestConfig, func() renderTestExpectation]{
		{
			Description: "Should render email verification template correctly",
			Config: renderTestConfig{
				Path: templates.MapEmailTemplateIdToPath(templates.SignupEmailTemplateId.String()),
				Data: map[string]interface{}{
					"firstName": "John",
					"lastName":  "Doe",
					"url":       template.URL("http://localhost:8080/email/verify?email=john_doe@example.com&code=ae60e10ca0a173c2ece3d5d693e1fa21084075aca85edc14a5ce8d58a6503fff"),
				},
			},
			Expectation: func() renderTestExpectation {
				expectedValue, err := os.ReadFile(path.Join(path.Dir(file), "test", "expected-signup.html"))
				s.Require().Nil(err)
				return renderTestExpectation{
					Result: string(expectedValue),
				}
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Description, func() {
			if testCase.BeforeCall != nil {
				testCase.BeforeCall(testCase.Config)
			}

			s.T().Log("Rendering", testCase.Config.Path)
			res, err := s.emailTemplateService.Render(context.TODO(), testCase.Config.Path, testCase.Config.Data)

			if testCase.AfterCall != nil {
				testCase.AfterCall()
			}

			s.Require().Nil(err)
			s.Require().Equal(testCase.Expectation().Result, res)
		})
	}
}
