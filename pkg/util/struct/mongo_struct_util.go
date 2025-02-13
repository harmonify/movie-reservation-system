package struct_util

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"github.com/harmonify/movie-reservation-system/pkg/tracer"
	util_shared "github.com/harmonify/movie-reservation-system/pkg/util/shared"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MongoStructUtil is the utility that provides functions to convert struct to bson.D.
type MongoStructUtil interface {
	// ConvertStructToBsonDocument converts a struct to bson.D.
	// The struct must have bson tags. The bson tags are used as the keys in the bson.D.
	// The struct can have nested structs and sql.Null* types.
	// The sql.Null* types are converted to their non-null values if they are valid.
	// Parameters:
	// - ctx: the context
	// - input: the struct to convert
	// - prefix: the prefix to use for the keys in the bson.D (optional). The prefix is used to construct the full key of the fields by joining the prefix and the field name with a dot.
	ConvertStructToBsonDocument(ctx context.Context, input interface{}, prefix ...string) (result bson.D, err error)
}

type mongoStructUtilImpl struct {
	logger logger.Logger
	tracer tracer.Tracer
}

func NewMongoStructUtil(logger logger.Logger, tracer tracer.Tracer) MongoStructUtil {
	return &mongoStructUtilImpl{
		logger: logger,
		tracer: tracer,
	}
}

func (u *mongoStructUtilImpl) ConvertStructToBsonDocument(ctx context.Context, input interface{}, prefix ...string) (result bson.D, err error) {
	result = bson.D{}

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

		tag := typ.Field(i).Tag.Get("bson")
		tagParts := strings.Split(tag, ",")
		fieldName := tagParts[0]

		// Handle nested structs
		if field.Kind() == reflect.Struct {
			nestedPrefix := append(prefix, fieldName)
			nestedFilter, err := u.ConvertStructToBsonDocument(ctx, field.Interface(), nestedPrefix...)
			if err != nil {
				return result, fmt.Errorf("failed to convert nested struct to filter: %w", err)
			}
			result = append(result, bson.E{Key: fieldName, Value: nestedFilter})
			continue
		}

		// Handle sql.Null* types
		fullFieldName := strings.Join(append(prefix, fieldName), ".")
		switch value := field.Interface().(type) {
		case sql.NullBool:
			if value.Valid {
				result = append(result, bson.E{Key: fullFieldName, Value: value.Bool})
			}
		case sql.NullByte:
			if value.Valid {
				result = append(result, bson.E{Key: fullFieldName, Value: value.Byte})
			}
		case sql.NullFloat64:
			if value.Valid {
				result = append(result, bson.E{Key: fullFieldName, Value: value.Float64})
			}
		case sql.NullInt16:
			if value.Valid {
				result = append(result, bson.E{Key: fullFieldName, Value: value.Int16})
			}
		case sql.NullInt32:
			if value.Valid {
				result = append(result, bson.E{Key: fullFieldName, Value: value.Int32})
			}
		case sql.NullInt64:
			if value.Valid {
				result = append(result, bson.E{Key: fullFieldName, Value: value.Int64})
			}
		case sql.NullString:
			if value.Valid {
				result = append(result, bson.E{Key: fullFieldName, Value: value.String})
			}
		case sql.NullTime:
			if value.Valid {
				result = append(result, bson.E{Key: fullFieldName, Value: value.Time})
			}
		default:
			result = append(result, bson.E{Key: fullFieldName, Value: field.Interface()})
		}
	}

	return
}
