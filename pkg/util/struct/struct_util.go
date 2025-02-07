package struct_util

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	util_shared "github.com/harmonify/movie-reservation-system/pkg/util/shared"
	"github.com/iancoleman/strcase"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gorm.io/gorm"
)

type StructUtil interface {
	// SetNonPrimitiveDefaultValue sets the default value for non-primitive data types.
	// It returns nil if the data is nil or zero value.
	// It returns the data if the data is not nil or zero value.
	SetNonPrimitiveDefaultValue(ctx context.Context, data interface{}) (result interface{})

	// ConvertSqlStructToMap converts a struct with SQL nullable types to a map with snake_case keys,
	// excluding invalid nullable types.
	ConvertSqlStructToMap(ctx context.Context, input interface{}) (map[string]interface{}, error)

	// ConvertProtoToMap converts a proto message to a map with camelCase keys.
	ConvertProtoToMap(ctx context.Context, msg proto.Message) (map[string]interface{}, error)
}

type StructUtilParam struct {
	fx.In

	Logger logger.Logger
	Tracer tracer.Tracer
}

type StructUtilResult struct {
	fx.Out

	StructUtil StructUtil
}

type structUtilImpl struct {
	logger logger.Logger
	tracer tracer.Tracer
}

func NewStructUtil(p StructUtilParam) StructUtilResult {
	return StructUtilResult{
		StructUtil: &structUtilImpl{
			logger: p.Logger,
			tracer: p.Tracer,
		},
	}
}

func (u *structUtilImpl) SetNonPrimitiveDefaultValue(ctx context.Context, data any) (result any) {
	ctx, span := u.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if ok {
				u.logger.WithCtx(ctx).Error(fmt.Sprintf("Recovered from panic: %v", err))
			} else {
				u.logger.WithCtx(ctx).Error(fmt.Sprintf("Recovered from panic: %v", r))
			}
			result = nil
		}
	}()

	val := reflect.ValueOf(data)
	// Cast the data to non-pointer if it is a pointer
	if val.Kind() == reflect.Pointer || val.Kind() == reflect.Interface {
		val = val.Elem()
	}
	if val.Kind() != reflect.Invalid && !val.IsZero() {
		return data
	}

	rt := reflect.TypeOf(data)
	if rt == nil {
		result = make(map[string]interface{}, 0)
		return
	}
	// Cast the data to non-pointer if it is a pointer
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	// Get default value based on the data type
	switch expression := rt.Kind(); expression {
	case reflect.Interface:
		result = make(map[string]interface{}, 0)
	case reflect.Struct:
		result = make(map[string]interface{}, 0)
	case reflect.Map:
		result = make(map[string]interface{}, 0)
	case reflect.Slice:
		result = make([]interface{}, 0)
	case reflect.Array:
		result = make([]interface{}, 0)
	case reflect.Invalid:
		result = nil
	case reflect.Chan:
		result = nil
	case reflect.Func:
		result = nil
	case reflect.Uintptr:
		result = nil
	case reflect.UnsafePointer:
		result = nil
	default:
		result = nil
	}

	return
}

func (u *structUtilImpl) ConvertSqlStructToMap(ctx context.Context, input interface{}) (result map[string]interface{}, err error) {
	ctx, span := u.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(error)
			if ok {
				u.logger.WithCtx(ctx).Error(fmt.Sprintf("Recovered from panic: %v", e))
				err = e
			} else {
				u.logger.WithCtx(ctx).Error(fmt.Sprintf("Recovered from panic: %v", r))
				err = fmt.Errorf("recovered from panic: %v", r)
			}
		}
	}()

	result = make(map[string]interface{})

	if input == nil {
		return
	}

	// Ensure the input is a struct
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		err = &util_shared.UtilInvalidError{
			Params: []util_shared.InvalidParam{
				{
					Key:    "input",
					Value:  input,
					Reason: "input must be a struct or a pointer to struct",
				},
			},
		}
		return
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

	return
}

func (u *structUtilImpl) ConvertProtoToMap(ctx context.Context, msg proto.Message) (result map[string]interface{}, err error) {
	ctx, span := u.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(error)
			if ok {
				u.logger.WithCtx(ctx).Error(fmt.Sprintf("Recovered from panic: %v", e))
				err = e
			} else {
				u.logger.WithCtx(ctx).Error(fmt.Sprintf("Recovered from panic: %v", r))
				err = fmt.Errorf("recovered from panic: %v", r)
			}
		}
	}()

	result = make(map[string]interface{})

	if msg == nil {
		return
	}

	msgReflect := msg.ProtoReflect()

	msgReflect.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		name := string(fd.Name())

		// Handle nested message
		if fd.Kind() == protoreflect.MessageKind {
			if fd.IsList() { // Handle repeated messages
				list := v.List()
				arr := make([]interface{}, list.Len())
				for i := 0; i < list.Len(); i++ {
					res, err := u.ConvertProtoToMap(ctx, list.Get(i).Message().Interface())
					if err != nil {
						err = fmt.Errorf("failed to convert proto to map: %w", err)
						return false
					}
					arr[i] = res
				}
				result[name] = arr
			} else if fd.IsMap() { // Handle map fields
				mapValues := v.Map()
				mapResult := make(map[string]interface{})
				mapValues.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
					res, err := u.ConvertProtoToMap(ctx, v.Message().Interface())
					if err != nil {
						err = fmt.Errorf("failed to convert proto to map: %w", err)
						return false
					}
					mapResult[k.String()] = res
					return true
				})
				result[name] = mapResult
			} else { // Handle single nested message
				res, err := u.ConvertProtoToMap(ctx, v.Message().Interface())
				if err != nil {
					err = fmt.Errorf("failed to convert proto to map: %w", err)
					return false
				}
				result[name] = res
			}
		} else if fd.IsList() { // Handle repeated primitive fields
			list := v.List()
			arr := make([]interface{}, list.Len())
			for i := 0; i < list.Len(); i++ {
				arr[i] = list.Get(i).Interface()
			}
			result[name] = arr
		} else if fd.IsMap() { // Handle map of primitive types
			mapValues := v.Map()
			mapResult := make(map[string]interface{})
			mapValues.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
				mapResult[k.String()] = v.Interface()
				return true
			})
			result[name] = mapResult
		} else { // Handle primitive fields
			result[name] = v.Interface()
		}
		return true
	})

	return
}
