package database

import (
	"fmt"
	"log"
	"os"
	"time"

	config "github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewPostgresqlDatabase(p DatabaseParam, cfg *DatabaseConfig) (*Database, error) {
	master := fmt.Sprintf(
		`host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta`,
		cfg.DbHost,
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbName,
		cfg.DbPort,
	)

	client, err := gorm.Open(pgDriver.Open(master), &gorm.Config{
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
	if cfg.Env == config.EnvironmentProduction {
		client.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	if err != nil {
		p.Logger.Error(">> Database connection error: " + err.Error())
		return nil, err
	}

	if err := client.Use(otelgorm.NewPlugin()); err != nil {
		p.Logger.Warn("otelgorm.NewPlugin() error: " + err.Error())
	}

	p.Logger.Info(">> Database connected to " + cfg.DbHost)

	sqlDb, err := client.DB()
	if err != nil {
		p.Logger.Error(">> Database connection error: " + err.Error())
		return nil, err
	}

	sqlDb.SetMaxIdleConns(cfg.DbMaxIdleConn)
	sqlDb.SetMaxOpenConns(cfg.DbMaxOpenConn)
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.DbMaxLifetimeInMinute) * time.Minute)

	return &Database{
		DB:     client,
		Logger: p.Logger,
	}, nil
}
