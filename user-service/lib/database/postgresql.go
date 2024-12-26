package database

import (
	"fmt"
	"log"
	"os"
	"time"

	config_constant "github.com/harmonify/movie-reservation-system/user-service/lib/config/constant"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func newPostgresqlDatabase(p DatabaseParam) (DatabaseResult, error) {
	master := fmt.Sprintf(
		`host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta`,
		p.Config.DbHost,
		p.Config.DbUser,
		p.Config.DbPassword,
		p.Config.DbName,
		p.Config.DbPort,
	)

	db, err := gorm.Open(pgDriver.Open(master), &gorm.Config{
		Logger: gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			gormLogger.Config{
				SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
				LogLevel:                  gormLogger.Error,       // Log level
				IgnoreRecordNotFoundError: true,                   // Ignore gorm.ErrRecordNotFound error for logger
				ParameterizedQueries:      false,                  // False, meaning include params in the SQL log
				Colorful:                  true,                   // Colorful for non-production environment is fine
			},
		),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if p.Config.Env == config_constant.EnvironmentProduction {
		db.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	result := DatabaseResult{
		Database: &Database{
			DB:     db,
			Logger: p.Logger,
		},
	}

	if err != nil {
		p.Logger.Error(">> Database connection error: " + err.Error())
		return result, err
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		p.Logger.Warn("otelgorm.NewPlugin() error: " + err.Error())
	}

	p.Logger.Info(">> Database connected to " + p.Config.DbHost)

	if p.Config.DbMigration {
		p.Logger.Info(">> Database migration started")
		result.Database.Migrate()
	}

	sqlDb, err := db.DB()
	if err != nil {
		p.Logger.Error(">> Database connection error: " + err.Error())
		return result, err
	}

	sqlDb.SetMaxIdleConns(p.Config.DbMaxIdleConn)
	sqlDb.SetMaxOpenConns(p.Config.DbMaxOpenConn)
	sqlDb.SetConnMaxLifetime(time.Duration(p.Config.DbMaxLifetimeInMinute) * time.Minute)

	return result, nil
}
