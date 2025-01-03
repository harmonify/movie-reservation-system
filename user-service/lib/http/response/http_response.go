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
	error_constant "github.com/harmonify/movie-reservation-system/user-service/lib/error/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/harmonify/movie-reservation-system/user-service/lib/tracer"
	struct_util "github.com/harmonify/movie-reservation-system/user-service/lib/util/struct"
	"github.com/harmonify/movie-reservation-system/user-service/lib/util/validation"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

type HttpResponse interface {
	Send(c *gin.Context, data interface{}, err error)
	SendWithResponseCode(c *gin.Context, httpCode int, data interface{}, err error)
	Build(ctx context.Context, httpCode int, data interface{}, err error) (int, BaseResponseSchema, error)
	BuildError(code string, err error) *HttpErrorHandlerImpl
	BuildValidationError(code string, err error, errorFields []validation.BaseValidationErrorSchema) *HttpErrorHandlerImpl
}

type httpResponseImpl struct {
	logger         logger.Logger
	tracer         tracer.Tracer
	structUtil     struct_util.StructUtil
	customErrorMap *error_constant.CustomErrorMap
}

func NewHttpResponse(logger logger.Logger, tracer tracer.Tracer, structUtil struct_util.StructUtil, customHttpErrorMap *error_constant.CustomErrorMap) HttpResponse {
	return &httpResponseImpl{
		logger:         logger,
		tracer:         tracer,
		structUtil:     structUtil,
		customErrorMap: customHttpErrorMap,
	}
}

func (r *httpResponseImpl) Send(c *gin.Context, data interface{}, err error) {
	ctx := c.Request.Context()
	_, span := r.tracer.Start(ctx, "Response.Send")
	defer span.End()

	code, response, responseError := r.Build(ctx, http.StatusOK, data, err)

	r.logResponse(ctx, code, response, responseError)

	c.JSON(code, response)
}

func (r *httpResponseImpl) SendWithResponseCode(c *gin.Context, httpCode int, data interface{}, err error) {
	ctx := c.Request.Context()
	_, span := r.tracer.Start(ctx, "Response.SendWithResponseCode")
	defer span.End()

	code, response, responseError := r.Build(ctx, httpCode, data, err)

	r.logResponse(ctx, code, response, responseError)

	c.JSON(code, response)
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
			// responseError = &HttpErrorHandlerImpl{
			// 	Code:     constant.InternalServerError,
			// 	Original: err,
			// }
			responseError = r.buildErrorV2(err.Error(), err)
		}

		if _, ok := (*r.customErrorMap)[responseError.Code]; !ok {
			(*r.customErrorMap)[responseError.Code] = error_constant.DefaultCustomErrorMap[error_constant.InternalServerError]
		}

		if responseError.Errors == nil {
			responseError.Errors = r.structUtil.SetValueIfNotEmpty([]validation.BaseValidationErrorSchema{})
		}

		response.Error = BaseErrorResponseSchema{
			Code:    responseError.Code,
			Message: (*r.customErrorMap)[responseError.Code].Message,
			Errors:  responseError.Errors,
		}

		return (*r.customErrorMap)[responseError.Code].HttpCode, response, responseError
	}

	return responseCode, response, nil
}

func (r *httpResponseImpl) logResponse(ctx context.Context, httpCode int, response BaseResponseSchema, responseError error) {
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
			Code:     error_constant.InternalServerError,
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

	if httpCode >= http.StatusInternalServerError || r.logger.Level() == zapcore.DebugLevel {
		byteResponse, _ := json.Marshal(response)
		stringResponse := string(byteResponse)

		if httpCode >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, stringResponse)
			span.RecordError(responseError)
		}

		r.logger.WithCtx(ctx).Debug(stringResponse, fields...)
	}
}

func (r *httpResponseImpl) BuildError(code string, err error) *HttpErrorHandlerImpl {
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

func (r *httpResponseImpl) BuildValidationError(code string, err error, errorFields []validation.BaseValidationErrorSchema) *HttpErrorHandlerImpl {
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

func (r *httpResponseImpl) buildErrorV2(code string, err error) *HttpErrorHandlerImpl {
	source, fn, ln, path, stack := r.getSource(runtime.Caller(3))

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
