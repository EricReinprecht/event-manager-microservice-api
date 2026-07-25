package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/reinp/event-platform/backend/internal/models"
)

func TestDatabase() (*gorm.DB, error) {
	return testDatabase(false)
}

func TestDatabaseSilent() (*gorm.DB, error) {
	return testDatabase(true)
}

func testDatabase(silent bool) (*gorm.DB, error) {

	err := LoadTestEnv()

	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	logLevel := logger.Error

	if silent {
		logLevel = logger.Silent
	}

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.New(
				log.New(
					os.Stdout,
					"\r\n",
					log.LstdFlags,
				),
				logger.Config{
					SlowThreshold:             time.Second,
					LogLevel:                  logLevel,
					IgnoreRecordNotFoundError: true,
					Colorful:                  true,
				},
			),
		},
	)

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Party{},
		&models.PartyMember{},
		&models.TicketCategory{},
		&models.TicketAccessWindow{},
		&models.Ticket{},
		&models.TicketScan{},
		&models.Purchase{},
		&models.PurchaseItem{},
		&models.PaymentEvent{},
		&models.PartyMemberRole{},
	)

	if err != nil {
		return nil, err
	}

	err = db.Exec(`
		DROP INDEX IF EXISTS idx_ticket_window_active;

		CREATE UNIQUE INDEX IF NOT EXISTS idx_ticket_window_pending
		ON ticket_scans(ticket_id, ticket_access_window_id)
		WHERE status = 'PENDING';
	`).Error

	if err != nil {
		return nil, err
	}

	return db, nil
}
