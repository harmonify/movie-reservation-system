package mailgun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/harmonify/movie-reservation-system/notification-service/internal/core/shared"
	"github.com/harmonify/movie-reservation-system/notification-service/internal/driven/config"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"github.com/mailgun/mailgun-go"
	"go.uber.org/fx"
)

var (
	rateLimitErrorMessagePattern = `try again after (\d+) seconds`
	rateLimitErrorMessageRegex   = regexp.MustCompile(rateLimitErrorMessagePattern)
)

type (
	MailgunEmailProviderParam struct {
		fx.In

		Config *config.NotificationServiceConfig
		Logger logger.Logger
		Tracer tracer.Tracer
	}

	MailgunEmailProviderResult struct {
		fx.Out

		EmailProvider shared.EmailProvider
	}

	mailgunEmailProviderImpl struct {
		cfg    *config.NotificationServiceConfig
		logger logger.Logger
		tracer tracer.Tracer
		mg     mailgun.Mailgun // https://github.com/mailgun/mailgun-go/blob/master/examples/examples.go
	}

	sendHttpResponse struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	}
)

func NewMailgunEmailProvider(p MailgunEmailProviderParam) MailgunEmailProviderResult {
	mg := mailgun.NewMailgun(p.Config.MailgunDomain, p.Config.MailgunApiKey)
	return MailgunEmailProviderResult{
		EmailProvider: &mailgunEmailProviderImpl{
			cfg:    p.Config,
			mg:     mg,
			logger: p.Logger,
			tracer: p.Tracer,
		},
	}
}

// https://mailgun-docs.redoc.ly/docs/mailgun/api-reference/openapi-final/tag/Messages/#tag/Messages/operation/POST-v3--domain-name--messages
func (m *mailgunEmailProviderImpl) Send(ctx context.Context, message shared.EmailMessage) (string, error) {
	ctx, span := m.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	_, id, err := m.mg.Send(m.mg.NewMessage(
		m.cfg.MailgunDefaultSender,
		message.Subject,
		message.Body,
		message.Recipients...,
	))
	if err != nil {
		// Check for HTTP 429 error
		var httpErr *mailgun.UnexpectedResponseError
		if errors.As(err, &httpErr) {
			if httpErr.Actual == http.StatusTooManyRequests {
				var resp sendHttpResponse
				err := json.Unmarshal(httpErr.Data, &resp)
				if err != nil {
					fmt.Printf("Failed to unmarshal HTTP error message: %v\n", err)
					return "", err
				}
				retryAfter, parseErr := m.extractRetryAfter(resp.Message)
				if parseErr != nil {
					fmt.Printf("Failed to extract retry-after: %v\n", parseErr)
					return "", parseErr
				}
				return "", error_pkg.NewRateLimitExceededError(retryAfter)
			} else if httpErr.Actual >= 400 && httpErr.Actual < 500 {
				return "", error_pkg.InternalServerError
			} else if httpErr.Actual >= 500 {
				return "", error_pkg.BadGatewayError
			}
		}

		// Handle other errors
		fmt.Printf("Failed to send email: %v\n", err)
		return "", error_pkg.InternalServerError
	}

	return id, nil
}

func (m *mailgunEmailProviderImpl) extractRetryAfter(message string) (int, error) {
	// Find the first match
	matches := rateLimitErrorMessageRegex.FindStringSubmatch(message)
	if len(matches) < 2 {
		return 0, errors.New("retry-after seconds not found in the message")
	}

	// Convert the captured group (matches[1]) to an integer
	retryAfter, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("failed to parse retry-after seconds: %w", err)
	}

	return retryAfter, nil
}
