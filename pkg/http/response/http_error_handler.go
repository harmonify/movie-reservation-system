package response

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

func NewHttpErrorHandler() error {
	return &HttpErrorHandlerImpl{}
}

func (e HttpErrorHandlerImpl) Error() string {
	return e.Original.Error()
}
