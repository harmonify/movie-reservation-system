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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type HttpResponse interface {
	Send(c *gin.Context, data interface{}, err error)
	SendWithResponseCode(c *gin.Context, successHttpCode int, data interface{}, err error)
	// Build is a function to build the response schema
	// Build takes in the ctx, success http code (only respected when err is nil), data, and err
	// Build returns http response code and body
	// You typically don't need to call this function directly
	// Instead, use Send or SendWithResponseCode
	Build(ctx context.Context, successHttpCode int, data interface{}, err *error_pkg.ErrorWithDetails) (responseHttpCode int, responseBody *ResponseBodySchema)
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
	code, response := r.Build(ctx, http.StatusOK, data, detailedError)
	r.logResponse(ctx, code, response, detailedError)

	// If aborted or already written, do not write response
	if c.IsAborted() || c.Writer.Written() {
		return
	}

	c.JSON(code, response)
}

func (r *httpResponseImpl) SendWithResponseCode(c *gin.Context, successHttpCode int, data interface{}, err error) {
	ctx, span := r.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	detailedError, _ := r.errorMapper.FromError(err)
	code, response := r.Build(ctx, successHttpCode, data, detailedError)
	r.logResponse(ctx, code, response, detailedError)

	// If aborted or already written, do not write response
	if c.IsAborted() || c.Writer.Written() {
		return
	}

	c.JSON(code, response)
}

func (r *httpResponseImpl) Build(ctx context.Context, successHttpCode int, data interface{}, err *error_pkg.ErrorWithDetails) (int, *ResponseBodySchema) {
	_, span := r.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	var (
		traceId  = span.SpanContext().TraceID().String()
		metadata interface{}
	)

	response := &ResponseBodySchema{
		Success:  true,
		TraceId:  traceId,
		Error:    &ResponseBodyErrorSchema{},
		Metadata: r.structUtil.SetOrDefault(ctx, metadata),
		Result:   r.structUtil.SetOrDefault(ctx, data),
	}

	if err != nil {
		response.Success = false
		response.Error = &ResponseBodyErrorSchema{
			Code:    err.Code.String(),
			Message: err.Message,
			Errors:  err.Errors,
		}
		return err.HttpCode, response
	}

	return successHttpCode, response
}

func (r *httpResponseImpl) logResponse(ctx context.Context, httpCode int, response *ResponseBodySchema, detailedError *error_pkg.ErrorWithDetails) {
	span := trace.SpanFromContext(ctx)

	fields := []zap.Field{}

	if detailedError != nil {
		fields = append(
			fields,
			zap.Object("error", detailedError),
		)
	}

	if httpCode >= http.StatusInternalServerError || r.logger.Level() == zap.DebugLevel {
		byteResponse, _ := json.Marshal(response)
		fields = append(fields, zap.ByteString("response", byteResponse))
	}

	if httpCode >= http.StatusInternalServerError {
		span.SetStatus(codes.Error, detailedError.Error())
		span.RecordError(detailedError, trace.WithStackTrace(true))
		fields = append(fields, zap.StackSkip("stack", 2))
		r.logger.WithCtx(ctx).Error("Response error", fields...)
	} else if r.logger.Level() == zap.DebugLevel {
		r.logger.WithCtx(ctx).Debug("Response debug", fields...)
	}
}
