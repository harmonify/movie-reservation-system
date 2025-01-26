package http_pkg

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	error_pkg "github.com/harmonify/movie-reservation-system/pkg/error"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	struct_util "github.com/harmonify/movie-reservation-system/pkg/util/struct"
	"github.com/harmonify/movie-reservation-system/pkg/util/validation"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type HttpResponse interface {
	Send(c *gin.Context, data interface{}, err error)
	SendWithResponseCode(c *gin.Context, successHttpCode int, data interface{}, err error)
	// Build is a function to build the response schema
	// Build takes in the ctx, success http code (only respected when err is nil), data, and err
	// Build returns the http code, response schema, and error
	// The error is an instance of HttpError and is useful for debugging
	// You typically don't need to call this function directly
	// Instead, use Send or SendWithResponseCode
	Build(ctx context.Context, successHttpCode int, data interface{}, err *error_pkg.ErrorWithDetails) (responseHttpCode int, responseBody *Response, responseError *HttpError)
}

type httpResponseImpl struct {
	logger      logger.Logger
	tracer      tracer.Tracer
	structUtil  struct_util.StructUtil
	errorMapper error_pkg.ErrorMapper
}

func NewHttpResponse(logger logger.Logger, tracer tracer.Tracer, structUtil struct_util.StructUtil, errorMapper error_pkg.ErrorMapper) HttpResponse {
	return &httpResponseImpl{
		logger:      logger,
		tracer:      tracer,
		structUtil:  structUtil,
		errorMapper: errorMapper,
	}
}

func (r *httpResponseImpl) Send(c *gin.Context, data interface{}, err error) {
	ctx, span := r.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	detailedError, _ := r.errorMapper.FromError(err)
	code, response, responseError := r.Build(ctx, http.StatusOK, data, detailedError)
	r.logResponse(ctx, code, response, responseError)

	c.JSON(code, response)
}

func (r *httpResponseImpl) SendWithResponseCode(c *gin.Context, successHttpCode int, data interface{}, err error) {
	ctx, span := r.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	detailedError, _ := r.errorMapper.FromError(err)
	code, response, responseError := r.Build(ctx, successHttpCode, data, detailedError)
	r.logResponse(ctx, code, response, responseError)

	c.JSON(code, response)
}

func (r *httpResponseImpl) Build(ctx context.Context, successHttpCode int, data interface{}, err *error_pkg.ErrorWithDetails) (int, *Response, *HttpError) {
	_, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	var (
		traceId     = span.SpanContext().TraceID().String()
		errMetaData interface{}
		metadata    interface{}
	)

	response := &Response{
		Success:  true,
		TraceId:  traceId,
		Error:    r.structUtil.SetNonPrimitiveDefaultValue(ctx, errMetaData),
		Metadata: r.structUtil.SetNonPrimitiveDefaultValue(ctx, metadata),
		Result:   r.structUtil.SetNonPrimitiveDefaultValue(ctx, data),
	}

	if err != nil {
		response.Success = false
		response.Error = ErrorResponse{
			Code:    err.Code.String(),
			Message: err.Message,
			Errors:  r.structUtil.SetNonPrimitiveDefaultValue(ctx, []validation.ValidationError{}),
		}

		errWithStack := error_pkg.NewErrorWithStack(err, error_pkg.InvalidRequestBodyError)
		httpErr := NewHttpError(errWithStack)

		return err.HttpCode, response, httpErr
	}

	return successHttpCode, response, nil
}

func (r *httpResponseImpl) logResponse(ctx context.Context, httpCode int, response *Response, httpError *HttpError) {
	span := trace.SpanFromContext(ctx)

	fields := []zap.Field{
		zap.String("traceId", span.SpanContext().TraceID().String()),
		zap.Int("statusCode", httpCode),
	}

	if httpError != nil {
		stacks, _ := json.Marshal(httpError.Stack)
		fields = append(
			fields,
			zap.String("original", httpError.Original.Error()),
			zap.String("error", httpError.Error()),
			zap.String("source", httpError.Source),
			zap.String("functionName", httpError.Fn),
			zap.Int("line", httpError.Line),
			zap.String("path", httpError.Path),
			zap.String("stack", string(stacks)),
		)
	}

	var stringResponse string
	if httpCode >= http.StatusInternalServerError || r.logger.Level() == zap.DebugLevel {
		byteResponse, _ := json.Marshal(response)
		stringResponse = string(byteResponse)
		fields = append(fields, zap.Any("response", stringResponse))
	}

	if httpCode >= http.StatusInternalServerError {
		span.SetStatus(codes.Error, stringResponse)
		if httpError != nil {
			span.RecordError(httpError)
		}
		r.logger.WithCtx(ctx).Error("Response error", fields...)
	} else if r.logger.Level() == zap.DebugLevel {
		r.logger.WithCtx(ctx).Debug("Response debug", fields...)
	}
}
