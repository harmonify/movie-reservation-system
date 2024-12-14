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
	"github.com/harmonify/movie-reservation-system/pkg/constant"
	logger_shared "github.com/harmonify/movie-reservation-system/pkg/logger/shared"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	struct_util "github.com/harmonify/movie-reservation-system/pkg/util/struct"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

type Response interface {
	Send(c *gin.Context, httpCode int, data interface{}, err error) (int, BaseResponseSchema, error)
	Build(ctx context.Context, httpCode int, data interface{}, err error) (int, BaseResponseSchema, error)
	BuildError(code string, err error) ErrorHandler
	BuildValidationError(code string, err error, errorFields interface{}) ErrorHandler
}

type ResponseImpl struct {
	logger             logger_shared.Logger
	tracer             tracer.Tracer
	structUtil         struct_util.StructUtil
	customHttpErrorMap *constant.CustomHttpErrorMap
}

func NewResponse(logger logger_shared.Logger, tracer tracer.Tracer, structUtil struct_util.StructUtil, customHttpErrorMap *constant.CustomHttpErrorMap) Response {
	return &ResponseImpl{
		logger:             logger,
		tracer:             tracer,
		structUtil:         structUtil,
		customHttpErrorMap: customHttpErrorMap,
	}
}

func (r *ResponseImpl) Send(c *gin.Context, httpCode int, data interface{}, err error) (int, BaseResponseSchema, error) {
	ctx := c.Request.Context()
	_, span := r.tracer.Start(ctx, "Response.Send")
	defer span.End()

	code, response, responseError := r.Build(ctx, httpCode, data, err)

	r.log(ctx, code, response, responseError)

	c.JSON(code, response)

	return code, response, responseError
}

func (r *ResponseImpl) log(ctx context.Context, httpCode int, response BaseResponseSchema, responseError error) {
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

	var respError *ErrorHandlerImpl
	if !errors.As(responseError, &respError) {
		respError = &ErrorHandlerImpl{
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

func (r *ResponseImpl) Build(ctx context.Context, responseCode int, data interface{}, err error) (httpCode int, response BaseResponseSchema, responseError error) {
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

		var responseError *ErrorHandlerImpl
		if !errors.As(err, &responseError) {
			responseError = &ErrorHandlerImpl{
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

func (r *ResponseImpl) BuildError(code string, err error) ErrorHandler {
	source, fn, ln, path, stack := r.getSource(runtime.Caller(1))

	return &ErrorHandlerImpl{
		Code:     code,
		Original: err,
		source:   source,
		fn:       fn,
		line:     ln,
		path:     path,
		stack:    stack,
	}
}

func (r *ResponseImpl) getSource(pc uintptr, file string, line int, ok bool) (source string, fn string, ln int, path string, stack []string) {
	if details := runtime.FuncForPC(pc); details != nil {
		titles := strings.Split(details.Name(), ".")
		fn = titles[len(titles)-1]
	}

	if ok {
		source = fmt.Sprintf("Called from %s, line #%d, func: %v", file, line, fn)
	}

	return source, fn, line, file, r.stackTrace(3)
}

func (r *ResponseImpl) stackTrace(skip int) []string {
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

func (r *ResponseImpl) BuildValidationError(code string, err error, errorFields interface{}) ErrorHandler {
	source, fn, ln, path, stack := r.getSource(runtime.Caller(1))

	return &ErrorHandlerImpl{
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
