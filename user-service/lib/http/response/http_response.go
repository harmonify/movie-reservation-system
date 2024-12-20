package response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	constant "github.com/harmonify/movie-reservation-system/user-service/lib/http/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	struct_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/struct"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

type HttpResponse interface {
	Send(c *gin.Context, data interface{}, err error) (int, BaseResponseSchema)
	SendWithResponseCode(c *gin.Context, httpCode int, data interface{}, err error) (int, BaseResponseSchema)
	Build(ctx context.Context, httpCode int, data interface{}, err error) (int, BaseResponseSchema, error)
	BuildError(code string, err error) HttpErrorHandler
	BuildValidationError(code string, err error, errorFields interface{}) HttpErrorHandler
}

type httpResponseImpl struct {
	logger             logger.Logger
	tracer             tracer.Tracer
	structUtil         struct_util.StructUtil
	customHttpErrorMap *constant.CustomHttpErrorMap
}

func NewHttpResponse(logger logger.Logger, tracer tracer.Tracer, structUtil struct_util.StructUtil, customHttpErrorMap *constant.CustomHttpErrorMap) HttpResponse {
	return &httpResponseImpl{
		logger:             logger,
		tracer:             tracer,
		structUtil:         structUtil,
		customHttpErrorMap: customHttpErrorMap,
	}
}

func (r *httpResponseImpl) Send(c *gin.Context, data interface{}, err error) (int, BaseResponseSchema) {
	ctx := c.Request.Context()
	_, span := r.tracer.Start(ctx, "Response.Send")
	defer span.End()

	return r.SendWithResponseCode(c, http.StatusOK, data, err)
}

func (r *httpResponseImpl) SendWithResponseCode(c *gin.Context, httpCode int, data interface{}, err error) (int, BaseResponseSchema) {
	ctx := c.Request.Context()
	_, span := r.tracer.Start(ctx, "Response.SendWithResponseCode")
	defer span.End()

	code, response, responseError := r.Build(ctx, httpCode, data, err)

	r.log(ctx, code, response, responseError)

	c.JSON(code, response)

	return code, response
}

func (r *httpResponseImpl) log(ctx context.Context, httpCode int, response BaseResponseSchema, responseError error) {
	span := trace.SpanFromContext(ctx)

	fields := []zapcore.Field{
		{
			Key:    "traceId",
			Type:   zapcore.StringType,
			String: span.SpanContext().TraceID().String(),
		},
		{
			Key:     "statusCode",
			Type:    zapcore.Int64Type,
			Integer: int64(httpCode),
		},
	}

	var respError *HttpErrorHandlerImpl
	if !errors.As(responseError, &respError) {
		respError = &HttpErrorHandlerImpl{
			Code:     constant.InternalServerError,
			Original: responseError,
		}
	}
	if respError.Original != nil {
		stacks, _ := json.Marshal(respError.stack)
		fields = append(fields, zapcore.Field{
			Key:    "source",
			Type:   zapcore.StringType,
			String: respError.source,
		}, zapcore.Field{
			Key:    "functionName",
			Type:   zapcore.StringType,
			String: respError.fn,
		}, zapcore.Field{
			Key:     "line",
			Type:    zapcore.Int64Type,
			Integer: int64(respError.line),
		}, zapcore.Field{
			Key:    "path",
			Type:   zapcore.StringType,
			String: respError.path,
		}, zapcore.Field{
			Key:    "stack",
			Type:   zapcore.StringType,
			String: string(stacks),
		})
	}

	// TODO: add conditionals parsing this to improve perf
	byteResponse, _ := json.Marshal(response)
	stringResponse := string(byteResponse)

	if httpCode >= http.StatusInternalServerError {
		span.SetStatus(codes.Error, stringResponse)
		span.RecordError(responseError)
	}

	r.logger.Debug(stringResponse, fields...)
}

func (r *httpResponseImpl) Build(ctx context.Context, responseCode int, data interface{}, err error) (httpCode int, response BaseResponseSchema, responseError error) {
	_, span := r.tracer.Start(ctx, "Response.Build")
	defer span.End()

	var (
		traceId     = span.SpanContext().TraceID().String()
		errMetaData interface{}
		metadata    interface{}
	)

	response = BaseResponseSchema{
		Success:  true,
		TraceId:  traceId,
		Error:    r.structUtil.SetValueIfNotEmpty(errMetaData),
		Metadata: r.structUtil.SetValueIfNotEmpty(metadata),
		Result:   r.structUtil.SetValueIfNotEmpty(data),
	}

	if err != nil {
		response.Success = false

		var responseError *HttpErrorHandlerImpl
		if !errors.As(err, &responseError) {
			responseError = &HttpErrorHandlerImpl{
				Code:     constant.InternalServerError,
				Original: err,
			}
		}

		if _, ok := (*r.customHttpErrorMap)[responseError.Code]; !ok {
			(*r.customHttpErrorMap)[responseError.Code] = constant.DefaultCustomHttpErrorMap[constant.InternalServerError]
		}

		if responseError.Errors == nil {
			responseError.Errors = r.structUtil.SetValueIfNotEmpty([]BaseErrorValidationSchema{})
		}

		response.Error = BaseErrorResponseSchema{
			Code:    responseError.Code,
			Message: (*r.customHttpErrorMap)[responseError.Code].Message,
			Errors:  responseError.Errors,
		}

		return (*r.customHttpErrorMap)[responseError.Code].HttpCode, response, responseError
	}

	return responseCode, response, nil
}

func (r *httpResponseImpl) BuildError(code string, err error) HttpErrorHandler {
	source, fn, ln, path, stack := r.getSource(runtime.Caller(1))

	return &HttpErrorHandlerImpl{
		Code:     code,
		Original: err,
		source:   source,
		fn:       fn,
		line:     ln,
		path:     path,
		stack:    stack,
	}
}

func (r *httpResponseImpl) getSource(pc uintptr, file string, line int, ok bool) (source string, fn string, ln int, path string, stack []string) {
	if details := runtime.FuncForPC(pc); details != nil {
		titles := strings.Split(details.Name(), ".")
		fn = titles[len(titles)-1]
	}

	if ok {
		source = fmt.Sprintf("Called from %s, line #%d, func: %v", file, line, fn)
	}

	return source, fn, line, file, r.stackTrace(3)
}

func (r *httpResponseImpl) stackTrace(skip int) []string {
	var stacks []string
	for {
		pc, path, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)

		stacks = append(stacks, fmt.Sprintf("%s:%d %s()", path, line, fn.Name()))
		skip++
	}

	return stacks
}

func (r *httpResponseImpl) BuildValidationError(code string, err error, errorFields interface{}) HttpErrorHandler {
	source, fn, ln, path, stack := r.getSource(runtime.Caller(1))

	return &HttpErrorHandlerImpl{
		Code:     code,
		Original: err,
		Errors:   errorFields,
		source:   source,
		fn:       fn,
		line:     ln,
		path:     path,
		stack:    stack,
	}
}
