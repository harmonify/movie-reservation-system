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
	Build(ctx context.Context, successHttpCode int, data interface{}, err *error_pkg.ErrorWithDetails) (responseHttpCode int, responseBody *Response)
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
	c.JSON(code, response)
}

func (r *httpResponseImpl) SendWithResponseCode(c *gin.Context, successHttpCode int, data interface{}, err error) {
	ctx, span := r.tracer.StartSpanWithCaller(c.Request.Context())
	defer span.End()

	detailedError, _ := r.errorMapper.FromError(err)
	code, response := r.Build(ctx, successHttpCode, data, detailedError)
	r.logResponse(ctx, code, response, detailedError)
	c.JSON(code, response)
}

func (r *httpResponseImpl) Build(ctx context.Context, successHttpCode int, data interface{}, err *error_pkg.ErrorWithDetails) (int, *Response) {
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
			Errors:  r.structUtil.SetNonPrimitiveDefaultValue(ctx, err.Errors),
		}
		return err.HttpCode, response
	}

	return successHttpCode, response
}

func (r *httpResponseImpl) logResponse(ctx context.Context, httpCode int, response *Response, detailedError *error_pkg.ErrorWithDetails) {
	fields := []zap.Field{}

	if detailedError != nil {
		fields = append(
			fields,
			zap.Object("error", detailedError),
		)
	}

	var stringResponse string
	if httpCode >= http.StatusInternalServerError || r.logger.Level() == zap.DebugLevel {
		byteResponse, _ := json.Marshal(response)
		stringResponse = string(byteResponse)
		fields = append(fields, zap.Any("response", stringResponse))
	}

	if httpCode >= http.StatusInternalServerError {
		r.logger.WithCtx(ctx).Error("Response error", fields...)
	} else if r.logger.Level() == zap.DebugLevel {
		r.logger.WithCtx(ctx).Debug("Response debug", fields...)
	}
}
