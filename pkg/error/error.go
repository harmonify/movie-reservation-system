package error_pkg

import (
	"fmt"
	"runtime"
	"strings"

	"google.golang.org/grpc/codes"
)

type ErrorCode string

func (e ErrorCode) String() string {
	return string(e)
}

// ErrorWithDetails is a custom error type that implements the error interface.
type ErrorWithDetails struct {
	Code     ErrorCode   `json:"code"`      // error code
	HttpCode int         `json:"http_code"` // associated HTTP status code
	GrpcCode codes.Code  `json:"grpc_code"` // associated gRPC status code
	Message  string      `json:"message"`   // user-friendly message
	Data     interface{} `json:"data"`      // additional data
}

// Error returns the user-friendly error message.
func (e *ErrorWithDetails) Error() string {
	return e.Message
}

func (e *ErrorWithDetails) As(target interface{}) bool {
	_, ok := target.(*ErrorWithDetails)
	return ok
}

type ErrorWithStack struct {
	*ErrorWithDetails
	Original error    `json:"original"`
	Source   string   `json:"source"`
	Fn       string   `json:"fn"`
	Line     int      `json:"line"`
	Path     string   `json:"path"`
	Stack    []string `json:"stack"`
}

func NewErrorWithStack(original error, details *ErrorWithDetails) *ErrorWithStack {
	source, fn, ln, path, stack := getSource(runtime.Caller(1))
	return &ErrorWithStack{
		ErrorWithDetails: details,
		Original:         original,
		Source:           source,
		Fn:               fn,
		Line:             ln,
		Path:             path,
		Stack:            stack,
	}
}

func (e *ErrorWithStack) Error() string {
	return e.ErrorWithDetails.Error()
}

func (e *ErrorWithStack) As(target interface{}) bool {
	return e.ErrorWithDetails.As(target)
}

func (e *ErrorWithStack) Unwrap() error {
	return e.ErrorWithDetails
}

func getSource(pc uintptr, file string, line int, ok bool) (source string, fn string, ln int, path string, stack []string) {
	if details := runtime.FuncForPC(pc); details != nil {
		titles := strings.Split(details.Name(), ".")
		fn = titles[len(titles)-1]
	}

	if ok {
		source = fmt.Sprintf("Called from %s, line #%d, func: %v", file, line, fn)
	}

	return source, fn, line, file, stackTrace(0)
}

func stackTrace(skip int) []string {
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
