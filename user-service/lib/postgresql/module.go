package app

import (
	"context"
	"fmt"
	"time"

	"github.com/harmonify/movie-reservation-system/user-service/lib/config"
	config_constant "github.com/harmonify/movie-reservation-system/user-service/lib/config/constant"
	"github.com/harmonify/movie-reservation-system/user-service/lib/logger"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/fx"
	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var PostgresqlModule = fx.Module(
	"postgresql",
	fx.Provide(NewDatabase),
)

type DatabaseParam struct {
	fx.In

	cfg    *config.Config
	logger logger.Logger
}

type DatabaseResult struct {
	fx.Out

	Database *Database
}

type Database struct {
	*gorm.DB
	logger logger.Logger
}

func NewDatabase(p DatabaseParam) (DatabaseResult, error) {
	master := fmt.Sprintf(
		`host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta`,
		p.cfg.DbHostMaster,
		p.cfg.DbUserMaster,
		p.cfg.DbPasswordMaster,
		p.cfg.DbName,
		p.cfg.DbPortMaster,
	)

	db, err := gorm.Open(pgDriver.Open(master), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if p.cfg.Env == config_constant.EnvironmentProduction {
		db.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	result := DatabaseResult{
		Database: &Database{
			DB: db,
		},
	}

	if err != nil {
		p.logger.Error(">> Database connection error: " + err.Error())
		return result, err
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		p.logger.Warn("otelgorm.NewPlugin() error: " + err.Error())
	}

	p.logger.Info(">> Database connected to " + p.cfg.DbHostMaster)

	if p.cfg.DbMigration {
		p.logger.Info(">> Database migration started")
		result.Database.Migrate()
	}

	sqlDb, err := db.DB()
	if err != nil {
		p.logger.Error(">> Database connection error: " + err.Error())
		return result, err
	}

	sqlDb.SetMaxIdleConns(p.cfg.DbMaxIdleConn)
	sqlDb.SetMaxOpenConns(p.cfg.DbMaxOpenConn)
	sqlDb.SetConnMaxLifetime(time.Duration(p.cfg.DbMaxLifetimeInMinute) * time.Minute)

	return result, nil
}

func (d *Database) Migrate() {
	err := d.DB.WithContext(context.Background()).AutoMigrate()

	if err != nil {
		d.logger.Error(">> Database migration failed" + err.Error())
	}
}
