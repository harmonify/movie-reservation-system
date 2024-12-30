package database

import (
	"gorm.io/gorm"
)

type Transaction struct {
	DB *gorm.DB
}
