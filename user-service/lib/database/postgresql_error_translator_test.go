package database_test

import (
	"errors"
	"testing"

	"github.com/harmonify/movie-reservation-system/user-service/lib/database"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPostgresqlErrorTranslator_Translate(t *testing.T) {
	translator := database.NewPostgresqlErrorTranslator()

	tests := []struct {
		name        string
		inputError  error
		expectedErr error
	}{
		{
			name: "ForeignKeyViolation",
			inputError: &pgconn.PgError{
				Code:           database.ForeignKeyViolation,
				Detail:         "Key (uuid)=(55f39ef0-cc97-4f37-9efc-1d0af0ad89c8) is still referenced from table \"user_keys\".",
				TableName:      "user_keys",
				ConstraintName: "user_keys_fk",
			},
			expectedErr: &database.ForeignKeyViolatedError{
				ReferencedColumnName: "uuid",
				ReferencedValue:      "55f39ef0-cc97-4f37-9efc-1d0af0ad89c8",
				ReferrerTableName:    "user_keys",
				ConstraintName:       "user_keys_fk",
			},
		},
		{
			name: "UniqueViolation",
			inputError: &pgconn.PgError{
				Code:   database.UniqueViolation,
				Detail: "Key (phone_number)=(+6281234567890) already exists.",
			},
			expectedErr: &database.DuplicatedKeyError{
				ColumnName: "phone_number",
				Value:      "+6281234567890",
			},
		},
		{
			name: "InvalidField",
			inputError: &pgconn.PgError{
				Code:    database.InvalidField,
				Message: "column \"expired_atz\" of relation \"user_sessions\" does not exist",
			},
			expectedErr: &database.InvalidFieldError{
				FieldName: "expired_atz",
				TableName: "user_sessions",
			},
		},
		{
			name: "CheckViolation",
			inputError: &pgconn.PgError{
				Code:           database.CheckViolation,
				ConstraintName: "age_positive",
			},
			expectedErr: &database.CheckConstraintViolatedError{
				ConstraintName: "age_positive",
			},
		},
		{
			name:        "RecordNotFound",
			inputError:  gorm.ErrRecordNotFound,
			expectedErr: &database.RecordNotFoundError{},
		},
		{
			name:        "UnmatchedError",
			inputError:  errors.New("generic error"),
			expectedErr: errors.New("generic error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := translator.Translate(test.inputError)

			// Use testify's assert to validate the results
			assert.IsType(t, test.expectedErr, result)

			// Additional field validation for structured errors
			switch e := result.(type) {
			case *database.ForeignKeyViolatedError:
				expected := test.expectedErr.(*database.ForeignKeyViolatedError)
				assert.Equal(t, expected.ReferrerTableName, e.ReferrerTableName)
				assert.Equal(t, expected.ReferencedColumnName, e.ReferencedColumnName)
				assert.Equal(t, expected.ReferencedValue, e.ReferencedValue)
				assert.Equal(t, expected.ConstraintName, e.ConstraintName)
			case *database.DuplicatedKeyError:
				expected := test.expectedErr.(*database.DuplicatedKeyError)
				assert.Equal(t, expected.ColumnName, e.ColumnName)
				assert.Equal(t, expected.Value, e.Value)
			case *database.InvalidFieldError:
				expected := test.expectedErr.(*database.InvalidFieldError)
				assert.Equal(t, expected.FieldName, e.FieldName)
				assert.Equal(t, expected.TableName, e.TableName)
			case *database.CheckConstraintViolatedError:
				expected := test.expectedErr.(*database.CheckConstraintViolatedError)
				assert.Equal(t, expected.ConstraintName, e.ConstraintName)
			case *database.RecordNotFoundError:
				assert.Equal(t, "record not found", e.Error())
			default:
				assert.Equal(t, test.expectedErr, result)
			}
		})
	}
}
