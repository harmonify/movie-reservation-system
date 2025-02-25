package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	config_pkg "github.com/harmonify/movie-reservation-system/pkg/config"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func NewMysqlDatabase(p DatabaseParam, cfg *DatabaseConfig) (*Database, error) {
	if err := validator.New(validator.WithRequiredStructEnabled()).Struct(cfg); err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
	)

	client, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn, // data source name
		DefaultStringSize: 256, // default size for string fields
	}), &gorm.Config{
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
	if cfg.Env == config_pkg.EnvironmentProduction {
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
