package database

import (
	"context"

	"github.com/harmonify/movie-reservation-system/user-service/internal/driven/database/postgresql/model"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"gorm.io/gorm"
)

type Database struct {
	DB     *gorm.DB
	Logger logger.Logger
}

func (d *Database) Migrate() error {
	err := d.DB.WithContext(context.Background()).AutoMigrate(model.User{}, model.UserKey{}, model.UserSession{})

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
