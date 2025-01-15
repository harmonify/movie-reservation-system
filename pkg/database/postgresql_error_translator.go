package database

import (
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type PostgresqlErrorTranslator interface {
	Translate(err error) error
}

type postgresqlErrorTranslatorImpl struct{}

func NewPostgresqlErrorTranslator() PostgresqlErrorTranslator {
	return &postgresqlErrorTranslatorImpl{}
}

// Translate it will translate the error to custom errors.
// Since currently gorm supporting both pgx and pg drivers, only checking for pgx PgError types is not enough for translating errors, so we have additional error json marshal fallback.
func (t *postgresqlErrorTranslatorImpl) Translate(err error) error {
	if err == nil {
		return nil
	}

	// Translate gorm generic error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewRecordNotFoundError(err)
	}

	// pgx driver -----
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if translatedErr, found := t.translateErrorCode(pgErr, pgErr.Code); found {
			return translatedErr
		}
		return err
	}
	// pgx driver ----- end

	// pg driver -----
	parsedErr, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		return err
	}

	var errMsg ErrMessage
	unmarshalErr := json.Unmarshal(parsedErr, &errMsg)
	if unmarshalErr != nil {
		return err
	}
	errMsg.Original = err

	if translatedErr, found := t.translateErrorCode(errMsg, errMsg.Code); found {
		return translatedErr
	}
	return err
	// pg driver ----- end
}

func (t *postgresqlErrorTranslatorImpl) translateErrorCode(err error, code string) (error, bool) {
	switch code {
	case ForeignKeyViolation:
		{
			return NewForeignKeyViolatedError(err), true
		}
	case UniqueViolation:
		{
			return NewDuplicatedKeyError(err), true
		}
	case CheckViolation:
		{
			return NewCheckConstraintViolatedError(err), true
		}
	case InvalidField:
		{
			return NewInvalidFieldError(err), true
		}
	default:
		{
			return err, false
		}
	}
}
