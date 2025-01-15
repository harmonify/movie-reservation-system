package util_shared

var (
	UtilInvalidParam = "UTIL_INVALID_PARAM"
)

type (
	UtilInvalidError struct {
		Params []InvalidParam
	}

	InvalidParam struct {
		Key    string
		Value  any
		Reason string
	}
)

func (e *UtilInvalidError) Error() string {
	return UtilInvalidParam
}
