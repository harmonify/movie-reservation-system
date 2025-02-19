package database

import (
	"gorm.io/gorm"
)

type Transaction struct {
	DB *gorm.DB
}

func NewTransaction(db *gorm.DB) *Transaction {
	return &Transaction{DB: db}
}
