package database

import (
	"context"

	"gorm.io/gorm"
)

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
