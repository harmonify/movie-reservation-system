package database

import (
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

type Transaction struct {
	DB *gorm.DB
}
