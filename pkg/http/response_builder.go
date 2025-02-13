package http_pkg

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HttpResponseBuilder interface {
	// New creates a new instance of response
	New() HttpResponseV2
}

type HttpResponseV2 interface {
	// WithCtx extracts necessary information from context and sets it to the response
	// It currently only extracts traceId from the context
	WithCtx(ctx context.Context) HttpResponseV2
	// WithTraceId sets the traceId to the response
	WithTraceId(traceId string) HttpResponseV2
	// WithStatusCode sets the http status code to the response
	WithStatusCode(statusCode int) HttpResponseV2
	// WithError sets the error to the response. WithError sets the sets response body's `error` and `success` to false.
	// If the error is an instance of `error_pkg.ErrorWithDetails`, the HTTP status code will be set to the error's HTTP code.
	WithError(err error) HttpResponseV2
	// WithMetadata sets metadata to the response. WithMetadata appends the metadata to the existing set of metadata
	WithMetadata(metadata map[string]interface{}) HttpResponseV2
	// WithMetadataFromStruct sets metadata to the response from a struct. WithMetadataFromStruct appends the metadata to the existing set of metadata
	WithMetadataFromStruct(metadata interface{}) HttpResponseV2
	// WithPaginationMetadata sets the pagination metadata to the response
	WithPaginationMetadata(page, limit, total, totalPages int) HttpResponseV2
	// WithResult sets the result to the response
	WithResult(result interface{}) HttpResponseV2
	// Send sends the response to the client. Send also logs the response and captures error message and stack trace if the response error is an uncaptured exception
	Send(c *gin.Context)
}

type httpResponseBuilderImpl struct {
	logger      logger.Logger
	tracer      tracer.Tracer
	errorMapper error_pkg.ErrorMapper
}

type HttpResponseBuilderParam struct {
	fx.In
	logger.Logger
	tracer.Tracer
	error_pkg.ErrorMapper
}

type httpResponseV2Impl struct {
	logger      logger.Logger
	tracer      tracer.Tracer
	errorMapper error_pkg.ErrorMapper
	// Data
	ctx                  context.Context
	responseHttpCode     int
	responseBodySuccess  bool
	responseBodyTraceId  string
	responseBodyMetadata map[string]interface{}
	responseBodyError    error
	responseBodyResult   map[string]interface{}
}

func NewHttpResponseBuilder(p HttpResponseBuilderParam) HttpResponseBuilder {
	return &httpResponseBuilderImpl{
		logger:      p.Logger,
		tracer:      p.Tracer,
		errorMapper: p.ErrorMapper,
	}
}

func (b *httpResponseBuilderImpl) New() HttpResponseV2 {
	return &httpResponseV2Impl{
		logger:      b.logger,
		tracer:      b.tracer,
		errorMapper: b.errorMapper,
		// Default values
		responseHttpCode:     http.StatusOK,
		responseBodySuccess:  true,
		responseBodyTraceId:  "",
		responseBodyMetadata: make(map[string]interface{}),
		responseBodyError:    nil,
		responseBodyResult:   make(map[string]interface{}),
	}
}

func (b *httpResponseV2Impl) WithCtx(ctx context.Context) HttpResponseV2 {
	b.ctx = ctx
	traceId := trace.SpanContextFromContext(ctx).TraceID().String()
	return b.WithTraceId(traceId)
}

func (b *httpResponseV2Impl) WithTraceId(traceId string) HttpResponseV2 {
	b.responseBodyTraceId = traceId
	return b
}

func (b *httpResponseV2Impl) WithStatusCode(statusCode int) HttpResponseV2 {
	b.responseHttpCode = statusCode
	return b
}

func (b *httpResponseV2Impl) WithError(err error) HttpResponseV2 {
	if err == nil {
		return b
	}

	detailedError, valid := b.errorMapper.FromError(err)

	if !valid {
		b.logger.WithCtx(b.ctx).Debug("Uncatched error", zap.Error(err), zap.String("error_type", reflect.TypeOf(err).String()))
	}

	b.responseBodySuccess = false
	b.responseHttpCode = detailedError.HttpCode
	b.responseBodyError = detailedError

	return b
}

