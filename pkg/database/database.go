package database

import (
	"github.com/go-playground/validator/v10"
	"github.com/harmonify/movie-reservation-system/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DatabaseParam struct {
	fx.In

	Logger logger.Logger
}

type DatabaseResult struct {
	fx.Out

	Database                  *Database
	PostgresqlErrorTranslator PostgresqlErrorTranslator
}

type DatabaseConfig struct {
	Env                   string `validate:"required,oneof=dev test prod"`
	DbType                string `validate:"required,oneof=postgresql mysql"`
	DbHost                string `validate:"required"`
	DbPort                int    `validate:"required,min=1,max=65535"`
	DbUser                string `validate:"required"`
	DbPassword            string `validate:"required"`
	DbName                string `validate:"required"`
	DbMaxIdleConn         int    `validate:"required,min=1"`
	DbMaxOpenConn         int    `validate:"required,min=1"`
	DbMaxLifetimeInMinute int    `validate:"required,min=1"`
}

func NewDatabase(p DatabaseParam, cfg *DatabaseConfig) (DatabaseResult, error) {
	err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg)
	if err != nil {
		return DatabaseResult{}, err
	}

	var db *Database
	switch cfg.DbType {
	case "mysql":
		db, err = NewMysqlDatabase(p, cfg)
	default:
		db, err = NewPostgresqlDatabase(p, cfg)
	}

	if err != nil {
		return DatabaseResult{}, err
	}

	return DatabaseResult{
		Database:                  db,
		PostgresqlErrorTranslator: NewPostgresqlErrorTranslator(),
	}, nil
}

type Database struct {
	DB     *gorm.DB
	Logger logger.Logger
}

func (d *Database) WithTx(tx *Transaction) *Database {
	if tx == nil {
		return d
	} else {
		return &Database{
			DB:     tx.DB,
			Logger: d.Logger,
		}
	}
}

func (d *Database) Transaction(fc func(tx *Transaction) error) error {
	err := d.DB.Transaction(func(_tx *gorm.DB) error {
		tx := NewTransaction(_tx)
		return fc(tx)
	})

	return err
}
