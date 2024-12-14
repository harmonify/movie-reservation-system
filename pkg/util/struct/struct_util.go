package struct_util

import (
	"reflect"
)

type StructUtil interface {
	SetValueIfNotEmpty(data any) (result any)
}
type StructImpl struct{}

func NewStructUtil() StructUtil {
	return &StructImpl{}
}

func (s *StructImpl) SetValueIfNotEmpty(data any) (result any) {
	if data == nil {
		return struct{}{}
	}

	rt := reflect.TypeOf(data)

	if !reflect.ValueOf(data).IsZero() {
		return data
	}

	if rt != nil {
		switch expression := rt.Kind(); expression {
		case reflect.Struct:
			result = make(map[string]interface{}, 0)
		case reflect.Ptr:
			dataInterface := reflect.New(rt.Elem()).Interface()
			dataKind := reflect.TypeOf(dataInterface).Kind()
			switch dataKind {
			case reflect.Struct:
				result = make(map[string]interface{}, 0)
			default:
				result = make([]interface{}, 0)
			}
		default:
			result = make([]interface{}, 0)
		}
	} else {
		result = make(map[string]interface{}, 0)
	}

	return
}
