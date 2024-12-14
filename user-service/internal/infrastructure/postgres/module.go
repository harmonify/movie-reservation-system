package app

import (
	"context"
	"fmt"
	"time"

	config "github.com/harmonify/movie-reservation-system/user-service/cmd/config"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/fx"
	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var PostgresModule = fx.Module("pgdb", fx.Provide(NewPostgresDB))

func NewPostgresDB(cfg *config.Config, logger logger_shared.Logger) (*gorm.DB, error) {
	master := fmt.Sprintf(
		`host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta`,
		cfg.DbHostMaster,
		cfg.DbUserMaster,
		cfg.DbPasswordMaster,
		cfg.DbName,
		cfg.DbPortMaster,
	)

	db, err := gorm.Open(pgDriver.Open(master), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if cfg.Env == constants.EnvironmentProduction {
		db.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		logger.Warn("otelgorm.NewPlugin() error: " + err.Error())
	}

	if cfg.DbMigration {
		logger.Info(">> Database Migration Started")
		DbMigration(db, logger)
	}

	if err != nil {
		logger.Error(">> MySQL Connection error master: " + err.Error())
		return nil, err
	}

	logger.Info(">> MySQL Connected to master " + cfg.DbHostMaster)

	if cfg.DbUseReplica {
		slave := fmt.Sprintf(
			`host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta`,
			cfg.DbHostReplica,
			cfg.DbUserReplica,
			cfg.DbPasswordReplica,
			cfg.DbName,
			cfg.DbPortReplica,
		)

		replica := dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{pgDriver.Open(slave)},
			Policy:   dbresolver.RandomPolicy{},
		})

		replica.SetMaxOpenConns(cfg.DbMaxOpenConn)
		replica.SetMaxIdleConns(cfg.DbMaxIdleConn)
		replica.SetConnMaxLifetime(time.Duration(cfg.DbMaxLifetimeInMinute) * time.Minute)

		if err = db.Use(replica); err != nil {
			logger.Error(">> MySQL Connection replica error: " + err.Error())
			return db, err
		}
		logger.Info(">> MySQL Connected to replica " + cfg.DbHostReplica)
	}

	sqlDb, err := db.DB()
	if err != nil {
		logger.Error(">> MySQL Connection error master: " + err.Error())
		return db, err
	}

	sqlDb.SetMaxIdleConns(cfg.DbMaxIdleConn)
	sqlDb.SetMaxOpenConns(cfg.DbMaxOpenConn)
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.DbMaxLifetimeInMinute) * time.Minute)

	return db, nil
}

func DbMigration(db *gorm.DB, logger logger_shared.Logger) {
	err := db.WithContext(context.Background()).AutoMigrate()

	if err != nil {
		logger.Error(">> Database migration failed" + err.Error())
	}
}
