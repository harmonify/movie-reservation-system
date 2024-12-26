package database

import (
	"regexp"

	"github.com/jackc/pgx/v5/pgconn"
)

// https://www.postgresql.org/docs/current/errcodes-appendix.html
var (
	// Class 22 — Data Exception
	// InvalidTextRepresentation = "22P02" // occurs when trying to convert a string to a data type that is not compatible with the string's format.

	// Class 23 — Integrity Constraint Violation
	NotNullViolation    = "23502"
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
	CheckViolation      = "23514"
	// Class 40 — Transaction Rollback
	// TransactionRollback                     = "40000"
	// TransactionIntegrityConstraintViolation = "40002"
	// DeadlockDetected                        = "40P01"

	// Class 42 — Syntax Error or Access Rule Violation
	// SyntaxError = "42601"
	// InsufficientPrivilege = "42501"

	InvalidField = "42703"
)

// ------------------------

type ErrMessage struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func (e ErrMessage) Error() string {
	return e.Message
}

// ------------------------

// RecordNotFoundError record not found error
type RecordNotFoundError struct {
	Original error
}

func (e RecordNotFoundError) Error() string {
	return "record not found"
}

func NewRecordNotFoundError(err error) *RecordNotFoundError {
	return &RecordNotFoundError{
		Original: err,
	}
}

// ------------------------

// DuplicatedKeyError occurs when there is a unique key constraint violation
type DuplicatedKeyError struct {
	Original   error           `json:"original"`
	PgError    *pgconn.PgError `json:"pg_error"`
	ColumnName string          `json:"column_name"` // available for pgx driver
	Value      string          `json:"value"`       // available for pgx driver
}

func (e DuplicatedKeyError) Error() string {
	return "duplicated key not allowed"
}

func NewDuplicatedKeyError(err error) *DuplicatedKeyError {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		// Compile the regex
		// Example PgError.Detail: "Key (phone_number)=(+6281234567890) already exists."
		re := regexp.MustCompile(`Key \((.*?)\)=\((.*?)\)`)

		// Find matches
		matches := re.FindStringSubmatch(pgErr.Detail)

		// Check if we have the required matches
		if len(matches) > 2 {
			return &DuplicatedKeyError{
				Original:   err,
				PgError:    pgErr,
				ColumnName: matches[1],
				Value:      matches[2],
			}
		}
	}

	return &DuplicatedKeyError{
		Original: err,
	}
}

// ------------------------

// ForeignKeyViolatedError occurs when there is a unique key constraint violation
type ForeignKeyViolatedError struct {
	Original             error           `json:"original"`
	PgError              *pgconn.PgError `json:"pg_error"`
	ReferrerTableName    string          `json:"referrer_table_name"`
	ReferencedColumnName string          `json:"referenced_column_name"`
	ReferencedValue      string          `json:"referenced_value"`
	ConstraintName       string          `json:"constraint_name"`
}

func (e ForeignKeyViolatedError) Error() string {
	return "violates foreign key constraint"
}

func NewForeignKeyViolatedError(err error) *ForeignKeyViolatedError {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		// Compile the regex
		// "Key (uuid)=(55f39ef0-cc97-4f37-9efc-1d0af0ad89c8) is still referenced from table \"user_keys\"."
		re := regexp.MustCompile(`Key \((.*?)\)=\((.*?)\)`)

		// Find matches
		matches := re.FindStringSubmatch(pgErr.Detail)

		// Check if we have the required matches
		if len(matches) > 2 {
			return &ForeignKeyViolatedError{
				Original:             err,
				PgError:              pgErr,
				ReferrerTableName:    pgErr.TableName,
				ReferencedColumnName: matches[1],
				ReferencedValue:      matches[2],
				ConstraintName:       pgErr.ConstraintName,
			}
		}
	}

	return &ForeignKeyViolatedError{
		Original: err,
	}
}

// ------------------------

// InvalidFieldError occurs when there is an invalid field
type InvalidFieldError struct {
	Original  error           `json:"original"`
	PgError   *pgconn.PgError `json:"pg_error"`
	FieldName string          `json:"field_name"` // invalid column/field
	TableName string          `json:"table_name"`
}

func (e InvalidFieldError) Error() string {
	return "invalid field"
}

func NewInvalidFieldError(err error) *InvalidFieldError {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		// Compile the regex
		// "column \"expired_atz\" of relation \"user_sessions\" does not exist"
		re := regexp.MustCompile(`column "(.*?)" of relation "(.*?)"`)

		// Find matches
		matches := re.FindStringSubmatch(pgErr.Message)

		// Check if we have the required matches
		if len(matches) > 2 {
			return &InvalidFieldError{
				Original:  err,
				PgError:   pgErr,
				FieldName: matches[1],
				TableName: matches[2],
			}
		}
	}

	return &InvalidFieldError{
		Original: err,
	}
}

// ------------------------

// CheckConstraintViolatedError occurs when there is a check constraint violation
type CheckConstraintViolatedError struct {
	Original       error           `json:"original"`
	PgError        *pgconn.PgError `json:"pg_error"`
	ConstraintName string          `json:"constraint_name"`
}

func (e CheckConstraintViolatedError) Error() string {
	return "violates check constraint"
}

func NewCheckConstraintViolatedError(err error) *CheckConstraintViolatedError {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return &CheckConstraintViolatedError{
			Original:       err,
			PgError:        pgErr,
			ConstraintName: pgErr.ConstraintName,
		}
	}

	return &CheckConstraintViolatedError{
		Original: err,
	}
}

// ------------------------
