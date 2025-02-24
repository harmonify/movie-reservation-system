package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type OpaClient interface {
	// HasAccess checks if the input satisfies the policy.
	// policyId is the ID of the policy to check.
	// input is the input to the policy. It can be any type that can be marshaled to JSON.
	// Returns true if the input satisfies the policy, false otherwise.
	// Returns an error if there was an error checking the policy.
	HasAccess(ctx context.Context, policyId string, input OpaRequestBody) (bool, error)
}

type OpaClientParam struct {
	fx.In
	Logger logger.Logger
	Tracer tracer.Tracer
}

type OpaClientConfig struct {
	OpaServerUrl string `validate:"required,url"`
}

type OpaClientResult struct {
	fx.Out
	OpaClient OpaClient
}

type opaClientImpl struct {
	logger logger.Logger
	tracer tracer.Tracer
	config *OpaClientConfig
}

func NewOpaClient(p OpaClientParam, cfg *OpaClientConfig) OpaClientResult {
	return OpaClientResult{
		OpaClient: &opaClientImpl{
			logger: p.Logger,
			tracer: p.Tracer,
			config: cfg,
		},
	}
}

func (o *opaClientImpl) HasAccess(ctx context.Context, policyId string, input OpaRequestBody) (bool, error) {
	ctx, span := o.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	// Replace the dots in the policy ID with slashes to match the OPA server's policy path.
	policyPath := strings.ReplaceAll(policyId, ".", "/")
	url := fmt.Sprintf("%s/v1/data/%s/allow", o.config.OpaServerUrl, policyPath)

	jsonBody, err := json.Marshal(input)
	if err != nil {
		o.logger.Error("Failed to marshal json", zap.Error(err))
		return false, error_pkg.InternalServerError
	}

	span.SetAttributes(
		attribute.String("opa_policy_id", policyId),
		attribute.String("opa_policy_url", url),
		attribute.String("opa_input", string(jsonBody)),
	)

	// Call the OPA server REST API.
	result, err := http.Post(url, "application/json", strings.NewReader(string(jsonBody)))
	if err != nil {
		o.logger.Error("Failed to call OPA server", zap.Error(err))
		return false, error_pkg.InternalServerError
	}
	defer result.Body.Close()

	span.SetAttributes(
		attribute.Int("opa_response_status_code", result.StatusCode),
	)

	// Decode the response using OpaResponseBody
	var opaResponseBody OpaResponseBody
	err = json.NewDecoder(result.Body).Decode(&opaResponseBody)
	if err != nil {
		o.logger.Error("Failed to decode OPA response", zap.Error(err))
		return false, error_pkg.InternalServerError
	}

	return opaResponseBody.Result, nil
}
