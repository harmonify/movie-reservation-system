package response

type HttpErrorHandler interface {
	Error() string
}

type HttpErrorHandlerImpl struct {
	Code     string
	Original error
	Errors   interface{}
	source   string
	fn       string
	line     int
	path     string
	stack    []string
}

func NewHttpErrorHandler() HttpErrorHandler {
	return &HttpErrorHandlerImpl{}
}

func (e *HttpErrorHandlerImpl) Error() string {
	return e.Original.Error()
}
