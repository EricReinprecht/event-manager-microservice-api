package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/reinp/event-platform/backend/internal/config"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: gormLogger.New(
				log.New(
					os.Stdout,
					"\r\n",
					log.LstdFlags,
				),
				gormLogger.Config{
					SlowThreshold:             time.Second,
					LogLevel:                  gormLogger.Error,
					IgnoreRecordNotFoundError: true,
					Colorful:                  true,
				},
			),
		},
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
