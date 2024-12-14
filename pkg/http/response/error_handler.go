package response

type ErrorHandler interface {
	Error() string
}

type ErrorHandlerImpl struct {
	Code     string
	Original error
	Errors   interface{}
	source   string
	fn       string
	line     int
	path     string
	stack    []string
}

func NewErrorHandler() ErrorHandler {
	return &ErrorHandlerImpl{}
}

func (e *ErrorHandlerImpl) Error() string {
	return e.Original.Error()
}
