package database

import (
	"context"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Database struct {
	DB     *gorm.DB
	Logger logger.Logger
}

type DatabaseParam struct {
	fx.In

	Config *config.Config
	Logger logger.Logger
}

type DatabaseResult struct {
	fx.Out

	Database *Database
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

func (d *Database) Migrate() error {
	err := d.DB.WithContext(context.Background()).AutoMigrate()

	if err != nil {
		d.Logger.Error(">> Database migration failed" + err.Error())
	}

	return err
}

func (d *Database) Transaction(fc func(tx *Transaction) error) error {
	err := d.DB.Transaction(func(_tx *gorm.DB) error {
		tx := &Transaction{
			DB: _tx,
		}
		return fc(tx)
	})

	return err
}