func (b *httpResponseV2Impl) WithMetadata(metadata map[string]interface{}) HttpResponseV2 {
	// Append to existing metadata
	for key, value := range metadata {
		b.responseBodyMetadata[key] = value
	}

	return b
}

func (b *httpResponseV2Impl) WithMetadataFromStruct(metadata interface{}) HttpResponseV2 {
	v := reflect.ValueOf(metadata)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return b
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" {
			b.responseBodyMetadata[fieldType.Name] = field.Interface()
		} else {
			tags := strings.Split(jsonTag, ",")
			fieldName := tags[0]
			if fieldName == "-" {
				continue
			} else if fieldName == "" {
				fieldName = fieldType.Name
			} else if slices.Contains(tags, "omitempty") && reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
				continue
			}
			b.responseBodyMetadata[fieldName] = field.Interface()
		}
	}

	return b
}

func (b *httpResponseV2Impl) WithPaginationMetadata(page, limit, total, totalPages int) HttpResponseV2 {
	b.responseBodyMetadata["page"] = page
	b.responseBodyMetadata["limit"] = limit
	b.responseBodyMetadata["total"] = total
	b.responseBodyMetadata["totalPages"] = totalPages
	return b
}

func (b *httpResponseV2Impl) WithResult(result interface{}) HttpResponseV2 {
	b.responseBodyResult = make(map[string]interface{})
	b.responseBodyResult["data"] = result
	return b
}

func (r *httpResponseV2Impl) Send(c *gin.Context) {
	ctx, _ := r.tracer.StartSpanWithCaller(c.Request.Context())
	response := r.build()
	r.logResponse(ctx, response)
	c.JSON(response.HttpCode, response.Body)
}

func (b *httpResponseV2Impl) build() *ResponseSchema {
	return &ResponseSchema{
		HttpCode: b.responseHttpCode,
		Body: &ResponseBodySchema{
			Success:  b.responseBodySuccess,
			TraceId:  b.responseBodyTraceId,
			Error:    buildResponseError(b.responseBodyError),
			Metadata: b.responseBodyMetadata,
			Result:   b.responseBodyResult,
		},
	}
}

func buildResponseError(err error) *ResponseBodyErrorSchema {
	errs := make([]error, 0)

	if err == nil {
		return &ResponseBodyErrorSchema{
			Errors: errs,
		}
	}

	var detailedError *error_pkg.ErrorWithDetails
	if !errors.As(err, &detailedError) {
		detailedError = error_pkg.InternalServerError
		errs = append(errs, detailedError.Errors...)
	}

	return &ResponseBodyErrorSchema{
		Original: err,
		Code:     detailedError.Code.String(),
		Message:  detailedError.Message,
		Errors:   detailedError.Errors,
	}
}

func (r *httpResponseV2Impl) logResponse(ctx context.Context, response *ResponseSchema) {
	span := trace.SpanFromContext(ctx)
	fields := []zap.Field{}

	if response.Body.Error.Original != nil {
		// Capture error message
		span.SetStatus(codes.Error, response.Body.Error.Original.Error())

		// Capture error stack
		var dErr *error_pkg.ErrorWithDetails
		if !errors.As(response.Body.Error.Original, &dErr) {
			span.RecordError(response.Body.Error.Original, trace.WithStackTrace(true))
			fields = append(fields, zap.StackSkip("stack", 2))
		} else if dErr.HttpCode >= http.StatusInternalServerError {
			span.RecordError(dErr, trace.WithStackTrace(true))
			fields = append(fields, zap.StackSkip("stack", 2))
		}
	}

	// Log response body, if it's an error or debug level
	if response.HttpCode >= http.StatusInternalServerError || r.logger.Level() == zap.DebugLevel {
		byteResponse, _ := json.Marshal(response.Body)
		fields = append(fields, zap.ByteString("response", byteResponse))
	}

	if response.HttpCode >= http.StatusInternalServerError {
		r.logger.WithCtx(ctx).Error("Response error", fields...)
	} else if r.logger.Level() == zap.DebugLevel {
		r.logger.WithCtx(ctx).Debug("Response debug", fields...)
	}
}
