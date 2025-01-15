package struct_util

import (
	"database/sql"
	"reflect"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	util_shared "github.com/harmonify/movie-reservation-system/pkg/util/shared"
	"github.com/iancoleman/strcase"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type StructUtil interface {
	SetValueIfNotEmpty(data any) (result any)
	ConvertSqlStructToMap(input interface{}) (map[string]interface{}, error)
}

type StructUtilParam struct {
	fx.In

	Logger logger.Logger
}

type StructUtilResult struct {
	fx.Out

	StructUtil StructUtil
}

type structUtilImpl struct {
	logger logger.Logger
}

func NewStructUtil(p StructUtilParam) StructUtilResult {
	return StructUtilResult{
		StructUtil: &structUtilImpl{
			logger: p.Logger,
		},
	}
}

func (u *structUtilImpl) SetValueIfNotEmpty(data any) (result any) {
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

// ConvertSqlStructToMap converts a struct with SQL nullable types to a map with snake_case keys,
// excluding fields with empty values.
func (u *structUtilImpl) ConvertSqlStructToMap(input interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	val := reflect.ValueOf(input)

	// Ensure the input is a struct or a pointer to struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return result, &util_shared.UtilInvalidError{
			Params: []util_shared.InvalidParam{
				{
					Key:    "input",
					Value:  input,
					Reason: "input must be a struct or a pointer to struct",
				},
			},
		}
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Convert field name to snake_case
		snakeCaseKey := strcase.ToSnake(fieldType.Name)

		// Extract value based on SQL nullable type
		switch value := field.Interface().(type) {
		case sql.NullBool:
			if value.Valid {
				result[snakeCaseKey] = value.Bool
			}
		case sql.NullByte:
			if value.Valid {
				result[snakeCaseKey] = value.Byte
			}
		case sql.NullFloat64:
			if value.Valid {
				result[snakeCaseKey] = value.Float64
			}
		case sql.NullInt16:
			if value.Valid {
				result[snakeCaseKey] = value.Int16
			}
		case sql.NullInt32:
			if value.Valid {
				result[snakeCaseKey] = value.Int32
			}
		case sql.NullInt64:
			if value.Valid {
				result[snakeCaseKey] = value.Int64
			}
		case sql.NullString:
			if value.Valid {
				result[snakeCaseKey] = value.String
			}
		case sql.NullTime:
			if value.Valid {
				result[snakeCaseKey] = value.Time
			}
		case gorm.DeletedAt:
			if value.Valid {
				result[snakeCaseKey] = value.Time
			}
		}
	}
	return result, nil
}
