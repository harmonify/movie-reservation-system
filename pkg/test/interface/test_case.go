package test_interface

import (
	"net/http"
	"net/http/httptest"
)

type (
	TestCase[Config, Expectation any] struct {
		Description string
		Config      Config
		Expectation Expectation
		BeforeCall  func(config Config)
		AfterCall   func()
	}

	HttpTestCase[RequestBody, ResponseBody any] struct {
		Description string
		Config      Request[RequestBody]
		Expectation ResponseExpectation[ResponseBody]
		BeforeCall  func(req *http.Request)
		AfterCall   func(w *httptest.ResponseRecorder)
	}

	Request[RequestBody any] struct {
		RequestHeader []RequestHeaderConfig
		RequestQuery  []RequestQueryConfig
		RequestBody   RequestBody
	}
	RequestHeaderConfig struct {
		Key   string
		Value string
	}
	RequestQueryConfig struct {
		Key   string
		Value string
	}

	ResponseExpectation[ResponseBody any] struct {
		ResponseStatusCode       NullInt
		ResponseBodyStatus       NullBool
		ResponseBodyResult       ResponseBody
		ResponseBodyErrorCode    NullString
		ResponseBodyErrorMessage NullString
		ResponseBodyErrorObject  []interface{}
	}

	// similar to [database/sql#NullInt]
	NullInt struct {
		Int   int
		Valid bool
	}

	// similar to [database/sql#NullString]
	NullString struct {
		String string
		Valid  bool
	}
	// similar to [database/sql#NullBool]
	NullBool struct {
		Bool  bool
		Valid bool // Valid should be set to `true` if `Bool` is not null
	}
)
